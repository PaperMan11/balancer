package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":1331", "input http listen addr")
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println(r.Header)
		fmt.Fprintf(w, "%s: hello", *addr)
	})
	http.ListenAndServe(*addr, nil)
}
