package usecase

import (
	"bytes"
	"image/color"
	"log"

	"github.com/OlegLaban/geo-flag/internal/config"
	"github.com/OlegLaban/geo-flag/internal/domain"
	"github.com/fogleman/gg"
)


type traydata struct {
	conf *config.Config
	geoData domain.GeoData
}


func NewTrayData(confg *config.Config, geoData domain.GeoData) *traydata {
	return &traydata{conf: confg, geoData: geoData}
}

func (td *traydata) GetIcon() []byte {
	var icon []byte
	switch td.conf.IconType {
	case config.ImgIcon:
		icon = td.geoData.Flag
	case config.CodeIcon:
		icon = GenerateIcon(td.geoData.CountryCode)
	}
	return  icon
}

func (td *traydata) GetTitle() string {
	return  td.geoData.CountryName
}

func GenerateIcon(text string) []byte {
    const size = 24
    const fontSize = 14

    dc := gg.NewContext(size, size)
    dc.SetRGBA(0, 0, 0, 1)
    dc.Clear()

    err := dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", fontSize)
    if err != nil {
        log.Println("Font error:", err)
        return nil
    }

    dc.SetColor(color.White)
    dc.DrawStringAnchored(text, size/2, size/2, 0.5, 0.5)

    var buf bytes.Buffer
    err = dc.EncodePNG(&buf)
    if err != nil {
        log.Println("PNG encode error:", err)
        return nil
    }

    return buf.Bytes()
}

