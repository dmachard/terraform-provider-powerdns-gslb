# Terraform Provider PowerDNS GLSB records

![pdns-auth 4.4](https://img.shields.io/badge/pdns_auth%204.4-tested-green) ![pdns-auth 4.5](https://img.shields.io/badge/pdns_auth%204.5-tested-green)
![pdns-auth 4.6](https://img.shields.io/badge/pdns_auth%204.6-tested-green)

A Terraform provider for PowerDNS server to manage LUA records through DNS updates (RFC2136).
This provider can be to used to have a dynamic behaviour of your PowerDNS server, such as Global Server Load Balancing.

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) > 0.12
-	[Go](https://golang.org/doc/install) >= 1.18

## Using the Provider

```hcl
terraform {
  required_providers {
    powerdns-gslb = {
      version = "1.3.1"
      source  = "dmachard/powerdns-gslb"
    }
  }
}

# Configure the DNS Provider
provider "powerdns-gslb" {
    server        = "10.0.0.210"
    key_name      = "test."
    key_algo      = "hmac-sha256"
    key_secret    = "SxEKov9vWTM+c7k9G6ho5nKX1cJN.....ND5BOHzE6ybvy0+dw=="
}

# Generic LUA record
resource "powerdns-gslb_lua" "foo" {
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

## PowerDNS tuning

Update your `pdns.conf` configuration file  to enable LUA records and DNS update features.

```
enable-lua-records=yes
dnsupdate=yes
```

Then enable the `TSIG` mechanism, `AXFR` and `DNSUPDATE` on your dns zone `test.internal`

```
pdnsutil create-tsig-key tsigkey hmac-sha256
pdnsutil set-meta test.internal TSIG-ALLOW-DNSUPDATE tsigkey
pdnsutil set-meta test.internal TSIG-ALLOW-AXFR tsigkey
pdnsutil set-meta test.internal ALLOW-DNSUPDATE-FROM 0.0.0.0/0
```
