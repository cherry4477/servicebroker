package api

import "fmt"

type dcosOption struct {
	auth      Auth
	host      string
	Acs_token string `json:"token"`
}

type authConfig string

func (c *dcosOption) validate() error {
	if len(c.host) == 0 {
		return fmt.Errorf("application config host must not be nil.")
	}

	if len(c.Acs_token) == 0 {
		return fmt.Errorf("application config acs_token must not be nil.")
	}

	return nil
}

//Create and start a new application.
//Note: This operation will create a deployment. The operation finishes, if the deployment succeeds.
//You can query the deployments endoint to see the status of the deployment.
//{
//	"id": "/mysql",
//	"cpus": 1.5,
//	"mem": 2000,
//	"disk": 10000,
//	"instances": 1,
//	"container": {
//		"type": "DOCKER",
//		"volumes": [
//			{
//				"containerPath": "mysql_data",
//				"mode": "RW",
//				"persistent": {
//					"size": 100
//				}
//			},
//			{
//				"containerPath": "/var/lib/mysql",
//				"hostPath": "mysql_data",
//				"mode": "RW"
//			}
//		],
//		"docker": {
//			"image": "mysql:5.7.12",
//			"network": "BRIDGE",
//			"portMappings": [
//				{
//					"containerPort": 3306,
//					"hostPort": 16666,
//					"servicePort": 10000,
//					"protocol": "tcp"
//				}
//			],
//		"forcePullImage": false
//		}
//	},
//	"env": {
//		"MYSQL_USER": "wordpress",
//		"MYSQL_PASSWORD": "secret",
//		"MYSQL_ROOT_PASSWORD": "supersecret",
//		"MYSQL_DATABASE": "wordpress"
//	},
//	"upgradeStrategy": {
//		"minimumHealthCapacity": 0,
//		"maximumOverCapacity": 0
//	},
//	"deployments": [
//		{
//			"id": "95c520f9-ed2d-4488-a966-fa9cc6bb8513"
//		}
//	]
//}

type App struct {
	Id              string            `json:"id"`
	Cpus            float32           `json:"cpus"`
	Mem             uint32            `json:"mem"`
	Disk            uint32            `json:"disk"`
	Instances       uint32            `json:"instances"`
	Container       *Container        `json:"container,omitempty"`
	Env             map[string]string `json:"env,omitempty"`
	UpgradeStrategy *UpgradeStrategy  `json:"upgradestrategy,omitempty"`
	Deployments     []deploy          `json:"deployments,omitempty"`
}

type Container struct {
	Type    string   `json:"type"`
	Volumes []Volume `json:"volumes, omitempty"`
	Docker  Docker   `json:"docker"`
}

type Docker struct {
	Image          string        `json:"image"`
	NetWork        string        `json:"network"`
	PortMappings   []PortMapping `json:"portmappings"`
	ForcePullImage bool          `json:"forcepullimage"`
}

type PortMapping struct {
	ContainerPort uint32 `json:"containerport"`
	HostPort      uint32 `json:"hostport"`
	ServicePort   uint32 `json:"servicePort"`
	Protocol      string `json:"protocol"`
}

type Volume struct {
	ContainerPath string      `json:"containerpath"`
	Mode          string      `json:"mode"`
	Persistent    *Persistent `json:"persistent,omitempty"`
	HostPath      string      `json:"hostpath,omitempty"`
}

type Persistent struct {
	Size uint `json:"size"`
}

type UpgradeStrategy struct {
	MinimumHealthCapacity uint32 `json:"minimumhealthcapacity"`
	MaximumOverCapacity   uint32 `json:"maximumovercapacity"`
}

type deploy struct {
	Id string `json:"id"`
}

// --------------------------deployment------------------------------
//{
//	"id": "3ca0eefc-7e4f-4b1d-9bf0-d29dc004a776",
//	"version": "2016-06-23T02:41:03.683Z",
//	"affectedApps": [
//		"/xxxx"
//	],
//	"steps": [
//		{
//			"actions": [
//					{
//						"type": "StartApplication",
//						"app": "/xxxx"
//					}
//				]
//		},
//		{
//			"actions": [
//					{
//						"type": "ScaleApplication",
//						"app": "/xxxx"
//					}
//				]
//		}
//		],
//	"currentActions": [
//		{
//			"action": "ScaleApplication",
//			"app": "/xxxx",
//			"readinessCheckResults": []
//		}
//	],
//	"currentStep": 2,
//	"totalSteps": 2
//}

type Deployment struct {
	Id             string          `json:"id"`
	Version        string          `json:"version"`
	AffectedApps   []string        `json:"affectedApps"`
	Steps          []step          `json:"steps"`
	CurrentActions []currentAction `json:"currentActions"`
	CurrentStep    int             `json:"currentStep"`
	TotalSteps     int             `json:"totalStep"`
}

type step struct {
	Actions []action `json:"actions"`
}

type action struct {
	Type string `json:"type"`
	App  string `json:"app"`
}

type currentAction struct {
	action
	ReadinessCheckResults []string `json:"readinessCheckResults"`
}

// --------------------------task------------------------------
//{
//	"tasks": [
//			{
//				"id": "mysql.29183f57-38f0-11e6-a6fb-0242d9620757",
//				"slaveId": "a400f5b7-d734-4922-98a6-1d35ae7179d3-S1",
//				"host": "192.168.12.51",
//				"startedAt": "2016-06-23T03:11:23.076Z",
//				"stagedAt": "2016-06-23T03:11:22.205Z",
//				"ports": [
//					16666
//				],
//				"version": "2016-06-23T03:11:03.559Z",
//				"ipAddresses": [
//					{
//						"ipAddress": "172.17.0.2",
//						"protocol": "IPv4"
//					}
//				],
//				"localVolumes": [
//					{
//						"containerPath": "mysql_data",
//						"persistenceId": "mysql#mysql_data#29183f56-38f0-11e6-a6fb-0242d9620757"
//					}
//				],
//				"appId": "/mysql",
//				"servicePorts": [
//					10000
//				]
//			}
//	]
//}

type Tasks struct {
	Tasks []task `json:"tasks"`
}

type task struct {
	Id           string        `json:"id"`
	SlaveId      string        `json:"slaveId"`
	Host         string        `json:"host"`
	StartedAt    string        `json:"startedAt"`
	StagedAt     string        `json:"stagedAt"`
	Ports        []int         `json:"ports"`
	Version      string        `json:"version"`
	IpAddresses  []ipAddress   `json:"ipAddresses"`
	LocalVolumes []localVolume `json:"localVolumes"`
	AppId        string        `json:"appId"`
	ServicePorts []int         `json:"servicePorts"`
}

type ipAddress struct {
	IpAddress string
	Protocol  string
}

type localVolume struct {
	ContainerPath string `json:"containerPath"`
	PersistenceId string `json:"persistenceId"`
}
