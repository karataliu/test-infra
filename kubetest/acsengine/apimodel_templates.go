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
	"fmt"
	"os"
	"text/template"
)

var apimodelTemplates = map[string]string{
	"linuxCcm": linuxCcm,
	"windows":  windows,
}

func writeTemplate(templateName string, input *map[string]string, apimodelPath string) error {
	var templateContent string
	var ok bool
	if templateContent, ok = apimodelTemplates[templateName]; !ok {
		return fmt.Errorf("Unknown template %q, supported templates: %q", templateName, getApimodelTemplateKeys())
	}

	a := template.New("templateContent")
	b, err := a.Parse(templateContent)
	if err != nil {
		return err
	}

	f, err := os.Create(apimodelPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := b.Execute(f, *input); err != nil {
		return err
	}

	return nil
}

func getApimodelTemplateKeys() []string {
	keys := []string{}
	for k := range apimodelTemplates {
		keys = append(keys, k)
	}
	return keys
}

const linuxCcm string = `
{
    "apiVersion": "vlabs",
    "location": "{{.location}}",
    "properties": {
        "orchestratorProfile": {
            "orchestratorType": "Kubernetes",
            "orchestratorRelease": "1.9",
            "kubernetesConfig": {
                "useCloudControllerManager": true,
                "customCcmImage": "gcrio.azureedge.net/google_containers/cloud-controller-manager-amd64:v1.11.0-alpha.0",
                "customHyperkubeImage": "gcrio.azureedge.net/google_containers/hyperkube-amd64:v1.11.0-alpha.0",
                "networkPolicy": "none",
                "apiServerConfig": {
                    "--admission-control": "NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota,AlwaysPullImages"
                }
            }
        },
        "masterProfile": {
            "count": 1,
            "vmSize": "Standard_F2",
            "dnsPrefix": "{{.dnsPrefix}}"
        },
        "agentPoolProfiles": [
            {
                "name": "agentpool1",
                "count": 2,
                "vmSize": "Standard_F2",
                "availabilityProfile": "AvailabilitySet",
                "storageProfile": "ManagedDisks"
            }
        ],
        "linuxProfile": {
            "adminUsername": "k8s-ci",
            "ssh": {
                "publicKeys": [
                    {
                        "keyData": "{{.keyData}}"
                    }
                ]
            }
        },
        "servicePrincipalProfile": {
            "clientID": "{{.clientID}}",
            "secret": "{{.secret}}"
        }
    }
}
`

const windows string = `
{
    "apiVersion": "vlabs",
    "location": "{{.location}}",
    "properties": {
        "orchestratorProfile": {
            "orchestratorType": "Kubernetes",
            "orchestratorRelease": "1.9",
            "kubernetesConfig": {
                "useCloudControllerManager": true,
                "customCcmImage": "gcrio.azureedge.net/google_containers/cloud-controller-manager-amd64:v1.11.0-alpha.0",
                "customHyperkubeImage": "gcrio.azureedge.net/google_containers/hyperkube-amd64:v1.11.0-alpha.0",
                "networkPolicy": "none",
                "apiServerConfig": {
                    "--admission-control": "NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota,AlwaysPullImages"
                }
            }
        },
        "masterProfile": {
            "count": 1,
            "vmSize": "Standard_F2",
            "dnsPrefix": "{{.dnsPrefix}}"
        },
        "agentPoolProfiles": [
            {
                "name": "agentpool1",
                "count": 2,
                "vmSize": "Standard_F2",
                "availabilityProfile": "AvailabilitySet",
                "storageProfile": "ManagedDisks",
                "osType": "Windows"
            }
        ],
        "windowsProfile": {
            "adminUsername": "k8s-ci",
            "adminPassword": "{{.adminPassword}}"
        },
        "linuxProfile": {
            "adminUsername": "k8s-ci",
            "ssh": {
                "publicKeys": [
                    {
                        "keyData": "{{.keyData}}"
                    }
                ]
            }
        },
        "servicePrincipalProfile": {
            "clientID": "{{.clientID}}",
            "secret": "{{.secret}}"
        }
    }
}
`
