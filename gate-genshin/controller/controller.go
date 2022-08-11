package controller

import (
	"encoding/base64"
	"flswld.com/common/config"
	"flswld.com/light"
	"flswld.com/logger"
	"gate-genshin/dao"
	"gate-genshin/region"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	pb "google.golang.org/protobuf/proto"
	"net/http"
	"strconv"
)

type Controller struct {
	dao              *dao.Dao
	rpcUserConsumer  *light.Consumer
	regionListBase64 string
	regionCurrBase64 string
	signRsaKey       []byte
	encRsaKey        []byte
	pwdRsaKey        []byte
}

func NewController(dao *dao.Dao, rpcUserConsumer *light.Consumer) (r *Controller) {
	r = new(Controller)
	r.dao = dao
	r.rpcUserConsumer = rpcUserConsumer
	r.regionListBase64 = ""
	r.regionCurrBase64 = ""
	regionCurr, regionList := region.InitRegion(config.CONF.Genshin.KcpAddr, config.CONF.Genshin.KcpPort)
	r.signRsaKey, r.encRsaKey, r.pwdRsaKey = region.LoadRsaKey()
	regionCurrModify, err := pb.Marshal(regionCurr)
	if err != nil {
		logger.LOG.Error("Marshal QueryCurrRegionHttpRsp error")
		return nil
	}
	r.regionCurrBase64 = base64.StdEncoding.EncodeToString(regionCurrModify)
	regionListModify, err := pb.Marshal(regionList)
	if err != nil {
		logger.LOG.Error("Marshal QueryRegionListHttpRsp error")
		return nil
	}
	r.regionListBase64 = base64.StdEncoding.EncodeToString(regionListModify)
	r.runEngine()
	return r
}

func (c *Controller) runEngine() {
	if config.CONF.Logger.Level == "DEBUG" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	go func() {
		engine := c.registerRouter()
		engine.Use(c.tlsHandler())
		port := config.CONF.Genshin.HttpsPort
		addr := ":" + strconv.FormatInt(int64(port), 10)
		err := engine.RunTLS(addr, "static/server.crt", "static/server.key")
		if err != nil {
			logger.LOG.Error("gin run error: %v", err)
		}
	}()
	go func() {
		engine := c.registerRouter()
		port := config.CONF.Genshin.HttpPortA
		addr := ":" + strconv.FormatInt(int64(port), 10)
		err := engine.Run(addr)
		if err != nil {
			logger.LOG.Error("gin run error: %v", err)
		}
	}()
	go func() {
		engine := c.registerRouter()
		port := config.CONF.Genshin.HttpPortB
		addr := ":" + strconv.FormatInt(int64(port), 10)
		err := engine.Run(addr)
		if err != nil {
			logger.LOG.Error("gin run error: %v", err)
		}
	}()
}

