package cfg

import (
	"os"
	"testing"
)

var testTomlFile string = "test.toml"

func TestLoadCfgFile(t *testing.T) {
	err := Initialize(testTomlFile)
	if err != nil {
		t.Errorf("读取配置文件%q失败,err=%v\n", testTomlFile, err)
		os.Exit(1)
	}

	t.Run("TestExistedData", TestExistedData)
	t.Run("TestNotExistedData", TestNotExistedData)
}

/// 测试配置文件中存在的数据
func TestExistedData(t *testing.T) {
	/*
		global_String = "global_String"
		global_Int = 3
		global_Bool = true

		[server]
		server_String = "server_String"
		server_Int = 7
		server_Bool = false
	*/
	var existedData = map[string]interface{}{
		"global_String": "global_String",
		"global_Int":    3,
		"global_Bool":   true,

		"server.server_String": "server_String",
		"server.server_Int":    7,
		"server.server_Bool":   false,
	}

	for k, v := range existedData {
		switch v.(type) {
		case string:
			if MustGetString(k) != v {
				t.Errorf("MustGetString(%v) != %v, failed", k, v)
			}
		case int:
			if MustGetInt(k) != v {
				t.Errorf("MustGetInt(%v) != %v, failed", k, v)
			}
		case bool:
			if MustGetBool(k) != v {
				t.Errorf("MustGetBool(%v) != %v, failed", k, v)
			}
		default:
			t.Errorf("unknown data type(%v) != %v, failed", k, v)
		}
	}
}

/// 测试配置文件中不存在的数据，都是默认值
func TestNotExistedData(t *testing.T) {
	// 下面的key在配置文件中没有配置，都用golang默认值
	var notExistedData = map[string]interface{}{
		"no_global_String": "",
		"no_global_Int":    0,
		"no_global_Bool":   false,

		"no.server_String":      "",
		"no.server_Int":         0,
		"no.server_Bool":        false,
		"server.no_server_Bool": false,
	}

	for k, v := range notExistedData {
		switch v.(type) {
		case string:
			if GetString(k) != v {
				t.Errorf("GetString(%v) != %v, failed", k, v)
			}
		case int:
			if GetInt(k) != v {
				t.Errorf("GetInt(%v) != %v, failed", k, v)
			}
		case bool:
			if GetBool(k) != v {
				t.Errorf("GetBool(%v) != %v, failed", k, v)
			}
		default:
			t.Errorf("unknown data type(%v) != %v, failed", k, v)
		}
	}

}
