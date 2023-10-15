{
  description = "NixOS configurations for puercopop.com and hiippo.com";
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
      blog = pkgs.callPackage ./blog { };
      www = pkgs.callPackage ./www { };
      # TODO: Use callPackage
      ops = import ./ops {
        system = system;
        nixpkgs = nixpkgs;
        pkgs = pkgs;
        www = www;
        blog = blog;
      };
    in
    {
      formatter."${system}" = pkgs.nixpkgs-fmt;
      nixosConfigurations = {
        kraken = ops.kraken;
        sparrow = ops.sparrow;
      };
    };
}
