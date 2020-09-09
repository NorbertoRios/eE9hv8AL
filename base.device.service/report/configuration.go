package report

import (
	"encoding/xml"
	"io/ioutil"
	"log"

	"queclink-go/base.device.service/utils"
)

//ReportConfiguration contains protocol description for devices
var ReportConfiguration IReportConfiguration

//LoadReportConfiguration load report configuration for device protocol
func LoadReportConfiguration(dir, fileDest string, configInstance IReportConfiguration) error {
	filePath := utils.GetAbsFilePath(dir, fileDest)
	log.Println("Loading report configuration from:", filePath)
	reportXML, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(reportXML, configInstance)
	ReportConfiguration = configInstance
	return err
}
