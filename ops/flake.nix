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
        services.openssh = {
          enable = true;
          settings = {
            PermitRootLogin = "no";
            PasswordAuthentication = false;
          };
        };
          users = {
            mutableUsers = false;
            users.root.openssh.authorized.keys = [
              (builtins.fetchurl "https://github.com/PuercoPop.keys")
            ];
            users.nixos = {
              isNormalUser = true;
              initialPassword = "";
              extraGroups = [ "wheel" ];
              openssh.authorizedKeys.keys = [
                (builtins.fetchurl "https://github.com/PuercoPop.keys")
              ];
              
            };
          };
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
        formatter."${system}" = pkgs.nixfmt;
        packages = {
          # # nix run .#terraform
          # terraform = terraform;
          # "${system}" = {
          #   bootstrap-img = bootstrap-img;
          # };
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
