package config

import (
	"testing"
	"github.com/gwtony/dconf/test_tools/cases/service"
	"github.com/gwtony/dconf/test_tools/cases/group"
)

func TestConfigAddReadDelete(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigRead()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigDelete()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigReadNone()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}

func TestConfigUpdate(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigRead()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigUpdate()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigReadUpdate()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigDelete()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}

func TestConfigCopy(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigRead()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = group.GroupAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigCopy()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigReadCopy()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigDeleteCopy()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigDelete()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigReadNone()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}

func TestConfigCopyWildcard(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigAdd2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = group.GroupAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigCopyWildcard()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigReadCopyWildcard()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigDeleteCopyWildcard()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigDelete()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigDelete2()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigReadNoneCopyWildcard()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}

func TestConfigAddReadDeleteSlash(t *testing.T) {
	service.ServiceClean()
	msg := service.ServiceAdd()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigAddSlash()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigReadSlash()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigDeleteSlash()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = ConfigReadNoneSlash()
	if msg != "" {
		t.Fatal(msg)
	}
	msg = service.ServiceDelete()
	if msg != "" {
		t.Fatal(msg)
	}
}
