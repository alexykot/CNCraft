package network

import (
	"bytes"
	"compress/zlib"
	"io"
)

type compressor struct {
	enabled bool
}

func (c *compressor) Enable() {
	c.enabled = true
}
func (c *compressor) Disable() {
	c.enabled = false
}

func (c *compressor) Deflate(data []byte) []byte {
	if !c.enabled {
		return data
	}

	var out bytes.Buffer

	writer, _ := zlib.NewWriterLevel(&out, zlib.BestCompression)
	_, _ = writer.Write(data)
	_ = writer.Close()

	return out.Bytes()
}

func (c *compressor) Inflate(data []byte) []byte {
	if !c.enabled {
		return data
	}

	reader, _ := zlib.NewReader(bytes.NewReader(data)) // error should never happen with a plain bytes.Reader

	var out bytes.Buffer
	_, _ = io.Copy(&out, reader) // error should never happen with a plain bytes.Buffer
	_ = reader.Close()

	return out.Bytes()
}
