package pdnsgslb

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"pdnsgslb": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("PDNSGLSB_DNSUPDATE_SERVER"); err == "" {
		t.Fatal("PDNSGLSB_DNSUPDATE_SERVER must be set for acceptance tests")
	}
	if err := os.Getenv("PDNSGLSB_DNSUPDATE_PORT"); err == "" {
		t.Fatal("PDNSGLSB_DNSUPDATE_PORT must be set for acceptance tests")
	}
	if err := os.Getenv("PDNSGLSB_DNSUPDATE_KEYNAME"); err == "" {
		t.Fatal("PDNSGLSB_DNSUPDATE_KEYNAME must be set for acceptance tests")
	}
	if err := os.Getenv("PDNSGLSB_DNSUPDATE_KEYALGORITHM"); err == "" {
		t.Fatal("PDNSGLSB_DNSUPDATE_KEYALGORITHM must be set for acceptance tests")
	}
	if err := os.Getenv("PDNSGLSB_DNSUPDATE_KEYSECRET"); err == "" {
		t.Fatal("PDNSGLSB_DNSUPDATE_KEYSECRET must be set for acceptance tests")
	}
}
