{
  description = "Provision and deploy `site'";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };
  outputs = { self, nixpkgs}:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
      iou = "setup nixops";
      terraform = pkgs.writeShellScriptBin "terraform" ''
      export TF_VAR_VULTR_API_KEY=${iou}
      ${pkgs.terraform}/bin/terraform $@
      '';
    in
      {
        packages = {
          terraform = terraform;
        }
        devShells.${system}.default = pkgs.mkShell {
          buildInputs = [
            terraform
            pkgs.terraform-providers.vultr
            pkgs.terraform-ls
            pkgs.nixos-rebuild
          ];
        };
      };
}
