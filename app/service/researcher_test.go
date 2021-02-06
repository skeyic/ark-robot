package service

import (
	"flag"
	"github.com/golang/glog"
	"testing"
)

func Test_ResearcherInit(t *testing.T) {
	flag.Set("logtostderr", "true")
	flag.Set("v", "10")
	flag.Parse()

	glog.V(4).Infof("The RESEARCHER: %+v", TheResearcher)
}
