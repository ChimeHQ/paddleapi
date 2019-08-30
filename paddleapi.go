package paddleapi

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/elliotchance/phpserialize"
)

type FulfillmentRequest struct {
	fields map[string]string
}

const (
	FulfillmentSignatureField        = "p_signature"
	FulfillmentOrderIdField          = "p_order_id"
	FulfillmentCountryField          = "p_country"
	FulfillmentQuantityField         = "p_quantity"
	FulfillmentCouponField           = "p_coupon"
	FulfillmentCurrencyField         = "p_currency"
	FulfillmentProductIdField        = "p_product_id"
	FulfillmentMarketingConsentField = "marketing_consent"
	FulfillmentNameField             = "name"
	FulfillmentEmailField            = "email"
)

func NewFulfillmentRequest(req *http.Request) (*FulfillmentRequest, error) {
	err := req.ParseForm()
	if err != nil {
		return nil, err
	}

	fReq := FulfillmentRequest{fields: make(map[string]string)}

	for k, v := range req.Form {
		fReq.fields[k] = v[0]
	}

	return &fReq, nil
}

func (req FulfillmentRequest) FieldNamed(name string) string {
	return req.fields[name]
}

func (req FulfillmentRequest) ValidSignature(publicKey string) bool {
	encodedSignature := req.FieldNamed(FulfillmentSignatureField)

	signature, err := base64.StdEncoding.DecodeString(encodedSignature)
	if err != nil {
		fmt.Printf("could not decode signature \n%s\n", err.Error())
		return false
	}

	validatedFields := make(map[string]string)
	for k, v := range req.fields {
		validatedFields[k] = v
	}

	delete(validatedFields, FulfillmentSignatureField)

	serializedData, err := phpserialize.Marshal(validatedFields, nil)
	if err != nil {
		fmt.Printf("could not serialize validated fields")
		return false
	}

	return validateSignature(signature, publicKey, serializedData)
}

func validateSignature(signature []byte, publicKey string, data []byte) bool {
	hash := sha1.New()
	hash.Write(data)
	digest := hash.Sum(nil)

	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return false
	}

	parseResult, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false
	}

	switch parseResult.(type) {
	case *rsa.PublicKey:
		break
	default:
		return false
	}

	publcKey := parseResult.(*rsa.PublicKey)

	err = rsa.VerifyPKCS1v15(publcKey, crypto.SHA1, digest, signature)

	return err == nil
}

func (req FulfillmentRequest) OrderId() (int, error) {
	orderIdValue := req.FieldNamed(FulfillmentOrderIdField)
	i, err := strconv.ParseInt(orderIdValue, 10, 32)

	return int(i), err
}

func (req FulfillmentRequest) Country() (string, error) {
	country := req.FieldNamed(FulfillmentCountryField)

	var err error = nil

	if len(country) != 2 {
		err = errors.New("p_country value is invalid")
	}

	return country, err
}

func (req FulfillmentRequest) Quantity() (int, error) {
	value := req.FieldNamed(FulfillmentQuantityField)
	i, err := strconv.ParseInt(value, 10, 32)

	return int(i), err
}

func (req FulfillmentRequest) Coupon() (string, error) {
	return req.FieldNamed(FulfillmentCouponField), nil
}

func (req FulfillmentRequest) Currency() (string, error) {
	value := req.FieldNamed(FulfillmentCurrencyField)

	var err error = nil

	if len(value) != 3 {
		err = errors.New("p_currency value is invalid")
	}

	return value, err
}

func (req FulfillmentRequest) ProductId() (int, error) {
	value := req.FieldNamed(FulfillmentProductIdField)
	i, err := strconv.ParseInt(value, 10, 32)

	return int(i), err
}

func (req FulfillmentRequest) Name() (string, error) {
	return req.FieldNamed(FulfillmentNameField), nil
}

func (req FulfillmentRequest) Email() (string, error) {
	return req.FieldNamed(FulfillmentEmailField), nil
}
