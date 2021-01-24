package ocr

import (
	"strings"

	"github.com/liudanking/gotext/cfg"

	"github.com/chenqinghe/baidu-ai-go-sdk/vision"
	"github.com/chenqinghe/baidu-ai-go-sdk/vision/ocr"
	log "github.com/liudanking/goutil/logutil"
)

type BaiduOCRRsp struct {
	LogID          int64 `json:"log_id"`
	Direction      int   `json:"direction"`
	WordsResultNum int   `json:"words_result_num"`
	WordsResult    []struct {
		Words string `json:"words"`
	} `json:"words_result"`
}

var _baiduOcrClient *ocr.OCRClient

func getBaiduOCRClient() *ocr.OCRClient {
	if _baiduOcrClient == nil {
		_baiduOcrClient = ocr.NewOCRClient(cfg.Get().BaiduAIConf.AppKey, cfg.Get().BaiduAIConf.AppSecret)
	}
	return _baiduOcrClient
}

func GetOCRTextWithBaiduAI(fn string) (string, error) {
	img, err := vision.FromFile(fn)
	if err != nil {
		log.Error("get image file from [%s] failed:%v", fn, err)
		return "", err
	}
	client := getBaiduOCRClient()
	rsp, err := client.AccurateRecognizeBasic(
		img,
		ocr.DetectDirection(),
		ocr.DetectLanguage(),
		ocr.LanguageType("CHN_ENG"),
	)
	if err != nil {
		log.Error("GetOCRTextWithBaiduAI [fn:%s] failed:%v", fn, err)
		return "", nil
	}

	ocrRsp := &BaiduOCRRsp{}
	rsp.ToJSON(ocrRsp)

	sb := &strings.Builder{}
	for _, item := range ocrRsp.WordsResult {
		sb.WriteString(item.Words)
		sb.WriteString("\n")
	}

	return sb.String(), nil

}
