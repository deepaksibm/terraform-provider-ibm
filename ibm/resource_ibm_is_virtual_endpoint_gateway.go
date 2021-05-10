// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/internal/hashcode"
	"github.com/IBM/go-sdk-core/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isVirtualEndpointGatewayName               = "name"
	isVirtualEndpointGatewayResourceType       = "resource_type"
	isVirtualEndpointGatewayResourceGroupID    = "resource_group"
	isVirtualEndpointGatewayCreatedAt          = "created_at"
	isVirtualEndpointGatewayIPs                = "ips"
	isVirtualEndpointGatewayIPsID              = "id"
	isVirtualEndpointGatewayIPsAutoDelete      = "auto_delete"
	isVirtualEndpointGatewayIPsAddress         = "address"
	isVirtualEndpointGatewayIPsName            = "name"
	isVirtualEndpointGatewayIPsSubnet          = "subnet"
	isVirtualEndpointGatewayIPsResourceType    = "resource_type"
	isVirtualEndpointGatewayHealthState        = "health_state"
	isVirtualEndpointGatewayLifecycleState     = "lifecycle_state"
	isVirtualEndpointGatewayTarget             = "target"
	isVirtualEndpointGatewayTargetName         = "name"
	isVirtualEndpointGatewayTargetCRN          = "crn"
	isVirtualEndpointGatewayTargetResourceType = "resource_type"
	isVirtualEndpointGatewayVpcID              = "vpc"
	isVirtualEndpointGatewayTags               = "tags"
)

func resourceIBMISEndpointGateway() *schema.Resource {
	targetNameFmt := fmt.Sprintf("%s.0.%s", isVirtualEndpointGatewayTarget, isVirtualEndpointGatewayTargetName)
	targetCRNFmt := fmt.Sprintf("%s.0.%s", isVirtualEndpointGatewayTarget, isVirtualEndpointGatewayTargetCRN)
	return &schema.Resource{
		Create:   resourceIBMisVirtualEndpointGatewayCreate,
		Read:     resourceIBMisVirtualEndpointGatewayRead,
		Update:   resourceIBMisVirtualEndpointGatewayUpdate,
		Delete:   resourceIBMisVirtualEndpointGatewayDelete,
		Exists:   resourceIBMisVirtualEndpointGatewayExists,
		Importer: &schema.ResourceImporter{},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return resourceTagsCustomizeDiff(diff)
			},
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			isVirtualEndpointGatewayName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateISName,
				Description:  "Endpoint gateway name",
			},
			isVirtualEndpointGatewayResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway resource type",
			},
			isVirtualEndpointGatewayResourceGroupID: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The resource group id",
			},
			isVirtualEndpointGatewayCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway created date and time",
			},
			isVirtualEndpointGatewayHealthState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway health state",
			},
			isVirtualEndpointGatewayLifecycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway lifecycle state",
			},
			isVirtualEndpointGatewayIPs: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Endpoint gateway resource group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVirtualEndpointGatewayIPsID: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The IPs id",
						},
						isVirtualEndpointGatewayIPsAutoDelete: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Automatically delete reserved IP when the target is deleted or when the reserved IP is unbound.",
						},
						isVirtualEndpointGatewayIPsName: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The IPs name",
						},
						isVirtualEndpointGatewayIPsSubnet: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Subnet id",
						},
						isVirtualEndpointGatewayIPsResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC Resource Type",
						},
						isVirtualEndpointGatewayIPsAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address",
						},
					},
				},
			},
			isVirtualEndpointGatewayTarget: {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "Endpoint gateway target",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVirtualEndpointGatewayTargetName: {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							AtLeastOneOf: []string{
								targetNameFmt,
								targetCRNFmt,
							},
							Description: "The target name",
						},
						isVirtualEndpointGatewayTargetCRN: {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							AtLeastOneOf: []string{
								targetNameFmt,
								targetCRNFmt,
							},
							Description: "The target crn",
						},
						isVirtualEndpointGatewayTargetResourceType: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The target resource type",
						},
					},
				},
			},
			isVirtualEndpointGatewayVpcID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPC id",
			},
			isVirtualEndpointGatewayTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_is_virtual_endpoint_gateway", "tag")},
				Set:         resourceIBMVPCHash,
				Description: "List of tags for VPE",
			},
		},
	}
}

func resourceIBMISEndpointGatewayValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "tag",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmEndpointGatewayResourceValidator := ResourceValidator{ResourceName: "ibm_is_virtual_endpoint_gateway", Schema: validateSchema}
	return &ibmEndpointGatewayResourceValidator
}

func resourceIBMisVirtualEndpointGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	name := d.Get(isVirtualEndpointGatewayName).(string)

	// target opiton
	targetOpt := &vpcv1.EndpointGatewayTargetPrototype{}
	targetNameFmt := fmt.Sprintf("%s.0.%s", isVirtualEndpointGatewayTarget, isVirtualEndpointGatewayTargetName)
	targetCRNFmt := fmt.Sprintf("%s.0.%s", isVirtualEndpointGatewayTarget, isVirtualEndpointGatewayTargetCRN)
	targetResourceTypeFmt := fmt.Sprintf("%s.0.%s", isVirtualEndpointGatewayTarget, isVirtualEndpointGatewayTargetResourceType)
	targetOpt.ResourceType = core.StringPtr(d.Get(targetResourceTypeFmt).(string))
	if v, ok := d.GetOk(targetNameFmt); ok {
		targetOpt.Name = core.StringPtr(v.(string))
	}
	if v, ok := d.GetOk(targetCRNFmt); ok {
		targetOpt.CRN = core.StringPtr(v.(string))
	}

	// vpc option
	vpcID := d.Get(isVirtualEndpointGatewayVpcID).(string)
	vpcOpt := &vpcv1.VPCIdentity{
		ID: core.StringPtr(vpcID),
	}

	// update option
	opt := sess.NewCreateEndpointGatewayOptions(targetOpt, vpcOpt)
	opt.SetName(name)
	opt.SetTarget(targetOpt)
	opt.SetVPC(vpcOpt)

	// IPs option
	if ips, ok := d.GetOk(isVirtualEndpointGatewayIPs); ok {
		opt.SetIps(expandIPs(ips.(*schema.Set).List()))
	}

	// Resource group option
	if resourceGroup, ok := d.GetOk(isVirtualEndpointGatewayResourceGroupID); ok {
		resourceGroupID := resourceGroup.(string)

		resourceGroupOpt := &vpcv1.ResourceGroupIdentity{
			ID: core.StringPtr(resourceGroupID),
		}
		opt.SetResourceGroup(resourceGroupOpt)

	}
	result, response, err := sess.CreateEndpointGateway(opt)
	if err != nil {
		log.Printf("Create Endpoint Gateway failed: %v", response)
		return err
	}

	d.SetId(*result.ID)
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isVirtualEndpointGatewayTags); ok || v != "" {
		oldList, newList := d.GetChange(isVirtualEndpointGatewayTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *result.CRN)
		if err != nil {
			log.Printf(
				"Error on create of VPE (%s) tags: %s", d.Id(), err)
		}
	}
	return resourceIBMisVirtualEndpointGatewayRead(d, meta)
}

func resourceIBMisVirtualEndpointGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	if d.HasChange(isVirtualEndpointGatewayName) {
		name := d.Get(isVirtualEndpointGatewayName).(string)

		// create option
		endpointGatewayPatchModel := new(vpcv1.EndpointGatewayPatch)
		endpointGatewayPatchModel.Name = core.StringPtr(name)
		endpointGatewayPatchModelAsPatch, _ := endpointGatewayPatchModel.AsPatch()
		opt := sess.NewUpdateEndpointGatewayOptions(d.Id(), endpointGatewayPatchModelAsPatch)
		_, response, err := sess.UpdateEndpointGateway(opt)
		if err != nil {
			log.Printf("Update Endpoint Gateway failed: %v", response)
			return err
		}

	}

	if d.HasChange(isVirtualEndpointGatewayIPs) {
		o, n := d.GetChange(isVirtualEndpointGatewayIPs)
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))
		gatewayID := d.Id()
		for _, ips := range ors.List() {
			ipsMap := ips.(map[string]interface{})

			ipID := ipsMap["id"].(string)

			opt := sess.NewRemoveEndpointGatewayIPOptions(gatewayID, ipID)
			response, err := sess.RemoveEndpointGatewayIP(opt)
			if err != nil && response.StatusCode != 404 {
				log.Printf("Remove Endpoint Gateway IP failed: %v", response)
				return err
			}
		}

		// Make sure we save the state of the currently configured rules
		routes := o.(*schema.Set).Intersection(n.(*schema.Set))
		d.Set(isVirtualEndpointGatewayIPs, routes)

		// Then loop through all the newly configured routes and create them
		for _, ips := range nrs.List() {
			ipsMap := ips.(map[string]interface{})
			ipID := ipsMap["id"].(string)
			if ipID == "" {
				subnetID := ipsMap[isSubNetID].(string)
				options := sess.NewCreateSubnetReservedIPOptions(subnetID)

				nameStr := ipsMap[isReservedIPName].(string)

				if nameStr != "" {
					options.Name = &nameStr
				}

				autoDeleteBool := ipsMap[isReservedIPAutoDelete].(bool)
				options.AutoDelete = &autoDeleteBool

				rip, response, err := sess.CreateSubnetReservedIP(options)
				if err != nil || response == nil || rip == nil {
					return fmt.Errorf("Error creating the reserved IP: %s\n%s", err, response)
				}
				ipID = *rip.ID
			}

			opt := sess.NewAddEndpointGatewayIPOptions(gatewayID, ipID)
			_, response, err := sess.AddEndpointGatewayIP(opt)
			if err != nil {
				log.Printf("Add Endpoint Gateway failed: %v", response)
				return err
			}
		}

	}

	if d.HasChange(isVirtualEndpointGatewayTags) {
		opt := sess.NewGetEndpointGatewayOptions(d.Id())
		result, response, err := sess.GetEndpointGateway(opt)
		if err != nil {
			return fmt.Errorf("Error getting VPE: %s\n%s", err, response)
		}
		oldList, newList := d.GetChange(isVirtualEndpointGatewayTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *result.CRN)
		if err != nil {
			log.Printf(
				"Error on update of VPE (%s) tags: %s", d.Id(), err)
		}
	}
	return resourceIBMisVirtualEndpointGatewayRead(d, meta)
}

func resourceIBMisVirtualEndpointGatewayRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	// read option
	opt := sess.NewGetEndpointGatewayOptions(d.Id())
	result, response, err := sess.GetEndpointGateway(opt)
	if err != nil {
		log.Printf("Get Endpoint Gateway failed: %v", response)
		return err
	}
	d.Set(isVirtualEndpointGatewayName, result.Name)
	d.Set(isVirtualEndpointGatewayHealthState, result.HealthState)
	d.Set(isVirtualEndpointGatewayCreatedAt, result.CreatedAt.String())
	d.Set(isVirtualEndpointGatewayLifecycleState, result.LifecycleState)
	d.Set(isVirtualEndpointGatewayResourceType, result.ResourceType)
	d.Set(isVirtualEndpointGatewayIPs, flattenIPs(result.Ips))
	d.Set(isVirtualEndpointGatewayResourceGroupID, result.ResourceGroup.ID)
	d.Set(isVirtualEndpointGatewayTarget,
		flattenEndpointGatewayTarget(result.Target.(*vpcv1.EndpointGatewayTarget)))
	d.Set(isVirtualEndpointGatewayVpcID, result.VPC.ID)
	tags, err := GetTagsUsingCRN(meta, *result.CRN)
	if err != nil {
		log.Printf(
			"Error on get of VPE (%s) tags: %s", d.Id(), err)
	}
	d.Set(isVirtualEndpointGatewayTags, tags)
	return nil
}

func resourceIBMisVirtualEndpointGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	opt := sess.NewDeleteEndpointGatewayOptions(d.Id())
	response, err := sess.DeleteEndpointGateway(opt)
	if err != nil {
		log.Printf("Delete Endpoint Gateway failed: %v", response)
	}
	return nil
}

func resourceIBMisVirtualEndpointGatewayExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	opt := sess.NewGetEndpointGatewayOptions(d.Id())
	_, response, err := sess.GetEndpointGateway(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("Endpoint Gateway does not exist.")
			return false, nil
		}
		log.Printf("Error : %s", response)
		return false, err
	}
	return true, nil
}

func expandIPs(ipsSet []interface{}) (ipsOptions []vpcv1.EndpointGatewayReservedIPIntf) {
	ipsList := ipsSet
	for _, item := range ipsList {
		ips := item.(map[string]interface{})
		// IPs option
		ipsID := ips[isVirtualEndpointGatewayIPsID].(string)
		ipsName := ips[isVirtualEndpointGatewayIPsName].(string)
		var ipsautoDeleteBool bool
		if ipsautoDeleteBoolIntf := ips[isVirtualEndpointGatewayIPsAutoDelete]; ipsautoDeleteBoolIntf != nil {
			ipsautoDeleteBool = ipsautoDeleteBoolIntf.(bool)
		}

		// IPs subnet option
		ipsSubnetID := ips[isVirtualEndpointGatewayIPsSubnet].(string)

		ipsSubnetOpt := &vpcv1.SubnetIdentity{
			ID: &ipsSubnetID,
		}

		ipsOpt := &vpcv1.EndpointGatewayReservedIP{
			ID:         core.StringPtr(ipsID),
			Name:       core.StringPtr(ipsName),
			Subnet:     ipsSubnetOpt,
			AutoDelete: core.BoolPtr(ipsautoDeleteBool),
		}
		ipsOptions = append(ipsOptions, ipsOpt)
	}
	return ipsOptions
}

func flattenIPs(ipsList []vpcv1.ReservedIPReference) interface{} {
	ipsListOutput := make([]interface{}, 0)
	for _, item := range ipsList {
		ips := make(map[string]interface{}, 0)
		ips[isVirtualEndpointGatewayIPsID] = *item.ID
		ips[isVirtualEndpointGatewayIPsName] = *item.Name
		ips[isVirtualEndpointGatewayIPsSubnet] = strings.Split(*item.Href, "/")[5]
		ips[isVirtualEndpointGatewayIPsResourceType] = *item.ResourceType
		ips[isVirtualEndpointGatewayIPsAddress] = *item.Address
		ipsListOutput = append(ipsListOutput, ips)
	}
	return ipsListOutput
}

func flattenEndpointGatewayTarget(target *vpcv1.EndpointGatewayTarget) interface{} {
	targetSlice := []interface{}{}
	targetOutput := map[string]string{}
	if target == nil {
		return targetOutput
	}
	if target.Name != nil {
		targetOutput[isVirtualEndpointGatewayTargetName] = *target.Name
	}
	if target.CRN != nil {
		targetOutput[isVirtualEndpointGatewayTargetCRN] = *target.CRN
	}
	if target.ResourceType != nil {
		targetOutput[isVirtualEndpointGatewayTargetResourceType] = *target.ResourceType
	}
	targetSlice = append(targetSlice, targetOutput)
	return targetSlice
}

func resourceIbmVirtualEndPointGatewayHash(v interface{}) int {
	var buf bytes.Buffer
	m, castOk := v.(map[string]interface{})
	if !castOk {
		return 0
	}

	if v, ok := m[isVirtualEndpointGatewayIPsID]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m[isVirtualEndpointGatewayIPsAutoDelete]; ok {
		buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
	}

	if v, ok := m[isVirtualEndpointGatewayIPsName]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m[isVirtualEndpointGatewayIPsSubnet]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m[isVirtualEndpointGatewayIPsResourceType]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m[isVirtualEndpointGatewayIPsAddress]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}
