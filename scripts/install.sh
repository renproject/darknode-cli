#!/usr/bin/sh

# creating working directory
mkdir -p $HOME/.darknode
mkdir -p $HOME/.darknode/darknodes
cd $HOME/.darknode
curl -s 'https://darknode.republicprotocol.com/darknode.zip' > darknode.zip
unzip darknode.zip

ostype="$(uname -s)"
cputype="$(uname -m)"

# Download deployer binanry and terraform
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

curl -s "$TERRAFORM_URL" > terraform.zip

# unzip terraform
mv terraform* terraform.zip
unzip terraform
chmod +x terraform
mv terraform bin/terraform

# clean up all the zip files
rm darknode.zip
rm terraform.zip

# make sure the binary is installed in the path
if ! [ -x "$(command -v darknode)" ]; then
  if test -n $BASH_VERSION  &&  [ -f "$HOME/.bash_profile" ] ; then
    echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.bash_profile
    source $HOME/.bash_profile
  elif test -n $ZSH_VERSION && [ -f "$HOME/.zprofile" ] ; then
    echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.zprofile
    source $HOME/.zprofile
  elif [ -f "$HOME/.profile" ]; then
    echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.profile
    source $HOME/.profile
  fi

  echo ''
  echo 'If you are using a custom shell, make sure you update your PATH.'
  echo '$ export PATH=$PATH:$HOME/.darknode/bin'
fi

echo ''
echo 'Done! Restart terminal and run the command below to begin.'
echo 'darknode up --help'