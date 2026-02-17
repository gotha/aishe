# Session 1: Basic CLI Client - TypeScript Setup

## Your Task

- [ ] Implement the TODO comments inside the `main()` function in [starter.ts](./starter.ts)

## Prerequisites

Before starting, ensure you have:

1. **AISHE Server Running**: The AISHE server must be accessible. See [session-1 README](../README.md) for server setup instructions.
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

## Running the CLI

### Starter Code

```bash
npm start
```

### Reference Solution

```bash
npm run solution
```

## Troubleshooting

### node/npm command not found

If `node` or `npm` is not found after installation:

1. **Restart your shell** by spawning a new shell process:

   ```bash
   exec $SHELL
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

### Local AISHE server doesn't start

If your local server can't start on port 8000:

1. **Try a different port** and update the environment variable:

   Update the `AISHE_API_URL` in the `.env` file to the new port, then run:

   ```bash
   npm start
   ```

2. **Use the hosted AISHE server** (if available):

   Update the `AISHE_API_URL` in the `.env` file to the hosted server URL, then run:

   ```bash
   npm start
   ```
