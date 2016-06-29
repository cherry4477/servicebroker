package servicebroker

func (c *Catalog) GetService(id string) *Service {
	for _, svc := range c.Services {
		if svc.ID == id {
			return svc
		}
	}

	return nil
}

func (c *Catalog) RangeCatalogFunc(fn func(*Service)) {
	for _, svc := range c.Services {
		fn(svc)
	}
}

func (s *Service) GetPlan(id string) *Plan {
	for _, plan := range s.Plans {
		if plan.ID == id {
			return plan
		}
	}
	return nil
}
