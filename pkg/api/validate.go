package api

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/microsoft/acc-vm-engine/pkg/api/common"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	validate        *validator.Validate
	keyvaultIDRegex *regexp.Regexp
	labelValueRegex *regexp.Regexp
	labelKeyRegex   *regexp.Regexp
)

const (
	labelKeyPrefixMaxLength = 253
	labelValueFormat        = "^([A-Za-z0-9][-A-Za-z0-9_.]{0,61})?[A-Za-z0-9]$"
	labelKeyFormat          = "^(([a-zA-Z0-9-]+[.])*[a-zA-Z0-9-]+[/])?([A-Za-z0-9][-A-Za-z0-9_.]{0,61})?[A-Za-z0-9]$"
)

func init() {
	validate = validator.New()
	keyvaultIDRegex = regexp.MustCompile(`^/subscriptions/\S+/resourceGroups/\S+/providers/Microsoft.KeyVault/vaults/[^/\s]+$`)
	labelValueRegex = regexp.MustCompile(labelValueFormat)
	labelKeyRegex = regexp.MustCompile(labelKeyFormat)
}

// Validate implements APIObject
func (p *Properties) Validate(vmconf VMConfigurator, isUpdate bool) error {
	if e := validate.Struct(p); e != nil {
		return handleValidationErrors(e.(validator.ValidationErrors))
	}
	if e := p.validateVMProfile(vmconf); e != nil {
		return e
	}
	if e := p.validateVnetProfile(); e != nil {
		return e
	}
	if e := p.validateDiagnosticsProfile(); e != nil {
		return e
	}
	return nil
}

func handleValidationErrors(e validator.ValidationErrors) error {
	// Override any version specific validation error message
	// common.HandleValidationErrors if the validation error message is general
	return common.HandleValidationErrors(e)
}

