package meme

import (
	"io"
	"os"
)

func readFirstBytesFromFile(file string, byteCount uint8) ([512]byte, error) {
	var header [512]byte

	r, err := os.Open(file) // lint:allow_possible_insecure
	if err != nil {
		return [512]byte{}, err
	}

	defer func() {
		_ = r.Close()
	}()

	_, err = io.ReadFull(r, header[:])
	if err != nil {
		return [512]byte{}, err
	}

	return header, nil
}
