package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/renproject/darknode-cli/darknode"
	"github.com/renproject/darknode-cli/util"
)

type awsTerraform struct {
	Name            string
	Network         string
	Region          string
	InstanceType    string
	AccessKey       string
	SecretKey       string
	NodePath        string
	LatestVersion   string
	DarknodeService string
	GethService     string
}

// tfConfig generates the terraform config file for deploying to AWS.
func (p providerAws) tfConfig(name, region, instance, latestVersion string, network darknode.Network) error {
	tf := awsTerraform{
		Name:            name,
		Network:         string(network),
		Region:          region,
		InstanceType:    instance,
		NodePath:        fmt.Sprintf("~/.darknode/darknodes/%v", name),
		AccessKey:       p.accessKey,
		SecretKey:       p.secretKey,
		DarknodeService: darknodeService,
		GethService:     gethService(network, p.Name()),
		LatestVersion:   latestVersion,
	}

	t, err := template.New("aws").Parse(awsTemplate)
	if err != nil {
		return err
	}
	tfFile, err := os.Create(filepath.Join(util.NodePath(name), "main.tf"))
	if err != nil {
		return err
	}
	return t.Execute(tfFile, tf)
}

var awsTemplate = `
provider "aws" {
  region     = "{{.Region}}"
  access_key = "{{.AccessKey}}"
  secret_key = "{{.SecretKey}}"
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_security_group" "darknode" {
  name        = "darknode-sg-{{.Name}}"
  description = "Allow inbound SSH and REN project traffic"

  // SSH
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  // ren project
  ingress {
    from_port   = 18514
    to_port     = 18515
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  // ren evm 
  ingress {
    from_port   = 8545
    to_port     = 8545
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 30301
    to_port     = 30301
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 30301
    to_port     = 30301
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_key_pair" "darknode" {
  key_name   = "{{.Name}}"
  public_key = file("{{.NodePath}}/ssh_keypair.pub")
}

resource "aws_instance" "darknode" {
  ami             = data.aws_ami.ubuntu.id
  instance_type   = "{{.InstanceType}}"
  key_name        = aws_key_pair.darknode.key_name
  security_groups = [aws_security_group.darknode.name]
  monitoring      = true 

  tags = {
    Name = "{{.Name}}"
  }

  root_block_device {
    volume_type = "gp3"
    volume_size = 15
  }

  provisioner "remote-exec" {

	inline = [
      "set -x",
      "until sudo apt update; do sleep 4; done",
      "sudo adduser darknode --gecos \",,,\" --disabled-password",
      "sudo rsync --archive --chown=darknode:darknode ~/.ssh /home/darknode",
      "sudo DEBIAN_FRONTEND=noninteractive apt-get -y update",
      "sudo DEBIAN_FRONTEND=noninteractive apt-get -y upgrade",
      "sudo DEBIAN_FRONTEND=noninteractive apt-get -y dist-upgrade",
      "sudo DEBIAN_FRONTEND=noninteractive apt-get -y autoremove",
      "sudo apt-get install ufw",
      "sudo ufw limit 22/tcp",
      "sudo ufw allow 18514/tcp", 
      "sudo ufw allow 18515/tcp", 
      "sudo ufw allow 8545/tcp", 
      "sudo ufw allow 30301", 
      "sudo ufw --force enable",
	]

    connection {
      host        = coalesce(self.public_ip, self.private_ip)
      type        = "ssh"
      user        = "ubuntu"
      private_key = file("{{.NodePath}}/ssh_keypair")
    }
  }

  provisioner "file" {

    source      = "{{.NodePath}}/config.json"
    destination = "$HOME/config.json"

    connection {
      host        = coalesce(self.public_ip, self.private_ip)
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
      host        = coalesce(self.public_ip, self.private_ip)
      type        = "ssh"
      user        = "darknode"
      private_key = file("{{.NodePath}}/ssh_keypair")
    }
  }

  provisioner "file" {

    source      = "{{.NodePath}}/key.prv"
    destination = "$HOME/key.prv"

    connection {
      host        = coalesce(self.public_ip, self.private_ip)
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
      host        = coalesce(self.public_ip, self.private_ip)
      type        = "ssh"
      user        = "darknode"
      private_key = file("{{.NodePath}}/ssh_keypair")
    }
  }
}

output "provider" {
  value = "aws"
}

output "ip" {
  value = aws_instance.darknode.public_ip
}`
