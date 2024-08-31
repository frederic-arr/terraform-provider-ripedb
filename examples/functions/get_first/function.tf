data "ripedb_object" "test" {
  class = "organisation"
  key   = "ORG-TT1-TEST"
}

output "mnt_by" {
  value = provider::ripedb::get_all(data.ripedb_object.test.attributes, "mnt-by")
}
