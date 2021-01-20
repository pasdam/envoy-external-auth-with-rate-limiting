package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var resp []string

	log.Println("Request received")

	resp = append(resp, fmt.Sprintf("Hello %s!", r.Header.Get("x-user-id")))

	for name, headers := range r.Header {
		for _, h := range headers {
			resp = append(resp, fmt.Sprintf("%v: %v", name, h))
		}
	}

	w.Write([]byte(strings.Join(resp, "\n")))
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8123", nil)
}
