# PowerDNS GSLB Provider

A Terraform provider for PowerDNS server to manage LUA records through DNS updates (RFC2136).
This provider can be to used to have a dynamic behaviour of your PowerDNS server, such as Global Server Load Balancing.

## Requirements

The following features must be enabled on the PDNS server
- LUA records feature enabled
- DNS UPDATE enabled 
- TSIG (RFC 2845) which is required authentication
- DNS zone transfer, which is required for listing of LUA records.

## Example Usage

```terraform

terraform {
  required_providers {
    pdnsgslb = {
      version = "1.0.0"
      source  = "dmachard/powerdns-gslb"
    }
  }
}

# Configure the provider
provider "pdnsglsb" {
    server        = "10.0.0.210"
    key_name      = "test."
    key_algo      = "hmac-sha256"
    key_secret    = "SxEKov9vWTM+c7k9G6ho5nK.....n5nND5BOHzE6ybvy0+dw=="
}

# Create a LUA DNS record
resource "powerdns-gslb_lua" "foo" {
  # ...
}
```

## Argument Reference

### Required

- **server** (String) The hostname or IP address of the DNS server to send updates to. This can also be specified with `PDNSGLSB_DNSUPDATE_SERVER` environment variable.
- **key_algo** (String) The algorithm to use for HMAC TSIG authentication. This can also be specified with `PDNSGLSB_DNSUPDATE_KEYALGORITHM` environment variable.
- **key_name** (String) The name of the TSIG key used to sign the DNS update messages. This can also be specified with `PDNSGLSB_DNSUPDATE_KEYNAME` environment variable.
- **key_secret** (String) A Base64-encoded string containing the shared secret to be used for TSIG. This can also be specified with `PDNSGLSB_DNSUPDATE_SECRET` environment variable.

### Optional

- **port** (String) The target UDP port on the server where updates are sent to. Defaults to `53`. This can also be specified with `PDNSGLSB_DNSUPDATE_PORT` environment variable.
- **transport** (String) Transport to use for DNS queries. Valid values are udp, udp4, udp6, tcp, tcp4, or tcp6. Defaults to `tcp`. This can also be specified with `PDNSGLSB_DNSUPDATE_TRANSPORT` environment variable.
- **retries** (String) How many times to retry on connection timeout. Defaults to `2`. Optional parameter