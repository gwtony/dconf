package member

import (
	"testing"
	"github.com/gwtony/dconf/test_tools/cases/service"
	"github.com/gwtony/dconf/test_tools/cases/group"
)

func TestMemberAdd(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = MemberAdd1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = MemberRead()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}

func TestMemberDelete(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = MemberAdd1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = MemberRead()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = MemberDelete1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = MemberReadNone()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}

func TestMemberMove(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = MemberAdd1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = MemberRead()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = group.GroupAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = MemberMove()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = MemberReadTestgroup()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = MemberMoveBack()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = MemberRead()
	if msg != "" {
		t.Fatal(msg)
	}

	msg = MemberDelete1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}

