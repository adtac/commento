package main

import (
	"bytes"
	"compress/gzip"
)

func gzipStatic(b []byte) ([]byte, error) {
	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	if _, err := g.Write(b); err != nil {
		g.Close()
		return []byte{}, err
	}

	g.Close()
	return buf.Bytes(), nil
}
