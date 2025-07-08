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

var loftitSites = []SiteStruct {
	{URL: "https://loftit.ru/catalog/3/filter/category-is-%D0%BB%D1%8E%D1%81%D1%82%D1%80%D1%8B/apply/?PAGEN_1=%d", Category: "alldata\\data\\lystri.json", Pages: 5},
	{URL: "https://loftit.ru/catalog/3/filter/category-is-%D0%B1%D1%80%D0%B0/apply/?PAGEN_1=%d", Category: "alldata\\data\\bra.json", Pages: 3},
	{URL: "https://loftit.ru/catalog/3/filter/category-is-%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/in_category-is-%D0%BD%D0%B0%D1%81%D1%82%D0%B5%D0%BD%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?PAGEN_1=%d", Category: "alldata\\data\\bra.json", Pages: 3},
	{URL: "https://loftit.ru/catalog/3/filter/category-is-%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/in_category-is-%D0%BF%D0%BE%D0%B4%D0%B2%D0%B5%D1%81%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?PAGEN_1=%d", Category: "alldata\\data\\podves.json", Pages: 25},
	{URL: "https://loftit.ru/catalog/3/filter/category-is-%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/in_category-is-%D0%BF%D0%BE%D1%82%D0%BE%D0%BB%D0%BE%D1%87%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?PAGEN_1=%d", Category: "alldata\\data\\potol.json", Pages: 7},
	{URL: "https://loftit.ru/catalog/3/filter/category-is-%D0%BD%D0%B0%D1%81%D1%82%D0%BE%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5%20%D0%BB%D0%B0%D0%BC%D0%BF%D1%8B/apply/?PAGEN_1=%d", Category: "alldata\\data\\nastol.json", Pages: 3},
	{URL: "https://loftit.ru/catalog/3/filter/category-is-%D1%82%D0%BE%D1%80%D1%88%D0%B5%D1%80%D1%8B/apply/?PAGEN_1=%d", Category: "alldata\\data\\torsher.json", Pages: 2},
	{URL: "https://loftit.ru/catalog/3/filter/in_category-is-%D1%82%D1%80%D0%B5%D0%BA%D0%BE%D0%B2%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?PAGEN_1=%d", Category: "alldata\\data\\magn.json", Pages: 1},
	{URL: "https://loftit.ru/catalog/3/filter/category-is-%D1%82%D0%B5%D1%85%D0%BD%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B9%20%D1%81%D0%B2%D0%B5%D1%82/in_category-is-%D0%B2%D1%81%D1%82%D1%80%D0%B0%D0%B8%D0%B2%D0%B0%D0%B5%D0%BC%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8-or-%D0%B2%D1%81%D1%82%D1%80%D0%B0%D0%B8%D0%B2%D0%B0%D0%B5%D0%BC%D1%8B%D0%B5%20%D1%82%D0%B5%D1%85%D0%BD%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?PAGEN_1=%d", Category: "alldata\\data\\tochvstr.json", Pages: 3},
	{URL: "https://loftit.ru/catalog/3/filter/category-is-%D1%82%D0%B5%D1%85%D0%BD%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B9%20%D1%81%D0%B2%D0%B5%D1%82/in_category-is-%D0%BD%D0%B0%D0%BA%D0%BB%D0%B0%D0%B4%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8-or-%D0%BD%D0%B0%D0%BA%D0%BB%D0%B0%D0%B4%D0%BD%D1%8B%D0%B5%20%D1%82%D0%B5%D1%85%D0%BD%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?PAGEN_1=%d", Category: "alldata\\data\\tochnakl.json", Pages: 2},
	{URL: "https://loftit.ru/catalog/3/filter/category-is-%D1%82%D0%B5%D1%85%D0%BD%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B9%20%D1%81%D0%B2%D0%B5%D1%82/in_category-is-%D0%BD%D0%B0%D1%81%D1%82%D0%B5%D0%BD%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?PAGEN_1=%d", Category: "alldata\\data\\bra.json", Pages: 4},
	{URL: "https://loftit.ru/catalog/3/filter/category-is-%D1%83%D0%BB%D0%B8%D1%87%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/in_category-is-%D0%BB%D0%B0%D0%BD%D0%B4%D1%88%D0%B0%D1%84%D1%82%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?PAGEN_1=%d", Category: "alldata\\data\\land.json", Pages: 2},
	{URL: "https://loftit.ru/catalog/3/filter/category-is-%D1%83%D0%BB%D0%B8%D1%87%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/in_category-is-%D1%83%D0%BB%D0%B8%D1%87%D0%BD%D1%8B%D0%B5%20%D0%BD%D0%B0%D1%81%D1%82%D0%B5%D0%BD%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?PAGEN_1=%d", Category: "alldata\\data\\bra.json", Pages: 1},
	{URL: "https://loftit.ru/catalog/3/filter/category-is-%D1%83%D0%BB%D0%B8%D1%87%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/in_category-is-%D1%83%D0%BB%D0%B8%D1%87%D0%BD%D1%8B%D0%B5%20%D0%BF%D0%BE%D0%B4%D0%B2%D0%B5%D1%81%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?PAGEN_1=%d", Category: "alldata\\data\\podves.json", Pages: 1},
}



