package service

import (
	//"fmt"
	"testing"
	//"encoding/json"
	//lconf "git.lianjia.com/lianjia-sysop/lconf/handler"
	//"git.lianjia.com/lianjia-sysop/lconf/test_tools/utils"
)

func TestServiceAdd(t *testing.T) {
	if msg := ServiceAdd(); msg != "" {
		t.Fatal(msg)
	}
}

func TestSeviceRead(t *testing.T) {
	if msg := ServiceRead(); msg != "" {
		t.Fatal(msg)
	}
}

func TestServiceDelete(t *testing.T) {
	if msg := ServiceDelete(); msg != "" {
		t.Fatal(msg)
	}
}
