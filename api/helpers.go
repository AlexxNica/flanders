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
