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
	"strings"

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
	{URL: "https://artelamp.ru/catalog/tochechnyie-svetilniki/tochechnyie-podvesnyie-svetilniki/page_", Category: "alldata\\data\\tochpodves.json", Pages: 2},
	{URL: "https://artelamp.ru/catalog/ulichnoe-osveshhenie/gruntovyie-svetilniki/page_", Category: "alldata\\data\\groont.json", Pages: 1},
	{URL: "https://artelamp.ru/catalog/ulichnoe-osveshhenie/landshaftnyie-svetilniki/page_", Category: "alldata\\data\\land.json", Pages: 2},
	{URL: "https://artelamp.ru/catalog/ulichnoe-osveshhenie/parkovyie-svetilniki/page_", Category: "alldata\\data\\parkovie.json", Pages: 1},
	{URL: "https://artelamp.ru/catalog/ulichnoe-osveshhenie/trotuarnyie-svetilniki/page_", Category: "alldata\\data\\trotuarnie.json", Pages: 1},
	{URL: "https://artelamp.ru/catalog/ulichnoe-osveshhenie/ulichnyie-nastolnyie-svetilniki/page_", Category: "alldata\\data\\nastol.json", Pages: 4},
	{URL: "https://artelamp.ru/catalog/ulichnoe-osveshhenie/fasadnyie-svetilniki/page_", Category: "alldata\\data\\fasad.json", Pages: 4},
}

var crystalluxSites = []SiteStruct{
	{URL: "https://crystallux.ru/produktsiya/dekorativnyy-svet/lyustry/?PAGEN_1=", Category: "alldata\\data\\lystri.json", Pages: 26},
	{URL: "https://crystallux.ru/produktsiya/dekorativnyy-svet/torshery/?PAGEN_1=", Category: "alldata\\data\\torsher.json", Pages: 2},
	{URL: "https://crystallux.ru/produktsiya/dekorativnyy-svet/bra/?PAGEN_1=", Category: "alldata\\data\\bra.json", Pages: 22},
	{URL: "https://crystallux.ru/produktsiya/dekorativnyy-svet/nastolnye-lampy/?PAGEN_1=", Category: "alldata\\data\\nastol.json", Pages: 3},
	{URL: "https://crystallux.ru/produktsiya/dekorativnyy-svet/svetilniki/?PAGEN_1=", Category: "alldata\\data\\podves.json", Pages: 24},
	{URL: "https://crystallux.ru/produktsiya/dekorativnyy-svet/dekorativnyy-svet-potolochnye-svetilniki/?PAGEN_1=", Category: "alldata\\data\\potol.json", Pages: 5},
	{URL: "https://crystallux.ru/produktsiya/crystal-lux-technical-clt/clt-nastennye-svetilniki-bra/?PAGEN_1=", Category: "alldata\\data\\bra.json", Pages: 12},
	{URL: "https://crystallux.ru/produktsiya/crystal-lux-technical-clt/vstraivaemye/?PAGEN_1=", Category: "alldata\\data\\tochvstr.json", Pages: 13},
	{URL: "https://crystallux.ru/produktsiya/crystal-lux-technical-clt/trekovye-sistemy/svetilniki2/?PAGEN_1=", Category: "alldata\\data\\magn.json", Pages: 6},
	{URL: "https://crystallux.ru/produktsiya/crystal-lux-technical-clt/clt-potolochnye-svetilniki/?PAGEN_1=", Category: "alldata\\data\\tochpotol.json", Pages: 6},
	{URL: "https://crystallux.ru/produktsiya/crystal-lux-technical-clt/odnofaznaya-trekovaya-sistema/odnofaznaya-svetilniki/?PAGEN_1=", Category: "alldata\\data\\track.json", Pages: 4},
	{URL: "https://crystallux.ru/produktsiya/crystal-lux-technical-clt/technical-clt-podvesnye-svetilniki/?PAGEN_1=", Category: "alldata\\data\\tochpodv.json", Pages: 6},
}

