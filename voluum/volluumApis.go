package voluum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bhambri94/voluum-apis/configs"

	config "github.com/bhambri94/voluum-apis/configs"
)

type AuthApiResponse struct {
	Token               string    `json:"token"`
	ExpirationTimestamp time.Time `json:"expirationTimestamp"`
	Inaugural           bool      `json:"inaugural"`
}

type AuthApiRequest struct {
	AccessID  string `json:"accessId"`
	AccessKey string `json:"accessKey"`
}

type DailyReport struct {
	TotalRows int `json:"totalRows"`
	Rows      []struct {
		CampaignID        string  `json:"campaignId"`
		CampaignName      string  `json:"campaignName"`
		Cost              float64 `json:"cost"`
		Revenue           float64 `json:"revenue"`
		TrafficSourceID   string  `json:"trafficSourceId"`
		TrafficSourceName string  `json:"trafficSourceName"`
	} `json:"rows"`
}

type CustomVariableReport struct {
	TotalRows int `json:"totalRows"`
	Rows      []struct {
		CampaignID         string  `json:"campaignId"`
		CampaignName       string  `json:"campaignName"`
		CustomVariable1    string  `json:"customVariable1"`
		CustomVariable1TS  string  `json:"customVariable1-TS"`
		CustomVariable10   string  `json:"customVariable10"`
		CustomVariable10TS string  `json:"customVariable10-TS"`
		Revenue            float64 `json:"revenue"`
		TrafficSourceID    string  `json:"trafficSourceId"`
		TrafficSourceName  string  `json:"trafficSourceName"`
	}
}

var VoluumApiAccessToken AuthApiResponse
var (
	CustomVariableTSUpdateDone = false
)

func getAccessToken() string {

	if VoluumApiAccessToken.Token != "" {
		return VoluumApiAccessToken.Token
	}

	authApiRequest := AuthApiRequest{
		AccessID:  config.Configurations.VoluumAccessId,
		AccessKey: config.Configurations.VoluumAccessKey,
	}

	byteArray, err := json.Marshal(authApiRequest)
	if err != nil {
		fmt.Println(err)
	}
	reader := bytes.NewReader(byteArray)
	fmt.Println("Calling Voluum Access Token api")
	req, err := http.NewRequest("POST", "https://api.voluum.com/auth/access/session", reader)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Accessid", config.Configurations.VoluumAccessId)
	req.Header.Set("Accesskey", config.Configurations.VoluumAccessKey)
	req.Header.Set("Authorization", "Basic dm9sdXVtZGVtb0B2b2x1dW0uY29tOjFxYXohUUFaIn0=")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(body, &VoluumApiAccessToken)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return VoluumApiAccessToken.Token
}

