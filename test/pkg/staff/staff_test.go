package staff

import (
	"github.com/wadeling/linear-webhook/pkg/staff"
	"testing"
)

var (
	testConfigFile = "./staff.config"
)

func TestLoadStaffConfig(t *testing.T) {
	t.Log("stat")
	s := staff.Instance()
	err := s.LoadStaffInfoFromFile(testConfigFile)
	if err != nil {
		t.Fatalf("load file err:%v", err)
	}
	t.Log("load end")

	s.DumpUserInfo()
	t.Log("end")
}
