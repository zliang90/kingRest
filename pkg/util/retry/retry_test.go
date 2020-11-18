package retry

import (
	"fmt"
	"testing"
	"time"

	"github.com/zliang90/kingRest/pkg/log"
)

func TestRetry(t *testing.T) {
	fn := func() error {
		log.Infof("no func args")
		return fmt.Errorf("timeout")
	}
	if err := Retry(1, 2*time.Second, fn); err != nil {
		log.Fatal(err)
	}
}

func TestRetryWithFuncArgs(t *testing.T) {
	// closure fn
	fn := func() error {
		name := "tom"

		// args function
		return func(name string) error {
			log.Infof("name: %s", name)

			return fmt.Errorf("timeout")

			// stop retry
			// return Stop{fmt.Errorf("stop")}
		}(name)
	}

	if err := Retry(1, 2*time.Second, fn); err != nil {
		log.Fatal(err)
	}
}
