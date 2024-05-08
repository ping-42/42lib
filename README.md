# Ping-42 Function Library

This is the codebase of the 42 library and where the majority of infrastructure hacking takes place.

- [Ping-42 Function Library](#ping-42-function-library)
  - [Codespace Configuration](#codespace-configuration)
    - [Option 1 - Cloud Codespace _recommended_](#option-1---cloud-codespace-recommended)
    - [Option 2 - Local Image](#option-2---local-image)
  - [Ping42 DNS Telemetry Daemon](#ping42-dns-telemetry-daemon)
    - [Calculating socket RTT](#calculating-socket-rtt)
    - [Running under Docker](#running-under-docker)
  - [Authenticate to GCP (unused)](#authenticate-to-gcp-unused)
  - [Debugging Raw Sockets in VS Code](#debugging-raw-sockets-in-vs-code)

## Codespace Configuration

The recommended way to work on the project is to run within a Visual Studio Devcontainer.
Please follow the general guidelines below to get started.

### Option 1 - Cloud Codespace _recommended_

In this setup, Github runs a VM and the VS Code client connects to it via the internet.
For normal internet connections that have good latency, this is the recommended method.

Simply start a codespace in the `main` branch and start hacking away. Once it boots, open the workspace file called `ping42.code-workspace` to laod the appropriate setup and repos as well. A workspace file can be opened via either the command menu (Cmd+Shift+P, start typing "Open workspace from File") or via the File menu of VSCode.

Once the codespace is up and running, clicking on Run & Debug (Shift-Cmd-D) allows for running the server and sensor both with debuggers attached.

Happy hacking!

### Option 2 - Local Image

Make sure you have docker running locally, and then try the following steps.

Create a codespace from the image:

![Alt text](https://i.ibb.co/Xs2yYz8/icon.jpg)

And then it can be run/used in several ways:

- through [the codespaces page](https://github.com/codespaces)
- install codespaces and save it applications (on MacOS)
- install the vscode extension and run from vscode locally

Default image configuration is located in the `devcontainer.json`. It has common basic tools and vscode extensions. Notably the linter:
Trunk. It highlights all issues as you type. The trunk icon appears on the left-most side panel. It supports golangci-lint among many others see `.trunk` for detailed configuration.

Rebuilding the image doc - [Link](https://github.com/github/docs/blob/ceb80203edd27c259c6da1b3d53310614780a56a/content/codespaces/developing-in-codespaces/rebuilding-the-container-in-a-codespace.md)

## Ping42 DNS Telemetry Daemon

Further considerations and notes.

### Calculating socket RTT

Calculating RTT of client-server communication is possible on several levels/
First, for TCP connections, Linux maintains a list of kernel counters that [measure TCP rtt](https://github.com/torvalds/linux/blob/master/include/uapi/linux/tcp.h#L244) for every active TCP connection.
This effectively limits the operating system compatibility of our application [to Linux](https://stackoverflow.com/questions/71787548/how-to-measure-rtt-latency-through-tcp-clients-created-in-golang-from-a-tcp-se), 
but allows us to [reduce the hackiness](https://linuxgazette.net/136/pfeiffer.html) of sending various socket commands to generate a packet between the client and the server so we can measure the real RTT of the TCP connection.

### Running under Docker

> Note: This is old stuff and unless you have a good reason, it should not be used.

In order to be able to run the project under docker, the following is possible:

```bash
# Create a docker volume to cache our build artifacts
docker volume create golang_cache

# Build and run the codebase
docker run -ti --rm -v golang_cache:/go -v $PWD:/go/src/github.com/ping42/sensor golang:1.21 bash -c "cd src/github.com/ping42/sensor; go run ."
```
## Authenticate to GCP (unused)

In order to authenticate the current workspace to be able to access resources inside Google Cloud, we probably need to authenticate:

```bash
gcloud init --project ping42-xxxxxxxx
```

This requires one to authenticate and provide permission for the current workspace to access the GCP account. To test that everything is working, try `gcloud functions list`.
In order to then interact with the BigQuery dataset, a simple `gcloud alpha bq datasets list` or `gcloud alpha bq datasets describe events`. The `gcloud` util allows querying and enumerating the datasets and tables with queries directly from the CLI.

## Debugging Raw Sockets in VS Code
To debug the traceroute package and others that require raw socket operations, you can run VS Code with root privileges (preferably within a virtual environment to minimize security risks). Use the following on Linux:

Open a terminal and navigateto the lib directory.
Use this to launch vscode as root:

```bash
sudo code . --no-sandbox --user-data-dir ~/vsdata
```