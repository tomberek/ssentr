{
  description = "Send SSE when files change";

  outputs = {
    self,
    nixpkgs,
  }: let
    systems = {
      "x86_64-linux" = {};
      "x86_64-darwin" = {};
      "aarch64-linux" = {};
      "aarch64-darwin" = {};
    };
    generate = arg: builtins.mapAttrs arg (builtins.intersectAttrs systems nixpkgs.legacyPackages);
  in rec {
    packages = generate (
      system: pkgs:
        rec {
          ssentr = pkgs.callPackage blueprints.ssentr {};
          default = ssentr;
        }
        // (
          if builtins ? fetchClosure
          then {
            ssentr-baked = builtins.fetchClosure {
              fromPath = "/nix/store/sazl2bwfrj9sm5qm6qxiqag5ghmd24dp-ssentr";
              fromStore = "https://cache.floxdev.com";
            };
          }
          else {}
        )
    );

    blueprints.ssentr = {
      lib,
      buildGo118Module,
    }:
      buildGo118Module {
        name = "ssentr";

        # Ignore nix file changes during build
        src = builtins.path {
          filter = path: type: ! builtins.elem (builtins.baseNameOf path) ["flake.nix" "flake.lock"];
          path = self;
          name = "ssentr";
        };
        vendorSha256 = "sha256-3d5iPPz6iccXq1kJyp6IgyQBGlKI0yZUKZIedeuDzz8=";
        meta = with lib; {
          maintainers = [maintainers.tomberek];
          platforms = builtins.attrNames systems;
          license = licenses.mit;
        };
      };

    apps = generate (system: pkgs: {
      update = {
        type = "app";
        program =
          (with pkgs;
            writeScript "update.sh" ''
              ${nix-prefetch}/bin/nix-prefetch --file 'fetchTarball "channel:nixos-unstable"' \
                "{ sha256 }: (builtins.getFlake \"$PWD\").defaultPackage.x86_64-linux.go-modules.overrideAttrs
                      (_: { modSha256 = sha256; })" \
                --option experimental-features 'nix-command flakes'
            '')
          .outPath;
      };
    });

    hydraJobs = generate (system: pkgs: {ssentr = self.packages.${system}.ssentr;});
  };
}
