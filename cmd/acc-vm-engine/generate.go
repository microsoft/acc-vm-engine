package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/microsoft/acc-vm-engine/pkg/api"
	"github.com/microsoft/acc-vm-engine/pkg/engine"
	"github.com/pkg/errors"
)

type generateCmd struct {
	configFile string
	outputDir  string
	sshPubKeys []string
	Disable_SNP bool
	DiskRP_ver string

	// derived
	vm *api.APIModel
}

func NewGenerateCmd() *generateCmd {
	return &generateCmd{
		sshPubKeys: []string{},
	}
}

func (h *generateCmd) Run() error {
	if err := h.validate(); err != nil {
		return errors.Wrap(err, "failed to validate 'generate'")
	}

	if err := h.loadAPIModel(); err != nil {
		return errors.Wrap(err, "failed to load API model in 'generate'")
	}
	if(h.DiskRP_ver == "1") || (h.DiskRP_ver == "2") {
		h.vm.Properties.DiskRPversion = h.DiskRP_ver
	} else {
		h.vm.Properties.DiskRPversion = "2"
	}
	return h.run()
}

func (h *generateCmd) validate() error {
	if _, err := os.Stat(h.configFile); os.IsNotExist(err) {
		return err
	}
	for i, keyPath := range h.sshPubKeys {
		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			return err
		}
		b, err := ioutil.ReadFile(keyPath)
		if err != nil {
			return err
		}
		h.sshPubKeys[i] = strings.TrimSpace(string(b))
	}
	return nil
}

func (h *generateCmd) loadAPIModel() error {
	var err error

	apiloader := &api.Apiloader{}
	if(h.Disable_SNP) {
		apiloader.DisableSNP = true
	}

	h.vm, err = apiloader.LoadVMFromFile(h.configFile, true, false, h.sshPubKeys)
	if err != nil {
		return errors.Wrap(err, "failed to parse config file")
	}

	if h.outputDir == "" {
		h.outputDir = "_output"
	}
	return nil
}

func (h *generateCmd) run() error {
	fmt.Printf("Generating assets into %s...\n", h.outputDir)

	templateGenerator, err := engine.InitializeTemplateGenerator(h.vm.VMConfigurator)
	if err != nil {
		return err
	}

	template, parameters, err := templateGenerator.GenerateTemplate(h.vm, api.DefaultGeneratorCode)
	if err != nil {
		log.Fatalf("error generating template %s: %s", h.configFile, err.Error())
		os.Exit(1)
	}
	/*
		if !gc.noPrettyPrint {
			if template, err = transform.PrettyPrintArmTemplate(template); err != nil {
				log.Fatalf("error pretty printing template: %s \n", err.Error())
			}
			if parameters, err = transform.BuildAzureParametersFile(parameters); err != nil {
				log.Fatalf("error pretty printing template parameters: %s \n", err.Error())
			}
		}
	*/
	writer := &engine.ArtifactWriter{}
	if err = writer.WriteTLSArtifacts(h.vm, template, parameters, h.outputDir, false); err != nil {
		log.Fatalf("error writing artifacts: %s \n", err.Error())
	}

	return nil
}
