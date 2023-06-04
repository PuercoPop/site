{
  description = "Provision and deploy `site'";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };
  outputs = { self, nixpkgs}:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
      {
        devShells.${system}.default = pkgs.mkShell {
          buildInputs = [
            pkgs.terraform
            pkgs.terraform-providers.vultr
            pkgs.terraform-ls
            pkgs.nixos-rebuild
          ];
        };
      };
}
