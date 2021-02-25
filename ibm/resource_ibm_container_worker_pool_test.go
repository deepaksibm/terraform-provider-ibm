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

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIBMContainerWorkerPool_basic(t *testing.T) {

	workerPoolName := fmt.Sprintf("terraform-%d", acctest.RandIntRange(10, 100))
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerWorkerPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerWorkerPoolBasic(clusterName, workerPoolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_worker_pool.test_pool", "worker_pool_name", workerPoolName),
					resource.TestCheckResourceAttr(
						"ibm_container_worker_pool.test_pool", "size_per_zone", "1"),
					resource.TestCheckResourceAttr(
						"ibm_container_worker_pool.test_pool", "labels.%", "2"),
					resource.TestCheckResourceAttr(
						"ibm_container_worker_pool.test_pool", "state", "active"),
					resource.TestCheckResourceAttr(
						"ibm_container_worker_pool.test_pool", "disk_encryption", "true"),
					resource.TestCheckResourceAttr(
						"ibm_container_worker_pool.test_pool", "hardware", "shared"),
				),
			},
			{
				Config: testAccCheckIBMContainerWorkerPoolUpdate(clusterName, workerPoolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_worker_pool.test_pool", "worker_pool_name", workerPoolName),
					resource.TestCheckResourceAttr(
						"ibm_container_worker_pool.test_pool", "size_per_zone", "2"),
					resource.TestCheckResourceAttr(
						"ibm_container_worker_pool.test_pool", "labels.%", "2"),
					resource.TestCheckResourceAttr(
						"ibm_container_worker_pool.test_pool", "state", "active"),
					resource.TestCheckResourceAttr(
						"ibm_container_worker_pool.test_pool", "disk_encryption", "true"),
					resource.TestCheckResourceAttr(
						"ibm_container_worker_pool.test_pool", "hardware", "shared"),
				),
			},
		},
	})
}

func TestAccIBMContainerWorkerPool_importBasic(t *testing.T) {
	workerPoolName := fmt.Sprintf("terraform-%d", acctest.RandIntRange(10, 100))
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerWorkerPoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMContainerWorkerPoolBasic(clusterName, workerPoolName),
			},

			resource.TestStep{
				ResourceName:      "ibm_container_worker_pool.test_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIBMContainerWorkerPool_InvalidSizePerZone(t *testing.T) {
	workerPoolName := fmt.Sprintf("terraform-%d", acctest.RandIntRange(10, 100))
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerWorkerPoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccCheckIBMContainerWorkerPoolInvalidSizePerZone(clusterName, workerPoolName),
				ExpectError: regexp.MustCompile("must be greater than 0"),
			},
		},
	})
}

func testAccCheckIBMContainerWorkerPoolDestroy(s *terraform.State) error {

	csClient, err := testAccProvider.Meta().(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_container_worker_pool" {
			continue
		}

		parts, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}
		cluster := parts[0]
		workerPoolID := parts[1]

		target := v1.ClusterTargetHeader{
			Region: csRegion,
		}

		// Try to find the key
		_, err = csClient.WorkerPools().GetWorkerPool(cluster, workerPoolID, target)

		if err == nil {
			return fmt.Errorf("Worker pool still exists: %s", rs.Primary.ID)
		} else if !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("Error waiting for worker pool (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIBMContainerWorkerPoolBasic(clusterName, workerPoolName string) string {
	return fmt.Sprintf(`

resource "ibm_container_cluster" "testacc_cluster" {
  name            = "%s"
  datacenter      = "%s"
  machine_type    = "%s"
  hardware        = "shared"
  public_vlan_id  = "%s"
  private_vlan_id = "%s"
  kube_version    = "%s"
}

resource "ibm_container_worker_pool" "test_pool" {
  worker_pool_name = "%s"
  machine_type     = "%s"
  cluster          = ibm_container_cluster.testacc_cluster.id
  size_per_zone    = 1
  hardware         = "shared"
  disk_encryption  = true
  labels = {
    "test"  = "test-pool"
    "test1" = "test-pool1"
  }
}`, clusterName, datacenter, machineType, publicVlanID, privateVlanID, kubeVersion, workerPoolName, machineType)
}

func testAccCheckIBMContainerWorkerPoolUpdate(clusterName, workerPoolName string) string {
	return fmt.Sprintf(`

resource "ibm_container_cluster" "testacc_cluster" {
  name            = "%s"
  datacenter      = "%s"
  machine_type    = "%s"
  hardware        = "shared"
  public_vlan_id  = "%s"
  private_vlan_id = "%s"
  kube_version    = "%s"
}

resource "ibm_container_worker_pool" "test_pool" {
  worker_pool_name = "%s"
  machine_type     = "%s"
  cluster          = ibm_container_cluster.testacc_cluster.id
  size_per_zone    = 2
  hardware         = "shared"
  disk_encryption  = true
  labels = {
    "test"  = "test-pool"
    "test1" = "test-pool1"
  }
}`, clusterName, datacenter, machineType, publicVlanID, privateVlanID, kubeVersion, workerPoolName, machineType)
}

func testAccCheckIBMContainerWorkerPoolInvalidSizePerZone(clusterName, workerPoolName string) string {
	return fmt.Sprintf(`
resource "ibm_container_worker_pool" "test_pool" {
  worker_pool_name = "%s"
  machine_type     = "%s"
  cluster          = "%s"
  size_per_zone    = 0
  hardware         = "shared"
  disk_encryption  = true

  labels = {
    "test"  = "test-pool"
    "test1" = "test-pool1"
  }
}`, workerPoolName, machineType, clusterName)
}
