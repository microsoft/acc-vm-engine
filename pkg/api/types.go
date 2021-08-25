package api

// VMCategory represents VM category
type VMCategory string

// APIModel complies with the ARM model of
// resource definition in a JSON template.
type APIModel struct {
	VMCategory VMCategory  `json:"vm_category"`
	Properties *Properties `json:"properties,omitempty"`

	VMConfigurator VMConfigurator
}

// OSType represents OS type
type OSType string

// OSName represents pre-set OS name
type OSName string

// OSImageName represents pre-set OS image name
type OSImageName string

// SecurityProfile represents VM security profile
type SecurityProfile struct {
	SecureBoot string `json:"secure_boot_enabled,omitempty"`
	VTPM       string `json:"vtpm_enabled,omitempty"`
}

// VMProfile represents the definition of a VM
type VMProfile struct {
	Name            string           `json:"name"`
	OSImageName     OSImageName      `json:"os_image_name,omitempty"`
	OSType          OSType           `json:"os_type,omitempty"`
	OSName          OSName           `json:"os_name,omitempty"`
	OSDiskType      string           `json:"os_disk_type"`
	OSImage         *OSImage         `json:"os_image,omitempty"`
	DiskSizesGB     []int            `json:"disk_sizes_gb,omitempty"`
	VMSize          string           `json:"vm_size"`
	Ports           []int            `json:"ports,omitempty" validate:"dive,min=1,max=65535"`
	HasDNSName      bool             `json:"has_dns_name"`
	SecurityProfile *SecurityProfile `json:"security_profile,omitempty"`

	TipNodeSessionID string `json:"tip_node_session_id,omitempty"`
	ClusterName      string `json:"cluster_name,omitempty"`
}

// Properties represents the ACS cluster definition
type Properties struct {
	VnetProfile        *VnetProfile        `json:"vnet_rofile"`
	VMProfile          *VMProfile          `json:"vm_profile"`
	LinuxProfile       *LinuxProfile       `json:"linux_profile,omitempty"`
	WindowsProfile     *WindowsProfile     `json:"windows_profile,omitempty"`
	DiagnosticsProfile *DiagnosticsProfile `json:"diagnostics_profile,omitempty"`
}

// OSImage represents OS Image from Azure Image Gallery
type OSImage struct {
	URL       string `json:"url,omitempty"`
	Publisher string `json:"publisher"`
	Offer     string `json:"offer"`
	SKU       string `json:"sku"`
	Version   string `json:"version,omitempty"`
}

// LinuxProfile represents the linux parameters passed to the cluster
type LinuxProfile struct {
	AdminUsername string       `json:"admin_username" validate:"required"`
	AdminPassword string       `json:"admin_password"`
	SSHPubKeys    []*PublicKey `json:"ssh_public_keys"`
}

// WindowsProfile represents the windows parameters passed to the cluster
type WindowsProfile struct {
	AdminUsername string `json:"admin_username" validate:"required"`
	AdminPassword string `json:"admin_password" validate:"required"`
}

// VnetProfile represents the definition of a vnet
type VnetProfile struct {
	VnetResourceGroup string `json:"vnetResourceGroup,omitempty"`
	VnetName          string `json:"vnetName,omitempty"`
	VnetAddress       string `json:"vnetAddress,omitempty"`
	SubnetName        string `json:"subnetName,omitempty"`
	SubnetAddress     string `json:"subnetAddress,omitempty"`
}

// DiagnosticsProfile contains settings to on/off boot diagnostics collection
// in RD Host
type DiagnosticsProfile struct {
	Enabled             bool   `json:"true"`
	StorageAccountName  string `json:"storageAccountName"`
	IsNewStorageAccount bool   `json:"isNewStorageAccount"`
}

// PublicKey contains puvlic SSH key
type PublicKey struct {
	KeyData string `json:"key_data"`
}

// IsCustomVNET returns true if the customer brought their own VNET
func (p *VnetProfile) IsCustomVNET() bool {
	return len(p.VnetResourceGroup) > 0 && len(p.VnetName) > 0 && len(p.SubnetName) > 0
}

// HasAzureGalleryImage returns true if Azure Image Gallery is used
func (h *VMProfile) HasAzureGalleryImage() bool {
	return h.OSImage != nil && len(h.OSImage.Publisher) > 0 && len(h.OSImage.Offer) > 0 && len(h.OSImage.SKU) > 0
}

// HasCustomOsImage returns true if there is a custom OS image url specified
func (h *VMProfile) HasCustomOsImage() bool {
	return h.OSImage != nil && len(h.OSImage.URL) > 0
}

// HasDisks returns true if the customer specified disks
func (h *VMProfile) HasDisks() bool {
	return len(h.DiskSizesGB) > 0
}
