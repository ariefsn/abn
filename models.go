package abn

type M map[string]interface{}
type A []map[string]interface{}

type AbnAddressModel struct {
	Date     string `json:"date"`
	Postcode string `json:"postcode"`
	State    string `json:"state"`
}

type AbnEntityTypeModel struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type AbnEntityModel struct {
	Name string             `json:"name"`
	Type AbnEntityTypeModel `json:"type"`
}

type AbnModel struct {
	Abn                    string          `json:"abn"`
	Status                 string          `json:"status"`
	AbnStatusEffectiveFrom string          `json:"statusEffectiveFrom"`
	Acn                    string          `json:"acn"`
	Address                AbnAddressModel `json:"address"`
	BusinessNames          []string        `json:"businessNames"`
	Entity                 AbnEntityModel  `json:"entity"`
	Gst                    string          `json:"gst"`
}

type AbnSearchModel struct {
	Abn       string  `json:"abn"`
	Status    string  `json:"status"`
	IsCurrent bool    `json:"isCurrent"`
	Name      string  `json:"name"`
	NameType  string  `json:"nameType"`
	Postcode  string  `json:"postcode"`
	Score     float64 `json:"score"`
	State     string  `json:"state"`
}
