{
  description = "TinyClaw - Verified agent framework";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        version =
          if self ? rev
          then builtins.substring 0 8 self.rev
          else "dev";

        ldflags = [
          "-X github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal.version=${version}"
          "-X github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal.gitCommit=${version}"
          "-s" "-w"
        ];
      in
      {
        packages = {
          # Go gateway binary
          tinyclaw = pkgs.buildGoModule {
            pname = "tinyclaw";
            inherit version;
            src = ./.;
            vendorHash = "sha256-K3VY1oBTfb0suCHDYvR9zmSvXMNW31qiRH0R5BFsY9A=";
            env.CGO_ENABLED = "0";
            tags = [ "stdjson" ];
            inherit ldflags;
            subPackages = [ "cmd/tinyclaw" ];

            preBuild = ''
              # go:generate copies workspace/ into onboard package for embedding
              cp -r workspace cmd/tinyclaw/internal/onboard/workspace
            '';

            # Skip tests that require network
            doCheck = false;

            meta = {
              description = "Ultra-lightweight personal AI agent";
              license = pkgs.lib.licenses.mit;
            };
          };

          # Dhall config package - renders all configs to JSON
          dhall-config = pkgs.stdenv.mkDerivation {
            pname = "tinyclaw-dhall-config";
            inherit version;
            src = ./dhall;

            nativeBuildInputs = with pkgs; [ dhall dhall-json ];

            buildPhase = ''
              # Type-check all Dhall files
              find . -name '*.dhall' -exec dhall type --file {} \; > /dev/null

              # Render examples
              mkdir -p rendered
              for example in examples/*.dhall; do
                name=$(basename "$example" .dhall)
                dhall-to-json --file "$example" --output "rendered/$name.json"
              done
            '';

            installPhase = ''
              mkdir -p $out/share/tinyclaw
              cp -r rendered/* $out/share/tinyclaw/
              cp -r types $out/share/tinyclaw/types
            '';
          };

          # Full bundle: gateway + dhall config + default rendered configs
          tinyclaw-bundle = pkgs.symlinkJoin {
            name = "tinyclaw-bundle-${version}";
            paths = [
              self.packages.${system}.tinyclaw
              self.packages.${system}.dhall-config
            ];
            postBuild = ''
              # Verify both components are present
              test -x $out/bin/tinyclaw || (echo "Missing tinyclaw binary" && exit 1)
              test -d $out/share/tinyclaw || (echo "Missing dhall config" && exit 1)
            '';
          };

          # Docker image via pkgs.dockerTools
          tinyclaw-docker = pkgs.dockerTools.buildLayeredImage {
            name = "tinyclaw";
            tag = version;
            contents = [
              self.packages.${system}.tinyclaw
              self.packages.${system}.dhall-config
              pkgs.cacert
              pkgs.tzdata
            ];
            config = {
              Cmd = [ "/bin/tinyclaw" "gateway" ];
              Env = [
                "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
                "TZDIR=${pkgs.tzdata}/share/zoneinfo"
              ];
              ExposedPorts = {
                "18790/tcp" = {};
              };
            };
          };

          # F*-extracted verified core binary (OCaml)
          tinyclaw-core = pkgs.stdenv.mkDerivation {
            pname = "tinyclaw-core";
            inherit version;
            src = ./fstar/extracted;

            nativeBuildInputs = with pkgs; [
              ocaml
              dune_3
              ocamlPackages.findlib
              ocamlPackages.yojson
            ];

            buildPhase = ''
              dune build
            '';

            installPhase = ''
              mkdir -p $out/bin
              cp _build/default/bin/main.exe $out/bin/tinyclaw-core
            '';

            meta = {
              description = "TinyClaw verified core (F*-extracted)";
              license = pkgs.lib.licenses.mit;
            };
          };

          # Full verified bundle: gateway + core + dhall config
          tinyclaw-verified-bundle = pkgs.symlinkJoin {
            name = "tinyclaw-verified-bundle-${version}";
            paths = [
              self.packages.${system}.tinyclaw
              self.packages.${system}.tinyclaw-core
              self.packages.${system}.dhall-config
            ];
            postBuild = ''
              test -x $out/bin/tinyclaw || (echo "Missing gateway binary" && exit 1)
              test -x $out/bin/tinyclaw-core || (echo "Missing verified core binary" && exit 1)
              test -d $out/share/tinyclaw || (echo "Missing dhall config" && exit 1)
            '';
          };

          # Docker image with verified core
          tinyclaw-verified-docker = pkgs.dockerTools.buildLayeredImage {
            name = "tinyclaw-verified";
            tag = version;
            contents = [
              self.packages.${system}.tinyclaw
              self.packages.${system}.tinyclaw-core
              self.packages.${system}.dhall-config
              pkgs.cacert
              pkgs.tzdata
            ];
            config = {
              Cmd = [ "/bin/tinyclaw" "gateway" "--verified" ];
              Env = [
                "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
                "TZDIR=${pkgs.tzdata}/share/zoneinfo"
              ];
              ExposedPorts = {
                "18790/tcp" = {};
              };
            };
          };

          default = self.packages.${system}.tinyclaw;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go gateway
            go
            golangci-lint
            goreleaser

            # Dhall config
            dhall
            dhall-json
            dhall-lsp-server

            # F* / OCaml (verified core)
            ocaml
            dune_3
            ocamlPackages.findlib
            ocamlPackages.yojson

            # Futhark (parallel compute kernels)
            futhark

            # Build system
            just
            jq

            # Nix tools
            direnv
            nix-direnv
          ];

          shellHook = ''
            echo "tinyclaw dev shell"
            echo "  just --list    # available targets"
          '';
        };

        # Flake checks
        checks = {
          dhall-typecheck = pkgs.stdenv.mkDerivation {
            pname = "tinyclaw-dhall-check";
            inherit version;
            src = ./dhall;
            nativeBuildInputs = with pkgs; [ dhall dhall-json ];
            buildPhase = ''
              find . -name '*.dhall' -exec dhall type --file {} \; > /dev/null
            '';
            installPhase = "mkdir -p $out && touch $out/ok";
          };

          go-tests = pkgs.stdenv.mkDerivation {
            pname = "tinyclaw-go-tests";
            inherit version;
            src = ./.;
            nativeBuildInputs = with pkgs; [ go ];
            buildPhase = ''
              export HOME=$TMPDIR
              export GOFLAGS="-tags=stdjson"
              go generate ./...
              go test ./...
            '';
            installPhase = "mkdir -p $out && touch $out/ok";
          };
        };
      }
    );
}
