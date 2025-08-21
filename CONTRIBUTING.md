# Contributing to repo-slice

First off, thank you for considering contributing. This project adheres to a strict set of standards to ensure code quality and a clean version history. This document outlines the development workflow, coding style, and commit message format that all contributors must follow.

## Development Workflow

To ensure a clean version history and maintain code quality, all work must be done on a feature branch and merged via a Pull Request. Direct commits to the `main` branch are forbidden.

Our development process follows these steps:

1.  **Create an Issue**: All work must be traceable to a requirement. Before starting, ensure an issue exists that describes the feature or bug.
2.  **Create a Branch**: Create a new branch from `main` named according to our Branching Strategy.
3.  **Commit Your Work**: Make small, atomic commits that follow our Commit Message Standard.
4.  **Open a Pull Request**: When the work is complete, open a Pull Request to merge your branch into `main`.
5.  **Code Review**: The PR must be reviewed and approved by at least one other developer before it can be merged. The reviewer is responsible for ensuring the changes adhere to all project standards. **All pull requests will be automatically scanned by SonarCloud for code quality and Coveralls for test coverage.**

## Project Architecture

This project is composed of two main parts: a Go command-line tool and a GitHub Action that wraps it.

* **The Go CLI Tool (`./cmd/repo-slice`)**: This is the core engine of the project. It acts as a user-friendly wrapper around the powerful `rsync` command. Its responsibilities are to validate inputs, construct the correct `rsync` command with the appropriate filter rules, and execute it. The core file remapping logic also lives here.
* **The GitHub Action (`action.yml`)**: This is the primary, user-facing interface for the project. It's a `composite` action that orchestrates the entire workflow. Its responsibilities include parsing user inputs, downloading the correct binary, running the CLI tool, and performing workflow-specific tasks like validating the output and pushing the slice to a new branch.

When contributing, changes to the core slicing logic will likely involve modifying how the `rsync` command is constructed, while changes to the user workflow or CI/CD orchestration should be made in the `action.yml` file.

## Branching Strategy

All work, without exception, must be done on a feature branch. This ensures that the `main` branch is always stable. The following rules must be adhered to:

* **Branch Naming**: Branches must be named descriptively and reference the ticket they resolve.  The mandatory format is `[type]/[ticket-id]-short-description`. 
    * **Example (feature):** `feature/2-add-testing-framework` 
    * **Example (bug fix):** `fix/1-prevent-self-target` 
* **Short-Lived Branches**: Branches should be small in scope and short-lived.  They should be merged into `main` as soon as their single, focused task is complete and has been reviewed. 

## Environment Setup

To ensure a consistent and reproducible development environment, this project uses a `dev.nix` file to define its dependencies.

If you are working in a Nix-compatible environment like Firebase Studio, the required tools (e.g., Go, `golangci-lint`) will be automatically installed when you open the workspace. You do not need to install them manually.

For all other environments, please ensure you have a recent version of Go and `golangci-lint` installed.

## Commit Message Standard

To create an explicit and descriptive version history, we follow the **Conventional Commits specification**.  This is a mandatory requirement for all changes. 

### Format

A commit message must consist of a title and a body, separated by a blank line. 

```

\<type\>: \<A short, imperative-tense description of the change\>

\<A detailed explanation of the "why" behind the change. This body is
required for all non-trivial changes.\>

\<Signed-off-by: Author Name [author.email@example.com](mailto:author.email@example.com)\>

```

* **The Title is the "What"**: The title should be short, in the present tense, and kept under 50 characters. 
* **The Body is the "Why"**: The body is for providing context and explaining your reasoning.  It must be wrapped at 72 characters. 
* **Sign-off Trailer**: Every commit must be signed off by its author.  This certifies that you have the right to submit the work and agree to the project's terms.  This can be added automatically by using the `-s` flag when committing (e.g., `git commit -s`). 

### Commit Types

The `<type>` in the title must be one of the following: 

