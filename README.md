# fsscan - Advanced File System Scanner

A high-performance Go application that scans the entire file system starting from the root directory `/` and counts all files while displaying real-time progress.

## Features

- **Real-time Progress Display**: Shows live updates of scanned files, directories, errors, and more
- **Concurrent Processing**: Uses multiple goroutines for efficient scanning
- **Graceful Shutdown**: Handle Ctrl+C interrupts cleanly
- **Comprehensive Statistics**: File count, directory count, total size, scan speed, and error tracking
- **Smart Error Handling**: Continues scanning even when encountering permission errors
- **System Directory Skipping**: Automatically skips problematic system directories like `/proc`, `/sys`, `/dev`

## Requirements

- Go 1.21 or later
- Unix-like operating system (Linux, macOS)
- Root/sudo privileges recommended for complete system access

>[!NOTE]
> I did a horrible job at making this, since its a rewrite of an old project I made, this may compile weirdly, you might need to figure out how to add the command "fs" to your path, its stupidly complex at least for nushell, but it should work for bash, zsh, and fish without any issues, I also think because I use nix home-manager, that is probably why I had issues exporting the path, you don't need nix to use this, I don't think anyways. If you have any issues, please open an issue on github :)

## Installation

1. Clone or download this project
2. Navigate to the project directory:
   ```bash
   cd "fsscan"
   ```
3. Build the applications:
   ```bash
   go build -o fs .
   ```

## Usage

### Full System Scan
```bash
./fs              # Basic scan
sudo ./fs         # With root privileges (recommended if your system is large or your running it from a root directory)
```


Root privileges provide access to all system files and directories that would otherwise be restricted.

## Output Example

```
Scanning directory: /Users/playfairs/.nix

                    FILE SYSTEM SCAN RESULTS

SCANNED PATH         /Users/playfairs/.nix
TOTAL FILES          104 files
TOTAL DIRECTORIES    25 directories
TOTAL SIZE           14.9 MB
SCAN DURATION        4ms
AVERAGE FILE SIZE    146.9 KB
PROCESSING SPEED     23192.28
ERRORS               0 errors

LARGEST FILE         /Users/playfairs/.nix/wallpapers/feild.jpg (4.1 MB)
SMALLEST FILE        /Users/playfairs/.nix/modules/home/analygits/.gitignore (7 B)
OLDEST FILE          /Users/playfairs/.nix/.gitattributes (2025-11-04 15:20:05)
NEWEST FILE          /Users/playfairs/.nix/modules/home/shells/default.nix (2025-11-15 14:33:47)

#    EXTENSION      CATEGORY  COUNT   TOTAL SIZE 
---- -------------- --------- ------- -----------
1    .nix           Nix       64      96.2 KB    
3    .png           Image     7       6.4 MB     
4    .jpg           Image     7       8.3 MB     
5    .lock          Lockfile  5       84.7 KB    

#    DIRECTORY                                    FILES   DIRS    TOTAL SIZE 
---- -------------------------------------------- ------- ------- -----------
1    /Users/playfairs/.nix/wallpapers             8       1       8.6 MB     
2    /Users/playfairs/.nix/.res                   5       1       6.0 MB     
3    /Users/playfairs/.nix/modules/home/fastfetch 2       1       128.0 KB   
4    /Users/playfairs/.nix                        9       1       48.0 KB    
5    /Users/playfairs/.nix/modules/home           20      1       22.6 KB    
```

## Understanding the Output

- **Scanned Files**: Total number of regular files found
- **Dirs**: Total number of directories processed
- **Errors**: Files/directories that couldn't be accessed (usually permission issues)
- **Skipped**: System directories automatically skipped for safety
- **Size**: Total size of all scanned files
- **Current**: The file/directory currently being processed
- **Last Error**: Most recent error encountered

## Performance Considerations

- **Memory Usage**: The application uses minimal memory as it doesn't store file lists
- **CPU Usage**: Uses multiple goroutines (default: 2x CPU cores) for parallel processing
- **I/O Performance**: Optimized for fast directory traversal
- **Large File Systems**: Can handle millions of files efficiently

## Safety Features

- **System Directory Protection**: Automatically skips dangerous system directories
- **Graceful Interruption**: Ctrl+C stops the scan cleanly and shows partial results
- **Error Resilience**: Continues scanning even when individual files cause errors
- **Permission Handling**: Gracefully handles permission denied errors

## Troubleshooting

### Permission Errors
If you see many permission errors, run with sudo:
```bash
sudo ./fs
```

### Testing Before Full Scan
Use the demo version to test on smaller directories first:
```bash
./fs ~/Downloads
```

### Slow Performance
- Close unnecessary applications to free up I/O resources
- Consider running on an SSD for better performance
- The scan speed depends on your storage device and file system structure

### High Memory Usage
The application is designed to use minimal memory, but scanning very large directories with millions of files might use more system resources.

## Project Structure

```
.
├── build.sh
├── flake.lock
├── flake.nix
├── fs
├── go.mod
├── install.sh
├── main.go
├── Makefile
├── pkg
│   └── scanner
│       ├── analyzer
│       │   └── analyzer.go
│       ├── scanner.go
│       └── types
│           └── types.go
├── README.md
└── USAGE.md

5 directories, 13 files
```

## Technical Details

- **Language**: Go 1.21+
- **Concurrency**: Worker pool pattern with configurable goroutines
- **File System API**: Uses Go's `filepath.Walk` and `os.Lstat`
- **Progress Updates**: Real-time updates every 50ms
- **Architecture**: Concurrent producer-consumer pattern

## Contributing

Feel free to submit issues or pull requests to improve the application, not sure what you could possibly wanna do to this, but yea 😭

## License

This project is licensed under the Do What The Fuck You Want To Public License (WTFPL), yea, that means you can do whatever you want with it, I genuinely don't care 😭✌️

## Warning

If you run this on a root directory, it will scan the entire file system starting from root `/`. On systems with large amounts of data, this can:
- Take several hours to complete
- Generate significant disk I/O
- Require substantial system permissions

Always ensure you have adequate time and system resources before starting a full system scan.
