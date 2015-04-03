/*
Package GoContextIO provides a simple way to sign API requests for http://Context.IO.


*/
package contextio

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/garyburd/go-oauth/oauth"
)

// ContextIO is a struct containing the authentication information and a pointer to the oauth client
type ContextIO struct {
	key    string
	secret string
	client *oauth.Client
}

// NewContextIO returns a ContextIO struct based on your CIO User and Secret
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

// Do signs the request and returns an *http.Response. The body is a standard response.
// Body and must have defer response.Body.close().
// This is 2 legged authentication, and will not currently work with 3 legged authentication.
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

	err = c.client.SetAuthorizationHeader(req.Header, nil, req.Method, req.URL, nil)
	if err != nil {
		return
	}
	return http.DefaultClient.Do(req)
}

// DoJSON passes the request to Do and then returns the json in a []byte array
func (c *ContextIO) DoJSON(method, q string, params url.Values, body io.Reader) (json []byte, err error) {
	response, err := c.Do(method, q, params, body)
	defer response.Body.Close()
	json, err = ioutil.ReadAll(response.Body)
	return json, err
}
