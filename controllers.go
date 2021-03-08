package main

import (
	"fmt"
	"net/http"
)

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	if _, err := fmt.Fprint(w, "Hello world!"); err != nil {
		fmt.Println(err)
	}
}
