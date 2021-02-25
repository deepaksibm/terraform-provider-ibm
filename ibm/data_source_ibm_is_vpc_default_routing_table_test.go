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

func TestAccIBMISVPCDefaultRoutingTableDataSource_basic(t *testing.T) {
	node := "data.ibm_is_vpc_default_routing_table.def_route_table"
	vpcname := fmt.Sprintf("tf-vpcname-%d", acctest.RandIntRange(100, 200))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISVPCDefaultRoutingTableDataSourceConfig(vpcname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "id"),
					resource.TestCheckResourceAttrSet(node, "name"),
					resource.TestCheckResourceAttrSet(node, "lifecycle_state"),
				),
			},
		},
	})
}

func testAccCheckIBMISVPCDefaultRoutingTableDataSourceConfig(vpcname string) string {
	return fmt.Sprintf(`

	resource "ibm_is_vpc" "test_vpc" {
  		name = "%s"
	}
	
	data "ibm_is_vpc_default_routing_table" "def_route_table" {
		vpc = ibm_is_vpc.test_vpc.id
	}
	`, vpcname)
}
