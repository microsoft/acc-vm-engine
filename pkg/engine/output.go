package engine

import (
	"github.com/microsoft/acc-vm-engine/pkg/api"
)

// ArtifactWriter represents the object that writes artifacts
type ArtifactWriter struct {
}

// WriteTLSArtifacts saves TLS certificates and keys to the server filesystem
func (w *ArtifactWriter) WriteTLSArtifacts(vm *api.APIModel, template, parameters, artifactsDir string, parametersOnly bool) error {
	if len(artifactsDir) == 0 {
		artifactsDir = "_output"
	}

	f := &FileSaver{}

	// convert back the API object, and write it
	var b []byte
	var err error
	if !parametersOnly {
		apiloader := &api.Apiloader{}
		b, err = apiloader.SerializeVM(vm)

		if err != nil {
			return err
		}

		if e := f.SaveFile(artifactsDir, "apimodel.json", b); e != nil {
			return e
		}

		if e := f.SaveFileString(artifactsDir, "azuredeploy.json", template); e != nil {
			return e
		}
	}

	if e := f.SaveFileString(artifactsDir, "azuredeploy.parameters.json", parameters); e != nil {
		return e
	}

	return nil
}
