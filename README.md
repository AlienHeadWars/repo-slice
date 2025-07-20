# repo-slice

Automate the creation of streamlined, context-specific branches for your AI assistants.

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

*(This section will use a numbered list to describe the step-by-step CI/CD workflow: Define -> Slice -> Push, as per rule:styleguide:lists:20250720-001550-008).*

## Getting Started

*(This parent section will house the initial steps for a new user).*

### Installation

*(This section will contain the steps required to install the repo-slice CLI, to be filled in upon first release).*

### Basic Usage

*(This section will provide a simple, self-contained code snippet showing the most common command-line invocation, as required by rule:styleguide:code-snippets:20250720-001550-010).*

## Command-Line Reference

*(This section will serve as the formal API documentation for the tool).*

### Arguments

*(This section will use a table to detail all command-line flags, their purpose, and whether they are required, as per rule:styleguide:tables:20250720-001550-009).*

### Exit Codes

*(This section will formally document the tool's exit codes to aid in scripting and debugging, as per rule:styleguide:formal-syntax:20250720-001550-014).*

## Contributing

*(This section will briefly explain the project's openness to contributions and link directly to the CONTRIBUTING.md file for detailed standards and procedures).*