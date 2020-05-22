package configuration

import (
	"installer/pkg/cloudprovider"

	"gopkg.in/yaml.v2"
)

type UseCases interface {
	New(file []byte) (*Configuration, error)
}

type Cluster struct {
}

type Configuration struct {
	Namespace  string
	Cluster    *cloudprovider.Cloudprovider
	SecretName string
	Database   *DatabaseConfiguration
}

type DatabaseConfiguration struct {
	Hostname string
	Port     int
	Name     string
	Username string
	Password string
}

func New(file []byte) (*Configuration, error) {
	configuration := Configuration{}
	err := yaml.Unmarshal(file, &configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}

func (configuration *Configuration) validate() {

}
