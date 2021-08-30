# Azure template generator for ACC VMs

### Supported VMs
- Confidential VM (CVM)

#### Pre-Requisite
If you want to build on Linux or MacOS, make sure you have golang installed. If not, you can download it by following these instructions:
```sh
# update your system
sudo apt update
sudo apt upgrade

# download the go binary
wget https://dl.google.com/go/go1.15.5.linux-amd64.tar.gz

# extract binaries tarball
sudo tar -C /usr/local/ -xzf go1.15.5.linux-amd64.tar.gz

# set the right path
cd /usr/local/
echo $PATH
sudo nano $HOME/.profile
```
inside your profile, append the following:
```sh
export PATH=$PATH:/usr/local/go/bin
```
save and apply changes:
```sh
source $HOME/.profile
```
check that it installed correctly:
```sh
go version
```

#### Build
On Linux or MacOS, execute:

```sh
make build
```

Alternatively, you can build with a Docker container:
```sh
docker run \
  -v $PWD:/go/src/github.com/microsoft/acc-vm-engine \
  -w /go/src/github.com/microsoft/acc-vm-engine \
  golang:1.15-alpine go build ./cmd/acc-vm-engine
```

#### Generate template
A sample configuration file for CVM deployment is located in the [test/](./test/) folder.
If your machine runs Linux or MacOS, you can execute the following command to generate templates for a windows vm deployment: 
```sh
./acc-vm-engine generate -c ./test/cvm-win.json
```

Please note that for a linux deployment, you will have to put your ssh public key into the [test/cvm-ubuntu.json](./test/cvm-ubuntu.json) parameter file under 'admin_password_or_key' or you will not be able to connect. Then execute: 
```sh
./acc-vm-engine generate -c ./test/cvm-ubuntu.json
```

In case you need to generate an ssh public key, you can use the following command and use the public key content that it outputs into your ~/.ssh folder: 
```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
```

Alternatively, use Docker container:
```sh
docker run \
  -v $PWD:/go/src/github.com/microsoft/acc-vm-engine \
  -w /go/src/github.com/microsoft/acc-vm-engine \
  golang:1.15-alpine ./acc-vm-engine generate -c ./test/cvm-win.json
```
The template and parameter files will be generated in `_output` directory (by default).

#### Deploy the VM
Use [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) or PowerShell to deploy the VM.

When using Azure CLI, you may want to log in to Azure and set default subscription. This is a one-time operation:
```sh
SUB=<subscription ID>

az login

az account set --subscription ${SUB}
```

Create a resource group and deploy the VM:
```sh
RGROUP=<resource group name>
LOC=<deployment region>

az group create -n ${RGROUP} -l ${LOC}

az deployment group create \
  --resource-group ${RGROUP} \
  --name MyDeployment \
  --template-file ./_output/azuredeploy.json \
  --parameters @./_output/azuredeploy.parameters.json
```
