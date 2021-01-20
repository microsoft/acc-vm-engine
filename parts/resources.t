    {
      "apiVersion": "2018-05-01",
      "name": "[concat('ResourceGroupDeployment-', uniqueString(deployment().name))]",
      "type": "Microsoft.Resources/deployments",
      "properties": {
        "mode": "Incremental",
        "template": {
          "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
          "contentVersion": "1.0.0.0",
          "resources": []
        }
      }
    },
    {
      "type": "Microsoft.Network/publicIPAddresses",
      "apiVersion": "2019-02-01",
      "name": "[variables('publicIPAddressName')]",
      "location": "[parameters('location')]",
      "properties": {
{{if HasDNSName .}}
        "dnsSettings": {
          "domainNameLabel": "[parameters('vmName')]"
        },
{{end}}
        "publicIPAllocationMethod": "Dynamic"
      }
    },
    {
      "type": "Microsoft.Network/networkSecurityGroups",
      "apiVersion": "2019-02-01",
      "name": "[variables('nsgName')]",
      "location": "[parameters('location')]",
      "properties": {
        "securityRules": [
          {{GetSecurityRules .VMProfile.Ports}}
        ]
      }
    },
    {
      "condition": "[equals(parameters('vnetNewOrExisting'), 'new')]",
      "type": "Microsoft.Network/virtualNetworks",
      "apiVersion": "2019-09-01",
      "name": "[parameters('vnetName')]",
      "location": "[parameters('location')]",
      "properties": {
        "addressSpace": {
          "addressPrefixes": [
            "[parameters('vnetAddress')]"
          ]
        },
        "subnets": [
          {
            "name": "[parameters('subnetName')]",
            "properties": {
              "addressPrefix": "[parameters('subnetAddress')]"
            }
          }
        ]
      }
    },
    {
      "type": "Microsoft.Network/networkInterfaces",
      "apiVersion": "2019-07-01",
      "name": "[variables('nicName')]",
      "location": "[parameters('location')]",
      "dependsOn": [
        "[variables('publicIPAddressName')]",
        "[parameters('vnetName')]",
        "[variables('nsgName')]"
      ],
      "properties": {
        "ipConfigurations": [
          {
            "name": "ipConfigNode",
            "properties": {
              "privateIPAllocationMethod": "Dynamic",
              "subnet": {
                "id": "[variables('vnetSubnetId')]"
              },
              "publicIpAddress": {
                "id": "[resourceId('Microsoft.Network/publicIPAddresses',variables('publicIPAddressName'))]"
              }
            }
          }
        ]
        ,"networkSecurityGroup": {
          "id": "[variables('nsgId')]"
        }
      }
    },
    {
      "condition": "[equals(variables('diagnosticsStorageAction'), 'new')]",
      "type": "Microsoft.Storage/storageAccounts",
      "apiVersion": "2019-06-01",
      "name": "[parameters('diagnosticsStorageAccountName')]",
      "location": "[parameters('location')]",
      "kind": "[parameters('diagnosticsStorageAccountKind')]",
      "sku": {
        "name": "[parameters('diagnosticsStorageAccountType')]"
      }
    },
{{if HasCustomOsImage}}
     {
      "type": "Microsoft.Compute/images",
      "apiVersion": "2018-06-01",
      "name": "CustomImage",
      "location": "[parameters('location')]",
      "properties": {
        "storageProfile": {
          "osDisk": {
            "osType": "[parameters('osType')]",
            "osState": "Generalized",
            "blobUri": "[parameters('osImageURL')]",
            "storageAccountType": "Standard_LRS"
          }
        }
      }
    },
{{end}}
{{if HasAttachedOsDisk}}
    {
      "type": "Microsoft.Compute/disks",
      "apiVersion": "2020-09-30",
      "name": "CustomDisk",
      "location": "[parameters('location')]",
{{if HasAttachedOsDiskVMGS}}
      "tags": {
        "VmgsBlobUri": "[parameters('osDiskVmgsURL')]"
     },
{{end}}
      "sku": {
        "name": "Standard_LRS"
      },
      "properties": {
        "osType": "[parameters('osType')]",
        "hyperVGeneration": "V2",
        "securityType" : "{{GetSecurityType}}",
        "creationData": {
          "createOption": "Import",
          "storageAccountId": "[parameters('osDiskStorageAccountID')]",
          "sourceUri": "[parameters('osDiskURL')]"
        }
      }
    },
{{end}}
{{if HasTipNodeSession}}
    {
      "type": "Microsoft.Compute/availabilitySets",
      "apiVersion": "2020-06-01",
      "name": "[variables('availabilitySetName')]",
      "location": "[resourceGroup().location]",
      "properties": {
        "platformUpdateDomainCount": "1",
        "platformFaultDomainCount": "1",
        "internalData": {
          "pinnedFabricCluster": "[parameters('clusterName')]"
        }
      },
      "tags": {
        "TipNode.SessionId": "[parameters('tipNodeSessionId')]"
      },
      "sku": {
        "name": "aligned"
      }
    },
{{end}}
    {
      "type": "Microsoft.Compute/virtualMachines",
      "apiVersion": "2020-12-01",
      "name": "[parameters('vmName')]",
      "location": "[parameters('location')]",
      "dependsOn": [
{{if HasCustomOsImage}}
        "CustomImage",
{{end}}
{{if HasAttachedOsDisk}}
        "CustomDisk",
{{end}}
{{if HasTipNodeSession}}
        "[variables('availabilitySetName')]",
{{end}}
        "[concat('Microsoft.Network/networkInterfaces/', variables('nicName'))]"
      ],
      "tags":
      {
        "Platform.SecurityType": "{{GetSecurityType}}",
        "creationSource" : "['acc-vm-engine']"
      },
      "properties": {
        "hardwareProfile": {
          "vmSize": "[parameters('vmSize')]"
        },
        "securityProfile": {
          "uefiSettings": {
            "secureBootEnabled": "[parameters('secureBootEnabled')]",
            "vTPMEnabled": "[parameters('vTPMEnabled')]"
          }
        },
{{if not HasTipNodeSession}}
        "osProfile": "[variables('osProfile')]",
{{end}}
        "storageProfile": "[variables('storageProfile')]",
        "networkProfile": {
          "networkInterfaces": [
            {
              "id": "[resourceId('Microsoft.Network/networkInterfaces', variables('nicName'))]"
            }
          ]
        },
{{if HasTipNodeSession}}
        "availabilitySet": {
          "id": "[resourceId('Microsoft.Compute/availabilitySets', variables('availabilitySetName'))]"
        },
{{end}}
        "diagnosticsProfile": {
          "bootDiagnostics": {
            "enabled": "[equals(parameters('bootDiagnostics'), 'true')]",
            "storageUri": "[if(equals(parameters('bootDiagnostics'), 'true'), reference(resourceId(parameters('diagnosticsStorageAccountResourceGroupName'), 'Microsoft.Storage/storageAccounts', parameters('diagnosticsStorageAccountName')), '2018-02-01').primaryEndpoints['blob'], json('null'))]"
          }
        }
      }
    }
