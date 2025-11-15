package types

import (
	"time"
)

type FileInfo struct {
	Path      string
	Size      int64
	ModTime   time.Time
	IsDir     bool
	Extension string
}

type DirectoryStats struct {
	Path        string
	FileCount   int64
	DirCount    int64
	TotalSize   int64
	AverageSize float64
}

type ExtensionStats struct {
	Extension   string
	Count       int64
	TotalSize   int64
	AverageSize float64
	Percentage  float64
}

type ScanResult struct {
	TotalFiles   int64
	TotalDirs    int64
	TotalSize    int64
	TotalErrors  int64
	ScanDuration time.Duration

	LargestFile     FileInfo
	SmallestFile    FileInfo
	OldestFile      FileInfo
	NewestFile      FileInfo
	AverageFileSize float64

	TopExtensions []ExtensionStats

	TopDirectories []DirectoryStats

	FilesPerSecond float64
	BytesPerSecond float64
	DepthStats     map[int]int64
}

type FileAnalyzer interface {
	AnalyzeFile(path string, info FileInfo) error
	GetResults() *ScanResult
	Reset()
}

var ExtensionCategories = map[string]string{
	".txt":  "Document", // Text files
	".pdf":  "Document", // PDF files
	".doc":  "Document", // Microsoft Word documents
	".docx": "Document", // Microsoft Word documents
	".jpg":  "Image",    // JPEG images
	".jpeg": "Image",    // JPEG images
	".webp":  "Image",   // Modern web image format
	".bmp":   "Image",   // Bitmap
	".tiff":  "Image",   // High-quality photo format
	".tif":   "Image",   // TIFF alt
	".heic":  "Image",   // iPhone photos (HEIF)
	".heif":  "Image",   // HEIF general
	".ico":   "Image",   // Icons
	".psd":   "Image",   // Photoshop
	".xcf":   "Image",   // GIMP
	".png":  "Image",    // PNG images
	".gif":  "Image",    // GIF images
	".svg":  "Image",    // SVG images
	".webm": "Video",    // Web-friendly video
	".wmv":  "Video",    // Windows media
	".flv":  "Video",    // Old Flash video
	".mpeg": "Video",    // MPEG
	".mpg":  "Video",    // MPEG short extension
	".m4v":  "Video",    // MPEG-4 video (iTunes)
	".mp4":  "Video",    // MP4 videos
	".avi":  "Video",    // AVI videos
	".mov":  "Video",    // MOV videos
	".mkv":  "Video",    // MKV videos
	".mp3":  "Audio",    // MP3 audio
	".wav":  "Audio",    // WAV audio
	".flac": "Audio",    // FLAC audio
	".go":   "Code",     // Go source code
	".js":   "Code",     // JavaScript
	".ts":   "Code",     // TypeScript
	".py":   "Code",     // Python
	".java": "Code",     // Java
	".cpp":  "Code",     // C++
	".c":    "Code",     // C
	".h":    "Code",     // C header
	".rs":   "Code",     // Rust
	".rb":   "Code",     // Ruby
	".php":  "Code",     // PHP
	".pl":   "Code",     // Perl
	".asm":  "Code",     // Assembly
	".hs":   "Code",     // Haskell
	".swift": "Code",    // Swift
	".kt":    "Code",    // Kotlin
	".kts":   "Code",    // Kotlin script
	".cs":    "Code",    // C#
	".m":     "Code",    // Objective-C
	".mm":    "Code",    // Objective-C++
	".scala": "Code",    // Scala
	".dart":  "Code",    // Dart
	".sh":    "Code",    // Shell script (Bash/sh)
	".bash":  "Code",    // Bash
	".zsh":   "Code",    // Z shell script
	".fish":  "Code",    // Fish shell script
	".ps1":   "Code",    // PowerShell
	".lua":   "Code",    // Lua
	".r":     "Code",    // R
	".jl":    "Code",    // Julia
	".s":     "Code",     // Assembly (AT&T syntax)
	".v":     "Code",     // Verilog
	".vh":    "Code",     // Verilog header
	".sv":    "Code",     // SystemVerilog
	".svh":   "Code",     // SystemVerilog header
	".vhdl":  "Code",     // VHDL
	".tcl":  "Code",    // Tcl
	".awk":  "Code",    // Awk script
	".sed":  "Code",    // Sed script
	".nim":  "Code",    // Nim
	".cr":   "Code",    // Crystal
	".ex":   "Code",    // Elixir
	".exs":  "Code",    // Elixir script
	".erl":  "Code",    // Erlang
	".clj":  "Code",    // Clojure
	".cljs": "Code",    // ClojureScript
	".cljc": "Code",    // Clojure portable
	".jsx":  "Code",    // React (JavaScript XML)
	".tsx":  "Code",    // React TypeScript
	".vue":  "Code",    // Vue single-file component
	".svelte": "Code",  // Svelte component
	".ejs":   "Code",   // Embedded JS templates
	".erb":   "Code",   // Embedded Ruby templates
	".ipynb": "Code",    // Jupyter notebooks
	".rmd":   "Code",    // R Markdown with code chunks
	".ml":   "Code",     // OCaml
	".mli":  "Code",     // OCaml interface
	".fsi":  "Code",     // F# interface
	".fs":   "Code",     // F# source
	".fsx":  "Code",     // F# script
	".f90":  "Code",     // Fortran 90
	".f95":  "Code",     // Fortran 95
	".f03":  "Code",     // Fortran 2003
	".f08":  "Code",     // Fortran 2008
	".make": "Code",      // Makefile fragment
	".mk":   "Code",      // Makefile include
	".cmake": "Code",     // CMake script
	".gradle": "Code",    // Gradle script
	".groovy": "Code",    // Groovy (used by Gradle)
	".nut":  "Code",      // Squirrel
	".wat":  "Code",      // WebAssembly
	".md":   "Markup",   // Markdown
	".tex":  "Markup",   // LaTeX
	".bib":  "Markup",   // BibTeX
	".adoc": "Markup",   // Asciidoc
	".rst":  "Markup",   // reStructuredText
	".mdx":  "Markup",   // MDX (Markdown + JSX)
	".html": "Markup",   // HTML (already in your list as Web)
	".xhtml": "Markup",   // XHTML
	".css":  "Web",      // CSS
	".json": "Data",     // JSON
	".lock": "Lockfile", // Lockfile
	".xml":  "Data",     // XML
	".yaml": "Data",     // YAML
	".yml":  "Data",     // YAML
	".csv":  "Data",     // CSV
	".sql":  "Data",     // SQL
	".zip":  "Archive",  // ZIP archives
	".tar":  "Archive",  // TAR archives
	".gz":   "Archive",  // GZIP archives
	".tar.gz": "Archive", // TAR.GZ archives
	".tar.bz2": "Archive", // TAR.BZ2 archives
	".tar.xz": "Archive", // TAR.XZ archives
	".rar":  "Archive",  // RAR archives
	".7z":   "Archive",  // 7Z archives
	".exe":  "Executable", // Executable files (Windows)
	".dmg":  "Executable", // DMG files (macOS)
	".pkg":  "Executable", // PKG files (macOS)
	".deb":  "Executable", // Debian packages (Linux)
	".rpm":  "Executable", // Red Hat packages (Linux)
}

var SpecialLockFiles = map[string]string{
	"package-lock.json": "Data", // npm lockfile
	"pnpm-lock.yaml":    "Data", // pnpm lockfile
	"yarn.lock":         "Lockfile", // Yarn lockfile
	"Cargo.lock":        "Lockfile", // Cargo lockfile
	"flake.lock":        "Lockfile", // Nix flake lockfile
	"Pipfile.lock":      "Lockfile", // Pipenv lockfile
	"go.sum":            "Lockfile", // Go sum file
	"composer.lock":     "Lockfile", // Composer lockfile
	"Gemfile.lock":      "Lockfile", // Bundler lockfile
	"bun.lock":          "Lockfile", // Bun lockfile
	"bun.lockb":         "Lockfile", // Bun lockfile
}

var SpecialFilePatterns = map[string]string{
	"flake.": "Nix", // Nix flake (can be both flake.nix and flake.lock)
	".nix":   "Nix", // Nix files (mostly flakes or configs)
}
