/* IBM Confidential
*  Object Code Only Source Materials
*  5747-SM3
*  (c) Copyright IBM Corp. 2017,2021
*
*  The source code for this program is not published or otherwise divested
*  of its trade secrets, irrespective of what has been deposited with the
*  U.S. Copyright Office. */

package ibm

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
)

func TestAccIBMContainerCluster_basic(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerClusterBasic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "worker_num", "1"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "default_pool_size", "1"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "kube_version", kubeVersion),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "hardware", "shared"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "worker_pools.#", "1"),
					resource.TestCheckResourceAttrSet(
						"ibm_container_cluster.testacc_cluster", "resource_group_id"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "labels.%", "2"),
				),
			},
			{
				Config: testAccCheckIBMContainerClusterUpdate(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "worker_num", "2"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "default_pool_size", "2"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "kube_version", kubeUpdateVersion),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "workers_info.0.version", kubeUpdateVersion),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "is_trusted", "false"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "hardware", "shared"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "worker_pools.#", "1"),
					resource.TestCheckResourceAttrSet(
						"ibm_container_cluster.testacc_cluster", "resource_group_id"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "labels.%", "3"),
				),
			},
		},
	})
}

func TestAccIBMContainerClusterWaitTill(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerClusterWaitTill(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "default_pool_size", "1"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "kube_version", kubeVersion),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "hardware", "shared"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "worker_pools.#", "1"),
					resource.TestCheckResourceAttrSet(
						"ibm_container_cluster.testacc_cluster", "resource_group_id"),
				),
			},
		},
	})
}
func TestAccIBMContainerCluster_trusted(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerClusterTrusted(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "default_pool_size", "1"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "kube_version", kubeVersion),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "worker_pools.#", "1"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "is_trusted", "true"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "hardware", "dedicated"),
				),
			},
		},
	})
}

func TestAccIBMContainerCluster_KmsEnable(t *testing.T) {
	clusterName := fmt.Sprintf("terraform1_%d", acctest.RandIntRange(10, 100))
	kmsInstanceName := fmt.Sprintf("kmsInstance_%d", acctest.RandIntRange(10, 100))
	rootKeyName := fmt.Sprintf("rootKey_%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerClusterKmsEnable(clusterName, kmsInstanceName, rootKeyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "kms_config.#", "1"),
				),
			},
		},
	})
}
func TestAccIBMContainerCluster_nosubnet_false(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerClusterNosubnetFalse(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "default_pool_size", "1"),
					resource.TestCheckResourceAttrSet(
						"ibm_container_cluster.testacc_cluster", "ingress_hostname"),
					resource.TestCheckResourceAttrSet(
						"ibm_container_cluster.testacc_cluster", "ingress_secret"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "hardware", "dedicated"),
				),
			},
		},
	})
}

func TestAccIBMContainerCluster_worker_count(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerCluster_worker_count(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "default_pool_size", "1"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "worker_num", "1"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "workers_info.#", "2"),
				),
			},
			{
				Config: testAccCheckIBMContainerClusterWorkerCountUpdate(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "default_pool_size", "2"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "worker_num", "2"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "workers_info.#", "4"),
				),
			},
		},
	})
}

func TestAccIBMContainerCluster_with_worker_num_zero(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccCheckIBMContainerClusterWithWorkerNumZero(clusterName),
				ExpectError: regexp.MustCompile("must be greater than 0"),
			},
		},
	})
}

func TestAccIBMContainerCluster_diskEnc(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerClusterDiskEnc(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "default_pool_size", "1"),
				),
			},
		},
	})
}

//testAccCheckIBMContainerClusterOptionalOrgSpace_basic
func TestAccIBMContainerClusterOptionalOrgSpace_basic(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerClusterOptionalOrgSpaceBasic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "default_pool_size", "1"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "albs.#", "2"),
				),
			},
		},
	})
}

func TestAccIBMContainerCluster_private_subnet(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerClusterPrivateSubnet(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "default_pool_size", "1"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "ingress_hostname", ""),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "ingress_secret", ""),
				),
			},
		},
	})
}

func TestAccIBMContainerCluster_private_and_public_subnet(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerClusterPrivateAndPublicSubnet(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "default_pool_size", "1"),
					resource.TestCheckResourceAttrSet(
						"ibm_container_cluster.testacc_cluster", "ingress_hostname"),
					resource.TestCheckResourceAttrSet(
						"ibm_container_cluster.testacc_cluster", "ingress_secret"),
				),
			},
		},
	})
}

func TestAccIBMContainerCluster_Tag(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerClusterTag(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "default_pool_size", "1"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "tags.#", "1"),
				),
			},
			{
				Config: testAccCheckIBMContainerClusterUpdateTag(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "default_pool_size", "1"),
					resource.TestCheckResourceAttr(
						"ibm_container_cluster.testacc_cluster", "tags.#", "2"),
				),
			},
		},
	})
}

