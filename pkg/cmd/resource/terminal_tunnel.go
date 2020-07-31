package resource

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stripe/stripe-cli/pkg/config"
	"github.com/stripe/stripe-cli/pkg/terminal"
	"github.com/stripe/stripe-cli/pkg/validators"
	"github.com/stripe/stripe-cli/pkg/version"
)

// TunnelCmd starts an event listener that forwards PaymentIntent events to a Terminal reader in the local network
type TunnelCmd struct {
	cfg *config.Config
	cmd *cobra.Command
}

// NewTunnelCmd returns a new terminal tunnel command
func NewTunnelCmd(parentCmd *cobra.Command, config *config.Config) {
	tunnelCmd := &TunnelCmd{
		cfg: config,
	}

	tunnelCmd.cmd = &cobra.Command{
		Use:     "tunnel",
		Args:    validators.MaximumNArgs(0),
		Short:   "Forward PaymentIntent events to a Terminal reader in the local network",
		Example: `stripe terminal tunnel --api-key sk_123`,
		RunE:    tunnelCmd.runTunnelCmd,
	}

	parentCmd.AddCommand(tunnelCmd.cmd)
}

func (cc *TunnelCmd) runTunnelCmd(cmd *cobra.Command, args []string) error {
	version.CheckLatestVersion()

	err := terminal.Tunnel(cc.cfg)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}
