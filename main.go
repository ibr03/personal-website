package main

import (
	"fmt"
	"simpleSSG/builder"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: build | serve")
		return
	}

	switch os.Args[1] {
	case "build":
		err := builder.Build()
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println("site built!")

	case "serve":
		fmt.Println("serving at http://localhost:8080")
		http.ListenAndServe(":8080", http.FileServer(http.Dir("./build")))
	}
}