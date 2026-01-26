package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/wavy-cat/doki/internal/app"
)

func NewRootCommand(runner app.Runner) *cobra.Command {
	var (
		opts  app.Options
		ports portsFlag
	)

	cmd := &cobra.Command{
		Use:           "doki",
		Short:         "A minimalistic and fast port knocker.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Ports = ports
			return runner.Run(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Address, "address", "", "Target IP address")
	cmd.Flags().StringVar(&opts.Domain, "domain", "", "Target domain name")
	cmd.Flags().BoolVarP(&opts.ForceUseIPv4, "ipv4", "4", false, "Force use IPv4 protocol")
	cmd.Flags().BoolVarP(&opts.ForceUseIPv6, "ipv6", "6", false, "Force use IPv6 protocol")
	cmd.Flags().DurationVar(&opts.Timeout, "timeout", 30*time.Millisecond, "Maximum time to establish a connection")
	cmd.Flags().BoolVar(&opts.IgnoreErrors, "ignore-errors", false, "Ignore errors when establishing a connection")
	cmd.Flags().Var(&ports, "ports", "Comma-separated list of ports (0-65535 range)")

	return cmd
}
