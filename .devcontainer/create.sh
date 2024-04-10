#!/usr/bin/env bash
set +e

# Check if we ran before
if [[ -d "${HOMOE}/.ssh" ]]; then
	echo "No need to run - we're done here"
	exit 0
fi

# Shell quirks
mkdir -p "${HOME}/.ssh"

# Fixup Git aliases
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit
git config --global alias.st status
git config --global pull.rebase false
git config --global push.autoSetupRemote true

# Update Golang for private repos
go env -w "GOPRIVATE=github.com/ping-42/*"

# Get a default set of apt stuff
sudo apt update -qq