func GetVoluumReportsForMentionedDates(fromDate string, toDate string) (DailyReport, int) {
	token := getAccessToken()
	fmt.Println("Calling Get Vollum Report api for dates from: " + fromDate + " to: " + toDate)
	req, err := http.NewRequest("GET", "https://api.voluum.com/report?tz=America/Bogota&include="+configs.Configurations.IncludeTrafficSources+"&limit=10000&groupBy=traffic_source_id&groupBy=campaign_id&from="+fromDate+"&to="+toDate+"&column=traffic_source_id&column=traffic_source&column=campaign_id&column=campaign&column=cost&column=revenue", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Cwauth-Token", token)
	req.Header.Set("Authorization", "Basic dm9sdXVtZGVtb0B2b2x1dW0uY29tOjFxYXohUUFaIn0=")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var dailyReport DailyReport
	err = json.Unmarshal(body, &dailyReport)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return dailyReport, dailyReport.TotalRows
}

func GetVoluumReportsForCustomVariables(fromDate string, toDate string, ApiVariableName string, customVariableName string, TrafficSourceId string) CustomVariableReport {
	token := getAccessToken()
	fmt.Println("Calling Get Vollum Report api for custom variable from: " + fromDate + " to: " + toDate)

	req, err := http.NewRequest("GET", "https://api.voluum.com/report?tz=America/Bogota&groupBy=traffic_source_id&groupBy=campaign_id&groupBy="+ApiVariableName+"&from="+fromDate+"&to="+toDate+"&column=campaign_id&column=traffic_source_id&column=revenue&column="+ApiVariableName+"&column="+customVariableName+"&limit=10000&Include="+config.Configurations.IncludeTrafficSources+"&filter1=traffic_source_id&filter1Value="+TrafficSourceId, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Cwauth-Token", token)
	req.Header.Set("Authorization", "Basic dm9sdXVtZGVtb0B2b2x1dW0uY29tOjFxYXohUUFaIn0=")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var dailyRevenueReport CustomVariableReport
	err = json.Unmarshal(body, &dailyRevenueReport)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return dailyRevenueReport
}

func floatToString(inputNum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}

func createFinalReportForThisMonthData(dailyReport DailyReport, FinalRowsCount int, Day int, month string) ([][]interface{}, int) {
	var values [][]interface{}
	LocalDay := Day

	var firstRowOfSheetLabels []interface{}
	firstRowOfSheetLabels = append(firstRowOfSheetLabels, "Traffic Source Name", "Traffic Source ID", "Campaign Name", "Campaign ID", "CustomVariable-10TS")
	for LocalDay > 1 {
		firstRowOfSheetLabels = append(firstRowOfSheetLabels, "Cost - "+strconv.Itoa(LocalDay-1)+"/"+month, "Revenue - "+strconv.Itoa(LocalDay-1)+"/"+month)
		LocalDay--
	}
	values = append(values, firstRowOfSheetLabels)
	var secondBlankRow []interface{}
	secondBlankRow = append(secondBlankRow, "")
	values = append(values, secondBlankRow)

	ShortlistedTrafficSources = getShortlistedTrafficSources()
	fmt.Println("Preparing final sheet to be pushed to Google Sheets")

	rowID := 0
	for rowID < FinalRowsCount {
		LocalDay = Day
		if ShortlistedTrafficSources[strings.ToLower(dailyReport.Rows[rowID].TrafficSourceName)] || !config.Configurations.TrafficSourceFilteringEnabled && (strings.ToLower(dailyReport.Rows[rowID].TrafficSourceID) != strings.ToLower(config.Configurations.TSMappingViaCustomVariable.TrafficSourceId)) {
			var row []interface{}
			var customVariableTS string
			if val3, ok3 := finalMapCustomVariableTS[strings.ToLower(dailyReport.Rows[rowID].CampaignID)+strconv.Itoa(LocalDay)]; ok3 {
				customVariableTS = val3
			}
			row = append(row, dailyReport.Rows[rowID].TrafficSourceName, dailyReport.Rows[rowID].TrafficSourceID, dailyReport.Rows[rowID].CampaignName, dailyReport.Rows[rowID].CampaignID, customVariableTS)
			for LocalDay > 1 {
				var cost string
				var revenue string

				if val1, ok := finalMapCost[strings.ToLower(dailyReport.Rows[rowID].CampaignID)+strconv.Itoa(LocalDay)]; ok {
					if floatToString(val1) != "0.000000" {
						cost = "$" + floatToString(val1)
					} else {
						cost = ""
					}
				}
				if val2, ok2 := finalMapRevenue[strings.ToLower(dailyReport.Rows[rowID].CampaignID)+strconv.Itoa(LocalDay)]; ok2 {
					if floatToString(val2) != "0.000000" {
						revenue = "$" + floatToString(val2)
					} else {
						revenue = ""
					}
				}
				row = append(row, cost, revenue)
				LocalDay--
			}
			values = append(values, row)
		}
		rowID++
	}
	return values, rowID
}

var (
	finalMapCost                      = make(map[string]float64)
	finalMapRevenue                   = make(map[string]float64)
	finalMapCustomVariableTS          = make(map[string]string)
	finalRevenueMapCustomVariable10TS = make(map[string]float64)
	ShortlistedTrafficSources         = make(map[string]bool)
)

func getShortlistedTrafficSources() map[string]bool {
	configTrafficSources := config.Configurations.TrafficSourcesShortlisted
	for _, source := range configTrafficSources {
		ShortlistedTrafficSources[strings.ToLower(source)] = true
	}
	return ShortlistedTrafficSources
}

func addCostAndRevenueDayWiseToMap(dailyReport DailyReport, Day string, fromDate string, toDate string) {
	ShortlistedTrafficSources = getShortlistedTrafficSources()
	CustomVariableRevenueUpdateDone := false
	// CustomVariableTSUpdateDone := false

	fmt.Println("Saving Costs and Revenue Day wise")
	rowID := 0
	for rowID < len(dailyReport.Rows) {
		if ShortlistedTrafficSources[strings.ToLower(dailyReport.Rows[rowID].TrafficSourceName)] || !config.Configurations.TrafficSourceFilteringEnabled {
			if dailyReport.Rows[rowID].TrafficSourceName == config.Configurations.RevenueViaCustomVariable.Key && !CustomVariableRevenueUpdateDone {
				dailyRevenueReport := GetVoluumReportsForCustomVariables(fromDate, toDate, config.Configurations.RevenueViaCustomVariable.APIVariableName, config.Configurations.RevenueViaCustomVariable.CustomVariableName, config.Configurations.RevenueViaCustomVariable.TrafficSourceId)
				revenuerowID := 0
				for revenuerowID < len(dailyRevenueReport.Rows) {
					if dailyRevenueReport.Rows[revenuerowID].CustomVariable1TS == config.Configurations.RevenueViaCustomVariable.FieldName && IsValidCampaignId(dailyRevenueReport.Rows[revenuerowID].CustomVariable1) {
						finalMapRevenue[strings.ToLower(dailyRevenueReport.Rows[revenuerowID].CustomVariable1)+Day] = finalMapRevenue[strings.ToLower(dailyRevenueReport.Rows[revenuerowID].CustomVariable1)+Day] + dailyRevenueReport.Rows[revenuerowID].Revenue
					} else if dailyRevenueReport.Rows[revenuerowID].CustomVariable1TS != config.Configurations.RevenueViaCustomVariable.Key && dailyRevenueReport.Rows[revenuerowID].TrafficSourceName != config.Configurations.RevenueViaCustomVariable.Key {
						finalMapRevenue[strings.ToLower(dailyRevenueReport.Rows[revenuerowID].CampaignID)+Day] = finalMapRevenue[strings.ToLower(dailyRevenueReport.Rows[revenuerowID].CampaignID)+Day] + dailyRevenueReport.Rows[revenuerowID].Revenue
					} else if dailyRevenueReport.Rows[revenuerowID].CustomVariable1TS == config.Configurations.RevenueViaCustomVariable.FieldName && !IsValidCampaignId(dailyRevenueReport.Rows[revenuerowID].CustomVariable1) {
						finalMapRevenue[strings.ToLower(dailyRevenueReport.Rows[revenuerowID].CampaignID)+Day] = finalMapRevenue[strings.ToLower(dailyRevenueReport.Rows[revenuerowID].CampaignID)+Day] + dailyRevenueReport.Rows[revenuerowID].Revenue
					}
					revenuerowID++
				}
				CustomVariableRevenueUpdateDone = true
			} else if dailyReport.Rows[rowID].TrafficSourceName != config.Configurations.RevenueViaCustomVariable.Key {
				finalMapRevenue[strings.ToLower(dailyReport.Rows[rowID].CampaignID)+Day] = finalMapRevenue[strings.ToLower(dailyReport.Rows[rowID].CampaignID)+Day] + dailyReport.Rows[rowID].Revenue
			}
			finalMapCost[strings.ToLower(dailyReport.Rows[rowID].CampaignID)+Day] = dailyReport.Rows[rowID].Cost

			if dailyReport.Rows[rowID].TrafficSourceName == config.Configurations.TSMappingViaCustomVariable.Key && !CustomVariableTSUpdateDone {
				customVariableReport := GetVoluumReportsForCustomVariables(fromDate, toDate, config.Configurations.TSMappingViaCustomVariable.APIVariableName, config.Configurations.TSMappingViaCustomVariable.CustomVariableName, config.Configurations.TSMappingViaCustomVariable.TrafficSourceId)
				customVariableRowID := 0
				for customVariableRowID < len(customVariableReport.Rows) {
					if customVariableReport.Rows[customVariableRowID].CustomVariable10TS == config.Configurations.TSMappingViaCustomVariable.FieldName && IsValidCampaignId(customVariableReport.Rows[customVariableRowID].CustomVariable10) {
						finalMapCustomVariableTS[strings.ToLower(customVariableReport.Rows[customVariableRowID].CampaignID)+Day] = finalMapCustomVariableTS[strings.ToLower(customVariableReport.Rows[customVariableRowID].CampaignID)+Day] + ", " + customVariableReport.Rows[customVariableRowID].CustomVariable10
					}
					customVariableRowID++
				}
				CustomVariableTSUpdateDone = true
			}

		}
		rowID++
	}
}

func GetStandardVoluumReport() ([][]interface{}, int, string) {
	var finalValuesToSheet [][]interface{}
	var dailyReport DailyReport
	var RowCount int
	var monthYearDate string
	var EndOfMonthFlag bool
	var currentMonth string

	loc, _ := time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	// currentTime := time.Date(2020, time.July, 1, 18, 59, 59, 0, time.UTC) //This can be used to manually fill a sheet with from desired date
	currentDate := currentTime.Day()
	if currentDate == 1 {
		monthYearDate = currentTime.AddDate(0, -1, 0).Month().String() + strconv.Itoa(currentTime.Year()) //This will be used as Google Sheet name
		EndOfMonthFlag = true
		currentDate = 31
	} else {
		monthYearDate = currentTime.Month().String() + strconv.Itoa(currentTime.Year()) //This will be used as Google Sheet name
	}

	jdayIterator := 0
	for currentDate > 1 {
		fromDate := currentTime.AddDate(0, 0, jdayIterator-1).Format("2006-01-02T00")
		toDate := currentTime.AddDate(0, 0, jdayIterator).Format("2006-01-02T00")
		dailyReport, RowCount = GetVoluumReportsForMentionedDates(fromDate, toDate)
		addCostAndRevenueDayWiseToMap(dailyReport, strconv.Itoa(currentDate), fromDate, toDate)
		currentDate--
		jdayIterator--
	}
	if EndOfMonthFlag {
		fromDate := currentTime.AddDate(0, -1, 0).Format("2006-01-02T00")
		toDate := currentTime.Format("2006-01-02T00")
		dailyReport, RowCount = GetVoluumReportsForMentionedDates(fromDate, toDate)
		currentDate = currentTime.AddDate(0, 0, -1).Day() + 1
		currentMonth = strconv.Itoa(int(currentTime.AddDate(0, -1, 0).Month()))
	} else {
		fromDate := currentTime.AddDate(0, 0, -currentTime.Day()+1).Format("2006-01-02T00")
		toDate := currentTime.Format("2006-01-02T00")
		dailyReport, RowCount = GetVoluumReportsForMentionedDates(fromDate, toDate)
		currentDate = currentTime.Day()
		currentMonth = strconv.Itoa(int(currentTime.Month()))
	}

	finalValuesToSheet, RowCount = createFinalReportForThisMonthData(dailyReport, RowCount, currentDate, currentMonth)
	return finalValuesToSheet, RowCount, monthYearDate
}

func GetRevenueBasedOnCusomtVariable10() ([][]interface{}, int, string) {
	var finalValuesToSheet [][]interface{}
	var customVariableReport CustomVariableReport
	var RowCount int
	var monthYearDate string
	var EndOfMonthFlag bool

	loc, _ := time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	// currentTime := time.Date(2020, time.July, 1, 18, 59, 59, 0, time.UTC) //This can be used to manually fill a sheet with from desired date
	currentDate := currentTime.Day()
	if currentDate == 1 {
		monthYearDate = currentTime.AddDate(0, -1, 0).Month().String() + strconv.Itoa(currentTime.Year()) //This will be used as Google Sheet name
		EndOfMonthFlag = true
		currentDate = 31
	} else {
		monthYearDate = currentTime.Month().String() + strconv.Itoa(currentTime.Year()) //This will be used as Google Sheet name
	}

	jdayIterator := 0
	for currentDate > 1 {
		fromDate := currentTime.AddDate(0, 0, jdayIterator-1).Format("2006-01-02T00")
		toDate := currentTime.AddDate(0, 0, jdayIterator).Format("2006-01-02T00")
		customVariableReport = GetVoluumReportsForCustomVariables(fromDate, toDate, config.Configurations.TSMappingViaCustomVariable.APIVariableName, config.Configurations.TSMappingViaCustomVariable.CustomVariableName, config.Configurations.TSMappingViaCustomVariable.TrafficSourceId)
		addRevenueForCustomVaraibleDayWiseToMap(customVariableReport, strconv.Itoa(currentDate), fromDate, toDate)
		currentDate--
		jdayIterator--
	}
	if EndOfMonthFlag {
		currentDate = currentTime.AddDate(0, 0, -1).Day() + 1
	} else {
		currentDate = currentTime.Day()
	}

	finalValuesToSheet, RowCount = createFinalReportForCustomVariableData(currentDate)
	return finalValuesToSheet, RowCount, monthYearDate
}

var (
	UniqueCustomVariableValues [][]interface{}
	UniqueCustomVariableRows   = make(map[string]bool)
)

func addRevenueForCustomVaraibleDayWiseToMap(customVariableReport CustomVariableReport, Day string, fromDate string, toDate string) {
	fmt.Println("Saving Revenue based on custom variable 10")

	customVariableRowID := 0
	for customVariableRowID < len(customVariableReport.Rows) {
		if customVariableReport.Rows[customVariableRowID].CustomVariable10TS == config.Configurations.TSMappingViaCustomVariable.FieldName {
			if _, ok := UniqueCustomVariableRows[strings.ToLower(customVariableReport.Rows[customVariableRowID].CampaignID)+strings.ToLower(customVariableReport.Rows[customVariableRowID].CustomVariable10)]; !ok {
				var row []interface{}
				row = append(row, customVariableReport.Rows[customVariableRowID].TrafficSourceName, customVariableReport.Rows[customVariableRowID].TrafficSourceID, customVariableReport.Rows[customVariableRowID].CampaignName, customVariableReport.Rows[customVariableRowID].CampaignID, customVariableReport.Rows[customVariableRowID].CustomVariable10)
				UniqueCustomVariableValues = append(UniqueCustomVariableValues, row)
				UniqueCustomVariableRows[strings.ToLower(customVariableReport.Rows[customVariableRowID].CampaignID)+strings.ToLower(customVariableReport.Rows[customVariableRowID].CustomVariable10)] = true
			}
			finalRevenueMapCustomVariable10TS[strings.ToLower(customVariableReport.Rows[customVariableRowID].CampaignID)+strings.ToLower(customVariableReport.Rows[customVariableRowID].CustomVariable10)+Day] = finalRevenueMapCustomVariable10TS[strings.ToLower(customVariableReport.Rows[customVariableRowID].CampaignID)+strings.ToLower(customVariableReport.Rows[customVariableRowID].CustomVariable10)+Day] + customVariableReport.Rows[customVariableRowID].Revenue
		}
		customVariableRowID++
	}

}

func createFinalReportForCustomVariableData(Day int) ([][]interface{}, int) {
	var values [][]interface{}
	rowId := 0
	fmt.Println("Preparing final sheet For Custom Variable Data to be pushed to Google Sheets")

	for i := range UniqueCustomVariableValues {
		var row []interface{}
		var concatCampaignIdAndTS string
		for j := range UniqueCustomVariableValues[i] {
			row = append(row, UniqueCustomVariableValues[i][j])
			concatCampaignIdAndTS = strings.ToLower(UniqueCustomVariableValues[i][3].(string) + UniqueCustomVariableValues[i][4].(string))
		}
		LocalDay := Day
		for LocalDay > 1 {
			var cost string
			var revenue string
			if val, ok := finalRevenueMapCustomVariable10TS[concatCampaignIdAndTS+strconv.Itoa(LocalDay)]; ok {
				if floatToString(val) != "0.000000" {
					revenue = "$" + floatToString(val)
				} else {
					revenue = ""
				}
			}
			row = append(row, cost, revenue)
			LocalDay--
		}
		rowId++
		values = append(values, row)
	}
	return values, rowId
}

func IsValidCampaignId(u string) bool {
	length := len(u)
	if length > 16 {
		return true
	} else {
		return false
	}
}
