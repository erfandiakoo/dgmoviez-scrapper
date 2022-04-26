package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
)

func main() {
	res := GetDomain()
	GetLastSeries(res)
}

type ResponseModel struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Sub    string `json:"sub"`
	URL    string `json:"url"`
	Cover  string `json:"cover"`
}
type ResponseModel2 struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code"`
}

func GetLastSeries(domain string) {
	res := []ResponseModel{}
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{domain},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			r.HTMLDoc.Find("div.item_dubbled").Each(func(i int, s *goquery.Selection) {
				if s.Find("div.text_holder_series").Text() != "" {
					res = append(res, ResponseModel{
						Number: i,
						Title:  s.Find("img").AttrOr("alt", ""),
						Sub:    clearText(s.Find("div.text_holder_series").Text()),
						URL:    s.Find("a").AttrOr("href", ""),
						Cover:  s.Find("img").AttrOr("src", "")})
				}
			})

			result := ResponseModel2{
				Data:    res,
				Message: r.Request.Host,
				Code:    r.StatusCode,
			}
			g.Exports <- result
		},
		Exporters: []export.Exporter{&export.JSON{}},
	}).Start()
}
func GetDomain() string {
	domain := ""
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://zil.ink/Digimoviez/"},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			r.HTMLDoc.Find("a").Each(func(i int, s *goquery.Selection) {
				if i == 2 {
					domain = s.AttrOr("href", "")
				}
			})
		},
	}).Start()
	return domain
}
func clearText(text string) string {
	cleanText := strings.TrimSpace(text)
	return cleanText
}
