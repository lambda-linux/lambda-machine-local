# ![Lambda Machine Local](logo.png)

Lambda Machine Local provides a stable and secure local docker development
environment. It is designed to complement a stable and and secure production
environment on [AWS](https://aws.amazon.com/).

Lambda Machine Local
uses
[Lambda Linux VirtualBox flavor](https://github.com/lambda-linux/lambda-linux-vbox) as
its container host OS and is based
on [Amazon Linux](https://aws.amazon.com/amazon-linux-ami/), the Linux operating
system that powers AWS.

You can download Lambda Machine Local for Mac, Linux and
Windows [here](https://github.com/lambda-linux/lambda-machine-local/releases).

It works a bit like this &ndash;

```console
$ lambda-machine-local create default
Running pre-create checks...
[...]
(default) Creating VirtualBox VM...
(default) Creating SSH key...
(default) Starting the VM...
(default) Check network to re-create if needed...
(default) Waiting for an IP...
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with lambda...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine,
run: lambda-machine-local env default

$ lambda-machine-local ls
NAME      ACTIVE   DRIVER       STATE     URL                         DOCKER    ERRORS
default   -        virtualbox   Running   tcp://192.168.99.100:2376   v1.12.6

$ eval "$(lambda-machine-local env default)"

$ docker run --rm amazonlinux echo hello world
Unable to find image 'amazonlinux:latest' locally
latest: Pulling from library/amazonlinux
c9141092a50d: Pull complete
Digest: sha256:2010c88ac1e7c118d61793eec71dcfe0e276d72b38dd86bd3e49da1f8c48bf54
Status: Downloaded newer image for amazonlinux:latest
hello world
```

**Note:** In the above example we
used
[Amazon Linux Container Image](http://docs.aws.amazon.com/AmazonECR/latest/userguide/amazon_linux_container_image.html).
Lambda Machine Local is designed to be a local development primitive and we
support all major operating systems used inside containers. So, feel free to use
your preferred container image operating system with Lambda Machine Local. **We
will gladly support you!**

If you don't have a preferred operating system for your container images, we
recommend that you consider Amazon Linux. It is a proven operating system and is
used by many AWS services such as [AWS Lambda](https://aws.amazon.com/lambda/).

-----------------------------------------

**Related resources**:
  [Website](https://lambda-linux.io/) |
  [Slack](http://slack.lambda-linux.io/) |
  [Discussion Forum](https://groups.google.com/group/lambda-linux) |
  [Twitter](https://twitter.com/lambda_linux) |
  [Blog](https://lambda-linux.io/blog/) |
  [FAQs](http://lambda-linux.io/faqs/#!/lambda-machine-local-questions)

**Table of contents**

  * [What is Lambda Machine Local?](#what_is_lambda_machine_local)
  * [Installation](#installation)
  * [Getting Started](#getting_started)
    * [Prerequisite Information](#prerequisite_information)
    * [Using Lambda Machine Local to run containers](#using_lambda_machine_local_to_run_containers)
      * [Creating a container host virtual machine](#creating_a_container_host_virtual_machine)
      * [Running containers](#running_containers)
      * [Starting and stopping container host virtual machine](#starting_and_stopping_container_host_virtual_machine)
      * [Working with container host virtual machine without specifying name](#working_with_container_host_virtual_machine_without_specifying_name)
      * [Auto updates and crash reporting](#auto_updates_and_crash_reporting)
  * [Command line reference](#command_line_reference)
  * [Customizing Lambda Linux VirtualBox flavor](#customizing_lambda_linux_virtualbox_flavor)
  * [A note about VirtualBox shared folder support](#a_note_about_virtualbox_shared_folder_support)

-----------------------------------------

<a name="what_is_lambda_machine_local"></a>
## What is Lambda Machine Local?

Lambda Machine Local is a tool that provides a local docker development
environment with a high bar for stability and security. You can use Lambda
Machine Local to create a supported container host on your Mac, Linux or
Windows.

Using `lambda-machine-local` commands, you can start, inspect, stop, and restart
a managed local container host. By pointing Lambda Machine Local CLI at a
running managed local host you can run `docker` commands directly on that host.

For example &ndash; You can run `lambda-machine-local env default`. Then after
following on-screen instructions to complete environment setup, you can run
`docker ps`, `docker run --rm amazonlinux echo hello world`, and so forth.

Following components are used by Lambda Machine Local.

| Component        | Why is it included? / Remarks |
| ---------------- | ------------------- |
| Container OS | [Lambda Linux VirtualBox Flavor](https://github.com/lambda-linux/lambda-linux-vbox) is used as the container host OS. Lambda Linux is based on [Amazon Linux](https://aws.amazon.com/amazon-linux-ami/), the Linux operating system that powers [AWS](https://aws.amazon.com/) |
| Linux Kernel | LTS version [4.9.20](https://cdn.kernel.org/pub/linux/kernel/v4.x/ChangeLog-4.9.20) |
| Docker Engine | Version [1.12.6](https://github.com/docker/docker/releases/tag/v1.12.6) |
| Storage Driver| Device mapper in [direct-lvm](https://docs.docker.com/v1.12/engine/userguide/storagedriver/device-mapper-driver/#/configure-direct-lvm-mode-for-production) mode |
| Orchestration | Amazon ECS [CLI](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/ECS_CLI.html) and Docker Compose version [1](https://docs.docker.com/compose/compose-file/compose-file-v1/) and [2](https://docs.docker.com/compose/compose-file/compose-file-v2/) is supported. There is no Docker Swarm support in Lambda Machine Local. |
| VirtualBox Guest Additions | Version [5.1.20](http://download.virtualbox.org/virtualbox/5.1.10/) is included for shared filesystem support |
| Libmachine | Lambda Machine Local uses libmachine version [0.11.0](https://github.com/docker/machine/tree/v0.10.0/libmachine) to manage local container host OS |

<a name="installation"></a>
## Installation

Lambda Machine Local binaries for Mac, Linux and Windows along with installation
instructions and SHA-256 / MD5 checksums is
available [here](https://github.com/lambda-linux/lambda-machine-local/releases).

Also install Docker [client](https://github.com/docker/docker/releases) version
1.12.6 or later on your local Mac, Linux or Windows host. If you prefer to use
Docker client from your local container host, you can `lambda-machine-local ssh
<LOCAL_VIRTUAL_MACHINE_NAME>` and run `docker` command from there.

```console
$ lambda-machine-local ssh default

[ll-user@ip-192-168-99-100 ~]$ docker info
```

Depending on your container workflow, please consider installing

  * [Amazon ECS CLI](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/ECS_CLI.html)
  * [Docker Compose](https://docs.docker.com/v1.12/compose/overview/)

[AWS CLI](https://aws.amazon.com/cli/) is also available on local container host
via `yum install ...`

```console
$ lambda-machine-local ssh default

[ll-user@ip-192-168-99-100 ~]$ sudo yum install -y aws-cli
```

For Windows users, we recommend that you please consider
installing [MSYS2](http://www.msys2.org/)
or [Git Bash](https://git-for-windows.github.io/). With MSYS2 you can install
`winpty`, `rsync` and `openssh` packages using `pacman`. `pacman` is a package
management tool similar to `yum` and `apt-get`. With Git Bash you get `winpty`
and `openssh`, but [no](https://github.com/git-for-windows/git/issues/347)
`rsync` package.

```console
username@DESKTOP-xxxxxxx MINGW64 ~
$ pacman -S rsync openssh winpty
resolving dependencies...
looking for conflicting packages...

Packages (3) openssh-7.3p1-2  rsync-3.1.2-2  winpty-0.4.0-2

Total Installed Size:  7.58 MiB
Net Upgrade Size:      0.00 MiB

:: Proceed with installation? [Y/n] Y

[...]

username@DESKTOP-xxxxxxx MINGW64 ~
$
```

If `rsync` and `openssh` packages are available, Lambda Machine Local makes uses
them when running `lambda-machine-local ssh` and `lambda-machine-local scp`
commands.

We anticipate that you'll frequently use these commands in your workflow,
therefore installing `rsync` and `openssh` packages can be very helpful in
enhancing your developer experience with Lambda Machine Local.

Also on Windows, sometimes when trying to delete a container host virtual
machine, you might encounter `Can't remove "<virtual_machine_name>"` error. For
example &ndash;

```console
username@DESKTOP-xxxxxxx MINGW64 ~
$ lambda-machine-local ls
NAME      ACTIVE   DRIVER       STATE     URL                         DOCKER    ERRORS
default   -        virtualbox   Running   tcp://192.168.99.100:2376   v1.12.6

username@DESKTOP-xxxxxxx MINGW64 ~
$ lambda-machine-local rm -y default
About to remove default
WARNING: This action will delete both local reference and remote instance.
Can't remove "default"

username@DESKTOP-xxxxxxx MINGW64 ~
$ lambda-machine-local ls
NAME      ACTIVE   DRIVER      STATE   URL   DOCKER   ERRORS
default            not found   Error                  open C:\Users\username\.docker\lambda-machine-local\
machines\default\config.json: The system cannot find the file specified.
```

If you encounter this issue, you can use [`lambda-machine-local rm`](#cli_rm)
with `--force` (`-f`) flag to remove the `default` directory and correct the
error.

```console
username@DESKTOP-xxxxxxx MINGW64 ~
$ lambda-machine-local rm -f default
About to remove default
WARNING: This action will delete both local reference and remote instance.
Error removing host "default": open C:\Users\username\.docker\lambda-machine-local\machines\default\config.json:
The system cannot find the file specified.
Successfully removed default

username@DESKTOP-xxxxxxx MINGW64 ~
$ lambda-machine-local ls
NAME   ACTIVE   DRIVER   STATE   URL   DOCKER   ERRORS

username@DESKTOP-xxxxxxx MINGW64 ~
$
```

<a name="getting_started"></a>
## Getting Started

<a name="prerequisite_information"></a>
### Prerequisite Information

Lambda Machine Local requires [VirtualBox](https://www.virtualbox.org/)
hypervisor to run container host VMs. It is designed to be a good citizen and
can work alongside with other VirtualBox tools such
as
[Docker Machine](https://docs.docker.com/machine/),
[Minikube](https://github.com/kubernetes/minikube)
and [Vagrant](https://www.vagrantup.com/).

However, if you are
using [Docker for Windows](https://docs.docker.com/docker-for-windows/), it
defaults to using Hyper-V. Unfortunately on Windows, running Hyper-V disables
VirtualBox. Please
see
[this](http://www.hanselman.com/blog/switcheasilybetweenvirtualboxandhypervwithabcdeditbootentryinwindows81.aspx) link
for information on how to setup Windows to toggle between Hyper-V and
VirtualBox.

<a name="using_lambda_machine_local_to_run_containers"></a>
### Using Lambda Machine Local to run containers

To run containers, you'll need to &ndash;

  * Create a new container host VM or start an existing one
  * Switch your docker client environment to container host VM
  * Use docker client to create, load and manage containers

Once you create a container host VM, you can reuse it as often as you like. Just
like any VirtualBox VM, it maintains its configuration between uses.

<a name="creating_a_container_host_virtual_machine"></a>
#### Creating a container host virtual machine

1. Open a command shell or terminal window.

   The command examples are shown in a bash shell. For a different shell (such as a C shell) the commands are the same unless otherwise noted.

2. Use `lambda-machine-local ls` to list available container host virtual machines.

   ```console
   $ lambda-machine-local ls
   NAME   ACTIVE   DRIVER   STATE   URL   DOCKER   ERRORS
   ```

3. Create a container host virtual machine.

   Run `lambda-machine-local create` command and provide a name for container host virtual machine. If this is your first container host virtual machine, name it as `default` as shown in the example. If you already have a `default`, then choose another name.

   ```console
   $ lambda-machine-local create default
   Running pre-create checks...
   Creating machine...
   (default) Copying /Users/username/.docker/lambda-machine-local/cache/lambda-linux-vbox.iso to
   /Users/username/.docker/lambda-machine-local/machines/default/lambda-linux-vbox.iso...
   (default) Creating VirtualBox VM...
   (default) Creating SSH key...
   (default) Starting the VM...
   (default) Check network to re-create if needed...
   (default) Waiting for an IP...
   Waiting for machine to be running, this may take a few minutes...
   Detecting operating system of created instance...
   Waiting for SSH to be available...
   Detecting the provisioner...
   Provisioning with lambda...
   Copying certs to the local machine directory...
   Copying certs to the remote machine...
   Checking connection to Docker...
   Docker is up and running!
   To see how to connect your Docker Client to the Docker Engine running on this virtual machine,
   run: lambda-machine-local env default
   ```

4. List container host virtual machines again to see the newly minted `default` virtual machine.

   ```console
   $ lambda-machine-local ls
   NAME      ACTIVE   DRIVER       STATE     URL                         DOCKER    ERRORS
   default   -        virtualbox   Running   tcp://192.168.99.100:2376   v1.12.6
   ```

5. Get the environment commands for new container host virtual machine.

   As noted in the output of the `lambda-machine-local create` command, you need to tell Docker client to talk to the new container host virtual machine. You can do this with `lambda-machine-local env` command.

   ```console
   $ lambda-machine-local env default
   export DOCKER_TLS_VERIFY="1"
   export DOCKER_HOST="tcp://192.168.99.100:2376"
   export DOCKER_CERT_PATH="/Users/username/.docker/lambda-machine-local/machines/default"
   export DOCKER_MACHINE_NAME="default"
   # Run this command to configure your shell:
   # eval $(lambda-machine-local env default)
   ```

6. Connect your shell to new container host virtual machine.

   ```console
   $ eval $(lambda-machine-local env default)
   ```

   **Note:** If you are using a shell other than `bash` and `zsh` please see [`env`](#cli_env) command's documentation for details on configuring your shell.

<a name="running_containers"></a>
#### Running containers

Run a container with `docker run` to verify your set up.

1. Use `docker run` to download and run `busybox` with a simple `echo` command.

   ```console
   $ docker run busybox echo hello world
   Unable to find image 'busybox:latest' locally
   latest: Pulling from library/busybox
   7520415ce762: Pull complete
   Digest: sha256:32f093055929dbc23dec4d03e09dfe971f5973a9ca5cf059cbfb644c206aa83f
   Status: Downloaded newer image for busybox:latest
   hello world
   ```

2. Get the container host virtual machine's IP address.

   Any exposed ports is available on container host virtual machine's IP address. You can use `lambda-machine-local ip` command to get the IP address.

   ```console
   $ lambda-machine-local ip default
   192.168.99.100
   ```

3. Run a webserver ([nginx](https://nginx.org/)) in a container with the following command.

   ```console
   $ docker run -d -p 8000:80 nginx
   ```

   When the image is finished pulling, you can hit the server at port 8000 on the IP address given to you by `lambda-machine-local ip`.

   ```console
   $ curl $(lambda-machine-local ip default):8000
   <!DOCTYPE html>
   <html>
   <head>
   <title>Welcome to nginx!</title>
   <style>
       body {
           width: 35em;
           margin: 0 auto;
           font-family: Tahoma, Verdana, Arial, sans-serif;
       }
   </style>
   </head>
   <body>
   <h1>Welcome to nginx!</h1>
   <p>If you see this page, the nginx web server is successfully installed and
   working. Further configuration is required.</p>

   <p>For online documentation and support please refer to
   <a href="http://nginx.org/">nginx.org</a>.<br/>
   Commercial support is available at
   <a href="http://nginx.com/">nginx.com</a>.</p>

   <p><em>Thank you for using nginx.</em></p>
   </body>
   </html>
   ```

   Depending on your needs you can create and manage multiple container host virtual machines by running `lambda-machine-local create`. All container host virtual machines will appear in the output of `lambda-machine-local ls`.

<a name="starting_and_stopping_container_host_virtual_machine"></a>
#### Starting and stopping container host virtual machine

You can stop a container host virtual machine with `lambda-machine-local stop`
command. An existing container host virtual machine host can be started with
`lambda-machine-local start` command.

```console
$ lambda-machine-local stop default

$ lambda-machine-local start default
```

<a name="working_with_container_host_virtual_machine_without_specifying_name"></a>
#### Working with container host virtual machine without specifying name

If no container host virtual machine name is specified, some
`lambda-machine-local` commands will assume that the command needs to be run on
a virtual machine named `default`. Because naming container host virtual machine
is `default` is such a common pattern, this allows you to save a few strokes on
some of the most frequently used `lambda-machine-local` commands.

```console
$ lambda-machine-local stop
Stopping "default"...
Machine "default" was stopped.

$ lambda-machine-local start
Starting "default"...
(default) Check network to re-create if needed...
(default) Waiting for an IP...
Machine "default" was started.
Waiting for SSH to be available...
Detecting the provisioner...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Started machines may have new IP addresses. You may need to re-run the `lambda-machine-local env` command.
```

Following `lambda-machine-local` commands follow this pattern &ndash;

  * `lambda-machine-local config`
  * `lambda-machine-local env`
  * `lambda-machine-local inspect`
  * `lambda-machine-local ip`
  * `lambda-machine-local kill`
  * `lambda-machine-local provision`
  * `lambda-machine-local regenerate-certs`
  * `lambda-machine-local restart`
  * `lambda-machine-local ssh`
  * `lambda-machine-local start`
  * `lambda-machine-local status`
  * `lambda-machine-local stop`
  * `lambda-machine-local url`

<a name="auto_updates_and_crash_reporting"></a>
#### Auto updates and crash reporting

Like Amazon Linux, Lambda Linux VirtualBox flavor is a rolling release and is
continuously updated for stability and security. When you create a new container
host virtual machine, we check on GitHub to see if you are running the latest
version of Lambda Linux VirtualBox flavor. If your local ISO image is outdated,
we download the newer version and use that to create container host.

Lambda Machine Local also has a built-in crash reporting mechanism. Only in the
event of a crash, it allows us to collect information about your Lambda Machine
Local version, build, operating system, architecture, path to your current
shell, VirtualBox logs and the history of last command as seen with `--debug`
option. This data will help us debug issues with `lambda-machine-local` and
provide an appropriate fix.

If you wish to opt out of crash reporting, you can create a `no-error-report`
file in your `$HOME/.docker/lambda-machine-local` directory and Lambda Machine
Local will disable this behavior. Leaving the file empty is fine. Lambda Machine
Local only checks for its presence.

```console
$ mkdir -p ~/.docker/lambda-machine-local && touch ~/.docker/lambda-machine-local/no-error-report
```

You can also setup your own [BugSnag](https://www.bugsnag.com/) account and use
that to collect the crash report and send it to us when reporting an issue.
Please see `--bugsnag-api-token` command line flag for more details.

<a name="command_line_reference"></a>
### Command line reference

Lambda Machine Local CLI has the following sub-commands.

  * [active](#cli_active)
  * [config](#cli_config)
  * [create](#cli_create)
  * [env](#cli_env)
  * [help](#cli_help)
  * [inspect](#cli_inspect)
  * [ip](#cli_ip)
  * [kill](#cli_kill)
  * [ls](#cli_ls)
  * [provision](#cli_provision)
  * [regenerate-certs](#cli_regenerate_certs)
  * [restart](#cli_restart)
  * [rm](#cli_rm)
  * [scp](#cli_scp)
  * [ssh](#cli_ssh)
  * [start](#cli_start)
  * [status](#cli_status)
  * [stop](#cli_stop)
  * [url](#cli_url)

<a name="cli_active"></a>
#### active

List the local container host virtual machine that is _active_. A virtual
machine is considered active if the `DOCKER_HOST` environment variable points to
it.

```console
$ lambda-machine-local ls
NAME      ACTIVE   DRIVER       STATE     URL                         DOCKER    ERRORS
default   -        virtualbox   Running   tcp://192.168.99.103:2376   v1.12.6
etcd1     -        virtualbox   Running   tcp://192.168.99.100:2376   v1.12.6
etcd2     -        virtualbox   Running   tcp://192.168.99.101:2376   v1.12.6
etcd3     *        virtualbox   Running   tcp://192.168.99.102:2376   v1.12.6

$ echo $DOCKER_HOST
tcp://192.168.99.102:2376

$ lambda-machine-local active
etcd3
```

<a name="cli_config"></a>
#### config

```console
Usage: lambda-machine-local config [arg...]

Print the connection config for machine

Description:
   Argument is a machine name.
```

For example &ndash;

```console
$ lambda-machine-local config dev
--tlsverify
--tlscacert="/Users/username/.docker/machine/machines/dev/ca.pem"
--tlscert="/Users/username/.docker/machine/machines/dev/cert.pem"
--tlskey="/Users/username/.docker/machine/machines/dev/key.pem"
-H=tcp://192.168.99.100:2376
```

<a name="cli_create"></a>
#### create

Create a local container host virtual machine.

```console
$ lambda-machine-local create dev
Running pre-create checks...
(dev) Copying /Users/username/.docker/lambda-machine-local/cache/lambda-linux-vbox.iso to
/Users/username/.docker/lambda-machine-local/machines/dev/lambda-linux-vbox.iso...
(dev) Creating VirtualBox VM...
(dev) Creating SSH key...
(dev) Starting the VM...
(dev) Check network to re-create if needed...
(dev) Waiting for an IP...
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with lambda...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine,
run: lambda-machine-local env dev
```

**Note:** Libmachine has the concept of _driver_ which is used to manage
hypervisors or cloud instances. In Lambda Machine Local we use VirtualBox as our
default driver. It is also the only driver supported at this time. When
deploying containers to the cloud, we recommend that you use cloud provider
supported tools such as AWS CLI and ECS CLI.

```console
Usage: lambda-machine-local create [OPTIONS] [arg...]

Create a machine

Description:
   Run 'lambda-machine-local create --driver name' to include the create flags for that driver in the help text.

Options:

   --driver, -d "virtualbox"                    Driver to create machine with. [$MACHINE_DRIVER]
   --virtualbox-cpu-count "1"                   number of CPUs for the machine (-1 to use the number of CPUs
                                                available) [$VIRTUALBOX_CPU_COUNT]
   --virtualbox-disk-size "20000"               Size of disk for host in MB [$VIRTUALBOX_DISK_SIZE]
   --virtualbox-host-dns-resolver               Use the host DNS resolver [$VIRTUALBOX_HOST_DNS_RESOLVER]
   --virtualbox-hostonly-cidr "192.168.99.1/24" Specify the Host Only CIDR [$VIRTUALBOX_HOSTONLY_CIDR]
   --virtualbox-hostonly-nicpromisc "deny"      Specify the Host Only Network Adapter Promiscuous Mode
                                                [$VIRTUALBOX_HOSTONLY_NIC_PROMISC]
   --virtualbox-hostonly-nictype "82540EM"      Specify the Host Only Network Adapter Type
                                                [$VIRTUALBOX_HOSTONLY_NIC_TYPE]
   --virtualbox-hostonly-no-dhcp                Disable the Host Only DHCP Server
                                                [$VIRTUALBOX_HOSTONLY_NO_DHCP]
   --virtualbox-import-lambda-linux-vbox-vm     The name of a Lambda Linux VirtualBox VM to import
                                                [$VIRTUALBOX_LAMBDA_LINUX_VBOX_IMPORT_VM]
   --virtualbox-lambda-linux-vbox-url           The URL of the Lambda Linux VirtualBox image. Defaults to the latest
                                                available version [$VIRTUALBOX_LAMBDA_LINUX_VBOX_URL]
   --virtualbox-memory "2048"                   Size of memory for host in MB [$VIRTUALBOX_MEMORY_SIZE]
   --virtualbox-nat-nictype "82540EM"           Specify the Network Adapter Type [$VIRTUALBOX_NAT_NICTYPE]
   --virtualbox-no-dns-proxy                    Disable proxying all DNS requests to the host
                                                [$VIRTUALBOX_NO_DNS_PROXY]
   --virtualbox-no-share                        Disable the mount of your home directory
                                                [$VIRTUALBOX_NO_SHARE]
   --virtualbox-no-vtx-check                    Disable checking for the availability of hardware virtualization
                                                before the vm is started [$VIRTUALBOX_NO_VTX_CHECK]
   --virtualbox-share-folder                    Mount the specified directory instead of the default home location.
                                                Format: dir:name [$VIRTUALBOX_SHARE_FOLDER]
   --virtualbox-ui-type "headless"              Specify the UI Type: (gui|sdl|headless|separate)
                                                [$VIRTUALBOX_UI_TYPE]
```

<a name="cli_env"></a>
#### env

Set environment variables to configure docker client to run commands on a specific
container host.

```console
Usage: lambda-machine-local env [OPTIONS] [arg...]

Display the commands to set up the environment for the Docker client

Description:
   Argument is a machine name.

Options:

   --shell      Force environment to be configured for a specified shell: [fish, cmd, powershell, tcsh], default is
                auto-detect
   --unset, -u  Unset variables instead of setting them
   --no-proxy   Add machine IP to NO_PROXY environment variable
```

`lambda-machine-local env <LOCAL_VIRTUAL_MACHINE_NAME>` will print out `export`
commands which can be run in a sub-shell. Running `lambda-machine-local env -u`
will print out corresponding `unset` commands.

```console
$ env | grep DOCKER

$ eval "$(lambda-machine-local env dev)"

$ env | grep DOCKER
DOCKER_TLS_VERIFY=1
DOCKER_HOST=tcp://192.168.99.100:2376
DOCKER_CERT_PATH=/Users/username/.docker/lambda-machine-local/machines/dev
DOCKER_MACHINE_NAME=dev

$ # If you run a docker command, it will now run against "dev" container host

$ eval "$(lambda-machine-local env -u)"

$ env | grep DOCKER

$ # The environment variables have been unset
```

The output above is for `bash` and `zsh` shells. Lambda Machine Local can
auto-detect your environment and print appropriate commands for other shells
such as `cmd`, `fish`, `powershell` and `emacs`.

You can also use `--shell` flag to specify the shell for `lambda-machine-local
env` command.

```console
C:\Users\username\bin>lambda-machine-local.exe env --shell cmd dev
SET DOCKER_TLS_VERIFY=1
SET DOCKER_HOST=tcp://192.168.99.100:2376
SET DOCKER_CERT_PATH=C:\Users\username\.docker\lambda-machine-local\machines\dev
SET DOCKER_MACHINE_NAME=dev
SET COMPOSE_CONVERT_WINDOWS_PATHS=true
REM Run this command to configure your shell:
REM     @FOR /f "tokens=*" %i IN ('lambda-machine-local.exe env --shell cmd dev') DO @%i
```

##### Excluding the created container host virtual machines from proxies

The `env` command supports `--no-proxy` flag which outputs configuration to set
`NO_PROXY` environment variable. This is useful in network environments where a
HTTP proxy is required for Internet access.

```console
$ lambda-machine-local env --no-proxy dev
export DOCKER_TLS_VERIFY="1"
export DOCKER_HOST="tcp://192.168.99.100:2376"
export DOCKER_CERT_PATH="/Users/username/.docker/lambda-machine-local/machines/dev"
export DOCKER_MACHINE_NAME="dev"
export NO_PROXY="192.168.99.100"
# Run this command to configure your shell:
# eval $(lambda-machine-local env --no-proxy dev)
```

**Note:** When using `--no-proxy` flag, you might also need to
configure
[`HTTPS_PROXY`](https://docs.docker.com/v1.12/engine/reference/commandline/dockerd/#/running-a-docker-daemon-behind-an-httpsproxy) for
docker engine. Please
see
[docker configuration files](#customizing_lambda_linux_virtualbox_flavor_docker_configuration_files) for
details on how to customize docker engine running on container host virtual
machine.

<a name="cli_help"></a>
#### help

```console
Usage: lambda-machine-local help [arg...]

Shows a list of commands or help for one command
```

This command has the form `lambda-machine-local help <COMMAND>`. For example &ndash;

```console
Usage: lambda-machine-local config [arg...]

Print the connection config for machine

Description:
   Argument is a machine name.
```

<a name="cli_inspect"></a>
#### inspect

```console
Usage: lambda-machine-local inspect [OPTIONS] [arg...]

Inspect information about a machine

Description:
   Argument is a machine name.

   Options:

   --format, -f         Format the output using the given go template.
```

By default inspect will render information about local container host virtual
machine in JSON format. If a format is specified, the given template will be
executed for each result.
Go [`text/template`](http://golang.org/pkg/text/template/) package describes all
the details of the format accepted by `--format` flag.

##### List all the details of a local container host virtual machine

This is the default usage of `inspect`.

```console
$ lambda-machine-local inspect dev
{
    "ConfigVersion": 3,
    "Driver": {
        "IPAddress": "192.168.99.100",
        "MachineName": "dev",
        "SSHUser": "ll-user",
        "SSHPort": 61233,
        "SSHKeyPath": "/Users/username/.docker/lambda-machine-local/machines/dev/id_rsa",
        "StorePath": "/Users/username/.docker/lambda-machine-local",

        [...]

    },

    [...]
}
```

##### Get container host virtual machine's IP address

You can pick any field from JSON in a fairly straightforward manner. In this
example we will extract the IP Address.

```console
$ lambda-machine-local inspect --format='{{.Driver.IPAddress}}' dev
192.168.99.100
```

##### Formatting

To extract a subset of information in unformatted JSON, you can use `json`
function in the template.

```console
$ lambda-machine-local inspect --format='{{json .Driver}}' dev
{"Boot2DockerImportVM":"","Boot2DockerURL":"","CPU":1,"DNSProxy":true,"DiskSize":20000,"GithubAPIToken":"",
"HostDNSResolver":false,"HostInterfaces":{},"HostOnlyCIDR":"192.168.99.1/24","HostOnlyNicType":"82540EM",
"HostOnlyNoDHCP":false,"HostOnlyPromiscMode":"deny","IPAddress":"192.168.99.100","MachineName":"dev","Memory":2048,
"NatNicType":"82540EM","NoShare":false,"NoVTXCheck":false,"SSHKeyPath":
"/Users/username/.docker/lambda-machine-local/machines/dev/id_rsa","SSHPort":61233,"SSHUser":"ll-user", ...
```

For a more human-readable format, use `prettyjson` instead.

```console
$ lambda-machine-local inspect --format='{{prettyjson .Driver}}' dev
{
    "Boot2DockerImportVM": "",
    "Boot2DockerURL": "",
    "CPU": 1,
    "DNSProxy": true,
    "DiskSize": 20000,
    "GithubAPIToken": "",

    [...]

}
```

<a name="cli_ip"></a>
#### ip

Get the IP address of one or more container host virtual machines.

```console
$ lambda-machine-local ip dev1
192.168.99.100

$ lambda-machine-local ip dev1 dev2
192.168.99.101
192.168.99.100
```

<a name="cli_kill"></a>
#### kill

```console
Usage: lambda-machine-local kill [arg...]

Kill a machine

Description:
   Argument(s) are one or more machine names.
```

For example &ndash;

```console
$ lambda-machine-local ls
NAME   ACTIVE   DRIVER       STATE     URL                         DOCKER    ERRORS
dev    -        virtualbox   Running   tcp://192.168.99.100:2376   v1.12.6

$ lambda-machine-local kill dev
Killing "dev"...
Machine "dev" was killed.

$ lambda-machine-local ls
NAME   ACTIVE   DRIVER       STATE     URL   DOCKER    ERRORS
dev    -        virtualbox   Stopped         Unknown
```

<a name="cli_ls"></a>
#### ls

```console
Usage: lambda-machine-local ls [OPTIONS] [arg...]

List machines

Options:

   --quiet, -q                                  Enable quiet mode
   --filter [--filter option --filter option]   Filter output based on conditions provided
   --timeout, -t "10"                           Timeout in seconds, default to 10s
   --format, -f                                 Pretty-print machines using a Go template
```

##### Timeout

The `ls` command tries to reach each container host virtual machines in
parallel. If a given container host does not answer in less than 10 seconds, the
`ls` command will state that the host is in `Timeout` state. In some
circumstances (connection issues, high load, or while troubleshooting), you may
want to increase or decrease this value. You can use the `--timeout` flag for
this purpose with a numerical value in seconds.

```console
$ lambda-machine-local ls -t 12
NAME      ACTIVE   DRIVER       STATE     URL                         DOCKER    ERRORS
default   -        virtualbox   Running   tcp://192.168.99.100:2376   v1.12.6
```

##### Filtering

The filtering flag (`--filter`) format is a `key=value` pair. If there is more
than one filter, then pass multiple flags (`--filter "foo=bar" --filter
"bif=baz"`).

Currently supported filters are &ndash;

  * driver
  * state (`Running|Paused|Saved|Stopped|Stopping|Starting|Error`)
  * name (supports [golang style](https://github.com/google/re2/wiki/Syntax) regular expressions)

```console
$ lambda-machine-local ls
NAME      ACTIVE   DRIVER       STATE     URL                         DOCKER    ERRORS
default   -        virtualbox   Stopped                               Unknown
foo0      -        virtualbox   Running   tcp://192.168.99.101:2376   v1.12.6
foo1      -        virtualbox   Running   tcp://192.168.99.102:2376   v1.12.6
foo2      *        virtualbox   Running   tcp://192.168.99.103:2376   v1.12.6

$ lambda-machine-local ls -filter name=foo0
NAME   ACTIVE   DRIVER       STATE     URL                         DOCKER    ERRORS
foo0   -        virtualbox   Running   tcp://192.168.99.101:2376   v1.12.6

$ lambda-machine-local ls --filter driver=virtualbox --filter state=Stopped
NAME      ACTIVE   DRIVER       STATE     URL   DOCKER    ERRORS
default   -        virtualbox   Stopped         Unknown
```

##### Formatting

The formatting option (`--format`) will pretty-print container host virtual
machines using a Go template.

Valid placeholders for the Go template are listed below.

| Placeholder | Description |
| ----------- | ----------- |
| .Name | Container host virtual machine name |
| .Active | Is the machine active? |
| .DriverName | Driver name |
| .State | Container host virtual machine state (running, stopped, ...) |
| .URL | Container host virtual machine URL |
| .Error | Container host virtual machine error |
| .DockerVersion | Container host virtual machine docker daemon version |
| .ResponseTime | Time taken by container host virtual machine to respond |

When using the `--format` option, the `ls` command will output the data as the
template declares. You can use table directive to include column headers as
well.

The following example uses a template without headers and outputs the `Name` and
`State` entries separated by a colon.

```console
$ lambda-machine-local ls --format "{{.Name}}: {{.State}}"
default: Stopped
foo0: Running
foo1: Running
foo2: Running
```

To list the same information in table format you can use.

```console
$ lambda-machine-local ls --format "table {{.Name}}\t{{.State}}"
NAME      STATE
default   Stopped
foo0      Running
foo1      Running
foo2      Running
```

<a name="cli_provision"></a>
#### provision

Re-run the provisioning step on an existing container host virtual machine.

Sometimes it might be useful to re-run the provisioning step. Reasons for doing
so might include a failure of the original provisioning process.

Usage is `lambda-machine-local provision <LOCAL_VIRTUAL_MACHINE_NAME>`. Multiple
names may be specified.

```console
$ lambda-machine-local provision foo0 foo1
Waiting for SSH to be available...
Waiting for SSH to be available...
Detecting the provisioner...
Detecting the provisioner...
Copying certs to the local machine directory...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Copying certs to the remote machine...
```

<a name="cli_regenerate_certs"></a>
#### regenerate-certs

```console
Usage: lambda-machine-local regenerate-certs [OPTIONS] [arg...]

Regenerate TLS Certificates for a machine

Description:
   Argument(s) are one or more machine names.

   Options:

   --force, -f  Force rebuild and do not prompt
```

Regenerate TLS certificates and update container host virtual machine with new
certificates.

```console
$ lambda-machine-local regenerate-certs dev
Regenerate TLS machine certs?  Warning: this is irreversible. (y/n): y
Regenerating TLS certificates
Waiting for SSH to be available...
Detecting the provisioner...
Copying certs to the local machine directory...
Copying certs to the remote machine...
```

<a name="cli_restart"></a>
#### restart

```console
Usage: lambda-machine-local restart [arg...]

Restart a machine

Description:
   Argument(s) are one or more machine names.
```

Restart a container host virtual machine. This is equivalent to `lambda-machine-local stop` followed by `lambda-machine-local start`.

```console
$ lambda-machine-local restart dev
Restarting "dev"...
Stopping "dev"...
Machine "dev" was stopped.
Starting "dev"...
(dev) Check network to re-create if needed...
(dev) Waiting for an IP...
Machine "dev" was started.
Waiting for SSH to be available...
Detecting the provisioner...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Waiting for SSH to be available...
Detecting the provisioner...
Restarted machines may have new IP addresses. You may need to re-run the `lambda-machine-local env` command.
```

<a name="cli_rm"></a>
#### rm

Remove a container host virtual machine. This will remove local reference and
also remove it from VirtualBox.

```console
Usage: lambda-machine-local rm [OPTIONS] [arg...]

Remove a machine

Description:
   Argument(s) are one or more machine names.

Options:

   --force, -f  Remove local configuration even if machine cannot be removed, also implies an automatic yes (`-y`)
   -y           Assumes automatic yes to proceed with remove, without prompting further user confirmation
```

<a name="cli_scp"></a>
#### scp

Copy files between local host to virtual machine or between virtual machines using scp.

The notation is `virtualmachinename:/path/to/files` for the arguments. In case
of local host, you can just specify the path.

```console
$ cat foo.txt
cat: foo.txt: No such file or directory

$ lambda-machine-local ssh dev pwd
/home/ll-user

$ lambda-machine-local ssh dev 'echo A file created remotely! >foo.txt'

$ lambda-machine-local scp dev:/home/ll-user/foo.txt .
foo.txt                                                           100%   25     0.0KB/s   00:00

$ cat foo.txt
A file created remotely!
```

Just like how `scp` has a `-r` flag for copying files recursively,
`lambda-machine-local` has a corresponding `-r` flag.

When transferring files between virtual machines, it goes through the local
host's filesystem first (using `scp's -3` flag).

`lambda-machine-local scp` has `-d` flag which uses `rsync` to transfer deltas
instead of transferring all the files. This can be very useful when transferring
large files or updating directories with a lot of files.

When transferring directories and not just files, you can avoid `rsync`
surprises by using trailing slashes `(/)` on both source and destination. For
example &ndash;

```console
$ mkdir -p bar
$ touch bar/baz
$ lambda-machine-local scp -r -d bar/ dev:/var/lib/lambda-machine-local/home/ll-user/bar/
$ lambda-machine-local ssh dev ls /var/lib/lambda-machine-local/home/ll-user/bar/
baz
```

<a name="cli_ssh"></a>
#### ssh

Log into or run a command on a container host virtual machine using SSH.

To login, just run `lambda-machine-local ssh virtualmachinename`.

```console
$ lambda-machine-local ssh dev
Last login: Fri Apr 28 13:17:35 2017

*Thank you* for using
  _                  _        _         _     _
 | |    __ _  _ __  | |__  __| | __ _  | |   (_) _ _  _  _ __ __
 | |__ / _` || '  \ | '_ \/ _` |/ _` | | |__ | || ' \| || |\ \ /
 |____|\__,_||_|_|_||_.__/\__,_|\__,_| |____||_||_||_|\_,_|/_\_\

VirtualBox Release: 2017.03 | Twitter: @lambda_linux

[ll-user@ip-192-168-99-100 ~]$
```

You can also specify commands to run on container host virtual machine by
appending them directly to the `lambda-machine-local ssh` command, much like the
regular `ssh` program works.

```console
$ lambda-machine-local ssh dev free
             total       used       free     shared    buffers     cached
Mem:       1020432     953392      67040     852028       3244     881688
-/+ buffers/cache:      68460     951972
Swap:            0          0          0

$ lambda-machine-local ssh dev df -h
Filesystem      Size  Used Avail Use% Mounted on
tmpfs           499M     0  499M   0% /dev/shm
/dev/sda1       3.8G  8.1M  3.6G   1% /var/lib/lambda-machine-local
/Users          465G  214G  252G  46% /Users
```

##### Different types of SSH

When Lambda Machine Local is invoked, it will check for `ssh` binary and will
attempt to use that for SSH commands it needs to run. If it cannot find `ssh`
binary it will default to using a native Go implementation
from [`crypto/ssh`](https://godoc.org/golang.org/x/crypto/ssh). This is useful
on Windows when used without [MSYS2](http://www.msys2.org/)
or [Git Bash](https://git-for-windows.github.io/).

If you want to force Go version of SSH, you can use `--native-ssh` flag.

```console
$ lambda-machine-local --native-ssh ssh dev
```

<a name="cli_start"></a>
#### start

Start a container host virtual machine.

```console
Usage: lambda-machine-local start [arg...]

Start a machine

Description:
   Argument(s) are one or more machine names.
```

For example &ndash;

```console
$ lambda-machine-local start dev
Starting "dev"...
(dev) Check network to re-create if needed...
(dev) Waiting for an IP...
Machine "dev" was started.
Waiting for SSH to be available...
Detecting the provisioner...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Started machines may have new IP addresses. You may need to re-run the `lambda-machine-local env` command.
```

<a name="cli_status"></a>
#### status

```console
Usage: lambda-machine-local status [arg...]

Get the status of a machine

Description:
   Argument is a machine name.
```

For example &ndash;

```console
$ lambda-machine-local status dev
Running
```

<a name="cli_stop"></a>
#### stop

```console
Usage: lambda-machine-local stop [arg...]

Stop a machine

Description:
   Argument(s) are one or more machine names.
```

For example &ndash;

```console
$ lambda-machine-local ls
NAME   ACTIVE   DRIVER       STATE     URL                         DOCKER    ERRORS
dev    *        virtualbox   Running   tcp://192.168.99.100:2376   v1.12.6

$ lambda-machine-local stop dev
Stopping "dev"...
Machine "dev" was stopped.

$ lambda-machine-local ls
NAME   ACTIVE   DRIVER       STATE     URL   DOCKER    ERRORS
dev    -        virtualbox   Stopped         Unknown
```

<a name="cli_url"></a>
#### url

```console
Usage: lambda-machine-local stop [arg...]

Get the URL of a machine

Description:
   Argument is a machine name.
```

For example &ndash;

```console
$ lambda-machine-local url dev
tcp://192.168.99.100:2376
```

<a name="customizing_lambda_linux_virtualbox_flavor"></a>
## Customizing Lambda Linux VirtualBox flavor

Lambda Machine Local is designed to be a stable and secure **local development
primitive**. Our thinking behind primitives is very similar to **Primitives not
frameworks** lesson
from
[10 Lessons from 10 Years of Amazon Web Services](http://www.allthingsdistributed.com/2016/03/10-lessons-from-10-years-of-aws.html).

As such, we cannot anticipate all the use-cases where you might find using
Lambda Machine Local useful. Therefore we would like to provide you with
mechanisms to customize Lambda Linux VirtualBox flavor for your workflow needs.

Before we can describe the customization mechanisms available, we would like to
provide a high-level overview of how Lambda Machine Local and Lambda Linux
VirtualBox flavor works.

Lambda Machine Local consists of two components &ndash;

  * `lambda-machine-local` Command-line interface

    `lambda-machine-local` is responsible for managing the life-cycle of container host virtual machine. It is maintained in this repository.

  * `lambda-linux-vbox.iso` Bootable ISO image

    `lambda-linux-vbox.iso` is the operating system image that `lambda-machine-local` uses to create container host virtual machine. It is maintained in a separate [Lambda Linux VirtualBox flavor](https://github.com/lambda-linux/lambda-linux-vbox) GitHub repository.

<a name="customizing_lambda_linux_virtualbox_flavor_home_disk_vmdk"></a>
When a new container host virtual machine is created by `lambda-machine-local`
it sets `lambda-machine-vbox.iso` image as its primary boot device. In addition
a sparse disk image (`disk.vmdk`) is created to store persistent information.
The size of this disk can be adjusted using `--virtualbox-disk-size` flag
in [`lambda-machine-local create`](#cli_create) command.

You can find the ISO and disk image used by virtual machine at the following locations &ndash;

  * `~/username/.docker/lambda-machine-local/machines/virtualmachinename/lambda-linux-vbox.iso`
  * `~/username/.docker/lambda-machine-local/machines/virtualmachinename/disk.vmdk`

Upon booting Lambda Linux VirtualBox flavor
starts
[`/usr/bin/lml_init`](https://github.com/lambda-linux/lambda-linux-vbox/blob/master/install-root/usr-bin-lml_init).
`lml_init` first checks to see if `disk.vmdk` is already setup. If so, it
proceeds with the rest of the initialization process. If `disk.vmdk` is not
setup then `lml_init` partitions and formats the disk as follows.

`lml_init` creates
a
[GUID partition table (GPT)](https://en.wikipedia.org/wiki/GUID_Partition_Table)
with two partitions.

```console
[ll-user@ip-192-168-99-100 ~]$ sudo sgdisk -p /dev/sda

[...]

Number  Start (sector)    End (sector)  Size       Code  Name
   1            2048         8192000   3.9 GiB     8300  Linux
   2         8194048        40959966   15.6 GiB    8E00  Linux LVM
```

The first partition (`/dev/sda1`) is a standard Linux partition formatted with
an ext4 filesystem and mounted at `/var/lib/lambda-machine-local` . The second
partition (`/dev/sda2`) is configured as
a [Logical Volume Management](http://tldp.org/HOWTO/LVM-HOWTO/) device. The LVM
device is directly accessed by docker via its `devicemapper` storage backend.

```console
[ll-user@ip-192-168-99-100 ~]$ mount | grep sda
/dev/sda1 on /var/lib/lambda-machine-local type ext4 (rw)

[ll-user@ip-192-168-99-100 ~]$ sudo pvs
  PV         VG     Fmt  Attr PSize  PFree
  /dev/sda2  docker lvm2 a--  15.62g 144.00m

[ll-user@ip-192-168-99-100 ~]$ sudo vgs
  VG     #PV #LV #SN Attr   VSize  VFree
  docker   1   1   0 wz--n- 15.62g 144.00m

[ll-user@ip-192-168-99-100 ~]$ sudo lvs
  LV          VG     Attr       LSize  Pool Origin Data%  Meta%  Move Log Cpy%Sync Convert
  docker-pool docker twi-a-t--- 15.45g             1.88   0.37

[ll-user@ip-192-168-99-100 ~]$ docker info | grep "Data Space"
 Data Space Used: 312 MB
 Data Space Total: 16.59 GB
 Data Space Available: 16.28 GB

[ll-user@ip-192-168-99-100 ~]$ docker info | grep -A 4 "Storage Driver"
Storage Driver: devicemapper
 Pool Name: docker-docker--pool
 Pool Blocksize: 524.3 kB
 Base Device Size: 10.74 GB
 Backing Filesystem: ext4
```

All changes made via `docker` commands and changes to
`/var/lib/lambda-machine-local` are preserved across virtual machine reboots.
For your convenience, `/var/lib/lambda-machine-local` is organized as follows.

<a name="customizing_lambda_linux_virtualbox_flavor_home_ll_user"></a>
```console
[ll-user@ip-192-168-99-100 ~]$ sudo tree -L 2 /var/lib/lambda-machine-local/
/var/lib/lambda-machine-local/
|-- bootlocal.sh
|-- bootsync.sh
|-- etc
|   `-- sysconfig
|-- home
|   `-- ll-user
|-- lost+found
|-- root
`-- var
    `-- lib
```

We recommend that use `/var/lib/lambda-machine-local/home/ll-user` for your
local changes. This directory is configured with a `UID:GID` pair of `500:500`
which is the same `UID:GID` pair that we use to mount VirtualBox shared folders.
This allows you to easily move files between your host filesystem and
`/var/lib/lambda-machine-local/home/ll-user` directory.

```console
[ll-user@ip-192-168-99-100 ~]$ id ll-user
uid=500(ll-user) gid=500(ll-user) groups=500(ll-user),10(wheel),497(docker)

[ll-user@ip-192-168-99-100 ~]$ ls -la /var/lib/lambda-machine-local/home | grep ll-user
drwx------ 3 ll-user ll-user 4096 Mar 28 17:04 ll-user

[ll-user@ip-192-168-99-100 ~]$ mount | grep 500
/Users on /Users type vboxsf (uid=500,gid=500,iocharset=utf8,rw)
```

`lml_init` creates two customization files that it uses during subsequent
invocations. These are -

  * `/var/lib/lambda-machine-local/bootlocal.sh`

  * `/var/lib/lambda-machine-local/bootsync.sh`

#### `bootlocal.sh`

`lml_init` executes `bootlocal.sh` at the end of its initialization process. For
most users we recommend using this file for your customization.

If `bootlocal.sh` exits with an non-zero return status, `lml_init` will log a
message to `/var/log/messages`.

#### `bootsync.sh`

`bootsync.sh` is the first customizable hook executed by `lml_init`.
`bootsync.sh` is executed prior to starting the docker and ssh daemons.

**Note:** `lml_init` expects `bootsync.sh` to exit with posix return status of
zero. If a non-zero return status is encountered, then `lml_init` will **not**
proceed further with its initialization process. Therefore when customizing this
file please ensure that any errors in your program are gracefully handled.

#### Developing and debugging customization scripts

Developing and debugging customization scripts can sometimes get tricky. We
recommend that you use `--virtualbox-ui-type gui` flag when developing or
debugging your customization script.

With this flag, the virtual machine starts with a GUI console. You will be able
to log into GUI console using the user name `ll-user` and password `ll-user`.
You can then `sudo -i` and investigate the behavior of your customization
scripts.

<a name="customizing_lambda_linux_virtualbox_flavor_docker_configuration_files"></a>
#### Docker configuration files

Docker configuration files are maintained at
`/var/lib/lambda-machine-local/etc/sysconfig`. If you need to add additional
options to your docker daemon, you can edit these files.

```console
[ll-user@ip-192-168-99-100 ~]$ sudo tree /var/lib/lambda-machine-local/etc/sysconfig/
/var/lib/lambda-machine-local/etc/sysconfig/
|-- docker
|-- docker-storage
`-- docker-storage-setup

0 directories, 3 files
```

<a name="a_note_about_virtualbox_shared_folder_support"></a>
## A note about VirtualBox shared folder support

VirtualBox is a widely used, stable and
mature
[Type-2 hypervisor](https://en.wikipedia.org/wiki/Hypervisor#Classification).
Like most stable software VirtualBox has its share of known issues. The good
part of relying on mature hypervisor technology is that known issues also have
well-known workarounds.

There are three related VirtualBox issues that we would like to draw your
attention to, and describe workarounds. These issues are usually encountered
when using shared folder support as a part of _edit-compile-run_ cycle.

  * `sendfile` serving [cached file contents](https://www.virtualbox.org/ticket/9069)
  *  Shared folder [I/O performance](http://mitchellh.com/comparing-filesystem-performance-in-virtual-machines)
  * `vboxsf` does not support [inode notify (inotify)](https://www.virtualbox.org/ticket/10660)

In the following section we will describe two workarounds along with features
that we have in Lambda Machine Local that make it easy for you adopt the
suggested workarounds.

<a name="a_note_about_virtualbox_shared_folder_support_workaround_1"></a>
#### Workaround #1

In Lambda Machine Local, we support using Lambda Linux container host VM as a
development machine.

If your workflow allows you to do all your development inside a VM or a
container, then you can easily avoid VirtualBox shared folder issues. Please
see [workaround #2](#a_note_about_virtualbox_shared_folder_support_workaround_2)
in case you are unable to move all your development environment dependencies
into a container.

As mentioned
in
[customizing Lambda Linux VirtualBox flavor](#customizing_lambda_linux_virtualbox_flavor_home_ll_user),
on the container host virtual machine, we create
`/var/lib/lambda-machine-local/home/ll-user` for your local changes. You can use
this directory to store your code and other assets.

Then create containers as required by your workflow. One or more of these
containers would be your _development_ containers. _Development_ containers
would contain your editor and other tools that you would use for your
development.

When you need to share files between your development container and deployment
container you can bind mount relevant parts of
`/var/lib/lambda-machine-local/home/ll-user` directory using `--volumes / -v`
flag with `docker run` command.

`/var/lib/lambda-machine-local/home/ll-user` is backed
by [`disk.vmdk`](#customizing_lambda_linux_virtualbox_flavor_home_disk_vmdk)
sparse disk image. The filesystem partition backing
`/var/lib/lambda-machine-local/home/ll-user` is allocated 20% of `disk.vmdk`
disk image. By default the size of this partition is 3.9GiB and size of
`disk.vmdk` sparse disk image is 20GB. If your workflow requires more than
3.9GiB for `/var/lib/lambda-machine-local/home/ll-user`, you can use
`--virtualbox-disk-size` flag in [`lambda-machine-local create`](#cli_create)
command to increase the size of `disk.vmdk`.

When using development containers, you might want to consider allocating more
CPU and memory to the underlying container host virtual machine. By default 1
CPU and 2GiB of memory is allocated to each container host virtual machine. You
can use `--virtualbox-cpu-count` and `--virtualbox-memory` flag
in [`lambda-machine-local create`](#cli_create) to increase these defaults.

**Note:** When you run [`lambda-machine-local rm`](#cli_rm) command, it also
deletes `disk.vmdk` that is associated with the container host virtual machine.
Therefore, please ensure that you have backed up your code and assets stored on
`/var/lib/lambda-machine-local/home/ll-user` prior to running
`lambda-machine-local rm` command. [`lambda-machine-local stop`](#cli_stop)
and [`lambda-machine-local kill`](#cli_kill) commands does not affect
`disk.vmdk`.

<a name="a_note_about_virtualbox_shared_folder_support_workaround_2"></a>
#### Workaround #2

In a lot of workflows, you need to regularly synchronize files between your
VirtualBox host and containers as a part of your _edit-compile-run_ cycle.

Getting this workflow right initially can be tricky, but once you understand the
aspects involved, it is fairly straightforward to script a solution.

There are three aspects involved in this workaround &ndash;

1. Ensure that you have sufficient disk space on container host virtual machine

2. Set up synchronization mechanism

3. Triggering synchronization mechanism based on filesystem events

##### 1. Ensure that you have sufficient disk space on container host virtual machine

As we described
in [workaround #1](#a_note_about_virtualbox_shared_folder_support_workaround_1),
you can use `/var/lib/lambda-machine-local/home/ll-user` directory to store your
code and other assets. In this case you will need to ensure that you have
sufficient space for your synchronized code and assets.

However, if you are planning to be using NFS as your synchronization mechanism
(described next), then you do not need to think about space requirements of
`/var/lib/lambda-machine-local/home/ll-user` directory.

##### 2. Set up synchronization mechanism

There are two well known synchronization mechanisms that you can adopt &ndash;
_Rsync over SSH_ and _NFS_. In Lambda Machine Local, we have support for both
the mechanisms.

1. Rsync over SSH

   Rsync over SSH is more popular synchronization mechanism when compared to
   NFS. It is also easier to setup and works uniformly across Mac, Windows and
   Linux.

   `rsync` is available natively on Linux and Mac. On Windows we recommend
   installing `rsync` using [MSYS2](http://www.msys2.org/). MSYS2 comes with a
   package management tool called `pacman`. You can install `rsync` and
   `openssh` using the following command.

   ```console
   $ pacman -S rsync openssh
   ```

   Lambda Linux VirtualBox flavor also comes pre-installed with `rsync`. There
   is no additional package installation required on container host virtual
   machine.

   ```console
   [ll-user@ip-192-168-99-100 ~]$ which rsync
   /usr/bin/rsync
   ```

   `lambda-machine-local scp` command has `--delta` (`-d`) flag which uses
   `rsync` under the hood. To synchronize files (in the example below `app`
   directory) between your VirtualBox host and guest virtual machine you can
   &ndash;

   ```console
   $ lambda-machine-local scp -r -d app/ default:/var/lib/lambda-machine-local/home/ll-user/app/
   ```

   For most users, we recommend using `lambda-machine-local scp` with `-d` flag.

   However sometimes you might want to run `rsync` manually. You can run `rsync`
   manually as shown below.

   ```console
   $ rsync -av -e 'lambda-machine-local ssh default' --exclude '.git' \
       app :/var/lib/lambda-machine-local/home/ll-user/app
   ```

   There are two things to take note of in the above command.

     *  We explicitly specify virtual machine name when using `-e` or `--rsh`
        flag. In this case we use `lambda-machine-local ssh default` instead of
        the usual `lambda-machine-local ssh`.

     *  `:` is used to automatically infer remote using shell transport, which
        in this case is SSH

   `rsync` is a powerful command with a lot of features. You can `man rsync` for
   more details on how to customize `rsync`. If you would like to specify your
   own `ssh` command then you can find the RSA private key for container host
   virtual machine at
   `~/.docker/lambda-machine-local/machines/virtualmachinename/id_rsa`.

2. NFS

   Setting up NFS is a bit more complicated when compared to Rsync. There is
   also no mature NFS server available on Windows. However, with NFS you don't
   have to worry about disk space requirements on container host virtual
   machine.
   
   First step is to setup NFS server on Mac or Linux. Please refer to the
   respective OS documentation on how to securely export NFS shares. 

   **Note:** When exporting your NFS shares you **should** restrict access to
   the share to specific container host virtual machine IP address or to
   VirtualBox Host Only CIDR. You can get the IP address using the
   command [`lambda-machine-local ip`](#cmd_ip). The default VirtualBox Host
   Only CIDR is `192.168.99.1/24`

   Lambda Linux VirtualBox flavor also comes pre-installed with `nfs-utils`
   package. There is no additional package installation required on container
   host virtual machine.

   ```console
   [ll-user@ip-192-168-99-100 ~]$ rpm -q nfs-utils
   nfs-utils-1.3.0-0.21.ll1.x86_64
   ```

   Once the NFS share has been exported, you can mount it on container host
   virtual machine as follows.

   ```console
   [ll-user@ip-192-168-99-100 ~]$ sudo service rpcbind start
   Starting rpcbind:                                          [  OK  ]

   [ll-user@ip-192-168-99-100 ~]$ sudo mkdir /var/lib/lambda-machine-local/Users-via-nfs

   [ll-user@ip-192-168-99-100 ~]$ sudo mount -t nfs 192.168.99.1:/Users /var/lib/lambda-machine-local/Users-via-nfs

   [ll-user@ip-192-168-99-100 ~]$ mount | grep nfs
   192.168.99.1:/Users on /var/lib/lambda-machine-local/Users-via-nfs type nfs (rw,addr=192.168.99.1)
   ```

##### 3. Triggering synchronization mechanism based on filesystem events

Once the synchronization mechanism is setup either via rsync or NFS, you can
share files with deployment container using bind mount. You can bind mount
relevant parts of `/var/lib/lambda-machine-local/home/ll-user` directory using
`-volumes / -v` flag with `docker run` command.

Next step would be to run actions based on filesystem events. The actions that
you would run is workflow specific. For example in a _edit-compile-run_ cycle
&ndash;

  *  Using an editor on VirtualBox host you save a file
  
  *  This would trigger a filesystem event, based on which `rsync` would get
     executed
     
  *  `rsync` would update the file `/var/lib/lambda-machine-local/home/ll-user`
     directory, which is bind mounted in a container
     
  *  A filesystem event would be triggered inside the container based on which
     we can run the _compile_ action followed by a _run_ action.

This process might look complex in the beginning, but the general approach is
same for any workflow. It is only the actions that would differ between
workflows. For example &ndash; A Node.js front-end project would use a different
set of actions when compared to a Java backend project.

_NFS Caveat:_ When using NFS, filesystem event would not be triggered inside the
container. However you can use what is known as _polling_ mode to execute the
action inside the container.

There are a number of tools that let you monitor filesystem events and trigger
actions. Some popular ones are &ndash;

  * [`watchdog`](https://github.com/gorakhargosh/watchdog) with `watchmedo` command
  * [`guard`](https://github.com/guard/guard)
  * [`fswatch`](https://github.com/emcrisostomo/fswatch)

You can use any of the above tools or something that you are already comfortable
with to build your workflow.
