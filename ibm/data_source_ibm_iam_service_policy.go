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

	"github.com/IBM-Cloud/bluemix-go/api/iampap/iampapv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source to find all the policies for a serviceID
func dataSourceIBMIAMServicePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMIAMServicePolicyRead,

		Schema: map[string]*schema.Schema{
			"iam_service_id": {
				Description: "UUID of ServiceID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"sort": {
				Description: "Sort query for policies",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"roles": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Role names of the policy definition",
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service name of the policy definition",
									},
									"resource_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "ID of resource instance of the policy definition",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region of the policy definition",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource type of the policy definition",
									},
									"resource": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource of the policy definition",
									},
									"resource_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the resource group.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIAMServicePolicyRead(d *schema.ResourceData, meta interface{}) error {
	iamClient, err := meta.(ClientSession).IAMAPI()
	if err != nil {
		return err
	}

	serviceIDUUID := d.Get("iam_service_id").(string)

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	serviceID, err := iamClient.ServiceIds().Get(serviceIDUUID)
	if err != nil {
		return err
	}

	iampapClient, err := meta.(ClientSession).IAMPAPAPI()
	if err != nil {
		return err
	}

	query := iampapv1.SearchParams{
		AccountID: userDetails.userAccount,
		Type:      iampapv1.AccessPolicyType,
		IAMID:     serviceID.IAMID,
	}

	if v, ok := d.GetOk("sort"); ok {
		query.Sort = v.(string)
	}

	policies, err := iampapClient.V1Policy().List(query)
	if err != nil {
		return err
	}

	servicePolicies := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		roles := make([]string, len(policy.Roles))
		for i, role := range policy.Roles {
			roles[i] = role.Name
		}
		resources := flattenPolicyResource(policy.Resources)
		p := map[string]interface{}{
			"id":        fmt.Sprintf("%s/%s", serviceIDUUID, policy.ID),
			"roles":     roles,
			"resources": resources,
		}
		servicePolicies = append(servicePolicies, p)
	}
	d.SetId(serviceIDUUID)
	d.Set("policies", servicePolicies)
	return nil
}
