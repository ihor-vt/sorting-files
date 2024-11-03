package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"v1/analytics"
)

// DeleteEmptyFolders deletes empty folders and updates analytics
func DeleteEmptyFolders(folderPath string, analytics *analytics.Analytics) {
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
					analytics.IncreaseTotalEmptyFolders()
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
