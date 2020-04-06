package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	var port int
	var filename string
	var cors bool
	flag.IntVar(&port, "port", 8000, "The port to listen on (default 8000)")
	flag.StringVar(&filename, "file", "", "Optional. A file to echo")
	flag.BoolVar(&cors, "cors", false, "Whether to respond to CORS headers")
	flag.Parse()

	var filecontent string
	var err error
	if filename != "" {
		filecontent, err = readFile(filename)
		if err != nil {
			fmt.Printf("\033[31merr %s\n\033[m", err)
			os.Exit(1)
		}
	}

	s := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cors {
				origins := r.Header["Origin"]
				var origin string
				// try to avoid setting origin to * directly bc
				// chrome complains about setting origin to *
				// if request.credentials has been set to
				// 'include'
				// https://developer.mozilla.org/en-US/docs/Web/API/Request/credentials
				if len(origins) > 0 {
					origin = origins[0]
				} else {
					origin = "*"
				}
				w.Header().Set("Access-Control-Allow-Origin", origin)
				// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
				w.Header().Set("Access-Control-Allow-Credentials", "true")

				reqHeaders := r.Header["Access-Control-Request-Headers"]
				if len(reqHeaders) > 0 {
					w.Header().Set(
						"Access-Control-Allow-Headers",
						strings.Join(reqHeaders, ","))
				}
			}
			var s string
			if filename != "" {
				s = filecontent

			} else {
				bytes, err := ioutil.ReadAll(r.Body)
				s = string(bytes)
				if err != nil {
					fmt.Printf("\033[31merr %s\n\033[m", err)
					return
				}
			}
			fmt.Fprint(w, string(s))
		}),
	}
	fmt.Println(s.ListenAndServe())
}

func readFile(filename string) (string, error) {
	filename = strings.Trim(filename, " ")
	var bytes []byte
	var err error
	if filename != "-" {
		bytes, err = ioutil.ReadFile(filename)
		if err != nil {
			return "", err
		}
	} else {
		bytes, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
	}
	return string(bytes), nil
}
