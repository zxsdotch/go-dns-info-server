package main

import (
	"fmt"
	"log"
	insecureRand "math/rand"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/miekg/dns"
)

/**
 * DNS server which responds to TXT requests with various info:
 * - a counter
 * - some random value
 * - server's unix timestamp
 *
 * The purpose of this server is primarily for integration testing of DNS
 * clients -- specifically to ensure that caching is working as expected and/or
 * that cache busting is working as expected.
 *
 * The server is running at: dns-info.zxs.ch
 * Give it a try with e.g.: `dig TXT dns-info.zxs.ch`
 *
 * This code is mostly based off this gist:
 * https://gist.github.com/walm/0d67b4fb2d5daf3edd4fad3e13b162cb
 */

var counter atomic.Int64

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeTXT:
			log.Printf("Query for %s\n", q.Name)
			v := counter.Add(1)

			r := insecureRand.Int63()
			txts := []string{
				fmt.Sprintf("counter=%d", v),
				fmt.Sprintf("rand=%d", r),
				fmt.Sprintf("time=%d", time.Now().Unix()),
			}
			insecureRand.Shuffle(len(txts), func(i, j int) { txts[i], txts[j] = txts[j], txts[i] })

			for _, txt := range txts {
				rr := new(dns.TXT)
				rr.Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 3600}
				rr.Txt = []string{txt}
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}
	w.WriteMsg(m)
}

func main() {
	// The code uses math/rand which is not cryptographically secure but perfectly
	// fine for our purpose. An alternative would be to only use crypto/rand or
	// have math.rand source its randomness from crypto/rand (see
	// https://yourbasic.org/golang/crypto-rand-int/).
	insecureRand.Seed(time.Now().UnixNano())

	// Respond to all queries ("." matches everything).
	dns.HandleFunc(".", handleDnsRequest)

	// Start the server
	port := 53
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
