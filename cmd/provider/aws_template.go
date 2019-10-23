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
}

// tfConfig generates the terraform config file for deploying to AWS.
func (p providerAws) tfConfig(name, region, instance string) error {
	tf := awsTerraform{
		Name:          name,
		Region:        region,
		InstanceType:  instance,
		SshPubKey:     filepath.Join(util.NodePath(name), "ssh_keypair.pub"),
		SshPriKeyPath: filepath.Join(util.NodePath(name), "ssh_keypair"),
		AccessKey:     p.accessKey,
		SecretKey:     p.secretKey,
		ConfigPath:    filepath.Join(util.NodePath(name), "config.json"),
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
variable "name" {
  default = "{{.Name}}"
}

variable "region" {
  default = "{{.Region}}"
}

variable "instance_type" {
  default = "{{.InstanceType}}"
}

variable "public_key_path" {
  default = "{{.SshPubKey}}"
}

variable "private_key_path" {
  default = "{{.SshPriKeyPath}}"
}

variable "access_key" {
  default = "{{.AccessKey}}"
}

variable "secret_key" {
  default = "{{.SecretKey}}"
}

variable "config_path" {
  default = "{{.ConfigPath}}"
}

provider "aws" {
  region     = var.region
  access_key = var.access_key
  secret_key = var.secret_key
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
  name        = "darknode-sg-${var.name}"
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
  key_name   = var.name
  public_key = file(var.public_key_path)
}

resource "aws_instance" "darknode" {
  ami             = data.aws_ami.ubuntu.id
  instance_type   = var.instance_type
  key_name        = aws_key_pair.darknode.key_name
  security_groups = [aws_security_group.darknode.name]

  tags = {
    Name = var.name
  }

  provisioner "remote-exec" {
	inline = [
		"curl https://releases.renproject.io/darknode-cli/init.sh -sSf | sh"
	]

    connection {
      host        = coalesce(self.public_ip, self.private_ip)
      type        = "ssh"
      user        = "ubuntu"
      private_key = file("${var.private_key_path}")
    }
  }

  provisioner "file" {
    source      = var.config_path
    destination = "$HOME/darknode-config.json"

    connection {
      host        = coalesce(self.public_ip, self.private_ip)
      type        = "ssh"
      user        = "darknode"
      private_key = "${file("${var.private_key_path}")}"
    }
  }

  provisioner "remote-exec" {
	inline = [
		"curl https://releases.renproject.io/darknode-cli/install.sh -sSf | sh"
	]

    connection {
      host        = coalesce(self.public_ip, self.private_ip)
      type        = "ssh"
      user        = "darknode"
      private_key = file("${var.private_key_path}")
    }
  }
}

output "multiaddress" {
  value = "/ip4/${aws_instance.darknode.public_ip}/tcp/18514/ren/${var.address}"
}`

// {{if .AllocationID}}
// resource "aws_eip_association" "eip_assoc" {
// instance_id   = "${aws_instance.darknode.id}"
// allocation_id = "${var.allocation_id}"
// }{{else}}{{end}}
