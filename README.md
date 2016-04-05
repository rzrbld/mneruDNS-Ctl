# mneruDNS-Ctl

mneruDNS-Ctl is a bash script or Golang program, you can use it for control your domain zone, that resolves on mne.ru NS servers. You be able to add\remove\list with A,CNAME and other types, from your console

### Current version:
0.0.2

### Tech

There is two version of this app, one written on bash, other written on [Golang].

## About Golang version

### Configuring
change mne.ru login and password to your credentials
```sh
$ vim mnerudnsctl.go
...
#mne.ru credentials you need to change this!
username string = "myemail@mydomain.ru"
password string = "myS3cr#tP4$w0rd!"
...
```
compile
```sh
$ go version
go version go1.6 windows/amd64
$ go get .
$ go build mnerudnsctl.go
```
use
```sh
$ mnerudnsctl help
```
### Examples
list:
```sh
$ mnerudnsctl list
example.com       Expires: 2016-09-12
```
add:
```sh
$ mnerudnsctl add example.com testrr A 127.0.0.3
success
```
get:
```sh
$ mnerudnsctl get example.com
 Content: ns1.mne.ru. noc.mne.ru 0000000000 11111 2222 3333333 44444 Type: SOA
 Content: ns2.mne.ru Type: NS
 Content: ns1.mne.ru Type: NS
 Content: 192.168.1.101 Type: A
testrr Content: 127.0.0.3 Type: A
```
rm
```sh
$ mnerudnsctl rm example.com testrr
success
```


## About Bash version

mneruDNS-Ctl uses a number of open source projects to work properly:

* [curl] - curl is an open source command line tool and library for transferring data with URL syntax
* [jq] - jq is a lightweight and flexible command-line JSON processor.

### Configuring
change mne.ru login and password to your credentials
```sh
$ vim mnerudnsctl.sh
...
#mne.ru credentials you need to change this!
email='myemail@mydomain.ru'
passwd='myS3cr#tP4$w0rd!'
...
```
## Usage
Keep in mind if you choose to use Golang version, then you need command mnerudnsctl (mnerudnsctl.exe) instead of mnerudnsctl.sh.
During use curl creates a temporary file (by default cookies.txt, you able to change name and file location by editing #init section in ./mnerudnsctl.sh) in the same directory as the ./mnerudnsctl.sh, make sure that the user who runs ./mnerudnsctl.sh can create files in the current directory.
```sh
$ ./mnerudnsctl.sh list
mydomain1.ru
example.com
mydomain3.ru
```
```sh
$ ./mnerudnsctl.sh add example.com test5 A 192.16.0.1
ok
```
```sh
$ ./mnerudnsctl.sh get example.com
[
  {
    "id": "2323xxx",
    "type": "NS",
    "name": "",
    "content": "ns1.mne.ru",
    "prio": "0",
    "domain_id": "23xxxx",
    "ttl": "7200"
  },
  {
    "id": "2313xxx",
    "type": "A",
    "name": "",
    "content": "89.108.88.25",
    "prio": "0",
    "domain_id": "23xxxx",
    "ttl": "7200"
  }
]
```
```sh
$ ./mnerudnsctl.sh rm example.com test5
ok
```
License
----
The MIT License (MIT)
Copyright (c) 2012,2013,2014,2015
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.



   [curl]: <http://curl.haxx.se/>
   [jq]: <https://stedolan.github.io/jq/>
   [Golang]: <https://golang.org>
