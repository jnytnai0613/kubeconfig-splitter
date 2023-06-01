package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func generateConfig(config clientcmdapi.Config) {
	var fileIdx int

	// Iterate over the contexts in the kubeconfig file.
	// For each context, create a new kubeconfig file with the context as the default context.
	for k, v := range config.Contexts {
		clusterName := v.Cluster
		authInfo := v.AuthInfo
		contextName := k
		cmdContext := &clientcmdapi.Context{
			Cluster:  clusterName,
			AuthInfo: authInfo,
		}

		// Get the cluster, authinfo, and endpoint data from the parent kubeconfig file.
		apiEndpoint := string(config.Clusters[clusterName].Server)
		authData := config.Clusters[clusterName].CertificateAuthorityData
		certData := config.AuthInfos[authInfo].ClientCertificateData
		keyData := config.AuthInfos[authInfo].ClientKeyData

		// Create a new kubeconfig file with the cluster, authinfo, and endpoint data.
		cmdCluser := &clientcmdapi.Cluster{
			Server:                   apiEndpoint,
			CertificateAuthorityData: authData,
		}
		cmdAuthInfo := &clientcmdapi.AuthInfo{
			ClientCertificateData: certData,
			ClientKeyData:         keyData,
		}

		clusters := make(map[string]*clientcmdapi.Cluster)
		contexts := make(map[string]*clientcmdapi.Context)
		auths := make(map[string]*clientcmdapi.AuthInfo)

		// Add the cluster, authinfo, and context data to the new kubeconfig file.
		clusters[clusterName] = cmdCluser
		contexts[contextName] = cmdContext
		auths[authInfo] = cmdAuthInfo
		cmdConfig := &clientcmdapi.Config{
			Clusters:  clusters,
			Contexts:  contexts,
			AuthInfos: auths,
		}

		// Write the new kubeconfig file to a temporary location.
		// The kubeconfig file will be named "kubeconfig0", "kubeconfig1", "kubeconfig2", etc.
		if err := clientcmd.WriteToFile(
			*cmdConfig,
			filepath.Join(os.TempDir(), fmt.Sprintf("kubeconfig%d", fileIdx))); err != nil {
			log.Fatalln(err)
		}
		fileIdx++
	}
}

func ReadKubeconfig() (*clientcmdapi.Config, error) {
	// Specify the path of the kubeconfig file to be loaded in clientcmd.ClientConfigLoadingRules.
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.ExplicitPath = "./config"

	kubeconfig, err := loadingRules.Load()
	if err != nil {
		return kubeconfig, err
	}

	return kubeconfig, nil
}

func main() {
	config, err := ReadKubeconfig()
	if err != nil {
		log.Fatalln(err)
	}

	generateConfig(*config)
}
