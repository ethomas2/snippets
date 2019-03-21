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
	flag.IntVar(&port, "port", 8000, "The port to listen on (default 8000)")
	flag.StringVar(&filename, "file", "", "Optional. A file to echo")
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
