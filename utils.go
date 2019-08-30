package main

import (
	"fmt"
	"net/url"
)

const hostName string = "wikimon.hatnote.com"

var langCodePorts = map[string]int{
	"en": 9000,
	"de": 9010,
}

func MakeUrlForLangCode(langCode string) (string, error) {
	port, exists := langCodePorts[langCode]
	if !exists {
		return "", fmt.Errorf("Lang code '%s' not in port map!", langCode)
	}
	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", hostName, port),
		Path:   "/",
	}
	return u.String(), nil
}
