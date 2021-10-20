package storage

import (
	"testing"
	"time"
)

var url = "www.google.de/testing/plfpelfp?name=shorty&color=purple"

func TestInit(t *testing.T) {
	service, err := New()
	if err != nil {
		t.Errorf("Redis Init failed: %q", err)
	}
	service.Close()
}

func TestSaveLoad(t *testing.T) {
	service, err := New()
	if err != nil {
		t.Errorf("Redis Init failed: %q", err)
	}
	err = service.Save("123", url, time.Now().Add(time.Minute*time.Duration(1)), 0)
	if err != nil {
		t.Errorf("Redis Save failed: %q", err)
	}
	item, err := service.Load("123")
	if err != nil || item.URL != url {
		t.Errorf("Error Loading: %q", err)
	}
}
