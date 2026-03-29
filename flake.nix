{
  description = "File System Scan - CLI Tool written in GOLANG";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

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

        packages.default = pkgs.stdenv.mkDerivation {
          name = "fsscan";
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