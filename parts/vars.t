    "nicName": "[concat(parameters('vmName'), '-nic')]",
    "publicIPAddressName": "[concat(parameters('vmName'), '-ip')]",
    "nsgName": "[concat(parameters('vmName'), '-nsg')]",
    "nsgId": "[resourceId(resourceGroup().name, 'Microsoft.Network/networkSecurityGroups', variables('nsgName'))]",
    "vnetSubnetId": "[resourceId(parameters('vnetResourceGroupName'), 'Microsoft.Network/virtualNetworks/subnets/', parameters('vnetName'), parameters('subnetName'))]",
    "linuxConfiguration": {
      "disablePasswordAuthentication": "true",
      "ssh": {{GetLinuxPublicKeys}}
    },
    "singleQuote": "'",
    "windowsConfiguration": {
      "provisionVmAgent": "true"
    },
    "diagnosticsStorageAction": "[if(equals(parameters('bootDiagnostics'), 'false'), 'nop', parameters('diagnosticsStorageAccountNewOrExisting'))]",
    "osProfile": {
      "computername": "[parameters('vmName')]",
      "adminUsername": "[parameters('adminUsername')]",
      "adminPassword": "[parameters('adminPassword')]",
{{if IsLinux .}}
      "linuxConfiguration": "[if(equals(parameters('authenticationType'), 'password'), json('null'), variables('linuxConfiguration'))]"
{{else}}
      "windowsConfiguration": "[variables('windowsConfiguration')]"
{{end}}
    },
    "storageProfile": {
      "imageReference": {
{{if HasCustomImage}}
        "id": "[resourceId('Microsoft.Compute/images','CustomImage')]"
{{else}}
        "publisher": "[parameters('osImagePublisher')]",
        "offer": "[parameters('osImageOffer')]",
        "sku": "[parameters('osImageSKU')]",
        "version": "[parameters('osImageVersion')]"
{{end}}
      },
      {{GetDataDisks .}}
      "osDisk": {
        "caching": "ReadWrite",
        "createOption": "FromImage",
        "managedDisk": {
          "storageAccountType": "[parameters('osDiskType')]"
        }
      }
    }
