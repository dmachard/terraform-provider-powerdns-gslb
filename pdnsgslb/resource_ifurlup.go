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

func resourceIfUrlUp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIfUrlUpCreate,
		ReadContext:   resourceIfUrlUpRead,
		UpdateContext: resourceIfUrlUpUpdate,
		DeleteContext: resourceIfUrlUpDelete,
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
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"addresses": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"primary": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"backup": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"ttl": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"stringmatch": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
		},
	}
}

func resourceIfUrlUpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	rrset := ifUrlUpToLuaSnippet(records)

	// make dns update operation
	_, err := c.doCreate(recordId, rrset)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(recordId)

	return resourceIfUrlUpRead(ctx, d, m)
}

func resourceIfUrlUpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		re := regexp.MustCompile(`ifurlup\('(?P<param1>.*)\',\s*{(?P<param2>.*)},\s*{(?P<param3>.*)}\)`)
		matches_func := re.FindStringSubmatch(string(snippet))

		// no match, ignore record
		if len(matches_func) == 0 {
			continue
		}

		// get addresses paramters
		url := matches_func[re.SubexpIndex("param1")]

		// continue to decode addresses parameters
		re2 := regexp.MustCompile(`{(?P<param1>.*)},\s*{(?P<param2>.*)}`)
		param2 := matches_func[re.SubexpIndex("param2")]
		matches_addresses := re2.FindStringSubmatch(param2)

		re3 := regexp.MustCompile(`(?U)'(?P<ip>.*)'`)
		addrs_param1 := matches_addresses[re2.SubexpIndex("param1")]
		matches_addrs_param1 := re3.FindAllStringSubmatch(addrs_param1, -1)

		addrs_param2 := matches_addresses[re2.SubexpIndex("param2")]
		matches_addrs_param2 := re3.FindAllStringSubmatch(addrs_param2, -1)

		var addrs_primary []string
		for _, match := range matches_addrs_param1 {
			addrs_primary = append(addrs_primary, match[re3.SubexpIndex("ip")])
		}

		var addrs_backup []string
		for _, match := range matches_addrs_param2 {
			addrs_backup = append(addrs_backup, match[re3.SubexpIndex("ip")])
		}

		var addresses []interface{}
		map_addrs := make(map[string]interface{})
		map_addrs["primary"] = addrs_primary
		map_addrs["backup"] = addrs_backup
		addresses = append(addresses, map_addrs)

		// continue to decode settings
		re4 := regexp.MustCompile(`stringmatch='(?P<stringmatch>.*)'`)
		param3 := matches_func[re.SubexpIndex("param3")]
		matches_opts := re4.FindStringSubmatch(param3)

		// no match, ignore record
		if len(matches_opts) == 0 {
			continue
		}
		stringmatch := matches_opts[re4.SubexpIndex("stringmatch")]

		urr := make(map[string]interface{})
		urr["rrtype"] = rrtype
		urr["addresses"] = addresses
		urr["stringmatch"] = stringmatch
		urr["url"] = url
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

func resourceIfUrlUpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// dns client
	c := m.(*Client)

	// get ressource id
	recordId := d.Id()

	//return diag.Errorf("%s", d)

	if d.HasChange("record") {
		// get records and transform to lua snippets
		records := d.Get("record").([]interface{})
		rrset := ifUrlUpToLuaSnippet(records)

		// make dns update operation
		_, err := c.doUpdate(recordId, rrset)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIfUrlUpRead(ctx, d, m)
}

func resourceIfUrlUpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func ifUrlUpToLuaSnippet(records []interface{}) []interface{} {
	var rrset []interface{}

	for _, rr := range records {
		rec := rr.(map[string]interface{})

		url := rec["url"].(string)
		stringmatch := rec["stringmatch"].(string)

		addresses := rec["addresses"].([]interface{})[0].(map[string]interface{})

		primary_addrs := addresses["primary"].([]interface{})
		primary_list := make([]string, len(primary_addrs))
		for i, v := range primary_addrs {
			primary_list[i] = fmt.Sprintf("'%s'", v)
		}

		backup_addrs := addresses["backup"].([]interface{})
		backup_list := make([]string, len(backup_addrs))
		for i, v := range backup_addrs {
			backup_list[i] = fmt.Sprintf("'%s'", v)
		}

		// https://doc.powerdns.com/authoritative/lua-records/functions.html#ifportup
		snippet_lua := fmt.Sprintf("ifurlup(")
		snippet_lua += fmt.Sprintf("'%s', ", url)
		snippet_lua += "{{" + strings.Join(primary_list, ",") + "}, {" + strings.Join(backup_list, ",") + "} }"
		snippet_lua += fmt.Sprintf(",{stringmatch='%s'}", stringmatch)
		snippet_lua += ")"

		rr_new := map[string]interface{}{}
		rr_new["rrtype"] = rec["rrtype"].(string)
		rr_new["ttl"] = rec["ttl"].(int)
		rr_new["snippet"] = snippet_lua

		rrset = append(rrset, rr_new)
	}
	return rrset
}
