package util

import (
	"fmt"
	"gopkg.in/ini.v1"
)

/**
 * 获取签名密钥
 */
func GetSignKey() ( []byte,error) {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		Error(fmt.Sprintf("获取配置文件失败：%v\n",err.Error()))
		return nil,err
	}

	return []byte(cfg.Section("key").Key("KEY_SIGN").String()),nil
}
