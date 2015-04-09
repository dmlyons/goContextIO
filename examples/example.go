package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	contextio "github.com/dmlyons/goContextIO"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Println(os.Args[0], "-key=\"CioKey\" -secret=\"CioSecret\" [method] [endpoint] [query string]")
		fmt.Println("Make sure that your query string keys and values are properly escaped, such as with url.QueryEscape,")
		fmt.Println("as well as any values that are in the actual query")
		flag.PrintDefaults()
	}
	key := flag.String("key", "", "Your CIO User Key")
	secret := flag.String("secret", "", "Your CIO User Secret")
	body := flag.String("body", "", "The body of the request, ignored if method is not a POST/PUT")
	flag.Parse()
	c := contextio.NewContextIO(*key, *secret)
	if len(flag.Args()) < 2 {
		fmt.Println("Must provide at least 2 arguments")
		flag.Usage()
		os.Exit(1)
	}
	m := flag.Arg(0)
	q := flag.Arg(1)
	p := flag.Arg(2)
	params, err := url.ParseQuery(p)
	if err != nil {
		fmt.Printf("Unable to parse query string %s: %v\n", q, err)
	}
	j, err := c.DoJSON(m, q, params, strings.NewReader(*body))
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	var out bytes.Buffer
	err = json.Indent(&out, j, "", "  ")
	if err != nil {
		fmt.Println("JSON ERROR:", err)
	}
	fmt.Println(out.String())
}