func (c *Controller) registerRouter() *gin.Engine {
	engine := gin.Default()
	{
		// 404
		engine.NoRoute(func(context *gin.Context) {
			logger.LOG.Info("no route find, fallback to fuck mhy, url: %v", context.Request.RequestURI)
			context.Header("Content-type", "text/html; charset=UTF-8")
			context.Status(http.StatusNotFound)
			_, _ = context.Writer.WriteString("FUCK MHY")
		})
	}
	{
		// 调度
		// dispatchosglobal.yuanshen.com
		engine.GET("/query_security_file", c.query_security_file)
		engine.GET("/query_region_list", c.query_region_list)
		// osusadispatch.yuanshen.com
		engine.GET("/query_cur_region", c.query_cur_region)
	}
	{
		// 登录
		// hk4e-sdk-os.hoyoverse.com
		// 账号登录
		engine.POST("/hk4e_global/mdk/shield/api/login", c.apiLogin)
		// token登录
		engine.POST("/hk4e_global/mdk/shield/api/verify", c.apiVerify)
		// 获取combo token
		engine.POST("/hk4e_global/combo/granter/login/v2/login", c.v2Login)
	}
	{
		// BLK文件补丁下载
		// autopatchhk.yuanshen.com
		engine.HEAD("/client_design_data/2.6_live/output_6988297_84eeb1c18b/client_silence/General/AssetBundles/data_versions", c.headDataVersions)
		engine.GET("/client_design_data/2.6_live/output_6988297_84eeb1c18b/client_silence/General/AssetBundles/data_versions", c.getDataVersions)
		engine.HEAD("/client_design_data/2.6_live/output_6988297_84eeb1c18b/client_silence/General/AssetBundles/blocks/00/29342328.blk", c.headBlk)
		engine.GET("/client_design_data/2.6_live/output_6988297_84eeb1c18b/client_silence/General/AssetBundles/blocks/00/29342328.blk", c.getBlk)
	}
	{
		// 日志
		engine.POST("/sdk/dataUpload", c.sdkDataUpload)
		engine.GET("/perf/config/verify", c.perfConfigVerify)
		engine.POST("/perf/dataUpload", c.perfDataUpload)
		engine.POST("/log", c.log8888)
		engine.POST("/crash/dataUpload", c.crashDataUpload)
	}
	{
		// 返回固定数据
		// Windows
		engine.GET("/hk4e_global/mdk/agreement/api/getAgreementInfos", c.getAgreementInfos)
		engine.POST("/hk4e_global/combo/granter/api/compareProtocolVersion", c.postCompareProtocolVersion)
		engine.POST("/account/risky/api/check", c.check)
		engine.GET("/combo/box/api/config/sdk/combo", c.combo)
		engine.GET("/hk4e_global/combo/granter/api/getConfig", c.getConfig)
		engine.GET("/hk4e_global/mdk/shield/api/loadConfig", c.loadConfig)
		engine.POST("/data_abtest_api/config/experiment/list", c.list)
		engine.GET("/admin/mi18n/plat_oversea/m2020030410/m2020030410-version.json", c.version10Json)
		engine.GET("/admin/mi18n/plat_oversea/m2020030410/m2020030410-zh-cn.json", c.zhCN10Json)
		engine.GET("/geetestV2.html", c.geetestV2)
		// Android
		engine.POST("/common/h5log/log/batch", c.batch)
		engine.GET("/hk4e_global/combo/granter/api/getFont", c.getFont)
		engine.GET("/admin/mi18n/plat_oversea/m202003049/m202003049-version.json", c.version9Json)
		engine.GET("/admin/mi18n/plat_oversea/m202003049/m202003049-zh-cn.json", c.zhCN9Json)
		engine.GET("/hk4e_global/combo/granter/api/compareProtocolVersion", c.getCompareProtocolVersion)
		// Android geetest
		engine.GET("/gettype.php", c.gettype)
		engine.GET("/get.php", c.get)
		engine.POST("/ajax.php", c.ajax)
		engine.GET("/ajax.php", c.ajax)
		engine.GET("/static/appweb/app3-index.html", c.app3Index)
		engine.GET("/static/js/slide.7.8.6.js", c.slideJs)
		engine.GET("/favicon.ico", c.faviconIco)
		engine.GET("/static/js/gct.e7810b5b525994e2fb1f89135f8df14a.js", c.js)
		engine.GET("/static/ant/style_https.1.2.6.css", c.css)
		engine.GET("/pictures/gt/a330cf996/a330cf996.webp", c.webp)
		engine.GET("/pictures/gt/a330cf996/bg/86f9db021.webp", c.bgWebp)
		engine.GET("/pictures/gt/a330cf996/slice/86f9db021.png", c.slicePng)
		engine.GET("/static/ant/sprite2x.1.2.6.png", c.sprite2xPng)
	}
	return engine
}

func (c *Controller) tlsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		s := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":8888",
		})
		err := s.Process(context.Writer, context.Request)
		if err != nil {
			logger.LOG.Error("gin listen tls error")
			return
		}
		context.Next()
	}
}
