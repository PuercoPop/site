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
  # api_key = var.vultur_api_key
}

data "http" "github" {
  url = "https://github.com/PuercoPop.keys"
}

resource "vultr_ssh_key" "github" {
  name = "github"
  ssh_key = chomp(data.http.github.response_body)
}

resource "vultr_iso_private" "nix_iso" {
  url = "https://channels.nixos.org/nixos-23.05/latest-nixos-minimal-x86_64-linux.iso"
}

# resource "vultr_block_storage" "skull-island" {}

resource "vultr_instance" "kraken" {
  # TODO(javier): Show do I say the equivalent of nix_iso.id
  iso_id      = vultr_iso_private.nix_iso.id
  enable_ipv6 = true
  plan        = "vc2-1c-1gb"
  region      = "sea"
  hostname    = "kraken"
  ssh_key_ids = [vultr_ssh_key.github.id]
  # user_data = file("cloud-init.yml")
}

# How do I look at the data source
# https://registry.terraform.io/providers/vultr/vultr/latest/docs/data-sources/plan
