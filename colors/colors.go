package colors

import (
	"io/ioutil"

	"github.com/IgaguriMK/ledscreen/pixcels"

	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

type ColorSetting struct {
	Default string `yaml:"default"`
	Back    string `yaml:"back"`
	Colors  []struct {
		Name  string    `yaml:"name"`
		Value []float32 `yaml:"value"`
	} `yaml:"colors"`
}

func LoadTable(fileName string) (map[string]pixcels.Pixcel, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot read file")
	}

	var colorSetting ColorSetting
	err = yaml.Unmarshal(bytes, &colorSetting)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot parse file")
	}

	colorMap := make(map[string]pixcels.Pixcel)
	for _, color := range colorSetting.Colors {
		c := pixcels.FromArray(color.Value)
		colorMap[color.Name] = c
	}

	d, ok := colorMap[colorSetting.Default]
	if !ok {
		return nil, errors.New("Cannot find default fill color '" + colorSetting.Default + "'")
	}

	colorMap[":"] = d

	b, ok := colorMap[colorSetting.Back]
	if !ok {
		return nil, errors.New("Cannot find default background color '" + colorSetting.Back + "'")
	}

	colorMap["_"] = b

	return colorMap, nil
}
