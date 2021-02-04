# Terraform Provider

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.1+
-	[Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/IBM-Cloud/terraform-provider-ibm`

```sh
mkdir -p $GOPATH/src/github.com/IBM-Cloud; cd $GOPATH/src/github.com/IBM-Cloud
git clone git@github.com:IBM-Cloud/terraform-provider-ibm.git
```

Enter the provider directory and build the provider

```sh
cd $GOPATH/src/github.com/IBM-Cloud/terraform-provider-ibm
make build
```

## Docker Image For The Provider

You can also pull the docker image for the ibmcloud terraform provider :

```sh
docker pull ibmterraform/terraform-provider-ibm-docker
```

## Using the provider

If you want to run Terraform with the IBM Cloud provider plugin on your system, complete the following steps:

1. [Download and install Terraform for your system](https://www.terraform.io/intro/getting-started/install.html). 

2. [Download the IBM Cloud provider plugin for Terraform](https://github.com/IBM-Bluemix/terraform-provider-ibm/releases).

3. Unzip the release archive to extract the plugin binary (`terraform-provider-ibm_vX.Y.Z`).

4. Move the binary into the Terraform [plugins directory](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) for the platform.
    - Linux/Unix/OS X: `~/.terraform.d/plugins`
    - Windows: `%APPDATA%\terraform.d\plugins`

5. Export API credential tokens as environment variables. This can either be [IBM Cloud API keys](https://cloud.ibm.com/iam#/users) or Softlayer API keys and usernames, depending on the resources you are provisioning.

```sh
export IC_API_KEY="IBM Cloud API Key"
export IAAS_CLASSIC_API_KEY="IBM Cloud Classic Infrastructure API Key"
export IAAS_CLASSIC_USERNAME="IBM Cloud Classic Infrastructure username associated with Classic Infrastructure API KEY".
```

6. Add the plug-in provider to the Terraform configuration file.

```
provider "ibm" {}
```

See the [official documentation](https://cloud.ibm.com/docs/ibm-cloud-provider-for-terraform?topic=ibm-cloud-provider-for-terraform-getting-started) for more details on using the IBM provider.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
make build
...
$GOPATH/bin/terraform-provider-ibm
...
```

In order to test the provider, you can simply run `make test`.

```sh
make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
make testacc
```
In order to run a particular Acceptance test, export the variable `TESTARGS`. For example

```sh
export TESTARGS="-run TestAccIBMNetworkVlan_Basic"
```
Issuing `make testacc` will now run the testcase with names matching `TestAccIBMNetworkVlan_Basic`. This particular testcase is present in
`ibm/resource_ibm_network_vlan_test.go`

You will also need to export the following environment variables for running the Acceptance tests.
* `IC_API_KEY`- The IBM Cloud API Key
* `IAAS_CLASSIC_API_KEY` - The IBM Cloud Classic Infrastructure API Key
* `IAAS_CLASSIC_USERNAME` - The IBM Cloud Classic Infrastructure username associated with the Classic InfrastAPI Key.

Additional environment variables may be required depending on the tests being run. Check console log for warning messages about required variables. 


# IBM Cloud Ansible Modules

An implementation of generated Ansible modules using the
[IBM Cloud Terraform Provider].

## Prerequisites

1. Install [Python3]

2. [RedHat Ansible] 2.8+

    ```
    pip install "ansible>=2.8.0"
    ```


## Install

1. Download IBM Cloud Ansible modules from [release page]

2. Extract module archive.

    ```
    unzip ibmcloud_ansible_modules.zip
    ```

3. Add modules and module_utils to the [Ansible search path]. E.g.:

    ```
    cp build/modules/* $HOME/.ansible/plugins/modules/.
    cp build/module_utils/* $HOME/.ansible/plugins/module_utils/.

    ```

### Example Projects

1. [VPC Virtual Server Instance](examples/ansible/examples/simple-vm-ssh/)

2. [Power Virtual Server Instance](examples/ansible/examples/simple-vm-power-vs/)


[IBM Cloud Terraform Provider]: https://github.com/IBM-Cloud/terraform-provider-ibm
[Python3]: https://www.python.org/downloads/
[RedHat Ansible]: https://www.ansible.com/
[Ansible search path]: https://docs.ansible.com/ansible/latest/dev_guide/overview_architecture.html#ansible-search-path
[release page]:https://github.com/IBM-Cloud/terraform-provider-ibm/releases

