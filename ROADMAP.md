# Project Roadmap

This document outlines the development milestones to get `repo-slice` to its first major release.

## Milestone 1: Foundational Setup

**Goal**: To create a well-structured repository with all foundational documentation and standards in place.

- [x] Create `README.md` to define the project's purpose and scope (Issue #1).
- [x] Create `CONTRIBUTING.md` to define development standards (Issue #2).
- [x] Initialize the Go module and "hello world" application (Issue #5).
- [x] Create a placeholder `action.yml` file (Issue #6).
- [x] Add issue templates for bugs and features (Issue #7).
- [x] Add a `dev.nix` file for a consistent development environment (Issue #9).
- [x] Create this `ROADMAP.md` file (Issue #8).

## Milestone 2: TDD-Driven Core Logic & Automation

**Goal**: To build a fully functional and rigorously tested CLI tool using Test-Driven Development.

- Set up a CI pipeline to automatically run tests and a linter.
- Write tests for and then implement CLI argument handling.
- Write tests for and then implement input validation logic.
- Write tests for and then implement the `rsync` command execution.
- Write tests for and then implement the extension mapping feature.

## Milestone 3: Release & GitHub Action

**Goal**: To make the well-tested tool easily consumable by the community.

- Create a release workflow to automatically build and attach binaries to a GitHub Release.
- Finalize the `action.yml` to use the released binaries.
- Test the end-to-end GitHub Action workflow.
- Update the `README.md` with final installation instructions.