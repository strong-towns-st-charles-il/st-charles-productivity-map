package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

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

func main() {
	labels := make([]string, 0)
	values := make([]string, 0)

	scrapeURL := "https://kaneil.devnetwedge.com/parcel/view/0927391001/2021"

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

	c.Visit(scrapeURL)
	parcel := labelDescJson(labels, values)
	writeJSON(parcel)

}

func labelDescJson(labels, values []string) (parcelInfo Parcel) {
	for i, value := range values {
		switch labels[i] {
		case "Parcel Number":
			parcelInfo.ParcelNumber = value
		case "Site Address":
			parcelInfo.SiteAddress = value
		case "Owner Name & Address":
			parcelInfo.Owner = value
		case "Tax Year":
			parcelInfo.TaxYear = value
		case "Sale Status":
			parcelInfo.SaleStatus = value
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
			parcelInfo.Township = value
		case "Acres":
			parcelInfo.Acres = value
		case "Mail Address":
			parcelInfo.MailingAddress = value
		case "Legal Description (not for use in deeds or other transactional documents)":
			parcelInfo.LegalDescription = value
		}
	}
	return
}

func writeJSON(data Parcel) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	_ = ioutil.WriteFile("parcel.json", file, 0644)
}
