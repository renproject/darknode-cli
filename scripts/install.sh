#!/bin/sh

# define color escape codes
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# Install unzip if command not found
if ! [ -x "$(command -v unzip)" ];then
  sudo apt-get install unzip
fi

# creating working directory
mkdir -p $HOME/.darknode/darknodes
mkdir -p $HOME/.darknode/bin
cd $HOME/.darknode
curl -s 'https://darknode.republicprotocol.com/darknode.zip' > darknode.zip
unzip -o darknode.zip

# get system information
ostype="$(uname -s)"
cputype="$(uname -m)"

# download darknode binary depending on the system and architecture
if [ "$ostype" = 'Linux' -a "$cputype" = 'x86_64' ]; then
    TERRAFORM_URL='https://releases.hashicorp.com/terraform/0.11.10/terraform_0.11.10_linux_amd64.zip'
    curl -s 'https://darknode.republicprotocol.com/darknode_linux_amd64' > ./bin/darknode
elif [ "$ostype" = 'Darwin' -a "$cputype" = 'x86_64' ]; then
    TERRAFORM_URL='https://releases.hashicorp.com/terraform/0.11.10/terraform_0.11.10_darwin_amd64.zip'
    curl -s 'https://darknode.republicprotocol.com/darknode_darwin_amd64' > ./bin/darknode
else
   echo 'unsupported OS type or architecture'
   cd ..
   rm -rf .darknode
   exit 1
fi

chmod +x bin/darknode

# download terraform
curl -s "$TERRAFORM_URL" > terraform.zip
unzip terraform
chmod +x terraform
mv terraform bin/terraform


# clean up zip files
rm darknode.zip
rm terraform.zip

# make sure the binary is installed in the path
if ! [ -x "$(command -v darknode)" ]; then
  path=$SHELL
  shell=${path##*/}

  if [ "$shell" = 'zsh' ] ; then
    if [ -f "$HOME/.zprofile" ] ; then
      echo '\nexport PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.zprofile
    elif [ -f "$HOME/.zshrc" ] ; then
      echo '\nexport PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.zshrc
    elif [ -f "$HOME/.profile" ] ; then
      echo '\nexport PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.profile
    fi
  elif  [ "$shell" = 'bash' ] ; then
    if [ -f "$HOME/.bash_profile" ] ; then
      echo '\nexport PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.bash_profile
    elif [ -f "$HOME/.bashrc" ] ; then
      echo '\nexport PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.bashrc
    elif [ -f "$HOME/.profile" ] ; then
      echo '\nexport PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.profile
    else
      echo '\nexport PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.bash_profile
    fi
  elif [ -f "$HOME/.profile" ] ; then
    echo '\nexport PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.profile
  else
    echo '\nexport PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.profile
  fi

  echo ''
  echo 'If you are using a custom shell, make sure you update your PATH.'
  echo "${GREEN}export PATH=\$PATH:\$HOME/.darknode/bin ${NC}"
fi

echo ''
echo "${GREEN}Done! Restart terminal and run the command below to begin.${NC}"
echo ''
echo "${GREEN}darknode up --help ${NC}"