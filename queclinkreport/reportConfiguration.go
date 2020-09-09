package queclinkreport

import (
	"errors"
	"fmt"

	"queclink-go/base.device.service/report"

	"queclink-go/base.device.service/utils"
)

//ReportConfiguration root of report configuration
type ReportConfiguration struct {
	Devices []Device `xml:"device"`
}

//Device protocol description
type Device struct {
	Type              int                `xml:"type,attr"`
	Reports           []Report           `xml:"Reports>Report"`
	AckMessagesTypes  []AckMessagesType  `xml:"AckMessagesType>MessageType"`
	RespMessagesTypes []RespMessagesType `xml:"RespMessagesType>MessageType"`
}

//GetReport returns report by message header and message type
func (d *Device) GetReport(messageHeader string, messageType int) (*Report, error) {

	for _, r := range d.Reports {
		if r.ContainsClass(messageHeader) && r.ContainsType(messageType) {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("Not found report with class:%v; type:%v", messageHeader, messageType)
}

//FindRespMessagesType returns text report description
func (d *Device) FindRespMessagesType(command string) (*RespMessagesType, error) {
	for _, r := range d.RespMessagesTypes {
		if r.Command == command {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("Resp message type not found:%v", command)
}

//FindAckMessagesType returns ack report description
func (d *Device) FindAckMessagesType(command string) (*AckMessagesType, error) {
	for _, r := range d.AckMessagesTypes {
		if r.Command == command {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("Ack message type not found:%v", command)
}

//Report description
type Report struct {
	Location         bool            `xml:"Location,attr"`
	SupportedClasses []string        `xml:"SupportedClasses>Item"`
	SupportedTypes   []SupportedType `xml:"SupportedTypes>Item"`
	Header           []Item          `xml:"Header>ReportItem"`
	MultiPosition    []Item          `xml:"MultiPosition>ReportItem"`
	Tail             []Item          `xml:"Tail>ReportItem"`
}

//ContainsClass checks report contains message header
func (r *Report) ContainsClass(messageHeader string) bool {
	return utils.SliceContains(messageHeader, r.SupportedClasses)
}

//ContainsType checks report contains message header
func (r *Report) ContainsType(messageType int) bool {
	for _, v := range r.SupportedTypes {
		if v.Type == messageType {
			return true
		}
	}
	return false
}

//GetType search type by SupportedTypes
func (r *Report) GetType(messageType int32) (*SupportedType, bool) {
	for _, v := range r.SupportedTypes {
		if v.Type == int(messageType) {
			return &v, true
		}
	}
	return nil, false
}

//SupportedType struct
type SupportedType struct {
	Type              int  `xml:",chardata"`
	LastKnownPosition bool `xml:"lastKnownPosition,attr"`
}

//Item field description
type Item struct {
	ItemName string `xml:"ItemName"`
	Size     int    `xml:"Size"`
	MaskBit  *int   `xml:"MaskBit"`
	SizeIn   string `xml:"SizeIn"`
}

//AckMessagesType configuration item
type AckMessagesType struct {
	ID      int    `xml:"ID"`
	Command string `xml:"Command"`
}

//RespMessagesType configuration item
type RespMessagesType struct {
	ID      int    `xml:"ID"`
	Command string `xml:"Command"`
	Class   string `xml:"Class"`
}

//Find protocol description for device by message type and event
func (c *ReportConfiguration) Find(deviceType int, messageHeader string, messageType int) (*Report, error) {
	d, err := c.FindDeviceType(deviceType)
	if err != nil || d == nil {
		return nil, errors.New("Configuration not found")
	}
	return d.GetReport(messageHeader, messageType)
}

//FindDeviceType protocol description for device by message type and event
func (c *ReportConfiguration) FindDeviceType(deviceType int) (*Device, error) {
	for _, d := range c.Devices {
		if d.Type == deviceType {
			return &d, nil
		}
	}
	return nil, errors.New("Device configuration not found")
}

//LoadReportConfiguration proxy for base report
func LoadReportConfiguration(dir, fileDest string, instance interface{}) {
	report.LoadReportConfiguration(dir, fileDest, instance)
}
