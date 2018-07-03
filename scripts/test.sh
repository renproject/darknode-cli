#!/bin/sh

runTest()
{
  DIR=$1
  TYPE=$2

  vagrant up
  vagrant ssh -- -t '
  curl https://darknode.republicprotocol.com/install.sh -sSf | sh
  exit
  '
  vagrant ssh -- -t '
  source .profile
  darknode list
  sleep 5s
  exit
  '
  vagrant destory
}

cd "$HOME/vagrants" || exit 1
for dir in */ ; do
  cd "$dir"
  for type in */ ; do
    cd "$type"
    runTest $dir type
    cd ..
  done
  cd ..
done

