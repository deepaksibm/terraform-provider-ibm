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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	registryv1 "github.com/IBM-Cloud/bluemix-go/api/container/registryv1"
)

func dataIBMContainerRegistryNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataIBMContainerRegistryNamespacesRead,

		Schema: map[string]*schema.Schema{
			"namespaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Container Registry Namespaces",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Container Registry Namespace name",
						},
						"resource_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource Group to which namespace has to be assigned",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CRN of the Namespace",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created Date",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updated Date",
						},
					},
				},
			},
		},
	}
}

func dataIBMContainerRegistryNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	accountID := userDetails.userAccount

	crClient, err := meta.(ClientSession).ContainerRegistryAPI()
	if err != nil {
		return err
	}
	target := registryv1.NamespaceTargetHeader{
		AccountID: accountID,
	}

	crAPI := crClient.Namespaces()

	response, err := crAPI.GetDetailedNamespaces(target)
	if err != nil {
		return err
	}
	namespaces := []map[string]interface{}{}
	for _, ns := range response {
		namespace := map[string]interface{}{}
		namespace["name"] = ns.Name
		namespace["resource_group_id"] = ns.ResourceGroup
		namespace["crn"] = ns.CRN
		namespace["created_on"] = ns.CreatedDate
		namespace["updated_on"] = ns.UpdatedDate
		namespaces = append(namespaces, namespace)
	}
	d.Set("namespaces", namespaces)
	d.SetId(time.Now().UTC().String())
	return nil
}
