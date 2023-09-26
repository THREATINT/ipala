package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/netip"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	var (
		app *cli.App
		err error
	)

	app = &cli.App{
		Name:  "ipala",
		Usage: "IP Address List Aggregator, please see github.com/THREATINT/ipala",

		Action: func(ctx *cli.Context) error {
			var (
				filterList string
			)

			if ctx.NArg() > 1 {
				return errors.New("only one argument <filter-list> supported")
			}

			if ctx.NArg() == 1 {
				filterList = ctx.Args().Get(0)
			}

			return run(filterList)
		},
	}

	if err = app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(0xff)
	}
}

func run(filterList string) error {
	var (
		err            error
		readFile       *os.File
		scanner        *bufio.Scanner
		lines          []string
		line           string
		i, j           int
		ok             bool
		ip             netip.Addr
		network        netip.Prefix
		filterIps      []netip.Addr
		filterNetworks []netip.Prefix
		ips            []netip.Addr
		networks       []netip.Prefix
	)

	if filterList != "" {
		if readFile, err = os.Open(filterList); err != nil {
			return err
		}
		defer readFile.Close()
		scanner = bufio.NewScanner(readFile)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			lines = append(lines, strings.TrimSpace(scanner.Text()))
		}
		filterIps = make([]netip.Addr, 0)
		filterNetworks = make([]netip.Prefix, 0)
		for _, line = range lines {
			if ip, err = netip.ParseAddr(line); err == nil {
				filterIps = append(filterIps, ip)
			} else if network, err = netip.ParsePrefix(line); err == nil {
				filterNetworks = append(filterNetworks, network)
			}
		}
		lines = []string{}
	}

	scanner = bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		if line = strings.TrimSpace(scanner.Text()); line == "" {
			break
		}
		lines = append(lines, line)
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	ips = make([]netip.Addr, 0)
	networks = make([]netip.Prefix, 0)
	for _, line = range lines {
		if network, err = netip.ParsePrefix(line); err == nil {
			if network.IsSingleIP() {
				line = strings.TrimRight(line, "/32")
				line = strings.TrimRight(line, "/128")
			}
		}

		ok = true

		if ip, err = netip.ParseAddr(line); err == nil {
			for i = 0; i < len(filterIps); i++ {
				if filterIps[i].Compare(ip) == 0 {
					ok = false
					break
				}
			}
			if ok {
				for i = 0; i < len(filterNetworks); i++ {
					if filterNetworks[i].Contains(ip) {
						ok = false
						break
					}
				}

				if ok {
					for i = 0; i < len(ips); i++ {
						// the same ip address is already in the list -> skip
						if ips[i].Compare(ip) == 0 {
							ok = false
							break
						}
					}

					if ok {
						ips = append(ips, ip)
					}
				}
			}
		} else if network, err = netip.ParsePrefix(line); err == nil {
			for i = 0; i < len(filterNetworks); i++ {
				if filterNetworks[i].Overlaps(network) {
					ok = false
					break
				}
			}
			if ok {
				for i = 0; i < len(networks); i++ {
					if network.Overlaps(networks[i]) {
						// Our network / subnet overlaps with a network / subnet
						if networks[i].Bits() > network.Bits() {
							// The entry at the current position of the slice
							// is smaller than our network / subnet, to we replace the
							// entry in the slice
							networks[i] = network
						}
						ok = false
						break
					}
				}

				if ok {
					networks = append(networks, network)
				}
			}
		}
	}

	for i = 0; i < len(networks); i++ {
		fmt.Println(networks[i])
	}

	for i = 0; i < len(ips); i++ {
		ok = true

		for j = 0; j < len(networks); j++ {
			if networks[j].Contains(ips[i]) {
				ok = false
				break
			}
		}
		if ok {
			fmt.Println(ips[i])
		}
	}

	return nil
}
