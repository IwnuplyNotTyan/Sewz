{
  description = "🌑 I hate docker >:3";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs, ... }:
    let
      systems = [
        "x86_64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
	  version = "0.1.0";
        in
        {
          default = pkgs.buildGoModule {
            pname = "sewz";
            inherit version;
            src = self;
            modules = ./gomod2nix.toml;

            ldflags = [
	      "-s"
	      "-w"
            ];

            vendorHash = "sha256-N7imleQuAXMp5IyjkpyERNe1X9BuncOsGV8reR/G7o8=";

            meta = {
              description = "Pretty docker wrapper";
              homepage = "https://github.com/IwnuplyNotTyan/sewz";
              mainProgram = "sewz";
            };
          };
	});
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go
              gopls
              gotools
              golangci-lint
            ];
          };
        });
    };
}
