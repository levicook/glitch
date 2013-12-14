package main

import (
	"fmt"
	"log"
	"strings"
)

var (
	emptyStruct = struct{}{}

	hasSuffix = strings.HasSuffix
	contains  = strings.Contains

	sprintf = fmt.Sprintf
	errorf  = fmt.Errorf
	printf  = log.Printf
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func clearScrollBuffer() {
	print("\033c")
}
