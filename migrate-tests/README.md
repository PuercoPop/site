# Integration tests for `migrate`

## Run tests

```shell
nix-build
```

## Scenarios to test

1. Installing the migration schema
  a. If the DB is empty the migration schema is installed.
  b. If the DB contains the migration schema. The script runs successfully, with
     or without any migrations to run
2. A migration fails.
  a. Check the exit code is zero
  b. Check that previous migrations where run
  c. Check that subsequent migrations did not run.
3. Happy path
  a. Check that the migration table is up to date
  b. Check that the schema corresponds to what we want.


## Implementation Ideas

### NixOS tests

Nix Hour #20: https://www.youtube.com/watch?v=RgKl8Jue4qM

### dockertest

Or a podman API equivalent
