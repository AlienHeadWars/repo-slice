# Repo Slice CLI Reference

This document provides a detailed reference for the `repo-slice` command-line tool. For a high-level overview of the project and its primary use case as a GitHub Action, please see the main [README.md](../../README.md).

## Getting Started

### Prerequisites

To use `repo-slice`, you will need the following tools installed on your system:
* **Go**: Version 1.24 or newer. You can find the official installation instructions at [go.dev/doc/install](https://go.dev/doc/install).
* **`rsync`**: This tool is a required runtime dependency. You can install it using your system's package manager:
    * **Linux (Debian/Ubuntu):** `sudo apt-get update && sudo apt-get install rsync`
    * **macOS (Homebrew):** `brew install rsync`
    * **Windows:** `rsync` is included with [Git for Windows](https://git-scm.com/download/win). Ensure it is available in your `PATH`.

### Installation

Once you have the prerequisites, you can install `repo-slice` with a single command:

```bash
go install [github.com/AlienHeadwars/repo-slice/cmd/repo-slice@latest](https://github.com/AlienHeadwars/repo-slice/cmd/repo-slice@latest)
````

This will download the source code, compile it, and place the `repo-slice` executable in your Go binary path, ready to be used.

## Usage

### 1\. Create a Manifest File

The core of `repo-slice` is the manifest file. This is a text file that uses **`rsync`'s filter-rule syntax** to define which files to include or exclude. For a complete guide on the powerful features available, see the official `rsync` documentation on FILTER RULES.

**Format Rules:**

  * Each rule must be on a new line.
  * Lines starting with `#` are treated as comments and are ignored.
  * Use `+` to include a file or directory.
  * Use `-` to exclude a file or directory.
  * The first rule that matches a file is the one that takes effect.

**Example `allow-list.txt`:**

```
# Include the main application and the slicer utility.
+ /cmd/repo-slice/main.go
+ /internal/slicer/**

# Exclude all test files from the slicer utility.
- /internal/slicer/*_test.go

# Also include the project's license and README.
+ /LICENSE
+ /README.md

# Exclude everything else.
- *
```

### 2\. Run the Command

Use the `repo-slice` command, pointing to your manifest and specifying a source and output directory.

```bash
repo-slice --manifest="allow-list.txt" --source="./source-repo" --output="./sliced-repo"
```

To remap file extensions during the slice, use the `--extension-map` flag with a comma-separated list of `old:new` pairs.

```bash
repo-slice --manifest="allow-list.txt" --source="./source-repo" --output="./sliced-repo" --extension-map="tsx:ts,mdx:md"
```

## Command-Line Reference

### Arguments

| Flag | Description | Required | Default |
| :--- | :--- | :--- | :--- |
| `--manifest` | Path to the manifest file containing filter rules. | **Yes** | |
| `--source` | The source directory to read from. | No | `.` |
| `--output` | The destination directory where the filtered copy will be created. | **Yes**| |
| `--extension-map` | A comma-separated list of `old:new` extension pairs to remap (e.g., `tsx:ts,mdx:md`). | No | |


### Exit Codes

The tool uses the following exit codes to indicate success or failure, which can be used for scripting and debugging in a CI/CD environment.

| Code | Description |
| :--- | :--- |
| `0` | Success. The repository slice was created successfully. |
| `1` | General Error. The operation failed for a variety of reasons, such as invalid arguments, file system errors, or a failed `rsync` command. Check the standard error stream for a detailed message. |