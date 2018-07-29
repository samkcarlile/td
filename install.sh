#!/bin/bash

GO_INSTALLED=`which go`;
COLOR="$GOPATH/src/github.com/fatih/color"
PROMPTUI="$GOPATH/src/github.com/manifoldco/promptui"
CLI="$GOPATH/src/github.com/urfave/cli"

function check_for_go() {
  if [ -z $GO_INSTALLED ]; then
    echo "Go doesn't appear to be installed or it isn't in your path. Insall Go and try again."
    exit 1
  fi
}

function check_for_go_pkg() {
  if [ -d "$GOPATH/src/$1" ]; then
  echo "$1 is installed..."
else
  echo "$1 is NOT installed..."
  echo -n "Installing $1..."
  go get $1 && echo "done"
fi
}

function install_bash_autocomplete() {
  echo "Skipping..."
}

echo "Installing td..."
check_for_go
echo "1.) Installing go packages"
echo "--------------------------"
check_for_go_pkg "github.com/urfave/cli"
check_for_go_pkg "github.com/fatih/color"
check_for_go_pkg "github.com/manifoldco/promptui"
echo
echo "2.) Building td"
echo "---------------"
echo -n "Building..."
go build td.go && echo "done"
chmod +x ./td
echo "3.) Installing bash autocompletion..."
echo "-------------------------------------"
install_bash_autocomplete
echo
echo "⚡︎ Done Installing. Move td into somewhere in your path to use"