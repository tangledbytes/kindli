# Kindli

Kindli stands for KinD in Lima. Kindli is a CLI written ONLY for MacOS to help setup a VM and run upto 100 kind clusters in them with e2e networking setup. 

Kindli can setup the following:
1. Create upto 200 VMs (if you can 🤷‍♂️) with bidirectional networking enabled. VM supports running both Intel on ARM and ARM on Intel.
2. Create upto 100 KinD clusters per VM. For KinD, users can bring in custom KinD configs as well.
3. Setup MetalLB to get LoadBalancer working within the cluster. E2E networking is **not** supported for IPv6 type clusters.
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
2. Setup default VM with name "kindli"
3. Create a default KinD cluster
4. Setup docker to use docker daemon running in the VM

Kindli is written only for MacOS and hence the prequisite installation uses `brew` at places. If automatic prerequisite installation is not desirable then pass the flag `--skip-preq-install` to the above the command. To see the prerequisites, run `kindli preq check`.

`init` setup is just an elaborate way of the following
```bash
$ kindli vm start --cpu 4 --mem "16GiB" --mount $HOME:ro # Start default VM
...
$ kindli vm restart # Restart the VM - Required to setup docker sockets
...
$ kindli create # Create a default cluster
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
2. Start a VM
3. Create a default kind cluster in the VM

Usage:
  kindli init [flags]

Flags:
      --arch string         VM architecture
      --cpu int             specify number of cpu assigned to VM (default 4)
      --disk string         specify disk space assigned to the VM (default "100GiB")
  -h, --help                help for init
      --mem string          specify memory to be assigned to VM (default "16GiB")
      --mount strings       specify mounts in form of <PATH>:rw to make the mount available for read/write or in form of <PATH>:ro to make the mount available only for reading
      --skip-preq-install   if set to true, prerequisite install will be skipped

Global Flags:
      --cluster-name string   Name of the cluster (default "kindli")
      --vm-name string        Name of the VM (default "kindli")
```

### Create

Create command creates a new KinD cluster. It is similar to running `kind create cluster` but does some additional things like setting up metalLB for loadbalancers (can be skipped via `--skip-metallb`).

```
$ kindli create -h
Create a new kind cluster

Usage:
  kindli create [flags]

Flags:
  -c, --config string   kind configuration
  -h, --help            help for create
  -s, --skip-metallb    skip metallb setup

Global Flags:
      --cluster-name string   Name of the cluster (default "kindli")
      --vm-name string        Name of the VM (default "kindli")
```

### Delete

Delete command deletes the given KinD cluster. If no name is provided then the default cluster is deleted.

```
$ kindli delete -h
Delete given kind cluster

Usage:
  kindli delete [flags]

Flags:
  -h, --help   help for delete

Global Flags:
      --cluster-name string   Name of the cluster (default "kindli")
      --vm-name string        Name of the VM (default "kindli")
```

### Prune

Prune command will clearup all of the Kindli created entities on the user's system (except the prerequisites installed).

```
$ kindli prune
prune will prune kindli

Prune process will perform the following operations:
1. Stop all of the running VMs
2. Delete all of the VMs
3. Cleanup the ~/.kindli directory
4. Optionally clean up the lima cache

Usage:
  kindli prune [flags]

Flags:
  -a, --all          If true, prune will delete all of the VMs (default true)
      --clean-lima   If true, prune will clear lima cache
  -h, --help         help for prune

Global Flags:
      --cluster-name string   Name of the cluster (default "kindli")
      --vm-name string        Name of the VM (default "kindli")
```

### VM Start

`kindli vm start` starts the kindli VM with the given configuration. It should be noted that no more than 200 VMs might be created via Kindli. 

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
      --mem string      specify memory to be assigned to VM (default "16GiB")
      --mount strings   specify mounts in form of <PATH>:rw to make the mount available for read/write or in form of <PATH>:ro to make the mount available only for reading

Global Flags:
      --cluster-name string   Name of the cluster (default "kindli")
      --vm-name string        Name of the VM (default "kindli")
```

### VM Stop

`kindli vm stop` will stop the given (or default) VM. It should be noted that the VM should be in running state or else the command will throw error.

```
Stop running Kindli VM

Usage:
  kindli vm stop [flags]

