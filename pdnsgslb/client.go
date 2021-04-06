package pdnsgslb

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/bodgit/tsig"
	"github.com/miekg/dns"
)

const (
	TYPE_LUA = 65402
)

type Client struct {
	DNSClient *dns.Client
	SrvAddr   string
	Transport string
	KeyName   string
	KeyAlgo   string
	KeySecret string
}

func NewClient(server string, port int, transport string, keyname string, keysecret string, keyalgo string) (*Client, error) {
	c := Client{
		DNSClient: &dns.Client{},
		SrvAddr:   net.JoinHostPort(server, strconv.Itoa(port)),
		KeyName:   keyname,
		KeySecret: keysecret,
	}

	c.DNSClient.Net = transport
	c.DNSClient.TsigProvider = tsig.HMAC{keyname: keysecret}

	keyalgo, err := convertTsigAlgo(keyalgo)
	if err != nil {
		return nil, err
	}
	c.KeyAlgo = keyalgo

	return &c, nil
}

func (c *Client) doTransfer(record string) ([]*dns.RFC3597, error) {
	labels := dns.SplitDomainName(record)
	zone := dns.Fqdn(strings.Join(labels[1:], "."))

	dnstransfer := new(dns.Transfer)
	dnstransfer.TsigSecret = map[string]string{c.KeyName: c.KeySecret}

	// prepare DNS AXFR operation
	dnsmsg := new(dns.Msg)
	dnsmsg.SetAxfr(zone)
	dnsmsg.SetTsig(c.KeyName, c.KeyAlgo, 300, time.Now().Unix())

	in, err := dnstransfer.In(dnsmsg, c.SrvAddr)
	if err != nil {
		return nil, fmt.Errorf("Error axfr zone: %s", err)
	}

	var lua_records []*dns.RFC3597
	for env := range in {
		if env.Error != nil {
			return nil, fmt.Errorf("Error axfr zone: %s", env.Error)
		}
		for _, rr := range env.RR {
			if rr.Header().Rrtype == 65402 && rr.Header().Name == record {
				unknownRR := new(dns.RFC3597)
				err = unknownRR.ToRFC3597(rr)
				if err != nil {
					return nil, fmt.Errorf("Error to convert to rfc3597 representation: %s", env.Error)
				}
				lua_records = append(lua_records, unknownRR)
			}
		}
	}

	if len(lua_records) == 0 {
		return nil, fmt.Errorf("Error no LUA record retrieved for %s", record)
	}
	return lua_records, nil
}

func (c *Client) doCreate(record string, rrset []interface{}) (*dns.Msg, error) {
	labels := dns.SplitDomainName(record)
	zone := dns.Fqdn(strings.Join(labels[1:], "."))

	// prepare DNS UPDATE operation
	dnsmsg := new(dns.Msg)
	dnsmsg.SetUpdate(zone)

	for _, rr := range rrset {
		lua_rr := rr.(map[string]interface{})

		rrtype_int, err := convertRRType(lua_rr["rrtype"].(string))
		if err != nil {
			return nil, err
		}
		dns_rr := new(dns.RFC3597)
		dns_rr.Hdr.Name = record
		dns_rr.Hdr.Class = dns.ClassINET
		dns_rr.Hdr.Rrtype = 65402
		dns_rr.Hdr.Ttl = uint32(lua_rr["ttl"].(int))

		dns_rr.Rdata = fmt.Sprintf("%04x", rrtype_int)
		dns_rr.Rdata += fmt.Sprintf("%02x", len(lua_rr["snippet"].(string)))
		dns_rr.Rdata += hex.EncodeToString([]byte(lua_rr["snippet"].(string)))

		dnsmsg.Insert([]dns.RR{dns_rr})
	}
	// add tsig key
	dnsmsg.SetTsig(c.KeyName, c.KeyAlgo, 300, time.Now().Unix())

	// send dns query
	r, _, err := c.DNSClient.Exchange(dnsmsg, c.SrvAddr)
	if err != nil {
		return nil, fmt.Errorf("Error creating DNS LUA record: %s", err)
	}
	if r.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("Error creating DNS LUA record: %v (%s)", r.Rcode, dns.RcodeToString[r.Rcode])
	}

	return r, nil
}

