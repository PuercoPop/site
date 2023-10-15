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
, crane-repo ? fetchTarball {
    url = "https://github.com/ipetkov/crane/archive/refs/tags/v0.13.0.tar.gz";
    sha256 = "0k1qipknmy40wcndrmg4lmm6529k61qyn915k5y79vx2rj3jpj83";
  }
, crane ? import crane-repo { }
}:
let
  src = crane.cleanCargoSource (crane.path ./.);
  commonArgs = {
    inherit src;
    buildInputs = [
      pkgs.pkg-config
      pkgs.openssl
    ];
  };
  cargoArtifacts = crane.buildDepsOnly commonArgs;
  testDataFilter = path: _type: builtins.match ("testdata/*.md" path) != null;
  blog = crane.buildPackage (commonArgs // {
    inherit cargoArtifacts;
    src = pkgs.lib.cleanSourceWith {
      src = crane.path ./.;
      filter = path: type: (testDataFilter path type) || (crane.filterCargoSources path type);
    };
  });
in
{
  default = blog;
  schema = ./src/schema.sql;
  templates = ./templates;
  shell = pkgs.mkShell {
    packages = [
      pkgs.cargo
      pkgs.cargo-watch
      pkgs.rustc
      pkgs.rustfmt
      pkgs.clippy
      pkgs.rust-analyzer
      pkgs.sqitchPg
    ];
  };
}
