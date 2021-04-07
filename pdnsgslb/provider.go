package pdnsgslb

import (
	"context"
	"fmt"
	"os"
	"strconv"

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
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PDNSGLSB_DNSUPDATE_SERVER", nil),
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if envPortStr := os.Getenv("PDNSGLSB_DNSUPDATE_PORT"); envPortStr != "" {
						port, err := strconv.Atoi(envPortStr)
						if err != nil {
							err = fmt.Errorf("invalid PDNSGLSB_DNSUPDATE_PORT environment variable: %s", err)
						}
						return port, err
					}

					return defaultPort, nil
				},
			},
			"transport": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PDNSGLSB_DNSUPDATE_TRANSPORT", defaultTransport),
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
			"pdnsgslb_lua":        resourceLua(),
			"pdnsgslb_pickrandom": resourcePickRandom(),
			"pdnsgslb_ifportup":   resourceIfPortUp(),
			"pdnsgslb_ifurlup":    resourceIfUrlUp(),
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
