package provider

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/republicprotocol/darknode-cli/util"
)

type awsTerraform struct {
	Name          string
	Region        string
	InstanceType  string
	SshPubKey     string
	SshPriKeyPath string
	AccessKey     string
	SecretKey     string
	ConfigPath    string
	IPFS          string
}

// tfConfig generates the terraform config file for deploying to AWS.
func (p providerAws) tfConfig(name, region, instance, ipfs string) error {
	tf := awsTerraform{
		Name:          name,
		Region:        region,
		InstanceType:  instance,
		SshPubKey:     filepath.Join(util.NodePath(name), "ssh_keypair.pub"),
		SshPriKeyPath: filepath.Join(util.NodePath(name), "ssh_keypair"),
		AccessKey:     p.accessKey,
		SecretKey:     p.secretKey,
		ConfigPath:    filepath.Join(util.NodePath(name), "config.json"),
		IPFS:          ipfs,
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
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_key_pair" "darknode" {
  key_name   = "{{.Name}}"
  public_key = file("{{.SshPubKey}}")
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

  provisioner "remote-exec" {

	inline = [
      "sudo adduser darknode --gecos \",,,\" --disabled-password",
      "sudo rsync --archive --chown=darknode:darknode ~/.ssh /home/darknode",
      "sudo DEBIAN_FRONTEND=noninteractive apt-get -y update",
      "sudo DEBIAN_FRONTEND=noninteractive apt-get -y upgrade",
      "sudo DEBIAN_FRONTEND=noninteractive apt-get -y dist-upgrade",
      "sudo DEBIAN_FRONTEND=noninteractive apt-get -y auto-remove",
      "sudo apt-get update",
      "sudo apt-get -y install jq",
      "sudo apt-get install ufw",
      "sudo ufw limit 22/tcp",
      "sudo ufw allow 18514/tcp", 
      "sudo ufw allow 18515/tcp", 
      "sudo ufw --force enable",
	]

    connection {
      host        = coalesce(self.public_ip, self.private_ip)
      type        = "ssh"
      user        = "ubuntu"
      private_key = file("{{.SshPriKeyPath}}")
    }
  }

  provisioner "file" {

    source      = "{{.ConfigPath}}"
    destination = "$HOME/config.json"

    connection {
      host        = coalesce(self.public_ip, self.private_ip)
      type        = "ssh"
      user        = "darknode"
      private_key = file("{{.SshPriKeyPath}}")
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
      host        = coalesce(self.public_ip, self.private_ip)
      type        = "ssh"
      user        = "darknode"
      private_key = file("{{.SshPriKeyPath}}")
    }
  }
}

output "provider" {
  value = "aws"
}

output "ip" {
  value = "${aws_instance.darknode.public_ip}"
}`

// {{if .AllocationID}}
// resource "aws_eip_association" "eip_assoc" {
// instance_id   = "${aws_instance.darknode.id}"
// allocation_id = "${var.allocation_id}"
// }{{else}}{{end}}
