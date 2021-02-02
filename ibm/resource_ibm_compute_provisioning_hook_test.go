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
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
)

func TestAccIBMComputeProvisioningHook_Basic(t *testing.T) {
	var hook datatypes.Provisioning_Hook

	hookName1 := fmt.Sprintf("%s%s", "tfuathook", acctest.RandString(10))
	hookName2 := fmt.Sprintf("%s%s", "tfuathook", acctest.RandString(10))
	uri1 := "http://www.weather.com"
	uri2 := "https://www.ibm.com"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMComputeProvisioningHookDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMComputeProvisioningHookConfig(hookName1, uri1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMComputeProvisioningHookExists("ibm_compute_provisioning_hook.test-provisioning-hook", &hook),
					testAccCheckIBMComputeProvisioningHookAttributes(&hook, hookName1, uri1),
					resource.TestCheckResourceAttr(
						"ibm_compute_provisioning_hook.test-provisioning-hook", "name", hookName1),
					resource.TestCheckResourceAttr(
						"ibm_compute_provisioning_hook.test-provisioning-hook", "uri", uri1),
				),
			},

			resource.TestStep{
				Config: testAccCheckIBMComputeProvisioningHookConfig(hookName2, uri2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMComputeProvisioningHookExists("ibm_compute_provisioning_hook.test-provisioning-hook", &hook),
					resource.TestCheckResourceAttr(
						"ibm_compute_provisioning_hook.test-provisioning-hook", "name", hookName2),
					resource.TestCheckResourceAttr(
						"ibm_compute_provisioning_hook.test-provisioning-hook", "uri", uri2),
				),
			},
		},
	})
}

func TestAccIBMComputeProvisioningHookWithTag(t *testing.T) {
	var hook datatypes.Provisioning_Hook

	hookName1 := fmt.Sprintf("%s%s", "tfuathook", acctest.RandString(10))
	uri1 := "http://www.weather.com"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMComputeProvisioningHookDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMComputeProvisioningHookWithTag(hookName1, uri1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMComputeProvisioningHookExists("ibm_compute_provisioning_hook.test-provisioning-hook", &hook),
					testAccCheckIBMComputeProvisioningHookAttributes(&hook, hookName1, uri1),
					resource.TestCheckResourceAttr(
						"ibm_compute_provisioning_hook.test-provisioning-hook", "name", hookName1),
					resource.TestCheckResourceAttr(
						"ibm_compute_provisioning_hook.test-provisioning-hook", "uri", uri1),
					resource.TestCheckResourceAttr(
						"ibm_compute_provisioning_hook.test-provisioning-hook", "tags.#", "2"),
				),
			},

			resource.TestStep{
				Config: testAccCheckIBMComputeProvisioningHookWithUpdatedTag(hookName1, uri1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMComputeProvisioningHookExists("ibm_compute_provisioning_hook.test-provisioning-hook", &hook),
					resource.TestCheckResourceAttr(
						"ibm_compute_provisioning_hook.test-provisioning-hook", "name", hookName1),
					resource.TestCheckResourceAttr(
						"ibm_compute_provisioning_hook.test-provisioning-hook", "uri", uri1),
					resource.TestCheckResourceAttr(
						"ibm_compute_provisioning_hook.test-provisioning-hook", "tags.#", "3"),
				),
			},
		},
	})
}

func testAccCheckIBMComputeProvisioningHookDestroy(s *terraform.State) error {
	service := services.GetProvisioningHookService(testAccProvider.Meta().(ClientSession).SoftLayerSession())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_compute_provisioning_hook" {
			continue
		}

		hookId, _ := strconv.Atoi(rs.Primary.ID)

		// Try to find the provisioning hook
		_, err := service.Id(hookId).GetObject()

		if err == nil {
			return fmt.Errorf("Provisioning Hook still exists")
		}
	}

	return nil
}

func testAccCheckIBMComputeProvisioningHookAttributes(hook *datatypes.Provisioning_Hook, name, uri string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *hook.Name != name {
			return fmt.Errorf("Bad name: %s", *hook.Name)
		}

		if *hook.Uri != uri {
			return fmt.Errorf("Bad uri: %s", *hook.Uri)
		}

		return nil
	}
}

func testAccCheckIBMComputeProvisioningHookExists(n string, hook *datatypes.Provisioning_Hook) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		hookId, _ := strconv.Atoi(rs.Primary.ID)

		service := services.GetProvisioningHookService(testAccProvider.Meta().(ClientSession).SoftLayerSession())
		foundHook, err := service.Id(hookId).GetObject()

		if err != nil {
			return err
		}

		if strconv.Itoa(int(*foundHook.Id)) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*hook = foundHook

		return nil
	}
}

func testAccCheckIBMComputeProvisioningHookConfig(name, uri string) string {
	return fmt.Sprintf(`
resource "ibm_compute_provisioning_hook" "test-provisioning-hook" {
    name = "%s"
    uri = "%s"
}`, name, uri)
}

func testAccCheckIBMComputeProvisioningHookWithTag(name, uri string) string {
	return fmt.Sprintf(`
resource "ibm_compute_provisioning_hook" "test-provisioning-hook" {
    name = "%s"
	uri = "%s"
	tags = ["one", "two"]
}`, name, uri)
}

func testAccCheckIBMComputeProvisioningHookWithUpdatedTag(name, uri string) string {
	return fmt.Sprintf(`
resource "ibm_compute_provisioning_hook" "test-provisioning-hook" {
    name = "%s"
	uri = "%s"
	tags = ["one", "two", "three"]
}`, name, uri)
}
