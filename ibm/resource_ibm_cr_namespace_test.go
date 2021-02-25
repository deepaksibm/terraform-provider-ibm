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
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	registryv1 "github.com/IBM-Cloud/bluemix-go/api/container/registryv1"
)

func TestAccIBMCrNamespaceBasic(t *testing.T) {

	namespaceName := fmt.Sprintf("terraform-tf-%d", acctest.RandIntRange(10, 100))
	namespaceName1 := fmt.Sprintf("terraform-tf-%d", acctest.RandIntRange(10, 100))
	resourceName := "ibm_cr_namespace.test_namespace"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCrNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCrNamespaceBasic(namespaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", namespaceName),
					resource.TestCheckResourceAttrSet(
						resourceName, "crn"),
				),
			},
			{
				Config: testAccCheckIBMCrNamespaceBasic(namespaceName1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", namespaceName1),
					resource.TestCheckResourceAttrSet(
						resourceName, "crn"),
				),
			},
			{
				Config: testAccCheckIBMCrNamespaceRGBasic(namespaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", namespaceName),
					resource.TestCheckResourceAttrSet(
						resourceName, "crn"),
				),
			},
		},
	})
}

func TestAccIBMCrNamespaceImportBasic(t *testing.T) {
	namespaceName := fmt.Sprintf("terraform-tf-%d", acctest.RandIntRange(10, 100))
	resourceName := "ibm_cr_namespace.test_namespace"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCrNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCrNamespaceBasic(namespaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", namespaceName),
					resource.TestCheckResourceAttrSet(
						resourceName, "crn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMCrNamespaceDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_cr_namespace" {
			continue
		}
		userDetails, err := testAccProvider.Meta().(ClientSession).BluemixUserDetails()
		if err != nil {
			return err
		}
		accountID := userDetails.userAccount

		crClient, err := testAccProvider.Meta().(ClientSession).ContainerRegistryAPI()
		if err != nil {
			return err
		}
		namespace := rs.Primary.ID
		target := registryv1.NamespaceTargetHeader{
			AccountID: accountID,
		}

		crAPI := crClient.Namespaces()
		response, err := crAPI.GetDetailedNamespaces(target)
		if err == nil {
			for _, ns := range response {
				if ns.Name == namespace {
					return fmt.Errorf("Error checking if namespace (%s) has been destroyed", rs.Primary.ID)
				}
			}
		}
	}
	return nil
}

func testAccCheckIBMCrNamespaceBasic(namespaceName string) string {
	return fmt.Sprintf(`
	resource "ibm_cr_namespace" "test_namespace"{
		name = "%s"
	}
	`, namespaceName)
}
func testAccCheckIBMCrNamespaceRGBasic(namespaceName string) string {
	return fmt.Sprintf(`
	resource "ibm_resource_group" "test_group" {
		name="%s"
	}
	resource "ibm_cr_namespace" "test_namespace"{
		name = "%[1]s"
		resource_group_id=ibm_resource_group.test_group.id
	}
	`, namespaceName)
}
