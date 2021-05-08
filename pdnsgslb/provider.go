package pdnsgslb

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	defaultPort      = "53"
	defaultTransport = "tcp"
	defaultRetries   = "2"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PDNSGLSB_DNSUPDATE_SERVER", nil),
			},
			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PDNSGLSB_DNSUPDATE_PORT", defaultPort),
			},
			"transport": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PDNSGLSB_DNSUPDATE_TRANSPORT", defaultTransport),
			},
			"retries": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PDNSGLSB_DNSUPDATE_RETRIES", defaultRetries),
			},
			"key_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PDNSGLSB_DNSUPDATE_KEYNAME", nil),
			},
			"key_algo": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PDNSGLSB_DNSUPDATE_KEYALGORITHM", nil),
			},
			"key_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PDNSGLSB_DNSUPDATE_KEYSECRET", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"powerdns-gslb_lua":         resourceLua(),
			"powerdns-gslb_pickrandom":  resourcePickRandom(),
			"powerdns-gslb_pickwrandom": resourcePickWrandom(),
			"powerdns-gslb_ifportup":    resourceIfPortUp(),
			"powerdns-gslb_ifurlup":     resourceIfUrlUp(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	server := data.Get("server").(string)
	transport := data.Get("transport").(string)
	keyname := data.Get("key_name").(string)
	keyalgo := data.Get("key_algo").(string)
	keysecret := data.Get("key_secret").(string)

	// convert string port to int
	port := data.Get("port").(string)
	port_int, err := strconv.Atoi(port)
	if err != nil {
		return nil, diag.Errorf("invalid port: %s", port)
	}

	// convert string retries value to int
	retries := data.Get("retries").(string)
	retries_int, err := strconv.Atoi(retries)
	if err != nil {
		return nil, diag.Errorf("invalid retries: %s", retries)
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := NewClient(server, port_int, transport, keyname, keysecret, keyalgo, retries_int)
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
