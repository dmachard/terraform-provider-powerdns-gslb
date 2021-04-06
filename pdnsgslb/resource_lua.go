package pdnsgslb

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/miekg/dns"
)

func resourceLua() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLuaCreate,
		ReadContext:   resourceLuaRead,
		UpdateContext: resourceLuaUpdate,
		DeleteContext: resourceLuaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"record": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rrtype": {
							Type:     schema.TypeString,
							Required: true,
						},
						"snippet": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},
		},
	}
}

func resourceLuaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// dns client
	c := m.(*Client)

	// record id
	zone := d.Get("zone").(string)
	if !dns.IsFqdn(zone) {
		return diag.Errorf("Not a fully-qualified DNS name: %s", zone)
	}
	name := d.Get("name").(string)
	recordId := fmt.Sprintf("%s.%s", name, zone)

	rrset := d.Get("record").([]interface{})

	// make dns update operation
	_, err := c.doCreate(recordId, rrset)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(recordId)

	return resourceLuaRead(ctx, d, m)
}

func resourceLuaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// dns client
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// get ressource id
	recordId := d.Id()
	if !dns.IsFqdn(recordId) {
		return diag.Errorf("Not a fully-qualified DNS name: %s", recordId)
	}

	// make dns axfr operation
	labels := dns.SplitDomainName(recordId)
	zone := strings.Join(labels[1:], ".") + "."
	name := labels[0]

	rr_lua, err := c.doTransfer(recordId)
	if err != nil {
		return diag.FromErr(err)
	}

	var records []interface{}
	for _, rr := range rr_lua {
		// decode lua
		rrtype_int, _ := strconv.ParseInt(rr.Rdata[0:4], 16, 64)
		rrtype := dns.TypeToString[uint16(rrtype_int)]
		snippet, _ := hex.DecodeString(rr.Rdata[6:])

		urr := make(map[string]interface{})
		urr["rrtype"] = rrtype
		urr["snippet"] = string(snippet)
		urr["ttl"] = rr.Hdr.Ttl

		records = append(records, urr)
	}

	d.Set("zone", zone)
	d.Set("name", name)
	if err := d.Set("record", records); err != nil {
		return diag.Errorf("error setting records for %s: %s", d.Id(), err)
	}

	return diags
}

func resourceLuaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// dns client
	c := m.(*Client)

	// get ressource id
	recordId := d.Id()

	if d.HasChange("record") {
		records := d.Get("record").([]interface{})
		// make dns update operation
		_, err := c.doUpdate(recordId, records)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceLuaRead(ctx, d, m)
}

func resourceLuaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// dns client and variables
	c := m.(*Client)

	// get ressource id
	recordId := d.Id()

	// warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// make dns delete operation
	_, err := c.doDelete(recordId)
	if err != nil {
		return diag.FromErr(err)
	}

	// added here for explicitness, automatically called assuming delete returns no errors
	d.SetId("")

	return diags
}
