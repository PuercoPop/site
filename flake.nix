{
  description = "";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-23.05";
  };
  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs {
        inherit system;
        overlays = [ ];
        config = { };
      };
      ops = import ./ops {
        nixpkgs = nixpkgs;
      };
    in
    {
      formatter."${system}" = pkgs.nixpkgs-fmt;
      nixosConfigurations = {
        kraken = ops.config;
      };
      ops = ops;
    };
}
