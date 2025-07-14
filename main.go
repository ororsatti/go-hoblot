package main

import (
	"fmt"
	"os"

	"github.com/ororsatti/go-searchdex/radix"
)

func main() {
	smap := radix.NewSearchableMap()
	smap.Set("test", true)
	smap.Set("testing", true)

	smap.Print(os.Stdout)
	res := smap.FuzzyGet("test", 0)
	fmt.Println(res)
}
