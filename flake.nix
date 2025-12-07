{
  description = "RAG application with Wikipedia MCP server and Ollama";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # Node.js for JS/TS tooling (nixpkgs 25.05 provides nodejs_22 as LTS).
        # This will not be exactly v25.2.1, but is good enough for the tooling.
        nodejs = pkgs.nodejs_22;

        # Create a Python environment with dependencies
        pythonEnv = pkgs.python311.withPackages (ps:
          with ps; [
            fastapi
            mcp
            ollama
            pip
            virtualenv
            # Native Nix packages from requirements.txt
            platformdirs
            # Note: mcp, wikipedia-mcp, ollama, fastapi, uvicorn, httpx are not available in nixpkgs
            # These will be installed via pip in the shellHook
          ]);

        # Google Cloud SDK with GKE auth plugin
        gcloud = pkgs.google-cloud-sdk.withExtraComponents
          [ pkgs.google-cloud-sdk.components.gke-gcloud-auth-plugin ];
      in {
        apps.aishe = {
          type = "app";
          program = "${pkgs.writeShellScript "aishe" ''
            # Ensure we're in the project directory
            cd ${self}

            # Activate virtual environment if it exists
            if [ -d .venv ]; then
              source .venv/bin/activate
            fi

            # Run the CLI
            exec ${pythonEnv}/bin/python src/cli.py "$@"
          ''}";
        };

        apps.server = {
          type = "app";
          program = "${pkgs.writeShellScript "aishe-server" ''
            # Ensure we're in the project directory
            cd ${self}

            # Activate virtual environment if it exists
            if [ -d .venv ]; then
              source .venv/bin/activate
            fi

            # Run the server
            exec ${pythonEnv}/bin/python src/server.py "$@"
          ''}";
        };

        # JS CLI: builds everything and then runs `npm start` in bin/ts
        apps."aishe-js" = {
          type = "app";
          program = "${pkgs.writeShellScript "aishe-js" ''
            set -euo pipefail
            
            # echo "=== Building aishe-client (client-js) ==="
            cd clients/js/
            ${nodejs}/bin/npm install
            ${nodejs}/bin/npm run build

            echo "=== Running AIshe JS CLI (npm start) ==="
            exec ${nodejs}/bin/npm start "$@"
          ''}";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Python environment with native Nix packages
            pythonEnv

            # Python package manager
            uv

            # Ollama for local LLM
            ollama

            # Cloud and Kubernetes tools
            gcloud
            kubectl
            kubernetes-helm
            nodejs

            # Development tools
            git
          ];

          shellHook = ''
            # Create virtual environment if it doesn't exist
            if [ ! -d .venv ]; then
              echo "Creating Python virtual environment..."
              python -m venv .venv
            fi
            source .venv/bin/activate

            # Install remaining Python dependencies not available in nixpkgs
            # (mcp, wikipedia-mcp, ollama client, fastapi, uvicorn, httpx)
            if [ -f requirements.txt ]; then
              echo "Installing additional Python dependencies from requirements.txt..."
              pip install -q mcp wikipedia-mcp ollama fastapi 'uvicorn[standard]' httpx
            fi

            echo "========================================"
            echo "AISHE Development Environment Loaded"
            echo "========================================"
            echo "Python: $(python --version)"
            echo "uv: $(uv --version)"
            echo "Ollama: $(ollama --version)"
            echo "Virtual environment: .venv"
            echo "Node: $(${nodejs}/bin/node --version)"
            echo ""
            echo "Available commands:"
            echo "  nix run .#server   - Start the API server"
            echo "  nix run .#aishe    - Run the Python CLI client"
            echo "  nix run .#aishe-js - Build & run the JS CLI"
            echo ""
            echo "Workshop setup:"
            echo "  cd workshop/session-X/python && uv sync"
            echo ""
            echo "To start Ollama service, run:"
            echo "  ollama serve"
            echo "========================================"
          '';
        };
      });
}
