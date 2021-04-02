package pdnslua

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

func resourceRecordSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordSetCreate,
		ReadContext:   resourceRecordSetRead,
		UpdateContext: resourceRecordSetUpdate,
		DeleteContext: resourceRecordSetDelete,
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
			"lua": {
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

func resourceRecordSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// dns client
	c := m.(*Client)

	// record id
	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	recordId := fmt.Sprintf("%s.%s", name, zone)

	rrset := d.Get("lua").([]interface{})

	// make dns update operation
	_, err := c.doCreate(recordId, rrset)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(recordId)

	return resourceRecordSetRead(ctx, d, m)
}

func resourceRecordSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	zone := strings.Join(labels[1:], ".")
	name := labels[0]

	rr_lua, err := c.doTransfer(recordId)
	if err != nil {
		return diag.FromErr(err)
	}

	lua_list := make([]interface{}, len(rr_lua), len(rr_lua))

	for i, rr := range rr_lua {
		// decode lua
		rrtype_int, _ := strconv.ParseInt(rr.Rdata[0:4], 16, 64)
		rrtype := dns.TypeToString[uint16(rrtype_int)]
		snippet, _ := hex.DecodeString(rr.Rdata[6:])

		urr := make(map[string]interface{})
		urr["rrtype"] = rrtype
		urr["snippet"] = snippet
		urr["ttl"] = rr.Hdr.Ttl

		lua_list[i] = urr
	}

	d.Set("zone", zone)
	d.Set("name", name)
	d.Set("lua", lua_list)

	return diags
}

func resourceRecordSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// dns client
	c := m.(*Client)

	// get ressource id
	recordId := d.Id()

	if d.HasChange("lua") {
		lua_list := d.Get("lua").([]interface{})
		// make dns update operation
		_, err := c.doUpdate(recordId, lua_list)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceRecordSetRead(ctx, d, m)
}

func resourceRecordSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
