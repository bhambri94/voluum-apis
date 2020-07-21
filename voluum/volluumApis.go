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

var VoluumApiAccessToken AuthApiResponse

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
	fmt.Println(string(body))

	err = json.Unmarshal(body, &VoluumApiAccessToken)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return VoluumApiAccessToken.Token
}

func GetVoluumReportsForMentionedDates(fromDate string, toDate string) (DailyReport, int) {
	token := getAccessToken()

	req, err := http.NewRequest("GET", "https://api.voluum.com/report?include=ALL&limit=10000&groupBy=traffic_source_id&groupBy=campaign_id&from="+fromDate+"&to="+toDate+"&column=traffic_source_id&column=traffic_source&column=campaign_id&column=campaign&column=cost&column=revenue", nil)
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
	fmt.Println(string(body))

	var dailyReport DailyReport
	err = json.Unmarshal(body, &dailyReport)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	// values, rows := doCalculations(dailyReport)
	return dailyReport, dailyReport.TotalRows
}

func floatToString(inputNum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}

var finalValues [][]interface{}

func finalSheet(dailyReport DailyReport, FinalRowsCount int, Day int) ([][]interface{}, int) {
	var values [][]interface{}
	LocalDay := Day

	var firstRowOfSheetLabels []interface{}
	firstRowOfSheetLabels = append(firstRowOfSheetLabels, "Traffic Source Name", "Traffic Source ID", "Campaign Name", "Campaign ID")
	for LocalDay > 1 {
		firstRowOfSheetLabels = append(firstRowOfSheetLabels, "Cost-"+strconv.Itoa(LocalDay-1), "Revenue-"+strconv.Itoa(LocalDay-1))
		LocalDay--
	}
	values = append(values, firstRowOfSheetLabels)

	ShortlistedTrafficSources := make(map[string]bool)
	configTrafficSources := config.Configurations.TrafficSourcesShortlisted
	for _, source := range configTrafficSources {
		ShortlistedTrafficSources[strings.ToLower(source)] = true
	}

	rowID := 0
	for rowID < FinalRowsCount {
		LocalDay = Day
		if ShortlistedTrafficSources[strings.ToLower(dailyReport.Rows[rowID].TrafficSourceName)] || !config.Configurations.TrafficSourceFilteringEnabled {
			var row []interface{}
			row = append(row, dailyReport.Rows[rowID].TrafficSourceName, dailyReport.Rows[rowID].TrafficSourceID, dailyReport.Rows[rowID].CampaignName, dailyReport.Rows[rowID].CampaignID)
			for LocalDay > 1 {
				row = append(row, finalMapCost[dailyReport.Rows[rowID].CampaignID+strconv.Itoa(LocalDay)], finalMapRevenue[dailyReport.Rows[rowID].CampaignID+strconv.Itoa(LocalDay)])
				LocalDay--
			}
			fmt.Println(row)
			values = append(values, row)
		}
		rowID++
	}
	fmt.Println(values)
	return values, rowID
}

var (
	finalMapCost    = make(map[string]string)
	finalMapRevenue = make(map[string]string)
)

func addCostAndRevenueDayWise(dailyReport DailyReport, Day string) {
	ShortlistedTrafficSources := make(map[string]bool)
	configTrafficSources := config.Configurations.TrafficSourcesShortlisted
	for _, source := range configTrafficSources {
		ShortlistedTrafficSources[strings.ToLower(source)] = true
	}

	rowID := 0
	for rowID < len(dailyReport.Rows) {
		if ShortlistedTrafficSources[strings.ToLower(dailyReport.Rows[rowID].TrafficSourceName)] || !config.Configurations.TrafficSourceFilteringEnabled {
			finalMapCost[dailyReport.Rows[rowID].CampaignID+Day] = floatToString(dailyReport.Rows[rowID].Cost)
			finalMapRevenue[dailyReport.Rows[rowID].CampaignID+Day] = floatToString(dailyReport.Rows[rowID].Revenue)
		}
		rowID++
	}
}

func GetDailyVoluumReport() ([][]interface{}, int, string) {
	var values [][]interface{}
	var dailyReport DailyReport
	var HighestRows int
	// rows := 0
	finalRows := 0
	currentTime := time.Now()
	monthYearDate := currentTime.Month().String() + strconv.Itoa(currentTime.Year())
	currentDate := currentTime.Day()

	j := 0
	for currentDate > 1 {
		fromDate := currentTime.AddDate(0, 0, j-1).Format("2006-01-02T00")
		toDate := currentTime.AddDate(0, 0, j).Format("2006-01-02T00")
		dailyReport, HighestRows = GetVoluumReportsForMentionedDates(fromDate, toDate)
		if HighestRows > finalRows {
			finalRows = HighestRows
		}
		addCostAndRevenueDayWise(dailyReport, strconv.Itoa(currentDate))
		fmt.Println(fromDate, toDate)
		currentDate--
		j--
	}
	fromDate := currentTime.AddDate(0, 0, -currentTime.Day()+1).Format("2006-01-02T00")
	toDate := currentTime.Format("2006-01-02T00")
	dailyReport, HighestRows = GetVoluumReportsForMentionedDates(fromDate, toDate)
	currentDate = currentTime.Day()
	values, finalRows = finalSheet(dailyReport, HighestRows, currentDate)
	return values, finalRows, monthYearDate
}
