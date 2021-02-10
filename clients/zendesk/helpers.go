package zendesk

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

func setupRestyClient(email, password string) *resty.Client {
	// Setup resty client
	client := resty.New()

	// Setup Retries
	client.
		// Set retry count to non zero to enable retries
		SetRetryCount(3).
		// You can override initial retry wait time.
		// Default is 100 milliseconds.
		SetRetryWaitTime(5 * time.Second).
		// MaxWaitTime can be overridden as well.
		// Default is 2 seconds.
		SetRetryMaxWaitTime(20 * time.Second).
		// SetRetryAfter sets callback to calculate wait time between retries.
		// Default (nil) implies exponential backoff with jitter
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		})

	// Set retry condition
	client.AddRetryCondition(
		// RetryConditionFunc type is for retry condition function
		// input: non-nil Response OR request execution error
		func(r *resty.Response, err error) bool {
			return r.StatusCode() == http.StatusTooManyRequests
		},
	)

	// Assign Client Redirect Policy. Create one as per you need
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(15))

	// Set basic auth
	client.SetBasicAuth(email, password)

	return client
}
