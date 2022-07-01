package abn

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Abn struct {
	guid    string
	message *messages
}

// NewAbn for create new ABN instance with GUID
func NewAbn(guid string) *Abn {
	a := new(Abn)

	a.guid = guid
	a.message = NewMessages()

	return a
}

func (a *Abn) validateGuid() error {
	if a.guid == "" {
		return errors.New(a.message.GuidRequired)
	}

	return nil
}

// AbnSearch for searching with abn code, the results are the ABN Details, status code and the error
func (a *Abn) AbnSearch(abn string) (*AbnModel, int, error) {
	err := a.validateGuid()

	statusCode := http.StatusInternalServerError

	if err != nil {
		return nil, statusCode, err
	}

	if abn == "" {
		return nil, statusCode, errors.New(a.message.AbnRequired)
	}

	client := resty.New()
	client.SetBaseURL(baseUrl)

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"guid": a.guid,
			"abn":  abn,
		}).
		Get(abnPath)

	if resp != nil {
		statusCode = resp.StatusCode()
	}

	if err != nil {
		return nil, statusCode, err
	}

	body := string(resp.Body())

	body = strings.ReplaceAll(body, "callback({", "{")
	body = strings.ReplaceAll(body, "})", "}")

	m := M{}

	json.Unmarshal([]byte(body), &m)

	result, err := a.abnModelFromMap(m)

	return result, statusCode, err
}

// AcnSearch for searching with acn code, the results are the ABN Details, status code and the error
func (a *Abn) AcnSearch(acn string) (*AbnModel, int, error) {
	err := a.validateGuid()

	statusCode := http.StatusInternalServerError

	if err != nil {
		return nil, statusCode, err
	}

	if acn == "" {
		return nil, statusCode, errors.New(a.message.AcnRequired)
	}

	client := resty.New()
	client.SetBaseURL(baseUrl)

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"guid": a.guid,
			"acn":  acn,
		}).
		Get(acnPath)

	if resp != nil {
		statusCode = resp.StatusCode()
	}

	if err != nil {
		return nil, statusCode, err
	}

	body := string(resp.Body())

	body = strings.ReplaceAll(body, "callback({", "{")
	body = strings.ReplaceAll(body, "})", "}")

	m := M{}

	json.Unmarshal([]byte(body), &m)

	result, err := a.abnModelFromMap(m)

	return result, statusCode, err
}

// NameSearch for searching with name, the results are the list of ABN Details, status code and the error
func (a *Abn) NameSearch(name string, maxResults int) ([]AbnSearchModel, int, error) {
	err := a.validateGuid()

	statusCode := http.StatusInternalServerError

	if err != nil {
		return nil, statusCode, err
	}

	if name == "" {
		return nil, statusCode, errors.New(a.message.NameRequired)
	}

	client := resty.New()
	client.SetBaseURL(baseUrl)

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"guid":       a.guid,
			"name":       name,
			"maxResults": strconv.Itoa(maxResults),
		}).
		Get(namePath)

	if resp != nil {
		statusCode = resp.StatusCode()
	}

	if err != nil {
		return nil, statusCode, err
	}

	body := string(resp.Body())

	body = strings.ReplaceAll(body, "callback({", "{")
	body = strings.ReplaceAll(body, "})", "}")

	m := M{}

	json.Unmarshal([]byte(body), &m)

	result, err := a.abnSearchModelFromMap(m)

	return result, statusCode, err
}

func (a *Abn) abnModelFromMap(abnMap map[string]interface{}) (*AbnModel, error) {
	msg := abnMap["Message"].(string)

	if msg != "" {
		return nil, errors.New(msg)
	}

	businessNames := []string{}

	for _, v := range abnMap["BusinessName"].([]interface{}) {
		businessNames = append(businessNames, v.(string))
	}

	m := &AbnModel{
		Abn:                    abnMap["Abn"].(string),
		Status:                 abnMap["AbnStatus"].(string),
		AbnStatusEffectiveFrom: abnMap["AbnStatusEffectiveFrom"].(string),
		Acn:                    abnMap["Acn"].(string),
		Address: AbnAddressModel{
			Date:     abnMap["AddressDate"].(string),
			Postcode: abnMap["AddressPostcode"].(string),
			State:    abnMap["AddressState"].(string),
		},
		BusinessNames: businessNames,
		Entity: AbnEntityModel{
			Name: abnMap["EntityName"].(string),
			Type: AbnEntityTypeModel{
				Code: abnMap["EntityTypeCode"].(string),
				Name: abnMap["EntityTypeName"].(string),
			},
		},
		Gst: abnMap["Gst"].(string),
	}

	return m, nil
}

func (a *Abn) abnSearchModelFromMap(abnMap map[string]interface{}) ([]AbnSearchModel, error) {
	msg := abnMap["Message"].(string)

	if msg != "" {
		return nil, errors.New(msg)
	}

	abns := []AbnSearchModel{}

	for _, v := range abnMap["Names"].([]interface{}) {
		vMap := v.(map[string]interface{})
		abns = append(abns, AbnSearchModel{
			Abn:       vMap["Abn"].(string),
			Status:    vMap["AbnStatus"].(string),
			IsCurrent: vMap["IsCurrent"].(bool),
			Name:      vMap["Name"].(string),
			NameType:  vMap["NameType"].(string),
			Postcode:  vMap["Postcode"].(string),
			Score:     vMap["Score"].(float64),
			State:     vMap["State"].(string),
		})
	}

	return abns, nil
}

func (a *Abn) AbnValidation(abn string) error {
	// check 11 digits
	if len(abn) != 11 {
		return errors.New(a.message.AbnInvalidLength)
	}

	// check is number
	regex := regexp.MustCompile(`^[0-9]+$`)
	match := regex.MatchString(abn)
	if !match {
		return errors.New(a.message.AbnInvalidType)
	}

	// weighting factor
	weights := []int{10, 1, 3, 5, 7, 9, 11, 12, 15, 17, 19}

	// sum = digit*weight
	sum := 0
	for i, w := range weights {
		n, _ := strconv.Atoi(string(abn[i]))

		// substract 1 from first digit
		if i == 0 {
			n -= 1
		}

		sum += (w * n)
	}

	// divide sum by 89, if remainder is zero then its valid
	remainder := math.Remainder(float64(sum), 89)
	if remainder > 0 {
		return errors.New(a.message.AbnInvalid)
	}

	return nil
}
