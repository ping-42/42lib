#!/usr/bin/env bash

set +e

echo "[*] Updating codebase..."

# Pull in our code
# Warning: this gets cached on the pre-build if its not being put in the right place.
# https://docs.github.com/en/codespaces/prebuilding-your-codespaces/configuring-prebuilds

# All of our repos
REPOS=("sensor" "server" "scheduler" "admin-ui" "admin-api")

# Change the volume ownership for configs to work
sudo chown vscode: -R "${HOME}/.config"

# If the Github token is empty, try to find out if we're logged in
if [[ -z ${GITHUB_TOKEN} ]]; then
	echo "[*] Getting Github Token from GH cli..."
	GITHUB_TOKEN=$(gh auth token)
	export GITHUB_TOKEN
fi

# Ensure git is setup correctly
gh auth setup-git

# Massage the workspace
cd /workspaces || exit
sudo chown vscode: /workspaces

if [[ ! -d "./sensor/" ]]; then
	for REPO in "${REPOS[@]}"; do
		echo "[*] Checking out repo ${REPO}..."
		gh repo clone "ping-42/${REPO}" -- -q
	done
else
	for REPO in "${REPOS[@]}"; do
		echo "[*] Updating repo ${REPO}..."
		(cd "${REPO}" && git pull)
	done
fi

# Start the timescale and redis in detached mode
docker compose -f 42lib/.devcontainer/docker-compose.yml --progress plain up -d

# echo "[*] Updating Golang deps..."
# for REPO in "${REPOS[@]}"; do
#     echo "[*] Updating deps in repo ${REPO}..."
#     (cd "${REPO}" && go get ./...)
# done

# Fix permissions...
# NOTE: this is broken - it has to be fixed by Microsoft
# It tells the user that there is an error and the container needs to be rebuilt, which is a lie
# https://github.com/microsoft/vscode-remote-release/issues/4442

#sudo chmod o-w -R /workspaces/ "${HOME}"
#sudo chown vscode: -R /workspaces/ "${HOME}"

#sudo chmod o+w /workspaces/42lib/.devcontainer/devcontainer.json
