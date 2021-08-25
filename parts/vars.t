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
    "hasTipNodeSession": "[not(or(empty(parameters('tipNodeSessionId')), empty(parameters('clusterName'))))]",
    "imageReference": "[variables('imageList')[parameters('osImageName')]]",
    "imagePublisher": "[variables('imageReference')['publisher']]",
    "imageOffer": "[variables('imageReference')['offer']]",
    "imageSku": "[variables('imageReference')['sku']]",
    "imageVersion": "[variables('imageReference')['version']]",
    "imageDiskReferenceId": "[concat('/Subscriptions/', subscription().subscriptionId, '/Providers/Microsoft.Compute/Locations/', resourceGroup().location, '/Publishers/', variables('imagePublisher'), '/ArtifactTypes/VMImage/Offers/', variables('imageOffer'), '/Skus/', variables('imageSku'), '/Versions/', variables('imageVersion'))]",
    "diskName": "[concat(parameters('vmName'), '-osDisk')]",
    "nicName": "[concat(parameters('vmName'), '-nic')]",
    "publicIPAddressName": "[concat(parameters('vmName'), '-ip')]",
    "networkSecurityGroupName": "[concat(parameters('vmName'), '-nsg')]",
    "networkSecurityGroupId": "[resourceId(resourceGroup().name, 'Microsoft.Network/networkSecurityGroups', variables('networkSecurityGroupName'))]",
    "vnetSubnetId": "[resourceId(parameters('vnetResourceGroupName'), 'Microsoft.Network/virtualNetworks/subnets/', parameters('vnetName'), parameters('subnetName'))]",
    "isWindows": "[contains(parameters('osImageName'), 'Windows')]",
    "linuxConfiguration": {
      "disablePasswordAuthentication": "true",
      "ssh": {{GetLinuxPublicKeys}}
    },
    "windowsConfiguration": {
      "provisionVmAgent": "true"
    },
    "isMemoryUnencrypted": "[equals(parameters('securityType'), 'Unencrypted')]",
    "diskSecurityType": {
        "Unencrypted": "",
        "VMGuestStateOnly": "ConfidentialVM_VMGuestStateOnlyEncryptedWithPlatformKey",
        "DiskWithVMGuestState": "ConfidentialVM_DiskEncryptedWithPlatformKey"
    },
    "diskSecurityProfile": {
        "SecurityType": "[variables('diskSecurityType')[parameters('securityType')]]"
    },
    "vmStorageProfileManagedDisk": {
      "storageAccountType": "[parameters('osDiskType')]"
    },
    "vmStorageProfileManagedDiskEncrypted": {
      "storageAccountType": "[parameters('osDiskType')]",
      "securityProfile": {
          "confidentialDiskEncryptionType" : "[parameters('securityType')]"
      }
    },
    "diagnosticsStorageAction": "[if(equals(parameters('bootDiagnostics'), 'false'), 'nop', parameters('diagnosticsStorageAccountNewOrExisting'))]",
    "osProfile": {
      "computername": "[parameters('vmName')]",
      "adminUsername": "[parameters('adminUsername')]",
      "adminPasswordOrKey": "[parameters('adminPasswordOrKey')]",
{{if IsLinux .}}
      "linuxConfiguration": "[if(equals(parameters('authenticationType'), 'password'), json('null'), variables('linuxConfiguration'))]"
{{else}}
      "windowsConfiguration": "[variables('windowsConfiguration')]"
{{end}}
    }
