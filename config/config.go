package config

import (
	"demo-gogo/utils"
	"encoding/json"
	"flag"
	"github.com/jinzhu/configor"
	log "github.com/wonderivan/logger"
	"os"
	"strconv"
)

const (
	CONF_OSS_NGINX = "nginx"
	CONF_OSS_MINIO = "minio"
)

var Conf *Config

var DefaultConfig = Config{
	APP: APP{
		Name:               "demo-gogo",
		IP:                 "127.0.0.1",
		Port:               9094,
		Mode:               "release",
		SkipAuthentication: false,
		ContextPath:        "/api",
		UploadBasePath:     "files/any_files/",
		UploadFileSize:     10485760,
	},
	DB: DB{
		Name:            "demo-gogo",
		Host:            "120.46.48.255",
		User:            "root",
		Password:        "123456",
		Port:            5432,
		MaxIdleConnects: 10,
		MaxOpenConnects: 1024,
		InitTable:       true,
	},
	Redis: Redis{
		Host:               "120.46.48.255",
		Port:               6379,
		Password:           "",
		MaTeachProgressKey: "ma_teach_progress",
	},
	OSS: OSS{
		Type:         "nginx",
		Endpoint:     "http://120.46.48.255:38888",
		FileSavePath: "/mount/data5/hsr3_save_files",
		User:         "admin",
		Password:     "hsradmin",
		UseSSL:       false,
		Bucket:       "mateach",
	},
	HSR: HSR{
		IP:   "127.0.0.1",
		Port: 9999,
	},
	Compute: Compute{
		IP:           "120.46.48.255",
		Port:         50006,
		Url:          "/api/stop_compute",
		PullInterval: 0,
	},
	Nerf: Nerf{
		IP:         "120.46.48.255",
		Port:       7006,
		ViewerPort: 65500,
		Quality:    75,
		Threshold:  0.5,
	},
	Match: Match{
		IP:   "120.46.48.255",
		Port: 1909,
	},
	Robot: Robot{
		IP:   "120.46.48.255",
		Port: 1909,
	},
	Emq: Emq{
		Broker: "tcp://120.46.48.255:1883",
	},
}

type Config struct {
	APP     APP     `json:"app" yaml:"app"`
	DB      DB      `json:"db" yaml:"db"`
	Redis   Redis   `json:"redis" yaml:"redis"`
	OSS     OSS     `json:"oss" yaml:"oss"`
	HSR     HSR     `json:"hsr" yaml:"hsr"`
	Compute Compute `json:"compute" yaml:"compute"`
	Nerf    Nerf    `json:"nerf" yaml:"nerf"`
	Match   Match   `json:"match" yaml:"match"`
	Robot   Robot   `json:"robot" yaml:"robot"`
	Emq     Emq     `json:"emq" yaml:"emq"`
}

type APP struct {
	Name               string `yaml:"name" json:"name"`
	IP                 string `yaml:"ip" json:"ip"`
	Port               int    `yaml:"port" json:"port"`
	Mode               string `yaml:"mode" json:"mode"`
	SkipAuthentication bool   `yaml:"skip_authentication" json:"skip_authentication"`
	ContextPath        string `yaml:"context_path" json:"context_path"`
	UploadBasePath     string `yaml:"upload_base_path" json:"upload_base_path"`
	UploadFileSize     int    `yaml:"upload_file_size" json:"upload_file_size"`
}

type DB struct {
	Name            string `yaml:"name" json:"name"`
	Host            string `yaml:"host" json:"host"`
	User            string `yaml:"user" json:"user"`
	Password        string `yaml:"password" json:"password"`
	Port            uint   `yaml:"port" json:"port"`
	MaxIdleConnects int    `yaml:"max_idle_connects" json:"max_idle_connects"`
	MaxOpenConnects int    `yaml:"max_open_connects" json:"max_open_connects"`
	InitTable       bool   `yaml:"init_table" json:"init_table"`
}

type Redis struct {
	Host               string `yaml:"host" json:"host"`
	Port               int    `yaml:"port" json:"port"`
	Password           string `yaml:"password" json:"password"`
	MaTeachProgressKey string `yaml:"ma_teach_progress_key" json:"ma_teach_progress_key"`
}

type OSS struct {
	Type         string `yaml:"type" json:"type"`
	Endpoint     string `yaml:"endpoint" json:"endpoint"`
	FileSavePath string `yaml:"file_save_path" json:"file_save_path"`
	User         string `yaml:"user" json:"user"`
	Password     string `yaml:"password" json:"password"`
	UseSSL       bool   `yaml:"use_ssl" json:"use_ssl"`
	Bucket       string `yaml:"bucket" json:"bucket"`
}

type HSR struct {
	IP   string `yaml:"ip" json:"ip"`
	Port int    `yaml:"port" json:"port"`
}

type Compute struct {
	IP           string `yaml:"ip" json:"ip"`
	Port         int    `yaml:"port" json:"port"`
	Url          string `yaml:"url" json:"url"`
	PullInterval int    `yaml:"pull_interval" json:"pull_interval"`
}

type Nerf struct {
	IP         string  `yaml:"ip" json:"ip"`
	Port       int     `yaml:"port" json:"port"`
	ViewerPort int     `yaml:"viewer_port" json:"viewer_port"`
	Quality    int     `yaml:"quality" json:"quality"`
	Threshold  float64 `yaml:"threshold" json:"threshold"`
}

type Match struct {
	IP   string `yaml:"ip" json:"ip"`
	Port int    `yaml:"port" json:"port"`
}

type Robot struct {
	IP   string `yaml:"ip" json:"ip"`
	Port int    `yaml:"port" json:"port"`
}

