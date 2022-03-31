package electricitymap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type electricityMap struct {
	zoneKey string
	url     string
}

func New(zoneKey string) electricityMap {
	em := electricityMap{
		zoneKey: zoneKey,
		url:     "https://api.electricitymap.org/v3",
	}
	return em
}

func httpQueryBuilder(zoneKey string, params TypAPIParams) (header map[string]string, query map[string]string) {
	header = make(map[string]string)
	query = make(map[string]string)

	header["auth-token"] = zoneKey

	if params.Zone != "" {
		query["zone"] = params.Zone
	}
	if params.Lon != "" && params.Lat != "" {
		query["lon"] = params.Lon
		query["lat"] = params.Lat
	}
	if params.Datetime != "" {
		query["datetime"] = params.Datetime
	}
	if params.Start != "" {
		query["start"] = params.Start
	}
	if params.End != "" {
		query["end"] = params.End
	}
	if params.EstimationFallback == true {
		query["estimationFallback"] = strconv.FormatBool(params.EstimationFallback)
	}

	return
}

/*
This endpoint returns all zones available if no auth-token is provided.

If an auth-token is provided, it returns a list of zones and routes available with this token
*/
func (e electricityMap) GetZones() (map[string]typZone, error) {
	url := fmt.Sprintf("%v/zones", e.url)
	data := make(map[string]typZone)
	header := make(map[string]string)
	query := make(map[string]string)

	header["auth-token"] = e.zoneKey

	fmt.Println("Getting Electricity Map Zones")
	err := httpGet(url, &data, header, query)
	return data, err
}

/*
This endpoint retrieves the last known carbon intensity (in gCO2eq/kWh) of electricity consumed in an area. It can either be queried by zone identifier or by geolocation.

QUERY PARAMETERS

Parameter | Description

zone | A string representing the zone identifier

lon | Longitude (if querying with a geolocation)

lat | Latitude (if querying with a geolocation)
*/
func (e electricityMap) LiveCarbonIntensity(params TypAPIParams) (typCI, error) {
	url := fmt.Sprintf("%v/carbon-intensity/latest", e.url)
	var data typCI

	header, query := httpQueryBuilder(e.zoneKey, params)

	fmt.Println("Getting Electricity Map Live Carbon Intensity")
	err := httpGet(url, &data, header, query)
	return data, err

}

/*
This endpoint retrieves the last known data about the origin of electricity in an area.

 - "powerProduction" (in MW) represents the electricity produced in the zone, broken down by production type

 - "powerConsumption" (in MW) represents the electricity consumed in the zone, after taking into account imports and exports, and broken down by production type.

 - "powerExport" and "Power import" (in MW) represent the physical electricity flows at the zone border

 - "renewablePercentage" and "fossilFreePercentage" refers to the % of the power consumption breakdown coming from renewables or fossil-free power plants (renewables and nuclear) It can either be queried by zone identifier or by geolocation.

QUERY PARAMETERS

Parameter | Description

zone | A string representing the zone identifier

lon | Longitude (if querying with a geolocation)

lat | Latitude (if querying with a geolocation)
*/
func (e electricityMap) LivePowerBreakdown(params TypAPIParams) (typPB, error) {
	url := fmt.Sprintf("%v/power-breakdown/latest", e.url)
	var data typPB

	header, query := httpQueryBuilder(e.zoneKey, params)

	fmt.Println("Getting Electricity Map Live Power Breakdown")
	err := httpGet(url, &data, header, query)
	return data, err

}

/*
This endpoint retrieves the last 24h of carbon intensity (in gCO2eq/kWh) of an area. It can either be queried by zone identifier or by geolocation. The resolution is 60 minutes.

QUERY PARAMETERS

Parameter | Description

zone | A string representing the zone identifier

lon | Longitude (if querying with a geolocation)

lat | Latitude (if querying with a geolocation)
*/
func (e electricityMap) RecentCarbonIntensity(params TypAPIParams) (typRecentCI, error) {
	url := fmt.Sprintf("%v/carbon-intensity/history", e.url)
	var data typRecentCI

	header, query := httpQueryBuilder(e.zoneKey, params)

	fmt.Println("Getting Electricity Map Recent Carbon Intensity")
	err := httpGet(url, &data, header, query)
	return data, err

}

/*
This endpoint retrieves the last 24h of power consumption and production breakdown of an area, which represents the physical origin of electricity broken down by production type. It can either be queried by zone identifier or by geolocation. The resolution is 60 minutes.

QUERY PARAMETERS

Parameter | Description

zone | A string representing the zone identifier

lon | Longitude (if querying with a geolocation)

lat | Latitude (if querying with a geolocation)
*/
func (e electricityMap) RecentPowerBreakdown(params TypAPIParams) (typRecentPB, error) {
	url := fmt.Sprintf("%v/power-consumption-breakdown/history", e.url)
	var data typRecentPB

	header, query := httpQueryBuilder(e.zoneKey, params)

	fmt.Println("Getting Electricity Map Recent Power Breakdown")
	err := httpGet(url, &data, header, query)
	return data, err

}

