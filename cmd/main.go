package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/component-base/cli"

	"github.com/stlaz/psachecker/pkg/checker"
	"github.com/stlaz/psachecker/pkg/clusterinspect"
)

func main() {
	flags := pflag.NewFlagSet("psachecker", pflag.ExitOnError)
	pflag.CommandLine = flags

	validationCmd := newCmd()
	os.Exit(cli.Run(validationCmd))

}

func newCmd() *cobra.Command {
	o := checker.NewPSACheckerOptions()

	cmd := &cobra.Command{
		Use:          "psachecker resourceType resourceName [flags]",
		Short:        "get the least privileged PodSecurity level for your workload/namespace to keep current workloads running successfully",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			o.Complete(args)
			errs := o.Validate()
			if len(errs) > 0 {
				return fmt.Errorf("there were errors while setting up the command: %v", errs)
			}

			nsAggregatedResults, err := o.Run(context.Background())
			if err != nil {
				return err
			}

			for _, ns := range nsAggregatedResults.Keys() {
				fmt.Fprintf(c.OutOrStdout(), "%s: %s\n", ns, nsAggregatedResults.Get(ns))
			}
			return nil
		},
	}

	o.AddFlags(cmd)

	cmd.AddCommand(clusterinspect.NewClusterInspectCommand(o.ClientConfigOptions))
	return cmd
}
