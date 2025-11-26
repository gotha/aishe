{
  description = "RAG application with Wikipedia MCP server and Ollama";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Python environment
            python311
            python311Packages.pip
            python311Packages.virtualenv

            # Ollama for local LLM
            ollama

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

            # Install/update Python dependencies from requirements.txt
            if [ -f requirements.txt ]; then
              echo "Installing Python dependencies from requirements.txt..."
              pip install -q -r requirements.txt
            fi

            echo "=================================="
            echo "RAG Development Environment Loaded"
            echo "=================================="
            echo "Python: $(python --version)"
            echo "Ollama: $(ollama --version)"
            echo "Virtual environment: .venv"
            echo ""
            echo "To start Ollama service, run:"
            echo "  ollama serve"
            echo "=================================="
          '';
        };
      });
}
