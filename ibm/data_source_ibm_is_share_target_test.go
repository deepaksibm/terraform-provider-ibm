// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIbmIsShareTargetDataSourceAllArgs(t *testing.T) {
	vpcName := fmt.Sprintf("tf-vpc-%d", acctest.RandIntRange(10, 100))
	targetName := fmt.Sprintf("tf-share-target-%d", acctest.RandIntRange(10, 100))
	shareName := fmt.Sprintf("tf-fs-name-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIbmIsShareTargetDataSourceConfigBasic(shareName, vpcName, targetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_share_target.is_share_target", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_target.is_share_target", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_target.is_share_target", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_target.is_share_target", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_target.is_share_target", "mount_path"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_target.is_share_target", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_target.is_share_target", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_target.is_share_target", "vpc.#"),
				),
			},
		},
	})
}

func testAccCheckIbmIsShareTargetDataSourceConfigBasic(sname, vpcName, targetName string) string {
	return testAccCheckIbmIsShareTargetsDataSourceConfigBasic(sname, vpcName, targetName) + fmt.Sprintf(`
		
		data "ibm_is_share_target" "is_share_target" {
			share = ibm_is_share.is_share.id
			share_target = data.ibm_is_share_targets.is_share_targets.share_targets.0.id
		}
	`)
}
