data "ripe_object" "test" {
  class = "organisation"
  key   = "ORG-TT1-TEST"
}

output "name" {
  value = provider::ripe::get_first(data.ripe_object.test.attributes, "org-name")
}
