package product

type Option struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type PricingInterval struct {
	EbootisId     string  `json:"ebootisId,omitempty" bson:"ebootisId" mapstructure:"ebootisId"`
	MonthInterval int     `json:"monthInterval" bson:"monthInterval" mapstructure:"monthInterval"`
	StartMonth    int     `json:"startMonth" bson:"startMonth" mapstructure:"startMonth"`
	EndMonth      int     `json:"endMonth" bson:"endMonth" mapstructure:"endMonth"`
	Price         float64 `json:"price" bson:"price" mapstructure:"price"`
}

type ProductOption struct {
	Id                    string             `json:"id" bson:"id" mapstructure:"id"`
	Name                  string             `json:"name,omitempty" bson:"name" mapstructure:"name"`
	EbootisId             string             `json:"ebootisId,omitempty" bson:"ebootisId" mapstructure:"ebootisId"`
	Type                  string             `json:"type,omitempty" bson:"type" mapstructure:"type"`
	Category              string             `json:"category,omitempty" bson:"category" mapstructure:"category"`
	PricingIntervals      []*PricingInterval `json:"pricingIntervals,omitempty" bson:"pricingIntervals" mapstructure:"pricingsIntervals"`
	BillingFrequency      string             `json:"billingFrequency,omitempty" bson:"billingFrequency" mapstructure:"billingFrequency"`
	MonthlyPrice          float64            `json:"monthlyPrice,omitempty" bson:"monthlyPrice" mapstructure:"monthlyPrice"`
	OneTimePrice          float64            `json:"oneTimePrice,omitempty" bson:"oneTimePrice" mapstructure:"oneTimePricing"`
	Description           string             `json:"description,omitempty" bson:"description" mapstructure:"description"`
	InterfaceCode         string             `json:"interfaceCode,omitempty" bson:"interfaceCode" mapstructure:"interfaceCode"`
	InterfaceGroup        string             `json:"interfaceGroup,omitempty" bson:"interfaceGroup" mapstructure:"interfaceGroup"`
	Details               string             `json:"details,omitempty" bson:"details" mapstructure:"details"`
	ServiceItemCode       string             `json:"serviceItemCode,omitempty" bson:"serviceItemCode" mapstructure:"serviceItemCode"`
	FrontendName          string             `json:"frontendName,omitempty" bson:"frontendName" mapstructure:"frontendName"`
	PossiblePackgroupCode string             `json:"possiblePackgroupCode,omitempty" bson:"possiblePackgroupCode" mapstructure:"possiblePackgroupCode"`
}

type CustomField struct {
	Id      string `json:"id"`
	Label   string `json:"label"`
	Value   string `json:"value"`
	Section string `json:"section"` // === key from Bullet
}

type OptionLookup struct {
	Id           string `json:"id"`
	EbootsId     string `json:"ebootisId"`
	Type         string `json:"type"`
	Category     string `json:"category"`
	Name         string `json:"name"`
	FrontendName string `json:"frontendName"`
}
