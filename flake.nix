{
  description = "Send SSE when files change";
  outputs = { self, nixpkgs }: {

    packages.x86_64-linux.ssentr = nixpkgs.legacyPackages.x86_64-linux.buildGoModule {
      name = "ssentr";
      src = self;
      vendorSha256 = "sha256-L1TzJ6gAbMRVuALsGSRioXU6eQdMCHlbkEmeR4qz4Lg=";
      meta = with nixpkgs.legacyPackages.x86_64-linux.lib; {
        maintainers = [ maintainers.tomberek ];
        platforms = platforms.linux;
        license = licenses.mit;
      };
    };
    defaultPackage.x86_64-linux = self.packages.x86_64-linux.ssentr;

    apps.x86_64-linux.update = {
      type = "app";
      program = (with nixpkgs.legacyPackages.x86_64-linux; writeScript "update.sh" ''
        ${nix-prefetch}/bin/nix-prefetch --file 'fetchTarball "channel:nixos-unstable"' \
          "{ sha256 }: (builtins.getFlake \"$PWD\").defaultPackage.x86_64-linux.go-modules.overrideAttrs
                (_: { modSha256 = sha256; })" \
          --option experimental-features 'nix-command flakes'
      '').outPath;
      };
    hydraJobs.ssentr.x86_64-linux = self.defaultPackage.x86_64-linux;
  };
}
