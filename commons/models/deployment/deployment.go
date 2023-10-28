package deployment

import (
	"fmt"
	"sync"
	"time"
)

type Deployment struct {
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	//Name     string   `json:"name"`
	ConfigurationInitializedFromFirebaseSource bool             `json,omitempty"`
	Tiers                                      map[string]*Tier `json:"tiers" yaml:"tiers"`
	Secrets                                    Secrets          `json:"secrets" yaml:"secrets"`
	//Status                                     Status           `json:"status" yaml:"status"`
}

type Metadata struct {
	Locations []Location `json:"locations" yaml:"locations"`
	Name      string     `json:"name" yaml:"name"`
	Sources   Sources    `json:"sources" yaml:"sources"`
}

type Location string

type Sources struct {
	Mapping struct {
		API struct {
			Repository string   `json:"repository" yaml:"repository"`
			Branches   []string `json:"branches" yaml:"branches"`
		} `json:"api" yaml:"api"`
	} `json:"mapping" yaml:"mapping"`
	History struct {
		Builds  map[string]string `json:"builds" yaml:"builds"`
		Deploys map[string]string `json:"deploys" yaml:"deploys"`
	} `json:"history" yaml:"history"`
}

type Tier struct {
	initialYamlTemplate   []byte                `json:"-" yaml:"-"`
	Entrypoint            Entrypoint            `json:"entrypoint" yaml:"entrypoint"`
	LatestControlSequence LatestControlSequence `json:"latest_control_sequence" yaml:"latest_control_sequence"`
	Artifacts             []Artifact            `json:"artifacts" yaml:"artifacts"`
	HealthChecks          HealthChecks          `json:"health-checks" yaml:"health-checks"`
}

//type Credentials struct {
//	Aws Aws `json:"aws"`
//}

type Secrets map[string]any

type Entrypoint struct {
	Autostart bool     `json:"autostart" yaml:"autostart"`
	Command   []string `json:"command" yaml:"command"`
}

type LatestControlSequence struct {
	AutoDeploy            bool     `json:"auto_deploy" yaml:"auto_deploy"`
	CreatedAt             string   `json:"created_at" yaml:"created_at"`
	CommandSequence       []string `json:"command_sequence" yaml:"command_sequence"`
	ConfigurationSnapshot string   `json:"configuration_snapshot" yaml:"configuration_snapshot"`
}

//type Aws struct {
//	Token  string `json:"token"`
//	Secret string `json:"secret"`
//}

type Artifact struct {
	Source     []Source   `json:"source" yaml:"source"`
	Middleware Middleware `json:"middleware" yaml:"middleware"`
	Deploy     []Deploy   `json:"deploy" yaml:"deploy"`
}

type Source struct {
	From From   `json:"from" yaml:"from"`
	To   string `json:"to" yaml:"to"`
}

type From struct {
	S3 S3 `json:"s3" yaml:"s3"`
}

type S3 struct {
	Bucket string `json:"bucket" yaml:"bucket"`
	Region string `json:"region" yaml:"region"`
	Object string `json:"object" yaml:"object"`
}

type Middleware struct {
	Envsubst Envsubst `json:"envsubst" yaml:"envsubst"`
}

type Envsubst struct {
	Target    []string  `json:"target" yaml:"target"`
	Variables Variables `json:"variables" yaml:"variables"`
}

type Variables map[string]string

type Deploy struct {
	Scripts struct {
		BeforeDeploy Script `json:"before-deploy" yaml:"before-deploy"`
		AfterDeploy  Script `json:"after-deploy" yaml:"after-deploy"`
	} `json:"scripts" yaml:"scripts"`
	From string `json:"from" yaml:"from"`
	To   string `json:"to" yaml:"to"`
	Mode string `json:"mode" yaml:"mode"`
}

type HealthChecks struct {
	TCP   Tcp   `json:"tcp" yaml:"tcp"`
	Mssql Mssql `json:"mssql" yaml:"mssql"`
	Redis Redis `json:"redis" yaml:"redis"`
	HTTP  Http  `json:"http" yaml:"http"`
}

//type NetworkHealthCheck interface {
//	*Tcp | *Mssql | *Redis | *Http
//}

type Tcp struct {
	Targets []struct {
		Host string `json:"host" yaml:"host"`
		Port string `json:"port" yaml:"port"`
	} `json:"targets" yaml:"targets"`
	Options struct {
		Timeout            int `json:"timeout" yaml:"timeout"`
		Interval           int `json:"interval" yaml:"interval"`
		UnhealthyThreshold int `json:"unhealthy_threshold" yaml:"unhealthy_threshold"`
		HealthyThreshold   int `json:"healthy_threshold" yaml:"healthy_threshold"`
	} `json:"options" yaml:"options"`
}

type Mssql struct {
	Targets []struct {
		Host     string `json:"host" yaml:"host"`
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
		Database string `json:"database" yaml:"database"`
		Port     string `json:"port" yaml:"port"`
	} `json:"targets" yaml:"targets"`
	Options struct {
		Timeout            int `json:"timeout" yaml:"timeout"`
		Interval           int `json:"interval" yaml:"interval"`
		UnhealthyThreshold int `json:"unhealthy_threshold" yaml:"unhealthy_threshold"`
		HealthyThreshold   int `json:"healthy_threshold" yaml:"healthy_threshold"`
	} `json:"options" yaml:"options"`
}

