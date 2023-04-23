[![Github CI](https://github.com/ChimeHQ/paddleapi/workflows/CI/badge.svg)](https://github.com/ChimeHQ/paddleapi/actions)

# paddleapi

Package paddleapi aims to provide a Go implementation of the [Paddle](https://paddle.com) API and Webhooks.

Paddle has a pretty featureful API, but right now this package only suports the [fulfillment webhook](https://developer.paddle.com/webhook-reference/product-fulfillment/fulfillment-webhook).

## Usage

```go
fulfillmentReq, err := NewFulfillmentRequest(httpReq)
if err != nil {
	return
}

if !fulfillmentReq.ValidSignature(publicKey) {
	return
}

productId, err := fulfillmentReq.ProductId()
if err != nil {
	return
}
```

### Suggestions or Feedback

We'd love to hear from you! Get in touch via an issue, or a pull request.

Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.
