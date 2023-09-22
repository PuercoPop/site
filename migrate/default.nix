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
  migrate = pkgs.buildGoModule {
    pname = "migrate";
    version = "0.0.1";
    src = ./.;
    vendorHash = "sha256-I6dyihp9emHL1FR6CX8aLw4nWhwahnqOmGTKD/T6IG8=";
  };
in
{
  default = migrate;
  shell = pkgs.mkShell {
    nativeBuildInputs = [pkgs.postgresql_15 ];
  };
}
