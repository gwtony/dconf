package render

import (
	"time"
	"testing"
	"github.com/gwtony/dconf/test_tools/cases/service"
	"github.com/gwtony/dconf/test_tools/cases/member"
	"github.com/gwtony/dconf/test_tools/cases/config"
)

func TestRenderDo(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberAdd1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberRead()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderDo()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderRead()
	if msg != "" {
		t.Fatal(msg)
	}
	print("sleep\n")
	time.Sleep(time.Second * 120)
	msg = RenderDelete()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderReadNone()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigDelete()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}


func TestRenderReadWildcard(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigAdd2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberAdd1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderDo()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderDo2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderReadWildcard()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderDeleteWildcard()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigDelete()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigDelete2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}
