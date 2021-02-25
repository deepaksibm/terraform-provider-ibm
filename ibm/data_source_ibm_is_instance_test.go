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

func TestAccIBMISInstanceDataSource_basic(t *testing.T) {

	vpcname := fmt.Sprintf("tfins-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tfins-subnet-%d", acctest.RandIntRange(10, 100))
	sshname := fmt.Sprintf("tfins-ssh-%d", acctest.RandIntRange(10, 100))
	instanceName := fmt.Sprintf("tfins-name-%d", acctest.RandIntRange(10, 100))
	resName := "data.ibm_is_instance.ds_instance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMISInstanceDataSourceConfig(vpcname, subnetname, sshname, instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resName, "name", instanceName),
					resource.TestCheckResourceAttr(
						resName, "tags.#", "1"),
				),
			},
		},
	})
}

func testAccCheckIBMISInstanceDataSourceConfig(vpcname, subnetname, sshname, instanceName string) string {
	return fmt.Sprintf(`
resource "ibm_is_vpc" "testacc_vpc" {
  name = "%s"
}

resource "ibm_is_subnet" "testacc_subnet" {
  name            = "%s"
  vpc             = ibm_is_vpc.testacc_vpc.id
  zone            = "%s"
  ipv4_cidr_block = "%s"
}

resource "ibm_is_ssh_key" "testacc_sshkey" {
  name       = "%s"
  public_key = file("test-fixtures/.ssh/id_rsa.pub")
}

resource "ibm_is_instance" "testacc_instance" {
  name    = "%s"
  image   = "%s"
  profile = "%s"
  primary_network_interface {
    port_speed = "100"
    subnet     = ibm_is_subnet.testacc_subnet.id
  }
  vpc  = ibm_is_vpc.testacc_vpc.id
  zone = "%s"
  keys = [ibm_is_ssh_key.testacc_sshkey.id]
  network_interfaces {
    subnet = ibm_is_subnet.testacc_subnet.id
    name   = "eth1"
  }
  tags = ["tag1"]
}
data "ibm_is_instance" "ds_instance" {
  name        = ibm_is_instance.testacc_instance.name
  private_key = file("test-fixtures/.ssh/id_rsa")
  passphrase  = ""
}`, vpcname, subnetname, ISZoneName, ISCIDR, sshname, instanceName, isWinImage, instanceProfileName, ISZoneName)
}
