# ![Lambda Machine Local](logo.png)

Lambda Machine Local provides a stable and secure local container development
environment.

Lambda Machine Local
uses
[Lambda Linux VirtualBox flavor](https://github.com/lambda-linux/lambda-linux-vbox) as
its container host OS and is based on [Alpine Linux](https://alpinelinux.org/)
user space and a customized minimal kernel
from
[Yocto Project](http://www.yoctoproject.org/docs/2.4/kernel-dev/kernel-dev.html#kernel-dev-advanced).

You can download Lambda Machine Local for Mac, Linux or
Windows [here](https://github.com/lambda-linux/lambda-machine-local/releases).

It works a bit like this &ndash;

```console
$ lambda-machine-local create ll-default
Running pre-create checks...
[...]
(ll-default) Creating VirtualBox VM...
(ll-default) Creating SSH key...
(ll-default) Starting the VM...
(ll-default) Check network to re-create if needed...
(ll-default) Waiting for an IP...
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with lambda...
To connect to this virtual machine, run: lambda-machine-local ssh ll-default

$ lambda-machine-local ls
NAME         DRIVER       STATE     ERRORS
ll-default   virtualbox   Running

$ lambda-machine-local ssh ll-default

ip-192-168-99-101:~$ docker run --rm busybox echo hello world
Unable to find image 'busybox:latest' locally
latest: Pulling from library/busybox
0ffadd58f2a6: Pull complete
Digest: sha256:bbc3a03235220b170ba48a157dd097dd1379299370e1ed99ce976df0355d24f0
Status: Downloaded newer image for busybox:latest
hello world
```

-----------------------------------------

<a name="toc"></a>
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
      * [Running Lambda Machine Local alongside Docker Machine](#running_lambda_machine_local_alongside_docker_machine)
    * [Using Lambda Machine Local to run containers](#using_lambda_machine_local_to_run_containers)
      * [Creating a container host virtual machine](#creating_a_container_host_virtual_machine)
      * [Running containers](#running_containers)
      * [Starting and stopping container host virtual machine](#starting_and_stopping_container_host_virtual_machine)
      * [Working with container host virtual machine without specifying name](#working_with_container_host_virtual_machine_without_specifying_name)
      * [Auto updates and crash reporting](#auto_updates_and_crash_reporting)
  * [Command line reference](#command_line_reference)
  * [Customizing Lambda Linux VirtualBox flavor](#customizing_lambda_linux_virtualbox_flavor)
  * [VirtualBox shared folder support](#virtualbox_shared_folder_support)
  * [Additional Information](#additional_information)
    * [Windows Installation](#windows_installation)
    * [Running Lambda Machine Local alongside Docker Machine](#running_lambda_machine_local_alongside_docker_machine)

-----------------------------------------

<a name="what_is_lambda_machine_local"></a>
## What is Lambda Machine Local?

Lambda Machine Local is a tool that provides a local container development
environment with a high bar for stability and security. You can use Lambda
Machine Local to create a supported container host on your Mac, Linux or
Windows.

Using `lambda-machine-local` commands, you can start, inspect, stop, and restart
a managed local container host. By using Lambda Machine Local CLI you can also
ssh into a container host and run `docker` commands on the host.

For example &ndash; You can run `lambda-machine-local ssh default`. Then you can
run `docker ps`, `docker run --rm busybox echo hello world`, and so forth.

<a name="components"></a>
Following is the list of components supported by Lambda Machine Local.

| Component        | Why is it included? / Remarks |
| ---------------- | ------------------- |
| Container OS | [Lambda Linux VirtualBox Flavor](https://github.com/lambda-linux/lambda-linux-vbox) is used as the container host OS |
| Linux Kernel | LTS version [4.9.65](https://cdn.kernel.org/pub/linux/kernel/v4.x/ChangeLog-4.9.65) |
| Docker Engine | Version [17.09.1-ce](https://github.com/moby/moby/releases/tag/v17.09.1-ce)  without Swarm support |
| Storage Driver| [Overlay2](https://docs.docker.com/v17.06/engine/userguide/storagedriver/overlayfs-driver/) with ext4 as backing filesystem |
| VirtualBox Guest Additions | Version [5.1.30](http://download.virtualbox.org/virtualbox/5.1.30/) is included for shared filesystem support |
| Libmachine | Lambda Machine Local uses libmachine version [0.11.0](https://github.com/docker/machine/tree/v0.10.0/libmachine) to manage local container host OS |

<a name="installation"></a>
## Installation

Lambda Machine Local binaries for Mac, Linux and Windows along with installation
instructions and SHA-256 / MD5 checksums is
available [here](https://github.com/lambda-linux/lambda-machine-local/releases).

After creating local container host, you can `lambda-machine-local ssh
<LOCAL_VIRTUAL_MACHINE_NAME>` and run `docker` commands from there.

```console
$ lambda-machine-local ssh ll-default

ip-192-168-99-101:~$ docker info
```

Windows users please see this [link](#windows_installation) for additional
installation information.

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

Additionally if you are looking to run Lambda Machine Local alongside Docker
Machine, then please see
this [link](#running_lambda_machine_local_alongside_docker_machine) for
additional information.

<a name="using_lambda_machine_local_to_run_containers"></a>
### Using Lambda Machine Local to run containers

To run containers, you'll need to &ndash;

  * Create a new container host VM or start an existing one
  * Use Lambda Machine Local CLI to ssh into container host VM
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
   NAME   DRIVER   STATE   ERRORS
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
   To connect to this virtual machine, run: lambda-machine-local ssh default
   ```

4. List container host virtual machines again to see the newly minted `default` virtual machine.

   ```console
   $ lambda-machine-local ls
   NAME      DRIVER       STATE     ERRORS
   default   virtualbox   Running
   ```

5. Use Lambda Machine Local CLI to ssh into new container host virtual machine.

   As noted in the output of the `lambda-machine-local create` command, you can connect to the virtual machine using `lambda-machine-local ssh` command.

   ```console
   $ lambda-machine-local ssh default

   *Thank you* for using
     _                  _        _         _     _
    | |    __ _  _ __  | |__  __| | __ _  | |   (_) _ _  _  _ __ __
    | |__ / _` || '  \ | '_ \/ _` |/ _` | | |__ | || ' \| || |\ \ /
    |____|\__,_||_|_|_||_.__/\__,_|\__,_| |____||_||_||_|\_,_|/_\_\

   VirtualBox Release: 2018.01.0 | Twitter: @lambda_linux

   ip-192-168-99-101:~$
   ```

<a name="running_containers"></a>
#### Running containers

Run a container with `docker run` to verify your set up.

1. Use `docker run` to download and run `busybox` with a simple `echo` command.

   ```console
   ip-192-168-99-101:~$ docker run busybox echo hello world
   Unable to find image 'busybox:latest' locally
   latest: Pulling from library/busybox
   0ffadd58f2a6: Pull complete
   Digest: sha256:bbc3a03235220b170ba48a157dd097dd1379299370e1ed99ce976df0355d24f0
   Status: Downloaded newer image for busybox:latest
   hello world
   ```

2. Get the container host virtual machine's IP address.

   Any exposed ports is available on container host virtual machine's IP address. You can use `lambda-machine-local ip` command to get the IP address.

   ```console
   $ lambda-machine-local ip default
   192.168.99.101
   ```

3. Run a webserver ([nginx](https://nginx.org/)) in a container with the following command.

   ```console
   ip-192-168-99-101:~$ docker run -d -p 8000:80 nginx
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
```

Following `lambda-machine-local` commands follow this pattern &ndash;

  * `lambda-machine-local inspect`
  * `lambda-machine-local ip`
  * `lambda-machine-local kill`
  * `lambda-machine-local restart`
  * `lambda-machine-local ssh`
  * `lambda-machine-local start`
  * `lambda-machine-local status`
  * `lambda-machine-local stop`

<a name="auto_updates_and_crash_reporting"></a>
#### Auto updates and crash reporting

When you create a new container host virtual machine, we check on GitHub to see
if you are running the latest version of Lambda Linux VirtualBox flavor. If your
local ISO image is outdated, we download the newer version and use that to
create container host.

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

  * [create](#cli_create)
  * [help](#cli_help)
  * [inspect](#cli_inspect)
  * [ip](#cli_ip)
  * [kill](#cli_kill)
  * [ls](#cli_ls)
  * [restart](#cli_restart)
  * [rm](#cli_rm)
  * [scp](#cli_scp)
  * [ssh](#cli_ssh)
  * [start](#cli_start)
  * [status](#cli_status)
  * [stop](#cli_stop)

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
To connect to this virtual machine, run: lambda-machine-local ssh dev
```

**Note:** Libmachine has the concept of _driver_ which is used to manage
hypervisors or cloud instances. In Lambda Machine Local we use VirtualBox as our
default driver. It is also the only driver supported at this time.

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

<a name="cli_help"></a>
#### help

```console
Usage: lambda-machine-local help [arg...]

Shows a list of commands or help for one command
```

This command has the form `lambda-machine-local help <COMMAND>`. For example &ndash;

```console
$ lambda-machine-local help ip
Usage: lambda-machine-local ip [arg...]

Get the IP address of a machine

Description:
   Argument(s) are one or more machine names.
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
NAME   DRIVER       STATE     ERRORS
dev    virtualbox   Running

$ lambda-machine-local kill dev
Killing "dev"...
Machine "dev" was killed.

$ lambda-machine-local ls
NAME   DRIVER       STATE     ERRORS
dev    virtualbox   Stopped
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
NAME      DRIVER       STATE     ERRORS
default   virtualbox   Running
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
NAME      DRIVER       STATE     ERRORS
default   virtualbox   Stopped
foo0      virtualbox   Running
foo1      virtualbox   Running
foo2      virtualbox   Running

$ lambda-machine-local ls -filter name=foo0
NAME   DRIVER       STATE     ERRORS
foo0   virtualbox   Running

$ lambda-machine-local ls --filter driver=virtualbox --filter state=Stopped
NAME      DRIVER       STATE     ERRORS
default   virtualbox   Stopped
```

##### Formatting

The formatting option (`--format`) will pretty-print container host virtual
machines using a Go template.

Valid placeholders for the Go template are listed below.

| Placeholder | Description |
| ----------- | ----------- |
| .Name | Container host virtual machine name |
| .State | Container host virtual machine state (running, stopped, ...) |
| .Error | Container host virtual machine error |
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

*Thank you* for using
  _                  _        _         _     _
 | |    __ _  _ __  | |__  __| | __ _  | |   (_) _ _  _  _ __ __
 | |__ / _` || '  \ | '_ \/ _` |/ _` | | |__ | || ' \| || |\ \ /
 |____|\__,_||_|_|_||_.__/\__,_|\__,_| |____||_||_||_|\_,_|/_\_\

VirtualBox Release: 2018.01.0 | Twitter: @lambda_linux

ip-192-168-99-101:~$
```

You can also specify commands to run on container host virtual machine by
appending them directly to the `lambda-machine-local ssh` command, much like the
regular `ssh` program works.

```console
$ lambda-machine-local ssh dev free
total       used       free     shared    buffers     cached
Mem:       2051564     202412    1849152     107932       1084     143568
-/+ buffers/cache:      57760    1993804
Swap:            0          0          0

$ lambda-machine-local ssh dev sudo df -h
Filesystem                Size      Used Available Use% Mounted on
devtmpfs                 10.0M         0     10.0M   0% /dev
shm                    1001.7M         0   1001.7M   0% /dev/shm
tmpfs                  1001.7M    105.3M    896.5M  11% /

[...]
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
NAME   DRIVER       STATE     ERRORS
dev    virtualbox   Running

$ lambda-machine-local stop dev
Stopping "dev"...
Machine "dev" was stopped.

$ lambda-machine-local ls
NAME   DRIVER       STATE     ERRORS
dev    virtualbox   Stopped
```

<a name="customizing_lambda_linux_virtualbox_flavor"></a>
## Customizing Lambda Linux VirtualBox flavor

Lambda Machine Local is designed to be a stable and secure **local development
primitive**. As such, we cannot anticipate all the use-cases where you might
find using Lambda Machine Local useful. Therefore we would like to provide you
with mechanisms to customize Lambda Linux VirtualBox flavor for your workflow
needs.

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

Upon booting Lambda Linux VirtualBox flavor starts OpenRC service manager.
OpenRC local scripts first checks to see if `disk.vmdk` is already setup. If so,
it proceeds with the rest of the initialization process. If `disk.vmdk` is not
setup then `/etc/local.d/01-disk-setup.start` partitions and formats the disk as
follows.

`01-disk-setup.start` creates
a
[GUID partition table (GPT)](https://en.wikipedia.org/wiki/GUID_Partition_Table)
with one partitions.

```console
ip-192-168-99-101:~$ sudo sgdisk -p /dev/sda

[...]

Number  Start (sector)    End (sector)  Size       Code  Name
   1            2048        40959966   19.5 GiB    8300  Linux
```

The partition (`/dev/sda1`) is a standard Linux partition formatted with
an ext4 filesystem and mounted at `/var/lib/lambda-machine-local`.

All changes made via `docker` commands and changes to
`/var/lib/lambda-machine-local` are preserved across virtual machine reboots.
For your convenience, `/var/lib/lambda-machine-local` is organized as follows.

<a name="customizing_lambda_linux_virtualbox_flavor_home_ll_user"></a>
```console
/var/lib/lambda-machine-local/
|-- bootlocal.sh
|-- bootsync.sh
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
ip-192-168-99-101:~$ id ll-user
uid=500(ll-user) gid=500(ll-user) groups=500(ll-user),10(wheel),103(docker)

ip-192-168-99-101:~$ ls -la /var/lib/lambda-machine-local/home | grep ll-user
drwx------    3 ll-user  ll-user       4096 Dec 27 12:30 ll-user

ip-192-168-99-101:~$ ls -la / | grep Users
drwxr-xr-x    1 ll-user  ll-user        170 Feb 19  2016 Users
```

OpenRC local scripts creates two customization files that it uses during
subsequent invocations. These are -

  * `/var/lib/lambda-machine-local/bootlocal.sh`

  * `/var/lib/lambda-machine-local/bootsync.sh`

#### `bootlocal.sh`

`/etc/local.d/07-bootlocal.start` executes `bootlocal.sh` at the end of its
initialization process. For most users we recommend using this file for your
customization.

#### `bootsync.sh`

`bootsync.sh` is the first customizable hook executed by
`/etc/local.d/05-bootsync.start`. `bootsync.sh` is executed prior to starting
the docker and ssh daemons.

**Note:** OpenRC local scripts `bootsync.sh` to exit with posix return status of
zero. If a non-zero return status is encountered, then initialization will
**not** proceed. Therefore when customizing this file please ensure that any
errors in your program are gracefully handled.

#### Developing and debugging customization scripts

Developing and debugging customization scripts can sometimes get tricky. We
recommend that you use `--virtualbox-ui-type gui` flag when developing or
debugging your customization script.

With this flag, the virtual machine starts with a GUI console. You will be able
to log into GUI console using the user name `ll-user` and password `ll-user`.
You can then `sudo -i` and investigate the behavior of your customization
scripts.

<a name="virtualbox_shared_folder_support"></a>
## VirtualBox shared folder support

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
sparse disk image. By default the size of `disk.vmdk` sparse disk image is 20GB.
If your workflow requires more than 20GB, you can use `--virtualbox-disk-size`
flag in [`lambda-machine-local create`](#cli_create) command to increase the
size of `disk.vmdk`.

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

##### 2. Set up synchronization mechanism

In Lambda Machine Local, we have support using _Rsync over SSH_ synchronization
mechanism. It works uniformly across Mac, Windows and Linux.

`rsync` is available natively on Linux and Mac. On Windows we recommend
installing `rsync` using [MSYS2](http://www.msys2.org/). MSYS2 comes with a
package management tool called `pacman`. You can install `rsync` and `openssh`
using the following command.

```console
$ pacman -S rsync openssh
```

Lambda Linux VirtualBox flavor also comes pre-installed with `rsync`. There is
no additional package installation required on container host virtual machine.

```console
ip-192-168-99-101:~$ which rsync
/usr/bin/rsync
```

`lambda-machine-local scp` command has `--delta` (`-d`) flag which uses `rsync`
under the hood. To synchronize files (in the example below `app` directory)
between your VirtualBox host and guest virtual machine you can &ndash;

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
more details on how to customize `rsync`. If you would like to specify your own
`ssh` command then you can find the RSA private key for container host virtual
machine at `~/.docker/lambda-machine-local/machines/virtualmachinename/id_rsa`.

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

There are a number of tools that let you monitor filesystem events and trigger
actions. Some popular ones are &ndash;

  * [`watchdog`](https://github.com/gorakhargosh/watchdog) with `watchmedo` command
  * [`guard`](https://github.com/guard/guard)
  * [`fswatch`](https://github.com/emcrisostomo/fswatch)

You can use any of the above tools or something that you are already comfortable
with to build your workflow.

<a name="additional_information"></a>
## Additional Information

<a name="windows_installation"></a>
### Windows Installation

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
NAME         DRIVER       STATE     ERRORS
ll-default   virtualbox   Running

username@DESKTOP-xxxxxxx MINGW64 ~
$ lambda-machine-local rm -y ll-default
About to remove ll-default
WARNING: This action will delete both local reference and remote instance.
Can't remove "ll-default"

username@DESKTOP-xxxxxxx MINGW64 ~
$ lambda-machine-local ls
NAME         DRIVER      STATE   ERRORS
ll-default not found Error open C:\Users\username\.docker\lambda-machine-local\
machines\ll-default\config.json: The system cannot find the file specified.
```

If you encounter this issue, you can use [`lambda-machine-local rm`](#cli_rm)
with `--force` (`-f`) flag to remove the `ll-default` directory and correct the
error.

```console
username@DESKTOP-xxxxxxx MINGW64 ~
$ lambda-machine-local rm -f ll-default
About to remove ll-default
WARNING: This action will delete both local reference and remote instance.
Error removing host "ll-default": open C:\Users\username\.docker\lambda-machine-local\machines\ll-default\config.json:
The system cannot find the file specified.
Successfully removed ll-default

username@DESKTOP-xxxxxxx MINGW64 ~
$ lambda-machine-local ls
NAME   DRIVER   STATE   ERRORS

username@DESKTOP-xxxxxxx MINGW64 ~
$
```

<a name="running_lambda_machine_local_alongside_docker_machine"></a>
### Running Lambda Machine Local alongside Docker Machine

Lambda Machine Local and Docker Machine both use VirtualBox to manage their
respective container host virtual machines. In case of Docker Machine, it uses
boot2docker ISO image. In Lambda Machine Local we use Lambda Linux VirtualBox
Flavor ISO image.

When running Lambda Machine Local and Docker Machine at the same time, it is
**important** not to use the same container host virtual machine name. We
recommend that you use a prefix to distinguish between your Docker Machine and
Lambda Machine Local container virtual machines. In the following example we are
using the prefixes `b2d-` and `ll-`.

```console
$ docker-machine create b2d-default

$ lambda-machine-local create ll-default
```

If by mistake you create container host virtual machines with the same name,
you will encounter the following error.

```console
VBoxManage: error: The machine 'default' is already locked for a session (or being unlocked)
VBoxManage: error: Details: code VBOX_E_INVALID_OBJECT_STATE (0x80bb0007), component MachineWrap, interface
IMachine, callee nsISupports
VBoxManage: error: Context: "LockMachine(a->session, LockType_Write)" at line 493 of file VBoxManageModifyVM.cpp
```

You can use `VBoxManage` command to delete the virtual machines and remove the
metadata.

```console
$ VBoxManage list vms
"default" {136b3d00-d9bb-45ba-a927-9bca3087638e}
"default" {6b80915a-1c66-47c1-bc8c-fb6fdec42560}

$ VBoxManage unregistervm --delete 136b3d00-d9bb-45ba-a927-9bca3087638e
0%...10%...20%...30%...40%...50%...60%...70%...80%...90%...100%

$ VBoxManage unregistervm --delete 6b80915a-1c66-47c1-bc8c-fb6fdec42560
0%...10%...20%...30%...40%...50%...60%...70%...80%...90%...100%

$ rm -rf ~/.docker/machine/machines/default

$ rm -rf ~/.docker/lambda-machine-local/machines/default
```
