{ lib, config, options, ... }:
{
  options.services.blog = {
    enable = lib.mkEnableOption "Enable blog";
    # TODO: Should we make this an option?
    package = lib.mkOption {
      default = pkgs.callPackage ./default.nix { };
      type = lib.types.package;
      description = "The blog derivation to use";
    };
    options.user = lib.mkOption {
      default = "blog";
      type = types.str;
      description = "The user to run the blog service as";
    };
    options.dbname = lib.mkOption {
      type = lib.types.str;
      description = "The database name blog is deployed to.";
    };
  };

  config = mkIf cfg.enable {
    environment.systemPackages = [ cfg.package ];
    users = {
      users.${cfg.user} =
        {
          isSystemuser = true;
          group = "${cfg.user}";
        };
      group.${cfg.user} = { };
    };
    # TODO: Can I define a postgresql service here as well;

    systemd.services = {
      # https://github.com/serokell/systemd-nix
      # We need to define 3 systemd-units
      # 1. To run psql -f sql/schema.sql
      # 2. To run import-blog
      # 3. To run serve-blog
      import-blog-svc = {
        description = "Import all the posts";
        wantedBy = [ "multi-user.target" ];
        after = [
          "network.target"
          "postgresql.service"
        ];
        serviceConfig = {
          User = "blog";
          # man 7 systemd.directives
          # Type = "notify";
          ExecStart = ''${blog}/bin/import-blog'';
        };
      };
      serve-blog-svc = {
        description = "HTTP Server for the blog";
        wantedBy = [ "multi-user.target" ];
        wants = [ "network-online.target" ]; # TODO: Is this necessary?
        after = [
          "network.target"
          "postgresql.service"
        ];
        serviceConfig = {
          User = "blog";
          Group = "blog";
          Restart = "always";
          ExecStart = ''${blog}/bin/serve-blog'';
        };
      };
    };
  };
}
