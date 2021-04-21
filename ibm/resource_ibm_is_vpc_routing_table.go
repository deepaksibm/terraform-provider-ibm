// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/internal/hashcode"
	"github.com/IBM/go-sdk-core/v4/core"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	rtID                         = "routing_table"
	rtVpcID                      = "vpc"
	rtName                       = "name"
	rtRouteDirectLinkIngress     = "route_direct_link_ingress"
	rtRouteTransitGatewayIngress = "route_transit_gateway_ingress"
	rtRouteVPCZoneIngress        = "route_vpc_zone_ingress"
	rtCreateAt                   = "created_at"
	rtHref                       = "href"
	rtIsDefault                  = "is_default"
	rtResourceType               = "resource_type"
	rtLifecycleState             = "lifecycle_state"
	rtSubnets                    = "subnets"
	rtDestination                = "destination"
	rtAction                     = "action"
	rtNextHop                    = "next_hop"
	rtZone                       = "zone"
	rtOrigin                     = "origin"
	rtRoutes                     = "routes"
)

func resourceIBMISVPCRoutingTable() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISVPCRoutingTableCreate,
		Read:     resourceIBMISVPCRoutingTableRead,
		Update:   resourceIBMISVPCRoutingTableUpdate,
		Delete:   resourceIBMISVPCRoutingTableDelete,
		Exists:   resourceIBMISVPCRoutingTableExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			rtVpcID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPC identifier.",
			},
			rtRouteDirectLinkIngress: {
				Type:        schema.TypeBool,
				ForceNew:    false,
				Default:     false,
				Optional:    true,
				Description: "If set to true, this routing table will be used to route traffic that originates from Direct Link to this VPC.",
			},
			rtRouteTransitGatewayIngress: {
				Type:        schema.TypeBool,
				ForceNew:    false,
				Default:     false,
				Optional:    true,
				Description: "If set to true, this routing table will be used to route traffic that originates from Transit Gateway to this VPC.",
			},
			rtRouteVPCZoneIngress: {
				Type:        schema.TypeBool,
				ForceNew:    false,
				Default:     false,
				Optional:    true,
				Description: "If set to true, this routing table will be used to route traffic that originates from subnets in other zones in this VPC.",
			},
			rtName: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				Computed:     true,
				ValidateFunc: InvokeValidator("ibm_is_vpc_routing_table", rtName),
				Description:  "The user-defined name for this routing table.",
			},
			rtID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The routing table identifier.",
			},
			rtHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Routing table Href",
			},
			rtResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Routing table Resource Type",
			},
			rtCreateAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Routing table Created At",
			},
			rtLifecycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Routing table Lifecycle State",
			},
			rtIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this is the default routing table for this VPC",
			},
			rtSubnets: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						rtName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet name",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID",
						},
					},
				},
			},
			rtRoutes: {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						rtDestination: {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The destination of the route.",
						},
						rtZone: {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The zone to apply the route to. Traffic from subnets in this zone will be subject to this route.",
						},
						rtNextHop: {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "If action is deliver, the next hop that packets will be delivered to. For other action values, its address will be 0.0.0.0.",
						},
						rtAction: {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "deliver",
							Description:  "The action to perform with a packet matching the route.",
							ValidateFunc: InvokeValidator("ibm_is_vpc_routing_table", rAction),
						},
						rtName: {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "The user-defined name for this route.",
							ValidateFunc: InvokeValidator("ibm_is_vpc_routing_table", rName),
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this routing table route.",
						},
					},
				},
				Set: resourceIbmRouteTableHash,
			},
		},
	}
}

func resourceIBMISVPCRoutingTableValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 2)
	actionAllowedValues := "delegate, delegate_vpc, deliver, drop"

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 rtName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   false,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 rtAction,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   false,
			AllowedValues:              actionAllowedValues})

	ibmISVPCRoutingTableValidator := ResourceValidator{ResourceName: "ibm_is_vpc_routing_table", Schema: validateSchema}
	return &ibmISVPCRoutingTableValidator
}

func resourceIBMISVPCRoutingTableCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	vpcID := d.Get(rtVpcID).(string)
	rtName := d.Get(rtName).(string)

	createVpcRoutingTableOptions := sess.NewCreateVPCRoutingTableOptions(vpcID)
	createVpcRoutingTableOptions.SetName(rtName)
	if _, ok := d.GetOk(rtRouteDirectLinkIngress); ok {
		routeDirectLinkIngress := d.Get(rtRouteDirectLinkIngress).(bool)
		createVpcRoutingTableOptions.RouteDirectLinkIngress = &routeDirectLinkIngress
	}
	if _, ok := d.GetOk(rtRouteTransitGatewayIngress); ok {
		routeTransitGatewayIngress := d.Get(rtRouteTransitGatewayIngress).(bool)
		createVpcRoutingTableOptions.RouteTransitGatewayIngress = &routeTransitGatewayIngress
	}
	if _, ok := d.GetOk(rtRouteVPCZoneIngress); ok {
		routeVPCZoneIngress := d.Get(rtRouteVPCZoneIngress).(bool)
		createVpcRoutingTableOptions.RouteVPCZoneIngress = &routeVPCZoneIngress
	}

	if _, ok := d.GetOk("routes"); ok {
		routes := d.Get("routes").(*schema.Set).List()
		routesPrototypes := []vpcv1.RoutePrototype{}

		for _, routeItem := range routes {
			route := routeItem.(map[string]interface{})
			zone := route["zone"].(string)
			zoneidentity := &vpcv1.ZoneIdentity{
				Name: core.StringPtr(zone),
			}
			destination := route["destination"].(string)
			routePrototype := vpcv1.RoutePrototype{
				Zone:        zoneidentity,
				Destination: &destination,
			}

			nexthop := route["next_hop"].(string)
			if net.ParseIP(nexthop) == nil {
				nhConnectionID := &vpcv1.RouteNextHopPrototypeVPNGatewayConnectionIdentity{
					ID: core.StringPtr(nexthop),
				}
				routePrototype.NextHop = nhConnectionID
			} else {
				nh := &vpcv1.RouteNextHopPrototypeRouteNextHopIP{
					Address: core.StringPtr(nexthop),
				}
				routePrototype.NextHop = nh
			}

			if action, ok := route["action"]; ok {
				routeAction := action.(string)
				routePrototype.Action = &routeAction
			}

			if name, ok := route["name"]; ok {
				routeName := name.(string)
				routePrototype.Name = &routeName
			}
			routesPrototypes = append(routesPrototypes, routePrototype)
		}
		createVpcRoutingTableOptions.Routes = routesPrototypes
	}

	routeTable, response, err := sess.CreateVPCRoutingTable(createVpcRoutingTableOptions)
	if err != nil {
		log.Printf("[DEBUG] Create VPC Routing table err %s\n%s", err, response)
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", vpcID, *routeTable.ID))

	return resourceIBMISVPCRoutingTableRead(d, meta)
}

func resourceIBMISVPCRoutingTableRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	idSet := strings.Split(d.Id(), "/")
	getVpcRoutingTableOptions := sess.NewGetVPCRoutingTableOptions(idSet[0], idSet[1])
	routeTable, response, err := sess.GetVPCRoutingTable(getVpcRoutingTableOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting VPC Routing table: %s\n%s", err, response)
	}

	d.Set(rtID, routeTable.ID)
	d.Set(rtName, routeTable.Name)

	d.Set(rtHref, routeTable.Href)
	d.Set(rtLifecycleState, routeTable.LifecycleState)
	d.Set(rtCreateAt, routeTable.CreatedAt.String())
	d.Set(rtResourceType, routeTable.ResourceType)
	d.Set(rtRouteDirectLinkIngress, routeTable.RouteDirectLinkIngress)
	d.Set(rtRouteTransitGatewayIngress, routeTable.RouteTransitGatewayIngress)
	d.Set(rtRouteVPCZoneIngress, routeTable.RouteVPCZoneIngress)
	d.Set(rtIsDefault, routeTable.IsDefault)
	if routeTable.Routes != nil {
		routes := make([]map[string]interface{}, 0)
		for _, routesItem := range routeTable.Routes {
			routesMap := make(map[string]interface{})

			getVpcRoutingTableRouteOptions := sess.NewGetVPCRoutingTableRouteOptions(idSet[0], idSet[1], *routesItem.ID)
			route, response, err := sess.GetVPCRoutingTableRoute(getVpcRoutingTableRouteOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					d.SetId("")
					return nil
				}
				return fmt.Errorf("Error Getting VPC Routing table route: %s\n%s", err, response)
			}
			routesMap[rtDestination] = *route.Destination
			routesMap["id"] = *route.ID
			//routesMap[rtAction] =
			routesMap[rtName] = *route.Name
			routesMap[rtZone] = *route.Zone.Name
			if route.NextHop != nil {
				nexthop := route.NextHop.(*vpcv1.RouteNextHop)
				if nexthop.Address != nil {
					routesMap[rNextHop] = *nexthop.Address
				}
				if nexthop.ID != nil {
					routesMap[rNextHop] = *nexthop.ID
				}
			}
			routes = append(routes, routesMap)
		}

		routesinput := d.Get("routes").(*schema.Set).List()
		for routeItemId, routeItem := range routesinput {
			route := routeItem.(map[string]interface{})
			routes[routeItemId][rtAction] = route[rtAction].(string)
		}
		d.Set(rtRoutes, routes)

		if err != nil {
			return fmt.Errorf("Error setting routes %s", err)
		}
	}
	subnets := make([]map[string]interface{}, 0)

	for _, s := range routeTable.Subnets {
		subnet := make(map[string]interface{})
		subnet[ID] = *s.ID
		subnet["name"] = *s.Name
		subnets = append(subnets, subnet)
	}

	d.Set(rtSubnets, subnets)

	return nil
}

func resourceIBMIsRoutingTableFlattenRoutes(result []vpcv1.RouteReference) (routes []map[string]interface{}) {
	for _, routesItem := range result {
		routes = append(routes, resourceIBMIsRoutingTableRoutesToMap(routesItem))
	}

	return routes
}

func resourceIBMIsRoutingTableRoutesToMap(routesItem vpcv1.RouteReference) (routesMap map[string]interface{}) {
	routesMap = map[string]interface{}{}

	if routesItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := resourceIBMIsRoutingTableRoutesDeletedToMap(*routesItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		routesMap["deleted"] = deletedList
	}
	if routesItem.Href != nil {
		routesMap["href"] = routesItem.Href
	}
	if routesItem.ID != nil {
		routesMap["id"] = routesItem.ID
	}
	if routesItem.Name != nil {
		routesMap["name"] = routesItem.Name
	}

	return routesMap
}

