package analytics

import "fmt"

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

type Category struct {
	count int
	size  int64
}

func NewAnalytics() *Analytics {
	return &Analytics{categoryStats: make(map[string]Category)}
}

func (a *Analytics) IncreaseSkippedFolders() {
	a.skippedFolders++
}

func (a *Analytics) IncreaseSortedFolders() {
	a.sortedFolders++
}

func (a *Analytics) IncreaseTotalEmptyFolders() {
	a.totalEmptyFolders++
}

func (a *Analytics) UpdateAnalytics(category string, fileSize int64, fileName string) {
	a.totalFilesMoved++
	a.totalSizeMoved += fileSize

	catStats := a.categoryStats[category]
	catStats.count++
	catStats.size += fileSize
	a.categoryStats[category] = catStats

	if fileSize > a.largestFileSize {
		a.largestFileSize = fileSize
		a.largestFile = fileName
	}
}

func (a *Analytics) PrintSummary() {
	fmt.Println("\n--- Analytics Summary ---")
	fmt.Printf("Total files moved: %d\n", a.totalFilesMoved)
	fmt.Printf("Total size moved: %.2f MB\n", float64(a.totalSizeMoved)/(1024*1024))
	fmt.Printf("Largest file: %s (%.2f MB)\n", a.largestFile, float64(a.largestFileSize)/(1024*1024))
	fmt.Printf("Total empty folders removed: %d\n", a.totalEmptyFolders)
	fmt.Printf("Folders sorted: %d\n", a.sortedFolders)
	fmt.Printf("Folders skipped: %d\n", a.skippedFolders)

	if len(a.categoryStats) > 0 {
		fmt.Println("\nFiles moved by category:")
		for category, stats := range a.categoryStats {
			fmt.Printf("- %s: %d files, %.2f MB\n", category, stats.count, float64(stats.size)/(1024*1024))
		}
	}
}
