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
      "location": "[resourceGroup().location]",
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
      "name": "[variables('networkSecurityGroupName')]",
      "location": "[resourceGroup().location]",
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
      "location": "[resourceGroup().location]",
      "dependsOn": [
        "[variables('networkSecurityGroupId')]"
      ],
      "properties": {
        "addressSpace": {
          "addressPrefixes": [
            "[parameters('addressPrefix')]"
          ]
        },
        "subnets": [
          {
            "name": "[variables('subnetName')]",
            "properties": {
              "addressPrefix": "[parameters('subnetPrefix')]",
              "networkSecurityGroup": {
                "id": "[variables('networkSecurityGroupId')]"
              }
            }
          }
        ]
      }
    },
          "type": "Microsoft.Network/networkInterfaces",
      "apiVersion": "2019-07-01",
      "name": "[variables('networkInterfaceName')]",
      "location": "[resourceGroup().location]",
      "dependsOn": [
        "[variables('networkSecurityGroupId')]",
        "[concat('Microsoft.Network/virtualNetworks/', variables('virtualNetworkName'))]",
        "[concat('Microsoft.Network/publicIpAddresses/', variables('publicIpAddressName'))]"
      ],
      "properties": {
        "ipConfigurations": [
          {
            "name": "ipConfigNode",
            "properties": {
              "privateIPAllocationMethod": "Dynamic",
              "subnet": {
                "id": "[variables('subnetRef')]"
              },
              "publicIpAddress": {
                "id": "[resourceId('Microsoft.Network/publicIPAddresses',variables('publicIPAddressName'))]"
              }
            }
          }
        ],
        "networkSecurityGroup": {
          "id": "[variables('networkSecurityGroupId')]"
        }
      }
    },
    {
      "condition": "[equals(variables('diagnosticsStorageAction'), 'new')]",
      "type": "Microsoft.Storage/storageAccounts",
      "apiVersion": "2019-06-01",
      "name": "[parameters('diagnosticsStorageAccountName')]",
      "location": "[resourceGroup().location]",
      "kind": "[parameters('diagnosticsStorageAccountKind')]",
      "sku": {
        "name": "[parameters('diagnosticsStorageAccountType')]"
      }
    },
    {
      "type": "Microsoft.Compute/disks",
      "apiVersion": "2021-04-01",
      "name": "[variables('diskName')]",
      "location": "[resourceGroup().location]",
      "sku": {
        "name": "[parameters('osDiskType')]"
      },
      "properties": {
        "osType": "[if(variables('isWindows'), 'Windows', 'Linux')]",
        "SecurityProfile": "[if(variables('isMemoryUnencrypted'), json('null'), variables('diskSecurityProfile'))]",
        "creationData": {
          "createOption": "FromImage",
          "imageReference": {
            "id": "[variables('imageDiskReferenceId')]"
          }
        }
      }
    },
    {
      "type": "Microsoft.Compute/virtualMachines",
      "apiVersion": "2021-07-01",
      "name": "[parameters('vmName')]",
      "location": "[resourceGroup().location]",
      "dependsOn": [
        "[variables('diskName')]",
        "[concat('Microsoft.Network/networkInterfaces/', variables('nicName'))]"
      ],
      "properties": {
        "hardwareProfile": {
          "vmSize": "[parameters('vmSize')]"
        },
{{if HasSecurityProfile}}
        "securityProfile": {
          "uefiSettings": {
            "secureBootEnabled": "[parameters('secureBootEnabled')]",
            "vTPMEnabled": "[parameters('vTPMEnabled')]"
          },
          "securityType" : "{{GetVMSecurityType}}"
        },
{{end}}
        "osProfile": "[variables('osProfile')]",
        "storageProfile": {
          "osDisk": {
            "createOption": "fromImage",
            "managedDisk": "[if(variables('isMemoryUnencrypted'), variables('vmStorageProfileManagedDisk'), variables('vmStorageProfileManagedDiskEncrypted'))]"
          },
          "imageReference": "[variables('imageReference')]"
        },
        "networkProfile": {
          "networkInterfaces": [
            {
              "id": "[resourceId('Microsoft.Network/networkInterfaces', variables('nicName'))]"
            }
          ]
        },
        "diagnosticsProfile": {
          "bootDiagnostics": {
            "enabled": "[equals(parameters('bootDiagnostics'), 'true')]",
            "storageUri": "[if(equals(parameters('bootDiagnostics'), 'true'), reference(resourceId(parameters('diagnosticsStorageAccountResourceGroupName'), 'Microsoft.Storage/storageAccounts', parameters('diagnosticsStorageAccountName')), '2018-02-01').primaryEndpoints['blob'], json('null'))]"
          }
        }
      }
    }
