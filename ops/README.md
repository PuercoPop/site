# OPS

terraform for infrastructure provisioning nixos-rebuild for configuration.

# Overview

We store our secrets using sops via nix-ops. We make them available to Terraform
through a wrapper that sets environment variables before calling
`terraform`. Finally we use nix configurations to configure any machines that
are created using `nixos-rebuild`

# Provisioning a machine

vCPU/s: 1 vCPU
RAM: 1024.00 MB
Storage: 25 GB SSD
Bandwidth: 0 GB
$ mkdir -p ~/.ssh
$ curl -L https://github.com/PuercoPop.keys >~/.ssh/authorized_keys
$ scp configuration.nix nixos@$IP:/etc/configuration.nix

# References

- https://www.haskellforall.com/2023/01/announcing-nixos-rebuild-new-deployment.html?m=1
- https://jonascarpay.com/posts/2022-09-19-declarative-deployment.html#writing-and-building-the-application.
