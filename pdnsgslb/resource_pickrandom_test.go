package pdnsgslb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPdnsgslbPickrandom_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPdnsgslbPickrandomDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPdnsgslbPickrandomConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPdnsgslbPickrandomExists("powerdns_gslb_pickrandom.testpickrandom"),
				),
			},
		},
	})
}

func testAccCheckPdnsgslbPickrandomDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "powerdns_gslb_pickrandom" {
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

func testAccCheckPdnsgslbPickrandomExists(n string) resource.TestCheckFunc {
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

const testAccCheckPdnsgslbPickrandomConfig_basic = `
resource "powerdns_gslb_pickrandom" "testpickrandom" {
	zone = "test.internal."
	name = "testpickrandom"
	record {
	  rrtype = "A"
	  ttl = 5
	  addresses = [ 
		"127.0.0.1",
		"127.0.0.7",
	  ]
	}
}`
