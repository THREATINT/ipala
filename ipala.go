package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/netip"
	"os"

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

		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() != 0 {
				return errors.New("no arguments supported, please use Stdin")
			}

			return run()
		},
	}

	if err = app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(0xff)
	}
}

func run() error {
	var (
		err      error
		scanner  *bufio.Scanner
		lines    []string
		line     string
		i, j     int
		ip       netip.Addr
		network  netip.Prefix
		ips      []netip.Addr
		networks []netip.Prefix
	)

	scanner = bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		line := scanner.Text()

		if len(line) == 0 {
			break
		}
		lines = append(lines, line)
	}

	if scanner.Err() != nil {
		return err
	}

	ips = make([]netip.Addr, len(lines))
	networks = make([]netip.Prefix, len(lines))

	for _, line = range lines {
		if ip, err = netip.ParseAddr(line); err == nil {
			for i = 0; i < len(ips); i++ {
				if !ips[i].IsValid() {
					ips[i] = ip
					break
				}
			}
		} else if network, err = netip.ParsePrefix(line); err == nil {
			for i = 0; i < len(networks); i++ {
				if !networks[i].IsValid() {
					// we have reached the first entry in the slice that is
					// uninitialised. That meas we have gone through the list
					// without finding a network that matches our network, so
					// we store our network at the current position,
					// effectively adding it to the slice
					networks[i] = network
					break
				}

				if network.Overlaps(networks[i]) {
					// Our network / subnet overlaps with an network / subnet
					if networks[i].Bits() > network.Bits() {
						// The entry at the current position of the slice
						// is smaller than our network / subnet, to we replace the
						// entry in the slice
						networks[i] = network
						break
					}
				}
			}
		}
	}

	for i = 0; i < len(networks); i++ {
		if !networks[i].IsValid() {
			break
		}

		fmt.Println(networks[i].String())
	}

	for i = 0; i < len(ips); i++ {
		if !ips[i].IsValid() {
			break
		}

		for j = 0; j < len(networks); j++ {
			if !networks[j].IsValid() {
				// we have gone through the whole list and have not found
				// a network that contains our ip address, so we use our ip
				fmt.Println(ips[i].String())
				break
			}

			if networks[j].Contains(ips[i]) {
				// we have found a network that contains our ip address,
				// so we drop the ip address
				break
			}
		}
	}

	return nil
}