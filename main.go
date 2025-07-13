package main

import (
	"os"

	"github.com/ororsatti/go-searchdex/radix"
)

func main() {
	smap := radix.NewSearchableMap()
	smap.Set("", true)

	smap.Print(os.Stdout)
}
