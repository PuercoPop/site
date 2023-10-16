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
    templateDir = lib.mkOption {
      type = lib.types.path;
      description = "A directory containing the templates";
    };
    contentDir = lib.mkOption {
      type = lib.types.path;
      description = "A directory containing the posts";
    };
    dbSchema = lib.mkOption {
      type = lib.types.path;
      description = ''A file containing the database schema.
        It should be idempotent.'';
    };
    user = lib.mkOption {
      default = "blog";
      type = lib.types.str;
      description = "The user to run the blog service
as";
    };
    group = lib.mkOption {
      default = "blog";
      type = lib.types.str;
      description = "The group to run the blog service as";
    };
    postgresql = {
      package = lib.mkOption {
        type = lib.types.package;
        description = "The PostgreSQL package to use.";
      };
      dburl = lib.mkOption {
        default = "postgresql:///${cfg.postgresql.dbname}?user=${cfg.user}&host=/run/postgresql";
        description = "The connection string to use.
           See https://www.postgresql.org/docs/16/libpq-connect.html#LIBPQ-CONNSTRING";
      };
      dbname = lib.mkOption {
        default = "blog";
        type = lib.types.str;
        description = "The database to use";
      };
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
    services = {
      postgresql = {
        enable = true;
        package = cfg.postgresql.package;
        enableTCPIP = false;
        ensureDatabases = [ cfg.postgresql.dbname ];
        ensureUsers = [{
          name = cfg.user;
          ensurePermissions = {
            "DATABASE ${cfg.postgresql.dbname}" = "ALL PRIVILEGES";
          };
        }];
      };
      nginx = {
        enable = true;
        virtualHosts = {
          "blog.puercopop.com" = {
            forceSSL = true;
            enableACME = true;
            locations."/" = {
              proxyPass = "http://127.0.0.1:3000";
              # TODO: Set the X-Forwarded-For header
            };
          };
        };
      };
    };

    systemd.services = {
      # We need to do 3 things
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
          User = cfg.user;
          Type = "oneshot";
          ExecStart = [
            "${cfg.postgresql.package}/bin/psql -f ${cfg.dbSchema} -d ${cfg.postgresql.dburl}"
            "${cfg.package}/bin/import-blog -d ${cfg.postgresql.dburl} -D ${cfg.contentDir}"
          ];
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
          ExecStart = ''${cfg.package}/bin/serve-blog -d "${cfg.postgresql.dburl}" -D ${cfg.templateDir}'';
        };
      };
    };
  };
}
