---
page_title: "pdnsgslb_pickrandom Resource - terraform-provider-powerdns-gslb"
subcategory: ""
description: |-
  
---

# pdnsgslb_pickrandom (Resource)

Creates a [pickrandom](https://doc.powerdns.com/authoritative/lua-records/functions.html#pickrandom) LUA DNS record.

## Example Usage

```terraform
resource "pdnsgslb_pickrandom" "foo" {
  zone = "home.internal."
  name = "test_pickrandom"
  record {
    rrtype = "A"
    ttl = 5
    addresses = [ 
      "127.0.0.1",
      "127.0.0.2",
    ]
  }

   record {
    rrtype = "AAAA"
    ttl = 5
    addresses = [
      "::1",
      "fdb0:ccfe:81b8:6500:dc3d:bfff:feea:aa7c",
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
- **addresses** (List) A list of strings with the possible IP addresses.
- **ttl** (Number) The TTL of the record. Defaults to 0. Optional argument


## Import

Records can be imported using the Record Type and FQDN.

```
$ terraform import pdnsgslb_pickrandom.foo foo.example.com.
```