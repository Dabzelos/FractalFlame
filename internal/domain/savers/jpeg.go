package savers

import (
	"image"
	"image/jpeg"
	"os"
)

type JpegSaver struct {
}

func (j *JpegSaver) Save(img image.Image) error {
	file, err := os.Create("FractalFlame.jpg")
	if err != nil {
		return err
	}
	defer file.Close()

	options := &jpeg.Options{Quality: 100}
	if err := jpeg.Encode(file, img, options); err != nil {
		return err
	}
	return nil
}
