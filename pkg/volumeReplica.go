package pkg

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

var VolumeReplica = &cobra.Command{
	Use:   "volumereplica",
	Args:  cobra.ExactArgs(0),
	Short: "Manage the Hwameistor's volumereplica.",
	Long: "Manage the Hwameistor's volumereplica.Hwameistor provides LVM-based data volumes,\n" +
		"which offer read and write performance comparable to that of native disks.\n" +
		"These data voluvolumereplicames also provide advanced features such as data volume expansion,\n" +
		"migration, high availability, and more.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// root cmd will show help only
		return cmd.Help()
	},
}

func init() {
	// VolumeReplica sub commands
	VolumeReplica.AddCommand(VolumeReplicaGet, VolumeReplicaReset)
}

var VolumeReplicaGet = &cobra.Command{
	Use:     "get {volumeReplicaName}",
	Args:    cobra.ExactArgs(1),
	Short:   "Get the Hwameistor volume replica's detail information.",
	Long:    "Get the Hwameistor volume replica's detail information.",
	Example: "hwameidoc volume get pvc-1187f716-db92-47ac-a5fc-44fd19047a81-fsnjsf",
	RunE:    volumeReplicaGetRunE,
}

func volumeReplicaGetRunE(_ *cobra.Command, args []string) error {
	volumeReplicaName := args[0]

	hwameiCli, err := BuildHwameiStorageClient(filepath.Join(homedir.HomeDir(), ".kube/config"))
	if err != nil {
		return err
	}

	volumeReplica, err := hwameiCli.HwameistorV1alpha1().LocalVolumeReplicas().Get(context.Background(), volumeReplicaName, v1.GetOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("volumeReplica %s, status: %s\n", volumeReplica.Name, volumeReplica.Status.State)
	return nil
}

var VolumeReplicaReset = &cobra.Command{
	Use:     "reset {volumeReplicaName}",
	Args:    cobra.ExactArgs(1),
	Short:   "Reset the Hwameistor volume replica's status as NotReady.",
	Long:    "Reset the Hwameistor volume replica's status as NotReady.",
	Example: "hwameidoc volume reset pvc-1187f716-db92-47ac-a5fc-44fd19047a81-fsnjsf",
	RunE:    volumeReplicaResetRunE,
}

func volumeReplicaResetRunE(_ *cobra.Command, args []string) error {
	volumeReplicaName := args[0]

	hwameiCli, err := BuildHwameiStorageClient(filepath.Join(homedir.HomeDir(), ".kube/config"))
	if err != nil {
		return err
	}

	volumeReplica, err := hwameiCli.HwameistorV1alpha1().LocalVolumeReplicas().Get(context.Background(), volumeReplicaName, v1.GetOptions{})
	if err != nil {
		return err
	}

	if volumeReplica.Status.State == "NotReady" {
		fmt.Printf("volumeReplica %s is already in NotReady state.\n", volumeReplica.Name)
		return nil
	}

	stateBefore := volumeReplica.Status.State
	volumeReplica.Status.State = "NotReady"
	_, err = hwameiCli.HwameistorV1alpha1().LocalVolumeReplicas().UpdateStatus(context.Background(), volumeReplica, v1.UpdateOptions{})
	if err != nil {
		fmt.Printf("failed to update volumeReplica %s status from %s to NotReady: %v\n", volumeReplica.Name, stateBefore, err)
		return err
	}

	fmt.Printf("volumeReplica %s status is updated from %s to NotReady successfully.\n", volumeReplica.Name, stateBefore)
	return nil
}
