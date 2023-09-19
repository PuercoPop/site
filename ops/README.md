# OPS
## Provision and deploy puercopop.com

terraform for infrastructure provisioning nixos-rebuild for configuration.

# Overview

We store our secrets using sops via nix-ops. We make them available to Terraform
through a wrapper that sets environment variables before calling
`terraform`. Finally we use nix configurations to configure any machines that
are created using `nixos-rebuild`

# Runbooks

## Build boostrap ISO

```shell
nix-build -A bootstrap-img
cp result/nixos.qcow2.gz nix.iso
```

## Provision a new VM

```shell
nix-shell -A shell
terraform apply
nixos-rebuild switch --fast --flake ..#kraken --build-host root@puercopop.com --target-host root@puercopop.com
```

## Deploy to a VM

```shell
nix-shell -A shell
# This command needs to be updated
nixos-rebuild switch --fast --flake .#default --target-host root@puercopop.com --build-host root@puercopop.com
```

# References

- https://www.haskellforall.com/2023/01/announcing-nixos-rebuild-new-deployment.html?m=1
- https://jonascarpay.com/posts/2022-09-19-declarative-deployment.html#writing-and-building-the-application.

[cloud-init]: https://cloudinit.readthedocs.io/en/latest/
