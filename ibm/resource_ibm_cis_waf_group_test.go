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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMCisWAFGroup_Basic(t *testing.T) {
	name := "ibm_cis_waf_group." + "test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCis(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCisWAFGroupConfigBasic1("test", cisDomainStatic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "mode", "on"),
				),
			},
			{
				Config: testAccCheckCisWAFGroupConfigBasic2("test", cisDomainStatic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "mode", "off"),
				),
			},
			{
				Config: testAccCheckCisWAFGroupConfigBasic1("test", cisDomainStatic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "mode", "on"),
				),
			},
		},
	})
}

func TestAccIBMCisWAFGroup_Import(t *testing.T) {
	name := "ibm_cis_waf_group." + "test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCisWAFGroupConfigBasic2("test", cisDomainStatic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "mode", "off"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCisWAFGroupConfigBasic1(id string, cisDomainStatic string) string {
	return testAccCheckIBMCisDomainDataSourceConfigBasic1() + fmt.Sprintf(`
	resource "ibm_cis_waf_group" "%[1]s" {
		cis_id     = data.ibm_cis.cis.id
		domain_id  = data.ibm_cis_domain.cis_domain.domain_id
		package_id = "c504870194831cd12c3fc0284f294abb"
		group_id   = "3d8fb0c18b5a6ba7682c80e94c7937b2"
		mode       = "on"
	  }
`, id)
}
func testAccCheckCisWAFGroupConfigBasic2(id string, cisDomainStatic string) string {
	return testAccCheckIBMCisDomainDataSourceConfigBasic1() + fmt.Sprintf(`
	resource "ibm_cis_waf_group" "%[1]s" {
		cis_id     = data.ibm_cis.cis.id
		domain_id  = data.ibm_cis_domain.cis_domain.domain_id
		package_id = "c504870194831cd12c3fc0284f294abb"
		group_id   = "3d8fb0c18b5a6ba7682c80e94c7937b2"
		mode       = "off"
	  }
`, id)
}
