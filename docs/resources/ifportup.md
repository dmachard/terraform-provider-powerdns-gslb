---
page_title: "powerdns-gslb_ifportup Resource - terraform-provider-powerdns-gslb"
subcategory: ""
description: |-
  
---

# powerdns-gslb_ifportup (Resource)

Creates a [ifportup](https://doc.powerdns.com/authoritative/lua-records/functions.html#ifportup) LUA DNS record. 

## Example Usage

```terraform
resource "powerdns-gslb_ifportup" "foo" {
  zone = "home.internal."
  name = "test_ifportup"
  record {
    rrtype = "A"
    ttl = 5
    port = 443
    addresses = [ 
      "127.0.0.1",
      "127.0.0.2",
    ]
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
- **port** (Number) The port number to test connections to.
- **ttl** (Number) The TTL of the record. Defaults to 0. Optional argument
- **addresses** (List) A list of strings with the possible IP addresses.
- **timeout** (Number) Maximum time in seconds that you allow the check to take (default 5)

## Import

Records can be imported using the Record Type and FQDN.

```
$ terraform import powerdns-gslb_ifportup.foo foo.example.com.
```