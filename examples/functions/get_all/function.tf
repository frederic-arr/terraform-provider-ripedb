data "ripedb_object" "test" {
  class = "organisation"
  value = "ORG-TT1-TEST"
}

output "name" {
  value = provider::ripedb::get_first(data.ripedb_object.test.attributes, "org-name")
}
