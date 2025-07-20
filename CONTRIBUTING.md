# Contributing to repo-slice

First off, thank you for considering contributing. This project adheres to a strict set of standards to ensure code quality and a clean version history. This document outlines the development workflow, coding style, and commit message format that all contributors must follow.

## Development Workflow

To ensure a clean version history and maintain code quality, all work must be done on a feature branch and merged via a Pull Request. Direct commits to the `main` branch are forbidden.

Our development process follows these steps:

1.  **Create an Issue**: All work must be traceable to a requirement. Before starting, ensure an issue exists that describes the feature or bug.
2.  **Create a Branch**: Create a new branch from `main` named according to our Branching Strategy.
3.  **Commit Your Work**: Make small, atomic commits that follow our Commit Message Standard.
4.  **Open a Pull Request**: When the work is complete, open a Pull Request to merge your branch into `main`.
5.  **Code Review**: The PR must be reviewed and approved by at least one other developer before it can be merged. The reviewer is responsible for ensuring the changes adhere to all project standards.

## Branching Strategy

All work, without exception, must be done on a feature branch. This ensures that the `main` branch is always stable. The following rules must be adhered to:

* **Branch Naming**: Branches must be named descriptively and reference the ticket they resolve.  The mandatory format is `[type]/[ticket-id]-short-description`. 
    * **Example (feature):** `feature/TD-002-add-testing-framework` 
    * **Example (bug fix):** `fix/BUG-001-prevent-self-target` 
* **Short-Lived Branches**: Branches should be small in scope and short-lived.  They should be merged into `main` as soon as their single, focused task is complete and has been reviewed. 

## Commit Message Standard

*(This section will define the Conventional Commits specification as the mandatory format for all commit messages).*

## Coding Standards

*(This section will codify the project's coding standards, including the "Why, Not What" philosophy for comments and the required Go tooling like `gofmt` and `golangci-lint` for automated style enforcement).*

## Documentation Standards

*(This section will state that all documentation, including READMEs and code comments, must adhere to the project's "Internal Documentation Style Guide" to ensure clarity for both human and AI collaborators).*