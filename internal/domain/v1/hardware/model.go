package hardware

import (
	"fmt"
	// "microservice-backend/internal/domain/v1/colorgroup"
	// "microservice-backend/internal/domain/v1/manufacturer"
	// "microservice-backend/internal/domain/v1/os"
	// "microservice-backend/internal/domain/v1/product"
	"strings"
	"time"

	"github.com/Filipza/excel-mapping-tool/internal/domain/v1/product"
)

// type Hardware struct {
// 	product.ProductDetail
// 	Id               string                          `json:"id"` // MOVED _id -> id
// 	Name             string                          `json:"name"`
// 	DisplaySize      float64                         `json:"displaySize,omitempty"`
// 	Stock            int                             `json:"stock,omitempty"`
// 	WKZSum           float64                         `json:"wkzsum,omitempty"`
// 	Variants         []*Variant                      `json:"variants,omitempty"`
// 	ExtendedVariants []*ExtendedVariant              `json:"extendedVariants,omitempty"`
// 	PKBonus          bool                            `json:"pkBonus"`
// 	PKName           string                          `json:"pkName"`
// 	PKType           string                          `json:"pkType"`
// 	Conditions       string                          `json:"conditions"`
// 	Type             string                          `json:"type"`
// 	Options          []*product.ProductOption        `json:"options,omitempty"`
// 	PricingIntervals *product.PricingIntervalAverage `json:"pricingIntervals,omitempty"`
// 	RouterPrice      float64                         `json:"routerPrice"`
// 	DeliveryPrice    float64                         `json:"deliveryPrice"`
// }

// func (hw *Hardware) DefaultVariant() (*Variant, bool) {
// 	for _, variant := range hw.Variants {
// 		if variant.Default {
// 			return variant, true
// 		}
// 	}

// 	if len(hw.Variants) > 0 {
// 		return hw.Variants[0], true
// 	}

// 	return nil, false
// }

// func (hw *Hardware) AsLookup() *HardwareLookup {
// 	result := &HardwareLookup{
// 		Id:           hw.Id,
// 		Name:         hw.Name,
// 		Manufacturer: hw.Manufacturer,
// 		OS:           hw.OS,
// 		Variants:     make([]*VariantLookup, len(hw.Variants)),
// 		Type:         hw.Type,
// 	}

// 	for i, variant := range hw.Variants {
// 		result.Variants[i] = variant.AsLookup(result.Label())
// 	}

// 	return result
// }

// func (hw *Hardware) DefaultVariantInStock() (*Variant, bool) {
// 	if def, ok := hw.DefaultVariant(); ok && def.InStock() {
// 		return def, true
// 	}

// 	for _, variant := range hw.Variants {
// 		if variant.InStock() {
// 			return variant, true
// 		}
// 	}

// 	return nil, false
// }

// func (hw *Hardware) Variant(ebootisId string) (*Variant, bool) {
// 	for _, variant := range hw.Variants {
// 		if variant.EbootisId == ebootisId {
// 			return variant, true
// 		}
// 	}
// 	return nil, false
// }

// func (hw *Hardware) SelectedExtendedVariant() (*ExtendedVariant, bool) {
// 	for _, ev := range hw.ExtendedVariants {
// 		if ev.IsSelected {
// 			return ev, true
// 		}
// 	}
// 	// fallback if no variant is selected; return first
// 	if len(hw.ExtendedVariants) > 0 {
// 		return hw.ExtendedVariants[0], true
// 	}

// 	return nil, false
// }

// func (hw *Hardware) VariantIds() []string {
// 	result := make([]string, len(hw.Variants))
// 	for i, v := range hw.Variants {
// 		result[i] = v.EbootisId
// 	}

// 	return result
// }

// func (hw *Hardware) PossibleExtendedVariants() []*ExtendedVariant {
// 	results := make([]*ExtendedVariant, 0)
// 	if variant, ok := hw.SelectedExtendedVariant(); ok && variant.Variant != nil {
// 		// filter for variants with same storage
// 		for _, ev := range hw.ExtendedVariants {
// 			if ev.Variant != nil && ev.Variant.Storage == variant.Variant.Storage {
// 				results = append(results, ev)
// 			}
// 		}
// 	}
// 	return results
// }

