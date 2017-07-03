package render

import (
	//"time"
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
	//print("sleep\n")
	//time.Sleep(time.Second * 120)
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
	msg = member.MemberDelete1()
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
	msg = member.MemberDelete1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}

func TestRenderDoKeyWildcard(t *testing.T) {
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
	msg = RenderDoKeyWildcard()
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
	msg = member.MemberDelete1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}

func TestRenderDoGroupWildcard(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = group.GroupAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigAdd2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigAddGroup()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberAdd1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberAdd2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberMove()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderDoGroupWildcard()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderRead2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderReadGroup()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderDelete2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderDeleteGroup()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigDelete2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigDeleteGroup()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberMoveBack()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberDelete1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberDelete2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}


func TestRenderDoGroupWildcardDelete(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = group.GroupAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigAdd2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigAddGroup()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberAdd1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberAdd2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberMove()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderDoGroupWildcard()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderRead2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderReadGroup()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigDeleteGroup()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderDoGroupWildcard()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderReadGroupNone()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderDelete2()
	if msg != "" {
		t.Fatal(msg)
	}

	msg = config.ConfigDelete2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberMoveBack()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberDelete1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberDelete2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}


func TestRenderDoGroupKeyWildcard(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = group.GroupAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigAdd2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigAddGroup()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberAdd1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberAdd2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberMove()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderDoGroupKeyWildcard()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderRead2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderReadGroupNone()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = RenderDelete2()
	if msg != "" {
		t.Fatal(msg)
	}
	//TODO: delete
	//msg = RenderDeleteGroup()
	//if msg != "" {
	//	t.Fatal(msg)
	//}
	msg = config.ConfigDelete2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = config.ConfigDeleteGroup()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberMoveBack()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberDelete1()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = member.MemberDelete2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}
