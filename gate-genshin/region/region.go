package region

import (
	"encoding/base64"
	"flswld.com/common/utils/endec"
	"flswld.com/gate-genshin-api/api/proto"
	"flswld.com/logger"
	pb "google.golang.org/protobuf/proto"
	"io/ioutil"
)

func LoadRsaKey(log *logger.Logger) (signRsaKey []byte, encRsaKey []byte, pwdRsaKey []byte) {
	var err error = nil
	signRsaKey, err = ioutil.ReadFile("static/curr_region_sign_key.pem")
	if err != nil {
		log.Error("open curr_region_sign_key.pem error: %v", err)
		return nil, nil, nil
	}
	encRsaKey, err = ioutil.ReadFile("static/curr_region_enc_key.pem")
	if err != nil {
		log.Error("open curr_region_enc_key.pem error: %v", err)
		return nil, nil, nil
	}
	pwdRsaKey, err = ioutil.ReadFile("static/account_password_key.pem")
	if err != nil {
		log.Error("open account_password_key.pem error: %v", err)
		return nil, nil, nil
	}
	return signRsaKey, encRsaKey, pwdRsaKey
}

func InitRegion(log *logger.Logger, kcpAddr string, kcpPort int) (*proto.QueryCurrRegionHttpRsp, *proto.QueryRegionListHttpRsp) {
	dispatchKey, err := ioutil.ReadFile("static/dispatchKey.bin")
	if err != nil {
		log.Error("open dispatchKey.bin error: %v", err)
		return nil, nil
	}
	dispatchSeed, err := ioutil.ReadFile("static/dispatchSeed.bin")
	if err != nil {
		log.Error("open dispatchSeed.bin error: %v", err)
		return nil, nil
	}
	// RegionCurr
	regionCurr := new(proto.QueryCurrRegionHttpRsp)
	regionCurr.RegionInfo = &proto.RegionInfo{
		GateserverIp:   kcpAddr,
		GateserverPort: uint32(kcpPort),
		SecretKey:      dispatchSeed,
	}
	// RegionList
	customConfigStr := `
		{
			"sdkenv":        "2",
			"checkdevice":   "false",
			"loadPatch":     "false",
			"showexception": "false",
			"regionConfig":  "pm|fk|add",
			"downloadMode":  "0",
		}
	`
	customConfig := []byte(customConfigStr)
	endec.Xor(customConfig, dispatchKey)
	serverList := make([]*proto.RegionSimpleInfo, 0)
	server := &proto.RegionSimpleInfo{
		Name:        "os_usa",
		Title:       "America",
		Type:        "DEV_PUBLIC",
		DispatchUrl: "https://osusadispatch.yuanshen.com/query_cur_region",
	}
	serverList = append(serverList, server)
	regionList := new(proto.QueryRegionListHttpRsp)
	regionList.RegionList = serverList
	regionList.ClientSecretKey = dispatchSeed
	regionList.ClientCustomConfigEncrypted = customConfig
	regionList.EnableLoginPc = true
	return regionCurr, regionList
}

// 新的region构建方式已经不需要读取query_region_list和query_cur_region文件了 并跳过了版本校验的blk补丁下载
func _InitRegion(log *logger.Logger, kcpAddr string, kcpPort int) (*proto.QueryCurrRegionHttpRsp, *proto.QueryRegionListHttpRsp) {
	// TODO 总有一天要把这些烦人的数据全部自己构造别再他妈的读文件了
	// 加载文件
	regionListKeyFile, err := ioutil.ReadFile("static/query_region_list_key")
	if err != nil {
		log.Error("open query_region_list_key error")
		return nil, nil
	}
	regionCurrKeyFile, err := ioutil.ReadFile("static/query_cur_region_key")
	if err != nil {
		log.Error("open query_cur_region_key error")
		return nil, nil
	}
	regionListFile, err := ioutil.ReadFile("static/query_region_list")
	if err != nil {
		log.Error("open query_region_list error")
		return nil, nil
	}
	regionCurrFile, err := ioutil.ReadFile("static/query_cur_region")
	if err != nil {
		log.Error("open query_cur_region error")
		return nil, nil
	}
	regionListKeyBin, err := base64.StdEncoding.DecodeString(string(regionListKeyFile))
	if err != nil {
		log.Error("decode query_region_list_key error")
		return nil, nil
	}
	regionCurrKeyBin, err := base64.StdEncoding.DecodeString(string(regionCurrKeyFile))
	if err != nil {
		log.Error("decode query_cur_region_key error")
		return nil, nil
	}
	regionListBin, err := base64.StdEncoding.DecodeString(string(regionListFile))
	if err != nil {
		log.Error("decode query_region_list error")
		return nil, nil
	}
	regionCurrBin, err := base64.StdEncoding.DecodeString(string(regionCurrFile))
	if err != nil {
		log.Error("decode query_cur_region error")
		return nil, nil
	}
	regionListKey := new(proto.QueryRegionListHttpRsp)
	err = pb.Unmarshal(regionListKeyBin, regionListKey)
	if err != nil {
		log.Error("Unmarshal QueryRegionListHttpRsp error")
		return nil, nil
	}
	regionCurrKey := new(proto.QueryCurrRegionHttpRsp)
	err = pb.Unmarshal(regionCurrKeyBin, regionCurrKey)
	if err != nil {
		log.Error("Unmarshal QueryCurrRegionHttpRsp error")
		return nil, nil
	}
	regionList := new(proto.QueryRegionListHttpRsp)
	err = pb.Unmarshal(regionListBin, regionList)
	if err != nil {
		log.Error("Unmarshal QueryRegionListHttpRsp error")
		return nil, nil
	}
	regionCurr := new(proto.QueryCurrRegionHttpRsp)
	err = pb.Unmarshal(regionCurrBin, regionCurr)
	if err != nil {
		log.Error("Unmarshal QueryCurrRegionHttpRsp error")
		return nil, nil
	}
	secretKey, err := ioutil.ReadFile("static/dispatchSeed.bin")
	if err != nil {
		log.Error("open dispatchSeed.bin error")
		return nil, nil
	}
	// RegionCurr
	regionInfo := regionCurr.GetRegionInfo()
	regionInfo.GateserverIp = kcpAddr
	regionInfo.GateserverPort = uint32(kcpPort)
	regionInfo.SecretKey = secretKey
	regionCurr.RegionInfo = regionInfo
	// RegionList
	serverList := make([]*proto.RegionSimpleInfo, 0)
	server := &proto.RegionSimpleInfo{
		Name:        "os_usa",
		Title:       "America",
		Type:        "DEV_PUBLIC",
		DispatchUrl: "https://osusadispatch.yuanshen.com/query_cur_region",
	}
	serverList = append(serverList, server)
	regionList.RegionList = serverList
	regionList.ClientSecretKey = regionListKey.ClientSecretKey
	regionList.ClientCustomConfigEncrypted = regionListKey.ClientCustomConfigEncrypted
	return regionCurr, regionList
}
