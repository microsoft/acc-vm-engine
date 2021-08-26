package api

import (
	log "github.com/sirupsen/logrus"
)

// cvmConfigurator implements VMConfigurator interface
type cvmConfigurator struct{}

// NewCVMConfigurator returns VMConfigurator for CVM
func NewCVMConfigurator() (VMConfigurator, error) {
	return &cvmConfigurator{}, nil
}

func (h *cvmConfigurator) DefaultVMName() string {
	return "cvm"
}

// DefaultLinuxImage returns default Linux OS image
func (h *cvmConfigurator) OSImage() *OSImage {
	log.Fatal("OSName is not set")
	return nil
}

// DefaultLinuxImageName returns default Linux OS image name 
func (h *cvmConfigurator) OSImageName() string {
	log.Info("OSImageName is not set")
	return "Ubuntu 20.04 LTS Gen 2"
}

// DefaultOsDiskType returns default OS disk type
func (h *cvmConfigurator) DefaultOsDiskType() string {
	return "Premium_LRS"
}

// AllowedOsDiskTypes returns supported OS disk types
func (h *cvmConfigurator) AllowedOsDiskTypes() []string {
	return []string{
		"Premium_LRS",
		"StandardSSD_LRS",
	}
}

// AllowedVMSizes returns supported VM sizes
func (h *cvmConfigurator) AllowedVMSizes() []string {
	return []string{
		"Standard_DC2as_v5",
		"Standard_DC4as_v5",
		"Standard_DC8as_v5",
		"Standard_DC2ads_v5",
		"Standard_DC4ads_v5",
		"Standard_DC8ads_v5",
	}
}

func (h *cvmConfigurator) DefaultVMSize() string {
	return "Standard_DC2as_v5"
}
