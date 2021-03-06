package pdnsgslb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPdnsgslbPickwrandom_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPdnsgslbPickwrandomDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPdnsgslbPickwrandomConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPdnsgslbPickwrandomExists("powerdns-gslb_pickwrandom.testpickwrandom"),
				),
			},
		},
	})
}

func testAccCheckPdnsgslbPickwrandomDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "powerdns-gslb_pickwrandom" {
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

func testAccCheckPdnsgslbPickwrandomExists(n string) resource.TestCheckFunc {
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

const testAccCheckPdnsgslbPickwrandomConfig_basic = `
resource "powerdns-gslb_pickwrandom" "testpickwrandom" {
	zone = "test.internal."
	name = "testpickwrandom"
	record {
	  rrtype = "A"
	  ttl = 5
	  ipaddress {
        weight = 10
        ip = "192.168.1.1"
	  }
      ipaddress {
	    weight = 100
	    ip = "192.168.1.2"
	  }
	}
}`
