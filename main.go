package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func main() {
	fmt.Println("Start server...")
	fmt.Println("Listen tcp://localhost:8080")
	fasthttp.ListenAndServe("localhost:8080", nil)
}
