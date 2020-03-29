package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/renproject/darknode-cli/util"
)

type gcpTerraform struct {
	Name           string
	CredentialFile string
	Project        string
	Zone           string
	MachineType    string
	ConfigPath     string
	PubKeyPath     string
	PriKeyPath     string
	ServiceFile    string
	LatestVersion  string
}

func (p providerGcp) tfConfig(name, project, zone, machine, latestVersion string) error {
	tf := gcpTerraform{
		Name:           name,
		CredentialFile: p.credFile,
		Project:        project,
		Zone:           zone,
		MachineType:    machine,
		ConfigPath:     fmt.Sprintf("~/.darknode/darknodes/%v/config.json", name),
		PubKeyPath:     fmt.Sprintf("~/.darknode/darknodes/%v/ssh_keypair.pub", name),
		PriKeyPath:     fmt.Sprintf("~/.darknode/darknodes/%v/ssh_keypair", name),
		ServiceFile:    darknodeService,
		LatestVersion:  latestVersion,
	}

	t, err := template.New("gcp").Parse(gcpTemplate)
	if err != nil {
		return err
	}
	tfFile, err := os.Create(filepath.Join(util.NodePath(name), "main.tf"))
	if err != nil {
		return err
	}
	return t.Execute(tfFile, tf)
}

var gcpTemplate = `
provider "google" {
  credentials = file("{{.CredentialFile}}")
  project     = "{{.Project}}"
  zone        = "{{.Zone}}"
}

/*******************
**** NETWORKING ****
*******************/

resource "google_compute_network" "darknode_network" {
  name = "{{.Name}}"
}

resource "google_compute_firewall" "darknode_firewall" {
  network     = "${google_compute_network.darknode_network.name}"
  name        = "{{.Name}}"
  
  allow {
    protocol = "icmp" 
  }

  allow{
	protocol = "tcp"
    ports = ["22", "18514", "18515"]
  }	

  target_tags   = ["darknode"]
}

/***********
**** VM ****
***********/

resource "google_compute_instance" "darknode" {
  name         = "{{.Name}}"
  machine_type = "{{.MachineType}}"
  allow_stopping_for_update = true  

  boot_disk {
	initialize_params {
	  image = "ubuntu-os-cloud/ubuntu-1804-lts"
    }
  }

  scheduling {
    on_host_maintenance = "MIGRATE"
  }

  network_interface {
    network = "${google_compute_network.darknode_network.name}"
    access_config {}
  }

  service_account {
    scopes = ["service-control", "service-management", "storage-ro", "trace-append", "monitoring-write", "logging-write"]
  }

  tags = ["darknode"]

  metadata = {
    ssh-keys = "ubuntu:${file("{{.PubKeyPath}}")}"
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
      "until sudo apt-get install ufw; do sleep 4; done",
      "sudo ufw limit 22/tcp",
      "sudo ufw allow 18514/tcp", 
      "sudo ufw allow 18515/tcp", 
      "sudo ufw --force enable",
      "until sudo apt-get -y install jq; do sleep 4; done",
	]

    connection {
      type        = "ssh"
      user        = "ubuntu"
      private_key = file("{{.PriKeyPath}}")
      host        = self.network_interface[0].access_config[0].nat_ip
    }
  }

  provisioner "file" {

    source      = "{{.ConfigPath}}"
    destination = "$HOME/config.json"

    connection {
      type        = "ssh"
      user        = "darknode"
      private_key = file("{{.PriKeyPath}}")
      host        = self.network_interface[0].access_config[0].nat_ip
    }
  }

  provisioner "remote-exec" {
	
	inline = [
      "set -x",
	  "mkdir -p $HOME/.darknode/bin",
      "mkdir -p $HOME/.config/systemd/user",
      "mv $HOME/config.json $HOME/.darknode/config.json",
	  "curl -sL https://www.github.com/renproject/darknode-release/releases/latest/download/darknode > ~/.darknode/bin/darknode",
	  "chmod +x ~/.darknode/bin/darknode",
      "echo {{.LatestVersion}} > ~/.darknode/version.md",
	  <<EOT
	  echo "{{.ServiceFile}}" > ~/.config/systemd/user/darknode.service
      EOT
      ,
	  "loginctl enable-linger darknode",
      "systemctl --user enable darknode.service",
      "systemctl --user start darknode.service",
	]

    connection {
      type        = "ssh"
      user        = "darknode"
      private_key = file("{{.PriKeyPath}}")
      host        = self.network_interface[0].access_config[0].nat_ip
    }
  }
}

output "provider" {
  value = "gcp"
}

output "ip" {
  value = "${google_compute_instance.darknode.network_interface[0].access_config[0].nat_ip}"
}`
