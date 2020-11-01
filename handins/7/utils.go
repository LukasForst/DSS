package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
)

func encodeToBase64(v interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	defer encoder.Close()
	err := json.NewEncoder(encoder).Encode(v)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func decodeFromBase64(v interface{}, enc string) error {
	return json.NewDecoder(base64.NewDecoder(base64.StdEncoding, strings.NewReader(enc))).Decode(v)
}
