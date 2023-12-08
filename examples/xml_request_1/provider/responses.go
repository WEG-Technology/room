package provider

import "encoding/xml"

type TravelerinformationResponse struct {
	XMLName     xml.Name `xml:"TravelerinformationResponse"`
	Text        string   `xml:",chardata"`
	Xsd         string   `xml:"xsd,attr"`
	Xsi         string   `xml:"xsi,attr"`
	Page        string   `xml:"page"`
	PerPage     string   `xml:"per_page"`
	Totalrecord string   `xml:"totalrecord"`
	TotalPages  string   `xml:"total_pages"`
	Travelers   struct {
		Text                string `xml:",chardata"`
		Travelerinformation []struct {
			Text      string `xml:",chardata"`
			ID        string `xml:"id"`
			Name      string `xml:"name"`
			Email     string `xml:"email"`
			Adderes   string `xml:"adderes"`
			Createdat string `xml:"createdat"`
		} `xml:"Travelerinformation"`
	} `xml:"travelers"`
}
