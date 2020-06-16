package yaml

import (
	"errors"
	Yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

/*
	作者：宋言
	读写yaml
*/

// type ObjStorage map[interface{}]interface{}
type YamlConfig struct {
	ConfigFile string
	Lock       *sync.RWMutex
	Kvmap      map[interface{}]interface{}
}

// LibraT.yaml
type LibraT struct {
	TfCard struct {
		Device     string `yaml:"device"`
		Enable     bool   `yaml:"enable"`
		MountPoint string `yaml:"mountpoint"`
	}
}

func NewLibraTConf(path string) (*LibraT, error) {
	content, err := ReadFileToBytes(path)
	if err != nil {
		return nil, err
	}
	var libraConf LibraT
	err = Yaml.Unmarshal(content, &libraConf)
	return &libraConf, err
}

func NewYamlConfig(file string) (*YamlConfig, error) {
	content, err := ReadFileToBytes(file)
	if err != nil {
		return nil, err
	}
	kvmap := make(map[interface{}]interface{})
	err = Yaml.Unmarshal(content, &kvmap)
	yaml := &YamlConfig{file, new(sync.RWMutex), kvmap}
	return yaml, err
}

// TODO
// write protect
func (i *YamlConfig) writeYamlConfig() (err error) {
	content, _ := Yaml.Marshal(&i.Kvmap)
	_, err = WriteBytesToFile(i.ConfigFile, content)
	return
}

// get key-value
func (i *YamlConfig) GetValue(keyName string) (interface{}, error) {
	keyName = strings.Trim(keyName, " /")
	keysl := strings.Split(keyName, "/")
	if len(keysl) < 1 {
		return nil, errors.New("Request Key illegal!")
	}

	i.Lock.RLock()
	defer i.Lock.RUnlock()
	kvt := &i.Kvmap
	for n, key := range keysl {
		if v, ok := (*kvt)[key]; ok {
			switch vv := v.(type) {
			case map[interface{}]interface{}:
				kvt = &vv
				continue
			}
			if n == len(keysl)-1 {
				return v, nil
			}
		}
		return nil, errors.New("Request Key not match!")
	}
	//return nil, errors.New("Request Key not match!")
	return *kvt, nil
	// 	return a map ,not support yet
}

// set key-value
func (i *YamlConfig) SetValue(keyName string, value interface{}) (bool, error) {
	keyName = strings.Trim(keyName, " /")
	keysl := strings.Split(keyName, "/")
	if len(keysl) < 1 {
		return false, errors.New("Request Key illegal!")
	}

	i.Lock.Lock()
	defer i.Lock.Unlock()
	kvt := &i.Kvmap
	for n := 0; n < len(keysl); n++ {
		key := keysl[n]
		if v, ok := (*kvt)[key]; ok {
			switch vv := v.(type) {
			case map[interface{}]interface{}:
				kvt = &vv
				continue
			}
		}
		// lost key
		for ; n < len(keysl)-1; n++ {
			key := keysl[n]
			newmap := make(map[interface{}]interface{})
			(*kvt)[key] = newmap
			kvt = &newmap
		}
		break
	}
	(*kvt)[keysl[len(keysl)-1]] = value
	i.writeYamlConfig()
	return true, nil
}

func (i *YamlConfig) DeleteValue(keyName string) error {
	keyName = strings.Trim(keyName, " /")
	keysl := strings.Split(keyName, "/")

	i.Lock.Lock()
	defer i.Lock.Unlock()
	if len(keysl) < 1 {
		return errors.New("Request Key illegal!")
	} else if len(keysl) == 1 {
		delete(i.Kvmap, keysl[0])
		i.writeYamlConfig()
		return nil
	}

	kvt := &i.Kvmap
	for n := 0; n < len(keysl)-1; n++ {
		key := keysl[n]
		if v, ok := (*kvt)[key]; ok {
			switch vv := v.(type) {
			case map[interface{}]interface{}:
				kvt = &vv
				continue
			}
		}
		// not found the key
		return nil
	}
	delete(*kvt, keysl[len(keysl)-1])
	i.writeYamlConfig()
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// get key-value
func (i *YamlConfig) GetString(keyName string) (string, error) {
	it, err := i.GetValue(keyName)
	if err != nil {
		return "", err
	}
	if it == nil {
		return "", errors.New("Request Value not present!")
	}
	switch it.(type) {
	case string:
		return strings.TrimSpace(it.(string)), nil
	}
	return "", errors.New("Request Value not present!")
}

// get key-value
// nil/""/null interface is all set to def,thus "" is dangerous
func (i *YamlConfig) ValidString(keyName string, def string) string {
	it, err := i.GetValue(keyName)
	if err != nil {
		return def
	}
	if it == nil {
		return def
	}
	switch it.(type) {
	case string:
		if len(it.(string)) > 0 {
			return strings.TrimSpace(it.(string))
		}
	}
	return def
}

// get key-value
func (i *YamlConfig) GetInteger(keyName string) (int, error) {
	it, err := i.GetValue(keyName)
	if err != nil {
		return 0, err
	}
	if it == nil {
		return 0, errors.New("Request Value not present!")
	}
	switch it.(type) {
	case int:
		return it.(int), nil
	}
	return 0, errors.New("Request Value not present!")
}

func (i *YamlConfig) ValidInteger(keyName string, def int) int {
	it, err := i.GetValue(keyName)
	if err != nil {
		return def
	}
	if it == nil {
		return def
	}
	switch it.(type) {
	case int:
		return it.(int)
	}
	return def
}

// - key1: value1
// - key2: value2
// - key3: value3
func (i *YamlConfig) GetKVdata(keyName string) (map[string]string, error) {
	it, err := i.GetValue(keyName)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	switch vv := it.(type) {
	case []interface{}:
		for _, u := range vv {
			switch pp := u.(type) {
			case map[interface{}]interface{}:
				for k, v := range pp {
					ret[k.(string)] = strings.TrimSpace(v.(string))
				}
			}
		}
	}
	return ret, nil
}

// -
//   date: date1
//   value: value1
// -
//   date: date2
//   value: value2

func (i *YamlConfig) GetKVdataArray(keyName string) ([]map[string]string, error) {
	it, err := i.GetValue(keyName)
	if err != nil {
		return nil, err
	}
	if it == nil {
		return nil, errors.New("Request Value not present!")
	}

	ret := make([]map[string]string, 0)
	switch vv := it.(type) {
	case []interface{}:
		for _, u := range vv {
			switch pp := u.(type) {
			case map[interface{}]interface{}:
				m := make(map[string]string)
				for k, v := range pp {
					m[k.(string)] = strings.TrimSpace(v.(string))
				}
				ret = append(ret, m)
			}
		}
	}
	return ret, nil
}

/*
 - 1
 - 2
 - 3
	读取int数组
*/
func (i *YamlConfig) GetIntegerArray(keyName string) ([]int, error) {
	it, err := i.GetValue(keyName)
	if err != nil {
		return nil, err
	}
	if it == nil {
		return nil, errors.New("Request Value not present!")
	}

	ret := make([]int, 0)
	switch vv := it.(type) {
	case []interface{}:
		for _, u := range vv {
			switch u.(type) {
			case int:
				ret = append(ret, u.(int))
			}
		}
	}
	return ret, nil
}

/*
 - a
 - b
 - c
	读取字符串数组
	已知情况,在yaml自动识别为int或者其类型情况下读取识别
 - 1
 - "1"
	后一种处理才行
*/
func (i *YamlConfig) GetStringArray(keyName string) ([]string, error) {
	it, err := i.GetValue(keyName)
	if err != nil {
		return nil, err
	}
	if it == nil {
		return nil, errors.New("Request Value not present!")
	}

	ret := make([]string, 0)
	switch vv := it.(type) {
	case []interface{}:
		for _, u := range vv {
			switch u.(type) {
			case string:
				ret = append(ret, strings.TrimSpace(u.(string)))
			}
		}
	}
	return ret, nil
}

// WriteBytesToFile saves content type '[]byte' to file by given path.
// It returns error when fail to finish operation.
func WriteBytesToFile(filePath string, b []byte) (int, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(b)
}

// ReadFileToBytes reads data type '[]byte' from file by given path.
// It returns error when fail to finish operation.
func ReadFileToBytes(filePath string) ([]byte, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte(""), err
	}
	return b, nil
}
