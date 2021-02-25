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

func TestAccIBMDLProviderGWsDataSource_basic(t *testing.T) {
	name := "dl_provider_gws"
	resName := "data.ibm_dl_provider_gateways.test_dl_provider_gws"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMDLProviderGWsDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
		},
	})
}

func testAccCheckIBMDLProviderGWsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	   data "ibm_dl_provider_gateways" "test_%s" {
	   }
	  `, name)
}
