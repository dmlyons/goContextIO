package contextio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func ExampleNewContextIO() {
	c := NewContextIO("MyCIOUserKey", "MyCIOSecret")
	params := url.Values{}
	params.Set("limit", "10")

	// body is usually used for POSTs
	body := strings.NewReader("{some_param:1}")

	j, err := c.DoJSON("GET", "/2.0/accounts/", params, body)
	if err != nil {
		fmt.Println("DoJSON Error: %v", err)
		return
	}

	var out bytes.Buffer
	json.Indent(&out, j, "", "  ")
	fmt.Println(out.String())
}
