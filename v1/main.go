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
	totalFilesMoved       int                 // Total number of files moved
	totalSizeMoved        int64               // Total size of files moved in bytes
	totalEmptyFolders     int                 // Total empty folders removed
	sortedFolders         int                 // Number of folders sorted
	skippedFolders        int                 // Number of folders skipped
	categoryStats         map[string]Category // Stats per category
	largestFile           string              // Name of the largest file moved
	largestFileSize       int64               // Size of the largest file moved in bytes
}

// Category struct to store stats for each category
type Category struct {
	count int   // Number of files in this category
	size  int64 // Total size of files in this category
}

// checkArgs checks if a folder path is provided as a command-line argument
func checkArgs(args []string) bool {
	log.Println(args)
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

	// Traverse all files and folders in the specified directory using WalkDir
	err := filepath.WalkDir(folderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// If the entry is a directory, prompt to sort contents and handle empty folders
		if d.IsDir() && path != folderPath {
			if confirm(fmt.Sprintf("Do you want to sort files in the folder: %s?", path)) {
				analytics.sortedFolders++
			} else {
				analytics.skippedFolders++
				return filepath.SkipDir // Skip this directory
			}
		}

		// Process files within the folder
		if !d.IsDir() {
			ext := filepath.Ext(d.Name())
			category := categorizeFile(ext)

			// If the file has a recognized category
			if category != "" {
				// Create the category folder path
				categoryFolder := filepath.Join(folderPath, category)

				// Create the folder if it doesn't exist
				if _, err := os.Stat(categoryFolder); os.IsNotExist(err) {
					err := os.Mkdir(categoryFolder, 0755)
					if err != nil {
						return fmt.Errorf("failed to create directory %s: %v", categoryFolder, err)
					}
				}

				// Create the extension folder path
				categoryExtensionFolder := filepath.Join(categoryFolder, strings.ToLower(ext)[1:])

				// Create the folder if it doesn't exist
				if _, err := os.Stat(categoryExtensionFolder); os.IsNotExist(err) {
					err := os.Mkdir(categoryExtensionFolder, 0755)
					if err != nil {
						return fmt.Errorf("failed to create directory %s: %v", categoryExtensionFolder, err)
					}
				}

				// Define the new path for the file
				newPath := filepath.Join(categoryExtensionFolder, d.Name())

				// Get file information to calculate size
				fileInfo, err := d.Info()
				if err != nil {
					return fmt.Errorf("failed to get file info for %s: %v", path, err)
				}
				fileSize := fileInfo.Size()

				// Update analytics data
				analytics.totalFilesMoved++
				analytics.totalSizeMoved += fileSize

				// Update category statistics
				catStats := analytics.categoryStats[category]
				catStats.count++
				catStats.size += fileSize
				analytics.categoryStats[category] = catStats

				// Check if this file is the largest file encountered
				if fileSize > analytics.largestFileSize {
					analytics.largestFileSize = fileSize
					analytics.largestFile = d.Name()
				}

				// Move the file to the new path
				err = os.Rename(path, newPath)
				if err != nil {
					return fmt.Errorf("failed to move file %s to %s: %v", path, newPath, err)
				}
				// fmt.Printf("Moved %s to %s\n", path, newPath)
			}
		}
		return nil
	})

	if err != nil {
		log.Println("Error:", err)
	}

	// Delete empty folders
	err = filepath.WalkDir(folderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != folderPath {
			// Check if the directory is empty
			entries, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				err := os.Remove(path)
				if err != nil {
					return fmt.Errorf("failed to remove empty directory %s: %v", path, err)
				}
				analytics.totalEmptyFolders++
				fmt.Printf("Deleted empty folder: %s\n", path)
			}
		}
		return nil
	})

	if err != nil {
		log.Println("Error:", err)
	}

	// Print analytics summary
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
