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
    * **Example (feature):** `feature/2-add-testing-framework` 
    * **Example (bug fix):** `fix/1-prevent-self-target` 
* **Short-Lived Branches**: Branches should be small in scope and short-lived.  They should be merged into `main` as soon as their single, focused task is complete and has been reviewed. 

## Commit Message Standard

To create an explicit and descriptive version history, we follow the **Conventional Commits specification**.  This is a mandatory requirement for all changes. 

### Format

A commit message must consist of a title and a body, separated by a blank line. 

```
<type>: <A short, imperative-tense description of the change>

<A detailed explanation of the "why" behind the change. This body is
required for all non-trivial changes.>

<Signed-off-by: Author Name [author.email@dugoutsanddragons.com](mailto:author.email@dugoutsanddragons.com)>

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

## Coding Standards

A clean, predictable, and consistent codebase is the foundation of our development process. We automate the enforcement of these standards to eliminate debate and allow our focus to remain on functionality and architecture. 

* **Automated Formatting**: All Go code MUST be formatted with the standard `gofmt` tool.
* **Linting**: We use `golangci-lint` for identifying and fixing problematic patterns in our code.  All code must pass the linter's checks before being committed. We will use tooling like pre-commit hooks to automate this enforcement wherever possible. 
* **The "Why, Not What" Philosophy**: All comments and documentation must explain the 'why' behind the code, not the 'what' or 'when'.  The code itself explains what it's doing, and version control explains when it was changed. 
* **GoDoc Standard**: All exported functions, types, and interfaces must have a GoDoc-style comment block.  The block must provide a brief, one-sentence summary of the element's purpose and detail its parameters and what it returns.  This is a non-negotiable requirement. 

## Documentation Standards

*(This section will state that all documentation, including READMEs and code comments, must adhere to the project's "Internal Documentation Style Guide" to ensure clarity for both human and AI collaborators).*