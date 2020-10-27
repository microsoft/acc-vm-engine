package engine

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/microsoft/acc-vm-engine/pkg/api"
	log "github.com/sirupsen/logrus"
)

// setPropertiesDefaults for the container Properties
func setPropertiesDefaults(vm *api.APIModel) {
	if len(vm.Properties.VMProfile.Name) == 0 {
		log.Warnf("Missing VM Name. Setting to %s", vm.VMConfigurator.DefaultVMName())
		vm.Properties.VMProfile.Name = vm.VMConfigurator.DefaultVMName()
	}
	// set network defaults
	if vm.Properties.VnetProfile == nil {
		vm.Properties.VnetProfile = &api.VnetProfile{}
	}
	if !vm.Properties.VnetProfile.IsCustomVNET() {
		if len(vm.Properties.VnetProfile.VnetAddress) == 0 {
			vm.Properties.VnetProfile.VnetAddress = api.DefaultVnet
		}
		if len(vm.Properties.VnetProfile.SubnetAddress) == 0 {
			vm.Properties.VnetProfile.SubnetAddress = api.DefaultSubnet
		}
	}
	if vm.Properties.VMProfile.OSImage == nil {
		vm.Properties.VMProfile.OSImage = vm.VMConfigurator.OSImage()
	}
	if len(vm.Properties.VMProfile.OSDiskType) == 0 {
		vm.Properties.VMProfile.OSDiskType = vm.VMConfigurator.DefaultOsDiskType()
	}
}

func combineValues(inputs ...string) string {
	valueMap := make(map[string]string)
	for _, input := range inputs {
		applyValueStringToMap(valueMap, input)
	}
	return mapToString(valueMap)
}

func applyValueStringToMap(valueMap map[string]string, input string) {
	values := strings.Split(input, ",")
	for index := 0; index < len(values); index++ {
		// trim spaces (e.g. if the input was "foo=true, bar=true" - we want to drop the space after the comma)
		value := strings.Trim(values[index], " ")
		valueParts := strings.Split(value, "=")
		if len(valueParts) == 2 {
			valueMap[valueParts[0]] = valueParts[1]
		}
	}
}

func mapToString(valueMap map[string]string) string {
	// Order by key for consistency
	keys := []string{}
	for key := range valueMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	for _, key := range keys {
		buf.WriteString(fmt.Sprintf("%s=%s,", key, valueMap[key]))
	}
	return strings.TrimSuffix(buf.String(), ",")
}
