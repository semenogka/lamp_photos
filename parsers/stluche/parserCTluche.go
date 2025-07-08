package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

var stluceSites = []SiteStruct {
	{URL: "https://stluce.ru/catalog/dekorativnyy_svet/podvesnye_svetilniki_1/?PAGEN_1=%d", Category: "alldata\\data\\podves.json", Pages: 67},
	{URL: "https://stluce.ru/catalog/dekorativnyy_svet/potolochnye_svetilniki_3/?PAGEN_1=%d", Category: "alldata\\data\\potol.json", Pages: 17},
	{URL: "https://stluce.ru/catalog/dekorativnyy_svet/nastennye_svetilniki/?PAGEN_1=%d", Category: "alldata\\data\\bra.json", Pages: 25},
	{URL: "https://stluce.ru/catalog/dekorativnyy_svet/nastolnye_svetilniki/?PAGEN_1=%d", Category: "alldata\\data\\nastol.json", Pages: 7},
	{URL: "https://stluce.ru/catalog/dekorativnyy_svet/napolnye_svetilniki/?PAGEN_1=%d", Category: "alldata\\data\\torsher.json", Pages: 4},
	{URL: "https://stluce.ru/catalog/funktsionalnyy_svet/svetilniki_10/nastennye_svetilniki_3/?PAGEN_1=%d", Category: "alldata\\data\\bra.json", Pages: 6},
	{URL: "https://stluce.ru/catalog/funktsionalnyy_svet/svetilniki_10/podvesnye_svetilniki_3/?PAGEN_1=%d", Category: "alldata\\data\\podves.json", Pages: 8},
	{URL: "https://stluce.ru/catalog/funktsionalnyy_svet/svetilniki_10/potolochnye_svetilniki_2/vstraivaemye_1/?PAGEN_1=%d", Category: "alldata\\data\\tochvstr.json", Pages: 17},
	{URL: "https://stluce.ru/catalog/funktsionalnyy_svet/svetilniki_10/potolochnye_svetilniki_2/nakladnye_1/?PAGEN_1=%d", Category: "alldata\\data\\tochnakl.json", Pages: 11},
	{URL: "https://stluce.ru/catalog/ulichnyy_svet/landshaftnye_svetilniki/?PAGEN_1=%d", Category: "alldata\\data\\land.json", Pages: 4},
	{URL: "https://stluce.ru/catalog/ulichnyy_svet/nastennye_svetilniki_2/?PAGEN_1=%d", Category: "alldata\\data\\bra.json", Pages: 5},
	{URL: "https://stluce.ru/catalog/ulichnyy_svet/podvesnye_svetilniki_2/?PAGEN_1=%d", Category: "alldata\\data\\podves.json", Pages: 1},
	{URL: "https://stluce.ru/catalog/ulichnyy_svet/potolochnye_svetilniki_1/?PAGEN_1=%d", Category: "alldata\\data\\potol.json", Pages: 1},
}




func main() {
	for _, site := range stluceSites {
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
			c.OnHTML(".p-product__parameters-item", func(e *colly.HTMLElement) {
				title := e.ChildText(".p-product__parameters-name")
				if strings.HasPrefix(title, "Высота, мм") {
					value := e.ChildText(".p-product__parameters-value")
					product.Height = maxValue(value)
				}
				if strings.HasPrefix(title, "Ширина, мм") {
					value := e.ChildText(".p-product__parameters-value")
					product.Width = maxValue(value)
				}
				if strings.HasPrefix(title, "Длина, мм") {
					value := e.ChildText(".p-product__parameters-value")
					product.Length = maxValue(value)
				}
				if strings.HasPrefix(title, "Диаметр, мм") {
					value := e.ChildText(".p-product__parameters-value")
					product.Diameter = maxValue(value)
				}
				if strings.HasPrefix(title, "Материал плафона") {
					value := e.ChildText(".p-product__parameters-value")
					product.AbajurMaterial = replaceSeparatorsWithSlash(value)
				}
				if strings.HasPrefix(title, "Материал каркаса") {
					value := e.ChildText(".p-product__parameters-value")
					product.ArmaturMaterial = replaceSeparatorsWithSlash(value)
				}
				if strings.HasPrefix(title, "Цвет каркаса") {
					value := e.ChildText(".p-product__parameters-value")
					product.ArmaturColor = replaceSeparatorsWithSlash(value)
				}
				if strings.HasPrefix(title, "Цвет плафона") {
					value := e.ChildText(".p-product__parameters-value")
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




func extractName(url string) string {
	count := 0
	name := []rune{}
	for _, let := range url {
		if let == '/' {
			count += 1
		}

		if count == 6 && let != '/'{
			name = append(name, let)
		}
	}
	filename := fmt.Sprintf("%s.jpeg", string(name))
	return filename
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

func takeLinks(link string, pages int) ([]string, []string) {
	c := colly.NewCollector()
	var links []string
	var names []string
	
	c.OnHTML(".p-catalog__product-item", func(e *colly.HTMLElement) {
		link := e.ChildAttr("div a", "href")
		fullLink := e.Request.AbsoluteURL(link)
		//log.Println("Ссылка на товар:", fullLink)

		imgUrl := e.ChildAttr("div .product__image-wrapper .product__image-box  img", "data-src")
		fullImgURL := e.Request.AbsoluteURL(imgUrl)
		filename := fmt.Sprintf("CTLUCHE%s", extractName(fullImgURL))

		downloadImg(fullImgURL, filename)
		links = append(links, fullLink)
		names = append(names, filename)
		
		//log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, filename)
	})	

	for i := 1; i <= pages; i++{
		page := fmt.Sprintf(link, i)
		log.Println(page)
		err := c.Visit(page)
		if err != nil {
			log.Fatal(err)
		}
	}

	return links, names
}