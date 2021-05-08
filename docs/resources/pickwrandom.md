---
page_title: "powerdns-gslb_pickwrandom Resource - terraform-provider-powerdns-gslb"
subcategory: ""
description: |-
  
---

# powerdns-gslb_pickwrandom (Resource)

Creates a [pickwrandom](https://doc.powerdns.com/authoritative/lua-records/functions.html#pickwrandom) LUA DNS record.

## Example Usage

```terraform
resource "powerdns-gslb_pickwrandom" "foo" {
  zone = "home.internal."
  name = "test_pickwrandom"
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
}
```

## Argument Reference

### Required

- **zone** (String) DNS zone the record belongs to.
- **name** (String)  The name of the record. The zone argument will be appended to this value to create the full record path.
- **record** (List) LUA record set. See below for details

### Record set

- **rrtype** (String) The query type of the record (A, AAAA, ...)
- **ipaddress/weight** (Number) Weight for the associated ip address
- **ipaddress/ip** (String) Ip address 
- **ttl** (Number) The TTL of the record. Defaults to 0. Optional argument


## Import

Records can be imported using the Record Type and FQDN.

```
$ terraform import powerdns-gslb_pickwrandom.foo foo.example.com.
```