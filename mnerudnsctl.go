// dns ctl for mne.ru
// License: MIT
// Version: 0.0.1
// Author: rzrbld (Aleksandr Petruhin) https://github.com/rzrbld razblade@gmail.com
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "net/http"
	_ "net/http/cookiejar"
	_ "net/url"
	"os"
	// "log"
	"github.com/parnurzeal/gorequest"
)

const (
	username string = "myemail@mydomain.ru"
	password string = "myS3cr#tP4$w0rd!"
	hostBase string = "https://cp.mne.ru/"
)

var domainName = ""
var client = gorequest.New()

func getDomain(domainName string) (domainId string) {

	var listURL string = fmt.Sprint(hostBase, "dns_manager.php?json=true")
	_, contents, err := client.Get(listURL).End()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	var parsed map[string]interface{}
	JSONdata := json.Unmarshal([]byte(contents), &parsed)
	_ = JSONdata

	domainsArr, ok := parsed["domains"].([]interface{})
	if !ok {
		panic("error while parcing domains")
	}
	/*dnszonesArr, ok := parsed["dnszones"].([]interface{})
	if !ok {
		panic("error while parcing dnszones")
	}
	*/

	domainId = "all"

	for _, domainObj := range domainsArr {
		domainMap := domainObj.(map[string]interface{})
		if domainName == "all" {
			fmt.Println(domainMap["domain"], " Expires:", domainMap["expirydate"])
		} else {
			if domainMap["domain"].(string) == domainName {
				domainId = domainMap["id"].(string)
			}
		}
	}
	/*
		for _, dnszonesObj := range dnszonesArr {
			dnszonesMap := dnszonesObj.(map[string]interface{})
			if domainName == "all" {
				fmt.Println(dnszonesMap["domain"], " DNS Zone ")
			} else {
				if dnszonesMap["domain"].(string) == domainName {
					domainId = dnszonesMap["id"].(string)
				}
			}
		}
	*/
	return domainId
}

func getDomainInfo(domainName string, rrName string) (domainId string, rrId string) {
	domainId = getDomain(domainName)
	var getURL string = fmt.Sprint(hostBase, "dns_manager.php?type=domain&domain_id=", domainId, "&json=true")
	_, contents, err := client.Get(getURL).End()

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	var parsed map[string]interface{}
	JSONdata := json.Unmarshal([]byte(contents), &parsed)
	_ = JSONdata

	// log.Println(parsed)

	rrsArr, ok := parsed["rrs"].([]interface{})
	if !ok {
		fmt.Println("can't get [", domainName, "] Be sure that is your domain.")
		//panic("error while getting domain")
	}
	rrId = ""
	for _, rrsObj := range rrsArr {
		rrsMap := rrsObj.(map[string]interface{})
		if rrName == "" {
			fmt.Println(rrsMap["name"], " Content:", rrsMap["content"], "type:", rrsMap["type"])
		} else {
			// log.Println(rrsMap["name"].(string),"==",rrName)

			if rrsMap["name"].(string) == rrName {
				// log.Println("true")
				rrId = rrsMap["id"].(string)
			}

		}
	}
	if rrId == "" {
		if rrName != "" {
			fmt.Println("can't get [", domainName, " ", rrName, "].")
		}
	}

	return domainId, rrId
}

func usage() {
	execName := os.Args[0]
	fmt.Println("Usage: ", execName, "{list|add|get|rm}")
	fmt.Println("----")
	fmt.Println("list:  [", execName, "list ] output is list of your domains")
	fmt.Println("add:   [", execName, "add domain_name new_rr_name rr_type rr_ip ] example for (test2.rzrbld.ru):", execName, "add rzrbld.ru test2 A 127.0.0.1")
	fmt.Println("get:   [", execName, "get dommin_name ] example for (rzrbld.ru):", execName, "get rzrbld.ru output: rr for this domain is json format")
	fmt.Println("rm:    [", execName, "rm domain_name rr_name ] example for (test2.rzrbld.ru):", execName, "rm rzrbld.ru test2")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 1 {
		usage()
	}
	action := os.Args[1]

	var loginURL string = fmt.Sprint(hostBase, "dologin.php")

	// LOGIN
	_, _, err := client.Post(loginURL).
		Param("username", username).
		Param("password", password).
		End()

	if err != nil {
		fmt.Printf("error posting: %s", err)
		return
	}

	// ACTIONS
	switch action {
	case "list":
		getDomain("all")
	case "add":
		domainName := os.Args[2]
		newDomainName := os.Args[3]
		newDomainType := os.Args[4]
		newDomainIp := os.Args[5]
		domainId := getDomain(domainName)

		var addDomainURL string = fmt.Sprint(hostBase, "dns_manager.php")

		_, _, err := client.Post(addDomainURL).
			Param("action", "add_rr").
			Param("type", "domain").
			Param("domain_id", domainId).
			Param("name", newDomainName).
			Param("prio", "0").
			Param("content", newDomainIp).
			Param("rr_type", newDomainType).
			End()

		if err != nil {
			fmt.Printf("error posting: %s", err)
			return
		}

		fmt.Println("success")

	case "get":
		domainName := os.Args[2]

		getDomainInfo(domainName, "")

	case "rm":
		domainName := os.Args[2]
		rrName := os.Args[3]
		domainId, rrId := getDomainInfo(domainName, rrName)
		// log.Println("domainId:",domainId,"|rrid:",rrId)

		var getURL string = fmt.Sprint(hostBase, "dns_manager.php?type=domain&domain_id=", domainId, "&action=del_rr&rr=", rrId, "&json=true")
		response, _, err := client.Get(getURL).End()
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		} else {
			fmt.Println("success")
		}
		var parsed map[string]interface{}
		JSONdata := json.Unmarshal(contents, &parsed)
		_ = JSONdata

	default:
		usage()
	}
}