Flags:
  -h, --help   help for stop

Global Flags:
      --cluster-name string   Name of the cluster (default "kindli")
      --vm-name string        Name of the VM (default "kindli")
```

### VM Restart

`kindli vm restart` will restart the given (or default) VM. It expects the VM to be running state.

```
$ kindli vm restart -h
Restart Kindli VM

Usage:
  kindli vm restart [flags]

Flags:
  -h, --help   help for restart

Global Flags:
      --cluster-name string   Name of the cluster (default "kindli")
      --vm-name string        Name of the VM (default "kindli")
```

### VM Status

`kindli vm status` shows the status of all the VMs.

```
$ kindli vm status -h
Status of Kindli VM

Usage:
  kindli vm status [flags]

Flags:
  -A, --all    Show status of all VMs
  -h, --help   help for status

Global Flags:
      --cluster-name string   Name of the cluster (default "kindli")
      --vm-name string        Name of the VM (default "kindli")
```

### VM Shell

`kindli vm shell` starts a shell session with the given (or default) VM. It is like SSHing into the VM.

```
$ kindli vm shell -h
SSH into the Kindli VM

Usage:
  kindli vm shell [flags]

Flags:
  -h, --help   help for shell

Global Flags:
      --cluster-name string   Name of the cluster (default "kindli")
      --vm-name string        Name of the VM (default "kindli")
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

Global Flags:
      --cluster-name string   Name of the cluster (default "kindli")
      --vm-name string        Name of the VM (default "kindli")
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
  -h, --help   help for install

Global Flags:
      --cluster-name string   Name of the cluster (default "kindli")
      --vm-name string        Name of the VM (default "kindli")
```

### Network Setup

`kindli network setup` setups the networking between Host Machine, given (or default) VM and KinD Docker Network inside the VM. **Requires root access**. It should be noted that E2E networking although can be setup with multiple VMs at once but should ideally be established with one VM at a time only (and this is the only workflow that is tested). This command haults and will keep the connection alive till it is terminated by pressing `CTRL+C` (SIGINT).

```
$ kindli network setup -h
setup e2e networking with cluster

Usage:
  kindli network setup [flags]

Flags:
  -h, --help   help for setup

Global Flags:
      --cluster-name string   Name of the cluster (default "kindli")
      --vm-name string        Name of the VM (default "kindli")
```

### Image Load

`kindli image load` loads docker image from host machine into the given KinD cluster. If cluster name is not given then default cluster is selected. Use `kindli docker-env --vm-name <desired-vm>` to point to the right VM.

```
$ kindli image load -h
Load OCI images in the VM

Usage:
  kindli image load [flags]

Examples:
kindli image <image-name>

Flags:
  -h, --help   help for load

Global Flags:
      --cluster-name string   Name of the cluster (default "kindli")
      --vm-name string        Name of the VM (default "kindli")
```


### Docker Env Setup

`kindli docker-env --vm-name <vm-name>` can be used to point docker client on the host machine to the docker daemon running in the given VM.

```
$ kindli docker-env --vm-name <vm-name>
docker-env sets the docker context of the given VM

Usage:
  kindli docker-env [flags]

Flags:
  -h, --help   help for docker-env

Global Flags:
      --cluster-name string   Name of the cluster (default "kindli")
      --vm-name string        Name of the VM (default "kindli")
```

## FAQ

<details>
  <summary>
    I deleted one of my cluster: How to set kubernetes context to my another cluster?
  </summary>

  Running `kindli create --vm-name <name> --cluster-name <desired-cluster>` again will set the kubeneretes context to the given cluster.
</details>

<details>
  <summary>
    Can I setup E2E networking simultaneously to all of my VMs?
  </summary>

  Kind of. Theoretically, you can setup E2E networking simultaneously to all of your VMs however this flow is **not** tested.
</details>

<details>
  <summary>
    How do I use docker daemon of another VM?
  </summary>

  `kindli` uses docker contexts to switch between daemons. If you want to point docker client to one of the kindli VMs, just run `kindli docker-env --vm-name <vm-name>`
</details>


<details>
  <summary>
    How do I clear everything?
  </summary>

  Just run `kindli prune`. If you want to just prune one VM then provide the name of the VM via `--vm-name`.
</details>