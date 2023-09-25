{ modulesPath, ... }: {
  # TODO: Check if I can use qemu-guest.nix
  imports = [ "${modulesPath}/virtualisation/digital-ocean-image.nix" ];

  system.stateVersion = "23.05";
  time.timeZone = "America/Lima";
  networking.hostName = "sparrow";

  security = {
    sudo = {
      enable = true;
      wheelNeedsPassword = false;
    };

    acme = {
      acceptTerms = true;
      defaults.email = "pirata+hiippo@gmail.com";
    };
  };

  users.users.nixos = {
    isNormalUser = true;
    extraGroups = [ "wheel" "networkmanager" "video" ];
    # Allow the graphical user to login without password
    initialHashedPassword = "";
    openssh.authorizedKeys.keys = [
      (builtins.fetchurl {
        url = "https://github.com/PuercoPop.keys";
        sha256 = "012vifcnrnkw3jb31scgn2n53qf68cw03qmy2k5n3h2i8lla24f4";
      })
    ];
  };

  networking.firewall.allowedTCPPorts = [ 80 443 ];
  services = {
    jitsi-meet = {
      enable = true;
      hostName = "meet.hiippo.com";
    };
  };
}
