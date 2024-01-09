# migrate
## A barebones roll-forward database migration tool

Runs any pending migrations and exists.

# Principles

Give you the basic machinery but trust the developer knows what they are
doing. Keep things simple.

# Features

- Migrations only roll forward.
- The author of the migration is responsible for choosing and setting up the
  transactions they want to use.
- `migrate` keeps track of which migrations have been run against the
  database. Running the same version of `migrate` multiple times should be
  idempotent.
- Uses `psql` to run the migrations.

# Anti-Features

- Bundling the migrations. Copying the the migrations to the production
  environment is left up to the user. I use nixos-rebuild to take care of that


# FQA

## How do I run the migrations on deploy on NixOS?

Defining a Systemd oneshot unit. ⩰:

```nix
migrate-unit = {
  wantedBy = [ "multi-user.target" ];
  after = [
    "network.target"
    "postgresql.service"
  ];
  serviceConfig = {
    Type = "oneshot";
    ExecStart = "${pkgs.postgresql}/bin/psql" -d ${cfg.migrationDir} -D ${cfg.dburl}
  };
}

```

# Related approaches

- https://github.com/purcell/postgresql-migrations/
- https://nalanj.dev/posts/minimalist-postgresql-migrations/
- https://github.com/docteurklein/declarative-sql-migrations/#run-tests
