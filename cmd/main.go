package main

import (
	"github.com/spf13/cobra"
	"github.com/ssmallmonster/hwameistor-doctor/pkg"
	"io"
	"log"
	"os"
)

func main() {
	err := Hwameictl.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var Hwameictl = &cobra.Command{
	Use:   "hwameidoc",
	Args:  cobra.ExactArgs(0),
	Short: "Hwameictl is the command-line tool for Hwameistor.",
	Long: "Hwameictl is a tool that can manage all Hwameistor resources and their entire lifecycle.\n" +
		"Complete documentation is available at https://hwameistor.io/",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Root cmd will show help only
		return cmd.Help()
	},
}

var (
	localVolume        string
	localVolumeReplica string
	debug              bool
)

func init() {
	// Hwameictl flags
	Hwameictl.PersistentFlags().StringVar(&localVolumeReplica, "localVolumeReplica name", "", "Specify LocalVolumeReplica Name")

	// Sub commands
	Hwameictl.AddCommand(pkg.VolumeReplica)

	// Disable debug mode
	if debug == false {
		log.SetOutput(io.Discard)
	}
}
