provider "digitalocean" {
  token = "${var.do_token}"
}

resource "digitalocean_ssh_key" "darknode" {
  name       = "${var.name}"
  public_key = "${file("${var.pub_key}")}"

}

resource "digitalocean_droplet" "darknode" {
  provider = "digitalocean"
  image = "ubuntu-18-04-x64"
  name = "${var.name}"
  region = "${var.region}"
  size = "${var.size}"
  ssh_keys = [
    "${digitalocean_ssh_key.darknode.id}"
  ]

  provisioner "remote-exec" {
    inline = [
      "adduser ubuntu --gecos \",,,\" --disabled-password",
      "usermod -aG sudo ubuntu",
      "echo \"ubuntu ALL=(ALL) NOPASSWD: ALL\" >> /etc/sudoers",
      "rsync --archive --chown=ubuntu:ubuntu ~/.ssh /home/ubuntu"
    ]

    connection {
      type        = "ssh"
      user        = "root"
      private_key = "${file("${var.pvt_key}")}"
    }
  }

  provisioner "file" {
    source = "${var.path}/darknodes/${var.name}/config.json"
    destination = "$HOME/darknode-config.json"

    connection {
      type = "ssh"
      user = "ubuntu"
      private_key = "${file("${var.pvt_key}")}"
    }
  }

  provisioner "file" {
    source = "${var.path}/provisions"
    destination = "$HOME/provisions"

    connection {
      type = "ssh"
      user = "ubuntu"
      private_key = "${file("${var.pvt_key}")}"
    }
  }

  provisioner "file" {
    source = "${var.path}/scripts"
    destination = "$HOME/scripts"

    connection {
      type = "ssh"
      user = "ubuntu"
      private_key = "${file("${var.pvt_key}")}"
    }
  }

  provisioner "remote-exec" {
    script = "${var.path}/scripts/up.sh"

    connection {
      type = "ssh"
      user = "ubuntu"
      private_key = "${file("${var.pvt_key}")}"
    }
  }

  provisioner "local-exec" {
    command = "echo /ip4/${digitalocean_droplet.darknode.ipv4_address}/tcp/18514/republic/${var.id} > multiAddress.out"
  }
}

output "multiaddress" {
  value  = "/ip4/${digitalocean_droplet.darknode.ipv4_address}/tcp/18514/republic/${var.id}"
}
