package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AutoReload  bool
	ShowSwagger bool
	GinMode     string
	Domain      struct {
		Api     string
		H5      string
		Backend string
	}
	Http struct {
		Port int
	}
	Https struct {
		On   bool
		Port int
		Host string
	}
	DB struct {
		DriverName     string
		DataSourceName string
		ShowSQL        bool
		SQLLogFile     string
	}
	File struct {
		WebUploadRoot   string
		WebRelativePath string
	}
	Static map[string]string
	Openai struct {
		ApiKey string
	}
	JWT struct {
		AccessTokenExpires  int64
		RefreshTokenExpires int64
		PrivateKey          string
		PublicKey           string
		SkipPaths           []string
	}
	Casbin struct {
		Model string
	}
	Redis struct {
		Prefix   string
		Addr     string
		Password string
	}
	Nsq struct {
		Prefix string
		Addr   string
	}
	LogFilePath string
	Wechat      struct {
		OfficialAccountService struct { //公众服务号
			AppID         string
			Secret        string
			OriginalID    string
			Token         string
			EncodedAESKey string
		}
		OfficialAccountSubscription struct { //公众订阅号
			AppID         string
			Secret        string
			OriginalID    string
			Token         string
			EncodedAESKey string
		}
		MiniProgram struct { //小程序
			AppID         string
			Secret        string
			Token         string
			EncodedAESKey string
		}
		MiniProgramPay struct { //微信小程序支付
			AppID     string
			MchId     string
			ApiKey    string
			ApiKeyV3  string
			NotifyUrl string
		}
		H5 struct { //H5
			AppID     string
			ApiKey    string
			MchId     string
			NotifyUrl string
		}
		APP struct { //APP应用
			AppID     string
			ApiKey    string
			MchID     string
			NotifyUrl string
		}
	}
	Douyin struct {
		MiniProgram struct { // 抖音小程序
			AppID  string
			Secret string
		}
	}
	Aliyun struct {
		Env          string // 只有值为production时，才会真正发送短信
		RegionID     string // 地域ID
		AccessKeyID  string
		AccessSecret string
		SignName     string
		TemplateCode string
		Oss          struct {
			AccessKeyID     string //
			AccessKeySecret string //
			BucketName      string
			Endpoint        string
			Path            string
			Url             string
		}
	}
}

var CFG *Config

func NewConfig() *Config {
	config := &Config{}

	ReadConfig(config)

	return config
}

// ReadConfig 设置读取的配置文件
func ReadConfig(config *Config) {
	configName := "configs"
	if os.Getenv("CONFIGS_NAME") != "" {
		configName = os.Getenv("CONFIGS_NAME")
	}
	//fmt.Printf("当前环境变量CONFIG_NAME为%s，使用%s.json配置文件\n", configName, configName)
	configName = configName + ".json"

	configDirs := []string{"configs/", "./", "../../configs/"}
	for _, v := range configDirs {
		bytes, err := os.ReadFile(v + configName)
		if err != nil {
			continue
		}

		err = json.Unmarshal(bytes, &config)
		if err != nil {
			panic(err)
		}
		break
	}

	currentDir, _ := os.Getwd()
	config.File.WebUploadRoot = strings.ReplaceAll(config.File.WebUploadRoot, "${currentDir}", currentDir)
	if err := os.MkdirAll(config.File.WebUploadRoot, os.ModePerm); err != nil {
		fmt.Println(err)
	}
}
