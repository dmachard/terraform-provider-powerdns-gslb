---
page_title: "powerdns_gslb_ifurlup Resource - terraform-provider-powerdns_gslb"
subcategory: ""
description: |-
  
---

# powerdns-gslb_ifurlup (Resource)

Creates a [ifurlup](https://doc.powerdns.com/authoritative/lua-records/functions.html#ifurlup) LUA DNS record. 

## Example Usage

```terraform
resource "powerdns_gslb_ifurlup" "foo" {
  zone = "home.internal."
  name = "test_ifurlup2"
  record {
    rrtype = "A"
    ttl = 5
    url = "http://helloworld:8080/"
    addresses {
      primary = [ 
        "10.0.0.211",
      ]
      backup = [
        "10.0.0.210",
      ]
    }
    stringmatch="Google3"
  }
}
```

## Argument Reference

### Required

- **zone** (String) DNS zone the record belongs to.
- **name** (String)  The name of the record. The zone argument will be appended to this value to create the full record path.
- **record** (List) LUA record set. See below for details

### Record set

- **rrtype** (String) The query type of the record (A, AAAA, ...)
- **ttl** (Number) The TTL of the record. Defaults to 0. Optional argument
- **url** (String) The url to check.
- **addresses/primary** (List) First set of addresses to check, if an IP address from the first set is available, it will be returned. 
- **addresses/backup** (List) Second set of addresses to check when no addresses work in the first set.
- **stringmatch** (String) Check url for this string, only declare ‘up’ if found. Optional argument.


## Import

Records can be imported using the Record Type and FQDN.

```
$ terraform import powerdns-gslb_ifurlup.foo foo.example.com.
```