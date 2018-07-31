package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var testYaml = `
servers:
  -
    192.168.33.10:1111

  -
    192.168.33.10:2222

  -
    192.168.33.10:3333

balancing: ip-hash
`

func TestLoadConfig(t *testing.T) {
	const filename string = "test.yaml"
	defer deleteFile(filename)
	createFile(filename, []byte(testYaml))

	c, err := LoadConfig(filename)
	if err != nil {
		t.Errorf("LoadConfig is faild. error: %v, c: %v", err, c)
	}

	expected := &Config{
		Servers: Servers{
			"192.168.33.10:1111",
			"192.168.33.10:2222",
			"192.168.33.10:3333",
		},
		Balancing: "ip-hash",
	}

	if !reflect.DeepEqual(expected, c) {
		t.Errorf("LoadConfig is wrong. expected: %v, got: %v", expected, c)
	}
}

func TestToStringSlice(t *testing.T) {
	tests := []struct {
		servers  Servers
		expected []string
	}{
		{
			servers: Servers{
				"192.168.33.10:1111",
				"192.168.33.10:2222",
				"192.168.33.10:3333",
			},
			expected: []string{
				"192.168.33.10:1111",
				"192.168.33.10:2222",
				"192.168.33.10:3333",
			},
		},
	}

	for i, test := range tests {
		got := test.servers.ToStringSlice()

		if !reflect.DeepEqual(test.expected, got) {
			t.Errorf("tests[%d] - ToStringSlice is wrong. expected: %v, got: %v", i, test.expected, got)
		}
	}
}

func createFile(filename string, data []byte) {
	ioutil.WriteFile(filename, data, os.ModePerm)
}

func deleteFile(filename string) {
	os.Remove(filename)
}