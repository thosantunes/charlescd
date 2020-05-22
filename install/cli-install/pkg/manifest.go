package pkg

import (
	"path"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type ManifestGenerator struct {
}

func NewManifestGen() *ManifestGenerator {
	return &ManifestGenerator{}
}

func (manifestGen *ManifestGenerator) Generate(chartDirName string) (map[string]string, error) {

	newChart, err := chartutil.Load(chartDirName)
	Check(err)

	values, err := chartutil.ReadValuesFile(path.Join(chartDirName, "values.yaml"))
	Check(err)

	yaml, err := values.YAML()
	Check(err)

	config := &chart.Config{
		Raw: yaml,
	}

	renderedValues, err := chartutil.ToRenderValuesCaps(newChart, config, chartutil.ReleaseOptions{}, nil)
	Check(err)

	templateEngine := engine.New()

	templateRender, err := templateEngine.Render(newChart, renderedValues)
	Check(err)

	return templateRender, nil
}
