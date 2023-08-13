# ApproveDeny SDK for Go

![GitHub Actions](https://github.com/Ownage-FDT/approvedeny-sdk-go/actions/workflows/run-tests.yml/badge.svg)

The ApproveDeny SDK for Node.js provides an easy way to interact with the ApproveDeny API using Go.

## Installation
You can install the SDK using the `go get` command.

```bash
go get github.com/Ownage-FDT/approvedeny-sdk-go
```

## Usage
To use the SDK, you need to create an instance of the approve Client. You can do this by passing your API key to the constructor.

```go
import (
    "log"
    "github.com/Ownage-FDT/approvedeny-sdk-go"
)

 client, err := approvedeny.NewClient("your-api-key")

  if err != nil {
      log.Println(err)
  }
```

### Creating a new check request
To create a new check request, you need to call the `CreateCheckRequest` method on the client instance.
```go
requestPayload := approvedeny.CreateCheckRequestPayload{
    Description: "A description of the check request",
    Metadata: map[string]interface{}{
        "key": "value",
    },
}

checkRequest, err := client.CreateCheckRequest("check-id", requestPayload)

if err != nil {
  log.Println(err)
}

log.Println(checkRequest)
```

### Retrieving a check request
To retrieve a check request, you need to call the `GetCheckRequest` method on the client instance.
```go
checkRequestID := "check-request-id"

checkRequest, err := client.GetCheckRequest(checkRequestID)
if err != nil {
  log.Println(err)
}

log.Println(checkRequest)
```

### Retrieving a check request response
To retrieve a check request response, you need to call the `GetCheckRequestResponse` method on the client instance.
```go
checkRequestID := "check-request-id"

checkRequestResponse, err := client.GetCheckRequestResponse(checkRequestID)
if err != nil {
	log.Println(err)
}

log.Println(checkRequestResponse)
```

### Verifying webhook signatures
To verify webhook signatures, you need to call the `IsValidWebhookSignature` method on the client instance. This method returns a boolean value indicating whether the signature is valid or not.

```go

webhookPayload = approvedeny.WebhookPayload{
   	Event: "response.created",
		Data: map[string]interface{}{
			"foo": "bar",
		},
}

const isValidSignature = client.IsValidWebhookSignature("your-encryption-key", "signature", webhookPayload);

if isValidSignature {
  // The signature is valid
} else {
  // The signature is invalid
}
```

### Testing

```bash
go test
```

### Changelog

Please see [CHANGELOG](CHANGELOG.md) for more information what has changed recently.

## Contributing

Please see [CONTRIBUTING](CONTRIBUTING.md) for details.

### Security

If you discover any security related issues, please use the issue tracker.

## Credits

-   [Olayemi Olatayo](https://github.com/iamolayemi)
-   [All Contributors](../../contributors)

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
