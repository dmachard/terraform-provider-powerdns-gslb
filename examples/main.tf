
terraform {
  required_providers {
    powerdns-gslb = {
      version = "1.2.0"
      source  = "github.com/dmachard/powerdns-gslb"
    }
  }
}

# Configure the DNS Provider
provider "powerdns-gslb" {
    server        = "10.0.0.211"
    port = "5353"
    key_name      = "keytest."
    key_algo      = "hmac-sha256"
    key_secret    = "+0GMb2Hak8WiMETQjOCEpTTdRjgym6xJbU/W+Fv+jGsFi9Y0LxG82rRFxf4rdOKihUASU2BTCyCskSKFo+o6pw=="
}

resource "powerdns-gslb_lua" "res1" {
  zone = "test.internal."
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

resource "powerdns-gslb_pickrandom" "res1" {
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

resource "powerdns-gslb_ifportup" "res3" {
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
resource "powerdns-gslb_ifurlup" "res4" {
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

resource "powerdns-gslb_ifurlup" "res5" {
  zone = "home.internal."
  name = "ifurlup_backupgo"
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