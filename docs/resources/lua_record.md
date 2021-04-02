---
page_title: "pdnslua_record_set Resource - terraform-provider-pdnslua"
subcategory: ""
description: |-
  
---

# pdnslua_record_set (Resource)

Creates a LUA DNS record.

## Example Usage

```terraform
resource "pdnslua_record_set" "sample" {
  zone = "home.local."
  name = "test_lua"

  lua {
    rrtype = "TXT"
    ttl = 0
    snippet = "os.date()"
  }
}
```

## Argument Reference

### Required

- **zone** (String) DNS zone the record belongs to.
- **name** (String)  The name of the record. The zone argument will be appended to this value to create the full record path.
- **lua** (List) LUA record set. See below for details

### Lua record set

- **rrtype** (String) The query type of the record (A, AAAA, ...)
- **snippet** (String) Lua snippet, See PowerDNS [documentation](https://doc.powerdns.com/authoritative/lua-records/index.html#examples) for examples
- **ttl** (Number) The TTL of the record. Defaults to 0. Optional argument


## Import

Records can be imported using the Record Type and FQDN.

```
$ terraform import pdnslua_record_set.foo foo.example.com.
```