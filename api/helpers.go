package api

import (
	"net/url"
	"strings"
)

func sanatizedForm(f url.Values) url.Values {
	var sanatized = make(url.Values)
	for k, values := range f {
		k = strings.ToLower(k)
		sanatized[k] = values
	}
	return sanatized
}

//UniqueSlice returns a Unique, Compact, Trimmed sice of original
func UniqueSlice(in []string) []string {
	uniques := make(map[string]struct{})
	// put all supplied strings into map to force unique
	for _, s := range in {
		s = strings.Trim(s, " ")
		if s != "" {
			uniques[s] = struct{}{}
		}
	}

	// convert map keys back to slice
	out := make([]string, 0, len(uniques))
	for key := range uniques {
		out = append(out, key)
	}
	return out
}
