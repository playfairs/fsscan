{
  description = "fsscan: minimal dev shell + runnable CLI app";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }: 
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.gcc
          ];

          shellHook = ''
            echo "Go version: $(go version | cut -d' ' -f3)"
            echo "G++ version: $(${pkgs.gcc}/bin/g++ --version | head -n1)"
          '';
        };

        packages.default = pkgs.stdenv.mkDerivation {
          pname = "fsscan";
          version = "1.0.0";
          src = ./.;

          buildInputs = [ pkgs.bash ];

          installPhase = ''
            mkdir -p $out/bin
            cp fsscan.sh $out/bin/fsscan
            chmod +x $out/bin/fsscan
          '';
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/fsscan";
        };
      }
    );
}