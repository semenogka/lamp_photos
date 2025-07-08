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

var artelampSites = []SiteStruct{
	{URL: "https://artelamp.ru/catalog/lyustryi/page_", Category: "alldata\\data\\lystri.json", Pages: 15},
	{URL: "https://artelamp.ru/catalog/potolochnyie-svetilniki/page_", Category: "alldata\\data\\potol.json", Pages: 4},
	{URL: "https://artelamp.ru/catalog/podvesnyie-svetilniki/page_", Category: "alldata\\data\\podves.json", Pages: 11},
	{URL: "https://artelamp.ru/catalog/torsheryi/page_", Category: "alldata\\data\\torsher.json", Pages: 3},
	{URL: "https://artelamp.ru/catalog/trekovyie-sistemyi/trekovyie-svetilniki/page_", Category: "alldata\\data\\track.json", Pages: 6},
	{URL: "https://artelamp.ru/catalog/magnitnyie-trekovyie-sistemyi/magnitnyie-trekovyie-svetilniki/page_", Category: "alldata\\data\\magn.json", Pages: 8},
	{URL: "https://artelamp.ru/catalog/podsvetki/page_", Category: "alldata\\data\\podsvet.json", Pages: 5},
	{URL: "https://artelamp.ru/catalog/nastennyie-svetilniki-i-bra/page_", Category: "alldata\\data\\bra.json", Pages: 11},
	{URL: "https://artelamp.ru/catalog/nastolnyie-lampyi-i-nochniki/page_", Category: "alldata\\data\\nastol.json", Pages: 6},
	{URL: "https://artelamp.ru/catalog/tochechnyie-svetilniki/tochechnyie-vstraivaemyie-svetilniki/page_", Category: "alldata\\data\\tochvstr.json", Pages: 8},
	{URL: "https://artelamp.ru/catalog/tochechnyie-svetilniki/tochechnyie-nakladnyie-svetilniki/page_", Category: "alldata\\data\\tochnakl.json", Pages: 4},
	{URL: "https://artelamp.ru/catalog/tochechnyie-svetilniki/tochechnyie-podvesnyie-svetilniki/page_", Category: "alldata\\data\\tochpodv.json", Pages: 2},
	{URL: "https://artelamp.ru/catalog/ulichnoe-osveshhenie/gruntovyie-svetilniki/page_", Category: "alldata\\data\\groont.json", Pages: 1},
	{URL: "https://artelamp.ru/catalog/ulichnoe-osveshhenie/landshaftnyie-svetilniki/page_", Category: "alldata\\data\\land.json", Pages: 2},
	{URL: "https://artelamp.ru/catalog/ulichnoe-osveshhenie/parkovyie-svetilniki/page_", Category: "alldata\\data\\parkovie.json", Pages: 1},
	{URL: "https://artelamp.ru/catalog/ulichnoe-osveshhenie/trotuarnyie-svetilniki/page_", Category: "alldata\\data\\trotuarnie.json", Pages: 1},
	{URL: "https://artelamp.ru/catalog/ulichnoe-osveshhenie/ulichnyie-nastolnyie-svetilniki/page_", Category: "alldata\\data\\nastol.json", Pages: 4},
	{URL: "https://artelamp.ru/catalog/ulichnoe-osveshhenie/fasadnyie-svetilniki/page_", Category: "alldata\\data\\fasad.json", Pages: 4},
}



