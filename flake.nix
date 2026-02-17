{
  description = "Parse code fences from anywhere";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    systems.url = "github:nix-systems/default";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.inputs.systems.follows = "flake-parts";
    };

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;
      imports = [ inputs.treefmt-nix.flakeModule ];

      perSystem =
        {
          inputs',
          pkgs,
          lib,
          ...
        }:
        let
          inherit (inputs'.gomod2nix.legacyPackages) gomod2nix mkGoEnv buildGoApplication;

          goEnv = mkGoEnv { pwd = ./.; };

          version = "0.0.5";
          fenced = buildGoApplication {
            pname = "fenced";
            inherit version;

            modules = ./gomod2nix.toml;
            src = lib.cleanSource ./.;

            ldFlags = [
              "-X github.com/unstoppablemango/fenced/cmd.Version=${version}"
            ];

            nativeBuildInputs = [ pkgs.ginkgo ];

            checkPhase = ''
              ginkgo -r -v --race --trace
            '';

            meta = {
              description = "Parse code fences from anywhere";
              homepage = "https://github.com/UnstoppableMango/fenced";
              license = lib.licenses.mit;
              maintainers = with lib.maintainers; [ UnstoppableMango ];
              mainProgram = "fenced";
            };
          };

          ctr = pkgs.dockerTools.buildLayeredImage {
            name = "fenced";
            tag = version;

            contents = pkgs.buildEnv {
              name = "image-root";
              paths = [ fenced ];
              pathsToLink = [ "/bin" ];
            };

            config = {
              Workdir = "/data";
              Entrypoint = [ "/bin/fenced" ];
            };
          };
        in
        {
          packages = {
            inherit fenced ctr;
            default = fenced;
          };

          apps = {
            default = {
              type = "app";
              program = lib.getExe fenced;
              meta = fenced.meta;
            };

            version = {
              type = "app";
              program = lib.getExe (pkgs.writeShellScriptBin "version" ''
                echo -n "${version}"
              '');
            };
          };

          devShells.default = pkgs.mkShellNoCC {
            packages = with pkgs; [
              bashInteractive
              docker
              ginkgo
              go
              goEnv
              golangci-lint
              gomod2nix
              gopls
              goreleaser
              nil
              nix
              nixfmt
              podman
              watchexec
            ];

            DOCKER = "${pkgs.docker}/bin/docker";
            GINKGO = "${pkgs.ginkgo}/bin/ginkgo";
            GO = "${pkgs.go}/bin/go";
            GOMOD2NIX = "${gomod2nix}/bin/gomod2nix";
            GOLANGCI = "${pkgs.golangci-lint}/bin/golangci-lint";
            GORELEASER = "${pkgs.goreleaser}/bin/goreleaser";
            NIX = "${pkgs.nix}/bin/nix";
            PODMAN = "${pkgs.podman}/bin/podman";
            WATCHEXEC = "${pkgs.watchexec}/bin/watchexec";
          };

          treefmt = {
            programs.gofmt.enable = true;
            programs.nixfmt.enable = true;
          };
        };
    };
}
