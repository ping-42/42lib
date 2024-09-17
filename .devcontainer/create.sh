#!/usr/bin/env bash
set +e

# Fixup Git aliases
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit
git config --global alias.st status
git config --global pull.rebase false
git config --global push.autoSetupRemote true
git config --global --add safe.directory /workspaces/42lib

# Update Golang for private repos
go env -w "GOPRIVATE=github.com/ping-42/*"

mkdir -p /home/vscode/.local/bin

# Install Syft for Goreleaser SBOM support
curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /home/vscode/.local/bin

# Install Goreleaser itself
echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | sudo tee /etc/apt/sources.list.d/goreleaser.list
sudo apt-get update -qq
sudo apt-get install --no-install-recommends goreleaser

# Install Geth and friends
sudo add-apt-repository -y ppa:ethereum/ethereum
sudo apt-get install -y ethereum