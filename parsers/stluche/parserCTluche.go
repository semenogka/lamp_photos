package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gocolly/colly"
)

type Product struct {
	Name      string `json:"name"`
	Link      string `json:"link"`
}
var pages int = 1
var products []Product

func main() {
	c := colly.NewCollector()

	
	c.OnHTML(".p-catalog__product-item", func(e *colly.HTMLElement) {
		// title := e.ChildText("div .product__content .product__title")
		// log.Println(title)
		// if strings.Contains(title, "светильники"){
			link := e.ChildAttr("div a", "href")
			fullLink := e.Request.AbsoluteURL(link)
			log.Println("Ссылка на товар:", fullLink)

			imgUrl := e.ChildAttr("div .product__image-wrapper .product__image-box  img", "data-src")
			fullImgURL := e.Request.AbsoluteURL(imgUrl)
			filename := fmt.Sprintf("CTLUCHE%s", extractName(fullImgURL))

			downloadImg(fullImgURL, filename)

			products = append(products, Product{
				Name: filename,
				Link: fullLink,
			})

			

			log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, len(products), filename, imgUrl)
		// }


    	
	})	

	for i := 1; i <= pages; i++{
		url := "https://stluce.ru/catalog/ulichnyy_svet/potolochnye_svetilniki_1/?PAGEN_1=%d"


		page := fmt.Sprintf(url, i)
		log.Println(page)
		err := c.Visit(page)
		if err != nil {
			log.Fatal(err)
		}
	}

	
	saveToJson("alldata\\data\\potol.json", products)
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