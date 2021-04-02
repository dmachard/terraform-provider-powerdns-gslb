# PowerDNS LUA records Provider

A Terraform provider for PowerDNS server to manage LUA records through DNS updates (RFC2136).
This provider can be to used to have a dynamic behaviour of your PowerDNS server, such as Global Server Load Balancing.

## Requirements

The following features must be enabled on the PDNS server
-  TSIG (RFC 2845) which is required authentication
-  DNS zone transfer, which is required for listing of LUA records.

## Example Usage

```terraform
# Configure the PDNSLUA Provider
provider "pdnslua" {
    server        = "10.0.0.210"
    key_name      = "test."
    key_algo      = "hmac-sha256"
    key_secret    = "SxEKov9vWTM+c7k9G6ho5nKX1cJNN3EH9DaqSe8ClwIJezQTBtHrDn5ThGdC/o9my9n5nND5BOHzE6ybvy0+dw=="
}

# Create a LUA DNS record
resource "pdnslua_record_set" "svc1" {
  # ...
}
```

## Argument Reference

### Required

- **key_algo** (String) The algorithm to use for HMAC TSIG authentication.
- **key_name** (String) The name of the TSIG key used to sign the DNS update messages.
- **key_secret** (String) A Base64-encoded string containing the shared secret to be used for TSIG.
- **server** (String) The hostname or IP address of the DNS server to send updates to.

### Optional

- **port** (Number) The target UDP port on the server where updates are sent to. Defaults to 53.
- **transport** (String) Transport to use for DNS queries. Valid values are udp, udp4, udp6, tcp, tcp4, or tcp6. Defaults to tcp.
