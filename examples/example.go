package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

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
	verbose := flag.Bool("verbose", false, "Print extra information about the request")
	key := flag.String("key", "", "Your CIO User Key")
	secret := flag.String("secret", "", "Your CIO User Secret")
	body := flag.String("body", "", `The body of the request, ignored if method is not a POST/PUT, POST parameters should be url.QueryEscape'ed in here, not in the [query string]`)
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
	resp, err := c.Do(m, q, params, body)
	if err != nil {
		fmt.Println("Request Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	j, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	var out bytes.Buffer
	err = json.Indent(&out, j, "", "  ")
	if err != nil {
		fmt.Println("JSON INDENT ERROR:", err)
		os.Exit(1)
	}
	if *verbose {
		fmt.Println("Status:", resp.Status)
		err = resp.Header.Write(os.Stdout)
		if err != nil {
			fmt.Println("Header Write Error:", err)
		}
	}
	fmt.Println(out.String())
}
