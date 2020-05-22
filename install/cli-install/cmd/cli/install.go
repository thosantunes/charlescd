package cli

import (
	"context"
	"fmt"
	"installer/pkg/cloudprovider"
	"installer/pkg/configuration"
	"installer/pkg/deployer"
	"installer/pkg/git"
	"installer/pkg/template"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func newInstallationPrompt() ([]string, error) {
	components := []string{"butler", "circle-matcher", "moove", "octopipe", "ui", "darwin-villager"} //TODO: Meio dinamico para alterar essas lista
	selectedComponents := []string{}
	prompt := &survey.MultiSelect{
		Message: "Select CharlesCD modules to install/update:",
		Options: components,
	}

	survey.AskOne(prompt, &selectedComponents)

	return selectedComponents, nil
}

func getOverrideValues(name string, namespace string, component string, config *configuration.Configuration) map[string]string {
	overrideValues := map[string]string{}
	overrideValues["Name"] = name
	overrideValues["Namespace"] = namespace
	overrideValues["image.tag"] = "latest"

	switch component {
	case "villager":
		overrideValues["envVars.1.value"] = strings.Join([]string{"jdbc:postgresql://", config.Database.Hostname, ":", string(config.Database.Port), "/", config.Database.Name}, "")
		overrideValues["envVars.2.value"] = config.Database.Username
		overrideValues["envVars.3.value"] = config.Database.Password
	case "butler":
		// 	- name: DARWIN_NOTIFICATION_URL
		// 	value: "http://charlescd-butler.charles.svc.cluster.local:3000/notifications"
		//   - name: DARWIN_UNDEPLOYMENT_CALLBACK
		// 	value: "http://charlescd-butler.charles.svc.cluster.local:3000/notifications/undeployment"
		//   - name: DARWIN_DEPLOYMENT_CALLBACK
		// 	value: "http://charlescd-butler.charles.svc.cluster.local:3000/notifications/deployment"
		overrideValues["envVars.1.value"] = strings.Join([]string{"jdbc:postgresql://", config.Database.Hostname, ":", string(config.Database.Port), "/", config.Database.Name}, "")
		overrideValues["envVars.2.value"] = config.Database.Username
		overrideValues["envVars.3.value"] = config.Database.Password
	case "ui":
		// - name: REACT_APP_API_URI
		// value: https://charles-api.continuousplatform.com
	}

	return map[string]string{
		"Name":      name,
		"Namespace": namespace,
		"image.tag": "darwin-v-0-2-1",
	}
}

func install(components []string, config *configuration.Configuration) error {
	gitCmd := git.NewGitManager()
	templateCmd := template.NewTemplateManager()

	git, err := gitCmd.NewGit("GITHUB")
	if err != nil {
		return err
	}

	for _, component := range components {
		filesData, err := git.GetDataFromDefaultFiles(component, os.Getenv("GIT_TOKEN"), os.Getenv("TEMPLATE_REPOSITORY"))
		if err != nil {
			return err
		}

		templateProvider, err := templateCmd.NewTemplate("HELM")
		if err != nil {
			return err
		}

		fmt.Println(component)

		manifests, err := templateProvider.GetManifests(filesData[0], filesData[1], getOverrideValues(component, config.Namespace))
		if err != nil {
			return err
		}

		fmt.Println(manifests)

		// err = executeManifests(manifests, config)
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}

func executeManifests(manifests map[string]interface{}, config *configuration.Configuration) error {
	cloudProviderManager := cloudprovider.NewCloudproviderManager()
	deployerManager := deployer.NewDeployerManager()

	cloudConfig := cloudProviderManager.NewCloudProvider(config.Cluster)
	wg, _ := errgroup.WithContext(context.Background())

	for _, manifest := range manifests {
		func(manifest interface{}) {
			log.Println(manifest)
			wg.Go(func() error {
				resource := &deployer.Resource{
					Action:      "DEPLOY",
					Manifest:    deployer.ToUnstructured(manifest.(map[string]interface{})),
					ForceUpdate: false,
					Config:      cloudConfig,
					Namespace:   config.Namespace,
				}

				deployer, err := deployerManager.NewDeployer(resource)
				if err != nil {
					return err
				}

				err = deployer.Do()
				if err != nil {
					return err
				}

				return nil
			})
		}(manifest)
	}

	if err := wg.Wait(); err != nil {
		return err
	}
	return nil
}

func cmdInstallRun(cmd *cobra.Command, args []string) {
	fileName := args[0]
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
		return
	}

	config, err := configuration.New(file)
	if err != nil {
		log.Fatalln(err)
		return
	}

	components, err := newInstallationPrompt()
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = install(components, config)
	if err != nil {
		log.Fatal(err)
		return
	}
}
