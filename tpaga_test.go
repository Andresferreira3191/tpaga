package tpaga

import (
	"fmt"
	"testing"
)

const (
	PublicKey  = "i28jhne7stt5mf1vg8gf5ckknqi6q041"
	PrivateKey = "eu1s29no9ktjvo8suef34m2tjp3gjg7f"
)

func TestTpaga_CreateCustomer(t *testing.T) {
	customer := Customer{
		Email:              "probando@gmail.com",
		FirstName:          "Prueba general",
		MerchantCustomerID: "15",
	}

	tp := NewV1(false)
	tp.PublicKey = PublicKey
	tp.PrivateKey = PrivateKey

	c, err := tp.CreateCustomer(customer)
	if err != nil {
		fmt.Printf("%#v", err)
		t.Errorf("creando un cliente en TPaga: %#v", err)
	}

	fmt.Printf("%#v", c)
}

func TestTpaga_CreditCard(t *testing.T) {
	card := CreditCard{
		CardHolderLegalIDNumber: "1012457854",
		CardHolderLegalIDType:   "CC",
		CardHolderName:          "ETHAN DAY",
		CVC:                     "943",
		ExpirationMonth:         "06",
		ExpirationYear:          "2019",
		PrimaryAccountNumber:    "5458414955998751",
	}

	tp := NewV1(false)
	tp.PublicKey = PublicKey
	tp.PrivateKey = PrivateKey

	c, err := tp.CreditCard(card)
	if err != nil {
		fmt.Printf("%#v", err)
		t.Errorf("tokenizando una tarjeta de crédito: %#v", err)
	}

	fmt.Printf("%#v", c)
}

func TestTpaga_AssociateCreditCard(t *testing.T) {
	client := "tp8svu96dteji3fpmlb12uv1e4qk5t5p"
	card := "9sh9rekodard5lvcbptpqp8pglss6g6l"

	tp := NewV1(false)
	tp.URLDevelopment = "https://sandbox.tpaga.co/api/"
	tp.PublicKey = PublicKey
	tp.PrivateKey = PrivateKey

	c, err := tp.AssociateCreditCard(client, card)
	if err != nil {
		fmt.Printf("%#v", err)
		t.Errorf("asociando una tarjeta de crédito a un cliente: %#v", err)
	}

	fmt.Printf("%#v", c)
}

func TestTpaga_Charge(t *testing.T) {
	charge := ChargeRequest{
		CreditCard:   "a3nlbdv0lh1u80n18g68sdft3va2j9b6",
		Currency:     "COP",
		Description:  "Prueba desde Golang",
		Amount:       25000,
		IacAmount:    200,
		Installments: 3,
		OrderID:      "FAC016",
		TaxAmount:    4800,
		ThirdPartyID: "25",
		TipAmount:    0,
	}

	tp := NewV1(false)
	tp.PublicKey = PublicKey
	tp.PrivateKey = PrivateKey

	c, err := tp.Charge(charge)
	if err != nil {
		fmt.Printf("%#v", err)
		t.Errorf("creando un cargo a la tarjeta de crédito de un cliente: %#v", err)
	}

	fmt.Printf("%#v", c)
}
