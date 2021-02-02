package api

import (
	"fmt"
	"strings"
)

const (
	// DefaultGeneratorCode specifies the source generator of the cluster template.
	DefaultGeneratorCode = "acc-vm-engine"
	// DefaultVnet specifies default vnet address space
	DefaultVnet = "10.1.16.0/24"
	// DefaultSubnet specifies default subnet
	DefaultSubnet = "10.1.16.0/24"
)

const (
	// TVM VM category
	TVM VMCategory = "TVM"
	// CVM VM category
	CVM VMCategory = "CVM"
)

const (
	Linux   OSType = "Linux"
	Windows OSType = "Windows"

	Ubuntu1804                   OSName = "Ubuntu18.04"
	Windows10                    OSName = "Windows10"
	WindowsServer2016            OSName = "WindowsServer2016"
	WindowsServer2019            OSName = "WindowsServer2019"
	Windows10SecuredCore         OSName = "Windows10-SecuredCore"
	WindowsServer2016SecuredCore OSName = "WindowsServer2016-SecuredCore"
	WindowsServer2019SecuredCore OSName = "WindowsServer2019-SecuredCore"
)

func getAllowedValues(vals []string) string {
	strFormat := `"allowedValues": [
        "%s"
      ],
  `
	return fmt.Sprintf(strFormat, strings.Join(vals, "\",\n        \""))
}

func getDefaultValue(def string) string {
	strFormat := `"defaultValue": "%s",
	`
	return fmt.Sprintf(strFormat, def)
}

func getAllowedDefaultValues(vals []string, def string) string {
	return getAllowedValues(vals) + "    " + getDefaultValue(def)
}

// GetAllowedVMSizes returns allowed sizes for VM
func GetAllowedVMSizes(vmconf VMConfigurator) string {
	return getAllowedValues(vmconf.AllowedVMSizes())
}

// GetAllowedVMSizes returns allowed sizes for VM
func GetDefaultVMSize(vmconf VMConfigurator) string {
	return getDefaultValue(vmconf.DefaultVMSize())
}

// GetOsDiskTypes returns allowed and default OS disk types
func GetOsDiskTypes(vmconf VMConfigurator) string {
	return getAllowedDefaultValues(vmconf.AllowedOsDiskTypes(), vmconf.DefaultOsDiskType())
}
