package gpimg

import (
	"errors"
	"fmt"
)

var ErrImgTooLarge error = fmt.Errorf("image is too large must be at most 320x240")
var ErrUnrecognizedFormat error = errors.New("unrecognized format")

var FlagNames []string = []string{"opaque", "nocompression", "rle"}

var ColourFlags map[string]byte = map[string]byte{
	"opaque": FlagBgOpaque,
}

var CompressionFlags map[string]byte = map[string]byte{
	"nocompression": FlagNoCompression,
	"rle":           FlagRLECompression,
}

var AllFlags map[string]byte

const (
	FlagNoCompression byte = 1 << iota
	FlagRLECompression
	FlagBgOpaque // Mostly used for rasterisation
)

const HeaderSize int = 5

const (
	DefaultFlags byte = FlagNoCompression | FlagBgOpaque
)

func init() {

	AllFlags = make(map[string]byte)

	for k, v := range ColourFlags {
		AllFlags[k] = v
	}

	for k, v := range CompressionFlags {
		AllFlags[k] = v
	}

}
