package engine

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"text/template"

	"github.com/microsoft/acc-vm-engine/pkg/api"
	"github.com/microsoft/acc-vm-engine/pkg/helpers"
)

var (
	baseFile      = "base.t"
	templateFiles = []string{baseFile, "outputs.t", "params.t", "resources.t", "vars.t"}
)

// TemplateGenerator represents the object that performs the template generation.
type TemplateGenerator struct {
}

// InitializeTemplateGenerator creates a new template generator object
func InitializeTemplateGenerator(vmconfig api.VMConfigurator) (*TemplateGenerator, error) {
	t := &TemplateGenerator{}

	if err := t.verifyFiles(vmconfig); err != nil {
		return nil, err
	}

	return t, nil
}

// GenerateTemplate generates the template from the API Model
func (t *TemplateGenerator) GenerateTemplate(vm *api.APIModel, generatorCode string) (templateRaw string, parametersRaw string, err error) {
	// named return values are used in order to set err in case of a panic
	templateRaw = ""
	parametersRaw = ""
	err = nil

	var templ *template.Template

	properties := vm.Properties

	setPropertiesDefaults(vm)

	templ = template.New("vm template").Funcs(t.getTemplateFuncMap(vm))

	for _, file := range templateFiles {
		bytes, e := Asset(file)
		if e != nil {
			err = fmt.Errorf("Error reading file %s, Error: %s", file, e.Error())
			return templateRaw, parametersRaw, err
		}
		if _, err = templ.New(file).Parse(string(bytes)); err != nil {
			return templateRaw, parametersRaw, err
		}
	}
	// template generation may have panics in the called functions.  This catches those panics
	// and ensures the panic is returned as an error
	defer func() {
		if r := recover(); r != nil {
			s := debug.Stack()
			err = fmt.Errorf("%v - %s", r, s)

			// invalidate the template and the parameters
			templateRaw = ""
			parametersRaw = ""
		}
	}()

	var b bytes.Buffer
	if err = templ.ExecuteTemplate(&b, baseFile, properties); err != nil {
		return templateRaw, parametersRaw, err
	}
	templateRaw = b.String()

	var parametersMap paramsMap
	if parametersMap, err = getParameters(vm, generatorCode); err != nil {
		return templateRaw, parametersRaw, err
	}

	var parameterBytes []byte
	if parameterBytes, err = helpers.JSONMarshalIndent(parametersMap, "", "  ", false); err != nil {
		return templateRaw, parametersRaw, err
	}
	parametersRaw = string(parameterBytes)

	return templateRaw, parametersRaw, err
}

func (t *TemplateGenerator) verifyFiles(vmconfig api.VMConfigurator) error {
	for _, file := range templateFiles {
		if _, err := Asset(file); err != nil {
			return fmt.Errorf("template file %s does not exist", file)
		}
	}
	return nil
}

// getTemplateFuncMap returns all functions used in template generation
func (t *TemplateGenerator) getTemplateFuncMap(vm *api.APIModel) template.FuncMap {
	return template.FuncMap{
		"RequiresFakeAgentOutput": func() bool {
			return false
		},
		"IsPublic": func(ports []int) bool {
			return len(ports) > 0
		},

		"IsPrivateCluster": func() bool {
			return false
		},
		"GetLBRules": func(name string, ports []int) string {
			return getLBRules(name, ports)
		},
		"GetProbes": func(ports []int) string {
			return getProbes(ports)
		},
		"GetSecurityRules": func(ports []int) string {
			return getSecurityRules(ports)
		},
		"GetLinuxPublicKeys": func() string {
			if vm.Properties.LinuxProfile == nil {
				return `"json('null')"`
			}
			keyTempl := `          {
            "keyData": "%s",
            "path": "/home/%s/.ssh/authorized_keys"
          }`
			keyData := make([]string, len(vm.Properties.LinuxProfile.SSHPubKeys))
			for i, key := range vm.Properties.LinuxProfile.SSHPubKeys {
				keyData[i] = fmt.Sprintf(keyTempl, key.KeyData, vm.Properties.LinuxProfile.AdminUsername)
			}
			sshTempl := `{
        "publicKeys": [
%s
        ]
      }`
			return fmt.Sprintf(sshTempl, strings.Join(keyData, ",\n"))
		},
		"GetVMSizes": func() string {
			return api.GetVMSizes(vm.VMConfigurator)
		},
		"GetOsDiskTypes": func() string {
			return api.GetOsDiskTypes(vm.VMConfigurator)
		},
		"GetDataDisks": func(p *api.Properties) string {
			return getDataDisks(p.VMProfile)
		},
		"HasSecurityProfile": func() bool {
			return (vm.Properties.VMProfile.SecurityProfile != nil)
		},
		"GetSecurityType": func() string {
			switch vm.VMCategory {
			case api.TVM:
				return "SecureBoot"
			case api.CVM:
				return "ConfidentialVM_DiskEncryptedWithPlatformKey"
			default:
				return "None"
			}
		},
		"GetVMSecurityType": func() string {
			switch vm.VMCategory {
			case api.TVM:
				return ""
			case api.CVM:
				return "ConfidentialVM"
			default:
				return "None"
			}
		},
		"HasTipNodeSession": func() bool {
			return len(vm.Properties.VMProfile.TipNodeSessionID) > 0
		},
		"Base64": func(s string) string {
			return base64.StdEncoding.EncodeToString([]byte(s))
		},
		"WrapAsVariable": func(s string) string {
			return fmt.Sprintf("',variables('%s'),'", s)
		},
		"WrapAsVerbatim": func(s string) string {
			return fmt.Sprintf("',%s,'", s)
		},
		"IsLinux": func(p *api.Properties) bool {
			return p.VMProfile.OSType == api.Linux
		},
		"IsWindows": func(p *api.Properties) bool {
			return p.VMProfile.OSType == api.Windows
		},
		"HasCustomOsImage": func() bool {
			return vm.Properties.VMProfile.HasCustomOsImage()
		},
		"HasAttachedOsDisk": func() bool {
			return vm.Properties.VMProfile.HasAttachedOsDisk()
		},
		"HasAttachedOsDiskVMGS": func() bool {
			return vm.Properties.VMProfile.HasAttachedOsDiskVMGS()
		},
		"HasDNSName": func(p *api.Properties) bool {
			return p.VMProfile.HasDNSName
		},
		// inspired by http://stackoverflow.com/questions/18276173/calling-a-template-with-several-pipeline-parameters/18276968#18276968
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"loop": func(min, max int) []int {
			var s []int
			for i := min; i <= max; i++ {
				s = append(s, i)
			}
			return s
		},
		"subtract": func(a, b int) int {
			return a - b
		},
	}
}
