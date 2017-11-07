package tpaga

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Tpaga struct for authentication
type Tpaga struct {
	PublicKey        string
	publicKeyBase64  string
	PrivateKey       string
	privateKeyBase64 string

	// Define if is production or development
	isProduction bool

	// Define url api TPaga
	URLProduction  string
	URLDevelopment string
}

// City struct for the city of the customer
type City struct {
	Country string `json:"country,omitempty"`
	Name    string `json:"name,omitempty"`
	State   string `json:"state,omitempty"`
}

// Address struct for the address of the customer
type Address struct {
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	*City        `json:"city,omitempty"`
}

// Customer struct for the customers
type Customer struct {
	ID                 string `json:"id,omitempty"`
	Email              string `json:"email,omitempty"`
	FirstName          string `json:"firstName,omitempty"`
	LastName           string `json:"lastName,omitempty"`
	Gender             string `json:"gender,omitempty"`
	LegalIDNumber      string `json:"legalIdNumber,omitempty"`
	MerchantCustomerID string `json:"merchantCustomerId"`
	Phone              string `json:"phone,omitempty"`
	*Address           `json:"address,omitempty"`
}

// CreditCard struct for the credit cards
type CreditCard struct {
	CardHolderLegalIDNumber string `json:"cardHolderLegalIdNumber,omitempty"`
	CardHolderLegalIDType   string `json:"cardHolderLegalIdType,omitempty"`
	CardHolderName          string `json:"cardHolderName,omitempty"`
	CVC                     string `json:"cvc,omitempty"`
	ExpirationMonth         string `json:"expirationMonth,omitempty"`
	ExpirationYear          string `json:"expirationYear,omitempty"`
	PrimaryAccountNumber    string `json:"primaryAccountNumber,omitempty"`
}

// CreditCardToken struct with token from Tpaga
type CreditCardToken struct {
	Token string `json:"token"`
	Used  bool   `json:"used"`
}

// CreditCardPayment struct for create payment method
type CreditCardPayment struct {
	SkipLegalIDCheck bool   `json:"skipLegalIdCheck"`
	Token            string `json:"token"`
}

// CreditCardAssociated struct response from
// association credit card to costumer
type CreditCardAssociated struct {
	ID                      string `json:"id"`
	Bin                     string `json:"bin"`
	Type                    string `json:"type"`
	ExpirationMonth         string `json:"expirationMonth"`
	ExpirationYear          string `json:"expirationYear"`
	LastFour                string `json:"lastFour"`
	Customer                string `json:"customer"`
	CardHolderLegalIDNumber string `json:"cardHolderLegalIdNumber"`
	CardHolderLegalIDType   string `json:"cardHolderLegalIdType"`
	CardHolderName          string `json:"cardHolderName"`
	AddressLine1            string `json:"addressLine1"`
	AddressLine2            string `json:"addressLine2"`
	AddressCity             string `json:"addressCity"`
	AddressState            string `json:"addressState"`
	AddressPostalCode       string `json:"addressPostalCode"`
	AddressCountry          string `json:"addressCountry"`
	Fingerprint             string `json:"fingerprint"`
}

// ChargeRequest struct
type ChargeRequest struct {
	Amount       int    `json:"amount"`
	CreditCard   string `json:"creditCard"`
	Currency     string `json:"currency"`
	Installments int    `json:"installments"`
	OrderID      string `json:"orderId"`
	Description  string `json:"description"`
	IacAmount    int    `json:"iacAmount"`
	TaxAmount    int    `json:"taxAmount"`
	ThirdPartyID string `json:"thirdPartyId"`
	TipAmount    int    `json:"tipAmount"`
}

// TransactionInfo struct
type TransactionInfo struct {
	AuthorizationCode string `json:"authorizationCode"`
	Status            string `json:"status"`
}

// ResponseCharge struct
type ResponseCharge struct {
	ID                 string `json:"id"`
	Amount             string `json:"amount"`
	TaxAmount          string `json:"taxAmount"`
	NetAmount          string `json:"netAmount"`
	IacAmount          string `json:"iacAmount"`
	TipAmount          string `json:"tipAmount"`
	ReteRentaAmount    string `json:"reteRentaAmount"`
	ReteIvaAmount      string `json:"reteIvaAmount"`
	ReteIcaAmount      string `json:"reteIcaAmount"`
	TpagaFeeAmount     string `json:"tpagaFeeAmount"`
	Currency           string `json:"currency"`
	Paid               bool   `json:"paid"`
	Installments       int    `json:"installments"`
	OrderID            string `json:"orderId"`
	Description        string `json:"description"`
	DateCreated        string `json:"dateCreated"`
	ThirdPartyID       string `json:"thirdPartyId"`
	Customer           string `json:"customer"`
	CreditCard         string `json:"creditCard"`
	PaymentTransaction string `json:"paymentTransaction"`
	ErrorCode          string `json:"errorCode"`
	ErrorMessage       string `json:"errorMessage"`
	TransactionInfo    `json:"transactionInfo"`
}

// TpagaError struct for error from Tpaga
type TpagaError struct {
	Field         string `json:"field"`
	Message       string `json:"message"`
	Object        string `json:"object"`
	RejectedValue string `json:"rejected-value"`
}

