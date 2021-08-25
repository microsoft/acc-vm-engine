{
  "$schema": "https://schema.management.azure.com/schemas/2018-05-01/subscriptionDeploymentTemplate.json#",
  "contentVersion": "1.0.0.1",
  "parameters": {
    {{template "params.t" .}}
  },
  "variables": {
    {{template "vars.t" .}}
  },
  "resources": [
    {{template "resources.t" .}}
  ],
  "outputs": {
    {{template "outputs.t" .}}
  }
}
