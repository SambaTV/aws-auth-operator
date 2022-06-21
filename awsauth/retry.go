/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awsauth

import (
	"log"
	"time"

	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
)

const (
	defaultRetryerBackoffFactor float64 = 2.0
	defaultRetryerBackoffJitter         = true
)

// WithRetry runs the passed operation function with its arguments and retries
// on failures until success or max number of retry attempts have failed.
func WithRetry(fn func(*Arguments) error, args *Arguments) error {
	var (
		counter int
		err     error
		bkoff   = &backoff.Backoff{
			Min:    args.MinRetryTime,
			Max:    args.MaxRetryTime,
			Factor: defaultRetryerBackoffFactor,
			Jitter: defaultRetryerBackoffJitter,
		}
	)

	for {
		if counter >= args.MaxRetryCount {
			break
		}

		if err = fn(args); err != nil {
			d := bkoff.Duration()
			log.Printf("error: %v: will retry after %v", err, d)
			time.Sleep(d)
			counter++
			continue
		}
		return nil
	}
	return errors.Wrap(err, "waiter timed out")
}
