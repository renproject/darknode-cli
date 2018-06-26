#!/usr/bin/sh

if [ -d "$HOME/.darknode" ] && [ -d "$HOME/.darknode/darknodes" ]; then
    cd $HOME/.darknode
    curl -s 'https://darknode.republicprotocol.com/darknode.zip' > darknode.zip
    unzip -o darknode.zip
else
    echo "cannot find the darknode folder"
    echo "please install first"
    exit 1
fi

# Check the os type and cpu architecture
ostype="$(uname -s)"
cputype="$(uname -m)"

# Download darknode-deployer
if [ "$ostype" = 'Linux' -a "$cputype" = 'x86_64' ]; then
    curl -s 'https://darknode.republicprotocol.com/darknode_linux_amd64' > ./bin/darknode
elif [ "$ostype" = 'Darwin' -a "$cputype" = 'x86_64' ]; then
    curl -s 'https://darknode.republicprotocol.com/darknode_darwin_amd64' > ./bin/darknode
else
   echo 'unsupported OS type or architecture'
   exit 1
fi

chmod +x bin/darknode
rm darknode.zip

echo ''
echo 'If you are using a custom shell, make sure you update your PATH.'
echo '$ export PATH=$PATH:$HOME/.darknode/bin'

echo ''
echo 'Done! Your Darknode Deployer has been updated.'