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
)

type Parcel struct {
	ParcelNumber     string `json:"parcelNumber"`
	SiteAddress      string `json:"siteAddress"`
	Owner            string `json:"owner"`
	TaxYear          string `json:"taxYear"`
	SaleStatus       string `json:"saleStatus"`
	PropertyClass    string `json:"propertyClass"`
	TaxCode          string `json:"taxCode"`
	TaxStatus        string `json:"taxStatus"`
	NetTaxableValue  string `json:"netTaxableValue"`
	TaxRate          string `json:"taxRate"`
	TotalTax         string `json:"totalTax"`
	Township         string `json:"township"`
	Acres            string `json:"acres"`
	MailingAddress   string `json:"mailingAddress"`
	LegalDescription string `json:"legalDescription"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	dat, err := os.ReadFile("./data/sample-parcel.csv")
	check(err)
	years, parcelData := formatCSV(string(dat))
	var parcels []Parcel
	for i := range parcelData {
		URL := "https://kaneil.devnetwedge.com/parcel/view/" + parcelData[i] + "/" + years[i]
		p := scrapeParcels(URL)
		parcels = append(parcels, p)
	}

	writeJSON(parcels)

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
