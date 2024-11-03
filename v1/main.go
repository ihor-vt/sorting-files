package main

import (
	"v1/ui"
	"v1/cmd"
	"v1/utils"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Run the intro using Bubble Tea
	p := tea.NewProgram(ui.InitialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// Type assertion to access the model after Bubble Tea exits
	finalModel, ok := m.(ui.Model)
	if !ok {
		fmt.Println("Error asserting model type")
		os.Exit(1)
	}

	if finalModel.Cursor == 0 {
		// User selected Standard Sorting
		var folderPath string
		if len(os.Args) > 1 {
			folderPath = os.Args[1]
		} else {
			fmt.Println("Enter the path to the folder:")
			fmt.Scanln(&folderPath)
		}

		finalModel.FolderPath = folderPath
		cmd.TraditionalSorting(folderPath, finalModel.AnalyticsData)
		utils.DeleteEmptyFolders(folderPath, finalModel.AnalyticsData)
		finalModel.AnalyticsData.PrintSummary()
	}
}
