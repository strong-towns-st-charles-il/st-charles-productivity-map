package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

// all of the fields on https://kaneil.devnetwedge.com
type Parcel struct {
	ParcelNumber     string `json:"ParcelNumber"`
	SiteAddress      string `json:"SiteAddress"`
	Owner            string `json:"Owner"`
	TaxYear          string `json:"TaxYear"`
	SaleStatus       string `json:"SaleStatus"`
	PropertyClass    string `json:"PropertyClass"`
	TaxCode          string `json:"TaxCode"`
	TaxStatus        string `json:"TaxStatus"`
	NetTaxableValue  string `json:"NetTaxableValue"`
	TaxRate          string `json:"TaxRate"`
	TotalTax         string `json:"TotalTax"`
	Township         string `json:"Township"`
	Acres            string `json:"Acres"`
	MailingAddress   string `json:"MailingAddress"`
	LegalDescription string `json:"LegalDescription"`
}

// json format for 3D mapping
// Note: check type for Latitude, Longtitude, and Tax
type Data struct {
	PID       string  `json:"PID"`
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	LotSize   string  `json:"Loy_Size"`
	LotDesc   string  `json:"Lot_Desc"`
	Type      string  `json:"Type"`
	Tax       float64 `json:"Tax"`
}

func main() {
	// for loading the api key for geocoding - convert from address to lat/long
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	apiKey := os.Getenv("GEOCODING_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	filename := "sample-parcel.csv"
	path := "./data/"

	dat, err := os.ReadFile(path + filename)
	if err != nil {
		log.Fatalf("%s not found in %s: %v", filename, path, err)
	}

	years, parcelData := formatCSV(string(dat))
	var parcels []Parcel
	for i := range parcelData {
		URL := "https://kaneil.devnetwedge.com/parcel/view/" + parcelData[i] + "/" + years[i]
		p := scrapeParcels(URL)
		parcels = append(parcels, p)
	}

	// Write parsed data to file parcel.json
	writeJSON(parcels)

	// should be a new file here for main. Maybe split the scraper/parser and 3D Mapping Data?

	// Read the `parcel.json` file
	content, err := ioutil.ReadFile("./parcel.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Error here: cannot unmarshal array
	var payload Parcel
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
}

func formatCSV(table string) (years []string, parcels []string) {
	r := csv.NewReader(strings.NewReader(table))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		years = append(years, record[0])
		parcels = append(parcels, strings.Replace(string(record[1]), "-", "", -1))
	}
	return years[1:], parcels[1:]
}

func scrapeParcels(URL string) Parcel {
	labels := make([]string, 0)
	values := make([]string, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("kaneil.devnetwedge.com"),
	)

	c.OnHTML("table div.inner-label", func(e *colly.HTMLElement) {
		labels = append(labels, e.Text)
	})

	c.OnHTML("table div.inner-value", func(e *colly.HTMLElement) {
		values = append(values, e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.Visit(URL)
	parcel := labelDescJson(labels, values)
	return parcel
}

func labelDescJson(labels, values []string) (parcelInfo Parcel) {
	for i, value := range values {
		switch labels[i] {
		case "Parcel Number":
			parcelInfo.ParcelNumber = strings.Replace(value, "-", "", -1)
		case "Site Address":
			parcelInfo.SiteAddress = addressCleanUp(value)
		case "Owner Name & Address":
			parcelInfo.Owner = addressCleanUp(value)
		case "Tax Year":
			parcelInfo.TaxYear = strings.TrimSpace(value)[:4]
		case "Sale Status":
			parcelInfo.SaleStatus = strings.TrimSpace(value)
		case "Property Class":
			parcelInfo.PropertyClass = value
		case "Tax Code":
			parcelInfo.TaxCode = value
		case "Tax Status":
			parcelInfo.TaxStatus = value
		case "Net Taxable Value":
			parcelInfo.NetTaxableValue = value
		case "Tax Rate":
			parcelInfo.TaxRate = value
		case "Total Tax":
			parcelInfo.TotalTax = value
		case "Township":
			parcelInfo.Township = strings.TrimSpace(value)
		case "Acres":
			parcelInfo.Acres = value
		case "Mail Address":
			parcelInfo.MailingAddress = strings.TrimSpace(value)
		case "Legal Description (not for use in deeds or other transactional documents)":
			parcelInfo.LegalDescription = strings.TrimSpace(value)
		}
	}
	return
}

func addressCleanUp(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Replace(value, "\n", ",", -1)
	values := strings.Split(value, ",")
	value = ""
	for j := range values {
		values[j] = strings.TrimSpace(values[j])
	}
	return strings.Join(values, ", ")
}

func writeJSON(data []Parcel) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	_ = ioutil.WriteFile("parcel.json", file, 0644)
}
