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
    # /home/puercopop/src/nixpkgs/nixos/tests/web-apps/mastodon/remote-postgresql.nix
  };
  testScript = ''
    start_all()
    machine.succeed("sudo -u nixos whoami")
  '';
}
