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

resource "vultr_iso_private" "nix_iso" {
  url = "https://channels.nixos.org/nixos-23.05/latest-nixos-minimal-x86_64-linux.iso"
}

# resource "vultr_block_storage" "skull-island" {}

resource "vultr_instance" "kraken" {
  iso_id      = "nix_iso"
  enable_ipv6 = true
  plan        = "vc2-1c-1gb"
  region      = "seattle"
  hostname    = "kraken"
}
