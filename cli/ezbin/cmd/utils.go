package cmd

import (
	"fmt"
)

// prompt asks the user if they want to do X with a response of y or n
func promptConfirmation(promptString string) (bool, error) {
	fmt.Print(promptString, " [y/N]")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		if err.Error() != "unexpected newline" {
			return false, err
		}
	}

	if response == "y" || response == "Y" {
		return true, nil
	}

	return false, nil
}
