package text

import (
	"fmt"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/types"
	"time"
)

func Output(result chan types.Hostresult, size int) {
	cfg := options.GetCfg()
	for i := 0; i < size; i++ {
		select {
		case res := <-result:
			fmt.Printf("%+v", res)
		case <-time.After(time.Duration(cfg.Tasktimeout) * time.Second):
			break
		}
	}
}
