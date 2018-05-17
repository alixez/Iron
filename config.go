package Iron

import (
	"io/ioutil"
	"os"
	"path"
	"sort"
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
	configs     ConfigDict
}

/*
ConfigDict class
配置字典类
*/
type ConfigDict map[string]interface{}

func (cd ConfigDict) GetDict(field string) ConfigDict {
	return cd[field].(ConfigDict)
}

func (cd ConfigDict) GetString(field string) string {
	return cd[field].(string)
}

func (cd ConfigDict) GetInt(field string) int {
	return cd[field].(int)
}

func (cd ConfigDict) GetInt64(field string) int64 {
	return cd[field].(int64)
}

func (cd ConfigDict) GetBool(field string) bool {
	return cd[field].(bool)
}

/*
Init a env object
初始化一个配置对象
*/
func (env *Env) Init(configs ConfigDict) {
	env.AppName = configs.GetString("appname")
	env.Version = configs.GetString("version")
	env.Environment = configs.GetString("environment")
	env.Debug = configs.GetBool("debug")
	env.configs = configs
	env.HasLoaded = true
}

type fileListSlice []string

func (s fileListSlice) Len() int      { return len(s) }
func (s fileListSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s fileListSlice) Less(i, j int) bool {

	sortItems := map[int]int{}

	switch s[i] {
	case "default.yaml":
		sortItems[i] = 0
	case "default.dev.yaml":
		sortItems[i] = 1
	case "default.proc.yaml":
		sortItems[i] = 2
	default:
		sortItems[i] = len(s[i])
	}

	switch s[j] {
	case "default.yaml":
		sortItems[j] = 0
	case "default.dev.yaml":
		sortItems[j] = 1
	case "default.proc.yaml":
		sortItems[j] = 2
	default:
		sortItems[j] = len(s[i])
	}

	return sortItems[i] < sortItems[j]
}

/*
Get config like a.b.c
*/
func (env *Env) Get(query string) interface{} {
	querySlice := strings.Split(query, ".")
	value := env.configs[querySlice[0]]
	// if env.Environment == "development" {
	// 	value = env.devConfigs[querySlice[0]]
	// }
	// if len(querySlice) == 1 {
	// 	return value
	// }

	for _, field := range querySlice[1:] {
		if v, ok := value.(map[interface{}]interface{}); ok {
			value = v[field]
			continue
		}
		value = nil
	}
	return value
}

// func (env *Env) GetMap(query string) map[interface{}]interface{} {
// 	return env.Get(query).(map[interface{}]interface{})
// }

func (env *Env) GetDict(query string) ConfigDict {
	return env.Get(query).(ConfigDict)
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

func (env *Env) GetBool(query string) bool {
	return env.Get(query).(bool)
}

/*
Set config
*/
func (env *Env) Set(field string, value interface{}) {
	// if env.Environment == "development" {
	// 	env.devConfigs[field] = value
	// } else if env.Environment == "production" {
	// 	env.prodConfigs[field] = value
	// }
	env.configs[field] = value
}

var systemEnv = &Env{
	HasLoaded: false,
}

func LoadApplicationEnv(root string, environment string) (env *Env) {
	if systemEnv.HasLoaded {
		env = systemEnv
		return
	}
	// fmt.Println("...开始加载配置文件...")
	dirPath := path.Join(root, "config")
	filepathList, err := utils.ListDir(dirPath, "yaml")
	if err != nil {
		panic("获取配置文件列表时发生错误")
	}
	if !utils.ArrayContainer(filepathList, "default.yaml") {
		filepathList = append(filepathList, "default.yaml")
		f, _ := os.Create(path.Join(dirPath, "env.yaml"))
		f.WriteString("appname: iron\r\nversion: v1.0")
	}
	if !utils.ArrayContainer(filepathList, "default.dev.yaml") {
		filepathList = append(filepathList, "default.dev.yaml")
		f, _ := os.Create(path.Join(dirPath, "default.dev.yaml"))
		f.WriteString("# 在下面填写开发环境的配置, 此文件的配置会覆盖默认配置\r\ndebug: true")
	}

	if !utils.ArrayContainer(filepathList, "default.proc.yaml") {
		filepathList = append(filepathList, "default.proc.yaml")
		f, _ := os.Create(path.Join(dirPath, "default.proc.yaml"))
		f.WriteString("# 在下面填写生产环境的配置, 此文件的配置会覆盖默认配置\r\ndebug: false")
	}

	// filepathList := []string{
	// 	"env.yaml",
	// 	"default.yaml",
	// }
	// masterConfigs := map[interface{}]interface{}{
	// 	"environment": environment,
	// 	"version":     "v0.1",
	// 	"appname":     "demo",
	// }
	masterConfigs := ConfigDict{
		"environment": environment,
		"version":     "v0.1",
		"appname":     "demo",
	}

	// 对文件列表进行排序
	// default.yaml
	// default.dev.yaml
	// default.proc.yaml
	filepathSlice := fileListSlice{}
	for _, v := range filepathList {
		filepathSlice = append(filepathSlice, v)
	}

	sort.Stable(filepathSlice)

	for _, v := range filepathSlice {
		// 获取文件名
		masterConfigs = readConfigFromFile(path.Join(dirPath, v), masterConfigs)
	}
	systemEnv.Init(masterConfigs)
	env = systemEnv
	return
}

// func mergeConfig(m1 interface{}, m2 interface{}) interface{} {
// 	if m1 == nil {
// 		m1 = map[interface{}]interface{}{}
// 	}
// 	if m2 == nil {
// 		m2 = map[interface{}]interface{}{}
// 	}
// 	m3 := m2.(map[interface{}]interface{})
// 	for k, v := range m1.(map[interface{}]interface{}) {
// 		m3[k] = v
// 	}
// 	return m3
// }

func mergeConfig(m1 ConfigDict, m2 ConfigDict) ConfigDict {
	if m1 == nil {
		m1 = ConfigDict{}
	}

	if m2 == nil {
		m2 = ConfigDict{}
	}

	m3 := m2

	for k, v := range m1 {
		m3[k] = v
	}

	return m3
}

func readConfigFromFile(filepath string, out ConfigDict) ConfigDict {
	configs := ConfigDict{}
	configByte, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configByte, &configs)
	if err != nil {
		panic(err)
	}

	return mergeConfig(configs, out)
}

// func readConfigFromFile(filepath string, out map[interface{}]interface{}) interface{} {
// 	configs := make(map[interface{}]interface{})
// 	configByte, err := ioutil.ReadFile(filepath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = yaml.Unmarshal(configByte, &configs)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return mergeConfig(configs, out)
// 	// if filepath == "default.yaml" {
// 	// 	// out["environment"] = configs["environment"]
// 	// 	out["version"] = configs["version"]
// 	// 	out["appname"] = configs["appname"]
// 	// } else {
// 	// 	// out["production"] = mergeConfig(configs["production"], out["production"])
// 	// 	// out["development"] = mergeConfig(configs["development"], out["development"])
// 	// }
// }
