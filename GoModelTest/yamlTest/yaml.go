package yaml

import (
	"errors"
	"strings"
	"sync"

	"github.com/deepglint/flowservice/util/filetool"
	Yaml "gopkg.in/yaml.v2"
)

/*
	作者：宋言
	读写yaml
*/

// type ObjStorage map[interface{}]interface{}
type YamlConfig struct {
	ConfigFile string
	WriteLock  *sync.Mutex
	Kvmap      map[interface{}]interface{}
}

func NewYamlConfig(file string) (yaml *YamlConfig, err error) {
	content, err := filetool.ReadFileToBytes(file)
	if err != nil {
		return nil, err
	}
	kvmap := make(map[interface{}]interface{})
	err = Yaml.Unmarshal(content, &kvmap)
	yaml = &YamlConfig{file, new(sync.Mutex), kvmap}
	return
}

// TODO
// write protect
func (i *YamlConfig) WriteYamlConfig() (err error) {
	content, _ := Yaml.Marshal(&i.Kvmap)
	i.WriteLock.Lock()
	_, err = filetool.WriteBytesToFile(i.ConfigFile, content)
	i.WriteLock.Unlock()
	return
}

// get key-value
func (i *YamlConfig) GetValue(keyName string) (interface{}, error) {
	keyName = strings.Trim(keyName, " /")
	keysl := strings.Split(keyName, "/")
	if len(keysl) < 1 {
		return nil, errors.New("Request Key illegal!")
	}
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
	i.WriteYamlConfig()
	return true, nil
}

func (i *YamlConfig) DeleteValue(keyName string) error {
	keyName = strings.Trim(keyName, " /")
	keysl := strings.Split(keyName, "/")
	if len(keysl) < 1 {
		return errors.New("Request Key illegal!")
	} else if len(keysl) == 1 {
		delete(i.Kvmap, keysl[0])
		i.WriteYamlConfig()
		return nil
	}

	kvt := &i.Kvmap
	for n := 0; n < len(keysl); n++ {
		key := keysl[n]
		if v, ok := (*kvt)[key]; ok {
			switch vv := v.(type) {
			case map[interface{}]interface{}:
				kvt = &vv
				continue
			case interface{}:
				continue
			}
		}
		// not found the key
		return errors.New("Request Key not match!")
	}
	delete(*kvt, keysl[len(keysl)-1])
	i.WriteYamlConfig()
	return nil
}

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
		return it.(string), nil
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
			return it.(string)
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
