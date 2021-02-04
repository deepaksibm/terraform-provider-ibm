package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIBMPIVolumeDataSource_basic(t *testing.T) {
	name := fmt.Sprintf("tf-pi-volume-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPIVolumeDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_pi_volume.testacc_ds_volume", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMPIVolumeDataSourceConfig(name string) string {
	return testAccCheckIBMPIVolumeConfig(name) + fmt.Sprintf(`
data "ibm_pi_volume" "testacc_ds_volume" {
    pi_volume_name = ibm_pi_volume.power_volume.pi_volume_name
    pi_cloud_instance_id = "%s"
}`, pi_cloud_instance_id)

}
