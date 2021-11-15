#!/bin/sh

main(){
  # Update this when minimum terraform version is changed.
  min_terraform_ver="0.12.24"
  cur_terraform_ver="0.12.24"

  # Check if darknode-cli has been installed
  if ! check_cmd darknode; then
    echo "cannot find darknode-cli"
    err "please install darknode-cli first"
  fi

  echo "Updating darknode-cli ..."

  # Check terraform version
  if check_cmd terraform; then
    version="$(terraform --version | grep 'Terraform v' | cut -d "v" -f2)"
    major="$(echo $version | cut -d. -f1)"
    minor="$(echo $version | cut -d. -f2)"
    patch="$(echo $version | cut -d. -f3)"
    requiredMajor="$(echo $min_terraform_ver | cut -d. -f1)"
    requiredMinor="$(echo $min_terraform_ver | cut -d. -f2)"
    requiredPatch="$(echo $min_terraform_ver | cut -d. -f3)"
    if [ "$major" -lt "$requiredMajor" ]; then
      err "Please upgrade your terraform to version above $min_terraform_ver"
    fi
    if [ "$minor" -lt "$requiredMinor" ]; then
      err "Please upgrade your terraform to version above $min_terraform_ver"
    fi
    if [ "$patch" -lt "$requiredPatch" ]; then
      err "Please upgrade your terraform to version above $min_terraform_ver"
    fi
  else
    install_terraform $cur_terraform_ver
  fi
  progressBar 40 100

  # Update the binary
  current=$(darknode --version | grep "darknode-cli version" | cut -d ' ' -f 3)
  latest=$(get_latest_release "renproject/darknode-cli")
  vercomp $current $latest
  if [ "$?" -eq "2" ]; then
    ostype="$(uname -s | tr '[:upper:]' '[:lower:]')"
    cputype="$(uname -m | tr '[:upper:]' '[:lower:]')"
    if [ $cputype = "x86_64" ];then
      cputype="amd64"
    fi

    darknode_url="https://www.github.com/renproject/darknode-cli/releases/latest/download/darknode_${ostype}_${cputype}"
    ensure downloader "$darknode_url" "$HOME/.darknode/bin/darknode"
    ensure chmod +x "$HOME/.darknode/bin/darknode"

    progressBar 100 100
    sleep 1
    echo ''
    echo "Done! Your 'darknode-cli' has been updated to $latest."
  else
    progressBar 100 100
    echo ''
    echo "You're running the latest version"
  fi
}

install_terraform(){
  terraform_ver="$1"
  mkdir -p $HOME/.darknode/bin
  ostype="$(uname -s | tr '[:upper:]' '[:lower:]')"
  cputype="$(uname -m | tr '[:upper:]' '[:lower:]')"
  if [ $cputype = "x86_64" ];then
      cputype="amd64"
  fi
  terraform_url="https://releases.hashicorp.com/terraform/${terraform_ver}/terraform_${terraform_ver}_${ostype}_${cputype}.zip"
  ensure downloader "$terraform_url" "$HOME/.darknode/bin/terraform.zip"
  ensure unzip -qq "$HOME/.darknode/bin/terraform.zip" -d "$HOME/.darknode/bin"
  ensure chmod +x "$HOME/.darknode/bin/terraform"
  rm "$HOME/.darknode/bin/terraform.zip"
}

# Source: https://sh.rustup.rs
check_cmd() {
    command -v "$1" > /dev/null 2>&1
}

# This wraps curl or wget. Try curl first, if not installed, use wget instead.
# Source: https://sh.rustup.rs
downloader() {
    if check_cmd curl; then
        if ! check_help_for curl --proto --tlsv1.2; then
            echo "Warning: Not forcing TLS v1.2, this is potentially less secure"
            curl --silent --show-error --fail --location "$1" --output "$2"
        else
            curl --proto '=https' --tlsv1.2 --silent --show-error --fail --location "$1" --output "$2"
        fi
    elif check_cmd wget; then
        if ! check_help_for wget --https-only --secure-protocol; then
            echo "Warning: Not forcing TLS v1.2, this is potentially less secure"
            wget "$1" -O "$2"
        else
            wget --https-only --secure-protocol=TLSv1_2 "$1" -O "$2"
        fi
    else
        echo "Unknown downloader"   # should not reach here
    fi
}

# Source: https://sh.rustup.rs
check_help_for() {
    local _cmd
    local _arg
    local _ok
    _cmd="$1"
    _ok="y"
    shift

    for _arg in "$@"; do
        if ! "$_cmd" --help | grep -q -- "$_arg"; then
            _ok="n"
        fi
    done

    test "$_ok" = "y"
}

# Source: https://sh.rustup.rs
err() {
    echo ''
    echo "$1" >&2
    exit 1
}

# Source: https://sh.rustup.rs
ensure() {
    if ! "$@"; then err "command failed: $*"; fi
}

get_latest_release() {
  curl --silent "https://api.github.com/repos/$1/releases/latest" | # Get latest release from GitHub api
    grep '"tag_name":' |                                            # Get tag line
    sed -E 's/.*"([^"]+)".*/\1/'                                    # Pluck JSON value
}

# source : https://stackoverflow.com/questions/4023830/how-to-compare-two-strings-in-dot-separated-version-format-in-bash
vercomp () {
    if [[ $1 == $2 ]]
    then
        return 0
    fi
    local IFS=.
    local i ver1=($1) ver2=($2)
    # fill empty fields in ver1 with zeros
    for ((i=${#ver1[@]}; i<${#ver2[@]}; i++))
    do
        ver1[i]=0
    done
    for ((i=0; i<${#ver1[@]}; i++))
    do
        if [[ -z ${ver2[i]} ]]
        then
            # fill empty fields in ver2 with zeros
            ver2[i]=0
        fi
        if ((10#${ver1[i]} > 10#${ver2[i]}))
        then
            return 1
        fi
        if ((10#${ver1[i]} < 10#${ver2[i]}))
        then
            return 2
        fi
    done
    return 0
}

# Source: https://github.com/fearside/ProgressBar
progressBar() {
    _progress=$1
    _done=$((_progress*5/10))
    _left=$((50-_done))
    done=""
    if ! [ $_done = "0" ];then
        done=$(printf '#%.0s' $(seq $_done))
    fi
    left=""
    if ! [ $_left = "0" ];then
      left=$(printf '=%.0s' $(seq $_left))
    fi
    printf "\rProgress : [$done$left] ${_progress}%%"
}

main "$@" || exit 1
