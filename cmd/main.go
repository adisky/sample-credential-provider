/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/sample-credential-provider/provider"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	klog "k8s.io/klog/v2"

	"os"

	credentialproviderapi "k8s.io/kubelet/pkg/apis/credentialprovider/v1alpha1"
)

const (
	apiKind    = "CredentialProviderResponse"
	apiVersion = "credentialprovider.kubelet.k8s.io/v1alpha1"
)

func main() {
	defer klog.Flush()
	rootCmd := &cobra.Command{
		Use:   "sample-credential-provider",
		Short: "sample kubelet credential provider based on docker config",
	}
	credCmd, err := NewGetCredentialsCommand()
	if err != nil {
		klog.Errorf(err.Error())
		os.Exit(1)
	}
	rootCmd.AddCommand(credCmd)
	klog.InitFlags(nil)
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	if err := rootCmd.Execute(); err != nil {
		klog.Errorf(err.Error())
		os.Exit(1)
	}
}

// NewGetCredentialsCommand returns a cobra command that retrieves auth credentials after validating flags.
func NewGetCredentialsCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "get-credentials",
		Short: "Get authentication credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCredentials()
		},
	}
	return cmd, nil
}

func getCredentials() error {
	klog.V(2).Infof("sample kubelet credential provider get-credentials")

	unparsedRequest, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	var authRequest credentialproviderapi.CredentialProviderRequest
	err = json.Unmarshal(unparsedRequest, &authRequest)
	if err != nil {
		return fmt.Errorf("error unmarshaling auth credential request: %w", err)
	}

	// TODO: Add a flag for preferred docker config path
	// setting this path as it is used by node e2e
	provider.SetPreferredDockercfgPath("/var/lib/kubelet")

	dockercfg, err := provider.Provide(authRequest.Image)
	if err != nil {
		return err
	}

	authMap := map[string]credentialproviderapi.AuthConfig{}

	for k, v := range dockercfg {
		authMap[k] = credentialproviderapi.AuthConfig{
			Username: v.Username,
			Password: v.Password,
		}

	}

	response := &credentialproviderapi.CredentialProviderResponse{
		CacheKeyType: credentialproviderapi.RegistryPluginCacheKeyType,
		Auth:         authMap,
	}

	response.TypeMeta.Kind = apiKind
	response.TypeMeta.APIVersion = apiVersion

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		// The error from json.Marshal is intentionally not included so as to not leak credentials into the logs
		return fmt.Errorf("error marshaling credentials")
	}
	// Emit authentication response for kubelet to consume
	fmt.Println(string(jsonResponse))

	klog.V(2).Infof("sample credential provider credentials returned")

	return nil
}
