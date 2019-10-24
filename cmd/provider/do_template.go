package provider

type doTerraform struct {
	Name       string
	Token      string
	Region     string
	Size       string
	ConfigPath string
	PubKeyPath string
	PriKeyPath string
}

var doTemplate = `
variable "name" {
	default = "{{.Name}}"
}

variable "do_token" {
	default = "{{.Token}}"
}

variable "region" {
	default = "{{.Region}}"
}

variable "size" {
	default = "{{.Size}}"
}

variable "config_path" {
  default = "{{.ConfigPath}}"
}

variable "pub_key" {
  default = "{{.PubKeyPath}}"
}

variable "pri_key" {
  default = "{{.PriKeyPath}}"
}

provider "digitalocean" {
  token = var.do_token
}

resource "digitalocean_ssh_key" "darknode" {
   name       = var.name
   public_key = file(var.pub_key)
}

resource "digitalocean_droplet" "darknode" {
  provider   = digitalocean
  image      = "ubuntu-18-04-x64"
  name       = var.name
  region     = var.region
  size       = var.size
  monitoring = true

  ssh_keys   = [
    digitalocean_ssh_key.darknode.id
  ]

  provisioner "remote-exec" {
	inline = [
		"curl https://releases.renproject.io/darknode-cli/init.sh -sSf | sh"
	]

    connection {
      host        = self.ipv4_address
      type        = "ssh"
      user        = "root"
      private_key = file(var.pri_key)
    }
  }

  provisioner "file" {
    source = var.config_path
    destination = "$HOME/darknode-config.json"

    connection {
      host        = self.ipv4_address
      type        = "ssh"
      user        = "darknode"
      private_key = file(var.pri_key)
    }
  }

  provisioner "remote-exec" {
	inline = [
		"curl https://releases.renproject.io/darknode-cli/install.sh -sSf | sh"
	]

    connection {
      host        = self.ipv4_address
      type        = "ssh"
      user        = "darknode"
      private_key = file(var.pri_key)
    }
  }
}`

