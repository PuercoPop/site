{
  description = "Provision and deploy `site'";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-23.05";
    nixos-generators= {
      url = "github:nix-community/nixos-generators";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };
  outputs = { self, nixpkgs, nixos-generators }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
      # iou = "setup nixops";
      # terraform = pkgs.writeShellScriptBin "terraform" ''
      #   export TF_VAR_VULTR_API_KEY=${iou}
      #   ${pkgs.terraform}/bin/terraform $@
      # '';
      bootstrap-config-module = {
        system.stateVersion = "23.05";
        services.openssh.enable = true;
        users.users.root.openssh.authorized.keys = [
          "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKi6ih3rTLCwqlQnyOQHqyIUWHh8ipHLrFmjNH4rB5yP"
        ];
      };
      bootstrap-img-name = "nixos-bootstrap-${system}";
      bootstrap-img = nixos-generators.nixosGenerate {
        format = "qcow"; # or raw
        modules = [
          bootstrap-config-module
        ];
      };
    in
      {
        packages = {
          # # nix run .#terraform
          # terraform = terraform;
          x86_64-linux = {
            bootstrap-img = bootstrap-img;
          };
        };
        devShells.${system}.default = pkgs.mkShell {
          buildInputs = [
            # terraform
            pkgs.terraform
            pkgs.terraform-providers.vultr
            pkgs.terraform-ls
            pkgs.nixos-rebuild
            pkgs.age
          ];
        };
      };
}
