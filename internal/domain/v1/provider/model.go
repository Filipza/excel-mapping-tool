package provider

type ServiceProvider struct {
	Id                   string            `json:"id"`
	Name                 string            `json:"providerName"`
	Label                string            `json:"label"`
	Type                 string            `json:"type"`
	InternalName         string            `json:"internalName"`
	Categories           map[string]any    `json:"hardwareCategories"`
	OrderCategory        string            `json:"orderCategory"`
	TariffProvider       string            `json:"tariffProvider"`
	TariffNet            string            `json:"tariffNet"`
	AddressIdKey         string            `json:"addressIdKey"`
	NumberPortingCarrier string            `json:"numberPortingCarrier"`
	Key                  string            `json:"key"`
	Meta                 map[string]string `json:"meta,omitempty"`
	Image                string            `json:"image"`
	TileImage            string            `json:"tileImage"`
	Info                 string            `json:"info"`
	Order                int               `json:"order"`
}

func (sp *ServiceProvider) GetLabel() string {
	if sp.Label != "" {
		return sp.Label
	}
	return sp.Name
}
