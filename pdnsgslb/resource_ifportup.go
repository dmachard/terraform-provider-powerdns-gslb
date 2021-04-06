package pdnsgslb

import (
	"context"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/miekg/dns"
)

func resourceIfPortUp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIfPortUpCreate,
		ReadContext:   resourceIfPortUpRead,
		UpdateContext: resourceIfPortUpUpdate,
		DeleteContext: resourceIfPortUpDelete,
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
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"addresses": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourceIfPortUpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// dns client
	c := m.(*Client)

	// record id
	zone := d.Get("zone").(string)
	if !dns.IsFqdn(zone) {
		return diag.Errorf("Not a fully-qualified DNS name: %s", zone)
	}
	name := d.Get("name").(string)
	recordId := fmt.Sprintf("%s.%s", name, zone)

	// get records and transform to lua snippets
	records := d.Get("record").([]interface{})
	rrset := ifPortUpToLuaSnippet(records)

	// make dns update operation
	_, err := c.doCreate(recordId, rrset)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(recordId)

	return resourceIfPortUpRead(ctx, d, m)
}

func resourceIfPortUpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

		// search pickrandom function in snippet
		re := regexp.MustCompile(`ifportup\((?P<param1>\d+),\s*{(?P<param2>.*)}\)`)
		matches_func := re.FindStringSubmatch(string(snippet))

		// no match, ignore record
		if len(matches_func) == 0 {
			continue
		}

		// get addresses paramters
		param1 := matches_func[re.SubexpIndex("param1")]
		port_int, _ := strconv.Atoi(param1)

		// ok, continue to decode addresses parameters
		re2 := regexp.MustCompile(`(?U)'(?P<ip>.*)'`)
		param2 := matches_func[re.SubexpIndex("param2")]
		matches_opt := re2.FindAllStringSubmatch(param2, -1)

		var addresses []string
		for _, match := range matches_opt {
			addresses = append(addresses, match[re2.SubexpIndex("ip")])
		}

		urr := make(map[string]interface{})
		urr["rrtype"] = rrtype
		urr["addresses"] = addresses
		urr["port"] = port_int
		urr["ttl"] = rr.Hdr.Ttl

		records = append(records, urr)
	}

	if len(records) == 0 {
		return diag.Errorf("No LUA records detected")
	}

	d.Set("zone", zone)
	d.Set("name", name)
	if err := d.Set("record", records); err != nil {
		return diag.Errorf("error setting records for %s: %s", d.Id(), err)
	}

	return diags
}

func resourceIfPortUpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// dns client
	c := m.(*Client)

	// get ressource id
	recordId := d.Id()

	if d.HasChange("record") {
		// get records and transform to lua snippets
		records := d.Get("record").([]interface{})
		rrset := ifPortUpToLuaSnippet(records)

		// make dns update operation
		_, err := c.doUpdate(recordId, rrset)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIfPortUpRead(ctx, d, m)
}

func resourceIfPortUpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func ifPortUpToLuaSnippet(records []interface{}) []interface{} {
	var rrset []interface{}

	for _, rr := range records {
		rec := rr.(map[string]interface{})

		portnum := rec["port"].(int)
		addresses := rec["addresses"].([]interface{})
		addresses_list := make([]string, len(addresses))
		for i, v := range addresses {
			addresses_list[i] = fmt.Sprintf("'%s'", v)
		}

		// https://doc.powerdns.com/authoritative/lua-records/functions.html#ifportup
		snippet_lua := fmt.Sprintf("ifportup(")
		snippet_lua += fmt.Sprintf("%s, ", strconv.Itoa(portnum))
		snippet_lua += "{" + strings.Join(addresses_list, ",") + "}"
		snippet_lua += ")"

		rr_new := map[string]interface{}{}
		rr_new["rrtype"] = rec["rrtype"].(string)
		rr_new["ttl"] = rec["ttl"].(int)
		rr_new["snippet"] = snippet_lua

		rrset = append(rrset, rr_new)
	}
	return rrset
}
