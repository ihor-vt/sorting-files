package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// SelectSortingMode prompts the user to select a sorting mode
func SelectSortingMode() int {
	fmt.Println("Select sorting mode:")
	fmt.Println("1. Traditional Sorting")

	var mode int
	for {
		fmt.Print("Enter your choice (1): ")
		_, err := fmt.Scanf("%d", &mode)
		if err == nil && mode == 1 {
			break
		}
		fmt.Println("Invalid choice. Please enter 1.")
	}
	return mode
}

// Confirm prompts the user with a yes/no question
func Confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n): ", prompt)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y"
}
