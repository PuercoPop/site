# OPS

terraform for infrastructure provisioning
nixos-rebuild for configuration (https://www.haskellforall.com/2023/01/announcing-nixos-rebuild-new-deployment.html?m=1)

# Overview

We store our secrets using sops via nix-ops. We make them available to Terraform
through a wrapper that sets environment variables before calling
`terraform`. Finally we use nix configurations to configure any machines that
are created using `nixos-rebuild`.
