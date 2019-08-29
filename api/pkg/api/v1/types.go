package v1

import (
	"encoding/json"

	"github.com/Masterminds/semver"

	kubermaticv1 "github.com/kubermatic/kubermatic/api/pkg/crd/kubermatic/v1"
	ksemver "github.com/kubermatic/kubermatic/api/pkg/semver"

	cmdv1 "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// LegacyObjectMeta is an object storing common metadata for persistable objects.
// Deprecated: LegacyObjectMeta is deprecated use ObjectMeta instead.
type LegacyObjectMeta struct {
	Name            string `json:"name"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
	UID             string `json:"uid,omitempty"`

	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// ObjectMeta defines the set of fields that objects returned from the API have
// swagger:model ObjectMeta
type ObjectMeta struct {
	// ID unique value that identifies the resource generated by the server. Read-Only.
	ID string `json:"id,omitempty"`

	// Name represents human readable name for the resource
	Name string `json:"name"`

	// DeletionTimestamp is a timestamp representing the server time when this object was deleted.
	// swagger:strfmt date-time
	DeletionTimestamp *Time `json:"deletionTimestamp,omitempty"`

	// CreationTimestamp is a timestamp representing the server time when this object was created.
	// swagger:strfmt date-time
	CreationTimestamp Time `json:"creationTimestamp,omitempty"`
}

// DigitialoceanDatacenterSpec specifies a datacenter of DigitalOcean.
type DigitialoceanDatacenterSpec struct {
	Region string `json:"region"`
}

// HetznerDatacenterSpec specifies a datacenter of Hetzner.
type HetznerDatacenterSpec struct {
	Datacenter string `json:"datacenter"`
	Location   string `json:"location"`
}

// ImageList defines a map of operating system and the image to use
type ImageList map[string]string

// VSphereDatacenterSpec specifies a datacenter of VSphere.
type VSphereDatacenterSpec struct {
	Endpoint   string    `json:"endpoint"`
	Datacenter string    `json:"datacenter"`
	Datastore  string    `json:"datastore"`
	Cluster    string    `json:"cluster"`
	Templates  ImageList `json:"templates"`
}

// BringYourOwnDatacenterSpec specifies a data center with bring-your-own nodes.
type BringYourOwnDatacenterSpec struct{}

// AWSDatacenterSpec specifies a data center of Amazon Web Services.
type AWSDatacenterSpec struct {
	Region string `json:"region"`
}

// AzureDatacenterSpec specifies a datacenter of Azure.
type AzureDatacenterSpec struct {
	Location string `json:"location"`
}

// OpenstackDatacenterSpec specifies a generic bare metal datacenter.
type OpenstackDatacenterSpec struct {
	AvailabilityZone  string    `json:"availability_zone"`
	Region            string    `json:"region"`
	AuthURL           string    `json:"auth_url"`
	Images            ImageList `json:"images"`
	EnforceFloatingIP bool      `json:"enforce_floating_ip"`
}

// PacketDatacenterSpec specifies a datacenter of Packet.
type PacketDatacenterSpec struct {
	Facilities []string `json:"facilities"`
}

// GCPDatacenterSpec specifies a datacenter of GCP.
type GCPDatacenterSpec struct {
	Region       string   `json:"region"`
	ZoneSuffixes []string `json:"zone_suffixes"`
	Regional     bool     `json:"regional"`
}

// DatacenterSpec specifies the data for a datacenter.
type DatacenterSpec struct {
	Seed         string                       `json:"seed"`
	Country      string                       `json:"country,omitempty"`
	Location     string                       `json:"location,omitempty"`
	Provider     string                       `json:"provider,omitempty"`
	Digitalocean *DigitialoceanDatacenterSpec `json:"digitalocean,omitempty"`
	BringYourOwn *BringYourOwnDatacenterSpec  `json:"bringyourown,omitempty"`
	AWS          *AWSDatacenterSpec           `json:"aws,omitempty"`
	Azure        *AzureDatacenterSpec         `json:"azure,omitempty"`
	Openstack    *OpenstackDatacenterSpec     `json:"openstack,omitempty"`
	Packet       *PacketDatacenterSpec        `json:"packet,omitempty"`
	GCP          *GCPDatacenterSpec           `json:"gcp,omitempty"`
	Hetzner      *HetznerDatacenterSpec       `json:"hetzner,omitempty"`
	VSphere      *VSphereDatacenterSpec       `json:"vsphere,omitempty"`
}

// DatacenterList represents a list of datacenters
// swagger:model DatacenterList
type DatacenterList []Datacenter

// Datacenter is the object representing a Kubernetes infra datacenter.
// swagger:model Datacenter
type Datacenter struct {
	Metadata LegacyObjectMeta `json:"metadata"`
	Spec     DatacenterSpec   `json:"spec"`
	Seed     bool             `json:"seed,omitempty"`
}

// AWSSize represents a object of AWS size.
// swagger:model AWSSize
type AWSSize struct {
	Name       string  `json:"name"`
	PrettyName string  `json:"pretty_name"`
	Memory     float32 `json:"memory"`
	VCPUs      int     `json:"vcpus"`
	Price      float64 `json:"price"`
}

// AWSSizeList represents an array of AWS sizes.
// swagger:model AWSSizeList
type AWSSizeList []AWSSize

// AWSZone represents a object of AWS availability zone.
// swagger:model AWSZone
type AWSZone struct {
	Name string `json:"name"`
}

// AWSZoneList represents an array of AWS availability zones.
// swagger:model AWSZoneList
type AWSZoneList []AWSZone

// AWSSubnetList represents an array of AWS availability subnets.
// swagger:model AWSSubnetList
type AWSSubnetList []AWSSubnet

// AWSSubnet represents a object of AWS availability subnet.
// swagger:model AWSSubnet
type AWSSubnet struct {
	Name                    string   `json:"name"`
	ID                      string   `json:"id"`
	AvailabilityZone        string   `json:"availability_zone"`
	AvailabilityZoneID      string   `json:"availability_zone_id"`
	IPv4CIDR                string   `json:"ipv4cidr"`
	IPv6CIDR                string   `json:"ipv6cidr"`
	Tags                    []AWSTag `json:"tags,omitempty"`
	State                   string   `json:"state"`
	AvailableIPAddressCount int64    `json:"available_ip_address_count"`
	DefaultForAz            bool     `json:"default"`
}

// AWSTag represents a object of AWS tags.
// swagger:model AWSTag
type AWSTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AWSVPCList represents an array of AWS VPC's.
// swagger:model AWSVPCList
type AWSVPCList []AWSVPC

// AWSVPC represents a object of AWS VPC.
// swagger:model AWSVPC
type AWSVPC struct {
	// The primary IPv4 CIDR block for the VPC.
	CidrBlock string `json:"cidrBlock"`

	// Information about the IPv4 CIDR blocks associated with the VPC.
	CidrBlockAssociationSet []AWSVpcCidrBlockAssociation `json:"cidrBlockAssociationSet,omitempty"`

	// The ID of the set of DHCP options you've associated with the VPC (or default
	// if the default options are associated with the VPC).
	DhcpOptionsID string `json:"dhcpOptionsId"`

	// The allowed tenancy of instances launched into the VPC.
	InstanceTenancy string `json:"instanceTenancy"`

	// Information about the IPv6 CIDR blocks associated with the VPC.
	Ipv6CidrBlockAssociationSet []AWSVpcIpv6CidrBlockAssociation `json:"ipv6CidrBlockAssociationSet,omitempty"`

	// Indicates whether the VPC is the default VPC.
	IsDefault bool `json:"isDefault"`

	// The ID of the AWS account that owns the VPC.
	OwnerID string `json:"ownerId"`

	// The current state of the VPC.
	State string `json:"state"`

	// Any tags assigned to the VPC.
	Tags []AWSTag `json:"tags,omitempty"`

	Name string `json:"name"`

	// The ID of the VPC.
	VpcID string `json:"vpcId"`
}

// AWSVpcCidrBlockAssociation describes an IPv4 CIDR block associated with a VPC.
// swagger:model AWSVpcCidrBlockAssociation
type AWSVpcCidrBlockAssociation struct {
	// The association ID for the IPv4 CIDR block.
	AssociationID string `json:"associationId"`

	// The IPv4 CIDR block.
	CidrBlock string `json:"cidrBlock"`

	// The state of the CIDR block.
	State string `json:"state"`

	// A message about the status of the CIDR block, if applicable.
	StatusMessage string `json:"statusMessage"`
}

// AWSVpcIpv6CidrBlockAssociation describes an IPv6 CIDR block associated with a VPC.
// swagger:model AWSVpcIpv6CidrBlockAssociation
type AWSVpcIpv6CidrBlockAssociation struct {
	AWSVpcCidrBlockAssociation
}

// GCPDiskTypeList represents an array of GCP disk types.
// swagger:model GCPDiskTypeList
type GCPDiskTypeList []GCPDiskType

// GCPDiskType represents a object of GCP disk type.
// swagger:model GCPDiskType
type GCPDiskType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GCPMachineSizeList represents an array of GCP machine sizes.
// swagger:model GCPMachineSizeList
type GCPMachineSizeList []GCPMachineSize

// GCPMachineSize represents a object of GCP machine size.
// swagger:model GCPMachineSize
type GCPMachineSize struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Memory      int64  `json:"memory"`
	VCPUs       int64  `json:"vcpus"`
}

// GCPZone represents a object of GCP zone.
// swagger:model GCPZone
type GCPZone struct {
	Name string `json:"name"`
}

// GCPZoneList represents an array of GCP zones.
// swagger:model GCPZoneList
type GCPZoneList []GCPZone

// DigitaloceanSizeList represents a object of digitalocean sizes.
// swagger:model DigitaloceanSizeList
type DigitaloceanSizeList struct {
	Standard  []DigitaloceanSize `json:"standard"`
	Optimized []DigitaloceanSize `json:"optimized"`
}

// CredentialList represents a object for provider credential names.
// swagger:model CredentialList
type CredentialList struct {
	Names []string `json:"names,omitempty"`
}

// DigitaloceanSize is the object representing digitalocean sizes.
// swagger:model DigitaloceanSize
type DigitaloceanSize struct {
	Slug         string   `json:"slug"`
	Available    bool     `json:"available"`
	Transfer     float64  `json:"transfer"`
	PriceMonthly float64  `json:"price_monthly"`
	PriceHourly  float64  `json:"price_hourly"`
	Memory       int      `json:"memory"`
	VCPUs        int      `json:"vcpus"`
	Disk         int      `json:"disk"`
	Regions      []string `json:"regions"`
}

// AzureSizeList represents an array of Azure VM sizes.
// swagger:model AzureSizeList
type AzureSizeList []AzureSize

// AzureSize is the object representing Azure VM sizes.
// swagger:model AzureSize
type AzureSize struct {
	Name                 string `json:"name"`
	NumberOfCores        int32  `json:"numberOfCores"`
	OsDiskSizeInMB       int32  `json:"osDiskSizeInMB"`
	ResourceDiskSizeInMB int32  `json:"resourceDiskSizeInMB"`
	MemoryInMB           int32  `json:"memoryInMB"`
	MaxDataDiskCount     int32  `json:"maxDataDiskCount"`
}

// HetznerSizeList represents an array of Hetzner sizes.
// swagger:model HetznerSizeList
type HetznerSizeList struct {
	Standard  []HetznerSize `json:"standard"`
	Dedicated []HetznerSize `json:"dedicated"`
}

// HetznerSize is the object representing Hetzner sizes.
// swagger:model HetznerSize
type HetznerSize struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Cores       int     `json:"cores"`
	Memory      float32 `json:"memory"`
	Disk        int     `json:"disk"`
}

// PacketSizeList represents an array of Packet VM sizes.
// swagger:model PacketSizeList
type PacketSizeList []PacketSize

// PacketSize is the object representing Packet VM sizes.
// swagger:model PacketSize
type PacketSize struct {
	Name   string        `json:"name,omitempty"`
	CPUs   []PacketCPU   `json:"cpus,omitempty"`
	Memory string        `json:"memory,omitempty"`
	Drives []PacketDrive `json:"drives,omitempty"`
}

// PacketCPU represents an array of Packet CPUs. It is a part of PacketSize.
// swagger:model PacketCPU
type PacketCPU struct {
	Count int    `json:"count,omitempty"`
	Type  string `json:"type,omitempty"`
}

// PacketDrive represents an array of Packet drives. It is a part of PacketSize.
// swagger:model PacketDrive
type PacketDrive struct {
	Count int    `json:"count,omitempty"`
	Size  string `json:"size,omitempty"`
	Type  string `json:"type,omitempty"`
}

// SSHKey represents a ssh key
// swagger:model SSHKey
type SSHKey struct {
	ObjectMeta
	Spec SSHKeySpec `json:"spec"`
}

// SSHKeySpec represents the details of a ssh key
type SSHKeySpec struct {
	Fingerprint string `json:"fingerprint"`
	PublicKey   string `json:"publicKey"`
}

// User represent an API user
// swagger:model User
type User struct {
	ObjectMeta

	// Email an email address of the user
	Email string `json:"email"`

	// Projects holds the list of project the user belongs to
	// along with the group names
	Projects []ProjectGroup `json:"projects,omitempty"`
}

// ProjectGroup is a helper data structure that
// stores the information about a project and a group prefix that a user belongs to
type ProjectGroup struct {
	ID          string `json:"id"`
	GroupPrefix string `json:"group"`
}

// These are the valid statuses of a ServiceAccount.
const (
	// ServiceAccountActive means the ServiceAccount is available for use in the system
	ServiceAccountActive string = "Active"

	// ServiceAccountInactive means the ServiceAccount is inactive and requires further initialization
	ServiceAccountInactive string = "Inactive"

	// ServiceAccountTerminating means the ServiceAccount is undergoing graceful termination
	ServiceAccountTerminating string = "Terminating"
)

// ServiceAccount represent an API service account
// swagger:model ServiceAccount
type ServiceAccount struct {
	ObjectMeta
	// Status describes three stages of ServiceAccount life including Active, Inactive and Terminating
	Status string `json:"status"`
	// Group that a service account belongs to
	Group string `json:"group"`
}

// PublicServiceAccountToken represent an API service account token without secret fields
// swagger:model PublicServiceAccountToken
type PublicServiceAccountToken struct {
	ObjectMeta
	// Expiry is a timestamp representing the time when this token will expire.
	// swagger:strfmt date-time
	Expiry Time `json:"expiry,omitempty"`
}

// ServiceAccountToken represent an API service account token
// swagger:model ServiceAccountToken
type ServiceAccountToken struct {
	PublicServiceAccountToken
	// Token the JWT token
	Token string `json:"token,omitempty"`
}

// Project is a top-level container for a set of resources
// swagger:model Project
type Project struct {
	ObjectMeta
	Status string `json:"status"`
	// Owners an optional owners list for the given project
	Owners []User `json:"owners,omitempty"`
}

// Kubeconfig is a clusters kubeconfig
// swagger:model Kubeconfig
type Kubeconfig struct {
	cmdv1.Config
}

// OpenstackSize is the object representing openstack's sizes.
// swagger:model OpenstackSize
type OpenstackSize struct {
	// Slug holds  the name of the size
	Slug string `json:"slug"`
	// MemoryTotalBytes is the amount of memory, measured in MB
	Memory int `json:"memory"`
	// VCPUs indicates how many (virtual) CPUs are available for this flavor
	VCPUs int `json:"vcpus"`
	// Disk is the amount of root disk, measured in GB
	Disk int `json:"disk"`
	// Swap is the amount of swap space, measured in MB
	Swap int `json:"swap"`
	// Region specifies the geographic region in which the size resides
	Region string `json:"region"`
	// IsPublic indicates whether the size is public (available to all projects) or scoped to a set of projects
	IsPublic bool `json:"isPublic"`
}

// OpenstackSubnet is the object representing a openstack subnet.
// swagger:model OpenstackSubnet
type OpenstackSubnet struct {
	// Id uniquely identifies the subnet
	ID string `json:"id"`
	// Name is human-readable name for the subnet
	Name string `json:"name"`
}

// OpenstackTenant is the object representing a openstack tenant.
// swagger:model OpenstackTenant
type OpenstackTenant struct {
	// Id uniquely identifies the current tenant
	ID string `json:"id"`
	// Name is the name of the tenant
	Name string `json:"name"`
}

// OpenstackNetwork is the object representing a openstack network.
// swagger:model OpenstackNetwork
type OpenstackNetwork struct {
	// Id uniquely identifies the current network
	ID string `json:"id"`
	// Name is the name of the network
	Name string `json:"name"`
	// External set if network is the external network
	External bool `json:"external"`
}

// OpenstackSecurityGroup is the object representing a openstack security group.
// swagger:model OpenstackSecurityGroup
type OpenstackSecurityGroup struct {
	// Id uniquely identifies the current security group
	ID string `json:"id"`
	// Name is the name of the security group
	Name string `json:"name"`
}

// VSphereNetwork is the object representing a vsphere network.
// swagger:model VSphereNetwork
type VSphereNetwork struct {
	// Name is the name of the network
	Name string `json:"name"`
	// AbsolutePath is the absolute path inside vCenter
	AbsolutePath string `json:"absolutePath"`
	// RelativePath is the relative path inside the datacenter
	RelativePath string `json:"relativePath"`
	// Type defines the type of network
	Type string `json:"type"`
}

// VSphereFolder is the object representing a vsphere folder.
// swagger:model VSphereFolder
type VSphereFolder struct {
	// Path is the path of the folder
	Path string `json:"path"`
}

// MasterVersion describes a version of the master components
// swagger:model MasterVersion
type MasterVersion struct {
	Version *semver.Version `json:"version"`
	Default bool            `json:"default,omitempty"`

	// If true, then given version control plane version is not compatible
	// with one of the kubelets inside cluster and shouldn't be used.
	RestrictedByKubeletVersion bool `json:"restrictedByKubeletVersion,omitempty"`
}

// CreateClusterSpec is the structure that is used to create cluster with its initial node deployment
// swagger:model CreateClusterSpec
type CreateClusterSpec struct {
	Cluster        Cluster         `json:"cluster"`
	NodeDeployment *NodeDeployment `json:"nodeDeployment,omitempty"`
}

const (
	// OpenShiftClusterType defines the OpenShift cluster type
	OpenShiftClusterType string = "openshift"
	// KubernetesClusterType defines the Kubernetes cluster type
	KubernetesClusterType string = "kubernetes"
)

// Cluster defines the cluster resource
//
// Note:
// Cluster has a custom MarshalJSON method defined
// and thus the output may vary
//
// swagger:model Cluster
type Cluster struct {
	ObjectMeta `json:",inline"`
	Type       string        `json:"type"`
	Credential string        `json:"credential,omitempty"`
	Spec       ClusterSpec   `json:"spec"`
	Status     ClusterStatus `json:"status"`
}

// ClusterSpec defines the cluster specification
type ClusterSpec struct {
	// Cloud specifies the cloud providers configuration
	Cloud kubermaticv1.CloudSpec `json:"cloud"`

	// MachineNetworks optionally specifies the parameters for IPAM.
	MachineNetworks []kubermaticv1.MachineNetworkingConfig `json:"machineNetworks,omitempty"`

	// Version desired version of the kubernetes master components
	Version ksemver.Semver `json:"version"`

	// OIDC settings
	OIDC kubermaticv1.OIDCSettings `json:"oidc,omitempty"`

	// If active the PodSecurityPolicy admission plugin is configured at the apiserver
	UsePodSecurityPolicyAdmissionPlugin bool `json:"usePodSecurityPolicyAdmissionPlugin,omitempty"`

	// AuditLogging
	AuditLogging *kubermaticv1.AuditLoggingSettings `json:"auditLogging,omitempty"`
}

// MarshalJSON marshals ClusterSpec object into JSON. It is overwritten to control data
// that will be returned in the API responses (see: PublicCloudSpec struct).
func (cs *ClusterSpec) MarshalJSON() ([]byte, error) {
	ret, err := json.Marshal(struct {
		Cloud                               PublicCloudSpec                        `json:"cloud"`
		MachineNetworks                     []kubermaticv1.MachineNetworkingConfig `json:"machineNetworks,omitempty"`
		Version                             ksemver.Semver                         `json:"version"`
		OIDC                                kubermaticv1.OIDCSettings              `json:"oidc"`
		UsePodSecurityPolicyAdmissionPlugin bool                                   `json:"usePodSecurityPolicyAdmissionPlugin,omitempty"`
		AuditLogging                        *kubermaticv1.AuditLoggingSettings     `json:"auditLogging,omitempty"`
	}{
		Cloud: PublicCloudSpec{
			DatacenterName: cs.Cloud.DatacenterName,
			Fake:           newPublicFakeCloudSpec(cs.Cloud.Fake),
			Digitalocean:   newPublicDigitaloceanCloudSpec(cs.Cloud.Digitalocean),
			BringYourOwn:   newPublicBringYourOwnCloudSpec(cs.Cloud.BringYourOwn),
			AWS:            newPublicAWSCloudSpec(cs.Cloud.AWS),
			Azure:          newPublicAzureCloudSpec(cs.Cloud.Azure),
			Openstack:      newPublicOpenstackCloudSpec(cs.Cloud.Openstack),
			Packet:         newPublicPacketCloudSpec(cs.Cloud.Packet),
			Hetzner:        newPublicHetznerCloudSpec(cs.Cloud.Hetzner),
			VSphere:        newPublicVSphereCloudSpec(cs.Cloud.VSphere),
			GCP:            newPublicGCPCloudSpec(cs.Cloud.GCP),
		},
		Version:                             cs.Version,
		MachineNetworks:                     cs.MachineNetworks,
		OIDC:                                cs.OIDC,
		UsePodSecurityPolicyAdmissionPlugin: cs.UsePodSecurityPolicyAdmissionPlugin,
		AuditLogging:                        cs.AuditLogging,
	})

	return ret, err
}

// PublicCloudSpec is a public counterpart of apiv1.CloudSpec.
type PublicCloudSpec struct {
	DatacenterName string                       `json:"dc"`
	Fake           *PublicFakeCloudSpec         `json:"fake,omitempty"`
	Digitalocean   *PublicDigitaloceanCloudSpec `json:"digitalocean,omitempty"`
	BringYourOwn   *PublicBringYourOwnCloudSpec `json:"bringyourown,omitempty"`
	AWS            *PublicAWSCloudSpec          `json:"aws,omitempty"`
	Azure          *PublicAzureCloudSpec        `json:"azure,omitempty"`
	Openstack      *PublicOpenstackCloudSpec    `json:"openstack,omitempty"`
	Packet         *PublicPacketCloudSpec       `json:"packet,omitempty"`
	Hetzner        *PublicHetznerCloudSpec      `json:"hetzner,omitempty"`
	VSphere        *PublicVSphereCloudSpec      `json:"vsphere,omitempty"`
	GCP            *PublicGCPCloudSpec          `json:"gcp,omitempty"`
}

// PublicFakeCloudSpec is a public counterpart of apiv1.FakeCloudSpec.
type PublicFakeCloudSpec struct{}

func newPublicFakeCloudSpec(internal *kubermaticv1.FakeCloudSpec) (public *PublicFakeCloudSpec) {
	if internal == nil {
		return nil
	}

	return &PublicFakeCloudSpec{}
}

// PublicDigitaloceanCloudSpec is a public counterpart of apiv1.DigitaloceanCloudSpec.
type PublicDigitaloceanCloudSpec struct{}

func newPublicDigitaloceanCloudSpec(internal *kubermaticv1.DigitaloceanCloudSpec) (public *PublicDigitaloceanCloudSpec) {
	if internal == nil {
		return nil
	}

	return &PublicDigitaloceanCloudSpec{}
}

// PublicHetznerCloudSpec is a public counterpart of apiv1.HetznerCloudSpec.
type PublicHetznerCloudSpec struct{}

func newPublicHetznerCloudSpec(internal *kubermaticv1.HetznerCloudSpec) (public *PublicHetznerCloudSpec) {
	if internal == nil {
		return nil
	}

	return &PublicHetznerCloudSpec{}
}

// PublicAzureCloudSpec is a public counterpart of apiv1.AzureCloudSpec.
type PublicAzureCloudSpec struct{}

func newPublicAzureCloudSpec(internal *kubermaticv1.AzureCloudSpec) (public *PublicAzureCloudSpec) {
	if internal == nil {
		return nil
	}

	return &PublicAzureCloudSpec{}
}

// PublicVSphereCloudSpec is a public counterpart of apiv1.VSphereCloudSpec.
type PublicVSphereCloudSpec struct{}

func newPublicVSphereCloudSpec(internal *kubermaticv1.VSphereCloudSpec) (public *PublicVSphereCloudSpec) {
	if internal == nil {
		return nil
	}

	return &PublicVSphereCloudSpec{}
}

// PublicBringYourOwnCloudSpec is a public counterpart of apiv1.BringYourOwnCloudSpec.
type PublicBringYourOwnCloudSpec struct{}

func newPublicBringYourOwnCloudSpec(internal *kubermaticv1.BringYourOwnCloudSpec) (public *PublicBringYourOwnCloudSpec) {
	if internal == nil {
		return nil
	}

	return &PublicBringYourOwnCloudSpec{}
}

// PublicAWSCloudSpec is a public counterpart of apiv1.AWSCloudSpec.
type PublicAWSCloudSpec struct{}

func newPublicAWSCloudSpec(internal *kubermaticv1.AWSCloudSpec) (public *PublicAWSCloudSpec) {
	if internal == nil {
		return nil
	}

	return &PublicAWSCloudSpec{}
}

// PublicOpenstackCloudSpec is a public counterpart of apiv1.OpenstackCloudSpec.
type PublicOpenstackCloudSpec struct {
	FloatingIPPool string `json:"floatingIpPool"`
}

func newPublicOpenstackCloudSpec(internal *kubermaticv1.OpenstackCloudSpec) (public *PublicOpenstackCloudSpec) {
	if internal == nil {
		return nil
	}

	return &PublicOpenstackCloudSpec{
		FloatingIPPool: internal.FloatingIPPool,
	}
}

// PublicPacketCloudSpec is a public counterpart of apiv1.PacketCloudSpec.
type PublicPacketCloudSpec struct{}

func newPublicPacketCloudSpec(internal *kubermaticv1.PacketCloudSpec) (public *PublicPacketCloudSpec) {
	if internal == nil {
		return nil
	}

	return &PublicPacketCloudSpec{}
}

// PublicGCPCloudSpec is a public counterpart of apiv1.GCPCloudSpec.
type PublicGCPCloudSpec struct{}

func newPublicGCPCloudSpec(internal *kubermaticv1.GCPCloudSpec) (public *PublicGCPCloudSpec) {
	if internal == nil {
		return nil
	}

	return &PublicGCPCloudSpec{}
}

// ClusterStatus defines the cluster status
type ClusterStatus struct {
	// Version actual version of the kubernetes master components
	Version ksemver.Semver `json:"version"`

	// URL specifies the address at which the cluster is available
	URL string `json:"url"`
}

// ClusterHealth stores health information about the cluster's components.
// swagger:model ClusterHealth
type ClusterHealth struct {
	Apiserver                    kubermaticv1.HealthStatus `json:"apiserver"`
	Scheduler                    kubermaticv1.HealthStatus `json:"scheduler"`
	Controller                   kubermaticv1.HealthStatus `json:"controller"`
	MachineController            kubermaticv1.HealthStatus `json:"machineController"`
	Etcd                         kubermaticv1.HealthStatus `json:"etcd"`
	CloudProviderInfrastructure  kubermaticv1.HealthStatus `json:"cloudProviderInfrastructure"`
	UserClusterControllerManager kubermaticv1.HealthStatus `json:"userClusterControllerManager"`
}

// Addon represents a predefined addon that users may install into their cluster
// swagger:model Addon
type Addon struct {
	ObjectMeta `json:",inline"`

	Spec AddonSpec `json:"spec"`
}

// AddonSpec addon specification
// swagger:model AddonSpec
type AddonSpec struct {
	// Variables is free form data to use for parsing the manifest templates
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// ClusterList represents a list of clusters
// swagger:model ClusterList
type ClusterList []Cluster

// Node represents a worker node that is part of a cluster
// swagger:model Node
type Node struct {
	ObjectMeta `json:",inline"`
	Spec       NodeSpec   `json:"spec"`
	Status     NodeStatus `json:"status"`
}

// NodeCloudSpec represents the collection of cloud provider specific settings. Only one must be set at a time.
// swagger:model NodeCloudSpec
type NodeCloudSpec struct {
	Digitalocean *DigitaloceanNodeSpec `json:"digitalocean,omitempty"`
	AWS          *AWSNodeSpec          `json:"aws,omitempty"`
	Azure        *AzureNodeSpec        `json:"azure,omitempty"`
	Openstack    *OpenstackNodeSpec    `json:"openstack,omitempty"`
	Packet       *PacketNodeSpec       `json:"packet,omitempty"`
	Hetzner      *HetznerNodeSpec      `json:"hetzner,omitempty"`
	VSphere      *VSphereNodeSpec      `json:"vsphere,omitempty"`
	GCP          *GCPNodeSpec          `json:"gcp,omitempty"`
}

// UbuntuSpec ubuntu specific settings
// swagger:model UbuntuSpec
type UbuntuSpec struct {
	// do a dist-upgrade on boot and reboot it required afterwards
	DistUpgradeOnBoot bool `json:"distUpgradeOnBoot"`
}

// CentOSSpec contains CentOS specific settings
type CentOSSpec struct {
	// do a dist-upgrade on boot and reboot it required afterwards
	DistUpgradeOnBoot bool `json:"distUpgradeOnBoot"`
}

// ContainerLinuxSpec ubuntu linux specific settings
// swagger:model ContainerLinuxSpec
type ContainerLinuxSpec struct {
	// disable container linux auto-update feature
	DisableAutoUpdate bool `json:"disableAutoUpdate"`
}

// OperatingSystemSpec represents the collection of os specific settings. Only one must be set at a time.
// swagger:model OperatingSystemSpec
type OperatingSystemSpec struct {
	Ubuntu         *UbuntuSpec         `json:"ubuntu,omitempty"`
	ContainerLinux *ContainerLinuxSpec `json:"containerLinux,omitempty"`
	CentOS         *CentOSSpec         `json:"centos,omitempty"`
}

// NodeVersionInfo node version information
// swagger:model NodeVersionInfo
type NodeVersionInfo struct {
	Kubelet string `json:"kubelet"`
}

// TaintSpec defines a node taint
type TaintSpec struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Effect string `json:"effect"`
}

// NodeSpec node specification
// swagger:model NodeSpec
type NodeSpec struct {
	// required: true
	Cloud NodeCloudSpec `json:"cloud"`
	// required: true
	OperatingSystem OperatingSystemSpec `json:"operatingSystem"`
	// required: false
	SSHUserName string `json:"sshUserName,omitempty"`
	// required: true
	Versions NodeVersionInfo `json:"versions,omitempty"`
	// Map of string keys and values that can be used to organize and categorize (scope and select) objects.
	// It will be applied to Nodes allowing users run their apps on specific Node using labelSelector.
	// required: false
	Labels map[string]string `json:"labels,omitempty"`
	// List of taints to set on new nodes
	Taints []TaintSpec `json:"taints,omitempty"`
}

// DigitaloceanNodeSpec digitalocean node settings
// swagger:model DigitaloceanNodeSpec
type DigitaloceanNodeSpec struct {
	// droplet size slug
	// required: true
	Size string `json:"size"`
	// enable backups for the droplet
	Backups bool `json:"backups"`
	// enable ipv6 for the droplet
	IPv6 bool `json:"ipv6"`
	// enable monitoring for the droplet
	Monitoring bool `json:"monitoring"`
	// additional droplet tags
	Tags []string `json:"tags"`
}

// HetznerNodeSpec Hetzner node settings
// swagger:model HetznerNodeSpec
type HetznerNodeSpec struct {
	// server type
	// required: true
	Type string `json:"type"`
}

// AzureNodeSpec describes settings for an Azure node
// swagger:model AzureNodeSpec
type AzureNodeSpec struct {
	// VM size
	// required: true
	Size string `json:"size"`
	// should the machine have a publicly accessible IP address
	// required: false
	AssignPublicIP bool `json:"assignPublicIP"`
	// Additional metadata to set
	// required: false
	Tags map[string]string `json:"tags,omitempty"`
}

// VSphereNodeSpec VSphere node settings
// swagger:model VSphereNodeSpec
type VSphereNodeSpec struct {
	CPUs       int    `json:"cpus"`
	Memory     int    `json:"memory"`
	DiskSizeGB *int64 `json:"diskSizeGB,omitempty"`
	Template   string `json:"template"`
}

// OpenstackNodeSpec openstack node settings
// swagger:model OpenstackNodeSpec
type OpenstackNodeSpec struct {
	// instance flavor
	// required: true
	Flavor string `json:"flavor"`
	// image to use
	// required: true
	Image string `json:"image"`
	// Additional metadata to set
	// required: false
	Tags map[string]string `json:"tags,omitempty"`
	// Defines whether floating ip should be used
	// required: false
	UseFloatingIP bool `json:"useFloatingIP,omitempty"`
	// if set, the rootDisk will be a volume. If not, the rootDisk will be on ephemeral storage and its size will be derived from the flavor
	// required: false
	RootDiskSizeGB *int `json:"diskSize"`
}

// AWSNodeSpec aws specific node settings
// swagger:model AWSNodeSpec
type AWSNodeSpec struct {
	// instance type. for example: t2.micro
	// required: true
	InstanceType string `json:"instanceType"`
	// size of the volume in gb. Only one volume will be created
	// required: true
	VolumeSize int64 `json:"diskSize"`
	// type of the volume. for example: gp2, io1, st1, sc1, standard
	// required: true
	VolumeType string `json:"volumeType"`
	// ami to use. Will be defaulted to a ami for your selected operating system and region. Only set this when you know what you do.
	AMI string `json:"ami"`
	// additional instance tags
	Tags map[string]string `json:"tags"`
	// Availiability zone in which to place the node. It is coupled with the subnet to which the node will belong.
	AvailabilityZone string `json:"availabilityZone"`
	// The VPC subnet to which the node shall be connected.
	SubnetID string `json:"subnetID"`
}

// PacketNodeSpec specifies packet specific node settings
// swagger:model PacketNodeSpec
type PacketNodeSpec struct {
	// InstanceType denotes the plan to which the device will be provisioned.
	// required: true
	InstanceType string `json:"instanceType"`
	// additional instance tags
	// required: false
	Tags []string `json:"tags"`
}

// GCPNodeSpec gcp specific node settings
// swagger:model GCPNodeSpec
type GCPNodeSpec struct {
	Zone        string            `json:"zone"`
	MachineType string            `json:"machineType"`
	DiskSize    int64             `json:"diskSize"`
	DiskType    string            `json:"diskType"`
	Preemptible bool              `json:"preemptible"`
	Labels      map[string]string `json:"labels"`
	Tags        []string          `json:"tags"`
}

// NodeResources cpu and memory of a node
// swagger:model NodeResources
type NodeResources struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// NodeStatus is information about the current status of a node.
// swagger:model NodeStatus
type NodeStatus struct {
	// name of the actual Machine object
	MachineName string `json:"machineName"`
	// resources in total
	Capacity NodeResources `json:"capacity,omitempty"`
	// allocatable resources
	Allocatable NodeResources `json:"allocatable,omitempty"`
	// different addresses of a node
	Addresses []NodeAddress `json:"addresses,omitempty"`
	// node versions and systems info
	NodeInfo NodeSystemInfo `json:"nodeInfo,omitempty"`

	// in case of a error this will contain a short error message
	ErrorReason string `json:"errorReason,omitempty"`
	// in case of a error this will contain a detailed error explanation
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// NodeAddress contains information for the node's address.
// swagger:model NodeAddress
type NodeAddress struct {
	// address type. for example: ExternalIP, InternalIP, InternalDNS, ExternalDNS
	Type string `json:"type"`
	// the actual address. for example: 192.168.1.1, node1.my.dns
	Address string `json:"address"`
}

// NodeSystemInfo is a set of versions/ids/uuids to uniquely identify the node.
// swagger:model NodeSystemInfo
type NodeSystemInfo struct {
	KernelVersion           string `json:"kernelVersion"`
	ContainerRuntime        string `json:"containerRuntime"`
	ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
	KubeletVersion          string `json:"kubeletVersion"`
	OperatingSystem         string `json:"operatingSystem"`
	Architecture            string `json:"architecture"`
}

// ClusterMetrics defines a metric for the given cluster
// swagger:model ClusterMetrics
type ClusterMetrics struct {
	Name                string              `json:"name"`
	ControlPlaneMetrics ControlPlaneMetrics `json:"controlPlane"`
}

// ControlPlaneMetrics defines a metric for the user cluster control plane resources
// swagger:model ClusterMetrics
type ControlPlaneMetrics struct {
	// MemoryTotalBytes in bytes
	MemoryTotalBytes int64 `json:"memoryTotalBytes,omitempty"`
	// CPUTotalMillicores in m cores
	CPUTotalMillicores int64 `json:"cpuTotalMillicores,omitempty"`
}

// NodeMetric defines a metric for the given node
// swagger:model NodeMetric
type NodeMetric struct {
	Name string `json:"name"`
	// MemoryTotalBytes current memory usage in bytes
	MemoryTotalBytes int64 `json:"memoryTotalBytes,omitempty"`
	// MemoryAvailableBytes available memory for node
	MemoryAvailableBytes int64 `json:"memoryAvailableBytes,omitempty"`
	// MemoryUsedPercentage in percentage
	MemoryUsedPercentage int64 `json:"memoryUsedPercentage,omitempty"`
	// CPUTotalMillicores in m cores
	CPUTotalMillicores     int64 `json:"cpuTotalMillicores,omitempty"`
	CPUAvailableMillicores int64 `json:"cpuAvailableMillicores,omitempty"`
	// CPUUsedPercentage in percentage
	CPUUsedPercentage int64 `json:"cpuUsedPercentage,omitempty"`
}

// NodeDeployment represents a set of worker nodes that is part of a cluster
// swagger:model NodeDeployment
type NodeDeployment struct {
	ObjectMeta `json:",inline"`

	Spec   NodeDeploymentSpec               `json:"spec"`
	Status v1alpha1.MachineDeploymentStatus `json:"status"`
}

// NodeDeploymentSpec node deployment specification
// swagger:model NodeDeploymentSpec
type NodeDeploymentSpec struct {
	// required: true
	Replicas int32 `json:"replicas,omitempty"`
	// required: true
	Template NodeSpec `json:"template"`
	// required: false
	Paused *bool `json:"paused,omitempty"`
}

// Event is a report of an event somewhere in the cluster.
type Event struct {
	ObjectMeta `json:",inline"`

	// A human-readable description of the status of this operation.
	Message string `json:"message,omitempty"`

	// Type of this event (i.e. normal or warning). New types could be added in the future.
	Type string `json:"type,omitempty"`

	// The object reference that those events are about.
	InvolvedObject ObjectReference `json:"involvedObject"`

	// The time at which the most recent occurrence of this event was recorded.
	// swagger:strfmt date-time
	LastTimestamp Time `json:"lastTimestamp,omitempty"`

	// The number of times this event has occurred.
	Count int32 `json:"count,omitempty"`
}

// ObjectReference contains basic information about referred object.
type ObjectReference struct {
	// Type of the referent.
	Type string `json:"type,omitempty"`
	// Namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	Name string `json:"name,omitempty"`
}

// KubermaticVersions describes the versions of running Kubermatic components.
// swagger:model KubermaticVersions
type KubermaticVersions struct {
	// Version of the Kubermatic API server.
	API string `json:"api"`
}

const (
	// NodeDeletionFinalizer indicates that the nodes still need cleanup
	NodeDeletionFinalizer = "kubermatic.io/delete-nodes"
	// InClusterPVCleanupFinalizer indicates that the PVs still need cleanup
	InClusterPVCleanupFinalizer = "kubermatic.io/cleanup-in-cluster-pv"
	// InClusterLBCleanupFinalizer indicates that the LBs still need cleanup
	InClusterLBCleanupFinalizer = "kubermatic.io/cleanup-in-cluster-lb"
	// CredentialsSecretsCleanupFinalizer indicates that secrets for credentials still need cleanup
	CredentialsSecretsCleanupFinalizer = "kubermatic.io/cleanup-credentials-secrets"
)
