package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
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

var favouriteSites = []SiteStruct {
	{URL: "https://favourite-light.com/catalog/svetilniki_favourite/podvesnye_lyustry/?PAGEN_1=", Category: "alldata\\data\\lystri.json", Pages: 24},
	{URL: "https://favourite-light.com/catalog/svetilniki_favourite/potolochnye_lyustry/?PAGEN_1=", Category: "alldata\\data\\lystri.json", Pages: 7},
	{URL: "https://favourite-light.com/catalog/svetilniki_favourite/potolochnye_svetilniki/?PAGEN_1=", Category: "alldata\\data\\potol.json", Pages: 6},
	{URL: "https://favourite-light.com/catalog/svetilniki_favourite/podvesy/?PAGEN_1=", Category: "alldata\\data\\podves.json", Pages: 18},
	{URL: "https://favourite-light.com/catalog/svetilniki_favourite/bra/?PAGEN_1=", Category: "alldata\\data\\bra.json", Pages: 26},
	{URL: "https://favourite-light.com/catalog/svetilniki_favourite/nastolnye_lampy/?PAGEN_1=", Category: "alldata\\data\\nastol.json", Pages: 3},
	{URL: "https://favourite-light.com/catalog/svetilniki_favourite/torshery/?PAGEN_1=", Category: "alldata\\data\\torsher.json", Pages: 2},
	{URL: "https://favourite-light.com/catalog/funktsionalnyy_svet/nastennye_svetilniki/?PAGEN_1=", Category: "alldata\\data\\bra.json", Pages: 4},
	{URL: "https://favourite-light.com/catalog/funktsionalnyy_svet/podvesnye_svetilniki/?PAGEN_1=", Category: "alldata\\data\\podves.json", Pages: 1},
	{URL: "https://favourite-light.com/catalog/funktsionalnyy_svet/trekovye_sistemy/magnitnaya_sverkhtonkaya_trekovaya_sistema_unica/?PAGEN_1=", Category: "alldata\\data\\magn.json", Pages: 4},
	{URL: "https://favourite-light.com/catalog/funktsionalnyy_svet/trekovye_sistemy/magnitnaya_sverkhuzkaya_trekovaya_sistema_logica/?PAGEN_1=", Category: "alldata\\data\\magn.json", Pages: 3},
	{URL: "https://favourite-light.com/catalog/funktsionalnyy_svet/nakladnye_potolochnye_svetilniki/?PAGEN_1=", Category: "alldata\\data\\tochnakl.json", Pages: 1},
}



func main() {
	for _, site := range favouriteSites {
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
			c.OnHTML(".js-offers-prop", func(e *colly.HTMLElement) {
				e.DOM.Find("tr").Each(func(i int, s *goquery.Selection) {
					title := s.Find(".char_name .props_item span")
					text := cleanText(title.Text())
					
					if strings.HasPrefix(text, "Мин. Высота, мм") && product.Height == "None" {
						value := s.Find(".char_value span")
						valueText := maxValue(cleanText(value.Text()))
						product.Height = strconv.Itoa(valueText)
					}

					if strings.HasPrefix(text, "Диаметр") && product.Diameter == "None"{
						value := s.Find(".char_value span")
						valueText := maxValue(cleanText(value.Text()))
						product.Diameter = strconv.Itoa(valueText)
					}

					if strings.HasPrefix(text, "Ширина") && product.Width == "None" {
						value := s.Find(".char_value span")
						valueText := maxValue(cleanText(value.Text()))
						product.Width = strconv.Itoa(valueText)
					}

					if strings.HasPrefix(text, "Длина") && product.Length == "None" {
						value := s.Find(".char_value span")
						valueText := maxValue(cleanText(value.Text()))
						product.Length = strconv.Itoa(valueText)
					}

					if strings.HasPrefix(text, "Материал плафона") && product.AbajurMaterial == "None" {
						value := s.Find(".char_value span")
						valueText := cleanText(value.Text())
						product.AbajurMaterial = replaceSeparatorsWithSlash(valueText)
					}

					if strings.HasPrefix(text, "Материал арматуры") && product.ArmaturMaterial == "None" {
						value := s.Find(".char_value span")
						valueText := cleanText(value.Text())
						product.ArmaturMaterial = replaceSeparatorsWithSlash(valueText)
					}

					if strings.HasPrefix(text, "Цвет арматуры") && product.ArmaturColor == "None" {
						value := s.Find(".char_value span")
						valueText := cleanText(value.Text())
						product.ArmaturColor = replaceSeparatorsWithSlash(valueText)
					}

					if strings.HasPrefix(text, "Цвет плафона") && product.AbajurColor == "None" {
						value := s.Find(".char_value span")
						valueText := cleanText(value.Text())
						product.AbajurColor = replaceSeparatorsWithSlash(valueText)
					}
				})
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




func extractName(url string) string {
	count := 0
	name := []rune{}
	for _, let := range url {
		if let == '/' {
			count += 1
		}

		if count == 10 && let != '/'{
			name = append(name, let)
		}
	}
	filename := fmt.Sprintf("%s", string(name))
	return filename
}

func takeLinks(link string, pages int) ([]string, []string) {
	var links []string
	var names []string
	
	c := colly.NewCollector()

	
	c.OnHTML(".catalog_item_wrapp", func(e *colly.HTMLElement) {
	
    	link := e.ChildAttr(".inner_wrap .image_wrapper_block a", "href")
    	fullLink := e.Request.AbsoluteURL(link)
    	//log.Println("Ссылка на товар:", fullLink)

		imgUrl := e.ChildAttr(".inner_wrap .image_wrapper_block a img", "data-src")
		fullImgURL := e.Request.AbsoluteURL(imgUrl)
		filename := fmt.Sprintf("FAVOURITE%s", extractName(fullImgURL))
		links = append(links, fullLink)
		names = append(names, filename)
		downloadImg(fullImgURL, filename)

		

		

		//log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, filename)
	})	

	for i := 1; i <= pages; i++{
		url := "%s%d"


		page := fmt.Sprintf(url, link,  i)
		log.Println(page)
		err := c.Visit(page)
		if err != nil {
			log.Fatal(err)
		}
	}

	return links, names
}

func cleanText(raw string) string {
    // Заменяем неразрывный пробел на обычный пробел
    cleaned := strings.ReplaceAll(raw, "\u00A0", " ")
    // Убираем дефис и лишние пробелы
    cleaned = strings.ReplaceAll(cleaned, "-", "")
    // Удаляем все лишние пробелы
    cleaned = strings.TrimSpace(cleaned)
    return cleaned
}

func doSlash(text string) string {
	if strings.Contains(text, "/") {
        // Разбиваем строку по слэшу и берем первую часть
        parts := strings.Split(text, "/")
        res := parts[0]
		return res
    }
	return text
}

func maxValue(s string) int {
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
	return max
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