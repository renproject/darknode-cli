#!/bin/sh

runTest()
{
  DIR=$1
  TYPE=$2

  vagrant up
  if [ "$TYPE" = 'zsh' ] ; then
    vagrant ssh -- -t '
  echo 2
  curl https://darknode.republicprotocol.com/install.sh -sSf | $0
  exit
  '

  vagrant ssh -- -t '
  source $HOME/.profile
  darknode list
  exit
  '
  else
    vagrant ssh -- -t '
  curl https://darknode.republicprotocol.com/install.sh -sSf | sh
  exit
  '

  vagrant ssh -- -t '
  source $HOME/.bashrc
  darknode list
  exit
  '
  fi

  echo "if you see [cannot find any node], the cli is successfully installed in the vagrant box with $DIR and $TYPE"
  sleep 5s

  vagrant destroy -f
}

cd "$HOME/vagrant_boxes" || exit 1
for dir in */ ; do
  cd "$dir"
  for type in */ ; do
    cd "$type"
    runTest $dir $type
    cd ..
  done
  cd ..
done

