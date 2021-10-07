
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
      "name": "[variables('virtualNetworkName')]",
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
    {
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
      "apiVersion": "2021-07-01",
      "name": "[parameters('vmName')]",
      "location": "[resourceGroup().location]",
      "dependsOn": [
        {{if HasTipNodeSession}}
        "[variables('availabilitySetName')]",
        {{end}} 
        "[concat('Microsoft.Network/networkInterfaces/', variables('networkInterfaceName'))]"
      ],
      "properties": {
        "hardwareProfile": {
          "vmSize": "[parameters('vmSize')]"
        },
        "osProfile": {
          "computerName": "[parameters('vmName')]",
          "adminUsername": "[parameters('adminUsername')]",
          "adminPassword": "[parameters('adminPasswordOrKey')]",
          "linuxConfiguration": "[if(equals(parameters('authenticationType'), 'password'), json('null'), variables('linuxConfiguration'))]",
          "windowsConfiguration": "[if(variables('isWindows'), variables('windowsConfiguration'), json('null'))]"
        },
        "securityProfile": "[if(variables('isMemoryUnencrypted'), json('null'), variables('vmSecurityProfile'))]",
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
              "id": "[resourceId('Microsoft.Network/networkInterfaces', variables('networkInterfaceName'))]"
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
