// dns ctl for mne.ru
// License: MIT
// Version: 0.0.2
// Authors:
// Aleksandr Petruhin https://github.com/rzrbld razblade@gmail.com
// Max Nikolenko https://github.com/mephist

package main

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"os"
)

type DomainsResponse struct {
	Dnszones   interface{} `json:"dnszones"`
	DomainId   int64       `json:"domain_id"`
	DomainName string      `json:"domain_name"`
	Domains    []struct {
		Additionalnotes    string `json:"additionalnotes"`
		CreatedAt          string `json:"created_at"`
		Dnsmanagement      string `json:"dnsmanagement"`
		Domain             string `json:"domain"`
		Donotrenew         string `json:"donotrenew"`
		Emailforwarding    string `json:"emailforwarding"`
		Expirydate         string `json:"expirydate"`
		Firstpaymentamount string `json:"firstpaymentamount"`
		Id                 string `json:"id"`
		Idprotection       string `json:"idprotection"`
		Nextduedate        string `json:"nextduedate"`
		Nextinvoicedate    string `json:"nextinvoicedate"`
		Orderid            string `json:"orderid"`
		Paymentmethod      string `json:"paymentmethod"`
		Promoid            string `json:"promoid"`
		Recurringamount    string `json:"recurringamount"`
		Registrar          string `json:"registrar"`
		Registrationdate   string `json:"registrationdate"`
		Registrationperiod string `json:"registrationperiod"`
		Reminders          string `json:"reminders"`
		Status             string `json:"status"`
		Subscriptionid     string `json:"subscriptionid"`
		Synced             string `json:"synced"`
		Type               string `json:"type"`
		UpdatedAt          string `json:"updated_at"`
		Userid             string `json:"userid"`
	} `json:"domains"`
	Rrs []struct {
		Content  string `json:"content"`
		DomainId string `json:"domain_id"`
		Id       string `json:"id"`
		Name     string `json:"name"`
		Prio     string `json:"prio"`
		Ttl      string `json:"ttl"`
		Type     string `json:"type"`
	} `json:"rrs"`
	Type1 string `json:"type1"`
}

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

	response := DomainsResponse{}
	if err := json.Unmarshal([]byte(contents), &response); err != nil {
		panic("error while parcing domains")
	}

	domainId = "all"

	for _, domainObj := range response.Domains {
		if domainName == "all" {
			fmt.Printf("%s\tExpires: %s\n", domainObj.Domain, domainObj.Expirydate)
		} else if domainObj.Domain == domainName {
			domainId = domainObj.Id
		}
	}

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

	response := DomainsResponse{}
	if err := json.Unmarshal([]byte(contents), &response); err != nil {
		panic("can't get [" + domainName + "] Be sure that is your domain.")
	}

	rrId = ""
	for _, rrObj := range response.Rrs {
		if rrName == "" {
			fmt.Printf("%s Content: %s Type: %s\n", rrObj.Name, rrObj.Content, rrObj.Type)
		} else if rrObj.Name == rrName {
			rrId = rrObj.Id
		}
	}
	if rrId == "" && rrName != "" {
		fmt.Printf("can't get [%s %s].", domainName, rrName)
	}

	return domainId, rrId
}

func usage() {
	execName := os.Args[0]
	fmt.Println("Usage: ", execName, "{list|add|get|rm}")
	fmt.Println("----")
	fmt.Println("list:  [", execName, "list ] output is list of your domains")
	fmt.Println("add:   [", execName, "add domain_name new_rr_name rr_type rr_ip ] example for (test2.rzrbld.ru):", execName, "add rzrbld.ru test2 A 127.0.0.1")
	fmt.Println("get:   [", execName, "get domain_name ] example for (rzrbld.ru):", execName, "get rzrbld.ru output: rr for this domain is json format")
	fmt.Println("rm:    [", execName, "rm domain_name rr_name ] example for (test2.rzrbld.ru):", execName, "rm rzrbld.ru test2")
	os.Exit(1)
}

func main() {

	if len(os.Args) < 2 {
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

		var getURL string = fmt.Sprint(hostBase, "dns_manager.php?type=domain&domain_id=", domainId, "&action=del_rr&rr=", rrId, "&json=true")
		_, _, err := client.Get(getURL).End()
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		fmt.Println("success")
		// Unused
		/*
			var parsed map[string]interface{}
			_ := json.Unmarshal([]byte(contents), &parsed)
		*/
	default:
		usage()
	}
}
