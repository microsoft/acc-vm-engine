package api

import "fmt"

// cvmConfigurator implements VMConfigurator interface
type cvmConfigurator struct {
	osName OSName
}

var cvmOSImageMap map[OSName]*OSImage
var cvmOSTypeMap map[OSName]OSType

func init() {
	cvmOSImageMap = map[OSName]*OSImage{
		Ubuntu1804: &OSImage{
			Publisher: "Canonical",
			Offer:     "0003-com-ubuntu-server-trusted-vm",
			SKU:       "18_04-gen2",
			Version:   "18.04.202004290",
		},
	}
	cvmOSTypeMap = map[OSName]OSType{
		Ubuntu1804: Linux,
	}
}

func NewCVMConfigurator(osName OSName) (VMConfigurator, error) {
	if _, ok := cvmOSImageMap[osName]; !ok {
		return nil, fmt.Errorf("OSName %s is not supported", osName)
	}
	return &cvmConfigurator{osName: osName}, nil
}

func (h *cvmConfigurator) DefaultVMName() string {
	return "cvm"
}

// DefaultLinuxImage returns default Linux OS image
func (h *cvmConfigurator) OSImage() *OSImage {
	return cvmOSImageMap[h.osName]
}

// DefaultWindowsImage returns default Windows OS image
func (h *cvmConfigurator) DefaultWindowsImage() *OSImage {
	return &OSImage{
		Publisher: "MicrosoftWindowsServer",
		Offer:     "confidential-compute-preview",
		SKU:       "acc-windows-server-2016-datacenter",
	}
}

// DefaultOsDiskType returns default OS disk type
func (h *cvmConfigurator) DefaultOsDiskType() string {
	return "Premium_LRS"
}

// AllowedLocations returns supported azure regions
func (h *cvmConfigurator) AllowedLocations() []string {
	return []string{
		"eastus",
		"westeurope",
		"uksouth",
	}
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
		"Standard_DC2s",
		"Standard_DC4s",
		"Standard_DC1s_v2",
		"Standard_DC2s_v2",
		"Standard_DC4s_v2",
		"Standard_DC8_v2",
		"Standard_D2s_v3",
		"Standard_D4s_v3",
		"Standard_D8s_v3",
		"Standard_D16s_v3",
		"Standard_D32s_v3",
		"Standard_D64s_v3",
	}
}

func (h *cvmConfigurator) TemplateFiles() []string {
	return []string{"cvm/base.t", "cvm/outputs.t", "cvm/params.t", "cvm/resources.t", "cvm/vars.t"}
}
