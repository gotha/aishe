# Session 1: Basic CLI Client - Python Setup

## Prerequisites

Before starting, ensure you have:

1. **AISHE Server Running**: The AISHE server must be accessible
   - For local development: `http://localhost:8000`

2. **Python 3.12+**: Required for this project

## Setup Instructions

### Install uv 

install uv:

**On macOS/Linux:**
```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
```

**On Windows:**
```powershell
powershell -ExecutionPolicy ByPass -c "irm https://astral.sh/uv/install.ps1 | iex"
```

**Verify installation:**
```bash
uv --version
```

## Project Setup

### Sync Dependencies

Install all project dependencies defined in `pyproject.toml`:

```bash
# This will create a virtual environment and install dependencies
uv sync
```

This command will:
- Create a `.venv` directory with a Python virtual environment
- Install all dependencies from `pyproject.toml`
- Generate/update `uv.lock` file for reproducible builds

### 2. Activate the Virtual Environment

**On macOS/Linux:**
```bash
source .venv/bin/activate
```

**On Windows:**
```powershell
.venv\Scripts\activate
```

### 3. Verify Installation

```bash
# Check Python version
python --version

# List installed packages
uv pip list
```

## Running the CLI

### Using Python Directly

```bash
python main.py "What is the capital of France?"
```

### Testing Against Different Servers

**Local server:**
```bash
export AISHE_SERVER_URL=http://localhost:8000
python main.py "What is Redis?"
```

## Troubleshooting

### uv command not found

If `uv` is not found after installation:

1. **Restart your shell** or source your shell configuration:
   ```bash
   source ~/.bashrc  # or ~/.zshrc
   ```

2. **Check PATH**: Ensure `~/.cargo/bin` is in your PATH:
   ```bash
   echo $PATH | grep cargo
   ```

### Virtual environment not activating

Make sure you're in the correct directory:
```bash
cd workshop/session-1/python
ls -la .venv  # Should show the virtual environment directory
```

### Dependencies not installing

Try removing the virtual environment and starting fresh:
```bash
rm -rf .venv
uv sync
```

