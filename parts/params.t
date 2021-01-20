    "vmName": {
      "type": "string",
      "metadata": {
        "description": "Name of the VM."
      }
    },
    "vmSize": {
      "type": "string",
      {{GetAllowedVMSizes}}
      "metadata": {
        "description": "Size of the VM."
      }
    },
    "authenticationType": {
      "type": "string",
      "defaultValue": "password",
      "allowedValues": [
        "password",
        "sshPublicKey"
      ],
      "metadata": {
        "description": "Type of authentication to use on Linux virtual machine."
      }
    },
    "adminUsername": {
      "type": "string",
      "defaultValue": "azureuser",
      "metadata": {
        "description": "Username for the Virtual Machine."
      }
    },
    "adminPassword": {
      "type": "securestring",
      "defaultValue": "",
      "metadata": {
        "description": "Password for the Virtual Machine."
      }
    },
    "osType": {
      "type": "string",
      "defaultValue": "",
      "metadata": {
        "description": "OS type."
      }
    },
    "osImagePublisher": {
      "type": "string",
      "defaultValue": "",
      "metadata": {
        "description": "OS image publisher."
      }
    },
    "osImageOffer": {
      "type": "string",
      "defaultValue": "",
      "metadata": {
        "description": "OS image offer."
      }
    },
    "osImageSKU": {
      "type": "string",
      "defaultValue": "",
      "metadata": {
        "description": "OS image SKU."
      }
    },
    "osImageVersion": {
      "type": "string",
      "defaultValue": "latest",
      "metadata": {
        "description": "OS image version."
      }
    },
    "osImageURL": {
      "type": "string",
      "defaultValue": "",
      "metadata": {
        "description": "OS image URL."
      }
    },
    "osDiskURL": {
      "type": "string",
      "defaultValue": "",
      "metadata": {
        "description": "OS VHD URL."
      }
    },
    "osDiskVmgsURL": {
      "type": "string",
      "defaultValue": "",
      "metadata": {
        "description": "OS VMGS URL."
      }
    },
    "osDiskStorageAccountID": {
      "type": "string",
      "defaultValue": "",
      "metadata": {
        "description": "ID of the OS disk storage account."
      }
    },
    "osDiskType": {
      "type": "string",
      {{GetOsDiskTypes}}
      "metadata": {
        "description": "Type of managed disk to create."
      }
    },
    "vnetNewOrExisting": {
      "type": "string",
      "defaultValue": "new",
      "allowedValues": [
        "new",
        "existing"
      ],
      "metadata": {
        "description": "Determines whether or not a new virtual network should be provisioned"
      }
    },
    "vnetName": {
      "type": "string",
      "defaultValue": "[concat(resourceGroup().name, '-vnet')]",
      "metadata": {
        "description": "Name of the virtual network (alphanumeric, hyphen, underscore, period)."
      },
      "minLength": 2,
      "maxLength": 64
    },
    "vnetResourceGroupName": {
      "type": "string",
      "defaultValue": "[resourceGroup().name]",
      "metadata": {
        "description": "Name of the resource group for the existing virtual network."
      }
    },
    "vnetAddress": {
      "type": "string",
      "defaultValue": "{{.VnetProfile.VnetAddress}}",
      "metadata": {
        "description": "VNET address space"
      }
    },
    "subnetName": {
      "type": "string",
      "defaultValue": "[concat(resourceGroup().name, '-subnet')]",
      "metadata": {
        "description": "Name of the subnet."
      }
    },
    "subnetAddress": {
      "type": "string",
      "defaultValue": "{{.VnetProfile.SubnetAddress}}",
      "metadata": {
        "description": "Sets the subnet of the VM."
      }
    },
    "tipNodeSessionId": {
      "type": "string",
      "defaultValue": "",
      "metadata": {
        "description": "TIP Node session ID"
      }
    },
    "clusterName": {
      "type": "string",
      "defaultValue": "",
      "metadata": {
        "description": "Cluster"
      }
    },
    "secureBootEnabled": {
      "type": "string",
      "defaultValue": "true",
      "allowedValues": [
        "true",
         "false"
      ],
      "metadata": {
        "description": "Secure Boot setting of the VM."
      }
    },
    "vTPMEnabled": {
      "type": "string",
      "defaultValue": "true",
      "allowedValues": [
        "true",
        "false"
      ],
      "metadata": {
        "description": "vTPM setting of the VM."
      }
    },
    "bootDiagnostics": {
      "type": "string",
      "defaultValue": "false",
      "allowedValues": [
        "true",
        "false"
      ],
      "metadata": {
        "description": "Boot diagnostics setting of the VM."
      }
    },
    "diagnosticsStorageAccountNewOrExisting": {
      "type": "string",
      "defaultValue": "new",
      "allowedValues": [
        "new",
        "existing"
      ],
      "metadata": {
        "description": "Determines whether or not a new storage account should be provisioned."
      }
    },
    "diagnosticsStorageAccountName": {
      "type": "string",
      "defaultValue": "none",
      "metadata": {
        "description": "Name of the diagnostics storage account."
      }
    },
    "diagnosticsStorageAccountType": {
      "type": "string",
      "defaultValue": "Standard_LRS",
      "allowedValues": [
        "Standard_LRS",
        "Standard_GRS"
      ],
      "metadata": {
        "description": "Type of diagnostics storage account."
      }
    },
    "diagnosticsStorageAccountKind": {
      "type": "string",
      "defaultValue": "Storage",
      "allowedValues": [
        "Storage",
        "StorageV2"
      ],
      "metadata": {
        "description": "Kind of diagnostics storage account."
      }
    }