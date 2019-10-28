#!/bin/sh

# Check the os type and cpu architecture
ostype="$(uname -s)"
cputype="$(uname -m)"

if [ -d "$HOME/.darknode" ] && [ -d "$HOME/.darknode/darknodes" ]; then
    cd $HOME/.darknode
else
    echo "cannot find the darknode-cli"
    echo "please install darknode-cli first"
    exit 1
fi

# Download the latest binary darknode
if [ "$ostype" = 'Linux' -a "$cputype" = 'x86_64' ]; then
    curl -s -L 'https://www.github.com/renproject/darknode-cli/releases/latest/download/darknode_linux_amd64' > ./bin/darknode
elif [ "$ostype" = 'Darwin' -a "$cputype" = 'x86_64' ]; then
    curl -s -L 'https://www.github.com/renproject/darknode-cli/releases/latest/download/darknode_darwin_amd64' > ./bin/darknode
else
   echo 'unsupported OS type or architecture'
   exit 1
fi

chmod +x bin/darknode

echo ''
echo 'Done! Your darknode-cli has been updated.'