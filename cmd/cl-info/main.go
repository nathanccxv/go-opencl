package main

import (
	"github.com/kr/pretty"
	cl "github.com/nathanccxv/go-opencl"
)

func main() {
	info, _ := cl.Info()
	pretty.Println(info)
}
