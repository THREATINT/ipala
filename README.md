# ipala
IP Address List Agrregator

## Introduction
This command line utility take a list containing ip addresses and/or networks / subnets
from Stdin and prints out an aggregated list to Stdout.
The resulting (aggregated) list no longer contains ip addresses that are part of a network that 
is on the same list. The same goes for networks / subnets that are contained in a larger network / subnet

Both IPv4 and IPv6 are supported. 

IP networks / subnets need to be in CIDR notation 
(please see [RFC2317](https://www.rfc-editor.org/rfc/rfc2317.html)).

## Building and dependencies
ipala is written in Go (Golang), Go SDK v1.20+ is required.
You can either run `go build` or go for a platform build:
- `make amd64` for x86_64
- `make arm64` for arm64

## Usage
Let's jump into an example:

The file `example.txt` (part of this repository) contains the following lines:
```
10.1.1.1
10.1.1.0/24
10.1.0.0/16
192.168.1.0/24
192.168.2.0/24
192.168.1.1
192.168.2.1
192.168.3.1
```
If you `cat example.txt | ipala` you get the following list:
```
10.1.0.0/16
192.168.1.0/24
192.168.2.0/24
192.168.3.1
```
All three single ip addresses `10.1.1.1`, `192.168.1.1`, and `192.168.2.1` that are part of a network 
that is also listed are gone, so is the network / subnet `10.1.1.1.0/24` that is part of the larger 
network / subnet `10.1.0.0/16`. 

## FAQ
#### Does it have to be a tool of its own? Isn't it easier and more flexible to use regex and a small script to achieve the same results?
Yes and no. Regex (regular expressions) may work in most cases, but things get a bit messy when
it comes to debugging or when you have to deal with edge cases like ip addresses in CIDR notation (`/32` for IPv4, `/128` for IPv6),
#### Where is this tool used?
We use this tool in our pipelines to optimise e.g. our  
[Threat Data Feeds for SOHO](https://www.threatint.com/en/solutions/threat-data-feeds/soho). 
These feeds contain known spambots, forum spammers, ip scanners, etc. that should be blocked 
at the WAN side of a network.
The problem with SOHO (small office home office) network equipment is that it is notoriously low
on CPU power and RAM, so smaller more optimised lists are key for a successful deployment on 
small firewalls and similar devices.

## License
Release under the MIT License. (see LICENSE)

## QA
[![DeepSource](https://app.deepsource.com/gh/THREATINT/ipala.svg/?label=active+issues&show_trend=true&token=rvVy0Ld0yBaWKOZsRVfXiAZW)](https://app.deepsource.com/gh/THREATINT/ipala/?ref=repository-badge)