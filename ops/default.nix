{ system ? builtins.currentSystem
, nixpkgs ? fetchTarball {
    url = "https://github.com/NixOS/nixpkgs/archive/refs/tags/23.05.tar.gz";
    sha256 = "10wn0l08j9lgqcw8177nh2ljrnxdrpri7bp0g7nvrsn9rkawvlbf";
  }
, pkgs ? import nixpkgs {
    overlays = [ ];
    config = { };
    inherit system;
  }
}:
let
  terraform = pkgs.terraform.overrideAttrs (oldAttrs: {
    terraform-providers = [ pkgs.terraform-providers.digitalocean ];
  });
  bootstrap-config = { lib, modulesPath, ...}: {
    system.stateVersion = "23.05";
    imports = [ "${modulesPath}/virtualisation/digital-ocean-image.nix" ];
  };
  bootstrap-img = (pkgs.nixos bootstrap-config).digitalOceanImage;
  config = nixpkgs.lib.nixosSystem {
    system = "${system}";
    modules = [ ./kraken-configuration.nix ];
  }
in
{
  conf = conf;
  bootstrap-img = bootstrap-img;
  shell = pkgs.mkShell {
    buildInputs = [
      terraform
      pkgs.terraform-ls
      pkgs.nixos-rebuild
      # TODO: Figure out how to use it
      pkgs.age
      pkgs.agenix
    ];
  };
}
