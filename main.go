package main

import (
	"fmt"
	"strings"

	"github.com/jdkato/prose/v2"
	"github.com/ororsatti/go-hoblot/radix"
)

func main() {
	text := `
	The most saddening thought that arises after the perusal of this Volume,
	is, that no change has yet been made in the infamous Lunacy Laws, for
	which, in the main, we have to thank our Whig Rulers. Never was a more
	criminal or despotic Law passed than that which now enables a Husband to
	lock up his Wife in a Madhouse on the certificate of two medical men,
	who often in haste, frequently for a bribe, certify to madness where
	none exists. We believe that under these Statutes thousands of persons,
	perfectly sane, are now imprisoned in private asylums throughout the
	Kingdom; while strangers are in possession of their property; and the
	miserable prisoner is finally brought to a state of actual lunacy or
	imbecility—however rational he may have been when first immured. The
	Keepers of these Madhouse Dens, from long study in their diabolical art,
	can reduce, by certain drugs, the clearest brain to a state of stupor;
	and the Lunacy Commissioners take all for granted that they hear over the
	luxurious lunch with which the Mad Doctor regales them.
	`

	doc, _ := prose.NewDocument(text)

	type Data struct {
		Count int
	}

	smap := radix.NewSearchableMap[Data]()
	// Tokenize the text
	for _, tok := range doc.Tokens() {
		norm := strings.ToLower(tok.Text)
		prev := smap.Get(norm)

		var prevCount int
		if prev != nil {
			prevCount = prev.Data.Count
		}

		smap.Set(norm, &Data{
			Count: prevCount + 1,
		})
	}

	res := smap.FuzzyGet("king", 3)

	for key, val := range res {
		fmt.Println(key, val.Data, val.Distance)
	}

	// smap.Print()
}
