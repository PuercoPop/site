
output "kraken_ip" {
  value = vultr_instance.kraken.main_ip
  description = ""
}

output "kraken_hostname" {
  value = vultr_instance.kraken.hostname
}

output "kraken_password" {
  value = vultr_instance.kraken.default_password
  sensitive = true
  description = ""
}
