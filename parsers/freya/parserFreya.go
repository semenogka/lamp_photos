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
	"reflect"
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

var freayaSites = []SiteStruct {
	{URL: "https://freya-light.com/products/dekorativnyy_svet/lyustra/?SHOWALL_1=1", Category: "alldata\\data\\lystri.json", Pages: 0},
	{URL: "https://freya-light.com/products/dekorativnyy_svet/podvesnoy_svetilnik/?SHOWALL_1=1", Category: "alldata\\data\\podves.json", Pages: 0},
	{URL: "https://freya-light.com/products/dekorativnyy_svet/potolochnyy_svetilnik/?SHOWALL_1=1", Category: "alldata\\data\\potol.json", Pages: 0},
	{URL: "https://freya-light.com/products/dekorativnyy_svet/bra/?SHOWALL_1=1", Category: "alldata\\data\\bra.json", Pages: 0},
	{URL: "https://freya-light.com/products/dekorativnyy_svet/nastolnaya_lampa/?SHOWALL_1=1", Category: "alldata\\data\\nastol.json", Pages: 0},
	{URL: "https://freya-light.com/products/dekorativnyy_svet/torsher/?SHOWALL_1=1", Category: "alldata\\data\\torsher.json", Pages: 0},
}



func main() {
	for _, site := range freayaSites {
		var products []Product
		links, names := takeLinks(site.URL)

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
			c.OnHTML(".product-detail__chars-item", func(e *colly.HTMLElement) {
				name := e.DOM.Find(".product-detail__chars-item-name")
				text := name.Text()
				if strings.HasPrefix(text, "Мин. высота, мм") && product.Height == "None" {
					value := name.Next()
					product.Height = maxValue(value.Text())
				}
				if strings.HasPrefix(text, "Ширина, мм") && product.Width == "None" {
					value := name.Next()
					product.Width = maxValue(value.Text())
				}
				if strings.HasPrefix(text, "Длина, мм") && product.Length == "None" {
					value := name.Next()
					product.Length = maxValue(value.Text())
				}
				if strings.HasPrefix(text, "Диаметр светильника, мм") && product.Diameter == "None" {
					value := name.Next()
					product.Diameter = maxValue(value.Text())
				}
				if strings.HasPrefix(text, "Цвет") && product.ArmaturColor == "None" {
					value := name.Next()
					product.ArmaturColor = replaceSeparatorsWithSlash(value.Text())
				}
				if strings.HasPrefix(text, "Цвет абажура") && product.AbajurColor == "None" {
					value := name.Next()
					product.AbajurColor = replaceSeparatorsWithSlash(value.Text())
				}
				if strings.HasPrefix(text, "Материал абажура") && product.AbajurMaterial == "None" {
					value := name.Next()
					product.AbajurMaterial = replaceSeparatorsWithSlash(value.Text())
				}
				if strings.HasPrefix(text, "Материал арматуры") && product.ArmaturMaterial == "None" {
					value := name.Next()
					product.ArmaturMaterial = replaceSeparatorsWithSlash(value.Text())
				}
			})

			c.OnScraped(func(r *colly.Response) {
				products = append(products, product)
				log.Println(product)
			})

			err := c.Visit(link)
			if err != nil {
				log.Fatal(err)
			}
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

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			exData = append(exData, v.Index(i).Interface())
		}
	} else {
		exData = append(exData, data)
	}

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

func takeLinks(link string) ([]string, []string) {
	c := colly.NewCollector()
	var links []string
	var names []string
	c.OnHTML(".catalog-grid__item ", func(e *colly.HTMLElement) {
    	link := e.ChildAttr("a", "href")
    	fullLink := e.Request.AbsoluteURL(link)
    	log.Println("Ссылка на товар:", fullLink)

		imgURL := e.ChildAttr("a .catalog-card__picture picture source", "srcset")
		fullImgURL := e.Request.AbsoluteURL(imgURL)
		filename := fmt.Sprintf("FREYA%s", path.Base(imgURL))

		downloadImg(fullImgURL, filename)
		links = append(links, fullLink)
		names = append(names, filename)
		

		

		//log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, filename)
	})	

	err := c.Visit(link)

	if err != nil {
		log.Fatal(err)
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