package cluster

import (
	"github.com/huaweicloud/golangsdk"
)


type OutStatus struct {
	Status    string `json:"status"`
}

type MetaData struct {
	Name 	string `json:"name"`
	ID 		string `json:"uuid"`
	SpaceID string `json:"spaceuuid"`
}

type ClusterlistSpec struct {
	ClusterHostList 	ClusterHostList `json:"hostList"`
	AZ 					string 		    `json:"az"`
	CPU 				int 		    `json:"cpu"`
	Memory 				int 		    `json:"memory"`
	VPC				 	string		    `json:"vpc"`
	VpcId 				string		    `json:"vpcid"`
	Subnet 				string		    `json:"subnet"`
	Endpoint 			string		    `json:"endpoint"`
	ExternalEndpoint 	string		    `json:"external_endpoint"`
	SecurityGroupId		string		    `json:"security_group_id"`
	ClusterType 		string		    `json:"clustertype"`
}

type ClusterHostList struct {
	Kind 				string 		 `json:"kind"`
	ApiVersion 			string 		 `json:"apiVersion"`
	Metadata 			MetaData	 `json:"metadata"`
	HostListSpec  		HostListSpec `json:"spec"`
}

type HostListSpec struct {
	HostList []HostList `json:"hostList"`
}

type HostList struct {
	Kind       string  `json:"kind"`
	ApiVersion string  `json:"apiVersion"`
	Metadata   MetaData `json:"metadata"`
	Hostspec   HostSpec `json:"spec"`
	Replicas   int      `json:"replicas"`
	Status     string   `json:"status"`
	Message    string   `json:"message"`
}

type HostSpec struct {

	ClusterUuid string     `json:"clusteruuid"`
	ClusterName string     `json:"clustername"`
	PrivateIp   string     `json:"privateip"`
	PublicIp    string     `json:"publicip"`
	Flavor      string     `json:"flavor"`
	CPU         int        `json:"cpu"`
	Memory      int        `json:"memory"`
	HostId      string     `json:"hostid"`
	AZ          string     `json:"az"`
	Volume      []Volume   `json:"volume"`
	SshKey      string     `json:"sshkey"`
	Status      HostStatus `json:"status"`

}

type Volume struct {
	DiskType 	string  `json:"diskType"`
	DiskSize 	int		`json:"diskSize"`
	VolumeType  string	`json:"volumeType"`
}

type HostStatus struct {
	Capacity 		Capacity		`json:"capacity"`
	Allocatable 	Capacity		`json:"allocatable"`
	Conditions 		[]Conditions	`json:"conditions"`
	Addresses 		[]Addresses		`json:"addresses"`
	DaemonEndpoints DaemonEndpoints `json:"daemonEndpoints"`
	NodeInfo 		NodeInfo		`json:"nodeInfo"`
	Images 			[]Images		`json:"images"`
}

type Capacity struct {
	CPU 	string	`json:"cpu"`
	Memory 	string	`json:"memory"`
	Pods 	string	`json:"pods"`
}
type Conditions struct {
	Type 	string	`json:"type"`
	Status	string	`json:"status"`
	Reason  string	`json:"reason"`
}

type Addresses struct {
	Type	string	`json:"type"`
	Address	string	`json:"address"`
}

type DaemonEndpoints struct {
	KubeletEndpoint KubeletEndpoint	`json:"kubeletEndpoint"`
}

type KubeletEndpoint struct {
	Port 	int	`json:"Port"`
}

type NodeInfo struct {
	MachineID				string	`json:"machineID"`
	SystemUUID				string	`json:"systemUUID"`
	BootID					string	`json:"bootID"`
	KernelVersion			string	`json:"kernelVersion"`
	OsImage					string	`json:"osImage"`
	ContainerRuntimeVersion string	`json:"containerRuntimeVersion"`
	KubeletVersion			string	`json:"kubeletVersion"`
	KubeProxyVersion		string	`json:"kubeProxyVersion"`
}

type Images struct {
	Names		[]string	`json:"names"`
	SizeBytes	int			`json:"sizeBytes"`
}

//RetrievedCluster represents a Neutron Cluster.
//Manage and perform other operations on cluster,
//including querying Clusters as well as
//querying Cluster.
type RetrievedCluster struct {
	Kind string                   `json:"kind"`
	ApiVersion string             `json:"apiVersion"`
	Metadata MetaData             `json:"metadata"`
	Clusterspec ClusterlistSpec   `json:"spec"`
	ClusterStatus OutStatus       `json:"clusterStatus"`
	K8sVersion  string            `json:"k8s_version"`
}

type Certificate struct {
	ClusterID string         `json:"cluster_uuid"`
	ClusterName string       `json:"cluster_name"`
	EndPoint string          `json:"endpoint"`
	ExternalEndpoint string  `json:"external_endpoint"`
	Cacrt string             `json:"cacrt"`
	ClientCrt string         `json:"clientcrt"`
	ClientKey string         `json:"clientkey"`
}


type commonResult struct {
	golangsdk.Result
}

// ListResult represents the result of a get operation. Call its Extract
// method to interpret it as a Cluster.
type ListResult struct {
	commonResult
}

// ExtractCluster is a function that accepts a result and extracts a cluster.
func (r commonResult) ExtractCluster(opts ListOpts) ([]RetrievedCluster, error) {
	var s []RetrievedCluster
	err := r.ExtractInto(&s)
	if err!= nil{
		return nil,err
	}
	return FilterClusters(s,opts)
}


// Extract is a function that accepts a result and extracts a cluster.
func (r commonResult) Extract() (*Certificate, error) {
	var s Certificate
	err := r.ExtractInto(&s)
	return &s, err
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Cluster.
type GetResult struct {
	commonResult
}
