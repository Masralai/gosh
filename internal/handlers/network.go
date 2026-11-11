package handlers

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/tatsushid/go-fastping"
	"github.com/urfave/cli/v3"
)

func Ping() *cli.Command {
	return &cli.Command{
		Name:      "ping",
		Usage:     "Send ICMP echo requests to a host",
		UsageText: "ping <hostname>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: ping <hostname>")
			}
			select {
			case <-ctx.Done():
				return fmt.Errorf("operation cancelled: %w", ctx.Err())
			default:
			}
			host := c.Args().Get(0)
			p := fastping.NewPinger()
			ra, err := net.ResolveIPAddr("ip4:icmp", host)
			if err != nil {
				return fmt.Errorf("failed to resolve host %q: %v", host, err)
			}
			p.AddIPAddr(ra)
			p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
				fmt.Printf("Reply from %s: time=%v\n", addr.String(), rtt)
			}
			p.OnIdle = func() {}
			if err := p.Run(); err != nil {
				return fmt.Errorf("failed to ping %q: %w (ICMP ping may require elevated privileges on some systems; try running in an elevated shell or use the host 'ping' command)", host, friendlyPingError(err))
			}
			return nil
		},
	}
}

func friendlyPingError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, os.ErrPermission) {
		return err
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "permission") || strings.Contains(msg, "operation not permitted") || strings.Contains(msg, "access is denied") || strings.Contains(msg, "socket") {
		return err
	}
	return err
}
