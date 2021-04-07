
terraform {
  required_providers {
    pdnsgslb = {
      version = "0.0.1"
      source  = "github.com/dmachard/powerdns-gslb"
    }
  }
}

# Configure the DNS Provider
provider "pdnsgslb" {
    server        = "10.0.0.210"
    key_name      = "test."
    key_algo      = "hmac-sha256"
    key_secret    = "SxEKov9vWTM+c7k9G6ho5nKX1cJNN3EH9DaqSe8ClwIJezQTBtHrDn5ThGdC/o9my9n5nND5BOHzE6ybvy0+dw=="
}

resource "pdnsgslb_lua" "res1" {
  zone = "home.internal."
  name = "pdnslua"
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

resource "pdnsgslb_pickrandom" "res1" {
  zone = "home.internal."
  name = "pickrandom"
  record {
    rrtype = "A"
    ttl = 5
    addresses = [ 
      "127.0.0.1",
      "127.0.0.7",
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

resource "pdnsgslb_ifportup" "res3" {
  zone = "home.internal."
  name = "ifportup"
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

# actif/actif
resource "pdnsgslb_ifurlup" "res4" {
  zone = "home.internal."
  name = "ifurlup_aa"
  record {
    rrtype = "A"
    ttl = 300
    url = "https://www.facebook.com/"
    addresses {
      primary = [ 
        "10.0.0.210",
        "10.0.0.211",
      ]
      backup = [ 
      ]
    }
  }
}

resource "pdnsgslb_ifurlup" "res5" {
  zone = "home.internal."
  name = "ifurlup_backup"
  record {
    rrtype = "A"
    ttl = 5
    url = "https://www.google.com"
    addresses {
      primary = [ 
        "10.0.0.211",
      ]
      backup = [
        "10.0.0.210",
      ]
    }
    stringmatch="Google"
  }
}
