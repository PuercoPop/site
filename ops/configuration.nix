{ config, pkgs, ... }:

{
  imports = [ <nixpkgs/nixos/modules/installer/cd-dvd/installation-cd-minimal.nix> ];

  time.timeZone = "America/Lima";
  networking.hostName = "kraken";
  services.openssh = {
    enable = true;
    ports = [ 6969 ];
    settings = {
      PermitRootLogin = "no";
      PasswordAuthentication = false;
    };
  };

  security = {
    sudo = {
      enable = true;
      wheelNeedsPassword = false;
    };
  };

  users = {
    mutableUsers = false;
    users.nixos = {
      isNormalUser = true;
      initialPassword = "changeme";
      extraGroups = [ "wheel" ];
      openssh.authorizedKeys.keys = [
        (builtins.fetchurl "https://github.com/PuercoPop.keys")
      ];
    };
  };
}
