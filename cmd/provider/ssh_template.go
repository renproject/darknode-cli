package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/renproject/darknode-cli/util"
)

type sshTerraform struct {
	Name           string
	User           string
	Hostname       string
	PubKeyPath     string
	RootPriKeyPath string
	ConfigPath     string
	ServiceFile    string
	LatestVersion  string
}

func (p providerSsh) tfConfig(name, latestVersion string) error {
	tf := sshTerraform{
		Name:           name,
		User:           p.user,
		Hostname:       p.hostname,
		ConfigPath:     fmt.Sprintf("~/.darknode/darknodes/%v/config.json", name),
		PubKeyPath:     fmt.Sprintf("~/.darknode/darknodes/%v/ssh_keypair.pub", name),
		RootPriKeyPath: fmt.Sprintf(p.priKeyPath),
		ServiceFile:    darknodeService,
		LatestVersion:  latestVersion,
	}

	t, err := template.New("ssh").Parse(sshTemplate)
	if err != nil {
		return err
	}
	tfFile, err := os.Create(filepath.Join(util.NodePath(name), "main.tf"))
	if err != nil {
		return err
	}
	return t.Execute(tfFile, tf)
}

var sshTemplate = `
resource "null_resource" "darknode" {
	provisioner "remote-exec" {	
		inline = [
			"set -x",
			"until sudo apt update; do sleep 4; done",
			"sudo adduser darknode --gecos \",,,\" --disabled-password",
			"sudo rsync --archive --chown=darknode:darknode ~/.ssh /home/darknode",
			"until sudo apt-get -y update; do sleep 4; done",
			"until sudo apt-get -y upgrade; do sleep 4; done",
			"until sudo apt-get -y autoremove; do sleep 4; done",
			"until sudo apt-get install ufw; do sleep 4; done",
			"sudo ufw limit 22/tcp",
			"sudo ufw allow 18514/tcp",
			"sudo ufw allow 18515/tcp",
			"sudo ufw --force enable",
			"until sudo apt-get -y install jq; do sleep 4; done",
		]

		connection {
			host        = "{{.Hostname}}"
			type        = "ssh"
			user        = "{{.User}}"
			private_key = file("{{.RootPriKeyPath}}")
		}
	}

	provisioner "file" {
		source      = "{{.ConfigPath}}"
		destination = "$HOME/config.json"

		connection {
			host        = "{{.Hostname}}"
			type        = "ssh"
			user        = "darknode"
			private_key = file("{{.RootPriKeyPath}}")
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
			"echo {{.LatestVersion}} > ~/.darknode/version",
			<<EOT
			echo "{{.ServiceFile}}" > ~/.config/systemd/user/darknode.service
			EOT
			,
			"loginctl enable-linger darknode",
			"systemctl --user enable darknode.service",
			"systemctl --user start darknode.service",
		]

		connection {
			host        = "{{.Hostname}}"
			type        = "ssh"
			user        = "darknode"
			private_key = file("{{.RootPriKeyPath}}")
		}
	}

	provisioner "file" {
		source      = "{{.PubKeyPath}}"
		destination = "$HOME/.ssh/authorized_keys"

		connection {
			host        = "{{.Hostname}}"
			type        = "ssh"
			user        = "darknode"
			private_key = file("{{.RootPriKeyPath}}")
		}
	}
}

output "provider" {
	value = "ssh"
}

output "ip" {
	value = "{{.Hostname}}"
}`
