{
  description = "NixOS configurations for puercopop.com and hiippo.com";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-23.05";
    crane = {
      url = "github:ipetkov/crane";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };
  outputs = { self, nixpkgs, crane }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs {
        inherit system;
        overlays = [ ];
        config = { };
      };
      crane-pkgs = import crane {
        pkgs = pkgs;
      };
      blog = pkgs.callPackage ./blog { pkgs = pkgs; crane = crane-pkgs; };
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
