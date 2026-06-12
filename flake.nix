{
  description = "kavka";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        name = "kavka-dev-shell";

        buildInputs = with pkgs; [
          go
          gopls # Language server
          delve # Debugger
          golangci-lint # Linter
          glibc
          glibc.static
          gcc # Required for CGO
          pkg-config # For finding dependencies
        ];

        shellHook = ''
          echo "Go development environment ready"
          echo "Go version: $(go version)"
          export CGO_ENABLED=1
        '';
      };
    };
}
