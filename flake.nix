{
  description = "Send SSE when files change";
  outputs = { self, nixpkgs }: {

    packages.x86_64-linux.ssentr = nixpkgs.legacyPackages.x86_64-linux.buildGoModule {
      name = "ssentr";
      src = self;
      vendorSha256 = "sha256-L1TzJ6gAbMRVuALsGSRioXU6eQdMCHlbkEmeR4qz4Lg=";
    };
    defaultPackage.x86_64-linux = self.packages.x86_64-linux.ssentr;
    hydraJobs.ssentr.x86_64-linux = self.defaultPackage.x86_64-linux;
  };
}
