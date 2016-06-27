package api

type Interface interface {
	ApplicationsInterface
	DeploymentsInterface
	TasksInterface
}

func NewClientInterface(host, token string) (Interface, error) {
	c := authConfig(token)
	a := &dcosOption{auth: &c, host: host}
	if acs_token, err := a.auth.Exchange(host); err != nil {
		return nil, err
	} else {
		a.Acs_token = acs_token
	}

	if err := a.validate(); err != nil {
		return nil, err
	}

	return a, nil
}

func (o *dcosOption) Application() ApplicationInterface {
	return &applicationOption{o}
}

func (o *dcosOption) Deployment() DeploymentInterface {
	return &deploymentOption{o}
}

func (o *dcosOption) Task() TaskInterface {
	return &taskOption{o}
}
