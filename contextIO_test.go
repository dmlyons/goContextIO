package contextio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"
)

func ExampleNewContextIO() {
	c := NewContextIO("MyCIOUserKey", "MyCIOSecret")
	params := url.Values{}
	params.Set("limit", "10")

	// body is usually used for POSTs, ignored otherwise, but
	// its format looks like this:
	body := "some_param=1&some_other_param=2"

	j, err := c.DoJSON("GET", "/2.0/accounts", params, &body)
	if err != nil {
		fmt.Println("DoJSON Error:", err)
		return
	}

	var out bytes.Buffer
	json.Indent(&out, j, "", "  ")
	fmt.Println(out.String())
}

func TestTimer(t *testing.T) {
	c := NewContextIO("MyCIOUserKey", "MyCIOSecret")
	c.SetRate(30)
	fmt.Println(c.rate, "seconds per request")
	for i := 0; i < 3; i++ {
		before := time.Now()
		c.wait()
		fmt.Println(time.Now().Sub(before))
	}
}
