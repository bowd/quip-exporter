package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return fmt.Errorf("could not create folder: %s", err)
		}
	}
	return nil
}

func SaveBytesToFile(filename string, content []byte) error {
	if err := EnsureDir(path.Dir(filename)); err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create/open json file: %s: %s", filename, err)
	}
	_, err = f.Write(content)
	if err != nil {
		return fmt.Errorf("could not write json file: %s", err)
	}
	return nil
}

func SaveJSONToFile(filename string, object interface{}) error {
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return SaveBytesToFile(filename, data)
}
