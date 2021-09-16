    "imageList": {
      "Windows Server 2022 Gen 2": {
        "publisher": "microsoftwindowsserver",
        "offer": "windows-cvm",
        "sku":  "2022-datacenter-cvm",
        "version": "latest"
      },
      "Windows Server 2019 Gen 2": {
        "publisher": "microsoftwindowsserver",
        "offer":  "windows-cvm",
        "sku":  "2019-datacenter-cvm",
        "version": "latest"
      },
      "Ubuntu 20.04 LTS Gen 2": {
        "publisher": "Canonical",
        "offer": "0001-com-ubuntu-confidential-vm-experimental",
        "sku": "20_04-gen2",
        "version": "20.04.20210309"
      },
      "Ubuntu 18.04 LTS Gen 2": {
        "publisher": "Canonical",
        "offer": "0001-com-ubuntu-confidential-vm-experimental",
        "sku": "18_04-gen2",
        "version": "18.04.20210309"
      }
    },
    "imageReference": "[variables('imageList')[parameters('osImageName')]]",
    "networkInterfaceName": "[concat(parameters('vmName'), '-nic')]",
    "publicIPAddressName": "[concat(parameters('vmName'), '-ip')]",
    "networkSecurityGroupName": "[concat(parameters('vmName'), '-nsg')]",
    "networkSecurityGroupId": "[resourceId(resourceGroup().name, 'Microsoft.Network/networkSecurityGroups', variables('networkSecurityGroupName'))]",
    "virtualNetworkName": "[concat(parameters('vmName'), '-vnet')]",
    "virtualNetworkId": "[resourceId(resourceGroup().name, 'Microsoft.Network/virtualNetworks', variables('virtualNetworkName'))]",
    "subnetRef": "[concat(variables('virtualNetworkId'), '/subnets/', variables('subnetName'))]",
    "subnetName": "[concat(parameters('vmName'), 'Subnet')]",
    "vnetSubnetId": "[resourceId(parameters('vnetResourceGroupName'), 'Microsoft.Network/virtualNetworks/subnets/', variables('virtualNetworkName'), parameters('subnetName'))]",
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
    {{if HasTipNodeSession}}
    "availabilitySetName": "[concat(parameters('vmName'), '-availSet')]",
    {{end}}
    "isMemoryUnencrypted": "[equals(parameters('securityType'), 'Unencrypted')]",
    "vmStorageProfileManagedDisk": {
      "storageAccountType": "[parameters('osDiskType')]"
    },
    "vmStorageProfileManagedDiskEncrypted": {
      "storageAccountType": "[parameters('osDiskType')]",
      "securityProfile": {
          "securityEncryptionType" : "[parameters('securityType')]"
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
