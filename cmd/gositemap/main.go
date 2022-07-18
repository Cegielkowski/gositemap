package main

import (
	"gositemap/pkg/gositemap"
	"os"
)

func main() {
	os.Exit(gositemap.CLI(os.Args[1:]))
}
