package handlers

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/tatsushid/go-fastping"
	"github.com/urfave/cli/v3"
)

func Ping() *cli.Command {
	return &cli.Command{
		Name:     "ping",
		Usage:    "Send ICMP echo requests to a host",
		UsageText: "cli ping <hostname>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: ping <hostname>")
			}
			p := fastping.NewPinger()
			ra, err := net.ResolveIPAddr("ip4:icmp", c.Args().Get(0))
			if err != nil {
				return fmt.Errorf("failed to resolve IP address: %v", err)
			}
			p.AddIPAddr(ra)
			p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
				fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
			}
			p.OnIdle = func() {
				fmt.Println("finish")
			}
			if err := p.Run(); err != nil {
				return fmt.Errorf("failed to ping: %v", err)
			}
			return nil
		},
	}
}