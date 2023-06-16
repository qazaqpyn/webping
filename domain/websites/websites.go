package websites

import (
	"log"

	"github.com/qazaqpyn/webping/pkg/csvreader"
)

type Websites struct {
	list []string
}

func NewWebsites() (*Websites, error) {
	urls, err := csvreader.ReadCsvFile("./assets/websites.csv")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Websites{
		list: urls,
	}, nil
}

func (w *Websites) GetWebsites() []string {
	return w.list
}

func (w *Websites) CheckWebExist(url string) bool {
	for _, v := range w.list {
		if v == url {
			return true
		}
	}
	return false
}
