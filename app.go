package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", helloworld)
	http.ListenAndServe(":8080", nil)

}

func helloworld(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "hello world")
}
