package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/renproject/darknode-cli/util"
)

type doTerraform struct {
	Name       string
	Token      string
	Region     string
	Size       string
	ConfigPath string
	PubKeyPath string
	PriKeyPath string
	IPFS       string
}

func (p providerDo) tfConfig(name, region, droplet, ipfs string) error {
	tf := doTerraform{
		Name:       name,
		Token:      p.token,
		Region:     region,
		Size:       droplet,
		ConfigPath: fmt.Sprintf("~/.darknode/darknodes/%v/config.json", name),
		PubKeyPath: fmt.Sprintf("~/.darknode/darknodes/%v/ssh_keypair.pub", name),
		PriKeyPath: fmt.Sprintf("~/.darknode/darknodes/%v/ssh_keypair", name),
		IPFS:       ipfs,
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
   public_key = file("{{.PubKeyPath}}")
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
      "until sudo apt update; do sleep 2; done",
      "until sudo apt-get -y update; do sleep 2; done",
      "sudo adduser darknode --gecos \",,,\" --disabled-password",
      "sudo rsync --archive --chown=darknode:darknode ~/.ssh /home/darknode",
	  "curl -sSL https://repos.insights.digitalocean.com/install.sh | sudo bash",
      "sudo apt-get install ufw",
      "sudo ufw limit 22/tcp",
      "sudo ufw allow 18514/tcp", 
      "sudo ufw allow 18515/tcp", 
      "sudo ufw --force enable",
      "until sudo sudo apt-get -y install jq; do sleep 2; done",
	]

    connection {
      host        = self.ipv4_address
      type        = "ssh"
      user        = "root"
      private_key = file("{{.PriKeyPath}}")
    }
  }

  provisioner "file" {

    source      = "{{.ConfigPath}}"
    destination = "$HOME/config.json"

    connection {
      host        = self.ipv4_address
      type        = "ssh"
      user        = "darknode"
      private_key = file("{{.PriKeyPath}}")
    }
  }

  provisioner "remote-exec" {
	
	inline = [
      "wget -O darknode.gz {{.IPFS}}",
      "tar -zxvf darknode.gz",
	  "mkdir -p $HOME/.darknode",
      "mv $HOME/config.json $HOME/.darknode/config.json",
      "./install.sh",
      "rm -r darknode.gz bin config install.sh",
	]

    connection {
      host        = self.ipv4_address
      type        = "ssh"
      user        = "darknode"
      private_key = file("{{.PriKeyPath}}")
    }
  }
}

output "provider" {
  value = "do"
}

output "ip" {
  value = digitalocean_droplet.darknode.ipv4_address
}`
