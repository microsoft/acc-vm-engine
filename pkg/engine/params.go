package engine

import (
	"strconv"

	"github.com/microsoft/acc-vm-engine/pkg/api"
)

func getParameters(vm *api.APIModel, generatorCode string) (paramsMap, error) {
	properties := vm.Properties
	location := vm.Location
	parametersMap := paramsMap{}

	if len(vm.Location) > 0 {
		addValue(parametersMap, "location", location)
	}
	addValue(parametersMap, "vmName", properties.VMProfile.Name)
	addValue(parametersMap, "vmSize", properties.VMProfile.VMSize)
	if len(properties.VMProfile.OSDiskType) > 0 {
		addValue(parametersMap, "osDiskType", properties.VMProfile.OSDiskType)
	}
	if properties.VMProfile.OSImage.HasCustomImage() {
		addValue(parametersMap, "osImageSourceUrl", properties.VMProfile.OSImage.URL)
	} else {
		addValue(parametersMap, "osImagePublisher", properties.VMProfile.OSImage.Publisher)
		addValue(parametersMap, "osImageOffer", properties.VMProfile.OSImage.Offer)
		addValue(parametersMap, "osImageSKU", properties.VMProfile.OSImage.SKU)
		if len(properties.VMProfile.OSImage.Version) > 0 {
			addValue(parametersMap, "osImageVersion", properties.VMProfile.OSImage.Version)
		}
	}
	if properties.LinuxProfile != nil {
		addValue(parametersMap, "adminUsername", properties.LinuxProfile.AdminUsername)
		if len(properties.LinuxProfile.AdminPassword) > 0 {
			addValue(parametersMap, "authenticationType", "password")
			addValue(parametersMap, "adminPassword", properties.LinuxProfile.AdminPassword)
		} else {
			addValue(parametersMap, "authenticationType", "sshPublicKey")
		}
	}
	if properties.WindowsProfile != nil {
		addValue(parametersMap, "adminUsername", properties.WindowsProfile.AdminUsername)
		addValue(parametersMap, "adminPassword", properties.WindowsProfile.AdminPassword)
	}
	if properties.VMProfile.SecureBoot != nil {
		addValue(parametersMap, "secureBoot", strconv.FormatBool(*properties.VMProfile.SecureBoot))
	}
	if properties.VMProfile.VTPM != nil {
		addValue(parametersMap, "vTPM", strconv.FormatBool(*properties.VMProfile.VTPM))
	}
	if properties.VnetProfile.IsCustomVNET() {
		addValue(parametersMap, "vnetNewOrExisting", "existing")
		addValue(parametersMap, "vnetResourceGroupName", properties.VnetProfile.VnetResourceGroup)
		addValue(parametersMap, "vnetName", properties.VnetProfile.VnetName)
		addValue(parametersMap, "subnetName", properties.VnetProfile.SubnetName)
	} else {
		addValue(parametersMap, "vnetNewOrExisting", "new")
		addValue(parametersMap, "subnetAddress", properties.VnetProfile.SubnetAddress)
	}
	if properties.DiagnosticsProfile != nil && properties.DiagnosticsProfile.Enabled {
		addValue(parametersMap, "bootDiagnostics", "true")
		addValue(parametersMap, "diagnosticsStorageAccountName", properties.DiagnosticsProfile.StorageAccountName)
		if properties.DiagnosticsProfile.IsNewStorageAccount {
			addValue(parametersMap, "diagnosticsStorageAccountNewOrExisting", "new")
		} else {
			addValue(parametersMap, "diagnosticsStorageAccountNewOrExisting", "existing")
		}
	} else {
		addValue(parametersMap, "bootDiagnostics", "false")
	}
	return parametersMap, nil
}
