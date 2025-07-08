package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Product struct {
	Name      string `json:"name"`
	Link      string `json:"link"`
	Height    string `json:"height"`
	Width    string `json:"width"`
	Length    string `json:"length"`
	Diameter    string `json:"diameter"`
	ArmaturColor    string `json:"metalcolor"`
	AbajurColor string `json:"abajurcolor"`
	ArmaturMaterial    string `json:"armaturmaterial"`
	AbajurMaterial string `json:"abajurmaterial"`
}

type SiteStruct struct {
	URL      string
	Category string
	Pages    int
}


var isonexSItes = []SiteStruct {
	{URL: "https://isonex.ru/catalog/filter/dev_site_category-is-%D0%BB%D1%8E%D1%81%D1%82%D1%80%D1%8B/apply/?ajax_request=Y&PAGEN_1=%d", Category: "alldata\\data\\lystri.json", Pages: 21},
	{URL: "https://isonex.ru/catalog/filter/dev_site_category-is-%D0%BF%D0%BE%D1%82%D0%BE%D0%BB%D0%BE%D1%87%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?ajax_request=Y&PAGEN_1=%d", Category: "alldata\\data\\potol.json", Pages: 7},
	{URL: "https://isonex.ru/catalog/filter/dev_naznachenie_na_sayte-is-%D0%BD%D0%B0%D1%81%D1%82%D0%B5%D0%BD%D0%BD%D0%BE-%D0%BF%D0%BE%D1%82%D0%BE%D0%BB%D0%BE%D1%87%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?ajax_request=Y&PAGEN_1=%d", Category: "alldata\\data\\nastpot.json", Pages: 11},
	{URL: "https://isonex.ru/catalog/filter/dev_site_category-is-%D0%BD%D0%B0%D1%81%D1%82%D0%B5%D0%BD%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8%20%D0%B8%20%D0%B1%D1%80%D0%B0/apply/?ajax_request=Y&PAGEN_1=%d", Category: "alldata\\data\\bra.json", Pages: 27},
	{URL: "https://isonex.ru/catalog/filter/dev_site_category-is-%D0%BD%D0%B0%D1%81%D1%82%D0%BE%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?ajax_request=Y&PAGEN_1=%d", Category: "alldata\\data\\nastol.json", Pages: 7},
	{URL: "https://isonex.ru/catalog/filter/dev_site_category-is-%D0%BF%D0%BE%D0%B4%D0%B2%D0%B5%D1%81%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?ajax_request=Y&PAGEN_1=%d", Category: "alldata\\data\\podves.json", Pages: 25},
	{URL: "https://isonex.ru/catalog/filter/dev_site_category-is-%D0%BD%D0%B0%D0%BF%D0%BE%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?ajax_request=Y&PAGEN_1=%d", Category: "alldata\\data\\torsher.json", Pages: 3},
	{URL: "https://isonex.ru/catalog/filter/dev_site_category-is-%D1%82%D1%80%D0%B5%D0%BA%D0%BE%D0%B2%D1%8B%D0%B5%20%D1%81%D0%B8%D1%81%D1%82%D0%B5%D0%BC%D1%8B%20220v/apply/?ajax_request=Y&PAGEN_1=%d", Category: "alldata\\data\\track.json", Pages: 11},
	{URL: "https://isonex.ru/catalog/filter/dev_site_category-is-%D0%BD%D0%B0%D0%BA%D0%BB%D0%B0%D0%B4%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?ajax_request=Y&PAGEN_1=%d", Category: "alldata\\data\\tochnakl.json", Pages: 15},
	{URL: "https://isonex.ru/catalog/filter/dev_site_category-is-%D0%B2%D1%81%D1%82%D1%80%D0%B0%D0%B8%D0%B2%D0%B0%D0%B5%D0%BC%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?ajax_request=Y&PAGEN_1=%d", Category: "alldata\\data\\tochvstr.json", Pages: 15},
	{URL: "https://isonex.ru/catalog/filter/dev_naznachenie_na_sayte-is-%D0%BB%D0%B0%D0%BD%D0%B4%D1%88%D0%B0%D1%84%D1%82%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?ajax_request=Y&PAGEN_1=%d", Category: "alldata\\data\\land.json", Pages: 4},
	{URL: "https://isonex.ru/catalog/filter/dev_naznachenie_na_sayte-is-%D1%83%D0%BB%D0%B8%D1%87%D0%BD%D1%8B%D0%B5%20%D1%84%D0%B0%D1%81%D0%B0%D0%B4%D0%BD%D1%8B%D0%B5/apply/?ajax_request=Y&PAGEN_1=%d", Category: "alldata\\data\\fasad.json", Pages: 7},

}







