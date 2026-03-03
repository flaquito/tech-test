package image

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"

	_ "image/png"
)

func Normalize(r io.Reader) ([]byte, string, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, "", err
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return nil, "", err
	}

	return buf.Bytes(), ".jpg", nil
}
