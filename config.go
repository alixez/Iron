package Iron

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/alixez/Iron/utils"
	"gopkg.in/yaml.v2"
)

/*
Env class
环境配置类
*/
type Env struct {
	AppName     string
	Version     string
	Environment string
	Debug       bool
	HasLoaded   bool
	devConfigs  map[interface{}]interface{}
	prodConfigs map[interface{}]interface{}
}

/*
Init a env object
初始化一个配置对象
*/
func (env *Env) Init(configs map[interface{}]interface{}) {
	env.AppName = configs["appname"].(string)
	env.Version = configs["version"].(string)
	env.Environment = configs["environment"].(string)
	env.Debug = env.Environment != "production"
	env.devConfigs = configs["development"].(map[interface{}]interface{})
	env.prodConfigs = configs["production"].(map[interface{}]interface{})
	env.HasLoaded = true
}

/*
Get config like a.b.c
*/
func (env *Env) Get(query string) interface{} {
	querySlice := strings.Split(query, ".")
	value := env.prodConfigs[querySlice[0]]
	if env.Environment == "development" {
		value = env.devConfigs[querySlice[0]]
	}
	if len(querySlice) == 1 {
		return value
	}

	for _, field := range querySlice[1:] {
		if v, ok := value.(map[interface{}]interface{}); ok {
			value = v[field]
			continue
		}
		value = nil
	}
	return value
}

func (env *Env) GetMap(query string) map[interface{}]interface{} {
	return env.Get(query).(map[interface{}]interface{})
}

func (env *Env) GetString(query string) string {
	return env.Get(query).(string)
}

func (env *Env) GetInt(query string) int64 {
	return env.Get(query).(int64)
}

func (env *Env) GetFloat(query string) float64 {
	return env.Get(query).(float64)
}

/*
Set config
*/
func (env *Env) Set(field string, value interface{}) {
	if env.Environment == "development" {
		env.devConfigs[field] = value
	} else if env.Environment == "production" {
		env.prodConfigs[field] = value
	}
}

var systemEnv = &Env{
	HasLoaded: false,
}

func LoadApplicationEnv() (env *Env) {
	if systemEnv.HasLoaded {
		env = systemEnv
		return
	}
	fmt.Println("...开始加载配置文件...")
	filepathList, err := utils.ListDir("config", "yaml")
	if err != nil {
		panic("获取配置文件列表时发生错误")
	}
	if !utils.ArrayContainer(filepathList, "env.yaml") {
		filepathList = append(filepathList, "env.yaml")
		f, _ := os.Create(path.Join("config", "env.yaml"))
		f.WriteString("appname: arku\r\nversion: v1.0\r\nenvironment: development")
	}
	if !utils.ArrayContainer(filepathList, "default.yaml") {
		filepathList = append(filepathList, "default.yaml")
		f, _ := os.Create(path.Join("config", "default.yaml"))
		f.WriteString("production:\r\ndevelopment:\r\n")
	}

	// filepathList := []string{
	// 	"env.yaml",
	// 	"default.yaml",
	// }
	masterConfigs := map[interface{}]interface{}{
		"environment": "development",
		"version":     "v0.1",
		"appname":     "demo",
		"development": map[interface{}]interface{}{},
		"production":  map[interface{}]interface{}{},
	}
	for _, v := range filepathList {
		readConfigFromFile(v, masterConfigs)
	}
	systemEnv.Init(masterConfigs)
	env = systemEnv
	return
}

func mergeConfig(m1 interface{}, m2 interface{}) interface{} {
	if m1 == nil {
		m1 = map[interface{}]interface{}{}
	}
	if m2 == nil {
		m2 = map[interface{}]interface{}{}
	}
	m3 := m2.(map[interface{}]interface{})
	for k, v := range m1.(map[interface{}]interface{}) {
		m3[k] = v
	}
	return m3
}

func readConfigFromFile(filepath string, out map[interface{}]interface{}) {
	configs := make(map[interface{}]interface{})
	configByte, err := ioutil.ReadFile(path.Join("config", filepath))
	if err != nil {

		panic(err)
	}
	err = yaml.Unmarshal(configByte, &configs)
	if err != nil {
		panic(err)
	}
	if filepath == "env.yaml" {
		out["environment"] = configs["environment"]
		out["version"] = configs["version"]
		out["appname"] = configs["appname"]
	} else {
		out["production"] = mergeConfig(configs["production"], out["production"])
		out["development"] = mergeConfig(configs["development"], out["development"])
	}
}
