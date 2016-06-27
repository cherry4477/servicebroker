package servicebroker

// Catalog describes Cloud Foundry catalog
type Catalog struct {
	Services []*Service `json:"services"`
}

// Service describes Cloud Foundry service
type Service struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Bindable    bool         `json:"bindable"`
	Tags        []string     `json:"tags,omitempty"`
	Metadata    *ServiceMeta `json:"metadata,omitempty"`
	Requires    []string     `json:"requires,omitempty"`
	Plans       []*Plan      `json:"plans"`
	Dashboard   *Dashboard   `json:"dashboard_client,omitempty"`
}

// Plan describes Cloud Foundry plan structure
type Plan struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Metadata    *PlanMeta `json:"metadata,omitempty"`
	Free        bool      `json:"free,omitempty"`
}

// PlanMeta describers Cloud Foundry plan meta-data
type PlanMeta struct {
	Bullets     []string `json:"bullets"`
	Costs       []Cost   `json:"costs"`
	DisplayName string   `json:"displayName"`
}

// Cost describers Cloud Foundry plan Cost
type Cost struct {
	Amount *Amount `json:"amount"`
	Unit   string  `json:"unit"`
}

// Amount describers Cloud Foundry cost amount
type Amount struct {
	Usd float32 `json:"usd"`
	Eur float32 `json:"eur"`
}

// Dashboard describes Cloud Foundry dashboard
type Dashboard struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	URI    string `json:"redirect_uri"`
}

// ServiceMeta describers Cloud Foundry service meta-data
type ServiceMeta struct {
	DisplayName         string `json:"displayName"`
	ImageURL            string `json:"imageUrl"`
	Description         string `json:"longDescription"`
	ProviderDisplayName string `json:"providerDisplayName"`
	DocURL              string `json:"documentationUrl"`
	SupportURL          string `json:"supportUrl"`
}

type InstanceRequest struct {
	ServiceId        string `json:"service_id"`
	PlanId           string `json:"plan_id"`
	OrganizationGuid string `json:"organization_guid"`
	SpaceGuid        string `json:"space_guid"`
	//Incomplete       bool        `json:"accepts_incomplete, omitempty"`
	Parameters interface{} `json:"parameters, omitempty"`
}

type ServiceBindingRequest struct {
	InstanceID string `json:"-"`
	BindingID  string `json:"-"`
	ServiceID  string `json:"service_id"`
	PlanID     string `json:"plan_id"`
	AppGUID    string `json:"app_guid"`
}

type ServiceBindingResponse struct {
	Credentials map[string]string `json:"credentials"`
	//SyslogDrainURL string            `json:"syslog_drain_url"`
}
