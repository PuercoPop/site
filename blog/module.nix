{ lib, config, options, ... }:
let
  cfg = config.kraken.services.blog;
in
{
  options.services.blog = {
    enable = lib.mkEnableOption "Enable blog";
    package = lib.mkOption {
      type = lib.types.package;
      description = "The blog derivation to use";
    };
    options.user = lib.mkOption {
      default = "blog";
      type = lib.types.str;
      description = "The user to run the blog service as";
    };
    options.dbname = lib.mkOption {
      type = lib.types.str;
      description = "The database name blog is deployed to.";
    };
  };

  config = lib.mkIf cfg.enable {
    environment.systemPackages = [ cfg.package ];
    users = {
      users.${cfg.user} =
        {
          isSystemuser = true;
          group = "${cfg.user}";
        };
      group.${cfg.user} = { };
    };
    # TODO: Can I define a postgresql service here as well; Yes. Options are merged

    systemd.services = {
      postgresql = {
        enable = true;
        ensureDatabases = [ "blog" ];
      };

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
          User = "blog";
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
          User = "blog";
          Group = "blog";
          Restart = "always";
          ExecStart = ''${cfg.package}/bin/serve-blog'';
        };
      };
    };
  };
}
