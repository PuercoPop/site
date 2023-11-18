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
  # https://github.com/NixOS/nixpkgs/blob/master/doc/languages-frameworks/rust.section.md
  migrate = pkgs.rustPlatform.buildRustPackage rec {
    pname = "migrate";
    version = "0.1.0";
    src = ./.;
    cargoHash = "sha256-mftNtxQcab+5Mr8qvUpXbZSTk1ZAg8FPxR5qqZjHqf8=";
  };
in
{ default = migrate; }