type Emq struct {
	Broker string `yaml:"broker" json:"broker"`
}

func InitConfig() error {
	Conf = &DefaultConfig
	confPath := "./conf/config.yml"
	if utils.FileExist(confPath) {
		c := initConfLoader()
		log.Debug("加载用户自定义配置...")
		err := c.Load(Conf, confPath)
		if err != nil {
			return err
		}
	}
	// 启动命令参数覆盖默认配置
	appIP := flag.String("app_ip", "", "输入app的ip地址")
	appPort := flag.Int("app_port", 0, "输入app的端口号")
	dbHost := flag.String("db_host", "", "输入db的ip地址")
	hsrIP := flag.String("hsr_ip", "", "输入hsr的ip地址")
	flag.Parse()
	if *appIP != "" {
		Conf.APP.IP = *appIP
	}
	if *appPort != 0 {
		Conf.APP.Port = *appPort
	}
	// 现场部署，指定鹰眼ip（Nerf、中间组件都在鹰眼服务器上）
	if *hsrIP != "" {
		Conf.HSR.IP = *hsrIP
		Conf.Nerf.IP = *hsrIP
		Conf.Compute.IP = *hsrIP

		Conf.DB.Host = *hsrIP
		Conf.Redis.Host = *hsrIP
	}
	if *dbHost != "" {
		Conf.DB.Host = *dbHost
	}
	LoadConfFromEnv(Conf)
	log.Info("启动配置参数：")
	PrettyPrint(Conf)
	if !utils.Exists(Conf.APP.UploadBasePath) {
		err := os.MkdirAll(Conf.APP.UploadBasePath, 0777)
		if err != nil {
			log.Error("上传文件目录创建失败。err:[%#v]", err)
		}
	}
	return nil
}

func initConfLoader() *configor.Configor {
	config := configor.Config{
		AutoReload: true,
		AutoReloadCallback: func(config interface{}) {
			log.Info("配置文件热加载：")
			PrettyPrint(config)
		},
	}
	c := configor.New(&config)
	return c
}

func LoadConfFromEnv(conf *Config) {
	log.Debug("读取环境变量配置参数.env:[%#v]", os.Environ())
	if appIp, ok := os.LookupEnv("APP_IP"); ok {
		conf.APP.IP = appIp
	}
	if appPort, ok := os.LookupEnv("APP_PORT"); ok {
		port, err := strconv.Atoi(appPort)
		if err == nil {
			conf.APP.Port = port
		}
	}
	if dbHost, ok := os.LookupEnv("DB_HOST"); ok {
		conf.DB.Host = dbHost
	}
	if dbPort, ok := os.LookupEnv("DB_PORT"); ok {
		port, err := strconv.Atoi(dbPort)
		if err == nil {
			conf.DB.Port = uint(port)
		}
	}
	if dbInit, ok := os.LookupEnv("DB_INIT_TABLE"); ok {
		conf.DB.InitTable = dbInit == "true"
	}

	if redisHost, ok := os.LookupEnv("REDIS_HOST"); ok {
		conf.Redis.Host = redisHost
	}
	if redisPort, ok := os.LookupEnv("REDIS_PORT"); ok {
		port, err := strconv.Atoi(redisPort)
		if err == nil {
			conf.Redis.Port = port
		}
	}

	if fileSavePath, ok := os.LookupEnv("FILE_SAVE_PATH"); ok {
		conf.OSS.FileSavePath = fileSavePath
	}
	// todo 临时更新 不读环境变量
	if fileSaveNginx, ok := os.LookupEnv("FILE_SAVE_NGINX"); ok {
		conf.OSS.Endpoint = fileSaveNginx
		conf.OSS.Type = "nginx"
	}

	if hsrIp, ok := os.LookupEnv("HSR_IP"); ok {
		conf.HSR.IP = hsrIp
	}
	if hsrPort, ok := os.LookupEnv("HSR_PORT"); ok {
		port, err := strconv.Atoi(hsrPort)
		if err == nil {
			conf.HSR.Port = port
		}
	}

	if computeIp, ok := os.LookupEnv("COMPUTE_IP"); ok {
		conf.Compute.IP = computeIp
	}
	if computePort, ok := os.LookupEnv("COMPUTE_PORT"); ok {
		port, err := strconv.Atoi(computePort)
		if err == nil {
			conf.Compute.Port = port
		}
	}

	if nerfIp, ok := os.LookupEnv("NERF_IP"); ok {
		conf.Nerf.IP = nerfIp
	}

	if quality, ok := os.LookupEnv("NERF_QUALITY"); ok {
		qualityInt, err := strconv.Atoi(quality)
		if err == nil {
			conf.Nerf.Quality = qualityInt
		}
	}

	if threshold, ok := os.LookupEnv("NERF_THRESHOLD"); ok {
		thresholdFloat, err := strconv.ParseFloat(threshold, 64)
		if err == nil {
			conf.Nerf.Threshold = thresholdFloat
		}
	}

	if matchIp, ok := os.LookupEnv("MATCH_IP"); ok {
		conf.Match.IP = matchIp
	}

	if matchPort, ok := os.LookupEnv("MATCH_PORT"); ok {
		port, err := strconv.Atoi(matchPort)
		if err == nil {
			conf.Match.Port = port
		}
	}

	if broker, ok := os.LookupEnv("MMQ_BROKER"); ok {
		conf.Emq.Broker = broker
	}
}

func PrettyPrint(data interface{}) {
	p, _ := json.MarshalIndent(data, "", "\t")
	log.Info("%s \n", p)
}
