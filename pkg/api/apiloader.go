package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/microsoft/acc-vm-engine/pkg/helpers"
)

// Apiloader represents the object that loads api model
type Apiloader struct {
}

// VMConfigurator manages VM specific configuration
type VMConfigurator interface {
	DefaultVMName() string
	OSImageName() string
	DefaultOsDiskType() string
	AllowedOsDiskTypes() []string
	AllowedVMSizes() []string
	DefaultVMSize() string
}

// LoadVMFromFile loads an API Model from a JSON file
func (a *Apiloader) LoadVMFromFile(jsonFile string, validate, isUpdate bool, sshPubKeys []string) (*APIModel, error) {
	contents, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %s", jsonFile, err.Error())
	}
	return a.LoadVM(contents, validate, isUpdate, sshPubKeys)
}

// LoadVM loads and validates an API Model
func (a *Apiloader) LoadVM(contents []byte, validate, isUpdate bool, sshPubKeys []string) (*APIModel, error) {
	vm := &APIModel{}
	err := json.Unmarshal(contents, vm)
	if err != nil {
		return nil, err
	}
	err = checkJSONKeys(contents, reflect.TypeOf(*vm))
	if err != nil {
		return nil, err
	}
	var osImageName string
	if vm.Properties != nil && vm.Properties.VMProfile != nil {
		osImageName = vm.Properties.VMProfile.OSImageName 
	}
	if vm.VMConfigurator, err = getVMConfigurator(vm.VMCategory, osImageName); err != nil {
		return nil, err
	}
	// add SSH public keys from command line arguments
	if vm.Properties.LinuxProfile != nil {
		for _, key := range sshPubKeys {
			vm.Properties.LinuxProfile.SSHPubKeys = append(vm.Properties.LinuxProfile.SSHPubKeys, &PublicKey{KeyData: key})
		}
	}
	if err := vm.Properties.Validate(vm.VMConfigurator, isUpdate); validate && err != nil {
		return nil, err
	}
	return vm, nil
}

// SerializeVM takes an unversioned container service and returns the bytes
func (a *Apiloader) SerializeVM(vm *APIModel) ([]byte, error) {
	return helpers.JSONMarshalIndent(vm, "", "  ", false)
}

func getVMConfigurator(vmcat VMCategory, osImageName string) (VMConfigurator, error) {
	switch vmcat {
	case TVM:
		var newOsName OSName
		return NewTVMConfigurator(newOsName)
	case CVM:
		return NewCVMConfigurator()
	default:
		return nil, fmt.Errorf("unsupported VM category %q", vmcat)
	}
}
