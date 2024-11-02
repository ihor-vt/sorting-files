package analytics

import "fmt"

type Analytics struct {
	TotalFilesMoved   int
	TotalSizeMoved    int64
	TotalEmptyFolders int
	SortedFolders     int
	SkippedFolders    int
	CategoryStats     map[string]Category
	LargestFile       string
	LargestFileSize   int64
}

type Category struct {
	Count int
	Size  int64
}

func NewAnalytics() *Analytics {
	return &Analytics{CategoryStats: make(map[string]Category)}
}

func (a *Analytics) UpdateAnalytics(category string, fileSize int64, fileName string) {
	a.TotalFilesMoved++
	a.TotalSizeMoved += fileSize

	catStats := a.CategoryStats[category]
	catStats.Count++
	catStats.Size += fileSize
	a.CategoryStats[category] = catStats

	if fileSize > a.LargestFileSize {
		a.LargestFileSize = fileSize
		a.LargestFile = fileName
	}
}

func (a *Analytics) PrintSummary() {
	fmt.Println("\n--- Analytics Summary ---")
	fmt.Printf("Total files moved: %d\n", a.TotalFilesMoved)
	fmt.Printf("Total size moved: %.2f MB\n", float64(a.TotalSizeMoved)/(1024*1024))
	fmt.Printf("Largest file: %s (%.2f MB)\n", a.LargestFile, float64(a.LargestFileSize)/(1024*1024))
	fmt.Printf("Total empty folders removed: %d\n", a.TotalEmptyFolders)
	fmt.Printf("Folders sorted: %d\n", a.SortedFolders)
	fmt.Printf("Folders skipped: %d\n", a.SkippedFolders)

	fmt.Println("\nFiles moved by category:")
	for category, stats := range a.CategoryStats {
		fmt.Printf("- %s: %d files, %.2f MB\n", category, stats.Count, float64(stats.Size)/(1024*1024))
	}
}
