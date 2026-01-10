# Stereorepo CLI

A command-line tool to manage stereorepos (monorepos) by enabling efficient sparse checkouts. This allows working with specific subsets of a repository without downloading the entire codebase.

## Features

- **Clone**: Initialize a repository with specific modules.
- **Add**: Enable additional modules in the current workspace.
- **Sync**: Synchronize the sparse checkout configuration.

## Usage

### Clone a repository with specific modules

```bash
stereorepo clone <URL> <DESTINATION> --modules "frontend,backend"
```

### Add a module to an existing workspace

```bash
stereorepo add <MODULE_PATH>
```

## Installation

```bash
cargo install --path .
```
