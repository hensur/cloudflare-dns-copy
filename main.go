package main

import (
	"flag"
	"github.com/cloudflare/cloudflare-go"
	"log"
	"os"
	"strings"
)

var sourceDomain string
var targetDomain string
var apiToken string

func main() {
	flag.StringVar(&sourceDomain, "source", "", "domain to copy records from")
	flag.StringVar(&targetDomain, "target", "", "domain to copy records to")
	flag.Parse()

	apiToken = os.Getenv("CF_API_TOKEN")
	if apiToken == "" {
		log.Fatal("CF_API_TOKEN not found")
	}

	log.Printf("copying from %s to %s", sourceDomain, targetDomain)

	// Construct a new API object
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.Fatal(err)
	}

	sz, err := api.ZoneIDByName(sourceDomain)
	if err != nil {
		log.Fatal(err)
	}

	tz, err := api.ZoneIDByName(targetDomain)
	if err != nil {
		log.Fatal(err)
	}

	recs, err := api.DNSRecords(sz, cloudflare.DNSRecord{})
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range recs {
		log.Println()
		log.Println("ORIGINAL: -----------------")
		log.Printf("%#v", r)

		r.ID = ""
		r.Name = strings.Replace(r.Name, sourceDomain, targetDomain, -1)
		r.Content = strings.Replace(r.Content, sourceDomain, targetDomain, -1)
		r.ZoneID = tz
		r.ZoneName = targetDomain

		data, ok := r.Data.(map[string]interface{})
		if ok {
			for k := range data {
				if val, ok := data[k].(string); ok {
					data[k] = strings.Replace(val, sourceDomain, targetDomain, -1)
					log.Println(data[k])
				}
			}
			r.Data = data
		}

		log.Println("MODIFIED: --------------------")
		log.Printf("%#v", r)
		log.Println()

		_, err := api.CreateDNSRecord(tz, r)
		if err != nil {
			log.Fatal(err)
		}
	}
}
