# Mono-Clone

A CLI tool to clone ONLY specific directories from a large Git repository, saving bandwidth and disk space.

## Features
- **Bandwidth Efficient**: Uses `git clone --filter=blob:none` to avoid downloading history and unused files.
- **Sparse Checkout**: Downloads only the directories you specify.
- **Incremental Updates**: Run the command again to add more directories to an existing clone.
- **Cross-Platform**: Written in Python, works on Linux, macOS, and Windows (provided Git is installed).

## Prerequisites
- Python 3.6+
- Git installed and available in your system PATH.

## Usage

```bash
./mono-clone <REPO_URL> <DIRECTORY_1> [DIRECTORY_2 ...] [--branch BRANCH_NAME]
```

### Arguments
- `REPO_URL`: The URL of the Git repository (HTTPS or SSH).
- `DIRECTORY_n`: One or more paths within the repository to download.
- `--branch`: (Optional) The branch to clone. Defaults to `main`.

### Examples

**1. Clone a single directory:**
```bash
./mono-clone https://github.com/example/monorepo "backend/service-a"
```

**2. Add another directory to the existing clone:**
```bash
./mono-clone https://github.com/example/monorepo "frontend/ui"
```
*This will detect the existing folder and add "frontend/ui" to it.*

**3. Clone from a specific branch ('develop'):**
```bash
./mono-clone https://github.com/example/monorepo "docs" --branch develop
```
