package main

import (
	config "github.com/bhambri94/voluum-apis/configs"
	"github.com/bhambri94/voluum-apis/sheets"
	"github.com/bhambri94/voluum-apis/voluum"
)

func main() {
	config.SetConfig()
	valuesFromDailyReport, _, SheetName := voluum.GetStandardVoluumReport()
	valuesFromCustomVariableReport, _, _ := voluum.GetRevenueBasedOnCusomtVariable10()
	values := concattwoInterfaces(valuesFromDailyReport, valuesFromCustomVariableReport)
	sheets.ClearSheet(SheetName)
	sheets.BatchWrite(SheetName, values)

}

func concattwoInterfaces(interface1 [][]interface{}, interface2 [][]interface{}) [][]interface{} {
	var finalValues [][]interface{}
	for i := range interface1 {
		finalValues = append(finalValues, interface1[i])
	}
	for i := range interface2 {
		finalValues = append(finalValues, interface2[i])
	}
	return finalValues
}
