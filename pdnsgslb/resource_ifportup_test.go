package pdnsgslb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPdnsgslbIfportup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPdnsgslbIfportupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPdnsgslbIfportupConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPdnsgslbIfportupExists("powerdns-gslb_ifportup.testifportup"),
				),
			},
		},
	})
}

func testAccCheckPdnsgslbIfportupDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "powerdns-gslb_ifportup" {
			continue
		}

		recordId := rs.Primary.ID

		_, err := c.doDelete(recordId)
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckPdnsgslbIfportupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No record id set")
		}

		c := testAccProvider.Meta().(*Client)
		rr_lua, err := c.doTransfer(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(rr_lua) == 0 {
			return fmt.Errorf("Lua rr does not exist")
		}

		return nil
	}
}

const testAccCheckPdnsgslbIfportupConfig_basic = `
resource "powerdns-gslb_ifportup" "testifportup" {
	zone = "test.internal."
	name = "ifportup"
	record {
	  rrtype = "A"
	  ttl = 5
	  port = 443
	  addresses = [ 
		"127.0.0.1",
		"127.0.0.2",
	  ]
	  timeout = 10
	}
}`
