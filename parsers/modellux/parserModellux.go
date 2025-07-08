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

var modelluxSites = []SiteStruct {
	{URL: "https://modelux.ru/catalog/liustra?page=%d", Category: "alldata\\data\\lystri.json", Pages: 39},
	{URL: "https://modelux.ru/catalog/potolocnyi-svetilnik?page=%d", Category: "alldata\\data\\potol.json", Pages: 30},
	{URL: "https://modelux.ru/catalog/podvesnoi-svetilnik?page=%d", Category: "alldata\\data\\podves.json", Pages: 55},
	{URL: "https://modelux.ru/catalog/bra?page=%d", Category: "alldata\\data\\bra.json", Pages: 40},
	{URL: "https://modelux.ru/catalog/nastolnaia-lampa?page=%d", Category: "alldata\\data\\nastol.json", Pages: 6},
	{URL: "https://modelux.ru/catalog/napolnyi-svetilnik?page=%d", Category: "alldata\\data\\torsher.json", Pages: 3},
	{URL: "https://modelux.ru/catalog/magnitnye-trekovye-svetilniki?page=%d", Category: "alldata\\data\\magn.json", Pages: 2},
}




func main() {
	for _, site := range modelluxSites {
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
			c.OnHTML(".product-tech__label", func(e *colly.HTMLElement) {
				title := e.DOM.Text()
				if (strings.HasPrefix(title, "Минимальная высота мм") || strings.HasPrefix(title, "Высота мм")) && product.Height == "None" {
					value := e.DOM.Next().Text()
					product.Height = maxValue(value)
				}
				if strings.HasPrefix(title, "Ширина мм") && product.Width == "None" {
					value := e.DOM.Next().Text()
					product.Width = maxValue(value)
				}
				if strings.HasPrefix(title, "Длина мм") && product.Length == "None" {
					value := e.DOM.Next().Text()
					product.Length = maxValue(value)
				}
				if strings.HasPrefix(title, "Диаметр мм") && product.Diameter == "None" {
					value := e.DOM.Next().Text()
					product.Diameter = maxValue(value)
				}
				if strings.HasPrefix(title, "Цвет арматуры") && product.ArmaturColor == "None" {
					value := e.DOM.Next().Text()
					product.ArmaturColor = replaceSeparatorsWithSlash(value)
				}
				if strings.HasPrefix(title, "Материал арматуры") && product.ArmaturMaterial == "None" {
					value := e.DOM.Next().Text()
					product.ArmaturMaterial = replaceSeparatorsWithSlash(value)
				}
				if strings.HasPrefix(title, "Цвет рассеивателя плафона") && product.AbajurColor == "None" {
					value := e.DOM.Next().Text()
					product.AbajurColor = replaceSeparatorsWithSlash(value)
				}
				if strings.HasPrefix(title, "Материал рассеивателя плафона") && product.AbajurMaterial == "None" {
					value := e.DOM.Next().Text()
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

func extractImgUrl(style string) string{
	runes := []rune{}
	count := false
	var s string
	for _, r := range style {
		if r == '(' {
			count = true
		}
		if count && r != '('{
			runes = append(runes, r)
		}
		if len(runes) >= 4 {
			s = string(runes[:len(runes)-3]) // удаляем последние 3
			s = s[1:]    // удаляем первый
		}
	}
	return s
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

func takeLinks(link string, pages int)([]string, []string) {
	c := colly.NewCollector()
	var links []string
	var names []string
	count := 0
	c.OnHTML(".product-item", func(e *colly.HTMLElement) {
		if (count == 1) {
			link := e.Attr("href")
			fullLink := e.Request.AbsoluteURL(link)
			log.Println("Ссылка на товар:", fullLink)
			
			imgSrc := e.ChildAttr("img.product-item__image-img", "src")
			filename := fmt.Sprintf("MODELLUX%s", path.Base(imgSrc))

			if (filename != "MODELLUXdefault.JPG" && filename != "MODELLUXproduct-default.40dea826.png") {
				downloadImg(imgSrc, filename)

				links = append(links, fullLink)
				names = append(names, filename)
				log.Println("Картинка:", imgSrc, " ",  "Ссылка на товар:", fullLink, filename)
			}
		}else {
			count += 1
		}
    	
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