func resourceIBMIsRoutingTableRoutesDeletedToMap(deletedItem vpcv1.RouteReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func resourceIBMISVPCRoutingTableUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	idSet := strings.Split(d.Id(), "/")
	updateVpcRoutingTableOptions := new(vpcv1.UpdateVPCRoutingTableOptions)
	updateVpcRoutingTableOptions.VPCID = &idSet[0]
	updateVpcRoutingTableOptions.ID = &idSet[1]
	// Construct an instance of the RoutingTablePatch model
	routingTablePatchModel := new(vpcv1.RoutingTablePatch)

	if d.HasChange(rtName) {
		name := d.Get(rtName).(string)
		routingTablePatchModel.Name = core.StringPtr(name)
	}
	if d.HasChange(rtRouteDirectLinkIngress) {
		routeDirectLinkIngress := d.Get(rtRouteDirectLinkIngress).(bool)
		routingTablePatchModel.RouteDirectLinkIngress = core.BoolPtr(routeDirectLinkIngress)
	}
	if d.HasChange(rtRouteTransitGatewayIngress) {
		routeTransitGatewayIngress := d.Get(rtRouteTransitGatewayIngress).(bool)
		routingTablePatchModel.RouteTransitGatewayIngress = core.BoolPtr(routeTransitGatewayIngress)
	}
	if d.HasChange(rtRouteVPCZoneIngress) {
		routeVPCZoneIngress := d.Get(rtRouteVPCZoneIngress).(bool)
		routingTablePatchModel.RouteVPCZoneIngress = core.BoolPtr(routeVPCZoneIngress)
	}

	if d.HasChange("routes") {
		o, n := d.GetChange("routes")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		for _, route := range ors.List() {
			m := route.(map[string]interface{})

			rid := m["id"].(string)

			deleteVpcRoutingTableRouteOptions := sess.NewDeleteVPCRoutingTableRouteOptions(idSet[0], idSet[1], rid)
			response, err := sess.DeleteVPCRoutingTableRoute(deleteVpcRoutingTableRouteOptions)
			if err != nil {
				if response != nil && response.StatusCode != 404 {
					log.Printf("Error deleting VPC Routing table route : %s", response)

				}
				return err
			}
		}

		// Make sure we save the state of the currently configured rules
		routes := o.(*schema.Set).Intersection(n.(*schema.Set))
		d.Set("routes", routes)

		// Then loop through all the newly configured routes and create them
		for _, route := range nrs.List() {
			m := route.(map[string]interface{})
			destination := m[rtDestination].(string)
			zoneName := m[rtZone].(string)
			zone := &vpcv1.ZoneIdentityByName{
				Name: &zoneName,
			}
			createVpcRoutingTableRouteOptions := sess.NewCreateVPCRoutingTableRouteOptions(idSet[0], idSet[1], destination, zone)

			nexthop := m["next_hop"].(string)
			if net.ParseIP(nexthop) == nil {
				nhConnectionID := &vpcv1.RouteNextHopPrototypeVPNGatewayConnectionIdentity{
					ID: core.StringPtr(nexthop),
				}
				createVpcRoutingTableRouteOptions.NextHop = nhConnectionID
			} else {
				nh := &vpcv1.RouteNextHopPrototypeRouteNextHopIP{
					Address: core.StringPtr(nexthop),
				}
				createVpcRoutingTableRouteOptions.NextHop = nh
			}

			if action, ok := m["action"]; ok {
				routeAction := action.(string)
				createVpcRoutingTableRouteOptions.Action = &routeAction
			}

			if name, ok := m["name"]; ok {
				routeName := name.(string)
				createVpcRoutingTableRouteOptions.Name = &routeName
			}

			_, response, err := sess.CreateVPCRoutingTableRoute(createVpcRoutingTableRouteOptions)
			if err != nil {
				log.Printf("[DEBUG] Create VPC Routing table route err %s\n%s", err, response)
				return err
			}
		}

	}
	routingTablePatchModelAsPatch, asPatchErr := routingTablePatchModel.AsPatch()
	if asPatchErr != nil {
		return fmt.Errorf("Error calling asPatch for RoutingTablePatchModel: %s", asPatchErr)
	}
	updateVpcRoutingTableOptions.RoutingTablePatch = routingTablePatchModelAsPatch
	_, response, err := sess.UpdateVPCRoutingTable(updateVpcRoutingTableOptions)
	if err != nil {
		log.Printf("[DEBUG] Update VPC Routing table err %s\n%s", err, response)
		return err
	}
	return resourceIBMISVPCRoutingTableRead(d, meta)
}

func resourceIBMISVPCRoutingTableDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	idSet := strings.Split(d.Id(), "/")

	deleteTableOptions := sess.NewDeleteVPCRoutingTableOptions(idSet[0], idSet[1])
	response, err := sess.DeleteVPCRoutingTable(deleteTableOptions)
	if err != nil && response.StatusCode != 404 {
		log.Printf("Error deleting VPC Routing table : %s", response)
		return err
	}

	d.SetId("")
	return nil
}

func resourceIBMISVPCRoutingTableExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	idSet := strings.Split(d.Id(), "/")
	getVpcRoutingTableOptions := sess.NewGetVPCRoutingTableOptions(idSet[0], idSet[1])
	_, response, err := sess.GetVPCRoutingTable(getVpcRoutingTableOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return false, nil
		}
		return false, fmt.Errorf("Error Getting VPC Routing table : %s\n%s", err, response)
	}
	return true, nil
}

func resourceIbmRouteTableHash(v interface{}) int {
	var buf bytes.Buffer
	m, castOk := v.(map[string]interface{})
	if !castOk {
		return 0
	}

	if v, ok := m[rtDestination]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m[rtZone]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m[rtName]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m[rtAction]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m[rtNextHop]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}
