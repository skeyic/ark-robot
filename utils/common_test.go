package utils

import (
	"github.com/golang/glog"
	"testing"
)

func TestThousandFormatFloat64(t *testing.T) {
	EnableGlogForTesting()

	for _, num := range []float64{1234567890, -1234567890, -1234} {
		glog.V(4).Infof("%0.f is %s", num, ThousandFormatFloat64(num))
	}
}
