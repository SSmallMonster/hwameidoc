package pkg

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

var Volume = &cobra.Command{
	Use:   "volume",
	Args:  cobra.ExactArgs(0),
	Short: "Manage the Hwameistor's volume.",
	Long:  "Manage the Hwameistor's LocalVolume",
	RunE: func(cmd *cobra.Command, args []string) error {
		// root cmd will show help only
		return cmd.Help()
	},
}

func init() {
	// Volume sub commands
	Volume.AddCommand(VolumeGet, VolumeReset)
}

var VolumeGet = &cobra.Command{
	Use:     "get {volumeName}",
	Args:    cobra.ExactArgs(1),
	Short:   "Get the Hwameistor volume  detail information.",
	Long:    "Get the Hwameistor volume detail information.",
	Example: "hwameidoc volume get pvc-1187f716-db92-47ac-a5fc-44fd19047a81",
	RunE:    volumeReplicaGetRunE,
}

func volumeGetRunE(_ *cobra.Command, args []string) error {
	volumeName := args[0]

	hwameiCli, err := BuildHwameiStorageClient(filepath.Join(homedir.HomeDir(), ".kube/config"))
	if err != nil {
		return err
	}

	volume, err := hwameiCli.HwameistorV1alpha1().LocalVolumes().Get(context.Background(), volumeName, v1.GetOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("volume %s, status: %s\n", volume.Name, volume.Status.State)
	return nil
}

var VolumeReset = &cobra.Command{
	Use:     "reset {volumeName}",
	Args:    cobra.ExactArgs(1),
	Short:   "Reset the Hwameistor volume's status as NotReady.",
	Long:    "Reset the Hwameistor volume's status as NotReady.",
	Example: "hwameidoc volume reset pvc-1187f716-db92-47ac-a5fc-44fd19047a81",
	RunE:    volumeResetRunE,
}

func volumeResetRunE(_ *cobra.Command, args []string) error {
	volumeReplicaName := args[0]

	hwameiCli, err := BuildHwameiStorageClient(filepath.Join(homedir.HomeDir(), ".kube/config"))
	if err != nil {
		return err
	}

	volume, err := hwameiCli.HwameistorV1alpha1().LocalVolumes().Get(context.Background(), volumeReplicaName, v1.GetOptions{})
	if err != nil {
		return err
	}

	if volume.Status.State == "NotReady" {
		fmt.Printf("volume %s is already in NotReady state.\n", volume.Name)
		return nil
	}

	stateBefore := volume.Status.State
	volume.Status.State = "NotReady"
	_, err = hwameiCli.HwameistorV1alpha1().LocalVolumes().UpdateStatus(context.Background(), volume, v1.UpdateOptions{})
	if err != nil {
		fmt.Printf("failed to update volume %s status from %s to NotReady: %v\n", volume.Name, stateBefore, err)
		return err
	}

	fmt.Printf("volume %s status is updated from %s to NotReady successfully.\n", volume.Name, stateBefore)
	return nil
}