func main() {
	for _, site := range isonexSItes {
		links, names := takeLinks(site.URL, site.Pages)
		var products []Product

		for i, link := range links {
			c := colly.NewCollector()
			product := Product{
				Height:          "None",
				Width:          "None",
				Length: 		 "None",
				Diameter:        "None",
				ArmaturColor:    "None",
				AbajurColor:     "None",
				ArmaturMaterial: "None",
				AbajurMaterial:  "None",
			}
			product.Link = link
			product.Name = names[i]
			c.OnHTML(".characteristics-list-column", func(e *colly.HTMLElement) {
				e.DOM.Find("li").Each(func(i int, s *goquery.Selection) {
					title := s.Find(".text-muted")
					text := title.Text()
					
					if (strings.HasPrefix(text, "Высота без цепи, мм") || strings.HasPrefix(text, "Высота, мм")) && product.Height == "None" {
						value := title.Next().Find("span a")
						valueText := value.Text()
						product.Height = maxValue(valueText)
					}

					if strings.HasPrefix(text, "Диаметр, мм") {
						value := title.Next().Find("span a")
						valueText := value.Text()
						product.Diameter = maxValue(valueText)
					}

					if strings.HasPrefix(text, "Ширина, мм") {
						value := title.Next().Find("span a")
						valueText := value.Text()
						product.Width = maxValue(valueText)
					}

					if strings.HasPrefix(text, "Длина, мм") {
						value := title.Next().Find("span a")
						valueText := value.Text()
						product.Length = maxValue(valueText)
					}

					if strings.HasPrefix(text, "Материал арматуры") {
						value := title.Next().Find("span")
						valueText := value.Text()
						product.ArmaturMaterial =  replaceSeparatorsWithSlash(valueText)
					}
					if strings.HasPrefix(text, "Цвет арматуры") {
						value := title.Next().Find("span")
						valueText := value.Text()
						product.ArmaturColor =  replaceSeparatorsWithSlash(valueText)
					}
					if strings.HasPrefix(text, "Материал плафонов") {
						value := title.Next().Find("span")
						valueText := value.Text()
						product.AbajurMaterial =  replaceSeparatorsWithSlash(valueText)
					}
					if strings.HasPrefix(text, "Цвет плафона") {
						value := title.Next().Find("span")
						valueText := value.Text()
						product.AbajurColor = replaceSeparatorsWithSlash(valueText)
					}
				})
			})

			c.OnScraped(func(r *colly.Response) {
				log.Println(product)
				products = append(products, product)
			})

			c.Visit(link)
		}

		saveToJson(site.Category, products)
	}	
	
	

}


func downloadImg(src string, name string) {
	filepath := filepath.Join("allimgs", name)

	if _, err := os.Stat(filepath); err == nil {
		// Файл уже существует — не скачиваем
		log.Println("Файл уже существует, пропускаем:", name)
		return
	}

	res, err := http.Get(src)

	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	

	out, err := os.Create(filepath)
	if err != nil {
		log.Println(err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		log.Println(err)
	}
	
}

func saveToJson(filename string, data interface{}) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	var exData []interface{}
	
	stat, _ := file.Stat()
	if stat.Size() != 0 {
		json.NewDecoder(file).Decode(&exData)
	}

	exData = append(exData, data)


	if _, err := file.Seek(0, 0); err != nil {
		log.Println("Ошибка seek:", err)
		return
	}
	if err := file.Truncate(0); err != nil {
		log.Println("Ошибка truncate:", err)
		return
	}


	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	err = encoder.Encode(exData)

	if err != nil {
		log.Println(err)
	}

}

func extractName(url string) string {
	count := 0
	name := []rune{}
	for _, let := range url {
		if let == '/' {
			count += 1
		}

		if count == 8 && let != '/'{
			name = append(name, let)
		}
	}
	
	return string(name)
}

func takeLinks(link string, pages int) ([]string, []string) {
	c := colly.NewCollector()
	var links []string
	var names []string
	c.OnHTML(".product-card-inner", func(h *colly.HTMLElement) {
		
		
		url := h.ChildAttr(".flex-1 .px-1 a", "href")
		fullUrl := h.Request.AbsoluteURL(url)
		imgUrl := h.ChildAttr(".flex-1 .product-card-image .position-relative .list-card-slider img", "src")
		fullImgURL := h.Request.AbsoluteURL(imgUrl)
		filename := fmt.Sprintf("ISONEX%s", path.Base(imgUrl))
		links = append(links, fullUrl)
		names = append(names, filename)
		
		downloadImg(fullImgURL, filename)
		//log.Println("Ссылка: ", fullUrl, " Картинка: ", fullImgURL, " ", filename, " ")
		
		
	})	
	for i := 1; i <= pages; i++{
		temp := strings.ReplaceAll(link, "%", "%%")
		temp = strings.Replace(temp, "%%d", "%d", 1)
		page := fmt.Sprintf(temp, i)

		
		c.Visit(page)
	}

	return links, names
}

func maxValue(s string) string {
	parts := strings.Split(s, "/")
	if len(parts) == 1 {
		parts = strings.Split(s, "-")
		if len(parts) == 1 {
			parts = strings.Split(s, "–")
			if len(parts) == 1 {
				parts = strings.Split(s, "—")
			}
		}
	}
	max := 0
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		num, err := strconv.Atoi(extractNumber(part))
		
		if err == nil && num > max {
			max = num
		}
	}
	return strconv.Itoa(max)
}

func extractNumber(s string) string {
	re := regexp.MustCompile(`\d+`)
	return strings.Join(re.FindAllString(s, -1), "")
}

func replaceSeparatorsWithSlash(s string) string {
	reSep := regexp.MustCompile(`\s*(,|;)\s*`)
	s = reSep.ReplaceAllString(s, "/")

	reSlash := regexp.MustCompile(`\s*/\s*`)
	return reSlash.ReplaceAllString(s, "/")
}