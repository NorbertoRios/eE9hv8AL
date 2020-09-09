package models

import (
	"time"

	"queclink-go/base.device.service/config"
)

//DTCCodes represents device DTC codes
type DTCCodes struct {
	Codes               []string
	CodesUpdateDateTime time.Time
	LastRequest         time.Time
}

//NewDTCCodes returns new dtc codes struct
func NewDTCCodes() *DTCCodes {
	codes := &DTCCodes{}
	codes.Codes = []string{}
	return codes
}

//Count returns count of dtc codes
func (dtc *DTCCodes) Count() int {
	if dtc.Codes == nil {
		return 0
	}
	return len(dtc.Codes)
}

//NeedRefreshCodes checks ability to request DTC codes again
func (dtc *DTCCodes) NeedRefreshCodes() bool {
	interval := config.Config.GetBase().RefreshDtcIntervalH
	if time.Now().UTC().Before(dtc.CodesUpdateDateTime.Add(time.Duration(interval) * time.Hour)) {
		return false
	}

	if time.Now().UTC().Before(dtc.LastRequest.Add(5 * time.Minute)) {
		return false
	}
	dtc.LastRequest = time.Now().UTC()
	return true
}

//SetCodes assign codes to inner field
func (dtc *DTCCodes) SetCodes(codes []string) {
	if codes == nil {
		return
	}
	dtc.Codes = codes
	dtc.CodesUpdateDateTime = time.Now().UTC()
	dtc.LastRequest = time.Now().UTC().Add(1 * time.Hour)
}

//LoadCodes assign codes to inner field
func (dtc *DTCCodes) LoadCodes(codes []string) {
	if codes == nil {
		return
	}
	dtc.Codes = codes
}
