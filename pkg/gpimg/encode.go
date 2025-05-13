package gpimg

import (
	"encoding/binary"
	"image"
	"image/draw"
	"io"
	"math"
	"slices"
)

// In this case
func Encode(in io.Reader, out io.WriteSeeker, flags byte) error {
	// Load image

	imgFile, _, err := image.Decode(in)
	if err != nil {
		return err
	}

	img := image.NewRGBA(imgFile.Bounds())
	draw.Draw(img, imgFile.Bounds(), imgFile, imgFile.Bounds().Min, draw.Src)

	if img.Bounds().Dx() > 320 || img.Bounds().Dy() > 240 {
		return ErrImgTooLarge
	}

	// Write width and height as 4 bytes

	var size [4]byte = [4]byte{}

	binary.LittleEndian.PutUint16(
		size[:2],
		uint16(img.Bounds().Dx()),
	)

	binary.LittleEndian.PutUint16(
		size[2:],
		uint16(img.Bounds().Dy()),
	)

	out.Seek(0, io.SeekEnd)
	_, err = out.Write(size[:])
	if err != nil {
		return err
	}

	out.Seek(0, io.SeekEnd)
	_, err = out.Write([]byte{flags}[:])
	if err != nil {
		return err
	}

	// Write rest of image depending on flags

	if flags&FlagNoCompression != 0 {
		err = writeNoCompression(img, out, flags)
		if err != nil {
			return err
		}
	} else if flags&FlagRLECompression != 0 {
		err = writeRLE(img, out, flags)
		if err != nil {
			return err
		}
	} else {
		return ErrUnrecognizedFormat
	}

	return nil
}

func writeNoCompression(img *image.RGBA, out io.WriteSeeker, flags byte) error {

	out.Seek(0, io.SeekEnd)
	_, err := out.Write(img.Pix)
	if err != nil {
		return err
	}

	return nil

}

func writeRLE(img *image.RGBA, out io.WriteSeeker, flags byte) error {

	var comparisonPixel []byte = img.Pix[:4]

	var imgIndex int = 4
	var pixelSectionLength byte = 1

	for imgIndex < len(img.Pix) {

		var slicesEqual int = slices.Compare(comparisonPixel, img.Pix[imgIndex:imgIndex+4])

		if slicesEqual == 0 {
			pixelSectionLength++
		} else if slicesEqual != 0 || pixelSectionLength == math.MaxUint8 {

			out.Seek(0, io.SeekEnd)
			out.Write(
				[]byte{pixelSectionLength},
			)

			out.Seek(0, io.SeekEnd)
			out.Write(
				comparisonPixel,
			)

			pixelSectionLength = 1
			comparisonPixel = img.Pix[imgIndex : imgIndex+4]

		}

		imgIndex += 4

	}

	// Write the last set of pixels

	out.Seek(0, io.SeekEnd)
	out.Write(
		[]byte{pixelSectionLength},
	)

	out.Seek(0, io.SeekEnd)
	out.Write(
		comparisonPixel,
	)

	return nil

}
