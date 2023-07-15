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

resource "vultr_block_storage" "skull-island" {


}

# resource "vultr_instance" "kraken" {
#   enable_ipv6 = true
# }
