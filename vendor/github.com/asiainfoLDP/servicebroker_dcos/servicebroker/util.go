package servicebroker

import (
	"fmt"
	"strings"
)

func (c *Catalog) GetService(id string) *Service {
	for _, svc := range c.Services {
		if svc.ID == id {
			return svc
		}
	}

	return nil
}

func (s *Service) GetPlan(id string) *Plan {
	for _, plan := range s.Plans {
		if plan.ID == id {
			return plan
		}
	}
	return nil
}

func (o *InstanceRequest) Validate(c *Catalog) error {
	if len(strings.TrimSpace(o.ServiceId)) == 0 {
		return fmt.Errorf("invalid json field service_id")
	}
	if len(strings.TrimSpace(o.PlanId)) == 0 {
		return fmt.Errorf("invalid json field plan_id")
	}
	if len(strings.TrimSpace(o.OrganizationGuid)) == 0 {
		return fmt.Errorf("invalid json field organization_guid")
	}

	return nil
}
