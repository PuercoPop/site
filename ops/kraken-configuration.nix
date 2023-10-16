{ modulesPath, pkgs, www, blog, ... }: {
  # TODO: Check if I can use qemu-guest.nix
  imports = [
    "${modulesPath}/virtualisation/digital-ocean-image.nix"
    ../blog/module.nix
  ];

  system.stateVersion = "23.05";
  time.timeZone = "America/Lima";
  networking.hostName = "kraken";

  environment.systemPackages = [ blog.default ];

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

  users.users = {
    # TODO: Add a user for each service
    nixos = {
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
  };

  networking.firewall.allowedTCPPorts = [ 80 443 ];
  services = {
    nginx = {
      enable = true;
      virtualHosts = {
        "puercopop.com" = {
          forceSSL = true;
          enableACME = true;
          globalRedirect = "www.puercopop.com";
        };
        "www.puercopop.com" = {
          forceSSL = true;
          enableACME = true;
          locations."/" = {
            root = www.resources;
          };
        };
      };
    };
    blog = {
      enable = true;
      package = blog.default;
      templateDir = blog.templates;
      contentDir = blog.content;
      dbSchema = blog.schema;
      postgresql = {
        package = pkgs.postgresql_15;
      };
    };
  };
}
