---
page_title: "powerdns_gslb_lua Resource - terraform-provider-powerdns_gslb"
subcategory: ""
description: |-
  
---

# powerdns-gslb_lua (Resource)

Creates a [generic LUA](https://doc.powerdns.com/authoritative/lua-records/) DNS record.

## Example Usage

```terraform
resource "powerdns_gslb_lua" "svc1" {
  zone = "home.internal."
  name = "test_lua"
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

## Argument Reference

### Required

- **zone** (String) DNS zone the record belongs to.
- **name** (String)  The name of the record. The zone argument will be appended to this value to create the full record path.
- **record** (List) LUA record set. See below for details

### Record set

- **rrtype** (String) The query type of the record (A, AAAA, ...)
- **snippet** (String) Lua snippet. See PowerDNS [documentation](https://doc.powerdns.com/authoritative/lua-records/index.html#examples) for examples
- **ttl** (Number) The TTL of the record. Defaults to 0. Optional argument


## Import

Records can be imported using the Record Type and FQDN.

```
$ terraform import powerdns-gslb_lua.foo foo.example.com.
```