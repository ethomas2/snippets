package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8000, "The port to listen on (default 8000)")
	flag.Parse()
	s := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("\033[31merr %s\n\033[m", err)
			}
			fmt.Fprintln(w, string(s))
		}),
	}
	fmt.Println(s.ListenAndServe())
}
