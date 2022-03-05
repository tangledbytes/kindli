# Kindli

Kindli stands for KinD in Lima. Kindli is a CLI written ONLY for MacOS to help setup a VM and run upto 100 kind clusters in them with e2e networking setup. 

Kindli can setup the following:
1. Create a VM with bidirectional networking enabled. VM supports running both Intel on ARM and ARM on Intel.
2. Create upto 100 KinD clusters. For KinD, users can bring in custom KinD configs as well.
3. Setup MetalLB to get LoadBalancer working within the cluster.
4. Run `amd64` images on `arm` clusters and `arm` images on `amd64` clusters.

## Install

```
brew tap utkarsh-pro/utkarsh-pro
brew install kindli
```
## Getting Started

```
$ kindli init
```

Running the above will:
1. Install the prerequisites for running kindli
2. Setup the VM
3. Setup networking
4. Create a default KinD cluster
5. Setup docker to use docker daemon running in the VM

Kindli is written only for MacOS and hence the prequisite installation uses `brew` at places. If automatic prerequisite installation is not desirable then pass the flag `--skip-preq-install` to the above the command. To see the prerequisites, run `kindli preq check`.

**Kindli requires root access during the init process to setup the routing table.** If this is undesirable, then instead of running `kindli init`. Setup the enviroment manually:
```bash
$ kindli vm start --cpu 4 --mem "16GiB" --mount $HOME:ro # Start the VM
...
$ kindli vm restart # Restart the VM - Required to setup docker sockets
...
$ kindli create # Create a default cluster
...
$ kindli network setup # Setup e2e networking - Optional (Makes LoadBalancer IP would be reachable from Host)
...
```

## Command Reference

### Init

Init command combines multiple kindli commands to get the user going and hence should be the first command that user should run.

```
$ kindli init -h
init will initialize kindli

Initialization process will perform the following operations:
1. Install all the prerequisites (can be skipped, see the flags)
2. Start kindli default VM
3. Create a default kind cluster in the VM
4. Setup e2e networking to the KinD network

Usage:
  kindli init [flags]

Flags:
      --arch string         VM architecture
      --cpu int             specify number of cpu assigned to VM (default 4)
      --disk string         specify disk space assigned to the VM (default "100GiB")
  -h, --help                help for init
      --mem string          specify memory to be assigned to VM (default "8GiB")
      --mount strings       specify mounts in form of <PATH>:rw to make the mount available for read/write or in form of <PATH>:ro ti make the mount available only for reading
      --skip-preq-install   if set to true, prerequisite install will be skipped
```

### Create

Create command creates a new KinD cluster. It is similar to running `kind create cluster` but does some additional things like setting up networking.

```
$ kindli create -h
Create a new kind cluster

Usage:
  kindli create [flags]

Flags:
  -c, --config string   kind configuration
  -h, --help            help for create
      --name string     kind cluster name
```

### Delete

Delete command deletes the given KinD cluster. If no name is provided then the default cluster is deleted.

```
$ kindli delete -h
Delete given kind cluster

Usage:
  kindli delete <cluster-name> [flags]

Flags:
  -h, --help   help for delete
```

### Prune

Prune command will clearup all of the Kindli created entities on the user's system (except the prerequisites installed).

```
$ kindli prune
prune will prune kindli

Prune process will perform the following operations:
1. Stop the running VM
2. Delete the VM
3. Cleanup the ~/.kindli directory

Usage:
  kindli prune [flags]

Flags:
  -h, --help   help for prune
```

### VM Start

`kindli vm start` starts the kindli VM with the given configuration. It should be noted that the VM can be started only once.

```
$ kindli vm start -h
Start a new VM for Kindli.

NOTE: VM will be created using lima

Usage:
  kindli vm start [flags]

Flags:
      --arch string     VM architecture
      --cpu int         specify number of cpu assigned to VM (default 4)
      --disk string     specify disk space assigned to the VM (default "100GiB")
  -h, --help            help for start
      --mem string      specify memory to be assigned to VM (default "8GiB")
      --mount strings   specify mounts in form of <PATH>:rw to make the mount available for read/write or in form of <PATH>:ro ti make the mount available only for reading
```

### VM Stop

`kindli vm stop` will stop the running VM. It should be noted that the VM should be in running state or else the command will throw error.

```
$ kindli vm stop -h
Stop running Kindli VM

Usage:
  kindli vm stop [flags]

Flags:
  -h, --help   help for stop
```

### VM Restart

`kindli vm restart` will restart the running VM. It expects the VM to be running state.

```
$ kindli vm restart -h
Restart Kindli VM

Usage:
  kindli vm restart [flags]

Flags:
  -h, --help   help for restart
```

### VM Status

`kindli vm status` shows the status of the running VM.

```
$ kindli vm status -h
Status of Kindli VM

Usage:
  kindli vm status [flags]

Flags:
  -h, --help   help for status
```

### VM Shell

`kindli vm shell` starts a shell session with the running VM. It is like SSHing into the VM.

```
$ kindli vm shell -h
SSH into the Kindli VM

Usage:
  kindli vm shell [flags]

Flags:
  -h, --help   help for shell
```

### Preq Check

`kindli preq check` will check if the prerequisites for kindli are satisfied or not.

```
$ kindli preq check -h
Checks if prerequisites are avaialable on the system or not - Doesn't install missing prequisites though

Usage:
  kindli preq check [flags]

Aliases:
  check, i

Flags:
  -h, --help   help for check
```

### Preq Install

`kindli preq install` will install the required prequisites for kindli. (Uses Brew)

```
$ kindli preq install -h
Install missing prequisites

Usage:
  kindli preq install [flags]

Aliases:
  install, i

Flags:
  -h, --help   help for instal
```

### Network Setup

`kindli network setup` setups the networking between Host Machine, VM and KinD Docker Network. Requires root access.

```
$ kindli network setup -h
setup e2e networking with cluster

Usage:
  kindli network [command]

Available Commands:
  setup       setup e2e networking with cluster

Flags:
  -h, --help   help for network

Use "kindli network [command] --help" for more information about a command.
```

### Image Load

`kindli image load` loads docker image from host machine into the given KinD cluster. If not cluster name is given then default cluster is selected,

```
$ kindli image load -h
Load OCI images in the VM

Usage:
  kindli image load [flags]

Examples:
kindli image <image-name> <cluster-name>

Flags:
  -h, --help   help for load
```