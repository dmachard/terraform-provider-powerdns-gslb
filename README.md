# Terraform Provider PowerDNS GLSB records

A Terraform provider for PowerDNS server to manage LUA records through DNS updates (RFC2136).
This provider can be to used to have a dynamic behaviour of your PowerDNS server, such as Global Server Load Balancing.

## Using the Provider

```hcl
terraform {
  required_providers {
    pdnsgslb = {
      version = "1.0.0"
      source  = "dmachard/powerdns-gslb"
    }
  }
}

# Configure the DNS Provider
provider "pdnsgslb" {
    server        = "10.0.0.210"
    key_name      = "test."
    key_algo      = "hmac-sha256"
    key_secret    = "SxEKov9vWTM+c7k9G6ho5nKX1cJNN3EH9DaqSe8ClwIJezQTBtHrDn5ThGdC/o9my9n5nND5BOHzE6ybvy0+dw=="
}

resource "pdnsgslb_lua" "foo" {
  zone = "home.internal."
  name = "foo"
  record {
    rrtype = "A"
    ttl = 5
    snippet = "ifportup(8082, {'10.0.0.1', '10.0.0.2'})"
  }
  record {
    rrtype = "TXT"
    ttl = 15
    snippet = "os.date()"
  }
}
```

For detailed usage see [provider's documentation page](https://registry.terraform.io/providers/dmachard/powerdns-gslb/latest/docs)
