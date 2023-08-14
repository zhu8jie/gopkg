package yaml

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Parse(configPath string, dst interface{}) error {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errors.New(fmt.Sprintf("yaml file get err: %v ", err))
	}
	if err = yaml.Unmarshal(yamlFile, &dst); err != nil {
		return errors.New(fmt.Sprintf("yaml unmarshal: %v", err))
	}

	return nil
}