func (p *Properties) validateVMProfile(vmconf VMConfigurator) error {
	var isLinux bool
	vm := p.VMProfile
	if vm == nil {
		return fmt.Errorf("VMProfile is not specified")
	}
	if len(vm.OSType) == 0 {
		return fmt.Errorf("OSType is not specified")
	}
	hasOsImage := (len(vm.OSName) > 0 || vm.OSImage != nil || len(vm.OSImageName) > 0)
	hasOsDisk := (vm.OSDisk != nil)

	if hasOsImage {
		if len(vm.OSName) == 0 {
			if vm.OSImage == nil && vm.OSImageName == nil {
				return fmt.Errorf("Either OSName or OSImage should be specified")
			}
		} else {
			if vm.OSImage != nil {
				return fmt.Errorf("Cannot have OSName and OSImage both specified")
			}
		}
	}

	if hasOsImage && hasOsDisk {
		return fmt.Errorf("OS image and disk are mutually exclusive")
	}
	if !hasOsImage && !hasOsDisk {
		return fmt.Errorf("Neither OS image nor disk are specified")
	}

	switch vm.OSType {
	case Linux:
		isLinux = true
	case Windows:
		isLinux = false
	default:
		return fmt.Errorf("OS type '%s' is not supported", vm.OSType)
	}

	if e := validateOSImage(vm.OSImage); e != nil {
		return e
	}
	if e := validateOSDisk(vm.OSDisk); e != nil {
		return e
	}
	if (vm.SecurityProfile != nil) {
		if (vm.SecurityProfile.SecureBoot != "true") && (vm.SecurityProfile.SecureBoot != "false") && (vm.SecurityProfile.SecureBoot != "none"){
			return fmt.Errorf("Invalid Entry! Only the values \"true\", \"false\" and \"none\" are allowed for secure_boot_enabled")
		}
		if (vm.SecurityProfile.VTPM != "true") && (vm.SecurityProfile.VTPM != "false") {
			return fmt.Errorf("Invalid Entry! Only the values \"true\" and \"false\" are allowed for VTPM")
		}
		if (vm.SecurityProfile.SecureBoot == "none") && (vm.SecurityProfile.VTPM == "true") {
			return fmt.Errorf("Invalid Entry! vTPM cannot be \"true\" when secure-boot is \"none\"")
		}
	}
	if len(vm.OSDiskType) > 0 {
		found := false
		for _, t := range vmconf.AllowedOsDiskTypes() {
			if t == vm.OSDiskType {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("OS disk type '%s' is not included in supported [%s]", vm.OSDiskType, strings.Join(vmconf.AllowedOsDiskTypes(), ","))
		}
	}
	if len(vm.Ports) > 0 {
		if e := validateUniquePorts(vm.Ports, vm.Name); e != nil {
			return e
		}
	}
	if (len(vm.TipNodeSessionID) == 0 && len(vm.ClusterName) != 0) || (len(vm.TipNodeSessionID) != 0 && len(vm.ClusterName) == 0) {
		return fmt.Errorf("Must specify either both 'tip_node_session_id' and 'cluster_name', or neither")
	}
	if !hasOsDisk {
		if isLinux {
			if e := validateLinuxProfile(p.LinuxProfile); e != nil {
				return e
			}
		} else {
			if e := validateWindowsProfile(p.WindowsProfile); e != nil {
				return e
			}
		}
	}
	return nil
}

func validateLinuxProfile(p *LinuxProfile) error {
	if p == nil {
		return fmt.Errorf("LinuxProfile cannot be empty")
	}
	if len(p.AdminUsername) == 0 {
		return fmt.Errorf("LinuxProfile.AdminUsername cannot be empty")
	}
	if len(p.AdminPassword) > 0 && len(p.SSHPubKeys) > 0 {
		return fmt.Errorf("AdminPassword and SSH public keys are mutually exclusive")
	}
	if len(p.AdminPassword) == 0 && len(p.SSHPubKeys) == 0 {
		return fmt.Errorf("Must specify either AdminPassword or SSH public keys")
	}
	for i, key := range p.SSHPubKeys {
		if key == nil || len(key.KeyData) == 0 {
			return fmt.Errorf("SSH public key #%d cannot be empty", i)
		}
	}
	return nil
}

func validateWindowsProfile(p *WindowsProfile) error {
	if p == nil {
		return fmt.Errorf("WindowsProfile cannot be empty")
	}
	if e := validate.Var(p.AdminUsername, "required"); e != nil {
		return fmt.Errorf("WindowsProfile.AdminUsername cannot be empty")
	}
	if e := validate.Var(p.AdminPassword, "required"); e != nil {
		return fmt.Errorf("WindowsProfile.AdminPassword cannot be empty")
	}
	return nil
}

func validateOSImage(p *OSImage) error {
	if p == nil {
		return nil
	}
	if len(p.URL) > 0 {
		if len(p.Publisher) > 0 || len(p.Offer) > 0 || len(p.SKU) > 0 || len(p.Version) > 0 {
			return fmt.Errorf("OS image URL and Publisher/Offer/SKU are mutually exclusive")
		}
	} else {
		if len(p.Publisher) == 0 {
			return fmt.Errorf("OS image Publisher is not set")
		}
		if len(p.Offer) == 0 {
			return fmt.Errorf("OS image Offer is not set")
		}
		if len(p.SKU) == 0 {
			return fmt.Errorf("OS image SKU is not set")
		}
		// version is optional
	}
	return nil
}

func validateOSDisk(p *OSDisk) error {
	if p == nil {
		return nil
	}
	if len(p.VHD) == 0 {
		return fmt.Errorf("OS VHD URL is not set")
	}
	if len(p.StorageAccountID) == 0 {
		return fmt.Errorf("OS VHD storage account ID is not set")
	}
	return nil
}

func (p *Properties) validateDiagnosticsProfile() error {
	if p.DiagnosticsProfile == nil || !p.DiagnosticsProfile.Enabled {
		return nil
	}
	if len(p.DiagnosticsProfile.StorageAccountName) == 0 {
		return fmt.Errorf("DiagnosticsProfile.StorageAccountName cannot be empty string")
	}
	return nil
}

func (p *Properties) validateVnetProfile() error {
	h := p.VnetProfile
	if h == nil {
		return nil
	}
	// existing vnet is uniquely defined by resource group, vnet name, and subnet name
	if len(h.VnetResourceGroup) > 0 {
		if len(h.VnetName) == 0 {
			return fmt.Errorf("vnetProfile.vnetName cannot be empty for existing vnet")
		}
		if len(h.SubnetName) == 0 {
			return fmt.Errorf("vnetProfile.subnetName cannot be empty for existing vnet")
		}
		if len(h.VnetAddress) > 0 {
			return fmt.Errorf("vnetProfile.VnetResourceGroup and vnetProfile.vnetAddress are mutually exclusive")
		}
		if len(h.SubnetAddress) > 0 {
			return fmt.Errorf("vnetProfile.VnetResourceGroup and vnetProfile.subnetAddress are mutually exclusive")
		}
	}
	return nil
}

func validateUniquePorts(ports []int, name string) error {
	portMap := make(map[int]bool)
	for _, port := range ports {
		if _, ok := portMap[port]; ok {
			return fmt.Errorf("VM '%s' has duplicate port '%d', ports must be unique", name, port)
		}
		portMap[port] = true
	}
	return nil
}
