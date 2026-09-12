package datauri

import (
	"encoding/base64"
	"fmt"
	"mime"
	"net/url"
	"strings"

	"github.com/kellegous/poop"
)

func getMimeHeader(meta string) (string, error) {
	parts := strings.Split(meta, ";")
	isBase64 := false
	withoutBase64 := make([]string, 0, len(parts))

	for _, part := range parts {
		if strings.EqualFold(part, "base64") {
			isBase64 = true
			continue
		}
		withoutBase64 = append(withoutBase64, part)
	}

	if !isBase64 {
		return "", poop.New("data URI is not base64")
	}

	header := strings.Join(withoutBase64, ";")
	if header == "" {
		header = "text/plain;charset=US-ASCII"
	}
	return header, nil
}

func Parse(s string) (*URI, error) {
	meta, encodedData, ok := strings.Cut(s, ",")
	if !ok {
		return nil, poop.New("data URI is missing a comma")
	}

	// Decode percent-encoded metadata, but preserve '+' characters.
	meta, err := url.PathUnescape(meta)
	if err != nil {
		return nil, poop.Chain(err)
	}

	header, err := getMimeHeader(meta)
	if err != nil {
		return nil, poop.Chain(err)
	}

	mediaType, params, err := mime.ParseMediaType(header)
	if err != nil {
		return nil, fmt.Errorf("parse metadata: %w", err)
	}

	// PathUnescape is intentional: '+' is valid data and must not become a space.
	decodedData, err := url.PathUnescape(encodedData)
	if err != nil {
		return nil, poop.Chain(err)
	}

	data, err := base64.StdEncoding.DecodeString(decodedData)
	if err != nil {
		return nil, poop.Chain(err)
	}

	return &URI{
		MediaType: mediaType,
		Params:    params,
		Data:      data,
	}, nil
}