func (c *Client) doUpdate(record string, rrset []interface{}) (*dns.Msg, error) {
	labels := dns.SplitDomainName(record)
	zone := dns.Fqdn(strings.Join(labels[1:], "."))

	// prepare DNS UPDATE operation
	dnsmsg := new(dns.Msg)
	dnsmsg.SetUpdate(zone)

	// create remove rr
	rr_remove := new(dns.RFC3597)
	rr_remove.Hdr.Name = record
	rr_remove.Hdr.Class = dns.ClassANY
	rr_remove.Hdr.Rrtype = TYPE_LUA
	rr_remove.Hdr.Ttl = 0
	dnsmsg.RemoveRRset([]dns.RR{rr_remove})

	// re-create lua
	for _, rr := range rrset {
		lua_rr := rr.(map[string]interface{})

		rr_insert := new(dns.RFC3597)
		rr_insert.Hdr.Name = record
		rr_insert.Hdr.Class = dns.ClassINET
		rr_insert.Hdr.Rrtype = TYPE_LUA
		rr_insert.Hdr.Ttl = uint32(lua_rr["ttl"].(int))

		rrtype_int, err := convertRRType(lua_rr["rrtype"].(string))
		if err != nil {
			return nil, err
		}

		rr_insert.Rdata = fmt.Sprintf("%04x", rrtype_int)
		rr_insert.Rdata += fmt.Sprintf("%02x", len(lua_rr["snippet"].(string)))
		rr_insert.Rdata += hex.EncodeToString([]byte(lua_rr["snippet"].(string)))
		dnsmsg.Insert([]dns.RR{rr_insert})
	}

	// add tsig key
	dnsmsg.SetTsig(c.KeyName, c.KeyAlgo, 300, time.Now().Unix())

	// send dns update
	r, _, err := c.DNSClient.Exchange(dnsmsg, c.SrvAddr)
	if err != nil {
		return nil, fmt.Errorf("Error updating DNS LUA record: %s", err)
	}
	if r.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("Error updating DNS LUA record: %v (%s)", r.Rcode, dns.RcodeToString[r.Rcode])
	}

	return r, nil
}

func (c *Client) doDelete(record string) (*dns.Msg, error) {
	labels := dns.SplitDomainName(record)
	zone := dns.Fqdn(strings.Join(labels[1:], "."))

	// prepare DNS UPDATE operation
	dnsmsg := new(dns.Msg)
	dnsmsg.SetUpdate(zone)

	// prepare remove rr
	rr := new(dns.RFC3597)
	rr.Hdr.Name = record
	rr.Hdr.Class = dns.ClassANY
	rr.Hdr.Rrtype = TYPE_LUA
	rr.Hdr.Ttl = 0

	dnsmsg.RemoveRRset([]dns.RR{rr})

	// add tsig key
	dnsmsg.SetTsig(c.KeyName, c.KeyAlgo, 300, time.Now().Unix())

	// send dns delete
	r, _, err := c.DNSClient.Exchange(dnsmsg, c.SrvAddr)
	if err != nil {
		return nil, fmt.Errorf("Error deleting DNS LUA record: %s", err)
	}
	if r.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("Error deleting DNS LUA record: %v (%s)", r.Rcode, dns.RcodeToString[r.Rcode])
	}

	return r, nil
}

func convertTsigAlgo(name string) (string, error) {
	switch name {
	case "hmac-md5":
		return dns.HmacMD5, nil
	case "hmac-sha1":
		return dns.HmacSHA1, nil
	case "hmac-sha256":
		return dns.HmacSHA256, nil
	case "hmac-sha512":
		return dns.HmacSHA512, nil
	default:
		return "", fmt.Errorf("Unknown HMAC algorithm: %s", name)
	}
}

func convertRRType(name string) (uint16, error) {
	switch name {
	case "A":
		return dns.TypeA, nil
	case "AAAA":
		return dns.TypeAAAA, nil
	case "CNAME":
		return dns.TypeCNAME, nil
	case "TXT":
		return dns.TypeTXT, nil
	case "PTR":
		return dns.TypePTR, nil
	case "SOA":
		return dns.TypeSOA, nil
	case "LUA":
		return TYPE_LUA, nil
	default:
		return 0, fmt.Errorf("Unknown rrtype: %s", name)
	}
}