func main() {
	//links := []string{}
	for _, site := range artelampSites {
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
		
			c.OnHTML(".card_characters_list_content", func(e *colly.HTMLElement) {
				product.Link = link
				product.Name = names[i]
				
				e.DOM.Find(".card_characters_list_content_left").Each(func(_ int, p *goquery.Selection) {
					uls := e.DOM.Find("ul")
					e.DOM.Find("div").Each(func(_ int, d *goquery.Selection) {
					text := strings.TrimSpace(d.Text())
					gabariti := uls.Eq(2)
					if text == "Габариты светильника" {
						gabariti = d.Next()
					}
					
					
					gabariti.Find("li").Each(func(i int, s *goquery.Selection) {
						div := s.Find("div.name")
						name :=  strings.TrimSpace(div.Text())
						
						if name == "Высота" {
							div = s.Find("div.value")
							value := maxValue(div.Text())
							value = value * 10
							
							product.Height = strconv.Itoa(value)
						}
						if name == "Ширина" {
							div = s.Find("div.value")
							value := maxValue(div.Text())
							value = value * 10
							
							product.Width = strconv.Itoa(value)
						}
						if name == "Длина" {
							div = s.Find("div.value")
							value := maxValue(div.Text())
							value = value * 10
							
							product.Length = strconv.Itoa(value)
						}
						if name == "Диаметр" {
							div = s.Find("div.value")
							value := maxValue(div.Text())
							value = value * 10
							
							product.Diameter = strconv.Itoa(value)
						}
						})
					})
					
						
						
					
				})
				e.DOM.Find(".card_characters_list_content_right").Each(func(_ int, p *goquery.Selection) {
					uls := p.Find("ul")

					colors := uls.Eq(0)

					colors.Find("li").Each(func(i int, s *goquery.Selection) {
						div := s.Find("div.name")
						name :=  strings.TrimSpace(div.Text())
						if name == "Материал плафона" || name == "Материалы плафона" {
							div = s.Find("div.value")
							value := cleanValue(div.Text())
							
							product.AbajurMaterial = replaceSeparatorsWithSlash(value)
						}
						if name == "Цвет плафона" || name == "Цвета плафона" {
							div = s.Find("div.value")
							value := cleanValue(div.Text())
							
							product.AbajurColor = replaceSeparatorsWithSlash(value)
						}
						if name == "Материал арматуры" || name == "Материалы арматуры"{
							div = s.Find("div.value")
							value := cleanValue(div.Text())
							
							product.ArmaturMaterial = replaceSeparatorsWithSlash(value)
						}
						if name == "Цвет светильника"|| name == "Цвета светильника"{
							div := s.Find("div.value")
							value := cleanValue(div.Text())
							product.ArmaturColor = replaceSeparatorsWithSlash(value)
						}
					})
				})
					
				
			})
			c.OnScraped(func(r *colly.Response) {
				log.Println(product)
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

func takeLinks(link string, pages int) ([]string, []string) {
	var links []string
	var names []string
	
	c := colly.NewCollector()
	c.OnHTML(".listing_content_catalog_units", func(e *colly.HTMLElement) {
		e.ForEach(".unit", func(i int, h *colly.HTMLElement) {
			link := h.ChildAttr(".element .imgwr a", "href")
			fullLink := h.Request.AbsoluteURL(link)
			
			imgURL := h.ChildAttr(".element .imgwr a img", "data-src")
			fullImgURL := h.Request.AbsoluteURL(imgURL)
			
			filename := fmt.Sprintf("ARTELAMP%s", extractName(fullImgURL))

			downloadImg(fullImgURL, filename)

			

			//log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, " ", filename)

			names = append(names, filename)
			links = append(links, fullLink)
		})
				
	})	
	for i := 1; i <= pages; i++ {
		
		page := fmt.Sprintf("%s%d", link, i)
		err := c.Visit(page)

		if err != nil {
			log.Fatal(err)
		}
	}

	return links, names
}

func extractNumber(s string) string {
	re := regexp.MustCompile(`\d+`)
	return strings.Join(re.FindAllString(s, -1), "")
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

func cleanValue(value string) string {
    value = strings.TrimSpace(value)
    parts := strings.Split(value, ",")
    return strings.TrimSpace(parts[0])
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

func replaceSeparatorsWithSlash(s string) string {
	reSep := regexp.MustCompile(`\s*(,|;)\s*`)
	s = reSep.ReplaceAllString(s, "/")

	reSlash := regexp.MustCompile(`\s*/\s*`)
	return reSlash.ReplaceAllString(s, "/")
}