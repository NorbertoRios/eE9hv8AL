package queclinkreport

import (
	"testing"
)

func TestConfigurationFind(t *testing.T) {
	config := &ReportConfiguration{}
	LoadReportConfiguration("..", "/ReportConfiguration.xml", config)
	report, err := config.Find(47, "+RSP", 5)
	if err != nil {
		t.Error("[TestConfigurationFind]Error retrive report configuration")
	}
	if report == nil {
		t.Error("[TestConfigurationFind]Error retrive report configuration. Report is nil")
	}
}

func TestConfigurationNotFound(t *testing.T) {
	config := &ReportConfiguration{}
	LoadReportConfiguration("..", "/ReportConfiguration.xml", config)
	report, err := config.Find(16, "+RSP", 125)
	if err == nil {
		t.Error("[TestConfigurationFind]Error retrive report configuration")
	}
	if report != nil {
		t.Error("[TestConfigurationFind]Error retrive report configuration. Report is nil")
	}
}

func TestConfigurationLastLocationAttr(t *testing.T) {
	config := &ReportConfiguration{}
	LoadReportConfiguration("..", "/ReportConfiguration.xml", config)
	report, _ := config.Find(47, "+EVT", 21)
	it, found := report.GetType(21)
	if !found {
		t.Error("[TestConfigurationLastLocationAttr]Error retrive report Type")
	}
	if !it.LastKnownPosition {
		t.Error("[TestConfigurationLastLocationAttr]Invalid value for LastKnownPosition")
	}
}