var divinareSites = []SiteStruct{
	{URL: "https://divinare.ru/catalog/lyustryi/page_", Category: "alldata\\data\\lystri.json", Pages: 6},
	{URL: "https://divinare.ru/catalog/torsheryi/page_", Category: "alldata\\data\\torsher.json", Pages: 1},
	{URL: "https://divinare.ru/catalog/trekovyie-sistemyi/page_", Category: "alldata\\data\\track.json", Pages: 3},
	{URL: "https://divinare.ru/catalog/podvesnyie-svetilniki/page_", Category: "alldata\\data\\podves.json", Pages: 6},
	{URL: "https://divinare.ru/catalog/nastennyie-svetilniki-i-bra/page_", Category: "alldata\\data\\bra.json", Pages: 4},
	{URL: "https://divinare.ru/catalog/nastolnyie-lampyi-i-nochniki/page_", Category: "alldata\\data\\nastol.json", Pages: 1},
	{URL: "https://divinare.ru/catalog/tochechnyie-svetilniki/tochechnyie-vstraivaemyie-svetilniki/page_", Category: "alldata\\data\\tochvstr.json", Pages: 1},
	{URL: "https://divinare.ru/catalog/tochechnyie-svetilniki/tochechnyie-nakladnyie-svetilniki/page_", Category: "alldata\\data\\tochnakl.json", Pages: 1},
	{URL: "https://divinare.ru/catalog/tochechnyie-svetilniki/tochechnyie-podvesnyie-svetilniki/page_", Category: "alldata\\data\\tochpodv.json", Pages: 1},
	{URL: "https://divinare.ru/catalog/podsvetki/page_", Category: "alldata\\data\\podsvet.json", Pages: 1},
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

var freayaSites = []SiteStruct {
	{URL: "https://freya-light.com/products/dekorativnyy_svet/lyustra/?SHOWALL_1=1", Category: "alldata\\data\\lystri.json", Pages: 0},
	{URL: "https://freya-light.com/products/dekorativnyy_svet/podvesnoy_svetilnik/?SHOWALL_1=1", Category: "alldata\\data\\podves.json", Pages: 0},
	{URL: "https://freya-light.com/products/dekorativnyy_svet/potolochnyy_svetilnik/?SHOWALL_1=1", Category: "alldata\\data\\potol.json", Pages: 0},
	{URL: "https://freya-light.com/products/dekorativnyy_svet/bra/?SHOWALL_1=1", Category: "alldata\\data\\bra.json", Pages: 0},
	{URL: "https://freya-light.com/products/dekorativnyy_svet/nastolnaya_lampa/?SHOWALL_1=1", Category: "alldata\\data\\nastol.json", Pages: 0},
	{URL: "https://freya-light.com/products/dekorativnyy_svet/torsher/?SHOWALL_1=1", Category: "alldata\\data\\torsher.json", Pages: 0},
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
	artelampParser()
	crystallluxParser()
	divinareParser()
	favouriteParser()
	freyaParser()
	isonexParser()
	lightstarParser()
	loftitParser()
	maytoniParser()
	stluceParser()
}

func artelampParser() {
	for _, site := range artelampSites{
		var products []Product
		c := colly.NewCollector()
		c.OnHTML(".listing_content_catalog_units", func(e *colly.HTMLElement) {
			e.ForEach(".unit", func(i int, h *colly.HTMLElement) {
				link := h.ChildAttr(".element .imgwr a", "href")
				fullLink := h.Request.AbsoluteURL(link)

				imgURL := h.ChildAttr(".element .imgwr a img", "data-src")
				fullImgURL := h.Request.AbsoluteURL(imgURL)
				
				filename := fmt.Sprintf("ARTELAMP%s", extractName(fullImgURL))

				downloadImg(fullImgURL, filename)

				products = append(products, Product{
					Name: filename,
					Link: fullLink,
				})

				log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, len(products), " ", filename)
			})
					
		})	
		for i := 1; i <= site.Pages; i++ {
			
			page := fmt.Sprintf("%s%d", site.URL, i)
			err := c.Visit(page)

			if err != nil {
				log.Fatal(err)
			}
		}
		
		
		saveToJson(site.Category, products)
	}
	
}

