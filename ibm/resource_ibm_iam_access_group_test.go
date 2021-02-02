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

	"github.com/IBM-Cloud/bluemix-go/models"

	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIBMIAMAccessGroup_Basic(t *testing.T) {
	var conf models.AccessGroupV2
	name := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	updateName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMIAMAccessGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIAMAccessGroupBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIAMAccessGroupExists("ibm_iam_access_group.accgroup", conf),
					resource.TestCheckResourceAttr("ibm_iam_access_group.accgroup", "name", name),
					resource.TestCheckResourceAttr("ibm_iam_access_group.accgroup", "tags.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIAMAccessGroupUpdateWithSameName(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIAMAccessGroupExists("ibm_iam_access_group.accgroup", conf),
					resource.TestCheckResourceAttr("ibm_iam_access_group.accgroup", "name", name),
					resource.TestCheckResourceAttr("ibm_iam_access_group.accgroup", "description", "AccessGroup for test scenario1"),
					resource.TestCheckResourceAttr("ibm_iam_access_group.accgroup", "tags.#", "3"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIAMAccessGroupUpdate(updateName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_iam_access_group.accgroup", "name", updateName),
					resource.TestCheckResourceAttr("ibm_iam_access_group.accgroup", "description", "AccessGroup for test scenario2"),
					resource.TestCheckResourceAttr("ibm_iam_access_group.accgroup", "tags.#", "1"),
				),
			},
		},
	})
}

func TestAccIBMIAMAccessGroup_import(t *testing.T) {
	var conf models.AccessGroupV2
	name := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	resourceName := "ibm_iam_access_group.accgroup"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMIAMAccessGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIAMAccessGroupTag(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIAMAccessGroupExists(resourceName, conf),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "AccessGroup for test scenario2"),
				),
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIAMAccessGroupDestroy(s *terraform.State) error {
	accClient, err := testAccProvider.Meta().(ClientSession).IAMUUMAPIV2()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_iam_access_group" {
			continue
		}

		agID := rs.Primary.ID

		// Try to find the key
		_, _, err := accClient.AccessGroup().Get(agID)

		if err == nil {
			return fmt.Errorf("Access group still exists: %s", rs.Primary.ID)
		} else if !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("Error waiting for access group (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIBMIAMAccessGroupExists(n string, obj models.AccessGroupV2) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		accClient, err := testAccProvider.Meta().(ClientSession).IAMUUMAPIV2()
		if err != nil {
			return err
		}
		agID := rs.Primary.ID

		accgroup, _, err := accClient.AccessGroup().Get(agID)

		if err != nil {
			return err
		}

		obj = *accgroup
		return nil
	}
}

func testAccCheckIBMIAMAccessGroupBasic(name string) string {
	return fmt.Sprintf(`
		
		resource "ibm_iam_access_group" "accgroup" {
			name = "%s"
			tags = ["tag1", "tag2"]
	  	}
	`, name)
}

func testAccCheckIBMIAMAccessGroupUpdateWithSameName(name string) string {
	return fmt.Sprintf(`
		
		resource "ibm_iam_access_group" "accgroup" {
			name        = "%s"
			description = "AccessGroup for test scenario1"
			tags        = ["tag1", "tag2", "db"]
	  	}
	`, name)
}

func testAccCheckIBMIAMAccessGroupUpdate(updateName string) string {
	return fmt.Sprintf(`

		resource "ibm_iam_access_group" "accgroup" {
			name        = "%s"
			description = "AccessGroup for test scenario2"
			tags        = ["tag1"]
	 	}
	`, updateName)
}

func testAccCheckIBMIAMAccessGroupTag(name string) string {
	return fmt.Sprintf(`

		resource "ibm_iam_access_group" "accgroup" {
			name              = "%s"		
			description       = "AccessGroup for test scenario2"
		}
	`, name)
}
