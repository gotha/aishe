{
  description =
    "Development environment with pandoc and reveal.js presentation builder";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        revealjs = pkgs.fetchFromGitHub {
          owner = "hakimel";
          repo = "reveal.js";
          rev = "4.3.1";
          sha256 = "sha256-saCLoiP6ZtCPxlFV6exNoo9Lkq+zj0FPFPpnBRP4SPI=";
        };

        presentation = pkgs.stdenv.mkDerivation {
          name = "of25-nix-presentation";
          src = ./.;

          buildInputs = with pkgs; [ pandoc ];

          buildPhase = ''
            mkdir -p $out

            # Copy reveal.js to output
            cp -r ${revealjs} $out/reveal.js
            chmod -R +w $out/reveal.js

            cp -r img $out/

            # Find background image
            BACKGROUND_IMAGE="./img/background.png"

            echo "Using background image: $BACKGROUND_IMAGE"
            pandoc slides.md \
              -t revealjs \
              -s \
              --css slides.css \
              -o $out/index.html \
              --slide-level=2 \
              --variable theme=white \
              --variable transition=slide \
              --variable parallaxBackgroundImage="$BACKGROUND_IMAGE" \
              --variable parallaxBackgroundSize="cover" \
              --variable revealjs-url="./reveal.js/" \
              --variable notes=true

            # Add cache-busting meta tags to the HTML
            sed -i 's|<head>|<head>\n  <meta http-equiv="Cache-Control" content="no-cache, no-store, must-revalidate">\n  <meta http-equiv="Pragma" content="no-cache">\n  <meta http-equiv="Expires" content="0">|' $out/index.html

            # Add script to automatically open speaker notes
            # sed -i 's|</body>|<script>\n  Reveal.addEventListener("ready", function() {\n    // Auto-open speaker notes\n    if (!window.location.search.includes("receiver")) {\n      Reveal.getPlugin("notes").open();\n    }\n  });\n</script>\n</body>|' $out/index.html
          '';

          installPhase = ''
            cp slides.css $out/
          '';
        };

      in {
        packages.default = presentation;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ pandoc python312 ];

          shellHook = ''
            echo "Development environment with pandoc ready!"
            echo "Pandoc version: $(pandoc --version | head -n1)"
          '';
        };
      });
}
