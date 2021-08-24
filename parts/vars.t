    "imageList": {
      "Windows Server 2022 Gen 2": {
        "publisher": "AZURERT.PIRCORE.CAPSBVT",
        "offer": "longlivedconfidentialvm",
        "sku":  "WindowsServer2022",
        "version": "0.0.1"
      },
      "Windows Server 2019 Gen 2": {
        "publisher":  "AZURERT.PIRCORE.CAPSBVT",
        "offer":  "longlivedconfidentialvm",
        "sku":  "WindowsServer2019-2",
        "version": "0.0.1"
      },
      "Ubuntu 20.04 LTS Gen 2": {
        "publisher":  "AZURERT.PIRCORE.CAPSBVT",
        "offer":  "longlivedconfidentialvm",
        "sku":  "Ubuntu20.04",
        "version": "0.0.1"
      },
      "Ubuntu 18.04 LTS Gen 2": {
        "publisher":  "AZURERT.PIRCORE.CAPSBVT",
        "offer":  "longlivedconfidentialvm",
        "sku":  "Ubuntu18.04",
        "version": "0.0.1"
      }
    },
    "imageReference": "[variables('imageList')[parameters('osImageName')]]",
    "imagePublisher": "[variables('imageReference')['publisher']]",
    "imageOffer": "[variables('imageReference')['offer']]",
    "imageSku": "[variables('imageReference')['sku']]",
    "imageVersion": "[variables('imageReference')['version']]",
    "imageDiskReferenceId": "[concat('/Subscriptions/', subscription().subscriptionId, '/Providers/Microsoft.Compute/Locations/', resourceGroup().location, '/Publishers/', variables('imagePublisher'), '/ArtifactTypes/VMImage/Offers/', variables('imageOffer'), '/Skus/', variables('imageSku'), '/Versions/', variables('imageVersion'))]",
    "diskName": "[concat(parameters('vmName'), '-osDisk')]",
    "nicName": "[concat(parameters('vmName'), '-nic')]",
    "publicIPAddressName": "[concat(parameters('vmName'), '-ip')]",
    "nsgName": "[concat(parameters('vmName'), '-nsg')]",
    "nsgId": "[resourceId(resourceGroup().name, 'Microsoft.Network/networkSecurityGroups', variables('nsgName'))]",
    "vnetSubnetId": "[resourceId(parameters('vnetResourceGroupName'), 'Microsoft.Network/virtualNetworks/subnets/', parameters('vnetName'), parameters('subnetName'))]",
    "linuxConfiguration": {
      "disablePasswordAuthentication": "true",
      "ssh": {{GetLinuxPublicKeys}}
    },
    "windowsConfiguration": {
      "provisionVmAgent": "true"
    },
    "diagnosticsStorageAction": "[if(equals(parameters('bootDiagnostics'), 'false'), 'nop', parameters('diagnosticsStorageAccountNewOrExisting'))]",
{{if HasTipNodeSession}}
    "availabilitySetName": "[concat(parameters('vmName'), '-availSet')]",
{{end}}
{{if not HasAttachedOsDisk}}
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
{{end}}
    "storageProfile": {
{{if not HasAttachedOsDisk}}
      "imageReference": {"[variables('imageReference')]"},
{{end}}
      {{GetDataDisks .}}
      "osDisk": {
        "caching": "ReadWrite",
{{if HasAttachedOsDisk}}
        "osType": "[parameters('osType')]",
        "createOption": "Attach",
        "managedDisk": {
          "id": "[resourceId('Microsoft.Compute/disks','CustomDisk')]"
        }
{{else}}
        "createOption": "FromImage",
        "managedDisk": {
          "storageAccountType": "[parameters('osDiskType')]",
          "securityProfile": {
                    "confidentialDiskEncryptionType" : "DiskWithPlatformKey"
          }
        }
{{end}}
      }
    }
