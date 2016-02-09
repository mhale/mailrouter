package main

import (
	"io/ioutil"
	"os"
	"strconv"
)

func CreatePIDFile() error {
	pid := os.Getpid()
	data := []byte(strconv.Itoa(pid))
	err := ioutil.WriteFile(config.Options["PIDFile"], data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func RemovePIDFile() error {
	err := os.Remove(config.Options["PIDFile"])
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
