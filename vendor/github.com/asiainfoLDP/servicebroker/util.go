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

func (c *Catalog) Merge(newCatalog *Catalog) {
	if newCatalog != nil {
		c.Services = append(c.Services, newCatalog.Services...)
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
