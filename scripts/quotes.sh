#!/bin/sh

# define color escape codes
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN} $ export PATH=\$PATH:\$HOME/.darknode/bin ${NC}"
echo -e "${GREEN} $ darknode up --help ${NC}"

path=$SHELL
shell=${path##*/}
 echo $shell
  if [ "$shell" = 'zsh' ] ; then
    echo "zsh"
    if [ -f "$HOME/.zprofile" ] ; then
      echo 3
    elif [ -f "$HOME/.zshrc" ] ;then
      echo 4
    fi
  elif  [ "$shell" = 'bash' ] ; then
    echo "bash"
    if [ -f "$HOME/.bash_profile" ] ; then
      echo 1
    elif [ -f "$HOME/.bashrc" ] ; then
      echo 2
    fi
  elif [ -f "$HOME/.profile" ] ; then
    echo 5
fi


#
#  if [ "$ostype" == "Linux" ] ; then
#    if [ -f "$HOME/.zprofile" ] ; then
#      echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.zprofile
#      . $HOME/.zprofile
#    elif [ -f "$HOME/.bash_profile" ] ; then
#      echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.bash_profile
#      . $HOME/.bash_profile
#    elif [ -f "$HOME/.profile" ] ; then
#      echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.profile
#      . $HOME/.profile
#    fi
#  fi
#
#  if [ "$ostype" == "Darwin" ] ; then
#    if [ -f "$HOME/.zprofile" ] ; then
#      echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.zprofile
#      source $HOME/.zprofile
#    elif [ -f "$HOME/.bash_profile" ] ; then
#      echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.bash_profile
#      source $HOME/.bash_profile
#    elif [ -f "$HOME/.profile" ] ; then
#      echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.profile
#      source $HOME/.profile
#    fi
#  fi