* **feat**: A new feature for the user. 
* **fix**: A bug fix for the user. 
* **chore**: Updates to build scripts, dependencies, or other non-user-facing changes. 
* **docs**: Changes to documentation only. 
* **style**: Formatting changes that do not affect the meaning of the code. 
* **refactor**: A code change that neither fixes a bug nor adds a feature. 
* **test**: Adding missing tests or correcting existing tests. 

### Golden Rule: Atomic Commits

Commits must be small, atomic, and represent a single logical change.  Avoid large, multi-purpose commits.

### Automated Versioning

This project uses an automated release process that creates a new version tag on every merge to `main`. The type of version bump (major, minor, or patch) is determined by special hashtags in the commit messages of a pull request.

To control the version, include one of the following hashtags in the body of your commit message:

* `#major`: For any breaking change that is not backward-compatible.
* `#minor`: For any new user-facing feature (`feat:`).
* `#patch`: For any bug fix or minor improvement (`fix:`).
* `#none`: For any change that does not affect the user, such as documentation (`docs:`), CI/CD updates (`chore:`), or refactoring (`refactor:`).

If no hashtag is provided, the version will be bumped by a `patch` release by default. If multiple commits in a PR have different hashtags, the one with the highest precedence (`major` > `minor` > `patch`) will be used.


## Coding Standards

A clean, predictable, and consistent codebase is the foundation of our development process. We automate the enforcement of these standards to eliminate debate and allow our focus to remain on functionality and architecture. 

* **Automated Formatting**: All Go code MUST be formatted with the standard `gofmt` tool.
* **Linting**: We use `golangci-lint` for identifying and fixing problematic patterns in our code.  All code must pass the linter's checks before being committed. We will use tooling like pre-commit hooks to automate this enforcement wherever possible. 
* **The "Why, Not What" Philosophy**: All comments and documentation must explain the 'why' behind the code, not the 'what' or 'when'.  The code itself explains what it's doing, and version control explains when it was changed. 
* **GoDoc Standard**: All exported functions, types, and interfaces must have a GoDoc-style comment block.  The block must provide a brief, one-sentence summary of the element's purpose and detail its parameters and what it returns.  This is a non-negotiable requirement.

## Testing Standards

This project has two main types of tests: tests for the Go CLI tool and end-to-end tests for the GitHub Action.

### Go CLI Tool Tests
This project separates fast-running **unit tests** from slower **integration tests** that interact with the file system.

* **Unit Tests**: These are located in `_test.go` files and do not require any special configuration to run.
* **Integration Tests**: These are located in `_integration_test.go` files and are marked with a `//go:build integration` build tag.

To run only the unit tests, use the standard command:
`go test -v ./...`

To run all tests, including integration tests, use the `-tags` flag:
`go test -v ./... -tags=integration`

Our CI pipeline is configured to run the complete test suite.

### GitHub Action Tests
The GitHub Action is tested via dedicated workflows located in the `.github/workflows/` directory (e.g., `test-action.yml`). These workflows are configured to run automatically on every pull request to the `main` branch.

To test changes to the `action.yml` file, simply open a pull request with your changes. The test workflows will then run, using the version of the action in your branch (`uses: ./`).

## Documentation Standards

All project documentation, from the `README.md` to code comments and technical guides, must adhere to the project's **Internal Documentation Style Guide**. This ensures our knowledge base is clear for human developers and optimally structured for AI collaborators. 

The core principles of this guide are:

* **FAIR Principles**: All documentation must be Findable, Accessible, Interoperable, and Reusable by machines. 
* **Write for the RAG Chunker**: Every page or section must be self-contained and understandable in isolation, as if it were the first thing a reader sees. 
* **Prioritize Precision**: Language must be direct, explicit, and unambiguous to minimize the risk of misinterpretation by both humans and AI agents. 
* **Use a Strict Hierarchy**: All documents must use a strict and sequential heading hierarchy (e.g., H1 → H2 → H3) to maintain a clear, logical structure. 