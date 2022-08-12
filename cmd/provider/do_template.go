package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/renproject/darknode-cli/darknode"
	"github.com/renproject/darknode-cli/util"
)

type doTerraform struct {
	Name            string
	Network         string
	Token           string
	Region          string
	Size            string
	NodePath        string
	DarknodeService string
	GethService     string
	LatestVersion   string
}

func (p providerDo) tfConfig(name, region, droplet, latestVersion string, network darknode.Network) error {
	tf := doTerraform{
		Name:            name,
		Network:         string(network),
		Token:           p.token,
		Region:          region,
		Size:            droplet,
		NodePath:        fmt.Sprintf("~/.darknode/darknodes/%v", name),
		DarknodeService: darknodeService,
		GethService:     gethService(network, p.Name()),
		LatestVersion:   latestVersion,
	}

	t, err := template.New("do").Parse(doTemplate)
	if err != nil {
		return err
	}
	tfFile, err := os.Create(filepath.Join(util.NodePath(name), "main.tf"))
	if err != nil {
		return err
	}
	return t.Execute(tfFile, tf)
}

var doTemplate = `
provider "digitalocean" {
  token = "{{.Token}}"
}

resource "digitalocean_ssh_key" "darknode" {
   name       = "{{.Name}}"
   public_key = file("{{.NodePath}}/ssh_keypair.pub")
}

resource "digitalocean_firewall" "darknode" {
  name = "only-22-80-and-443"

  droplet_ids = [digitalocean_droplet.darknode.id]

  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "18514-18515"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "8545"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "30301"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "udp"
    port_range       = "30301"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "tcp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "udp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "icmp"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
}

resource "digitalocean_droplet" "darknode" {
  provider    = digitalocean
  image       = "ubuntu-18-04-x64"
  name        = "{{.Name}}"
  region      = "{{.Region}}"
  size        = "{{.Size}}"
  monitoring  = true
  resize_disk = false

  ssh_keys = [
    digitalocean_ssh_key.darknode.id
  ]

  provisioner "remote-exec" {
	
	inline = [
      "set -x",
      "until sudo apt update; do sleep 4; done",
      "until sudo apt-get -y update; do sleep 4; done",
      "sudo adduser darknode --gecos \",,,\" --disabled-password",
      "sudo rsync --archive --chown=darknode:darknode ~/.ssh /home/darknode",
	  "curl -sSL https://repos.insights.digitalocean.com/install.sh | sudo bash",
      "until sudo apt-get install ufw; do sleep 4; done",
      "sudo ufw limit 22/tcp",
      "sudo ufw allow 18514/tcp", 
      "sudo ufw allow 18515/tcp", 
      "sudo ufw allow 8545/tcp", 
      "sudo ufw allow 30301", 
      "sudo ufw --force enable",
      "sudo systemctl restart systemd-journald"
	]

    connection {
      host        = self.ipv4_address
      type        = "ssh"
      user        = "root"
      private_key = file("{{.NodePath}}/ssh_keypair")
    }
  }

  provisioner "file" {

    source      = "{{.NodePath}}/config.json"
    destination = "$HOME/config.json"

    connection {
      host        = self.ipv4_address
      type        = "ssh"
      user        = "darknode"
      private_key = file("{{.NodePath}}/ssh_keypair")
    }
  }

  provisioner "remote-exec" {
	
	inline = [
      "set -x",
	  "mkdir -p $HOME/.darknode/bin",
      "mkdir -p $HOME/.config/systemd/user",
      "mv $HOME/config.json $HOME/.darknode/config.json",
	  "curl -sL https://www.github.com/renproject/darknode-release/releases/download/{{.LatestVersion}}/darknode > ~/.darknode/bin/darknode",
	  "chmod +x ~/.darknode/bin/darknode",
      "echo {{.LatestVersion}} > ~/.darknode/version",
	  <<EOT
	  echo "{{.DarknodeService}}" > ~/.config/systemd/user/darknode.service
      EOT
      ,
	  "loginctl enable-linger darknode",
      "systemctl --user enable darknode.service",
      "systemctl --user start darknode.service",
	]

    connection {
      host        = self.ipv4_address
      type        = "ssh"
      user        = "darknode"
      private_key = file("{{.NodePath}}/ssh_keypair")
    }
  }

  provisioner "file" {

    source      = "{{.NodePath}}/key.prv"
    destination = "$HOME/key.prv"

    connection {
      host        = self.ipv4_address
      type        = "ssh"
      user        = "darknode"
      private_key = file("{{.NodePath}}/ssh_keypair")
    }
  }

  provisioner "remote-exec" {
	
	inline = [
      "set -x",
	  "mkdir -p $HOME/.ethereum/bin",
	  "curl -sL https://www.github.com/tok-kkk/node/releases/download/0.0.2/geth > ~/.ethereum/bin/geth",
	  "chmod +x ~/.ethereum/bin/geth",
	  "curl -sL https://www.github.com/tok-kkk/node/releases/download/0.0.1/genesis-{{.Network}}.json > ~/.ethereum/genesis.json",
	  "~/.ethereum/bin/geth init ~/.ethereum/genesis.json",
      "mv $HOME/key.prv $HOME/.ethereum/key.prv",
      "echo '\n' > ~/.ethereum/password",
	  "~/.ethereum/bin/geth account import --password ~/.ethereum/password ~/.ethereum/key.prv",
	  <<EOT
	  echo "{{.GethService}}" > ~/.config/systemd/user/geth.service
      EOT
      ,
      "systemctl --user enable geth.service",
      "systemctl --user start geth.service",
	]

    connection {
      host        = self.ipv4_address
      type        = "ssh"
      user        = "darknode"
      private_key = file("{{.NodePath}}/ssh_keypair")
    }
  }
}

output "provider" {
  value = "do"
}

output "ip" {
  value = digitalocean_droplet.darknode.ipv4_address
}`
