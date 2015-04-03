package contextio

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/garyburd/go-oauth/oauth"
)

type ContextIO struct {
	key    string
	secret string
	client *oauth.Client
}

func NewContextIO(key, secret string) *ContextIO {
	c := &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  key,
			Secret: secret,
		},
	}

	return &ContextIO{
		key:    key,
		secret: secret,
		client: c,
	}
}

const (
	apiHost = `api.context.io`
)

// Do signs the request and returns an *http.Response, the body must be defer response.Body.close()
func (c *ContextIO) Do(method, q string, params url.Values, body io.Reader) (response *http.Response, err error) {
	// Cannot use http.NewRequest because of the possibility of encoded data in the url
	req := &http.Request{
		Method: method,
		Host:   apiHost, // takes precendence over Request.URL.Host
		URL: &url.URL{
			Host:     apiHost,
			Scheme:   "https",
			Opaque:   q,
			RawQuery: params.Encode(),
		},
		Header: http.Header{
			"User-Agent": {"GoContextIO Simple Library"},
		},
	}

	fmt.Print("req.URL ")
	fmt.Println(req.URL)
	fmt.Print("req.URL.Opaque ")
	fmt.Println(req.URL.Opaque)
	if err != nil {
		return
	}
	err = c.client.SetAuthorizationHeader(req.Header, nil, req.Method, req.URL, nil)
	fmt.Println("HL:", req.Header)
	if err != nil {
		return
	}
	return http.DefaultClient.Do(req)
}

// DoJson passes the request to Do and then returns the json in a []byte array
func (c *ContextIO) DoJson(method, u string, params url.Values, body io.Reader) (j []byte, err error) {
	response, err := c.Do(method, u, params, body)
	defer response.Body.Close()
	j, err = ioutil.ReadAll(response.Body)
	return j, err
}
