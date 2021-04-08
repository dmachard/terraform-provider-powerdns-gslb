package pdnsgslb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPdnsgslbIfurlup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPdnsgslbIfurlupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPdnsgslbIfurlupConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPdnsgslbIfurlupExists("pdnsgslb_ifurlup.testifurlup"),
				),
			},
		},
	})
}

func testAccCheckPdnsgslbIfurlupDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pdnsgslb_ifurlup" {
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

func testAccCheckPdnsgslbIfurlupExists(n string) resource.TestCheckFunc {
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

const testAccCheckPdnsgslbIfurlupConfig_basic = `
resource "pdnsgslb_ifurlup" "testifurlup" {
	zone = "test.internal."
	name = "ifurlup"
	record {
	  rrtype = "A"
	  ttl = 300
	  url = "https://www.facebook.com/"
	  addresses {
		primary = [ 
		  "10.0.0.210",
		  "10.0.0.211",
		]
		backup = [ 
		]
	  }
	}
}`
