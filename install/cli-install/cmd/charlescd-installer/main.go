package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/manifoldco/promptui"
)

var selectTemplates = &promptui.SelectTemplates{
	Label:    "{{ . }}?",
	Active:   "\U000002C3 {{ . | cyan }}",
	Inactive: "  {{ . | white }}",
	Selected: "\U000002C3 {{ . | cyan | white }}",
}

var promptTemplates = &promptui.PromptTemplates{
	Prompt:  "{{ . }} ",
	Valid:   "{{ . | green }} ",
	Invalid: "{{ . | red }} ",
	Success: "{{ . | bold }} ",
}

func validateHost(hostOrIP string) error {
	hostRegEx, _ := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	ipRegex, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if !hostRegEx.MatchString(hostOrIP) && !ipRegex.MatchString(hostOrIP) {
		return errors.New("Not a valid hostname/ip")
	}
	return nil
}

func validatePort(portNumberS string) error {
	portNumber, err := strconv.Atoi(portNumberS)
	if err != nil {
		return err
	}
	if portNumber < 0 || portNumber > 65535 {
		return errors.New("Not a port number")
	}
	return nil
}

func validateResourceName(resourceName string) error {
	resourceNameRegEx, _ := regexp.Compile("^[a-zA-Z0-9_-]*$")
	if !resourceNameRegEx.MatchString(resourceName) {
		return errors.New("Not a valid resource name")
	}
	return nil
}

func database() map[string]string {
	databaseConfig := map[string]string{}
	prompt := promptui.Select{
		Label:     "Can we include a Postgres database of your installation?",
		Items:     []string{"Yes", "No, I want to use my own database"},
		Templates: selectTemplates,
	}

	i, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	switch i {
	case 1:
		subPrompt := promptui.Prompt{
			Label:     "Please enter database hostname or ip address:",
			Templates: promptTemplates,
			Validate:  validateHost,
		}
		hostname, err := subPrompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		subPrompt = promptui.Prompt{
			Label:     "Please enter database port number:",
			Templates: promptTemplates,
			Validate:  validatePort,
		}
		portNumber, err := subPrompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		subPrompt = promptui.Prompt{
			Label:     "Please enter database name:",
			Templates: promptTemplates,
			Validate:  validateResourceName,
		}
		databaseName, err := subPrompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		subPrompt = promptui.Prompt{
			Label:     "Please enter database username:",
			Templates: promptTemplates,
			Validate:  validateResourceName,
		}
		userName, err := subPrompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		subPrompt = promptui.Prompt{
			Label:     "Please enter database password:",
			Templates: promptTemplates,
			Mask:      '*',
		}
		password, err := subPrompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		databaseConfig["address"] = hostname
		databaseConfig["port"] = portNumber
		databaseConfig["database"] = databaseName
		databaseConfig["username"] = userName
		databaseConfig["password"] = password
	}

	return databaseConfig
}

func cluster() map[string]string {
	clusterConfig := map[string]string{}
	prompt := promptui.Select{
		Label:     "We need to know the type of authentication we will use to configure your cluster. What type of cloud provider are you using?",
		Items:     []string{"AWS", "Other"},
		Templates: selectTemplates,
	}

	i, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	switch i {
	case 0:
		subPrompt := promptui.Prompt{
			Label:     "Please enter AWS ID:",
			Templates: promptTemplates,
		}
		awsID, err := subPrompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		subPrompt = promptui.Prompt{
			Label:     "Please enter AWS secret:",
			Templates: promptTemplates,
		}
		awsSecret, err := subPrompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		subPrompt = promptui.Prompt{
			Label:     "Please enter AWS region:",
			Templates: promptTemplates,
		}
		awsRegion, err := subPrompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		subPrompt = promptui.Prompt{
			Label:     "Please enter AWS cluster name:",
			Templates: promptTemplates,
		}
		clusterName, err := subPrompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		clusterConfig["type"] = "AWS"
		clusterConfig["awsID"] = awsID
		clusterConfig["awsSecret"] = awsSecret
		clusterConfig["awsRegion"] = awsRegion
		clusterConfig["clusterName"] = clusterName
	case 1:
		subPrompt := promptui.Prompt{
			Label:     "Please enter the cluster hostname:",
			Templates: promptTemplates,
			Validate:  validateHost,
		}
		genericHostName, err := subPrompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		subPrompt = promptui.Prompt{
			Label:     "Please enter the kubeconfig CA Data:",
			Templates: promptTemplates,
		}
		genericCAData, err := subPrompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		subPrompt = promptui.Prompt{
			Label:     "Please enter the kubeconfig Client Certificate:",
			Templates: promptTemplates,
		}
		genericClientCertificate, err := subPrompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		subPrompt = promptui.Prompt{
			Label:     "Please enter the kubeconfig Client Key:",
			Templates: promptTemplates,
		}
		genericClientKey, err := subPrompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		clusterConfig["type"] = "Generic"
		clusterConfig["gHostName"] = genericHostName
		clusterConfig["gCAData"] = genericCAData
		clusterConfig["gClientCertificate"] = genericClientCertificate
		clusterConfig["gClientKey"] = genericClientKey
	}
	return clusterConfig
}

func main() {

	databaseConfig := database()
	fmt.Println("---------------------------")
	clusterConfig := cluster()
	fmt.Println("---------------------------")
	fmt.Println(clusterConfig)
	fmt.Println(databaseConfig)

	//_, result, err := prompt.Run()

	// app := &cli.App{
	// 	Name:   "charlescd-installer",
	// 	Usage:  "Start CharlesCD installation",
	// 	Action: testing,
	// }
	// app.EnableBashCompletion = true

	// err := app.Run(os.Args)

}
