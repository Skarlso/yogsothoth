package yogsot

import (
	"errors"

	"gopkg.in/yaml.v2"
)

// Parameter are the variables of the stack.
type Parameter struct {
	Default     string `yaml:"Default"`
	Type        string `yaml:"Type"`
	Description string `yaml:"Description"`
}

type createStackInput struct {
	Parameters map[string]Parameter              `yaml:"Parameters"`
	Resources  map[string]map[string]interface{} `yaml:"Resources"`
}

func parseTemplate(template []byte) (createStackInput, error) {
	csi := createStackInput{}
	err := yaml.Unmarshal(template, &csi)
	if err != nil {
		return csi, errors.New("error happened while unmarshaling template: " + err.Error())
	}

	return csi, nil
}
