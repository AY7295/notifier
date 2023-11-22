package notifier

import (
	"fmt"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	tick := time.NewTicker(2 * time.Second)
	tt := time.NewTicker(20 * time.Second)
	fmt.Println(time.Now())
	for {
		select {
		case <-tick.C:
			t.Log("tick", time.Now())
		case <-tt.C:
			return
		}
	}
}