func crystallluxParser() {
	for _, site := range crystalluxSites {
		var products []Product
		c := colly.NewCollector()

	
		c.OnHTML(".product-item", func(e *colly.HTMLElement) {
			link := e.ChildAttr("a", "href")
			fullLink := e.Request.AbsoluteURL(link)
			log.Println("Ссылка на товар:", fullLink)

			imgStyle := e.ChildAttr("a span .product-item-image-slide", "style")
			imgUrl := extractImgUrl(imgStyle)
			fullImgURL := e.Request.AbsoluteURL(imgUrl)
			filename := fmt.Sprintf("CRYSTALLUX%s", extractName(fullImgURL))

			//downloadImg(fullImgURL, filename)

			products = append(products, Product{
				Name: filename,
				Link: fullLink,
			})

			

			log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, len(products), filename, imgUrl)
		})	

		for i := 1; i <= site.Pages; i++{
			url := "%s%d"


			page := fmt.Sprintf(url, site.URL, i)
			log.Println(page)
			err := c.Visit(page)
			if err != nil {
				log.Fatal(err)
			}
		}

		
		//saveToJson(site.Category, products)
	}
	
}

func divinareParser() {
	for _, site := range divinareSites {
		var products []Product
		c := colly.NewCollector()
		c.OnHTML(".listing_content_catalog_units", func(e *colly.HTMLElement) {
			e.ForEach(".unit", func(i int, h *colly.HTMLElement) {
				link := h.ChildAttr(".element .imgwr a", "href")
				fullLink := h.Request.AbsoluteURL(link)

				imgURL := h.ChildAttr(".element .imgwr a img", "data-src")
				fullImgURL := h.Request.AbsoluteURL(imgURL)
				
				filename := fmt.Sprintf("DIVINARE%s", extractName(fullImgURL))

				downloadImg(fullImgURL, filename)

				products = append(products, Product{
					Name: filename,
					Link: fullLink,
				})

				log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, len(products), " ", filename)
				


			})
					
		})	
		for i := 1; i <= site.Pages; i++ {
			
			page := fmt.Sprintf("%s%d", site.URL,  i)
			err := c.Visit(page)

			if err != nil {
				log.Fatal(err)
			}
		}
		
		
		saveToJson(site.Category, products)
	}
}

func favouriteParser() {
	for _, site := range favouriteSites {
		var products []Product
		c := colly.NewCollector()

	
		c.OnHTML(".catalog_item_wrapp", func(e *colly.HTMLElement) {
		
			link := e.ChildAttr(".inner_wrap .image_wrapper_block a", "href")
			fullLink := e.Request.AbsoluteURL(link)
			log.Println("Ссылка на товар:", fullLink)

			imgUrl := e.ChildAttr(".inner_wrap .image_wrapper_block a img", "data-src")
			fullImgURL := e.Request.AbsoluteURL(imgUrl)
			filename := fmt.Sprintf("FAVOURITE%s", extractName(fullImgURL))

			downloadImg(fullImgURL, filename)

			products = append(products, Product{
				Name: filename,
				Link: fullLink,
			})

			

			log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, len(products), filename)
		})	

		for i := 1; i <= site.Pages; i++{
			url := "%s%d"


			page := fmt.Sprintf(url, site.URL, i)
			log.Println(page)
			err := c.Visit(page)
			if err != nil {
				log.Fatal(err)
			}
		}

		
		saveToJson(site.Category, products)
	}
}

func freyaParser() {
	for _, site := range freayaSites {
		var products []Product
		c := colly.NewCollector()
		c.OnHTML(".catalog-grid__item ", func(e *colly.HTMLElement) {
			link := e.ChildAttr("a", "href")
			fullLink := e.Request.AbsoluteURL(link)
			log.Println("Ссылка на товар:", fullLink)

			imgURL := e.ChildAttr("a .catalog-card__picture picture source", "srcset")
			fullImgURL := e.Request.AbsoluteURL(imgURL)
			filename := fmt.Sprintf("FREYA%s", path.Base(imgURL))

			downloadImg(fullImgURL, filename)

			products = append(products, Product{
				Name: filename,
				Link: fullLink,
			})

			

			log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, len(products), filename)
		})	

		err := c.Visit(site.URL)

		if err != nil {
			log.Fatal(err)
		}
	
	
		saveToJson(site.Category, products)
	}
}

