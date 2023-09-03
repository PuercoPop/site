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
      terraform = pkgs.terraform;
      bootstrap-config-module = { lib, modulesPath, ...}: {
        system.stateVersion = "23.05";
        # imports = [ <nixpkgs/nixos/modules/installer/cd-dvd/installation-cd-minimal.nix> ];
        # imports = [ "${nixpkgs}/nixos/modules/installer/cd-dvd/installation-cd-minimal.nix" ];
        imports = [ "${modulesPath}/installer/cd-dvd/installation-cd-minimal.nix" ];
        # imports = [ "${modulesPath}/virtualization/qemu-vm.nix" ];

        services.openssh = {
          enable = true;
          settings = {
            PermitRootLogin = lib.mkForce "no";
            PasswordAuthentication = false;
          };
        };
          users = {
            mutableUsers = false;
            users.root.openssh.authorizedKeys.keys = [
              "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKi6ih3rTLCwqlQnyOQHqyIUWHh8ipHLrFmjNH4rB5yP"
              # (builtins.fetchurl "https://github.com/PuercoPop.keys")
            ];
            users.nixos = {
              isNormalUser = true;
              initialPassword = "";
              extraGroups = [ "wheel" ];
              openssh.authorizedKeys.keys = [
                "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKi6ih3rTLCwqlQnyOQHqyIUWHh8ipHLrFmjNH4rB5yP"
                # (builtins.fetchurl "https://github.com/PuercoPop.keys")
              ];
              
            };
          };
      };
      bootstrap-img-name = "nixos-bootstrap-${system}";
      bootstrap-img = nixos-generators.nixosGenerate {
        format = "iso";
        system = system;
        modules = [
          bootstrap-config-module
        ];
      };
    in
      {
        formatter."${system}" = pkgs.nixfmt;
        packages = {
          "${system}" = {
            bootstrap-img = bootstrap-img;
            # nix run .#terraform
            terraform = terraform;
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
