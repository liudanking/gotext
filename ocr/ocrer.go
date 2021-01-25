package ocr

import (
	"errors"

	"github.com/liudanking/gotext/cfg"
	log "github.com/liudanking/goutil/logutil"
)

type OCRer interface {
	GetOCRText(fn string) (string, error)
}

var _ocrer OCRer

func InitOCRer(config *cfg.Config) error {
	switch config.OCRPlatform {
	case "baidu":
		_ocrer = newBaiduOCR(config.BaiduAIConf.AppKey, config.BaiduAIConf.AppSecret)
	default:
		log.Error("[ocr_platform:%s] not supported", config.OCRPlatform)
		return errors.New("ocr_platform invalid")
	}

	return nil
}

func GetOCRText(fn string) (string, error) {
	return _ocrer.GetOCRText(fn)
}
