# Ping-42 Function Library

This is the codebase of the 42 library and where the majority of infrastructure hacking takes place.

- [Ping-42 Function Library](#ping-42-function-library)
  - [Codespace Configuration](#codespace-configuration)
    - [Option 1 - Cloud Codespace _recommended_](#option-1---cloud-codespace-recommended)
    - [Option 2 - Local Image](#option-2---local-image)
  - [Further Notes](#further-notes)
    - [Calculating socket RTT](#calculating-socket-rtt)
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

## Further Notes

Further considerations and notes.

### Calculating socket RTT

Calculating RTT of client-server communication is possible on several levels/
First, for TCP connections, Linux maintains a list of kernel counters that [measure TCP rtt](https://github.com/torvalds/linux/blob/master/include/uapi/linux/tcp.h#L244) for every active TCP connection.
This effectively limits the operating system compatibility of our application [to Linux](https://stackoverflow.com/questions/71787548/how-to-measure-rtt-latency-through-tcp-clients-created-in-golang-from-a-tcp-se), 
but allows us to [reduce the hackiness](https://linuxgazette.net/136/pfeiffer.html) of sending various socket commands to generate a packet between the client and the server so we can measu

## Debugging Raw Sockets in VS Code
To debug the traceroute package and others that require raw socket operations, you can run VS Code with root privileges (preferably within a virtual environment to minimize security risks). Use the following on Linux:

Open a terminal and navigateto the lib directory.
Use this to launch vscode as root:

```bash
sudo code . --no-sandbox --user-data-dir ~/vsdata
```