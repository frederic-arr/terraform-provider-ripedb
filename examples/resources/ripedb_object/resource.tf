resource "ripedb_object" "person" {
  class = "person"
  value = "John Smith"
  attributes = [
    { name = "nic-hdl", value = "JS1-TEST" },
    { name = "address", value = "ACME, Inc." },
    { name = "phone", value = "+0" },
    { name = "mnt-by", value = "XYZ-MNT" },
  ]
}
