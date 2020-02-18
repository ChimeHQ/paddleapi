package paddleapi

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNewFulfillmentRequest(t *testing.T) {
	url := "https://test.com/something"
	bodyContent := "email=support%40chimehq.com&event_time=2019-08-30+15%3A34%3A56&marketing_consent=0&name=support&p_country=CA&p_coupon=&p_coupon_savings=0.00&p_currency=USD&p_earnings=%7B%22101267%22%3A%220.0000%22%7D&p_order_id=9181627&p_paddle_fee=0.00&p_price=0.00&p_product_id=569145&p_quantity=1&p_sale_gross=0.00&p_tax_amount=0.00&p_used_price_override=0&quantity=1&p_signature=k4oEEnk3vjXCTUczdaaFbLLjPnsbNisNwXuL2NApqyNJRCOqdh7zmjU2EtQq%2Bwe1GWaZxd0GmvEISpDWXccUmy%2FN3XFn3C59JbmXzRr%2FbnXK3J0cTj3jnZ7ovII3rwYcom0UcvKwC3lOZcUfUvBg8Ambmxz%2BYyV8vVRHb2FSOd0oEcmwb1J4L9%2B%2BMlLS6LUw9E34UGsw7AHZlwhQpxKV7XwdEavEv%2BGL4nB8soCiqvkYm3D3F7snRzLNlKrnSeIc3G%2BX%2FYc1JDQYN34KGHJWmOs6%2B793fm%2BDrpNfmN%2B5AZQS8v7tdbRtEZ4Im2bKDcrUvUklZ8bWIxAv6Ofsj7DEomQyn5b2D4iY8ekBVlQrjrqF02aiWqHyWnwd3tSkksZEaZBK92OKVbas0Xb8AZREgSrbJ0tHWqHBLxyDF6ImJWwr4zzegha11iuGHx3NfZkO3gWxua6%2BpRgtyU9WEk4rBD8fxgNdKj%2FJ4bCeFR3yKHzbA6heD1GGqFCEH5LdnW0Tmge5e3q1A3DcMOXO2dryij7vWmGO%2B2Z%2F5RUwUxyqhWDvhnvOHbcvjL0uakHqdb5gXBtN5eCT3%2BqVIvR2tNpJFAcvRj%2FXxLp47eEQ5guP3brgiJ70mcQTofwzuOk2nwdrHD0NVPFt6U%2Fzj0eTaBx9iqyIsXlRdXmTdODiUFl%2F4tw%3D"
	body := strings.NewReader(bodyContent)

	httpReq, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		t.Errorf("unable to create http request %s", err)
	}

	// this is necessary for correct form parsing
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	fulfilmentReq, err := NewFulfillmentRequest(httpReq)
	if err != nil {
		t.Errorf("unable to create fulfilment request %s", err)
		return
	}

	// This key is from a throw-away account
	const publicKey = `
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA2b7A8eStJXRrpIxmho05
hVCuMw12hu+Dm4uhlnnhQ4jSKLU+dlhDp9/wHgd74s3VODuwDL3jM43oabcGfBQ7
obwFxYsINIMShIuuv7TwplHsLLffTMoZhlyouUkroNhDcaZrehSaJPAVofBKc1tg
imLvWg+34UbR8Rh1CRRpIdzwCCtHL6VDCWgdCzAqD6AFRYYz2n9ZYAY1enGp125y
P12zYjmMwUayVGfpIMsDu22KKAwgP9Ilfk7K5AJrSiksGpXYCiy72AjeR9EusrC1
8F8S60Qc4bV/5hT9Qb/bnplwC7M0ydgjDV4ffCxccuMZpRw+wNX2u3mydqgjN6Qy
ztb3o+x4m4QuqqGs9K6f4HtCeTXyq34G2Jsbc8E3LKFulU805mv8kzMLaeh8jGDg
WYr1nLd2sA040cz5tLa4WM0PiJBDuTY4NuXNGOGijhci4g8rpnbwFWZ36Flewm2q
yQfM8QT0dUAldb41H1hTGxeVqmHUZMTieWwW+U6BdYbMBSy0rD2jiXjPkIgagxDx
oSJ5iby9YSMTt3e7Lt5ZNu0hQObFgwWKCGHqgpIWDrM5Dv4QuN4ptI8csgrXxq+n
HBrzfLy97R8jxxXBvt0Ah5yteSZTCUTzZDTDumqOVVSGSf/9HXrm83ocpQ2g2p6a
Jvdlm8I5nUhlRQTbz+S6csMCAwEAAQ==
-----END PUBLIC KEY-----
`

	if !fulfilmentReq.ValidSignature(publicKey) {
		t.Errorf("signature invalid")
	}

	if i, err := fulfilmentReq.OrderId(); err != nil || i != 9181627 {
		t.Errorf("order id not correct")
	}

	if c, err := fulfilmentReq.Country(); err != nil || c != "CA" {
		t.Errorf("country not correct")
	}

	if i, err := fulfilmentReq.Quantity(); err != nil || i != 1 {
		t.Errorf("quantity not correct")
	}

	if c, err := fulfilmentReq.Coupon(); err != nil || c != "" {
		t.Errorf("coupon not correct")
	}

	if c, err := fulfilmentReq.Currency(); err != nil || c != "USD" {
		t.Errorf("currency not correct")
	}

	if i, err := fulfilmentReq.ProductId(); err != nil || i != 569145 {
		t.Errorf("product id not correct")
	}
	
	if e, err := fulfilmentReq.Email(); err != nil || e != "support@chimehq.com" {
		t.Errorf("email not correct")
	}
	
	if n, err := fulfilmentReq.Name(); err != nil || n != "support" {
		t.Errorf("name not correct")
	}
	
	expectedDate := time.Date(2019, 8, 30, 15, 34, 56, 0, time.UTC)
	if et, err := fulfilmentReq.EventTime(); err != nil || et.Equal(expectedDate) == false {
		t.Errorf("time not correct %s", err.Error())
	}
}