// TpagaResponseError struct for handle errors from Tpaga
type TpagaResponseError struct {
	Errors []TpagaError `json:"errors"`
}

// NewV1 returns a instance of Tpaga struct
func NewV1(isProd bool) Tpaga {
	t := Tpaga{
		isProduction: isProd,
	}
	return t
}

// setPublicBase64 encode the public key to base64
func (t *Tpaga) setPublicBase64() {
	t.publicKeyBase64 = base64.StdEncoding.EncodeToString([]byte(t.PublicKey + ":"))
}

// setPrivateBase64 encode the private key to base64
func (t *Tpaga) setPrivateBase64() {
	t.privateKeyBase64 = base64.StdEncoding.EncodeToString([]byte(t.PrivateKey + ":"))
}

// CreateCustomer create a customer in Tpaga platform
func (t Tpaga) CreateCustomer(c Customer) (Customer, *TpagaResponseError) {

	bs, tre := t.requestPOST(c, "customer", true)
	if tre != nil {
		return c, tre
	}

	err := json.Unmarshal(bs, &c)
	if err != nil {
		tpe := TpagaError{Message: "Error al procesar el JSON" + err.Error()}
		tre.Errors = []TpagaError{tpe}
		return c, tre
	}
	return c, nil
}

// CreditCard creates a new Token which represents a CreditCard
// ************************************************************
// ** WARNING!!! **********************************************
// Don't use this method, this method must be used
// from your app client (Web (js) / App (Android / IOS))
func (t *Tpaga) CreditCard(c CreditCard) (CreditCardToken, *TpagaResponseError) {
	cct := CreditCardToken{}

	bs, tre := t.requestPOST(c, "tokenize/credit_card", false)
	if tre != nil {
		return cct, tre
	}

	err := json.Unmarshal(bs, &cct)
	if err != nil {
		tpe := TpagaError{Message: "Error al procesar el JSON" + err.Error()}
		tre.Errors = []TpagaError{tpe}
		return cct, tre
	}

	return cct, nil
}

// AssociateCreditCard associates a tokenized credit card to the customer
func (t *Tpaga) AssociateCreditCard(clientID, cardToken string) (CreditCardAssociated, *TpagaResponseError) {
	cca := CreditCardAssociated{}
	ccp := CreditCardPayment{
		SkipLegalIDCheck: false,
		Token:            cardToken,
	}

	bs, tre := t.requestPOST(ccp, "customer/"+clientID+"/credit_card_token", true)
	if tre != nil {
		return cca, tre
	}

	err := json.Unmarshal(bs, &cca)
	if err != nil {
		tpe := TpagaError{Message: "Error al procesar el JSON" + err.Error()}
		tre.Errors = []TpagaError{tpe}
		return cca, tre
	}

	return cca, nil
}

// Charge create a charge to a credit card
func (t *Tpaga) Charge(c ChargeRequest) (ResponseCharge, *TpagaResponseError) {
	rc := ResponseCharge{}

	bs, tre := t.requestPOST(c, "charge/credit_card", true)
	if tre != nil {
		return rc, tre
	}

	err := json.Unmarshal(bs, &rc)
	if err != nil {
		tpe := TpagaError{Message: "Error al procesar el JSON" + err.Error()}
		tre.Errors = []TpagaError{tpe}
		return rc, tre
	}

	return rc, nil
}

// requestPOST create a new Post Request to Tpaga platform
func (t *Tpaga) requestPOST(data interface{}, url string, private bool) ([]byte, *TpagaResponseError) {
	var urlTpaga string
	tre := &TpagaResponseError{}

	if t.isProduction {
		urlTpaga = t.URLProduction
	} else {
		urlTpaga = t.URLDevelopment
	}

	j, err := json.Marshal(data)
	if err != nil {
		tre.Errors = []TpagaError{TpagaError{Message: err.Error()}}
		return nil, tre
	}

	req, err := http.NewRequest("POST", urlTpaga+url, bytes.NewReader(j))
	if err != nil {
		tre.Errors = []TpagaError{TpagaError{Message: err.Error()}}
		return nil, tre
	}

	t.setPrivateBase64()
	t.setPublicBase64()

	req.Header.Set("Content-Type", "application/json")
	if private {
		req.Header.Set("Authorization", "Basic "+t.privateKeyBase64)
	} else {
		req.Header.Set("Authorization", "Basic "+t.publicKeyBase64)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tre.Errors = []TpagaError{TpagaError{Message: err.Error()}}
		return nil, tre
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tre.Errors = []TpagaError{TpagaError{Message: err.Error()}}
		return nil, tre
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return body, nil
	case http.StatusUnauthorized:
		tpe := TpagaError{Message: "No está autorizado, revise el token."}
		tre.Errors = []TpagaError{tpe}
		return nil, tre
	case http.StatusBadRequest:
		fallthrough
	case http.StatusUnprocessableEntity:
		err := json.Unmarshal(body, &tre)
		if err != nil {
			tre.Errors = []TpagaError{TpagaError{Message: err.Error()}}
		}

		return nil, tre
	}

	tre.Errors = []TpagaError{TpagaError{Message: "Error desconocido: " + resp.Status}}
	return nil, tre
}