/*
This endpoint retrieves a past carbon intensity (in gCO2eq/kWh) of an area. It can either be queried by zone identifier or by geolocation. The resolution is 60 minutes.

QUERY PARAMETERS

Parameter | Description

zone | A string representing the zone identifier

lon | Longitude (if querying with a geolocation)

lat | Latitude (if querying with a geolocation)

datetime | datetime in ISO format

estimationFallback | (optional) boolean (if estimated data should be included)
*/
func (e electricityMap) PastCarbonIntensity(params TypAPIParams) (typCI, error) {
	url := fmt.Sprintf("%v/carbon-intensity/past", e.url)
	var data typCI

	header, query := httpQueryBuilder(e.zoneKey, params)

	fmt.Println("Getting Electricity Map Past Carbon Intensity")
	err := httpGet(url, &data, header, query)
	return data, err

}

/*
This endpoint retrieves a past carbon intensity (in gCO2eq/kWh) of an area within a given date range. It can either be queried by zone identifier or by geolocation. The resolution is 60 minutes. The time range is limited to 10 days.

QUERY PARAMETERS

Parameter | Description

zone | A string representing the zone identifier

lon | Longitude (if querying with a geolocation)

lat | Latitude (if querying with a geolocation)

start | datetime in ISO format

end | datetime in ISO format (excluded)

estimationFallback | (optional) boolean (if estimated data should be included)
*/
func (e electricityMap) PastCarbonIntensityRange(params TypAPIParams) (map[string][]typCI, error) {
	url := fmt.Sprintf("%v/carbon-intensity/past-range", e.url)
	var data = make(map[string][]typCI)

	header, query := httpQueryBuilder(e.zoneKey, params)

	fmt.Println("Getting Electricity Map Past Carbon Intensity Range")
	err := httpGet(url, &data, header, query)
	return data, err

}

/*
This endpoint retrieves a past power breakdown of an area. It can either be queried by zone identifier or by geolocation. The resolution is 60 minutes.

QUERY PARAMETERS

Parameter | Description

zone | A string representing the zone identifier

lon | Longitude (if querying with a geolocation)

lat | Latitude (if querying with a geolocation)

datetime | datetime in ISO format

estimationFallback | (optional) boolean (if estimated data should be included)
*/
func (e electricityMap) PastPowerBreakdown(params TypAPIParams) (typPB, error) {
	url := fmt.Sprintf("%v/power-breakdown/past", e.url)
	var data typPB

	header, query := httpQueryBuilder(e.zoneKey, params)

	fmt.Println("Getting Electricity Map Past Power Breakdown")
	err := httpGet(url, &data, header, query)
	return data, err

}

/*
This endpoint retrieves a past power breakdown of an area within a given date range. It can either be queried by zone identifier or by geolocation. The resolution is 60 minutes. The time range is limited to 10 days.

QUERY PARAMETERS

Parameter | Description

zone | A string representing the zone identifier

lon | Longitude (if querying with a geolocation)

lat | Latitude (if querying with a geolocation)

start | datetime in ISO format

end | datetime in ISO format (excluded)

estimationFallback | (optional) boolean (if estimated data should be included)
*/
func (e electricityMap) PastPowerBreakdownRange(params TypAPIParams) (map[string][]typPB, error) {
	url := fmt.Sprintf("%v/power-breakdown/past-range", e.url)
	var data = make(map[string][]typPB)

	header, query := httpQueryBuilder(e.zoneKey, params)

	fmt.Println("Getting Electricity Map Past Power Breakdown Range")
	err := httpGet(url, &data, header, query)
	return data, err

}

func httpGet(url string, data interface{}, header map[string]string, query map[string]string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if hasError(err) {
		fmt.Println("http.NewRequest error")
		fmt.Println(err.Error())
		return err
	}

	// Add Headers
	for k := range header {
		// fmt.Printf("Adding header %v:%v\n", k, header[k])
		req.Header.Add(k, header[k])
	}

	// Get URL Query String
	q := req.URL.Query()

	for k := range query {
		q.Add(k, query[k])
	}

	// Add query string to URL
	req.URL.RawQuery = q.Encode()

	// fmt.Println(req.URL)
	response, err := client.Do(req)
	if hasError(err) {
		fmt.Println("client.Do error")
		fmt.Println(err.Error())
		return err
	}

	if response.StatusCode == 200 {
		responseData, err := ioutil.ReadAll(response.Body)
		if hasError(err) {
			fmt.Println("ioutil.ReadAll error")
			fmt.Println(err.Error())
			return err
		}

		json.Unmarshal(responseData, &data)
		return nil //no error
	} else {
		err = errors.New(response.Status)
		return err
	}
}

