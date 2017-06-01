package group

import (
	"testing"
	"github.com/gwtony/dconf/test_tools/cases/service"

)

func TestGroupAddReadDelete(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = GroupAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = GroupRead()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = GroupDelete()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = GroupReadNone()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}

func TestGroupList(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = GroupAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = GroupList()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = GroupDelete()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = GroupReadNone()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}

func TestGroupUpdate(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = GroupAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = GroupRead()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = GroupUpdate()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = GroupReadUpdate()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = GroupDelete()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = GroupReadNone()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}
