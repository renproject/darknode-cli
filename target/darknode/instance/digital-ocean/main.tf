variable "do_token" {}
variable "ssh_fingerprint" {}
variable "name" {}
variable "region" {}
variable "size" {}
variable "id" {}
variable "port" {}
variable "path" {}


provider "digitalocean" {
  token = "${var.do_token}"
}

resource "digitalocean_droplet" "${name}" {
  provider = "digitalocean"
  image = "ubuntu-18-04-x64"
  name = "${name}"
  region = "${region}"
  size = "${size}"
  ssh_keys = [
    "${var.ssh_fingerprint}"
  ]

  provisioner "file" {
    source      = "${var.path}/provisions"
    destination = "/home/ubuntu/provisions"

    connection {
      type        = "ssh"
      user        = "ubuntu"
      private_key = "${file("${var.ssh_private_key_location}")}"
    }
  }

  provisioner "file" {
    source      = "${var.path}/scripts"
    destination = "/home/ubuntu/scripts"

    connection {
      type        = "ssh"
      user        = "ubuntu"
      private_key = "${file("${var.ssh_private_key_location}")}"
    }
  }

  provisioner "remote-exec" {
    script = "${var.path}/scripts/up.sh"

    connection {
      type        = "ssh"
      user        = "ubuntu"
      private_key = "${file("${var.ssh_private_key_location}")}"
    }
  }

  provisioner "local-exec" {
    command = "echo /ip4/${digitalocean_droplet.${name}.ipv4_address}/tcp/${var.port}/republic/${var.id} > multiAddress.out"
  }
}