func main() {
	for _, site := range loftitSites {
		var products []Product
		links, names := takeLinks(site.URL, site.Pages)
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
			c.OnHTML(".product__characteristic-item", func(e *colly.HTMLElement) {
				title := e.ChildText(".product__characteristic-name")

				if (strings.HasPrefix(title, "Минимальная высота, мм") || strings.HasPrefix(title, "Высота, мм")) && product.Height == "None" {
					value := e.ChildText(".product__characteristic-value")
					product.Height = maxValue(value)

				}
				if strings.HasPrefix(title, "Ширина, мм") && product.Width == "None" {
					value := e.ChildText(".product__characteristic-value")
					product.Width = maxValue(value)
				}
				if strings.HasPrefix(title, "Длина, мм") && product.Length == "None" {
					value := e.ChildText(".product__characteristic-value")
					product.Length = maxValue(value)
				}
				if strings.HasPrefix(title, "Диаметр, мм") && product.Diameter == "None" {
					value := e.ChildText(".product__characteristic-value")
					product.Diameter = maxValue(value)
				}
				if strings.HasPrefix(title, "Материал основания")  && product.ArmaturMaterial == "None" {
					value := e.ChildText(".product__characteristic-value")
					product.ArmaturMaterial = replaceSeparatorsWithSlash(value)
				}
				if strings.HasPrefix(title, "Материал плафона")  && product.AbajurMaterial == "None" {
					value := e.ChildText(".product__characteristic-value")
					product.AbajurMaterial = replaceSeparatorsWithSlash(value)
				}
				if strings.HasPrefix(title, "Цвет основания")  && product.ArmaturColor == "None" {
					value := e.ChildText(".product__characteristic-value")
					product.ArmaturColor = replaceSeparatorsWithSlash(value)
				}
				if strings.HasPrefix(title, "Цвет плафона")  && product.AbajurColor == "None" {
					value := e.ChildText(".product__characteristic-value")
					product.AbajurColor = replaceSeparatorsWithSlash(value)
				}
			})
			c.OnScraped(func(r *colly.Response) {
				products = append(products, product)
				log.Println(product)
			})

			c.Visit(link)
		}
		saveToJson(site.Category, products)
	}

	
	
}

func downloadImg(src string, name string) {
	res, err := http.Get(src)

	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	filepath := filepath.Join("allimgs", name)

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

func takeLinks(link string, pages int) ([]string, []string){
	c := colly.NewCollector()
	var links []string
	var names []string
	
	c.OnHTML(".catalog-card__item ", func(e *colly.HTMLElement) {
    	link := e.ChildAttr("a", "href")
    	fullLink := e.Request.AbsoluteURL(link)
    	log.Println("Ссылка на товар:", fullLink)

		imgURL := e.ChildAttr("a img", "src")
		fullImgURL := e.Request.AbsoluteURL(imgURL)
		filename := fmt.Sprintf("LOFTIT%s", path.Base(imgURL))

		downloadImg(fullImgURL, filename)
		links = append(links, fullLink)
		names = append(names, filename)
		

		

		//log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, filename)
	})	

	for i := 1; i <= pages; i++{
		temp := strings.ReplaceAll(link, "%", "%%")
		temp = strings.Replace(temp, "%%d", "%d", 1)
		page := fmt.Sprintf(temp, i)
		log.Println(page)
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