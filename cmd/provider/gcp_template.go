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
	IPFS           string
}

func (p providerGcp) tfConfig(name, project, zone, machine, ipfs string) error {
	tf := gcpTerraform{
		Name:           name,
		CredentialFile: p.credFile,
		Project:        project,
		Zone:           zone,
		MachineType:    machine,
		ConfigPath:     fmt.Sprintf("~/.darknode/darknodes/%v/config.json", name),
		PubKeyPath:     fmt.Sprintf("~/.darknode/darknodes/%v/ssh_keypair.pub", name),
		PriKeyPath:     fmt.Sprintf("~/.darknode/darknodes/%v/ssh_keypair", name),
		IPFS:           ipfs,
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
      "until sudo apt update; do sleep 2; done",
      "sudo adduser darknode --gecos \",,,\" --disabled-password",
      "sudo rsync --archive --chown=darknode:darknode ~/.ssh /home/darknode",
      "sudo DEBIAN_FRONTEND=noninteractive apt-get -y update",
      "sudo DEBIAN_FRONTEND=noninteractive apt-get -y upgrade",
      "sudo DEBIAN_FRONTEND=noninteractive apt-get -y dist-upgrade",
      "sudo DEBIAN_FRONTEND=noninteractive apt-get -y autoremove",
      "sudo apt-get -y install jq",
      "sudo apt-get install ufw",
      "sudo ufw limit 22/tcp",
      "sudo ufw allow 18514/tcp", 
      "sudo ufw allow 18515/tcp", 
      "sudo ufw --force enable",
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
      "wget -O darknode.gz {{.IPFS}}",
      "tar -zxvf darknode.gz",
	  "mkdir -p $HOME/.darknode",
      "mv $HOME/config.json $HOME/.darknode/config.json",
      "./install.sh",
      "rm -r darknode.gz bin config install.sh",
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
