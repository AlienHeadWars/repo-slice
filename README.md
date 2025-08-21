# repo-slice

[![Coverage Status](https://coveralls.io/repos/github/AlienHeadWars/repo-slice/badge.svg)](https://coveralls.io/github/AlienHeadWars/repo-slice) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=AlienHeadWars_repo-slice&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=AlienHeadWars_repo-slice)

Automate the creation of streamlined, context-specific branches for your AI assistants.

## The Problem

As a project grows, the entire codebase quickly becomes too large to use as a single context for an AI assistant like a Gemini Gem. To work effectively, these assistants, each assigned to a specific role, need a streamlined and relevant slice of the repository.

Manually creating and maintaining these separate, context-specific branches is a tedious and error-prone process that needs to be repeated every time the codebase changes.

## The Solution

`repo-slice` is a GitHub Action that solves this problem by automating the entire workflow. It reads a manifest file, creates a clean, filtered copy of your repository, and pushes it to a dedicated branch.

This allows you to create and automatically maintain streamlined branches for each of your AI assistants, ensuring they always have the latest, most relevant context without any manual intervention.

## Example Workflow

Here is a complete, copy-paste-ready example of a GitHub Actions workflow. It generates a context for a "Go Backend Developer" AI, including only the main application, the slicer utility, and the project's license.

```yaml
# .github/workflows/update-ai-context.yml
name: Update AI Context Branches

on:
  push:
    branches:
      - 'main'

jobs:
  update-backend-developer-context:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Create Go Backend Slice
        uses: AlienHeadWars/repo-slice@v0.0.27 # Use the latest version
        with:
          manifest: |
            # Include all directories to allow for traversal.
            + **/
            # Include the main application and the slicer utility.
            + /cmd/repo-slice/**
            + /internal/slicer/**
            + /LICENSE
            # Exclude everything else.
            - *
          push-branch-name: 'context/backend-dev'
          commit-message: 'chore: Update backend-dev AI context'
```

## Using the Action

### Creating a Manifest File

The manifest file is the heart of `repo-slice`. It's a simple text file that uses **`rsync`'s filter-rule syntax** to define which files to include or exclude.

**Key Rules:**

  * **Include (`+`) and Exclude (`-`)**: Prefix each line with `+` to include a file/directory or `-` to exclude it.
  * **Order Matters**: `rsync` uses a "first match wins" logic. Place more specific rules (like excluding a single file) before more general rules (like including a whole directory).
  * **Directory Traversal**: To include a file in a subdirectory, you must also include its parent directories. The easiest way to do this is to include all directories with a `+ **/` rule at the start of your manifest.
  * **Exclude by Default**: To ensure your slice *only* contains the files you've explicitly included, you must end your manifest with a `- *` rule. This tells `rsync` to exclude everything else.

> **Warning**
> A common mistake is to forget to include the parent directories of a nested file. Without a rule like `+ **/`, `rsync` will exclude the parent directory and will never find the nested file you want to include.

For a complete guide on advanced features like inheriting rules from other files (`.`), see the official **[rsync documentation on FILTER RULES](https://download.samba.org/pub/rsync/rsync.1#FILTER_RULES)**.

### Inputs

| Input | Description | Required | Default |
| :--- | :--- | :--- | :--- |
| `manifest` | The manifest content, provided as an inline string. | No | |
| `manifestFile`| Path to the manifest file containing filter rules. | No | |
| `source` | The source directory to read from. | No | `.` |
| `output` | The destination directory. If not set, a temporary directory will be created. | No | |
| `extension-map`| A comma-separated list of `old:new` extension pairs to remap. | No | |
| `push-branch-name`| The name of the branch to push the sliced contents to. | No | |
| `commit-message`| The commit message to use when pushing the sliced branch. | No | `chore: Update repository slice` |
| `local-binary-path`| Path to a local binary. (For testing purposes). | No | |

**Note**: You must provide exactly one of `manifest` or `manifestFile`.

### Outputs

| Output | Description |
| :--- | :--- |
| `slice-path` | The path to the generated slice directory. |

## CLI Tool

This project also provides a command-line tool for local use. For detailed installation and usage instructions, please see the [CLI README](/cmd/repo-slice/README.md).

## Contributing

We welcome contributions\! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) for detailed standards and procedures.