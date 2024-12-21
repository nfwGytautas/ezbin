package shared

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

// Check if a file exists
func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Get the home directory
func HomeDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return home, nil
}

// Read a json file and parse it
func ReadJson(path string, data interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return err
	}

	return nil
}

// Marshal data to json and write it to a file
func WriteJson(path string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
