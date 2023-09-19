{ modulesPath, www,... }: {
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
      defaults.email = "pirata+puercopop@gmail.com";
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
        sha256 = "1q50qkjvcn57c5kn0al3097hisv7193jxzs2q1a4cg8hxxyicrp1";
      })
    ];
  };
  users.users.nginx.extraGroups = [ "acme" ];

  networking.firewall.allowedTCPPorts = [ 80 443 ];
  services = {
    # https://discourse.nixos.org/t/nixos-nginx-acme-ssl-certificates-for-multiple-domains/19608/3
    nginx = {
      enable = true;
      virtualHosts = {
        "www.puercopop.com" = {
          forceSSL = true;
          enableACME = true;
          locations."/" = {
            root = www.resources;
          };
        };
      };
    };
  };
}
