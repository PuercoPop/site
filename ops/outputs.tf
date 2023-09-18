output "kraken_ip" {
  value = digitalocean_droplet.kraken.ipv4_address
  description = "The public IP of the Kraken droplet"
}

output "droplet_region" {
  value = digitalocean_droplet.kraken.region
  description = "The region on which it is deployed"
}


output "kraken_price" {
  value = digitalocean_droplet.kraken.price_monthly
  description = "The month price in USDs"
}

output "dns_puercopop_ttl" {
  value = digitalocean_domain.dns_puercopop.ttl
}
