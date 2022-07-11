package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	pa "main/pkg/poweradapter"
	iac "main/pkg/infraascode"
	"os"
	"strings"
)

func main() {
	infraFileType := os.Getenv("IACType")
	infraFileName := os.Getenv("IACTemplateFile")
	electricityMapZoneKey := os.Getenv("ELECTRICITY_MAP_AUTH_TOKEN")
	cloudProvider := os.Getenv("CloudProvider")
	CarbonRateProvider := os.Getenv("CARBON_RATE_PROVIDER") // electricitymap or watttime
	wattTimeUser := os.Getenv("WATT_TIME_USER")
	wattTimePass := os.Getenv("WATT_TIME_PASS")
	var averageKwh float64
	var Totalco2perkwh float64
	var count int
	var qry TypCloudResourceQuery

	githubNoticeMessage("Starting carbon measure action.")

	// TODO: For terraform, we might need to accept a list of multiple files
	var param pa.TypCarbonQueryParams
	param.IacProvider = cloudProvider
	param.CarbonRateProvider = CarbonRateProvider
	param.ElectricityMapZoneKey = electricityMapZoneKey
	param.WattTimeUser = wattTimeUser
	param.WattTimePass = wattTimePass

	sumary := iac.GetIACSummary(iac.TypIACQuery{Filetype: infraFileType, Filename: infraFileName})

	for _, ts := range sumary {

		Sizes := ts.Sizes
		if ts.Resource == "Microsoft.Compute/virtualMachines" {
			for _, S := range Sizes {
				for _, D := range S.Details {
					count = count + 1
					param.IacLocation = D.Location
					qry.Type = ts.Resource
					qry.Provider = cloudProvider
					qry.SizeName = S.Size
					qry.Location = D.Location
					averageKwh = averageKwh + getCarbonIntensity(param)

					Totalco2perkwh = Totalco2perkwh + ((getCarbonIntensity(param) * float64(GetWattage(qry))) / 1000)
				}
			}
		}
	}

	fmt.Println("grams_carbon_equivalent_per_kwh", averageKwh)
	fmt.Println("grams_emitted_over_24h", Totalco2perkwh)
	fmt.Println("Successfully ran carbon measure action.")

	gitHubOutputVariable("set-output name=grams_carbon_equivalent_per_kwh", fmt.Sprint(averageKwh))
	gitHubOutputVariable("set-output name=grams_emitted_over_24h", fmt.Sprint(Totalco2perkwh))
	githubNoticeMessage("Successfully ran carbon measure action.")
}

func getCarbonIntensity(param pa.TypCarbonQueryParams) float64 {
	x := pa.LiveCarbonIntensity(param)
	return x.LiveCarbonIntensity
}

func readJSON(jsonPath string) []TypCloudResources {
	file, _ := ioutil.ReadFile(jsonPath)
	var cloudLoc []TypCloudResources
	err := json.Unmarshal([]byte(file), &cloudLoc)
	if err != nil {
		fmt.Println(err.Error())
	}
	return cloudLoc
}

func GetWattage(qry TypCloudResourceQuery) (watt int) {
	cloudLoc := readJSON("references/resources.json")
	for _, c := range cloudLoc {
		if strings.ToLower(c.Cloud) == strings.ToLower(qry.Provider) {

			for _, l := range c.Resouce {
				if strings.ToLower(l.Type) == strings.ToLower(qry.Type) {

					for _, s := range l.Sizes {
						if strings.ToLower(s.Name) == strings.ToLower(qry.SizeName) {
							watt = s.Wattage
						}
					}
				}
			}
		}
	}
	return

}

type TypCloudResourceQuery struct {
	Provider string
	Location string
	SizeName string
	Type     string
}
type TypCloudResources struct {
	Cloud   string       `json:"cloud"`
	Resouce []TypResouce `json:"resources"`
}

type TypResouce struct {
	Type  string    `json:"type"`
	Sizes []TypSize `json:"sizes"`
}
type TypSize struct {
	Name    string `json:"Name"`
	Wattage int    `json:"wattage"`
}
