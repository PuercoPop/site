terraform {
  required_providers {
    vultr = {
      source  = "vultr/vultr"
      version = "2.15.1"
    }
  }
}

provider "vultr" {
  # Ensure VULTR_API_KEY is set
  # api_key = ""
}

data "http" "github" {
  url = "https://github.com/PuercoPop.keys"
}

resource "vultr_ssh_key" "github" {
  name = "github"
  ssh_key = chomp(data.htttp.github.response_body)
}

resource "vultr_iso_private" "nix_iso" {
  url = "https://channels.nixos.org/nixos-23.05/latest-nixos-minimal-x86_64-linux.iso"
}

# resource "vultr_block_storage" "skull-island" {}

resource "vultr_instance" "kraken" {
  # TODO(javier): Show do I say the equivalent of nix_iso.id
  # iso_id      = vultr_iso_private.nix_iso.id
  iso_id      = "62950b1e-61e8-4020-ac4a-31c9f12b86ac"
  enable_ipv6 = true
  plan        = "vc2-1c-1gb"
  region      = "sea"
  hostname    = "kraken"
}

# How do I look at the data source
# https://registry.terraform.io/providers/vultr/vultr/latest/docs/data-sources/plan