type Redis struct {
	Targets []struct {
		Host     string `json:"host" yaml:"host"`
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
		Port     string `json:"port" yaml:"port"`
		Tls      string `json:"tls" yaml:"tls"`
	} `json:"targets" yaml:"targets"`
	Options struct {
		Timeout            int `json:"timeout" yaml:"timeout"`
		Interval           int `json:"interval" yaml:"interval"`
		UnhealthyThreshold int `json:"unhealthy_threshold" yaml:"unhealthy_threshold"`
		HealthyThreshold   int `json:"healthy_threshold" yaml:"healthy_threshold"`
	} `json:"options" yaml:"options"`
}

type Http struct {
	Targets []struct {
		URL    string `json:"url" yaml:"url"`
		Status string `json:"status" yaml:"status"`
	} `json:"targets" yaml:"targets"`
	Options struct {
		Timeout            int `json:"timeout" yaml:"timeout"`
		Interval           int `json:"interval" yaml:"interval"`
		UnhealthyThreshold int `json:"unhealthy_threshold" yaml:"unhealthy_threshold"`
		HealthyThreshold   int `json:"healthy_threshold" yaml:"healthy_threshold"`
	} `json:"options" yaml:"options"`
}

func GetNodeStatus() *NodeStatus {
	if nodeStatus == nil {
		nodeStatus = &NodeStatus{}
	}

	return nodeStatus
}

type NodeStatus struct {
	SyncMutext            *sync.Mutex            `json:",omitempty" yaml:",omitempty"`
	AgentHeartbeatAt      string                 `json:"agent_heartbeat_at" yaml:"agent_heartbeat_at"`
	Installation          Installation           `json:"installation,omitempty" yaml:"installation,omitempty"`
	ArtifactsDeployStatus []ArtifactDeployStatus `json:"artifacts_deploy" yaml:"artifacts_deploy"`
	HealthChecks          struct {
		TCP   []NetworkHealthCheckStatus `json:"tcp,omitempty" yaml:"tcp,omitempty"`
		Mssql []NetworkHealthCheckStatus `json:"mssql,omitempty" yaml:"mssql,omitempty"`
		Redis []NetworkHealthCheckStatus `json:"redis,omitempty" yaml:"redis,omitempty"`
		HTTP  []NetworkHealthCheckStatus `json:"http,omitempty" yaml:"http,omitempty"`
	} `json:"health_checks,omitempty" yaml:"health_checks,omitempty"`
	LatestControlSequence LatestControlSequence `json:"latest_control_sequence,omitempty" yaml:"latest_control_sequence,omitempty"`
}

func (ns *NodeStatus) Update(status NodeStatus) *NodeStatus {
	ns.SyncMutext.Lock()
	*ns = status
	ns.SyncMutext.Unlock()
	return ns
}

func (ns *NodeStatus) StartSubmittingToRedis(d *Deployment) *NodeStatus {

	ttl := 30 * time.Second

	go func() {
		for {
			// TODO: heartbeat
			PutIntoRedis(fmt.Sprintf(string(RDB_AGENT_HEARTBEAT_PATH), d.Name(), d.tierName(), agentId), GetTimeNowString(), ttl)

			// TODO: installation

			// TODO: artifactsDeploy
			for index, artifactStatus := range ns.ArtifactsDeployStatus {
				PutIntoRedis(fmt.Sprintf(string(RDB_ARTIFACTS_DEPLOY_STATUS_PATH), d.Name(), d.tierName(), agentId, index), artifactStatus, ttl)
			}

			// TODO: health checks
			for index, check := range ns.HealthChecks.TCP {
				PutIntoRedis(fmt.Sprintf(string(RDB_TCP_HEALTH_CHECKS_STATUS_PATH), d.Name(), d.tierName(), agentId, index), check, ttl)
			}
			for index, check := range ns.HealthChecks.Mssql {
				PutIntoRedis(fmt.Sprintf(string(RDB_MSSQL_HEALTH_CHECKS_STATUS_PATH), d.Name(), d.tierName(), agentId, index), check, ttl)
			}
			for index, check := range ns.HealthChecks.Redis {
				PutIntoRedis(fmt.Sprintf(string(RDB_REDIS_HEALTH_CHECKS_STATUS_PATH), d.Name(), d.tierName(), agentId, index), check, ttl)
			}
			for index, check := range ns.HealthChecks.HTTP {
				PutIntoRedis(fmt.Sprintf(string(RDB_HTTP_HEALTH_CHECKS_STATUS_PATH), d.Name(), d.tierName(), agentId, index), check, ttl)
			}

			// TODO: control sequence
			PutIntoRedis(fmt.Sprintf(string(RDB_LATEST_CONTROL_SEQUENCE_STATUS_PATH), d.Name(), d.tierName(), agentId), d.Tier().LatestControlSequence, ttl)

			time.Sleep(1 * time.Second)
		}
	}()

	return ns
}

type NetworkHealthCheckStatus struct {
	//Interval  int    `json:"interval,omitempty" yaml:"interval,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Target    string `json:"target,omitempty" yaml:"target,omitempty"`
	Status    string `json:"status" yaml:"status"`
	Message   string `json:"message,omitempty" yaml:"message,omitempty"`
}

type Installation struct {
	State   string `json:"state,omitempty" yaml:"state,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
}

type ArtifactDeployStatus struct {
	Started         time.Time `json:"started,omitempty" yaml:"started,omitempty"`
	Ended           time.Time `json:"ended,omitempty" yaml:"ended,omitempty"`
	DurationSeconds int64     `json:"duration,omitempty" yaml:"duration,omitempty"`
	Source          string    `json:"source,omitempty" yaml:"source,omitempty"`
	ControlSequence string    `json:"control_sequence" yaml:"control_sequence"`
	State           string    `json:"state,omitempty" yaml:"state,omitempty"`
	Message         string    `json:"message,omitempty" yaml:"message,omitempty"`
	Log             string    `json:"log,omitempty" yaml:"log,omitempty"`
}
