package util

import "os"

//FileExist return an error when a file does not exist
func FileExist(file string) error {
	_, err := os.Stat(file)
	return err
}
