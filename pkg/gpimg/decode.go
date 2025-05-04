package gpimg

import (
	"io"
)

// Given an encoded image (including header), it returns the raw RGBA data
func Decode(in io.ReaderAt, size int, offset int) ([]byte, error) {

	var header []byte = make([]byte, 5)
	_, err := in.ReadAt(header, int64(offset))
	if err != nil {
		return nil, err
	}

	var flags byte = header[4]
	offset += HeaderSize
	size -= HeaderSize

	if flags&FlagNoCompression != 0 {
		val, err := decodeNoCompression(in, offset, size)
		if err != nil {
			return nil, err
		} else {
			return val, nil
		}
	} else if flags&FlagRLECompression != 0 {
		val, err := decodeRLE(in, offset, size)
		if err != nil {
			return nil, err
		} else {
			return val, nil
		}
	} else {
		return nil, ErrUnrecognizedFormat
	}

}

func decodeNoCompression(in io.ReaderAt, offset int, size int) (res []byte, err error) {

	res = make([]byte, size)
	_, err = in.ReadAt(res, int64(offset))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func decodeRLE(in io.ReaderAt, offset int, size int) (res []byte, err error) {

	res = make([]byte, 0)

	for i := offset; i < offset+size; i += 5 {

		var currentColour []byte = make([]byte, 5)
		_, err := in.ReadAt(currentColour, int64(i))
		if err != nil {
			return nil, err
		}

		for j := 0; j < int(currentColour[0]); j++ {

			res = append(res, currentColour[1:]...)

		}

	}

	return res, nil
}
