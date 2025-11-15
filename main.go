package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"file-counter/pkg/scanner"
	"file-counter/pkg/scanner/analyzer"
	"file-counter/pkg/scanner/types"
)

func main() {
	var scanPath string
	if len(os.Args) > 1 {
		scanPath = os.Args[1]
	} else {
		scanPath = "."
		fmt.Printf("Scanning current directory: %s\n", scanPath)
	}

	if scanPath != "." {
		fmt.Printf("Scanning directory: %s\n", scanPath)
	}

	fileScanner := scanner.NewScanner()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	resultChan := make(chan *types.ScanResult, 1)
	go func() {
		result := fileScanner.Start(scanPath)
		resultChan <- result
	}()

	var result *types.ScanResult
	select {
	case <-sigChan:
		fmt.Println("\nReceived interrupt signal. Stopping scan...")
		fileScanner.Stop()
		select {
		case result = <-resultChan:
		case <-make(chan struct{}):
		}
	case result = <-resultChan:
	}

	displayResults(result, scanPath)
}

func displayResults(result *types.ScanResult, scanPath string) {
	// Display header
	fmt.Printf("\n                    FILE SYSTEM SCAN RESULTS\n\n")

	// Basic stats
	fmt.Printf("SCANNED PATH         %s\n", scanPath)
	fmt.Printf("TOTAL FILES          %d files\n", result.TotalFiles)
	fmt.Printf("TOTAL DIRECTORIES    %d directories\n", result.TotalDirs)
	fmt.Printf("TOTAL SIZE           %s\n", formatBytes(result.TotalSize))
	fmt.Printf("SCAN DURATION        %s\n", result.ScanDuration.Round(time.Millisecond).String())
	fmt.Printf("AVERAGE FILE SIZE    %s\n", formatBytes(int64(result.AverageFileSize)))
	fmt.Printf("PROCESSING SPEED     %.2f\n", result.FilesPerSecond)
	fmt.Printf("ERRORS               %d errors\n\n", result.TotalErrors)

	// File extremes
	if result.TotalFiles > 0 {
		fmt.Printf("LARGEST FILE         %s (%s)\n", result.LargestFile.Path, formatBytes(result.LargestFile.Size))
		fmt.Printf("SMALLEST FILE        %s (%s)\n", result.SmallestFile.Path, formatBytes(result.SmallestFile.Size))
		fmt.Printf("OLDEST FILE          %s (%s)\n", result.OldestFile.Path, result.OldestFile.ModTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("NEWEST FILE          %s (%s)\n\n", result.NewestFile.Path, result.NewestFile.ModTime.Format("2006-01-02 15:04:05"))
	}

	// Top extensions
	if len(result.TopExtensions) > 0 {
		// Calculate dynamic column widths
		maxRankWidth := len("RANK")
		maxExtWidth := len("EXTENSION")
		maxCatWidth := len("CATEGORY")
		maxCountWidth := len("COUNT")
		maxSizeWidth := len("TOTAL SIZE")

		for _, ext := range result.TopExtensions {
			rankStr := fmt.Sprintf("%d", len(result.TopExtensions))
			if len(rankStr) > maxRankWidth {
				maxRankWidth = len(rankStr)
			}
			if len(ext.Extension) > maxExtWidth {
				maxExtWidth = len(ext.Extension)
			}
			category := analyzer.GetExtensionCategory(ext.Extension)
			if len(category) > maxCatWidth {
				maxCatWidth = len(category)
			}
			countStr := fmt.Sprintf("%d", ext.Count)
			if len(countStr) > maxCountWidth {
				maxCountWidth = len(countStr)
			}
			sizeStr := formatBytes(ext.TotalSize)
			if len(sizeStr) > maxSizeWidth {
				maxSizeWidth = len(sizeStr)
			}
		}

		// Ensure minimum widths
		if maxRankWidth < 4 {
			maxRankWidth = 4
		}
		if maxExtWidth < 11 {
			maxExtWidth = 11
		}
		if maxCatWidth < 9 {
			maxCatWidth = 9
		}
		if maxCountWidth < 7 {
			maxCountWidth = 7
		}
		if maxSizeWidth < 11 {
			maxSizeWidth = 11
		}

		// Build format strings
		rankFormat := fmt.Sprintf("%%-%dd ", maxRankWidth)
		extFormat := fmt.Sprintf("%%-%ds ", maxExtWidth)
		catFormat := fmt.Sprintf("%%-%ds ", maxCatWidth)
		countFormat := fmt.Sprintf("%%-%dd ", maxCountWidth)
		sizeFormat := fmt.Sprintf("%%-%ds", maxSizeWidth)

		// Print header
		header := fmt.Sprintf("%s %s %s %s %s\n",
			fmt.Sprintf("%-*s", maxRankWidth, "#"),
			fmt.Sprintf("%-*s", maxExtWidth, "EXTENSION"),
			fmt.Sprintf("%-*s", maxCatWidth, "CATEGORY"),
			fmt.Sprintf("%-*s", maxCountWidth, "COUNT"),
			fmt.Sprintf("%-*s", maxSizeWidth, "TOTAL SIZE"))
		fmt.Print(header)

		// Print separator
		separator := fmt.Sprintf("%s %s %s %s %s\n",
			fmt.Sprintf("%-*s", maxRankWidth, strings.Repeat("-", maxRankWidth)),
			fmt.Sprintf("%-*s", maxExtWidth, strings.Repeat("-", maxExtWidth)),
			fmt.Sprintf("%-*s", maxCatWidth, strings.Repeat("-", maxCatWidth)),
			fmt.Sprintf("%-*s", maxCountWidth, strings.Repeat("-", maxCountWidth)),
			fmt.Sprintf("%-*s", maxSizeWidth, strings.Repeat("-", maxSizeWidth)))
		fmt.Print(separator)

		// Print data
		for i, ext := range result.TopExtensions {
			// Skip [no extension] entries
			if ext.Extension == "[no extension]" {
				continue
			}
			// Use GetFileCategory for better lock file detection
			category := analyzer.GetFileCategory(ext.Extension)
			extDisplay := ext.Extension
			if len(extDisplay) > maxExtWidth {
				extDisplay = extDisplay[:maxExtWidth-3] + "..."
			}
			fmt.Printf(rankFormat+extFormat+catFormat+countFormat+sizeFormat+"\n",
				i+1, extDisplay, category, ext.Count, formatBytes(ext.TotalSize))
		}
		fmt.Printf("\n")
	}

	// Top directories
	if len(result.TopDirectories) > 0 {
		// Calculate dynamic column widths
		maxRankWidth := len("RANK")
		maxPathWidth := len("DIRECTORY")
		maxFilesWidth := len("FILES")
		maxDirsWidth := len("DIRS")
		maxSizeWidth := len("TOTAL SIZE")

		for _, dir := range result.TopDirectories {
			rankStr := fmt.Sprintf("%d", len(result.TopDirectories))
			if len(rankStr) > maxRankWidth {
				maxRankWidth = len(rankStr)
			}
			if len(dir.Path) > maxPathWidth {
				maxPathWidth = len(dir.Path)
			}
			filesStr := fmt.Sprintf("%d", dir.FileCount)
			if len(filesStr) > maxFilesWidth {
				maxFilesWidth = len(filesStr)
			}
			dirsStr := fmt.Sprintf("%d", dir.DirCount)
			if len(dirsStr) > maxDirsWidth {
				maxDirsWidth = len(dirsStr)
			}
			sizeStr := formatBytes(dir.TotalSize)
			if len(sizeStr) > maxSizeWidth {
				maxSizeWidth = len(sizeStr)
			}
		}

		// Ensure minimum widths and reasonable maximum
		if maxRankWidth < 4 {
			maxRankWidth = 4
		}
		if maxPathWidth < 42 {
			maxPathWidth = 42
		}
		if maxPathWidth > 60 {
			maxPathWidth = 60
		} // Cap at reasonable width
		if maxFilesWidth < 7 {
			maxFilesWidth = 7
		}
		if maxDirsWidth < 7 {
			maxDirsWidth = 7
		}
		if maxSizeWidth < 11 {
			maxSizeWidth = 11
		}

		// Build format strings
		rankFormat := fmt.Sprintf("%%-%dd ", maxRankWidth)
		pathFormat := fmt.Sprintf("%%-%ds ", maxPathWidth)
		filesFormat := fmt.Sprintf("%%-%dd ", maxFilesWidth)
		dirsFormat := fmt.Sprintf("%%-%dd ", maxDirsWidth)
		sizeFormat := fmt.Sprintf("%%-%ds", maxSizeWidth)

		// Print header
		header := fmt.Sprintf("%s %s %s %s %s\n",
			fmt.Sprintf("%-*s", maxRankWidth, "#"),
			fmt.Sprintf("%-*s", maxPathWidth, "DIRECTORY"),
			fmt.Sprintf("%-*s", maxFilesWidth, "FILES"),
			fmt.Sprintf("%-*s", maxDirsWidth, "DIRS"),
			fmt.Sprintf("%-*s", maxSizeWidth, "TOTAL SIZE"))
		fmt.Print(header)

		// Print separator
		separator := fmt.Sprintf("%s %s %s %s %s\n",
			fmt.Sprintf("%-*s", maxRankWidth, strings.Repeat("-", maxRankWidth)),
			fmt.Sprintf("%-*s", maxPathWidth, strings.Repeat("-", maxPathWidth)),
			fmt.Sprintf("%-*s", maxFilesWidth, strings.Repeat("-", maxFilesWidth)),
			fmt.Sprintf("%-*s", maxDirsWidth, strings.Repeat("-", maxDirsWidth)),
			fmt.Sprintf("%-*s", maxSizeWidth, strings.Repeat("-", maxSizeWidth)))
		fmt.Print(separator)

		// Print data
		for i, dir := range result.TopDirectories {
			displayPath := dir.Path
			if len(displayPath) > maxPathWidth {
				displayPath = "..." + displayPath[len(displayPath)-maxPathWidth+3:]
			}
			fmt.Printf(rankFormat+pathFormat+filesFormat+dirsFormat+sizeFormat+"\n",
				i+1, displayPath, dir.FileCount, dir.DirCount, formatBytes(dir.TotalSize))
		}
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
