resource "ockam_identity" "example" {}

output "id" {
  value = ockam_identity.example.id
}
output "identity" {
  value     = ockam_identity.example.identity
  sensitive = true
}
output "vault" {
  value     = ockam_identity.example.vault
  sensitive = true
}
