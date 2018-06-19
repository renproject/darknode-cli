package main

import (
	"fmt"
	"os"
)

// StringInSlice checks whether the string is in the slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Mkdir tries to create folder with the given name. If not success,
// it will try appending number appended to the name.
func Mkdir(name string) (string, error) {
	// Try creating directory with the exact name we get
	if _, err := os.Stat(name); os.IsNotExist(err) {
		err = os.Mkdir(name, os.ModePerm)
		if err == nil {
			return name, nil
		}
	}

	// Try creating directory with number appended to the name.
	i := 1
	for {
		dirName := fmt.Sprintf("%v_%d", name, i)
		if _, err := os.Stat(dirName); os.IsNotExist(err) {
			err = os.Mkdir(dirName, os.ModePerm)
			if err != nil {
				return "", err
			}
			return dirName, nil
		}
		i++
	}
}
