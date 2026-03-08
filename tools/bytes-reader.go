package tools

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type BytesReader struct {
	reader *bytes.Reader
	position int
}

func NewBytesReader(data []byte) BytesReader {
	return BytesReader{
		reader: bytes.NewReader(data),
	}
}

func ToHex(data []byte) []string {
	hex := make([]string, len(data))
	for i, b := range data {
		hex[i] = fmt.Sprintf("%02x", b)
	}
	return hex
}

func (r *BytesReader) Seek(n int64, whence int) (int64, error) {
	return r.reader.Seek(n, whence)
}
func (r *BytesReader) ReadUint8() (byte, error) {
	b, err := r.reader.ReadByte()
	if err != nil {
		r.position += 1
	}
	return b, err
}
func (r *BytesReader) ReadUint16BE() (uint16, error) {
	var buf [2]byte
	_, err := r.reader.Read(buf[:])
	if err != nil {
		return 0, err
	}
	value := binary.BigEndian.Uint16(buf[:])
	r.position += 2
	return value, nil
}
func (r *BytesReader) ReadString(n int) (string, error) {
	data, err := r.Read(n)
	return string(data), err
}
func (r *BytesReader) ReadStringNT() (string, error) {
	var buf []byte
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return "", err
		}
		if b == 0 {
			break
		}
		buf = append(buf, b)
	}
	return string(buf), nil
}
func (r *BytesReader) Read(n int) ([]byte, error) {
	buf := make([]byte, n)
	size, err := r.reader.Read(buf)
	return buf[:size], err
}
func (r *BytesReader) ReadStringHex(n int) ([]string, error) {
	data, err := r.Read(n)
	if err != nil {
		return nil, err
	}
	hex := make([]string, len(data))
	for i, b := range data {
		hex[i] = fmt.Sprintf("%02x", b)
	}
	return hex, nil
}
func (r *BytesReader) ReadAll() ([]byte, error) {
	buf := make([]byte, r.reader.Len())
	size, err := r.reader.Read(buf)
	return buf[:size], err
}