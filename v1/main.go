package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Analytics struct to store information about moved files and directories
type Analytics struct {
	totalFilesMoved   int
	totalSizeMoved    int64
	totalEmptyFolders int
	sortedFolders     int
	skippedFolders    int
	categoryStats     map[string]Category
	largestFile       string
	largestFileSize   int64
}

// Category struct to store stats for each category
type Category struct {
	count int
	size  int64
}

// checkArgs checks if a folder path is provided as a command-line argument
func checkArgs(args []string) bool {
	return len(args) == 1
}

// categorizeFile determines the category of the file based on its extension
func categorizeFile(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg", ".png", ".gif", ".heic", ".svg":
		return "Images"
	case ".mp4", ".avi", ".mov":
		return "Videos"
	case ".mp3", ".wav":
		return "Music"
	case ".txt", ".pdf", ".docx", ".doc", ".xlsx":
		return "Documents"
	case ".epub":
		return "Books"
	case ".zip", ".rar", ".tar":
		return "Archives"
	case ".json", ".sql":
		return "Codes"
	default:
		return ""
	}
}

// confirm prompts the user with a yes/no question
func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n): ", prompt)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y"
}

// handleDirectory prompts to sort the contents of the directory
func handleDirectory(path string, analytics *Analytics) bool {
	if confirm(fmt.Sprintf("Do you want to sort files in the folder: %s?", path)) {
		analytics.sortedFolders++
		return true
	} else {
		analytics.skippedFolders++
		return false
	}
}

// processFile processes and moves a single file to its categorized location
func processFile(path string, d os.DirEntry, folderPath string, analytics *Analytics) {
	ext := filepath.Ext(d.Name())
	category := categorizeFile(ext)

	if category != "" {
		categoryFolder := filepath.Join(folderPath, category)
		createFolderIfNotExists(categoryFolder)

		extensionFolder := filepath.Join(categoryFolder, strings.ToLower(ext)[1:])
		createFolderIfNotExists(extensionFolder)

		newPath := filepath.Join(extensionFolder, d.Name())

		fileInfo, err := d.Info()
		if err != nil {
			log.Printf("Failed to get file info for %s: %v\n", path, err)
			return
		}

		// Move the file
		if err := os.Rename(path, newPath); err != nil {
			log.Printf("Failed to move file %s to %s: %v\n", path, newPath, err)
			return
		}

		// Update analytics
		updateAnalytics(analytics, category, fileInfo.Size(), d.Name())
	}
}

// createFolderIfNotExists creates a folder if it does not exist
func createFolderIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			log.Printf("Failed to create directory %s: %v\n", path, err)
		}
	}
}

// updateAnalytics updates the analytics data
func updateAnalytics(analytics *Analytics, category string, fileSize int64, fileName string) {
	analytics.totalFilesMoved++
	analytics.totalSizeMoved += fileSize

	catStats := analytics.categoryStats[category]
	catStats.count++
	catStats.size += fileSize
	analytics.categoryStats[category] = catStats

	if fileSize > analytics.largestFileSize {
		analytics.largestFileSize = fileSize
		analytics.largestFile = fileName
	}
}

// deleteEmptyFolders deletes empty folders and updates analytics
func deleteEmptyFolders(folderPath string, analytics *Analytics) {
	err := filepath.WalkDir(folderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != folderPath {
			entries, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				if err := os.Remove(path); err != nil {
					log.Printf("Failed to remove empty directory %s: %v\n", path, err)
				} else {
					analytics.totalEmptyFolders++
					fmt.Printf("Deleted empty folder: %s\n", path)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Println("Error deleting empty folders:", err)
	}
}

// printAnalytics prints the summary of the sorting operation
func printAnalytics(analytics Analytics) {
	fmt.Println("\n--- Analytics Summary ---")
	fmt.Printf("Total files moved: %d\n", analytics.totalFilesMoved)
	fmt.Printf("Total size moved: %.2f MB\n", float64(analytics.totalSizeMoved)/(1024*1024))
	fmt.Printf("Largest file: %s (%.2f MB)\n", analytics.largestFile, float64(analytics.largestFileSize)/(1024*1024))
	fmt.Printf("Total empty folders removed: %d\n", analytics.totalEmptyFolders)
	fmt.Printf("Folders sorted: %d\n", analytics.sortedFolders)
	fmt.Printf("Folders skipped: %d\n", analytics.skippedFolders)

	fmt.Println("\nFiles moved by category:")
	for category, stats := range analytics.categoryStats {
		fmt.Printf("- %s: %d files, %.2f MB\n", category, stats.count, float64(stats.size)/(1024*1024))
	}
}

func main() {
	var folderPath string
	analytics := Analytics{categoryStats: make(map[string]Category)}

	argsStatus := checkArgs(os.Args[1:])
	if argsStatus {
		folderPath = os.Args[1]
	} else {
		log.Println("Enter the path to the folder:")
		fmt.Scanln(&folderPath)
	}

	log.Println("Path:", folderPath)

	// Traverse and sort files
	err := filepath.WalkDir(folderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// If directory, confirm sorting and skip if not confirmed
		if d.IsDir() && path != folderPath {
			if !handleDirectory(path, &analytics) {
				return filepath.SkipDir
			}
		}

		// If file, process it
		if !d.IsDir() {
			processFile(path, d, folderPath, &analytics)
		}
		return nil
	})

	// Handle errors during traversal
	if err != nil {
		log.Println("Error:", err)
	}

	// Delete empty folders
	deleteEmptyFolders(folderPath, &analytics)

	// Print analytics summary
	printAnalytics(analytics)
}
