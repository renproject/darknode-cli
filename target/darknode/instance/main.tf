
variable "region" {}
variable "avz" {}
variable "ami" {}
variable "id" {}
variable "config" {}
variable "ec2_instance_type" {}
variable "ssh_public_key" {}
variable "ssh_private_key_location" {}
variable "access_key" {}
variable "secret_key" {}
variable "is_bootstrap" {}
variable "port" {}
variable "path" {}

provider "aws" {
  alias      = "darknode"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

resource "aws_security_group" "darknode" {
  provider    = "aws.darknode"
  name        = "falcon-sg-${var.id}"
  description = "Allow inbound SSH ,Republic Protocol traffic and logstash/kibana"

  // SSH
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  // Logstash
  ingress {
    from_port   = 9200
    to_port     = 9200
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  // Kibana
  ingress {
    from_port   = 5601
    to_port     = 5601
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  // Republic Protocol
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
  provider   = "aws.darknode"
  key_name   = "falcon-kp-${var.id}"
  public_key = "${var.ssh_public_key}"
}

output "multiaddress" {
  value       = "/ip4/${aws_instance.darknode.public_ip}/tcp/18514/republic/${var.id}"
}

resource "aws_instance" "darknode" {
  provider        = "aws.darknode"
  ami             = "${var.ami}"
  instance_type   = "${var.ec2_instance_type}"
  key_name        = "${aws_key_pair.darknode.key_name}"
  security_groups = ["${aws_security_group.darknode.name}"]

  provisioner "file" {
    source      = "${var.config}"
    destination = "/home/ubuntu/darknode-config.json"

    connection {
      type        = "ssh"
      user        = "ubuntu"
      private_key = "${file("${var.ssh_private_key_location}")}"
    }
  }

  provisioner "file" {
    source      = "${var.path}/provisions"
    destination = "/home/ubuntu/provisions"

    connection {
      type        = "ssh"
      user        = "ubuntu"
      private_key = "${file("${var.ssh_private_key_location}")}"
    }
  }

  provisioner "remote-exec" {
    script = "${var.path}/scripts/onCreate.sh"

    connection {
      type        = "ssh"
      user        = "ubuntu"
      private_key = "${file("${var.ssh_private_key_location}")}"
    }
  }

  provisioner "local-exec" {
      command = "echo /ip4/${aws_instance.darknode.public_ip}/tcp/${var.port}/republic/${var.id} > multiAddress.out"
  }
}
