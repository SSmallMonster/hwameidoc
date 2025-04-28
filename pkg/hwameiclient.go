package pkg

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"

	hwacli "github.com/hwameistor/hwameistor/pkg/apis/client/clientset/versioned"
	apisv1alpha1 "github.com/hwameistor/hwameistor/pkg/apis/hwameistor/v1alpha1"
)

var home = os.Getenv("HOME")

func BuildKubeClient(kubeConfigPath string) (*kubernetes.Clientset, client.Client, error) {
	loadingRules := clientcmd.NewDefaultPathOptions().LoadingRules

	loadingRules.ExplicitPath = kubeConfigPath
	if !Exists(kubeConfigPath) {
		return nil, nil, fmt.Errorf("kubeconfig is not exists at %v", kubeConfigPath)
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{
		Timeout: "15",
	})

	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	options := client.Options{
		Scheme: runtime.NewScheme(),
	}

	// Setup Scheme for resources
	if err = scheme.AddToScheme(options.Scheme); err != nil {
		return nil, nil, err
	}

	if err = apisv1alpha1.AddToScheme(options.Scheme); err != nil {
		return nil, nil, err
	}

	kClient, err := client.New(config, options)
	return clientSet, kClient, err
}

func BuildHwameiStorageClient(kubeConfigPath string) (*hwacli.Clientset, error) {
	loadingRules := clientcmd.NewDefaultPathOptions().LoadingRules

	loadingRules.ExplicitPath = kubeConfigPath
	if !Exists(kubeConfigPath) {
		return nil, fmt.Errorf("kubeconfig is not exists at %v", kubeConfigPath)
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{
		Timeout: "15",
	})

	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	hwameiCli, err := hwacli.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return hwameiCli, err
}

func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return true
}
