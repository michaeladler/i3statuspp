{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    utils.url = "github:numtide/flake-utils";
    naersk = {
      url = "github:nmattia/naersk";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, utils, naersk }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages."${system}";
        naersk-lib = naersk.lib."${system}";
      in
      rec {
        # `nix build`
        packages.i3statuspp = naersk-lib.buildPackage {
          pname = "i3statuspp";
          version = "0.1.0";
          root = ./.;

          meta = with pkgs.lib; {
            description = "Extend i3status functionality";
            homepage = "https://github.com/michaeladler/i3statuspp";
            license = licenses.asl20;
          };
        };
        defaultPackage = packages.i3statuspp;

        # `nix run`
        apps.i3statuspp = utils.lib.mkApp {
          drv = packages.i3statuspp;
        };
        defaultApp = apps.i3statuspp;

        # `nix develop`
        devShell = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [ rustc cargo ];
        };
      });
}
