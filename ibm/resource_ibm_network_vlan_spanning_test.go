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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMNetworkVlanSpan_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMNetworkVlanSpanOnConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_network_vlan_spanning.test_vlan", "vlan_spanning", "on"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMNetworkVlanSpanOffConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_network_vlan_spanning.test_vlan", "vlan_spanning", "off"),
				),
			},
		},
	})
}

const testAccCheckIBMNetworkVlanSpanOnConfig_basic = `
resource "ibm_network_vlan_spanning" "test_vlan" {
   "vlan_spanning" = "on"
}`
const testAccCheckIBMNetworkVlanSpanOffConfig_basic = `
resource "ibm_network_vlan_spanning" "test_vlan" {
   "vlan_spanning" = "off"
}`
