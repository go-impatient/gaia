package cache

import (
	"bytes"
	"encoding/gob"
)

var GobEncoding Encoding = NewGobEncoding()

func NewGobEncoding() *gobEncoding {
	return &gobEncoding{}
}

type gobEncoding struct{}

// Encode 编码器
func (e *gobEncoding) Encode(data interface{}) ([]byte, error) {
	encoded := &bytes.Buffer{}
	enc := gob.NewEncoder(encoded)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	return encoded.Bytes(), nil
}

// Decode 解码器
func (e *gobEncoding) Decode(b []byte, data interface{}) error {
	if !isPointer(data) {
		return ErrNotAPointer
	}
	dec := gob.NewDecoder(bytes.NewReader(b))
	return dec.Decode(data)
}
