{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    mach-nix.url = "github:davhau/mach-nix";
  };

  outputs = { self, nixpkgs, mach-nix }:
    let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
      mach = mach-nix.lib.x86_64-linux;
      pythonVersion = "python311";
      pythonEnv = mach.mkPython {
        python = pythonVersion;
        requirements = builtins.readFile ./requirements.txt;
      };
    in
    {
      devShells.x86_64-linux.default = pkgs.mkShellNoCC
        {
          packages = [ pythonEnv ];
          shellHook = ''
            export PYTHONPATH="${pythonEnv}/bin/python"
          '';
        };
    };
}
