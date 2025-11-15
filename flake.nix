{
  description = "Minimal dev shell with Go and g++";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.05";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }: 
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.gcc
          ];

          shellHook = ''
            go version | cut -d' ' -f3
            ${pkgs.gcc}/bin/g++ --version | head -n1
          '';
        };
      }
    );
}

