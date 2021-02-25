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

func TestAccIBMCisCacheSettings_Basic(t *testing.T) {
	name := "ibm_cis_cache_settings." + "test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCis(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCisCacheSettingsConfigBasic1("test", cisDomainStatic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "caching_level", "simplified"),
					resource.TestCheckResourceAttr(name, "browser_expiration", "7200"),
					resource.TestCheckResourceAttr(name, "development_mode", "on"),
					resource.TestCheckResourceAttr(name, "query_string_sort", "on"),
				),
			},
			{
				Config: testAccCheckCisCacheSettingsConfigBasic2("test", cisDomainStatic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "caching_level", "aggressive"),
					resource.TestCheckResourceAttr(name, "browser_expiration", "14400"),
					resource.TestCheckResourceAttr(name, "development_mode", "off"),
					resource.TestCheckResourceAttr(name, "query_string_sort", "off"),
				),
			},
		},
	})
}

func TestAccIBMCisCacheSettings_WithoutPurgeAction(t *testing.T) {
	name := "ibm_cis_cache_settings." + "test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCis(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCisCacheSettingsConfigWithoutPurgeAction("test", cisDomainStatic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "caching_level", "simplified"),
					resource.TestCheckResourceAttr(name, "browser_expiration", "7200"),
					resource.TestCheckResourceAttr(name, "development_mode", "on"),
					resource.TestCheckResourceAttr(name, "query_string_sort", "on"),
				),
			},
			{
				Config: testAccCheckCisCacheSettingsConfigBasic2("test", cisDomainStatic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "caching_level", "aggressive"),
					resource.TestCheckResourceAttr(name, "browser_expiration", "14400"),
					resource.TestCheckResourceAttr(name, "development_mode", "off"),
					resource.TestCheckResourceAttr(name, "query_string_sort", "off"),
				),
			},
		},
	})
}

func TestAccIBMCisCacheSettings_Import(t *testing.T) {
	name := "ibm_cis_cache_settings." + "test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCisCacheSettingsConfigBasic2("test", cisDomainStatic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "caching_level", "aggressive"),
					resource.TestCheckResourceAttr(name, "browser_expiration", "14400"),
					resource.TestCheckResourceAttr(name, "development_mode", "off"),
					resource.TestCheckResourceAttr(name, "query_string_sort", "off"),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"purge"},
			},
		},
	})
}

func testAccCheckCisCacheSettingsConfigBasic1(id string, cisDomainStatic string) string {
	return testAccCheckIBMCisDomainDataSourceConfigBasic1() + fmt.Sprintf(`
	resource "ibm_cis_cache_settings" "%[1]s" {
		cis_id          = data.ibm_cis.cis.id
		domain_id       = data.ibm_cis_domain.cis_domain.domain_id
		caching_level      = "simplified"
		browser_expiration = 7200
		development_mode   = "on"
		query_string_sort  = "on"
		purge_all          = true
	  }
`, id)
}
func testAccCheckCisCacheSettingsConfigBasic2(id string, cisDomainStatic string) string {
	return testAccCheckIBMCisDomainDataSourceConfigBasic1() + fmt.Sprintf(`
	resource "ibm_cis_cache_settings" "%[1]s" {
		cis_id          = data.ibm_cis.cis.id
		domain_id       = data.ibm_cis_domain.cis_domain.domain_id
		caching_level      = "aggressive"
		browser_expiration = 14400
		development_mode   = "off"
		query_string_sort  = "off"
		purge_by_urls      = ["http://test.com/index.html", "http://example.com/index.html"]
	  }
`, id)
}

func testAccCheckCisCacheSettingsConfigWithoutPurgeAction(id string, cisDomainStatic string) string {
	return testAccCheckIBMCisDomainDataSourceConfigBasic1() + fmt.Sprintf(`
	resource "ibm_cis_cache_settings" "%[1]s" {
		cis_id          = data.ibm_cis.cis.id
		domain_id       = data.ibm_cis_domain.cis_domain.domain_id
		caching_level      = "simplified"
		browser_expiration = 7200
		development_mode   = "on"
		query_string_sort  = "on"
	  }
`, id)
}
