/*
Package contextio provides a simple way to sign API requests for http://Context.IO.

The simplest usage is to use DoJSON() to return a json byte array that you can use elsewhere in your code.
For more advanced usage, you can use Do() and parse through the http.Response struct yourself. It is not
specific to an API version, so you can use it to make any request you would make through http://console.Context.IO.
*/
package contextio

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/garyburd/go-oauth/oauth"
)

const (
	defaultMaxMemory = 32 << 21 // 64 MB
)

// the default host that the library contacts
const defaultAPIHost = "api.context.io"

// ContextIO is a struct containing the authentication information and a pointer to the oauth client
type ContextIO struct {
	key     string
	secret  string
	client  *oauth.Client
	apiHost string
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
		key:     key,
		secret:  secret,
		client:  c,
		apiHost: defaultAPIHost,
	}
}

// SetAPIHost sets the domain (i.e. "api.context.io) for the requests, useful if you are mocking the API for testing
func (c *ContextIO) SetAPIHost(h string) *ContextIO {
	c.apiHost = h
	return c
}

// NewRequest generates a request and signs it
func (c *ContextIO) NewRequest(method, q string, queryParams url.Values, body *string) (req *http.Request, err error) {
	// make sure q has a slash in front of it
	if q[0:1] != "/" {
		q = "/" + q
	}

	query := c.apiHost + q
	if len(queryParams) > 0 {
		query = query + "?" + queryParams.Encode()
	}
	req, err = http.NewRequest(method, "https://"+query, strings.NewReader(*body))
	if err != nil {
		return nil, err
	}
	req.URL.Opaque = q
	req.Header.Set("User-Agent", "GoContextIO Simple Library v. 0.1")

	v := url.Values{}
	switch method {
	case "PUT", "POST", "DELETE":
		// need form data here if uploading
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		v, err = url.ParseQuery(*body)
		if err != nil {
			return nil, err
		}
	}

	err = c.client.SetAuthorizationHeader(req.Header, nil, req.Method, req.URL, v)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// AttachFile will create a file upload in the request, assumes NewRequest has already been called
func (c *ContextIO) AttachFile(req *http.Request, fieldName, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filepath.Base(fileName))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, f)
	if err != nil {
		return err
	}

	// transfer the existing post vals into the new body
	for key, valSlice := range req.PostForm {
		for _, val := range valSlice {
			err = writer.WriteField(key, val)
			if err != nil {
				return err
			}
		}
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	rc := ioutil.NopCloser(body)
	req.Body = rc
	// update the form
	err = req.ParseMultipartForm(defaultMaxMemory)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return nil
}

// Do signs the request and returns an *http.Response. The body is a standard response.Body
// and must have defer response.Body.close().  Does not support uploads, use NewRequest and AttachFile for that.
// This is 2 legged authentication, and will not currently work with 3 legged authentication.
func (c *ContextIO) Do(method, q string, params url.Values, body *string) (response *http.Response, err error) {
	req, err := c.NewRequest(method, q, params, body)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

// DoJSON passes the request to Do and then returns the json in a []byte array
func (c *ContextIO) DoJSON(method, q string, params url.Values, body *string) (json []byte, err error) {
	response, err := c.Do(method, q, params, body)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()
	json, err = ioutil.ReadAll(response.Body)
	return json, err
}
