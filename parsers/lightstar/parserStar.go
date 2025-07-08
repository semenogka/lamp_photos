package main

import (
	"encoding/json"
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

var lightstarSites = []SiteStruct {
	{URL: "https://lightstar.ru/kupit-lustry-svetilniki-optom/kupit-lustru-optom/?limit=all", Category: "alldata\\data\\lystri.json", Pages: 21},
	{URL: "https://lightstar.ru/kupit-lustry-svetilniki-optom/kupit-bra-optom/?limit=all", Category: "alldata\\data\\bra.json", Pages: 7},
	{URL: "https://lightstar.ru/kupit-lustry-svetilniki-optom/kupit-track-optom/trekovye-svetilniki/?limit=all", Category: "alldata\\data\\track.json", Pages: 11},
	{URL: "https://lightstar.ru/kupit-lustry-svetilniki-optom/trekovyye-sistemy-pro/trekovye-pro-svetilniki/?limit=all", Category: "alldata\\data\\track.json", Pages: 27},
	{URL: "https://lightstar.ru/kupit-lustry-svetilniki-optom/trekovyye-sistemy-linea/trekovye-linea-svetilniki/?limit=all", Category: "alldata\\data\\track.json", Pages: 7},
	{URL: "https://lightstar.ru/kupit-lustry-svetilniki-optom/trekovyye-sistemy-nove-led/trekovye-nove-svetilniki/?limit=all", Category: "alldata\\data\\track.json", Pages: 25},
	{URL: "https://lightstar.ru/kupit-lustry-svetilniki-optom/trekovyye-sistemy-uno/trekovye-uno-svetilniki/?limit=all", Category: "alldata\\data\\magn.json", Pages: 3},
	{URL: "https://lightstar.ru/kupit-lustry-svetilniki-optom/trekovyye-sistemy-due/trekovye-due-svetilniki/?limit=all", Category: "alldata\\data\\track.json", Pages: 11},
	{URL: "https://lightstar.ru/kupit-lustry-svetilniki-optom/kupit-nastolnye-lampi-optom/?limit=all", Category: "alldata\\data\\nastol.json", Pages: 15},
	{URL: "https://lightstar.ru/kupit-lustry-svetilniki-optom/kupit-torshery-optom/?limit=all", Category: "alldata\\data\\torsher.json", Pages: 15},
	{URL: "https://lightstar.ru/kupit-lustry-svetilniki-optom/vstraivaemye-tochechnye-svetilniki-pod-lampy/?limit=all", Category: "alldata\\data\\vstraivaem.json", Pages: 4},
	{URL: "https://lightstar.ru/kupit-lustry-svetilniki-optom/nakladnye-tochechnye-svetilniki-pod-lampy/?limit=all", Category: "alldata\\data\\tochnakl.json", Pages: 7},
	{URL: "https://lightstar.ru/kupit-lustry-svetilniki-optom/podvesnye-svetilniki/?limit=all", Category: "alldata\\data\\podves.json", Pages: 7},
	{URL: "https://lightstar.ru/kupit-lustry-svetilniki-optom/potolochnye-svetilniki/?limit=all", Category: "alldata\\data\\potol.json", Pages: 7},
}



func main() {
	for _, site := range lightstarSites {
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
			c.OnHTML(".specification__wrapper", func(e *colly.HTMLElement) {
				
				title := e.ChildText(".specification__term")
				if (strings.HasPrefix(title, "Высота min (H min), мм") ||  strings.HasPrefix(title, "Высота (H), мм")) && product.Height == "None" {
					value := e.ChildText(".specification__definition")
					product.Height = maxValue(value)
				}
				if (strings.HasPrefix(title, "Ширина (W), мм")) && product.Width == "None" {
					value := e.ChildText(".specification__definition")
					product.Width = maxValue(value)
				}
				if (strings.HasPrefix(title, "Диаметр (D), мм")) && product.Diameter == "None" {
					value := e.ChildText(".specification__definition")
					product.Diameter = maxValue(value)
				}
				if (strings.HasPrefix(title, "Длина (Глубина) (L), мм")) && product.Length == "None" {
					value := e.ChildText(".specification__definition")
					product.Length = maxValue(value)
				}
				if (strings.HasPrefix(title, "Цвет арматуры")) && product.ArmaturColor == "None" {
					value := e.ChildText(".specification__definition")
					product.ArmaturColor = replaceSeparatorsWithSlash(value)
				}
				if (strings.HasPrefix(title, "Цвет плафона")) && product.AbajurColor == "None" {
					value := e.ChildText(".specification__definition")
					product.AbajurColor = replaceSeparatorsWithSlash(value)
				}
				if (strings.HasPrefix(title, "Материал арматуры")) && product.ArmaturMaterial == "None" {
					value := e.ChildText(".specification__definition")
					product.ArmaturMaterial = replaceSeparatorsWithSlash(value)
				}
				if (strings.HasPrefix(title, "Материал плафона")) && product.AbajurMaterial == "None" {
					value := e.ChildText(".specification__definition")
					product.AbajurMaterial = replaceSeparatorsWithSlash(value)
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


func takeLinks(link string)([]string, []string) {
	c := colly.NewCollector()
	var links []string
	var names []string
	
	c.OnHTML(".grid__item", func(e *colly.HTMLElement) {
    	link := e.ChildAttr("header .card__title a ", "href")
    	fullLink := e.Request.AbsoluteURL(link)
    	log.Println("Ссылка на товар:", fullLink)

		imgURL := e.ChildAttr("header .card__image-wrapper img", "src")
		fullImgURL := e.Request.AbsoluteURL(imgURL)
		filename := path.Base(imgURL)

		downloadImg(fullImgURL, filename)
		links = append(links, link)
		names = append(names, filename)
		

		

		log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, filename)
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