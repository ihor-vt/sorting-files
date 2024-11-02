package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"v1/analytics"
)

// ProcessFile processes and moves a single file to its categorized location
func ProcessFile(path string, d os.DirEntry, folderPath string, analytics *analytics.Analytics) {
	ext := filepath.Ext(d.Name())
	category := CategorizeFile(ext)

	if category != "" {
		categoryFolder := filepath.Join(folderPath, category)
		CreateFolderIfNotExists(categoryFolder)

		extensionFolder := filepath.Join(categoryFolder, strings.ToLower(ext)[1:])
		CreateFolderIfNotExists(extensionFolder)

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
		analytics.UpdateAnalytics(category, fileInfo.Size(), d.Name())
	}
}

// CategorizeFile determines the category of the file based on its extension
func CategorizeFile(ext string) string {
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

// CreateFolderIfNotExists creates a folder if it does not exist
func CreateFolderIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			log.Printf("Failed to create directory %s: %v\n", path, err)
		}
	}
}
