# repo-slice

[![Coverage Status](https://coveralls.io/repos/github/AlienHeadWars/repo-slice/badge.svg)](https://coveralls.io/github/AlienHeadWars/repo-slice) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=AlienHeadWars_repo-slice&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=AlienHeadWars_repo-slice)

Automate the creation of streamlined, context-specific branches for your AI assistants.

## 🚧 Development Status

**This project is currently in active development and is not yet ready for production use.**

We are in the process of building out the core features as outlined in our [Project Roadmap](ROADMAP.md). We welcome feedback and contributions! If you're interested in helping, please see our [Contributing Guide](CONTRIBUTING.md).

## The Problem

As a project grows, the entire codebase quickly becomes too large to use as a single context for an AI assistant like a Gemini Gem. To work effectively, these assistants, each assigned to a specific role, need a streamlined and relevant slice of the repository.

Manually creating and maintaining these separate, context-specific branches is a tedious and error-prone process that needs to be repeated every time the codebase changes.

## The Solution

`repo-slice` is a simple, automation-focused CLI tool designed to solve this problem. It acts as the engine of a CI/CD pipeline that automatically maintains streamlined branches for your AI assistants.

The workflow is straightforward:
1.  **Define**: List the files and directories for a specific AI role in a simple `allow-list.txt` manifest.
2.  **Slice**: In a CI/CD job, `repo-slice` reads the manifest and creates a clean, filtered copy of your repository.
3.  **Push**: The job then pushes this filtered copy to a dedicated branch (e.g., `context/gem-ui-developer`).
4.  **Configure**: Point your AI assistant at this branch, and it will always have the latest, relevant context without any manual updates.

## Core Features

* **Declarative Manifests**: Instead of complex filter patterns, you use a simple, explicit list of files in a manifest. This is easier to read, maintain, and generate programmatically.
* **Automation-Focused**: Designed from the ground up to be a reliable and portable tool for any CI/CD environment like GitHub Actions or GitLab CI.
* **Branch-Ready Output**: Produces a clean directory structure ready to be committed to a new branch, unlike tools that generate a single text file for prompting.
* **Extension Mapping**: Optionally remap file extensions during the copy process. This is useful for improving compatibility with tools that don't recognize certain extensions, such as renaming `.tsx` files to `.ts` for better LLM interpretation.

## Workflow Overview

`repo-slice` is designed to be the engine of a fully automated CI/CD pipeline. The typical workflow is as follows:

1.  **Define**: In your main branch, create and maintain a manifest file (e.g., `roles/ai-docs-writer.allow.txt`). This file declaratively lists every file and directory that is relevant to a specific AI role.
2.  **Slice**: A scheduled CI/CD job checks out your repository. It then runs the `repo-slice` command, pointing to your manifest file, which creates a clean, filtered directory.
3.  **Push**: The CI/CD job forcefully pushes the contents of this newly created directory to a dedicated context branch (e.g., `context/ai-docs-writer`), overwriting its previous contents.
4.  **Configure**: Your AI assistant (e.g., a Gemini Gem) is configured to use this specific context branch as its knowledge source, ensuring it always has the latest, most relevant information without any manual intervention.

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
go install github.com/AlienHeadwars/repo-slice/cmd/repo-slice@latest
````

This will download the source code, compile it, and place the `repo-slice` executable in your Go binary path, ready to be used.

## Usage

### 1\. Create a Manifest File

The core of `repo-slice` is the manifest file. This is a simple text file (e.g., `allow-list.txt`) that explicitly lists every file and directory you want to include in your slice.

**Format Rules:**

  * Each entry must be on a new line.
  * Paths should be relative to the `--source` directory.
  * Lines starting with `#` are treated as comments and are ignored.
  * Empty lines are ignored.

**Example `allow-list.txt`:**

```
# Include the main package and the slicer utility.
cmd/repo-slice/main.go
internal/slicer/

# Also include the project's license and README.
LICENSE
README.md
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
| `--manifest` | Path to the "allow-list" file containing paths to include (one per line). | **Yes** | |
| `--source` | The source directory to read from. | No | `.` |
| `--output` | The destination directory where the filtered copy will be created. | **Yes**| |
| `--extension-map` | A comma-separated list of `old:new` extension pairs to remap (e.g., `tsx:ts,mdx:md`). | No | |

### Exit Codes

*(This section will formally document the tool's exit codes to aid in scripting and debugging).*

## Quality Assurance

This project is committed to a high standard of code quality and security. To ensure this, we have integrated the following tools into our development process:

* **Coveralls**: For tracking test coverage on every pull request and ensuring it remains high.
* **SonarCloud**: For continuous static analysis to detect bugs, vulnerabilities, and code smells.
* **Snyk**: For scanning dependencies against a database of known open-source vulnerabilities.
* **Dependabot**: For automatically keeping our dependencies up-to-date.

## Versioning

This project follows [Semantic Versioning](https://semver.org/). We use an automated release process that creates a new version tag on every merge to the `main` branch. See our [Contributing Guide](CONTRIBUTING.md#automated-versioning) for details on how commit messages influence the version number.

## Contributing

We welcome contributions! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) for detailed standards and procedures.