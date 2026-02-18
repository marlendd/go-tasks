package main

import (
	"fmt"
	"strings"
)

func slugify(src string) string {
    var b strings.Builder
    b.Grow(len(src))
    
    needsDash := false
    firstChar := true
    
    for _, r := range src {
        switch {
        case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
            if needsDash && !firstChar {
                b.WriteByte('-')
            }
            b.WriteRune(r)
            needsDash = false
            firstChar = false
        case r >= 'A' && r <= 'Z':
            if needsDash && !firstChar {
                b.WriteByte('-')
            }
            b.WriteRune(r + ('a' - 'A'))
            needsDash = false
            firstChar = false
        case r == '-':
            if needsDash && !firstChar {
                b.WriteByte('-')
            }
            b.WriteByte('-')
            needsDash = false
            firstChar = false
        default:
            if !firstChar {
                needsDash = true
            }
        }
    }
    
    return b.String()
}

func main() {
	const phrase = "A 100x Investment (2019)"
	slug := slugify(phrase)
	fmt.Println(slug)
	// a-100x-investment-2019
}
