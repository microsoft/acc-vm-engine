{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    {{template "tvm/params.t" .}}
  },
  "variables": {
    {{template "tvm/vars.t" .}}
  },
  "resources": [
    {{template "tvm/resources.t" .}}
  ],
  "outputs": {
    {{template "tvm/outputs.t" .}}
  }
}
