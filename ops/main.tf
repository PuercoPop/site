terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
    }

    http = {
      source = "hashicorp/http"
    }
  }
}

data "http" "github_keys" {
  url = "https://github.com/PuercoPop.keys"
}

provider "digitalocean" {
  token = var.do_token
  spaces_access_id = var.do_spaces_access_id
  spaces_secret_key = var.do_spaces_secret_key
}

resource "digitalocean_ssh_key" "github" {
  name = "github_keys"
  public_key = chomp(data.http.github_keys.response_body)
  lifecycle {
    precondition {
      condition     = contains([201, 200, 204], data.http.github_keys.status_code)
      error_message = "Could not fetch SSH keys from GitHub"
    }
  }
}

# resource "digitalocean_spaces_bucket" "custom-isos" {
#   name = "custom-isos"
#   region = "nyc3"
# }

# resource "digitalocean_spaces_bucket_object" "bootstrap-nix-iso" {
#   key = "nix.iso"
#   bucket = digitalocean_spaces_bucket.custom-isos.name
#   source = "./nix.iso"
#   # content = file("./nix.iso")
#   acl = "public-read"
#   etag = filemd5("./nix.iso")
#   region = "nyc3"
# }

resource "digitalocean_custom_image" "nixos-image" {
  distribution = "Unknown OS"
  name = "nixos-23.05"
  # url = "https://${digitalocean_spaces_bucket.custom-isos.bucket_domain_name}/${digitalocean_spaces_bucket_object.bootstrap-nix-iso.key}"
  url = "https://custom-isos.nyc3.digitaloceanspaces.com/nix.iso"
  regions = ["nyc3"]
}

resource "digitalocean_droplet" "kraken" {
  image = digitalocean_custom_image.nixos-image.id
  name = "kraken"
  size = "s-1vcpu-1gb"
  region = "nyc3"
  ipv6 = false # Custom images cannot use ipv6
  ssh_keys = [digitalocean_ssh_key.github.fingerprint]
  tags = ["kraken"]

}

resource "digitalocean_domain" "dns_puercopop" {
  name = "puercopop.com"
  ip_address = digitalocean_droplet.kraken.ipv4_address
}

resource "digitalocean_record" "www_subdomain" {
  domain = digitalocean_domain.dns_puercopop.id
  type = "A"
  name = "www"
  value = digitalocean_droplet.kraken.ipv4_address
}
