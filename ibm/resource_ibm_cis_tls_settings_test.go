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

func TestAccIBMCisTLSSettings_Basic(t *testing.T) {
	name := "ibm_cis_tls_settings." + "test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCis(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCisTLSSettingsConfigBasic1("test", cisDomainStatic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "tls_1_3", "off"),
					resource.TestCheckResourceAttr(name, "universal_ssl", "true"),
					resource.TestCheckResourceAttr(name, "min_tls_version", "1.1"),
				),
			},
			{
				Config: testAccCheckCisTLSSettingsConfigBasic2("test", cisDomainStatic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "universal_ssl", "false"),
					resource.TestCheckResourceAttr(name, "min_tls_version", "1.2"),
				),
			},
			{
				Config: testAccCheckCisTLSSettingsConfigBasic3("test", cisDomainStatic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "tls_1_3", "off"),
					resource.TestCheckResourceAttr(name, "universal_ssl", "false"),
					resource.TestCheckResourceAttr(name, "min_tls_version", "1.1"),
				),
			},
		},
	})
}

func TestAccIBMCisTLSSettings_Import(t *testing.T) {
	name := "ibm_cis_tls_settings." + "test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCisTLSSettingsConfigBasic3("test", cisDomainStatic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "tls_1_3", "off"),
					resource.TestCheckResourceAttr(name, "universal_ssl", "false"),
					resource.TestCheckResourceAttr(name, "min_tls_version", "1.1"),
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

func testAccCheckCisTLSSettingsConfigBasic1(id string, cisDomainStatic string) string {
	return testAccCheckIBMCisDomainDataSourceConfigBasic1() + fmt.Sprintf(`
	resource "ibm_cis_tls_settings" "%[1]s" {
		cis_id          = data.ibm_cis.cis.id
		domain_id       = data.ibm_cis_domain.cis_domain.id
		tls_1_3         = "off"
		min_tls_version = "1.1"
		universal_ssl   = true
	  }
`, id)
}
func testAccCheckCisTLSSettingsConfigBasic2(id string, cisDomainStatic string) string {
	return testAccCheckIBMCisDomainDataSourceConfigBasic1() + fmt.Sprintf(`
	resource "ibm_cis_tls_settings" "%[1]s" {
		cis_id          = data.ibm_cis.cis.id
		domain_id       = data.ibm_cis_domain.cis_domain.id
		tls_1_3         = "on"
		min_tls_version = "1.2"
		universal_ssl   = false
	  }
`, id)
}

func testAccCheckCisTLSSettingsConfigBasic3(id string, cisDomainStatic string) string {
	return testAccCheckIBMCisDomainDataSourceConfigBasic1() + fmt.Sprintf(`
	resource "ibm_cis_tls_settings" "%[1]s" {
		cis_id          = data.ibm_cis.cis.id
		domain_id       = data.ibm_cis_domain.cis_domain.domain_id
		tls_1_3         = "off"
		min_tls_version = "1.1"
		universal_ssl   = false
	  }
`, id)
}
