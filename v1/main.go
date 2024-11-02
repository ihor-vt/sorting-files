package main

import (
	"v1/cmd"
	"v1/analytics"
	"v1/utils"
	"fmt"
	"log"
	"os"
)

func main() {
	var folderPath string
	analyticsData := analytics.NewAnalytics()

	mode := utils.SelectSortingMode()

	if len(os.Args) > 1 {
		folderPath = os.Args[1]
	} else {
		log.Println("Enter the path to the folder:")
		fmt.Scanln(&folderPath)
	}

	if mode == 1 {
		cmd.TraditionalSorting(folderPath, analyticsData)
	} else {
		fmt.Println("Invalid mode. Exiting program.")
		return
	}

	utils.DeleteEmptyFolders(folderPath, analyticsData)
	analyticsData.PrintSummary()
}
