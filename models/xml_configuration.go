package models

type Configuration struct {
	Assets Assets `xml:"Assets"`
}

type Assets struct {
	Outputs     Outputs           `xml:"Outputs"`
	Dataset     ElementMetaData   `xml:"Dataset"`
	Subject     string            `xml:"Subject"`
	Body        string            `xml:"Body"`
	Attachments []ElementMetaData `xml:"Attachments>Attachment"`
}

type Outputs struct {
	LogFile ElementMetaData `xml:"LogFile"`
}

type ElementMetaData struct {
	Name string `xml:"Name"`
	Type string `xml:"Type"`
	Path string `xml:"Path"`
}
