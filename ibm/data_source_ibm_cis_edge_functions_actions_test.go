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

func TestAccIBMCisEdgeFunctionsActionsDataSource_basic(t *testing.T) {
	node := "data.ibm_cis_edge_functions_actions.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCisEdgeFunctionsActionsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "cis_edge_functions_actions.0.etag"),
				),
			},
		},
	})
}

func testAccCheckIBMCisEdgeFunctionsActionsDataSourceConfig() string {
	testName := "tf-acctest-basic"
	scriptName := "sample_script"

	return testAccCheckIBMCisEdgeFunctionsActionBasic(testName, scriptName) + fmt.Sprintf(`
	data "ibm_cis_edge_functions_actions" "test" {
		cis_id    = ibm_cis_edge_functions_action.tf-acctest-basic.cis_id
		domain_id = ibm_cis_edge_functions_action.tf-acctest-basic.domain_id
	  }`)
}
