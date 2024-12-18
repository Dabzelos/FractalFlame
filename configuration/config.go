package configuration

import (
	"encoding/json"
	"os"

	"FractalFlame/internal/domain/errors"
)

type LinearTransformationsConfig struct {
	Spherical    bool `json:"Spherical"`
	Sinusoidal   bool `json:"Sinusoidal"`
	Handkerchief bool `json:"Handkerchief"`
	Swirl        bool `json:"Swirl"`
	Horseshoe    bool `json:"Horseshoe"`
	Polar        bool `json:"Polar"`
	Disc         bool `json:"Disc"`
	Heart        bool `json:"Heart"`
	Linear       bool `json:"Linear"`
	EyeFish      bool `json:"EyeFish"`
}

type Configuration struct {
	Application struct {
		Width              int     `json:"width"`
		Height             int     `json:"height"`
		StartingPoints     int     `json:"startingPoints"`
		Iterations         int     `json:"iterations"`
		SingleThread       bool    `json:"singleThread"`
		Gamma              bool    `json:"gamma"`
		GammaCoeff         float64 `json:"gammaCoeff"`
		NumWorkers         int     `json:"numWorkers"`
		HorizontalSymmetry bool    `json:"horizontalSymmetry"`
		VerticalSymmetry   bool    `json:"verticalSymmetry"`
		Format             string  `json:"format"`
	} `json:"Application"`
	ListOfTransformations LinearTransformationsConfig `json:"LinearTransformations"`
}

func Read(filePath string) (*Configuration, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Configuration

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	if config.Application.Height == 0 || config.Application.Width == 0 || config.Application.StartingPoints == 0 {
		return nil, errors.ErrZeroSizeMatrix{}
	}

	return &config, nil
}
