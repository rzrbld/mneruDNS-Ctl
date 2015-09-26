#!/bin/bash
#dns ctl for mne.ru
#License: MIT
#Version: 0.0.1
#Author: rzrbld (Aleksandr Petruhin) https://github.com/rzrbld razblade@gmail.com
#Dependencies: bash, curl, jq


#mne.ru credentials you need to change this!
email='myemail@mydomain.ru'
passwd='myS3cr#tP4$w0rd!'


#init 
cookies='cookies.txt'
u_agent='Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1'

#login
curl -s -X POST 'https://cp.mne.ru/dologin.php' --user-agent '$u_agent' --data-urlencode "username=${email}" --data-urlencode "password=${passwd}" --cookie $cookies --cookie-jar $cookies


case "$1" in
list)
	curl -s -X GET 'https://cp.mne.ru/dns_manager.php?json=true' --user-agent '$u_agent' --cookie $cookies --cookie-jar $cookies |  jq -r '.domains' | jq -r '.[].domain'
;;
add)
	DOMAINID=`curl -s -X GET 'https://cp.mne.ru/dns_manager.php?json=true' --user-agent '$u_agent' --cookie $cookies --cookie-jar $cookies |  jq -r '.domains'  | jq -r '.[] | select(.domain == "'$2'") | .id' `;
    if [ "$DOMAINID" = "" ]
		then
			echo 'ERROR: Domain '$2' not found! try list command';
		else
			# add domain_id new_name new_type new_ip
			curl -s -X POST 'https://cp.mne.ru/dns_manager.php' --user-agent '$u_agent' --data-urlencode "action=add_rr" --data-urlencode "domain_id=$DOMAINID"  --data-urlencode "type=domain"  --data-urlencode "name=$3"  --data-urlencode "prio=0"  --data-urlencode "content=$5" --data-urlencode "rr_type=$4"  --cookie $cookies --cookie-jar $cookies 
	fi
;;
get)
   # get domain_name
   DOMAINID=`curl -s -X GET 'https://cp.mne.ru/dns_manager.php?json=true' --user-agent '$u_agent' --cookie $cookies --cookie-jar $cookies |  jq -r '.domains'  | jq -r '.[] | select(.domain == "'$2'") | .id' `;
   curl -s -X GET 'https://cp.mne.ru/dns_manager.php?type=domain&domain_id='$DOMAINID'&json=true' --user-agent '$u_agent' --cookie $cookies --cookie-jar $cookies | jq -r '.rrs'
;;
rm)
   # rm domain_name rr_name
   DOMAINID=`curl -s -X GET 'https://cp.mne.ru/dns_manager.php?json=true' --user-agent '$u_agent' --cookie $cookies --cookie-jar $cookies |  jq -r '.domains'  | jq -r '.[] | select(.domain == "'$2'") | .id' `;
   RRID=`curl -s -X GET 'https://cp.mne.ru/dns_manager.php?type=domain&domain_id='$DOMAINID'&json=true' --user-agent '$u_agent' --cookie $cookies --cookie-jar $cookies | jq -r '.rrs' | jq -r '.[] | select(.name=="'$3'") | .id'`
   curl -s -X GET 'https://cp.mne.ru/dns_manager.php?type=domain&domain_id='$DOMAINID'&action=del_rr&rr='$RRID'&json=true'  --user-agent '$u_agent' --cookie $cookies --cookie-jar $cookies | jq -r '.result'
;;
*)
   echo "Usage: $0 {list|add|get|rm}"
   echo "list:  $0 list output: list of your domains"
   echo "add:   $0 add domain_name new_rr_name rr_type rr_ip. example for (test2.rzrbld.ru): $0 add rzrbld.ru test2 A 127.0.0.1"
   echo "get:   $0 get dommin_name example for (rzrbld.ru): $0 get rzrbld.ru output: rr for this domain is json format" 
   echo "rm:    $0 rm domain_name rr_name example for (test2.rzrbld.ru): $0 rm rzrbld.ru test2"
exit 1
esac