// func (hw *Hardware) PossibleExtendedVariantInStock() (*ExtendedVariant, bool) {
// 	for _, ev := range hw.PossibleExtendedVariants() {
// 		if ev.Stock > 0 {
// 			return ev, true
// 		}
// 	}
// 	return hw.SelectedExtendedVariant()
// }

// type ExtendedVariant struct {
// 	IsSelected bool     `json:"isSelected"`
// 	Stock      int      `json:"stock"`
// 	Variant    *Variant `json:"variant"`
// 	// from extendedVariants-Call
// 	EyeCatcher      *product.EyeCatcher `json:"eyeCatchers,omitempty"`
// 	HardwareId      string              `json:"hardwareId,omitempty"`
// 	HardwareName    string              `json:"hardwareName,omitempty"`
// 	Manufacturer    string              `json:"manufacturer,omitempty"`
// 	ManufacturerUrl string              `json:"manufacturerUrl,omitempty"`
// 	ImageCount      int                 `json:"numberOfImages,omitempty"`
// 	OfferGroup      string              `json:"offerGroupUrl,omitempty"`
// 	OfferId         string              `json:"offerId,omitempty"`
// 	OfferType       string              `json:"offerType,omitempty"`
// 	Tariff          *tariff.Tariff      `json:"tariff,omitempty"`
// }

// type Variant struct {
// 	Color                 *Color                   `json:"color"`
// 	ColorGroup            *Color                   `json:"colorGroup"`
// 	Default               bool                     `json:"default"`
// 	DeliveryTime          *product.DeliveryTime    `json:"deliveryTime"`
// 	DeliveryPrice         *product.DeliveryPrice   `json:"deliveryPrice,omitempty"`
// 	DeliverySetting       *product.DeliverySetting `json:"deliverySettings,omitempty"`
// 	EbootisId             string                   `json:"ebootisId"`
// 	ImageCount            int                      `json:"numberOfImages"`
// 	Price                 float64                  `json:"price,omitempty"`
// 	BasisPrice            float64                  `json:"basisPrice,omitempty"`
// 	Stock                 int                      `json:"stock"`
// 	Storage               int                      `json:"storage"`
// 	TariffMap             interface{}              `json:"tariffMap"` // TODO define model
// 	Tariffs               interface{}              `json:"tariffs"`   // TODO define model / behavior
// 	Url                   string                   `json:"url"`
// 	Images                []*product.Images        `json:"images"`
// 	ExternalArticleNumber string                   `json:"externalArticlenumber,omitempty"`
// 	PkCouponName          string                   `json:"pkCouponName"`
// 	PkCouponValue         float64                  `json:"pkCouponValue"`
// 	TradeInActive         bool                     `json:"tradeInActive"`
// 	TradeInValue          string                   `json:"tradeInValue"`
// 	ProtectedVoucher      bool                     `json:"protectedVoucher"`
// }

// func (va *Variant) InStock() bool {
// 	return va.Stock > 0
// }

// func (va *Variant) ColorUrl() string {
// 	if va.Color != nil {
// 		return strings.ToLower(strings.ReplaceAll(va.Color.Name, " ", "-"))
// 	}
// 	return ""
// }

// func (va *Variant) StorageUrl() string {
// 	if va.Storage > 0 {
// 		return strconv.Itoa(va.Storage)
// 	}
// 	return ""
// }

// func (va *Variant) EbootisSuffix() int {
// 	if _, snum, ok := strings.Cut(va.EbootisId, "-"); ok {
// 		suffix, _ := strconv.Atoi(snum)
// 		return suffix
// 	}
// 	return 0
// }

// func (va *Variant) AsLookup(hardwareLabel string) *VariantLookup {
// 	name := fmt.Sprintf("%s - %s", va.EbootisId, hardwareLabel)
// 	if va.Color != nil {
// 		name = fmt.Sprintf("%s %s", name, va.Color.Name)
// 	}
// 	if va.Storage > 0 {
// 		name = fmt.Sprintf("%s %dGB", name, va.Storage)
// 	}
// 	price := va.Price
// 	if price < 0.01 {
// 		price = 0.0
// 	}
// 	name = fmt.Sprintf("%s    %.2f EUR", name, price)

// 	return &VariantLookup{
// 		Id:      va.EbootisId,
// 		Name:    name,
// 		Default: va.Default,
// 	}
// }

// type Color struct {
// 	Name string `json:"name"`
// 	Hex  string `json:"hex"`
// }

