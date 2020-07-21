package main

import (
	"fmt"

	config "github.com/bhambri94/voluum-apis/configs"
	"github.com/bhambri94/voluum-apis/sheets"
	"github.com/bhambri94/voluum-apis/voluum"
)

func main() {
	config.SetConfig()
	fmt.Println(config.Configurations.SpreadsheetId)

	values, _, SheetName := voluum.GetDailyVoluumReport()
	sheets.ClearSheet(SheetName)
	writeRange := SheetName + "!A1"
	sheets.BatchWrite(writeRange, values)
}
