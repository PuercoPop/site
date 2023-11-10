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
pkgs.nixosTest {
  name = "migrate tests";
  nodes.machine = {
    users.users.nixos = {
      isNormalUser = true;
      initialHashedPassword = "nixos";
    };
    services = {
      # /home/puercopop/src/nixpkgs/nixos/tests/web-apps/mastodon/remote-postgresql.nix
      postgresql = {
        enable = true;
      };
    };
  };
  testScript = ''
    machine.wait_for_unit("postgresql.service")
    machine.succeed("sudo -u nixos whoami")
    machine.succeed("sudo -u postgres psql -c 'SELECT 1;'")
  '';
}