func isonexParser() {
	for _, site := range isonexSItes {
		var products []Product
		c := colly.NewCollector()
	
		c.OnHTML(".product-card-inner", func(h *colly.HTMLElement) {
				url := h.ChildAttr(".flex-1 .px-1 a", "href")
				fullUrl := h.Request.AbsoluteURL(url)
				imgUrl := h.ChildAttr(".flex-1 .product-card-image .position-relative .list-card-slider img", "src")
				fullImgURL := h.Request.AbsoluteURL(imgUrl)
				filename := fmt.Sprintf("ISONEX%s", path.Base(fullImgURL))

				products = append(products, Product{
					Name: filename,
					Link: url,
				})

				downloadImg(fullImgURL, filename)
				log.Println(fullUrl, " ", fullImgURL, " ", filename, " ", len(products))
		})	
		for i := 1; i <= site.Pages; i++{
			url := site.URL
			temp := strings.ReplaceAll(url, "%", "%%")
			temp = strings.Replace(temp, "%%d", "%d", 1)
			page := fmt.Sprintf(temp, i)
			log.Println(page)
			c.Visit(page)
		}
		
		
		saveToJson(site.Category, products)
	}
}

func lightstarParser() {
	for _, site := range lightstarSites {
		var products []Product
		c := colly.NewCollector()

	
		c.OnHTML(".grid__item", func(e *colly.HTMLElement) {
			link := e.ChildAttr("header .card__title a ", "href")
			fullLink := e.Request.AbsoluteURL(link)
			log.Println("Ссылка на товар:", fullLink)

			imgURL := e.ChildAttr("header .card__image-wrapper img", "src")
			fullImgURL := e.Request.AbsoluteURL(imgURL)
			filename := path.Base(imgURL)

			downloadImg(fullImgURL, filename)

			products = append(products, Product{
				Name: filename,
				Link: fullLink,
			})

			

			log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, filename, len(products))
		})	

		err := c.Visit(site.URL)
		if err != nil {
			log.Fatal(err)
		}
		
		saveToJson(site.Category, products)
	}
}

func loftitParser() {
	for _, site := range loftitSites {
		var products []Product
		c := colly.NewCollector()

		c.OnHTML(".catalog-card__item ", func(e *colly.HTMLElement) {
			link := e.ChildAttr("a", "href")
			fullLink := e.Request.AbsoluteURL(link)
			log.Println("Ссылка на товар:", fullLink)

			imgURL := e.ChildAttr("a img", "src")
			fullImgURL := e.Request.AbsoluteURL(imgURL)
			filename := fmt.Sprintf("LOFTIT%s", path.Base(imgURL))

			downloadImg(fullImgURL, filename)

			products = append(products, Product{
				Name: filename,
				Link: fullLink,
			})

			

			log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, len(products), filename)
		})	

		for i := 1; i <= site.Pages; i++{
			url := site.URL
			temp := strings.ReplaceAll(url, "%", "%%")
			temp = strings.Replace(temp, "%%d", "%d", 1)
			page := fmt.Sprintf(temp, i)
			log.Println(page)
			c.Visit(page)
		}

		
		saveToJson(site.Category, products)
	}
}

func maytoniParser() {
	for _, site := range maytoniSites {
		var products []Product
		c := colly.NewCollector()

	
		c.OnHTML(".catalog-card", func(e *colly.HTMLElement) {
			link := e.ChildAttr(".catalog-card__link", "href")
			fullLink := e.Request.AbsoluteURL(link)
			log.Println("Ссылка на товар:", fullLink)

			imgURL := e.ChildAttr(".catalog-card__img .catalog-card__img-pages picture img", "src")
			fullImgURL := e.Request.AbsoluteURL(imgURL)
			filename := path.Base(imgURL)

			downloadImg(fullImgURL, filename)

			products = append(products, Product{
				Name: filename,
				Link: fullLink,
			})

			

			log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, filename)
		})	

		err := c.Visit(site.URL)

		if err != nil {
			log.Fatal(err)
		}
		
		saveToJson(site.Category, products)
	}
}

func stluceParser() {
	for _, site := range stluceSites {
		var products []Product

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

		for i := 1; i <= site.Pages; i++{
			url := site.URL


			page := fmt.Sprintf(url, i)
			log.Println(page)
			err := c.Visit(page)
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

func extractNameIsonex(url string) string {
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