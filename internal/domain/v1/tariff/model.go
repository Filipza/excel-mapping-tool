package tariff

import (
	"time"

	"github.com/Filipza/excel-mapping-tool/internal/domain/v1/carrier"
	"github.com/Filipza/excel-mapping-tool/internal/domain/v1/product"
	"github.com/Filipza/excel-mapping-tool/internal/domain/v1/provider"
)

type TariffLookup struct {
	Id              string `json:"id"`
	EbootisId       string `json:"ebootisId,omitempty"`
	Name            string `json:"name,omitempty"`
	Carrier         string `json:"carrier,omitempty"`
	ServiceProvider string `json:"serviceProvider,omitempty"`
	Type            string `json:"type,omitempty"`
}

type TariffCRUD struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	SystemName string `json:"systemName"`
	EbootisId  string `json:"ebootisId"`
	Type       string `json:"type"`
	LeadType   int    `json:"leadType"`
	// internet specific
	FrontendName   string  `json:"frontendName"`
	VariationCode  string  `json:"variationCode"`
	ContractId     string  `json:"contractId"`
	DownloadSpeed  float64 `json:"downloadSpeed"` // in Mbit/s
	UploadSpeed    float64 `json:"uploadSpeed"`   // in Mbit/s
	ConnectionType string  `json:"connectionType"`
	// payment
	BasicCharge        float64 `json:"basicCharge"`
	BasicChargeRenewal float64 `json:"basicChargeRenewal"`
	ConnectionFee      float64 `json:"connectionFee"`
	Subsidy            float64 `json:"subsidy"`
	PromotionBonus     float64 `json:"promotionBonus"`
	Provision          float64 `json:"provision"`
	XProvision         float64 `json:"xProvision"` // TODO what the hell ?!?!
	// duration
	ContractTerm    int `json:"contractTerm"`
	PromotionPeriod int `json:"promotionPeriod"`
	// description
	Highlights []string `json:"highlights"`
	LegalNote  string   `json:"legalNote"`
	PibLink    string   `json:"pibLink"`
	// condition / features
	DataVolume      float64 `json:"dataVolume"`
	FreeMinutes     int     `json:"freeMinutes"`
	AllnetFlat      bool    `json:"allnetFlat"`
	SmsFlat         bool    `json:"smsFlat"`
	Lte             bool    `json:"lte"`
	Students        bool    `json:"students"`
	FreelanceTariff bool    `json:"freelanceTariff"`
	// deprecated FamilyCode      string  `json:"familyCode"`
	// sim-cards
	SimEbootisId        string `json:"simEbootisId"`
	DoubleSim           bool   `json:"doubleSim"`
	Subcards            int    `json:"subCards"`
	MaxSubCards         int    `json:"maxSubCards"`
	MinSubCards         int    `json:"minSubCards"`
	SuperSelectMainCard string `json:"superselectMainCard"`
	SuperSelectSubCard  string `json:"superselectSubCard"`
	SecondEbootisId     string `json:"secondEbootisId"`
	DoubleCard          bool   `json:"doubleCard"`
	// relations
	Provider         *provider.ServiceProvider  `json:"provider"`
	Carrier          *carrier.Carrier           `json:"carrier"`
	Bullets          []*product.Option          `json:"bullets"`
	Wkz              []*product.Option          `json:"wkz"`
	PricingIntervals []*product.PricingInterval `json:"pricingIntervals"`
	Options          []*product.OptionLookup    `json:"options"`
	CustomFields     []*product.CustomField     `json:"customFields"`
	// logging
	LastEdit  string     `json:"lastEdit,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt"`
	CreatedAt *time.Time `json:"createdAt"`
}

func (tf *TariffCRUD) AsLookup() *TariffLookup {
	result := &TariffLookup{
		Id:        tf.Id,
		EbootisId: tf.EbootisId,
		Name:      tf.Name,
		Type:      tf.Type,
	}

	if carrier := tf.Carrier; carrier != nil {
		result.Carrier = carrier.Name
	}

	if provider := tf.Provider; provider != nil {
		result.ServiceProvider = provider.Name
	}

	return result
}
