package analyzer

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"file-counter/pkg/scanner/types"
)

type StatisticsCollector struct {
	mu sync.RWMutex

	totalFiles  int64
	totalDirs   int64
	totalSize   int64
	totalErrors int64
	startTime   time.Time

	largestFile  types.FileInfo
	smallestFile types.FileInfo
	oldestFile   types.FileInfo
	newestFile   types.FileInfo

	extensionStats map[string]*types.ExtensionStats
	directoryStats map[string]*types.DirectoryStats

	depthStats map[int]int64
}

func NewStatisticsCollector() *StatisticsCollector {
	return &StatisticsCollector{
		startTime:      time.Now(),
		extensionStats: make(map[string]*types.ExtensionStats),
		directoryStats: make(map[string]*types.DirectoryStats),
		depthStats:     make(map[int]int64),
		smallestFile: types.FileInfo{
			Size: int64(^uint64(0) >> 1),
		},
		oldestFile: types.FileInfo{
			ModTime: time.Now(),
		},
	}
}

func (sc *StatisticsCollector) AnalyzeFile(path string, info types.FileInfo) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if info.IsDir {
		sc.totalDirs++
		return nil
	}

	sc.totalFiles++
	sc.totalSize += info.Size

	depth := strings.Count(strings.TrimPrefix(path, "/"), "/")
	sc.depthStats[depth]++

	if info.Size > sc.largestFile.Size {
		sc.largestFile = info
	}
	if info.Size < sc.smallestFile.Size && info.Size > 0 {
		sc.smallestFile = info
	}

	if info.ModTime.Before(sc.oldestFile.ModTime) {
		sc.oldestFile = info
	}
	if info.ModTime.After(sc.newestFile.ModTime) {
		sc.newestFile = info
	}

	ext := strings.ToLower(info.Extension)
	if ext == "" {
		ext = "[no extension]"
	}

	if stat, exists := sc.extensionStats[ext]; exists {
		stat.Count++
		stat.TotalSize += info.Size
	} else {
		sc.extensionStats[ext] = &types.ExtensionStats{
			Extension: ext,
			Count:     1,
			TotalSize: info.Size,
		}
	}

	dirPath := filepath.Dir(path)
	if stat, exists := sc.directoryStats[dirPath]; exists {
		stat.FileCount++
		stat.TotalSize += info.Size
	} else {
		sc.directoryStats[dirPath] = &types.DirectoryStats{
			Path:      dirPath,
			FileCount: 1,
			DirCount:  0,
			TotalSize: info.Size,
		}
	}

	return nil
}

func (sc *StatisticsCollector) AnalyzeDirectory(path string, info types.FileInfo) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.totalDirs++

	dirPath := path
	if stat, exists := sc.directoryStats[dirPath]; exists {
		stat.DirCount++
	} else {
		sc.directoryStats[dirPath] = &types.DirectoryStats{
			Path:      dirPath,
			FileCount: 0,
			DirCount:  1,
			TotalSize: 0,
		}
	}

	return nil
}

func (sc *StatisticsCollector) IncrementError() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.totalErrors++
}

func (sc *StatisticsCollector) GetResults() *types.ScanResult {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	scanDuration := time.Since(sc.startTime)

	var avgFileSize float64
	if sc.totalFiles > 0 {
		avgFileSize = float64(sc.totalSize) / float64(sc.totalFiles)
	}
	var filesPerSecond, bytesPerSecond float64
	if scanDuration.Seconds() > 0 {
		filesPerSecond = float64(sc.totalFiles) / scanDuration.Seconds()
		bytesPerSecond = float64(sc.totalSize) / scanDuration.Seconds()
	}

	topExtensions := sc.getTopExtensions(5)
	topDirectories := sc.getTopDirectories(5)

	return &types.ScanResult{
		TotalFiles:      sc.totalFiles,
		TotalDirs:       sc.totalDirs,
		TotalSize:       sc.totalSize,
		TotalErrors:     sc.totalErrors,
		ScanDuration:    scanDuration,
		LargestFile:     sc.largestFile,
		SmallestFile:    sc.smallestFile,
		OldestFile:      sc.oldestFile,
		NewestFile:      sc.newestFile,
		AverageFileSize: avgFileSize,
		TopExtensions:   topExtensions,
		TopDirectories:  topDirectories,
		FilesPerSecond:  filesPerSecond,
		BytesPerSecond:  bytesPerSecond,
		DepthStats:      sc.depthStats,
	}
}

func (sc *StatisticsCollector) Reset() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.totalFiles = 0
	sc.totalDirs = 0
	sc.totalSize = 0
	sc.totalErrors = 0
	sc.startTime = time.Now()
	sc.extensionStats = make(map[string]*types.ExtensionStats)
	sc.directoryStats = make(map[string]*types.DirectoryStats)
	sc.depthStats = make(map[int]int64)
}

func (sc *StatisticsCollector) getTopExtensions(n int) []types.ExtensionStats {
	var extensions []types.ExtensionStats
	totalFiles := sc.totalFiles

	if totalFiles == 0 {
		return extensions
	}

	for ext, stat := range sc.extensionStats {
		stat.AverageSize = float64(stat.TotalSize) / float64(stat.Count)
		stat.Percentage = (float64(stat.Count) / float64(totalFiles)) * 100
		stat.Extension = ext
		extensions = append(extensions, *stat)
	}

	sort.Slice(extensions, func(i, j int) bool {
		return extensions[i].Count > extensions[j].Count
	})

	if len(extensions) > n {
		extensions = extensions[:n]
	}

	return extensions
}

func (sc *StatisticsCollector) getTopDirectories(n int) []types.DirectoryStats {
	var directories []types.DirectoryStats

	for _, stat := range sc.directoryStats {
		if stat.FileCount > 0 {
			stat.AverageSize = float64(stat.TotalSize) / float64(stat.FileCount)
		}
		directories = append(directories, *stat)
	}

	sort.Slice(directories, func(i, j int) bool {
		return directories[i].TotalSize > directories[j].TotalSize
	})

	if len(directories) > n {
		directories = directories[:n]
	}

	return directories
}

func GetExtensionCategory(ext string) string {
	if category, exists := types.ExtensionCategories[strings.ToLower(ext)]; exists {
		return category
	}
	return "Other"
}

func GetFileCategory(filename string) string {
	if category, exists := types.SpecialLockFiles[filename]; exists {
		return category
	}

	if strings.HasPrefix(filename, "flake.") {
		return "Nix"
	}
	if strings.HasSuffix(filename, ".nix") {
		return "Nix"
	}

	ext := filepath.Ext(filename)
	if ext == "" {
		return "Other"
	}

	return GetExtensionCategory(ext)
}

func FormatBytes(bytes int64) string {
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