func hasError(err error) bool {
	if err != nil {
		log.Fatal(err)
		return true
	}
	return false
}

type TypAPIParams struct {
	Zone               string
	Lon                string
	Lat                string
	Datetime           string
	Start              string
	End                string
	EstimationFallback bool
}

type typCI struct {
	Zone            string `json:"zone"`
	CarbonIntensity int    `json:"carbonIntensity"`
	Datetime        string `json:"datetime"`
	UpdatedAt       string `json:"updatedAt"`
	CreatedAt       string `json:"createdAt"`
}

type typPB struct {
	Zone                      string                       `json:"zone"`
	Datetime                  string                       `json:"datetime"`
	PowerProductionBreakdown  typPowerProductionBreakdown  `json:"powerProductionBreakdown"`
	PowerProductionTotal      int                          `json:"powerProductionTotal"`
	PowerConsumptionBreakdown typPowerConsumptionBreakdown `json:"powerConsumptionBreakdown"`
	PowerConsumptionTotal     int                          `json:"powerConsumptionTotal"`
	PowerImportBreakdown      typPowerImpExpBreakdown      `json:"powerImportBreakdown"`
	PowerImportTotal          int                          `json:"powerImportTotal"`
	PowerExportBreakdown      typPowerImpExpBreakdown      `json:"powerExportBreakdown"`
	PowerExportTotal          int                          `json:"powerExportTotal"`
	FossilFreePercentage      int                          `json:"fossilFreePercentage"`
	RenewablePercentage       int                          `json:"renewablePercentage"`
	UpdatedAt                 string                       `json:"updatedAt"`
	CreatedAt                 string                       `json:"createdAt"`
}

type typPowerConsumptionBreakdown struct {
	BatteryDischarge string // battery discharge `json:"batteryDischarge"`
	Biomass          int    `json:"biomass"`
	Coal             int    `json:"coal"`
	Gas              int    `json:"gas"`
	Geothermal       int    `json:"geothermal"`
	Hydro            int    `json:"hydro"`
	HydroDischarge   int    //hydro discharge `json:"hydroDischarge"`
	Nuclear          int    `json:"nuclear"`
	Oil              int    `json:"oil"`
	Solar            int    `json:"solar"`
	Unknown          int    `json:"unknown"`
	Wind             int    `json:"wind"`
}

type typPowerImpExpBreakdown struct {
	DE     int `json:"DE"`
	DK_DK1 int //DK-DK1 `json:"DK_DK1"`
	SE     int `json:"SE"`
}

type typPowerProductionBreakdown struct {
	Biomass    int `json:"biomass"`
	Coal       int `json:"coal"`
	Gas        int `json:"gas"`
	Geothermal int `json:"geothermal"`
	Hydro      int `json:"hydro"`
	Nuclear    int `json:"nuclear"`
	Oil        int `json:"oil"`
	Solar      int `json:"solar"`
	Unknown    int `json:"unknown"`
	Wind       int `json:"wind"`
}

type typZone struct {
	CountryName string   `json:"countryName"`
	ZoneName    string   `json:"zoneName"`
	Access      []string `json:"access"`
}

type typRecentCI struct {
	Zone    string `json:"zone"`
	History []struct {
		CarbonIntensity int    `json:"carbonIntensity"`
		Datetime        string `json:"datetime"`
		UpdatedAt       string `json:"updatedAt"`
		CreatedAt       string `json:"createdAt"`
	} `json:"history"`
}

type typRecentPB struct {
	Zone    string `json:"zone"`
	History []struct {
		Datetime                  string                       `json:"datetime"`
		FossilFreePercentage      string                       `json:"fossilFreePercentage"`
		PowerConsumptionBreakdown typPowerConsumptionBreakdown `json:"powerConsumptionBreakdown"`
		PowerConsumptionTotal     int                          `json:"powerConsumptionTotal"`
		PowerImportBreakdown      typPowerImpExpBreakdown      `json:"powerImportBreakdown"`
		PowerImportTotal          int                          `json:"powerImportTotal"`
		PowerExportBreakdown      typPowerImpExpBreakdown      `json:"powerExportBreakdown"`
		PowerExportTotal          int                          `json:"powerExportTotal"`
		PowerProductionBreakdown  typPowerProductionBreakdown  `json:"powerProductionBreakdown"`
		PowerProductionTotal      int                          `json:"powerProductionTotal"`
		RenewablePercentage       int                          `json:"renewablePercentage"`
	} `json:"history"`
}