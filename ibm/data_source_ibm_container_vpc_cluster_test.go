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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMContainerVPCClusterDataSource_basic(t *testing.T) {
	name := fmt.Sprintf("tf-vpc-cluster-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerVPCClusterDataSource(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_container_vpc_cluster.testacc_ds_cluster", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_container_cluster_config.testacc_ds_cluster", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMContainerVPCClusterDataSource(name string) string {
	return testAccCheckIBMContainerVpcClusterBasic(name) + fmt.Sprintf(`
data "ibm_container_vpc_cluster" "testacc_ds_cluster" {
    cluster_name_id = ibm_container_vpc_cluster.cluster.id
}
data "ibm_container_cluster_config" "testacc_ds_cluster" {
	cluster_name_id = ibm_container_vpc_cluster.cluster.id
  }
`)
}
