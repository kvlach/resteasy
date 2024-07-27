package resteasy

import (
	"fmt"
	"io"
	"log/slog"
	"math"
	"time"

	"net/http"
)

type Request struct {
	url             string
	method          string
	retriesMax      int
	retriesAttmpted int
	token           string
	json            bool
	//data
	//form
	query map[string]string
}

func NewRequest(url string) *Request {
	if url == "" {
		panic("url is empty")
	}
	return &Request{url: url}
}

func GET(url string) *Request {
	return NewRequest(url).Method("GET")
}

func HEAD(url string) *Request {
	return NewRequest(url).Method("HEAD")
}

func POST(url string) *Request {
	return NewRequest(url).Method("POST")
}

func PUT(url string) *Request {
	return NewRequest(url).Method("PUT")
}

func DELETE(url string) *Request {
	return NewRequest(url).Method("DELETE")
}

func CONNECT(url string) *Request {
	return NewRequest(url).Method("CONNECT")
}

func OPTIONS(url string) *Request {
	return NewRequest(url).Method("OPTIONS")
}

func TRACE(url string) *Request {
	return NewRequest(url).Method("TRACE")
}

func PATCH(url string) *Request {
	return NewRequest(url).Method("PATCH")
}

func (r *Request) Method(method string) *Request {
	r.method = method
	return r
}

// or possibly rename to Bearer()
func (r *Request) Token(token string) *Request {
	r.token = token
	return r
}

func (r *Request) Retry(attempts int) *Request {
	r.retriesMax = attempts
	return r
}

func (r *Request) JSON(parse bool) *Request {
	r.json = parse
	return r
}

func (r *Request) Query(pairs ...string) *Request {
	if len(pairs)%2 != 0 {
		panic("Expected pairs (even number of arguments)")
	}

	params := make(map[string]string)
	for i := 0; i < len(pairs); i += 2 {
		params[pairs[i]] = pairs[i+1]
	}
	r.query = params
	return r
}

// the status code's are taken from curl's --retry flag
func transientError(statusCode int) bool {
	switch statusCode {
	case 408, 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}

// If 204, ret will not be modified
func (r *Request) Do(ret any) {
	req, err := http.NewRequest(r.method, r.url, nil)
	if err != nil {
		panic(err)
	}
	if r.token != "" {
		req.Header.Add("Authorization", "Bearer "+r.token)
	}

	query := req.URL.Query()
	for k, v := range r.query {
		query.Set(k, v)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if transientError(resp.StatusCode) && r.retriesMax > r.retriesAttmpted {
		r.retriesAttmpted++
		wait := time.Second * time.Duration(math.Pow(2, float64(r.retriesAttmpted)))
		slog.Debug(
			fmt.Sprintf("Request failed, retrying in %v (%d/%d)",
				wait, r.retriesAttmpted, r.retriesMax),
			"StatusCode", resp.StatusCode)
		time.Sleep(wait)
		r.Do(ret)
		return
	}

	if resp.StatusCode >= 400 {
		panic(fmt.Errorf("Request failed, status code: %d, body: %s", resp.StatusCode, body))
	}

	if r.json {
		if err := StrictUnmarshalJSON(body, ret); err != nil {
			panic(err)
		}
	} else {
		switch ret.(type) {
		case *string:
			*ret.(*string) = string(body)
		default:
			panic("expected string pointer")
		}
	}
}
