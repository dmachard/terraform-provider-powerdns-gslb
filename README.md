# Terraform Provider PowerDNS LUA records

A Terraform provider for PowerDNS server to manage LUA records through DNS updates (RFC2136).
This provider can be to used to have a dynamic behaviour of your PowerDNS server, such as Global Server Load Balancing.

## Using the Provider

```hcl
terraform {
  required_providers {
    pdnslua = {
      version = "0.0.1"
      source  = "github.com/dmachard/pdnslua"
    }
  }
}

# Configure the DNS Provider
provider "pdnslua" {
    server        = "10.0.0.210"
    port          = 53
    transport     = "tcp"
    key_name      = "test."
    key_algo      = "hmac-sha256"
    key_secret    = "SxEKov9vWTM+c7k9G6ho5nKX1cJNN3EH9DaqSe8ClwIJezQTBtHrDn5ThGdC/o9my9n5nND5BOHzE6ybvy0+dw=="
}
```

For detailed usage see [provider's documentation page](https://www.terraform.io/docs/providers/pdnslua/index.html)
