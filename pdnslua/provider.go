package pdnslua

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	defaultPort      = 53
	defaultTransport = "tcp"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  defaultPort,
			},
			"transport": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  defaultTransport,
			},
			"key_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_algo": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_secret": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"pdnslua_record_set": resourceRecordSet(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	server := data.Get("server").(string)
	port := data.Get("port").(int)
	transport := data.Get("transport").(string)
	keyname := data.Get("key_name").(string)
	keyalgo := data.Get("key_algo").(string)
	keysecret := data.Get("key_secret").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := NewClient(server, port, transport, keyname, keysecret, keyalgo)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create PdnsGslb client",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return c, diags
}
