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
    "networkInterfaceName": "[concat(parameters('vmName'), '-nic')]",
    "publicIPAddressName": "[concat(parameters('vmName'), '-ip')]",
    "networkSecurityGroupName": "[concat(parameters('vmName'), '-nsg')]",
    "networkSecurityGroupId": "[resourceId(resourceGroup().name, 'Microsoft.Network/networkSecurityGroups', variables('networkSecurityGroupName'))]",
    "virtualNetworkName": "[concat(parameters('vmName'), '-vnet')]",
    "virtualNetworkId": "[resourceId(resourceGroup().name, 'Microsoft.Network/virtualNetworks', variables('virtualNetworkName'))]",
    "subnetRef": "[concat(variables('virtualNetworkId'), '/subnets/', variables('subnetName'))]",
    "subnetName": "[concat(parameters('vmName'), 'Subnet')]",
    "vnetSubnetId": "[resourceId(parameters('vnetResourceGroupName'), 'Microsoft.Network/virtualNetworks/subnets/', parameters('vnetName'), parameters('subnetName'))]",
    "isWindows": "[contains(parameters('osImageName'), 'Windows')]",
    "linuxConfiguration": {
      "disablePasswordAuthentication": "true",
      "ssh": {
        "publicKeys": [
          {
            "keyData": "[parameters('adminPasswordOrKey')]",
            "path": "[concat('/home/', parameters('adminUsername'), '/.ssh/authorized_keys')]"
          }
        ]
      }
    },
    "windowsConfiguration": {
      "enableAutomaticUpdates": "true",
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
    "vmSecurityProfile": {
      "uefiSettings" : {
        "secureBootEnabled": "[parameters('secureBootEnabled')]",
        "vTpmEnabled": "true"
      },
      "securityType" : "ConfidentialVM"
    }
