#!/usr/bin/sh

# Check commands are all available
if ! [ -x "$(command -v unzip)" ];then
  sudo apt-get install unzip
fi

# creating working directory
mkdir -p $HOME/.darknode/darknodes
mkdir -p $HOME/.darknode/bin
cd $HOME/.darknode
curl -s 'https://darknode.republicprotocol.com/darknode.zip' > darknode.zip
unzip darknode.zip


# get system information
ostype="$(uname -s)"
cputype="$(uname -m)"

# Download terraform
if [ "$ostype" = 'Linux' -a "$cputype" = 'x86_64' ]; then
    TERRAFORM_URL='https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_linux_amd64.zip'
    curl -s 'https://darknode.republicprotocol.com/darknode_linux_amd64' > ./bin/darknode
elif [ "$ostype" = 'Darwin' -a "$cputype" = 'x86_64' ]; then
    TERRAFORM_URL='https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_darwin_amd64.zip'
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
  if test -n "$ZSH_VERSION" ; then
    echo "zsh"
    if [ -f "$HOME/.zprofile" ] ; then
      echo 3
      echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.zprofile
      source $HOME/.zprofile
    elif [ -f "$HOME/.zshrc" ] ;then
      echo 4
      echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.zshrc
      source $HOME/.zshrc
    fi
  elif test -n "$BASH_VERSION"; then
    echo "bash"
    if [ -f "$HOME/.bash_profile" ] ; then
      echo 1
      echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.bash_profile
      source $HOME/.bash_profile
    elif [ -f "$HOME/.bashrc" ] ; then
      echo 2
      echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.bashrc
      source $HOME/.bashrc
    fi
  elif [ -f "$HOME/.profile" ] ; then
    echo 5
    echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.profile
    . $HOME/.profile
  fi

  echo ''
  echo 'If you are using a custom shell, make sure you update your PATH.'
  echo '$ export PATH=$PATH:$HOME/.darknode/bin'
fi

echo ''
echo 'Done! Restart terminal and run the command below to begin.'
echo 'darknode up --help'