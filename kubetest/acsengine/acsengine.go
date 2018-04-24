/*
Copyright 2018 The Kubernetes Authors.

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

package acsengine

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const azureProvider string = "azure"

var (
	acsengineLocation               = flag.String("acsengine-location", "westus2", "acsengine deployment location")
	acsengineApimodelTemplate       = flag.String("acsengine-apimodel-template", "", "acsengine apimodel template name")
	acsengineApimodelTemplateConfig = flag.String("acsengine-apimodel-template-config", "", "acsengine apimodel template config file")
)

// Deployer for acsengine
type Deployer struct {
	workspace     string
	resourceGroup string
}

// NewDeployer creates a deployer
func NewDeployer(provider, cluster string) (*Deployer, error) {
	if provider != azureProvider {
		return nil, fmt.Errorf("--provider must be %q for acsengine deployment, found %q", azureProvider, provider)
	}

	workspace, err := ioutil.TempDir("", "kubetest_acsengine_")
	if err != nil {
		return nil, err
	}

	if cluster == "" {
		return nil, fmt.Errorf("--cluster must be set for acsengine deployment")
	}

	reg := regexp.MustCompile("[\\.|_]")
	rg := strings.ToLower(reg.ReplaceAllString(cluster, ""))
	log.Printf("Using workspace: %q", workspace)
	log.Printf("Using resource group: %s", rg)

	return &Deployer{
		workspace:     workspace,
		resourceGroup: rg,
	}, nil
}

// Up setups cluster
func (d *Deployer) Up() error {
	log.Println("Up")

	if *acsengineApimodelTemplate == "" {
		return fmt.Errorf("--acsengine-apimodel-template must be set")
	}

	if *acsengineApimodelTemplateConfig == "" {
		return fmt.Errorf("--acsengine-apimodel-template-config must be set")
	}

	apimodelPath := filepath.Join(d.workspace, "apimodel.json")

	input := make(map[string]string)

	configFile, err := os.Open(*acsengineApimodelTemplateConfig)
	if err != nil {
		return err
	}
	defer configFile.Close()

	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		line := scanner.Text()
		index := strings.Index(line, "=")
		if index > 0 {
			key := line[0:index]
			value := line[index+1:]
			input[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	input["location"] = *acsengineLocation
	input["dnsPrefix"] = d.resourceGroup

	if err := writeTemplate(*acsengineApimodelTemplate, &input, apimodelPath); err != nil {
		return nil
	}

	log.Println("Run acs-engine")
	cmd := exec.Command("acs-engine", "deploy",
		"--location", *acsengineLocation,
		"--subscription-id", input["subscription_id"],
		"--auth-method", "client_secret",
		"--client-id", input["clientID"],
		"--client-secret", input["secret"],
		"-f",
		apimodelPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return fmt.Errorf("j")
	return nil
}

// IsUp checks cluster is up
func (d *Deployer) IsUp() error {
	log.Println("IsUp")
	return nil
}

// Down tears down cluster
func (d *Deployer) Down() error {
	log.Println("Down")
	return nil
}

// DumpClusterLogs dumps clusterlog
func (d *Deployer) DumpClusterLogs(localPath, gcsPath string) error {
	log.Println("DumpClusterLogs")
	return fmt.Errorf("Unimplemented")
}

// TestSetup prepares test
func (d *Deployer) TestSetup() error {
	log.Println("TestSetup")

	if err := os.Unsetenv("KUBERNETES_PROVIDER"); err != nil {
		return err
	}
	if err := os.Setenv("KUBERNETES_CONFORMANCE_TEST", "yes"); err != nil {
		return err
	}
	if err := os.Setenv("KUBERNETES_CONFORMANCE_PROVIDER", azureProvider); err != nil {
		return err
	}
	if err := os.Setenv("CLOUD_CONFIG", "/dev/null"); err != nil {
		return err
	}

	return nil
}

// GetClusterCreated gets createtime
func (d *Deployer) GetClusterCreated(_ string) (time.Time, error) {
	return time.Time{}, fmt.Errorf("Unimplemented")
}
