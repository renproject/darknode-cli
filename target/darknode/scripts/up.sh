#!/bin/sh

# Print commands before executing
set -x

# Do until not locked - will enter infinite loop if update fails
until sudo apt update; do sleep 2; done

# Install services
sudo mv ./provisions/darknode-updater.service /etc/systemd/system/darknode-updater.service
sudo mv ./provisions/darknode.service /etc/systemd/system/darknode.service
sudo mv ./provisions/logstash.service /etc/systemd/system/logstash.service
sudo mv ./scripts/updater.sh

# Install golang
wget https://dl.google.com/go/go1.10.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz
rm go1.10.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> $HOME/.profile
sudo ln -s /usr/local/go/bin/go /usr/bin/go

# Install logstash
wget https://artifacts.elastic.co/downloads/logstash/logstash-6.2.2.tar.gz
tar -xvf logstash-6.2.2.tar.gz
rm logstash-6.2.2.tar.gz
until sudo apt install -y default-jre; do sleep 2; done
mv ./provisions/logstash.conf ./logstash-6.2.2/darknode.conf

# Configure darknode and the updater
mkdir ./.darknode/
mv ./darknode-config.json ./.darknode/config.json
mv ./scripts/updater.sh ./.darknode/updater.sh

# Install metricbeat
curl -L -O https://artifacts.elastic.co/downloads/beats/metricbeat/metricbeat-6.2.2-amd64.deb
until sudo dpkg -i ./metricbeat-6.2.2-amd64.deb; do sleep 2; done
rm ./metricbeat-6.2.2-amd64.deb
sudo mv ./provisions/metricbeat.yml /etc/metricbeat/metricbeat.yml
sudo chown root /etc/metricbeat/metricbeat.yml
sudo chmod go-w /etc/metricbeat/metricbeat.yml
sudo metricbeat setup

# Install dep
mkdir -p $HOME/go/bin
export GOBIN=$HOME/go/bin
export GOPATH=$HOME/go
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Install darknode
until sudo apt install -y gcc; do sleep 2; done
mkdir -p ./go/src/github.com/republicprotocol
cd ./go/src/github.com/republicprotocol
git clone -b develop https://github.com/republicprotocol/republic-go.git
cd republic-go/cmd/darknode
$GOBIN/dep ensure
go install
cd $HOME

# Will fail if there are any files still in ./provisions/
rmdir ./provisions/
rm -rf ./scripts/

# Start services
sudo systemctl daemon-reload
sudo systemctl enable darknode-updater.service
sudo systemctl enable darknode.service 
sudo systemctl enable logstash.service
sudo systemctl start darknode-updater.service
sudo systemctl start darknode.service
sudo systemctl start logstash.service
sudo systemctl start metricbeat.service