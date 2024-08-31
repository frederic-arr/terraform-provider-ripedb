data "ripe_object" "test" {
  class = "organisation"
  key   = "ORG-TT1-TEST"
}

output "mnt_by" {
  value = provider::ripe::get_all(data.ripe_object.test.attributes, "mnt-by")
}
