package savers

import (
	"image"
	"image/png"
	"os"
)

type PngSaver struct{}

func (p *PngSaver) Save(img image.Image) error {
	file, err := os.Create("FractalFlame.png")
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}
