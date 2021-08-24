package api

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// tvmConfigurator implements VMConfigurator interface
type tvmConfigurator struct {
	osName OSName
}

var tvmOSImageMap map[OSName]*OSImage

func init() {
	tvmOSImageMap = map[OSName]*OSImage{
		Ubuntu1804: &OSImage{
			Publisher: "Canonical",
			Offer:     "0003-com-ubuntu-server-trusted-vm",
			SKU:       "18_04-gen2",
			Version:   "18.04.202004290",
		},
		Windows10: &OSImage{
			Publisher: "MicrosoftWindowsServer",
			Offer:     "windowsserver-gen2preview-preview",
			SKU:       "windows10-tvm",
			Version:   "18363.592.2001092016",
		},
		WindowsServer2016: &OSImage{
			Publisher: "MicrosoftWindowsServer",
			Offer:     "windowsserver-gen2preview-preview",
			SKU:       "windowsserver2016-tvm",
			Version:   "14393.3443.2001090113",
		},
		WindowsServer2019: &OSImage{
			Publisher: "MicrosoftWindowsServer",
			Offer:     "windowsserver-gen2preview-preview",
			SKU:       "windowsserver2019-tvm",
			Version:   "17763.973.2001110547",
		},
		Windows10SecuredCore: &OSImage{
			Publisher: "MicrosoftWindowsServer",
			Offer:     "windowsserver-gen2preview-preview",
			SKU:       "windows10-tvm-sc",
			Version:   "19041.329.2006042020",
		},
		WindowsServer2016SecuredCore: &OSImage{
			Publisher: "MicrosoftWindowsServer",
			Offer:     "windowsserver-gen2preview-preview",
			SKU:       "windowsserver2016-tvm-sc",
			Version:   "14393.3750.2006031549",
		},
		WindowsServer2019SecuredCore: &OSImage{
			Publisher: "MicrosoftWindowsServer",
			Offer:     "windowsserver-gen2preview-preview",
			SKU:       "windowsserver2019-tvm-sc",
			Version:   "19041.329.2006042019",
		},
	}
}

// NewTVMConfigurator creates VMConfigurator for TVM
func NewTVMConfigurator(osName OSName) (VMConfigurator, error) {
	if len(osName) > 0 {
		if _, ok := tvmOSImageMap[osName]; !ok {
			return nil, fmt.Errorf("OSName %q is not supported", osName)
		}
	}
	return &tvmConfigurator{osName: osName}, nil
}

func (h *tvmConfigurator) DefaultVMName() string {
	return "tvm"
}

func (h *tvmConfigurator) OSImage() *OSImage {
	if len(h.osName) == 0 {
		log.Fatal("OSName is not set")
	}
	return tvmOSImageMap[h.osName]
}

func (h *tvmConfigurator) OSImageName() *OSImageName {
	if len(h.osName) == 0 {
		log.Fatal("OSNameName is not set")
	}
	var osImageName OSImageName
	return osImageName
}

// DefaultOsDiskType returns default OS disk type
func (h *tvmConfigurator) DefaultOsDiskType() string {
	return "Premium_LRS"
}

// AllowedOsDiskTypes returns supported OS disk types
func (h *tvmConfigurator) AllowedOsDiskTypes() []string {
	return []string{
		"Premium_LRS",
		"StandardSSD_LRS",
	}
}

// AllowedVMSizes returns supported VM sizes
func (h *tvmConfigurator) AllowedVMSizes() []string {
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

func (h *tvmConfigurator) DefaultVMSize() string {
	return "Standard_DC2s_v2"
}
