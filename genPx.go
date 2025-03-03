package main

import (
	"bytes"
	"image"
	"image/png"
)

func generatePixel() ([]byte, error) {
	img := image.NewAlpha(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, image.Transparent)
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
