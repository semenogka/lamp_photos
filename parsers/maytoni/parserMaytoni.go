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

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type SiteStruct struct {
	URL      string
	Category string
	Pages    int
}

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

var maytoniSites = []SiteStruct {
	{URL: "https://maytoni.ru/catalog/decorative/lyustry/?SHOWALL=1#product_25", Category: "alldata\\data\\lystri.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/decorative/podvesy/?SHOWALL=1#product_25", Category: "alldata\\data\\podves.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/decorative/potolochnye-svetilniki/?SHOWALL=1#product_25", Category: "alldata\\data\\potol.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/decorative/torshery/?SHOWALL=1#product_25", Category: "alldata\\data\\torsher.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/decorative/nastennye-svetilniki/?SHOWALL=1#product_25", Category: "alldata\\data\\bra.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/decorative/nastolnye-svetilniki/?SHOWALL=1#product_25", Category: "alldata\\data\\nastol.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/decorative/podsvetka-decorative/?SHOWALL=1#product_25", Category: "alldata\\data\\podsvet.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/decorative/vstraivaemyy-svetilnik/", Category: "alldata\\data\\vstraivaem.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/street/outdoor-sistemy-osveshcheniya/ulichnaya-trekovaya-sistema-osveshcheniya-elasity-ip/?type[]=1719", Category: "alldata\\data\\track.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/street/podvesnye-svetilniki/", Category: "alldata\\data\\podves.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/street/bra/?SHOWALL=1#product_25", Category: "alldata\\data\\bra.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/street/potolochnye-svetilniki-outdoor/?SHOWALL=1#product_25", Category: "alldata\\data\\potol.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/street/landshaftnye-svetilniki/?SHOWALL=1#product_25", Category: "alldata\\data\\land.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/street/prozhektory/", Category: "alldata\\data\\projectors.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/street/vstraivaemye-svetilniki-street/", Category: "alldata\\data\\vstraivaem.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/trekovye-sistemy/odnofaznaya-trekovaya-sistema-unity/?type[]=80&SHOWALL=1", Category: "alldata\\data\\track.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/trekovye-sistemy/trekhfaznaya-trekovaya-sistema-trinity/?type%5B0%5D=102&SHOWALL=1#product_25", Category: "alldata\\data\\track.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/trekovye-sistemy/shtangovaya-sistema-osveshcheniya-axity/?type[]=1739", Category: "alldata\\data\\nastpot.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/trekovye-sistemy/magnitnaya-modulnaya-trekovaya-sistema-flarity/?type%5B0%5D=1728&SHOWALL=1#product_25", Category: "alldata\\data\\magn.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/trekovye-sistemy/gibkaya-trekovaya-sistema-elasity/?type[]=1746", Category: "alldata\\data\\track.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/trekovye-sistemy/gibkaya-trekovaya-sistema-flexity/?type[]=1712", Category: "alldata\\data\\track.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/trekovye-sistemy/magnitnaya-trekovaya-sistema-5mm-levity/?type%5B0%5D=1783&SHOWALL=1#product_25", Category: "alldata\\data\\magn.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/trekovye-sistemy/magnitnaya-trekovaya-sistema-exility/?type%5B0%5D=84&SHOWALL=1#product_25", Category: "alldata\\data\\magn.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/trekovye-sistemy/magnitnaya-trekovaya-sistema-radity/?type[]=352", Category: "alldata\\data\\magn.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/trekovye-sistemy/magnitnaya-trekovaya-sistema-gravity/?type[]=141&SHOWALL=1", Category: "alldata\\data\\magn.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/trekovye-sistemy/magnitnaya-trekovaya-sistema-s35/?type%5B0%5D=137&SHOWALL=1#product_25", Category: "alldata\\data\\magn.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/nastennye-svetilniki-func/?SHOWALL=1#product_25", Category: "alldata\\data\\bra.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/potolochnye-svetilniki-func/potolochnye-vstraivaemye-svetilniki/?type%5B0%5D=2032&SHOWALL=1#product_25", Category: "alldata\\data\\tochvstr.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/potolochnye-svetilniki-func/potolochnye-nakladnye-svetilniki/?type%5B0%5D=2034&SHOWALL=1#product_25", Category: "alldata\\data\\tochnakl.json", Pages: 5},
	{URL: "https://maytoni.ru/catalog/functional/potolochnye-svetilniki-func/potolochnye-podvesnye-svetilniki/?type%5B0%5D=1776&SHOWALL=1#product_25", Category: "alldata\\data\\tochpodv.json", Pages: 5},

}



func main() {
	//links := []string{}
	for _, site := range maytoniSites {
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
		
			c.OnHTML(".characteristic-list__item", func(e *colly.HTMLElement) {
				title := e.ChildText(".characteristic-list__item-name")
				
				product.Link = link
				product.Name = names[i]
				
				if title == "Высота" {
					value := e.ChildText(".characteristic-list__item-text")
					product.Height = maxValue(value)
				}
				if title == "Ширина" {
					value := e.ChildText(".characteristic-list__item-text")
					product.Width = maxValue(value)
				}
				if title == "Диаметр" {
					value := e.ChildText(".characteristic-list__item-text")
					product.Diameter = maxValue(value)
				}
				if title == "Длина" {
					value := e.ChildText(".characteristic-list__item-text")
					product.Diameter = maxValue(value)
				}
				if title == "Цвет арматуры" {
					value := e.ChildText(".characteristic-list__item-text")
					product.ArmaturColor = replaceSeparatorsWithSlash(value)
				}
				if title == "Цвет абажура" {
					value := e.ChildText(".characteristic-list__item-text")
					product.AbajurColor = replaceSeparatorsWithSlash(value)
				}
				if title == "Материал арматуры" {
					value := e.ChildText(".characteristic-list__item-text")
					product.ArmaturMaterial = replaceSeparatorsWithSlash(value)
				}
				if title == "Материал абажура" {
					value := e.ChildText(".characteristic-list__item-text")
					product.AbajurMaterial = replaceSeparatorsWithSlash(value)
				}
			})
			c.OnScraped(func(r *colly.Response) {
				log.Println(product.Link)
				products = append(products, product)
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

func takeLinks(link string) ([]string, []string) {
	var links []string
	var names []string

	c := colly.NewCollector()

	
	c.OnHTML(".catalog-card", func(e *colly.HTMLElement) {
		colors := e.DOM.Find(".catalog-card__colors")
		_, ok := colors.Find("a").Attr("href")
		if ok {
			
			colors.Find("a").Each(func(i int, s *goquery.Selection) {
				link, _ := s.Attr("href")
				fullLink := e.Request.AbsoluteURL(link)
				links = append(links, fullLink)
			})
			imgs := e.DOM.Find(".catalog-card__img .catalog-card__img-pages")
			
			imgs.Find("picture").Each(func(i int, s *goquery.Selection) {
				imgURL, _ := s.Find("img").Attr("src")
				fullImgURL := e.Request.AbsoluteURL(imgURL)
				filename := path.Base(imgURL)
				downloadImg(fullImgURL, filename)
				names = append(names, filename)
			})
		}else {
			link := e.ChildAttr(".catalog-card__link", "href")
			fullLink := e.Request.AbsoluteURL(link)
			links = append(links, fullLink)
			log.Println("Ссылка на товар:", fullLink)
			imgURL := e.ChildAttr(".catalog-card__img .catalog-card__img-pages picture img", "src")
			fullImgURL := e.Request.AbsoluteURL(imgURL)
			filename := path.Base(imgURL)
			downloadImg(fullImgURL, filename)
			names = append(names, filename)
		}

    	
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