type VariantLookup struct {
	Id      string `json:"id"`
	Name    string `json:"name,omitempty"`
	Default bool   `json:"default"`
}

type HardwareLookup struct {
	Id           string           `json:"id"`
	Name         string           `json:"name,omitempty"`
	Manufacturer string           `json:"manufacturer,omitempty"`
	OS           string           `json:"os,omitempty"`
	Variants     []*VariantLookup `json:"variants,omitempty"`
	Type         string           `json:"type,omitempty"`
}

func (hl *HardwareLookup) Label() string {
	return strings.TrimSpace(fmt.Sprintf("%s %s", hl.Manufacturer, hl.Name))
}

type HardwareCRUD struct {
	Id               string                     `json:"id"`
	Name             string                     `json:"name"`
	Type             string                     `json:"type"`
	PKBonus          bool                       `json:"pkBonus"`
	PKName           string                     `json:"pkName"`
	PKType           string                     `json:"pkType"`
	SimType          string                     `json:"simType"`
	DisplaySize      float64                    `json:"displaySize"`
	EyeCatcher       string                     `json:"eyeCatcher"`
	Bullets          []*product.Option          `json:"bullets"`
	Wkz              []*product.Option          `json:"wkz"`
	Variants         []*VariantCRUD             `json:"variants"`
	CustomFields     []*product.CustomField     `json:"customFields"`
	Conditions       string                     `json:"conditions"`
	Highlights       []string                   `json:"highlights"`
	PricingIntervals []*product.PricingInterval `json:"pricingIntervals"`
	RouterPrice      float64                    `json:"routerPrice"`
	DeliveryPrice    float64                    `json:"deliveryPrice"`
	Options          []*product.OptionLookup    `json:"options"`
	UpdatedAt        *time.Time                 `json:"updatedAt"`
	CreatedAt        *time.Time                 `json:"createdAt"`
}

func (hw *HardwareCRUD) AsLookup() *HardwareLookup {
	result := &HardwareLookup{
		Id:       hw.Id,
		Name:     hw.Name,
		Variants: make([]*VariantLookup, len(hw.Variants)),
		Type:     hw.Type,
	}

	for i, variant := range hw.Variants {
		result.Variants[i] = variant.AsLookup(result.Label())
	}

	return result
}

func (hw *HardwareCRUD) Variant(id string) (*VariantCRUD, bool) {
	for _, va := range hw.Variants {
		if va.EbootisId == id {
			return va, true
		}
	}
	return nil, false
}

func (hw *HardwareCRUD) VariantViaArticleNo(artNo string) (*VariantCRUD, bool) {
	for _, va := range hw.Variants {
		if va.ExternalArticleNumber == artNo {
			return va, true
		}
	}
	return nil, false
}

func (hw *HardwareCRUD) FirstVariant() (*VariantCRUD, bool) {
	if len(hw.Variants) > 0 {
		return hw.Variants[0], true
	}
	return nil, false
}

type VariantCRUD struct {
	EbootisId             string  `json:"ebootisId"`
	Default               bool    `json:"default"`
	ExternalArticleNumber string  `json:"externalArticleNumber"`
	EAN                   string  `json:"ean"`
	Storage               int     `json:"storage"`
	Price                 float64 `json:"price"`
	ColorName             string  `json:"colorName"`
	RiskPremium           float64 `json:"riskPremium"`
	DeliveryTimeDays      int     `json:"deliveryTimeDays"`
	DeliveryTimeText      string  `json:"deliveryTimeText"`
	PublicationDate       string  `json:"publicationDate"` // as YYYY-MM-DD
	PkCouponName          string  `json:"pkCouponName"`
	PkCouponValue         float64 `json:"pkCouponValue"`
}

func (va *VariantCRUD) AsLookup(hardwareLabel string) *VariantLookup {
	name := fmt.Sprintf("%s - %s", va.EbootisId, hardwareLabel)
	if va.ColorName != "" {
		name = fmt.Sprintf("%s %s", name, va.ColorName)
	}
	if va.Storage > 0 {
		name = fmt.Sprintf("%s %dGB", name, va.Storage)
	}
	price := va.Price
	if price < 0.01 {
		price = 0.0
	}
	name = fmt.Sprintf("%s    %.2f EUR", name, price)

	return &VariantLookup{
		Id:      va.EbootisId,
		Name:    name,
		Default: va.Default,
	}
}
