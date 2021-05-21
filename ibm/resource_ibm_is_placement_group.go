/**
 * (C) Copyright IBM Corp. 2021.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

const (
	isPlacementGroupDeleting   = "deleting"
	isPlacementGroupStable     = "stable"
	isPlacementGroupFailed     = "failed"
	isPlacementGroupDeleteDone = "done"
	isPlacementGroupPending    = "pending"
	isPlacementGroupWaiting    = "waiting"
	isPlacementGroupUpdating   = "updating"
	isPlacementGroupInUse      = "inuse"
	isPlacementGroupSuspended  = "suspended"
)

func resourceIbmIsPlacementGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIsPlacementGroupCreate,
		ReadContext:   resourceIbmIsPlacementGroupRead,
		UpdateContext: resourceIbmIsPlacementGroupUpdate,
		DeleteContext: resourceIbmIsPlacementGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"strategy": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: InvokeValidator("ibm_is_placement_group", "strategy"),
				Description:  "The strategy for this placement group- `host_spread`: place on different compute hosts- `power_spread`: place on compute hosts that use different power sourcesThe enumerated values for this property may expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the placement group on which the unexpected strategy was encountered.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_placement_group", "name"),
				Description:  "The unique user-defined name for this placement group. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the placement group was created.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this placement group.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this placement group.",
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the placement group.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

func resourceIbmIsPlacementGroupValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "strategy",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "host_spread, power_spread",
		},
		ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)

	resourceValidator := ResourceValidator{ResourceName: "ibm_is_placement_group", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmIsPlacementGroupCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createPlacementGroupOptions := &vpcv1.CreatePlacementGroupOptions{}

	createPlacementGroupOptions.SetStrategy(d.Get("strategy").(string))
	if pgnameIntf, ok := d.GetOk("name"); ok {
		createPlacementGroupOptions.SetName(pgnameIntf.(string))
	}
	if resourceGroupIntf, ok := d.GetOk("resource_group"); ok {
		resourceGroup := resourceGroupIntf.(string)
		resourceGroupIdentity := &vpcv1.ResourceGroupIdentity{
			ID: &resourceGroup,
		}
		createPlacementGroupOptions.SetResourceGroup(resourceGroupIdentity)
	}

	placementGroup, response, err := vpcClient.CreatePlacementGroupWithContext(context, createPlacementGroupOptions)
	if err != nil {
		log.Printf("[DEBUG] CreatePlacementGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(*placementGroup.ID)

	_, err = isWaitForPlacementGroupAvailable(vpcClient, d.Id(), d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceIbmIsPlacementGroupRead(context, d, meta)
}

func resourceIbmIsPlacementGroupRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getPlacementGroupOptions := &vpcv1.GetPlacementGroupOptions{}

	getPlacementGroupOptions.SetID(d.Id())

	placementGroup, response, err := vpcClient.GetPlacementGroupWithContext(context, getPlacementGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetPlacementGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	if err = d.Set("strategy", placementGroup.Strategy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting strategy: %s", err))
	}
	if err = d.Set("name", placementGroup.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if placementGroup.ResourceGroup != nil {
		if err = d.Set("resource_group", *placementGroup.ResourceGroup.ID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group: %s", err))
		}
	}
	if err = d.Set("created_at", placementGroup.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("crn", placementGroup.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("href", placementGroup.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", placementGroup.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("resource_type", placementGroup.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	return nil
}

func resourceIbmIsPlacementGroupResourceGroupIdentityToMap(resourceGroupIdentity vpcv1.ResourceGroupIdentity) map[string]interface{} {
	resourceGroupIdentityMap := map[string]interface{}{}

	resourceGroupIdentityMap["id"] = resourceGroupIdentity.ID

	return resourceGroupIdentityMap
}

func resourceIbmIsPlacementGroupResourceGroupIdentityByIDToMap(resourceGroupIdentityByID vpcv1.ResourceGroupIdentityByID) map[string]interface{} {
	resourceGroupIdentityByIDMap := map[string]interface{}{}

	resourceGroupIdentityByIDMap["id"] = resourceGroupIdentityByID.ID

	return resourceGroupIdentityByIDMap
}

func resourceIbmIsPlacementGroupUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updatePlacementGroupOptions := &vpcv1.UpdatePlacementGroupOptions{}

	updatePlacementGroupOptions.SetID(d.Id())

	hasChange := false

	placementGroupPatchModel := &vpcv1.PlacementGroupPatch{}
	if d.HasChange("name") {
		plName := d.Get("name").(string)
		placementGroupPatchModel.Name = &plName
		hasChange = true
	}
	if hasChange {
		placementGroupPatch, err := placementGroupPatchModel.AsPatch()
		if err != nil {
			log.Printf("[DEBUG] Error calling AsPatch for PlacementGroupPatch %s", err)
			return diag.FromErr(err)
		}
		updatePlacementGroupOptions.SetPlacementGroupPatch(placementGroupPatch)
		_, response, err := vpcClient.UpdatePlacementGroupWithContext(context, updatePlacementGroupOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdatePlacementGroupWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
	}

	return resourceIbmIsPlacementGroupRead(context, d, meta)
}

func resourceIbmIsPlacementGroupDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deletePlacementGroupOptions := &vpcv1.DeletePlacementGroupOptions{}

	deletePlacementGroupOptions.SetID(d.Id())

	response, err := vpcClient.DeletePlacementGroupWithContext(context, deletePlacementGroupOptions)
	if err != nil {
		if response.StatusCode == 409 {
			_, err = isWaitForPlacementGroupDeleteRetry(vpcClient, d, d.Id())
		}
	}
	_, err = isWaitForPlacementGroupDelete(vpcClient, d, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func isWaitForPlacementGroupDelete(vpcClient *vpcv1.VpcV1, d *schema.ResourceData, id string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isPlacementGroupDeleting, isPlacementGroupStable, isPlacementGroupPending, isPlacementGroupWaiting, isPlacementGroupUpdating},
		Target:  []string{isPlacementGroupDeleteDone, ""},
		Refresh: func() (interface{}, string, error) {
			getPlacementGroupOptions := &vpcv1.GetPlacementGroupOptions{
				ID: &id,
			}

			placementGroup, response, err := vpcClient.GetPlacementGroup(getPlacementGroupOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return placementGroup, isPlacementGroupDeleteDone, nil
				} else if response != nil && response.StatusCode == 409 {
					return placementGroup, *placementGroup.LifecycleState, fmt.Errorf("The  PLacementGroup %s failed to delete: %v", id, err)
				}
				return nil, "", fmt.Errorf("Error Getting PLacementGroup: %s\n%s", err, response)
			}
			if *placementGroup.LifecycleState == isPlacementGroupFailed {
				return placementGroup, *placementGroup.LifecycleState, fmt.Errorf("The  PLacementGroup %s failed to delete: %v", id, err)
			}
			return placementGroup, isPlacementGroupDeleting, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isWaitForPlacementGroupDeleteRetry(vpcClient *vpcv1.VpcV1, d *schema.ResourceData, id string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isPlacementGroupInUse},
		Target:  []string{isPlacementGroupDeleting, isPlacementGroupDeleteDone, ""},
		Refresh: func() (interface{}, string, error) {
			deletePlacementGroupOptions := &vpcv1.DeletePlacementGroupOptions{}

			deletePlacementGroupOptions.SetID(id)

			response, err := vpcClient.DeletePlacementGroup(deletePlacementGroupOptions)
			if err != nil {
				if response != nil && response.StatusCode == 409 {
					return nil, isPlacementGroupInUse, err
				} else if response != nil && response.StatusCode == 404 {
					return nil, isPlacementGroupDeleteDone, nil
				}
				return nil, "", fmt.Errorf("Error deleting PLacementGroup: %s\n%s", err, response)
			}
			return nil, isPlacementGroupDeleting, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isWaitForPlacementGroupAvailable(vpcClient *vpcv1.VpcV1, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for placement group (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isPlacementGroupPending, isPlacementGroupWaiting, isPlacementGroupUpdating},
		Target:     []string{isPlacementGroupFailed, isPlacementGroupStable, isPlacementGroupSuspended, ""},
		Refresh:    isPlacementGroupRefreshFunc(vpcClient, id, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isPlacementGroupRefreshFunc(vpcClient *vpcv1.VpcV1, id string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getinsOptions := &vpcv1.GetPlacementGroupOptions{
			ID: &id,
		}
		placementGroup, response, err := vpcClient.GetPlacementGroup(getinsOptions)
		if placementGroup == nil || err != nil {
			return nil, "", fmt.Errorf("Error getting placementGroup : %s\n%s", err, response)
		}

		d.Set("lifecycle_state", *placementGroup.LifecycleState)

		if *placementGroup.LifecycleState == isPlacementGroupSuspended || *placementGroup.LifecycleState == isPlacementGroupFailed {

			return placementGroup, *placementGroup.LifecycleState, fmt.Errorf("status of placement group is %s : \n%s", *placementGroup.LifecycleState, response)

		}
		return placementGroup, *placementGroup.LifecycleState, nil
	}
}
