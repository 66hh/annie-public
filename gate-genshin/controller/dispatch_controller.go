package controller

import (
	"bytes"
	"encoding/base64"
	"flswld.com/common/utils/endec"
	"flswld.com/logger"
	"gate-genshin/entity/api"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
)

func (c *Controller) query_security_file(context *gin.Context) {
	file, err := ioutil.ReadFile("static/security_file")
	if err != nil {
		logger.LOG.Error("open security_file error")
		return
	}
	context.Header("Content-type", "text/html; charset=UTF-8")
	_, _ = context.Writer.WriteString(string(file))
}

func (c *Controller) query_region_list(context *gin.Context) {
	context.Header("Content-type", "text/html; charset=UTF-8")
	_, _ = context.Writer.WriteString(c.regionListBase64)
}

func (c *Controller) query_cur_region(context *gin.Context) {
	versionName := context.Query("version")
	response := "CAESGE5vdCBGb3VuZCB2ZXJzaW9uIGNvbmZpZw=="
	if len(context.Request.URL.RawQuery) > 0 {
		response = c.regionCurrBase64
	}
	if strings.Contains(versionName, "2.7.5") || strings.Contains(versionName, "2.8.") {
		logger.LOG.Debug("do genshin 2.8 rsa logic")
		if context.Query("dispatchSeed") == "" {
			rsp := &api.QueryCurRegionRspJson{
				Content: response,
				Sign:    "TW9yZSBsb3ZlIGZvciBVQSBQYXRjaCBwbGF5ZXJz",
			}
			context.JSON(http.StatusOK, rsp)
			return
		}
		keyId := context.Query("key_id")
		var encPubPrivKey []byte = nil
		if keyId == "3" {
			// 国际服
			encPubPrivKey = c.encRsaKey
		} else {
			// 国服
			logger.LOG.Error("current region enc key not exist")
			return
		}
		regionInfo, err := base64.StdEncoding.DecodeString(response)
		if err != nil {
			logger.LOG.Error("decode region info error: %v", err)
			return
		}
		chunkSize := 256 - 11
		regionInfoLength := len(regionInfo)
		numChunks := int(math.Ceil(float64(regionInfoLength) / float64(chunkSize)))
		encryptedRegionInfo := make([]byte, 0)
		for i := 0; i < numChunks; i++ {
			from := i * chunkSize
			to := int(math.Min(float64((i+1)*chunkSize), float64(regionInfoLength)))
			chunk := regionInfo[from:to]
			pubKey, err := endec.RsaParsePubKeyByPrivKey(encPubPrivKey)
			if err != nil {
				logger.LOG.Error("parse rsa pub key error: %v", err)
				return
			}
			privKey, err := endec.RsaParsePrivKey(encPubPrivKey)
			if err != nil {
				logger.LOG.Error("parse rsa priv key error: %v", err)
				return
			}
			encrypt, err := endec.RsaEncrypt(chunk, pubKey)
			if err != nil {
				logger.LOG.Error("rsa enc error: %v", err)
				return
			}
			decrypt, err := endec.RsaDecrypt(encrypt, privKey)
			if err != nil {
				logger.LOG.Error("rsa dec error: %v", err)
				return
			}
			if bytes.Compare(decrypt, chunk) != 0 {
				logger.LOG.Error("rsa dec test fail")
				return
			}
			encryptedRegionInfo = append(encryptedRegionInfo, encrypt...)
		}
		signPrivkey, err := endec.RsaParsePrivKey(c.signRsaKey)
		if err != nil {
			logger.LOG.Error("parse rsa priv key error: %v", err)
			return
		}
		signData, err := endec.RsaSign(regionInfo, signPrivkey)
		if err != nil {
			logger.LOG.Error("rsa sign error: %v", err)
			return
		}
		ok, err := endec.RsaVerify(regionInfo, signData, &signPrivkey.PublicKey)
		if err != nil {
			logger.LOG.Error("rsa verify error: %v", err)
			return
		}
		if !ok {
			logger.LOG.Error("rsa verify test fail")
			return
		}
		rsp := &api.QueryCurRegionRspJson{
			Content: base64.StdEncoding.EncodeToString(encryptedRegionInfo),
			Sign:    base64.StdEncoding.EncodeToString(signData),
		}
		context.JSON(http.StatusOK, rsp)
		return
	} else {
		context.Header("Content-type", "text/html; charset=UTF-8")
		_, _ = context.Writer.WriteString(response)
	}
}