func testAccCheckIBMContainerClusterDestroy(s *terraform.State) error {
	csClient, err := testAccProvider.Meta().(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_container_cluster" {
			continue
		}
		targetEnv := containerv1.ClusterTargetHeader{}
		// targetEnv := getClusterTargetHeaderTestACC()
		// Try to find the key
		_, err := csClient.Clusters().Find(rs.Primary.ID, targetEnv)

		if err == nil {
			return fmt.Errorf("Cluster still exists: %s", rs.Primary.ID)
		} else if !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("Error waiting for cluster (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

// func getClusterTargetHeaderTestACC() v1.ClusterTargetHeader {
// 	org := cfOrganization
// 	space := cfSpace
// 	c := new(bluemix.Config)
// 	sess, err := session.New(c)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	client, err := mccpv2.New(sess)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	orgAPI := client.Organizations()
// 	myorg, err := orgAPI.FindByName(org, BluemixRegion)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	spaceAPI := client.Spaces()
// 	myspace, err := spaceAPI.FindByNameInOrg(myorg.GUID, space, BluemixRegion)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	accClient, err := accountv2.New(sess)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	accountAPI := accClient.Accounts()
// 	myAccount, err := accountAPI.FindByOrg(myorg.GUID, c.Region)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	target := v1.ClusterTargetHeader{
// 		OrgID:     myorg.GUID,
// 		SpaceID:   myspace.GUID,
// 		AccountID: myAccount.GUID,
// 	}

// 	return target
// }

func testAccCheckIBMContainerClusterBasic(clusterName string) string {
	return fmt.Sprintf(`

data "ibm_resource_group" "testacc_ds_resource_group" {
  is_default = "true"
}

resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"
  resource_group_id = data.ibm_resource_group.testacc_ds_resource_group.id

  default_pool_size = 1

  hardware        = "shared"
  kube_version    = "%s"
  machine_type    = "%s"
  public_vlan_id  = "%s"
  private_vlan_id = "%s"
  no_subnet       = true
  region          = "%s"
}	`, clusterName, datacenter, kubeVersion, machineType, publicVlanID, privateVlanID, csRegion)
}
func testAccCheckIBMContainerClusterWaitTill(clusterName string) string {
	return fmt.Sprintf(`

data "ibm_resource_group" "testacc_ds_resource_group" {
  is_default = "true"
}

resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"
  resource_group_id = data.ibm_resource_group.testacc_ds_resource_group.id

  default_pool_size = 1
  wait_till       = "MasterNodeReady"
  hardware        = "shared"
  kube_version    = "%s"
  machine_type    = "%s"
  public_vlan_id  = "%s"
  private_vlan_id = "%s"
  region          = "%s"
}	`, clusterName, datacenter, kubeVersion, machineType, publicVlanID, privateVlanID, csRegion)
}

func testAccCheckIBMContainerClusterKmsEnable(clusterName, kmsInstanceName, rootKeyName string) string {
	return fmt.Sprintf(`
	
	data "ibm_resource_group" "testacc_ds_resource_group" {
		name = "default"
	}
	
	resource "ibm_resource_instance" "kms_instance1" {
		name              = "%s"
		service           = "kms"
		plan              = "tiered-pricing"
		location          = "us-south"
	}
	  
	resource "ibm_kms_key" "test" {
		instance_id = "${ibm_resource_instance.kms_instance1.guid}"
		key_name = "%s"
		standard_key =  false
		force_delete = true
	}
	
	resource "ibm_container_cluster" "testacc_cluster" {
		name              = "%s"
		datacenter        = "%s"
		no_subnet         = true
		default_pool_size = 2
		hardware          = "shared"
		resource_group_id = data.ibm_resource_group.testacc_ds_resource_group.id
		machine_type      = "%s"
		public_vlan_id    = "%s"
		private_vlan_id   = "%s"
		kms_config {
			instance_id = ibm_resource_instance.kms_instance1.guid
			crk_id = ibm_kms_key.test.key_id
			private_endpoint = false
		}
	}

`, kmsInstanceName, rootKeyName, clusterName, datacenter, machineType, publicVlanID, privateVlanID)
}

func testAccCheckIBMContainerClusterTrusted(clusterName string) string {
	return fmt.Sprintf(`


resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"



  default_pool_size = 1

  kube_version      = "%s"
  machine_type      = "%s"
  hardware          = "dedicated"
  public_vlan_id    = "%s"
  private_vlan_id   = "%s"
  no_subnet         = true
  is_trusted        = true
  wait_time_minutes = 1440
}	`, clusterName, datacenter, kubeVersion, trustedMachineType, publicVlanID, privateVlanID)
}

func testAccCheckIBMContainerClusterNosubnetFalse(clusterName string) string {
	return fmt.Sprintf(`


resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"



  machine_type    = "%s"
  hardware        = "dedicated"
  public_vlan_id  = "%s"
  private_vlan_id = "%s"
  no_subnet       = false
}	`, clusterName, datacenter, machineType, publicVlanID, privateVlanID)
}

func testAccCheckIBMContainerClusterWithWorkerNumZero(clusterName string) string {
	return fmt.Sprintf(`


resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"

  account_guid      = data.ibm_account.acc.id
  default_pool_size = 0
  machine_type      = "%s"
  hardware          = "shared"
  public_vlan_id    = "%s"
  private_vlan_id   = "%s"
  no_subnet         = true
}	`, clusterName, datacenter, machineType, publicVlanID, privateVlanID)
}

func testAccCheckIBMContainerClusterOptionalOrgSpaceBasic(clusterName string) string {
	return fmt.Sprintf(`

resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"

  machine_type    = "%s"
  hardware        = "shared"
  public_vlan_id  = "%s"
  private_vlan_id = "%s"
  disk_encryption = true
}	`, clusterName, datacenter, machineType, publicVlanID, privateVlanID)
}

func testAccCheckIBMContainerClusterDiskEnc(clusterName string) string {
	return fmt.Sprintf(`


resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"



  machine_type    = "%s"
  hardware        = "shared"
  public_vlan_id  = "%s"
  private_vlan_id = "%s"
  no_subnet       = true
  disk_encryption = false
}	`, clusterName, datacenter, machineType, publicVlanID, privateVlanID)
}

func testAccCheckIBMContainerClusterUpdate(clusterName string) string {
	return fmt.Sprintf(`

data "ibm_resource_group" "testacc_ds_resource_group" {
  is_default = "true"
}

resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"
  worker_num = 2

  default_pool_size = 2

  hardware           = "shared"
  resource_group_id  = data.ibm_resource_group.testacc_ds_resource_group.id
  kube_version       = "%s"
  machine_type       = "%s"
  public_vlan_id     = "%s"
  private_vlan_id    = "%s"
  no_subnet          = true
  update_all_workers = true
  region             = "%s"
}	`, clusterName, datacenter, kubeUpdateVersion, machineType, publicVlanID, privateVlanID, csRegion)
}

func testAccCheckIBMContainerClusterPrivateAndPublicSubnet(clusterName string) string {
	return fmt.Sprintf(`


resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"



  machine_type    = "%s"
  hardware        = "shared"
  public_vlan_id  = "%s"
  private_vlan_id = "%s"
  no_subnet       = true
  subnet_id       = ["%s", "%s"]
}	`, clusterName, datacenter, machineType, publicVlanID, privateVlanID, privateSubnetID, publicSubnetID)
}

func testAccCheckIBMContainerClusterPrivateSubnet(clusterName string) string {
	return fmt.Sprintf(`

resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"
  machine_type    = "%s"
  hardware        = "shared"
  public_vlan_id  = "%s"
  private_vlan_id = "%s"
  no_subnet       = true
  subnet_id       = ["%s"]
}	`, clusterName, datacenter, machineType, publicVlanID, privateVlanID, privateSubnetID)
}

func testAccCheckIBMContainerClusterTag(clusterName string) string {
	return fmt.Sprintf(`

resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"
  machine_type    = "%s"
  hardware        = "shared"
  public_vlan_id  = "%s"
  private_vlan_id = "%s"
  tags            = ["test"]
}	`, clusterName, datacenter, machineType, publicVlanID, privateVlanID)
}

func testAccCheckIBMContainerClusterUpdateTag(clusterName string) string {
	return fmt.Sprintf(`

resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"

  machine_type    = "%s"
  hardware        = "shared"
  public_vlan_id  = "%s"
  private_vlan_id = "%s"
  tags            = ["test", "once"]
}	`, clusterName, datacenter, machineType, publicVlanID, privateVlanID)
}

func testAccCheckIBMContainerCluster_worker_count(clusterName string) string {
	return fmt.Sprintf(`


resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"

  account_guid = data.ibm_account.acc.id

  worker_num = 1

  machine_type    = "%s"
  hardware        = "shared"
  public_vlan_id  = "%s"
  private_vlan_id = "%s"
  no_subnet       = true
}	`, clusterName, datacenter, machineType, publicVlanID, privateVlanID)
}

func testAccCheckIBMContainerClusterWorkerCountUpdate(clusterName string) string {
	return fmt.Sprintf(`

resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"

  default_pool_size = 2
  machine_type      = "%s"
  hardware          = "shared"
  public_vlan_id    = "%s"
  private_vlan_id   = "%s"
  no_subnet         = true
}	`, clusterName, datacenter, machineType, publicVlanID, privateVlanID)
}
