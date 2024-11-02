package cmd

import (
	"v1/analytics"
	"v1/utils"
	"path/filepath"
	"os"
	"log"
	"fmt"
)

func TraditionalSorting(folderPath string, analytics *analytics.Analytics) {
	err := filepath.WalkDir(folderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != folderPath {
			if !utils.Confirm(fmt.Sprintf("Do you want to sort files in the folder: %s?", path)) {
				analytics.SkippedFolders++
				return filepath.SkipDir
			}
			analytics.SortedFolders++
		}
		if !d.IsDir() {
			utils.ProcessFile(path, d, folderPath, analytics)
		}
		return nil
	})

	if err != nil {
		log.Println("Error:", err)
	}
}
