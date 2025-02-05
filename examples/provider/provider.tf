provider "ripe" {}

data "ripedb_object" "mntner_ripe" {
  class = "mntner"
  value = "RIPE-MNT"
}
