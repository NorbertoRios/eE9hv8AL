package models

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"queclink-go/base.device.service/utils"
)

func TestMarshalingSoftware(t *testing.T) {
	str := "{\"FirmwareVersion\":\"X0z1-1137CD3\",\"OBDFirmwareVersion\":\"0402\",\"UpdatedAt\":\"2018-04-30T21:06:31Z\"}"
	software, _ := UnMarshalSoftware(str)
	if software.FirmwareVersion != "X0z1-1137CD3" {
		t.Error("Invalid FirmwareVersion value. Should:X0z1-1137CD3; Current:", software.FirmwareVersion)
	}

	if software.OBDFirmwareVersion != "0402" {
		t.Error("Invalid Platform value. Should:0402; Current:", software.OBDFirmwareVersion)
	}

	datetime, _ := time.Parse("2006-01-02T15:04:05Z", "2018-04-30T21:06:31Z")

	if software.UpdatedAt.Time != datetime {
		t.Error("Invalid UpdatedAt value. Should:2018-04-30T21:06:31Z; Current:", software.UpdatedAt.String())
	}

	jsoft := software.Marshal()
	fmt.Println(jsoft)
	if !strings.Contains(jsoft, "UpdatedAt\":\"2018-04-30T21:06:31Z\"}") {
		t.Error("Invalid marshaling for software. Should  contains:UpdatedAt\":\"2018-04-30T21:06:31Z\"}; Current:", jsoft)
	}
}

func TestSoftwareEqual(t *testing.T) {
	software := &Software{
		FirmwareVersion:    "X0z1-1137CD3",
		OBDFirmwareVersion: "0402",
		UpdatedAt:          &utils.JSONTime{Time: time.Now().UTC()}}
	equalSoftware := &Software{
		FirmwareVersion:    "X0z1-1137CD3",
		OBDFirmwareVersion: "0402",
		UpdatedAt:          &utils.JSONTime{Time: time.Now().UTC()}}
	if !software.compare(equalSoftware) {
		t.Error("Invalid compare software result")
	}
}

func TestSoftwareNotEqual(t *testing.T) {
	software := &Software{
		FirmwareVersion:    "X0z1-1137CD3",
		OBDFirmwareVersion: "0402",
		UpdatedAt:          &utils.JSONTime{Time: time.Now().UTC()}}
	notEqualSoftware := &Software{
		FirmwareVersion:    "X0z1-1137CD3",
		OBDFirmwareVersion: "0404",
		UpdatedAt:          &utils.JSONTime{Time: time.Now().UTC()}}
	if software.compare(notEqualSoftware) {
		t.Error("Invalid compare software result")
	}
}

func TestExpiredSoftware(t *testing.T) {
	expiredSoftware := &Software{
		FirmwareVersion:    "X0z1-1137CD3",
		OBDFirmwareVersion: "0402",
		UpdatedAt:          &utils.JSONTime{Time: time.Now().UTC().Add(-86500 * time.Second)}}
	software := &Software{
		FirmwareVersion:    "X0z1-1137CD3",
		OBDFirmwareVersion: "0402",
		UpdatedAt:          &utils.JSONTime{Time: time.Now().UTC()}}
	if expiredSoftware.compare(software) {
		t.Error("Invalid compare software result")
	}
}

func TestSoftwareNullValue(t *testing.T) {
	software := &Software{
		FirmwareVersion:    "X0z1-1137CD3",
		OBDFirmwareVersion: "0402",
		UpdatedAt:          nil}
	equalSoftware := &Software{
		FirmwareVersion:    "X0z1-1137CD3",
		OBDFirmwareVersion: "0402",
		UpdatedAt:          &utils.JSONTime{Time: time.Now().UTC()}}
	if !software.compare(equalSoftware) {
		t.Error("Invalid compare software result")
	}
}
