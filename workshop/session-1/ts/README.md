# Session 1: Basic CLI Client - TypeScript Setup

## Prerequisites

Before starting, ensure you have:

1. **AISHE Server Running**: The AISHE server must be accessible
   - For local development: `http://localhost:8000`
2. **Node.js 20+**: Required for this project
3. **npm 10+**: Required for this project

## Setup Instructions

### Install Node.js (recommended via NVM)

See: https://www.nvmnode.com/guide/download.html#download-nvm

**On macOS (using Homebrew):**

```bash
brew install node@24
```

**On Linux (using nvm):**

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 24
nvm use 24
```

**On Windows:**
Download and install from [nodejs.org](https://nodejs.org/en)

**Verify installation:**

```bash
node --version
npm --version
```

## Project Setup

### 1. Install Dependencies

Navigate to the project directory and install all dependencies:

```bash
cd workshop/session-1/ts
npm install
```

### 2. Configure Environment

Copy the example environment file:

```bash
cp .env.example .env
```

### 3. Verify Installation

```bash
# Check Node version
node --version

# List installed packages
npm list
```

## Running the CLI

### Starter Code

```bash
npm run start
```

### Reference Solution

```bash
npm run solution
```

### Testing Against Different Servers

**Local server:**

```bash
export AISHE_API_URL=http://localhost:8000
npm run start
```

## Troubleshooting

### node/npm command not found

If `node` or `npm` is not found after installation:

1. **Restart your shell** or source your shell configuration:

   ```bash
   source ~/.bashrc  # or ~/.zshrc
   ```

2. **Check PATH**: Ensure Node.js is in your PATH:

   ```bash
   which node
   ```

### Dependencies not installing

Try removing node_modules and starting fresh:

```bash
rm -rf node_modules
npm install
```
