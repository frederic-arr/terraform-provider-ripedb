provider "ripe" {
  certificate = file("$PATH_TO/cert.pem")
  key         = file("$PATH_TO/key.pem")
}
