{ lib, config, options, ... }:
let
  cfg = config.services.blog;
in
{
  options.services.blog = {
    enable = lib.mkEnableOption "Enable blog";
    package = lib.mkOption {
      type = lib.types.package;
      description = "The blog derivation to use";
    };
    user = lib.mkOption {
      default = "blog";
      type = lib.types.str;
      description = "The user to run the blog service as";
    };
    group = lib.mkOption {
      default = "blog";
      type = lib.types.str;
      description = "The group to run the blog service as";
    };
    dbname = lib.mkOption {
      type = lib.types.str;
      description = "The database name blog is deployed to.";
    };
  };

  config = lib.mkIf cfg.enable {
    environment.systemPackages = [ cfg.package ];
    users = {
      users.${cfg.user} =
        {
          isSystemUser = true;
          group = "${cfg.user}";
        };
      groups.${cfg.group} = { };
    };
    # TODO: Can I define a postgresql service here as well; Yes. Options are merged

    services = {
      postgresql = {
        enable = true;
        ensureDatabases = [ cfg.dbname ];
        ensureUsers = [ cfg.user ];
        # ensurePermissions = { };
      };
    };

    systemd.services = {
      # https://github.com/serokell/systemd-nix
      # We need to define 3 systemd-units
      # 1. To run psql -f sql/schema.sql
      # 2. To run import-blog
      # 3. To run serve-blog
      #
      # blog-schema = pkgs.writeShellScriptBin "blog-schema" ''
      #   psql -f ${./sql/schema.sql} -d adress
      # '';
      import-blog-svc = {
        description = "Import all the posts";
        wantedBy = [ "multi-user.target" ];
        after = [
          "network.target"
          "postgresql.service"
        ];
        serviceConfig = {
          User = cfg.services.blog.user;
          # man 7 systemd.directives
          # Type = "notify";
          ExecStart = ''${cfg.package}/bin/import-blog'';
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
          User = cfg.user;
          Group = cfg.group;
          Restart = "always";
          ExecStart = ''${cfg.package}/bin/serve-blog'';
        };
      };
    };
  };
}
