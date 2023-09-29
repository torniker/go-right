package repeat

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/torniker/go-right/pkg/logger"
)

type Repeater func(c context.Context) error

var Repeat string

func Tick(r Repeater) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if Repeat == "" {
			ctx := context.Background()
			err := r(ctx)
			if err != nil {
				// log.Error(ctx, err)
				fmt.Println(err)
			}
			return
		}

		repeatDuration, err := time.ParseDuration(Repeat)
		if err != nil {
			logger.Errorf("bad repeat duration: %s", Repeat)
			return
		}

		for {
			ctx := context.Background()
			err := r(ctx)
			if err != nil {
				logger.Error(err)
			}
			time.Sleep(repeatDuration)
		}
	}
}

