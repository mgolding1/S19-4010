package ReadConfig

import (
	"io/ioutil"
	"os"
	"testing"
)

// GlobalConfigData is the gloal configuration data.
// It holds all the data from the cfg.json file.
type GlobalConfigData struct {
	ExampeWithDefault string `default:"dflt-1"`
	SomePassword      string `default:"dflt-2"`
	CheckDefault      string `default:"dflt-3"`
}

var gCfg GlobalConfigData // global configuration data.

func TestMineBlock(t *testing.T) {

	tests := []struct {
		SetEnvName       string
		SetEnvVal        string
		FileName         string
		ExpectedPassword string
	}{
		{
			SetEnvName:       "MyPassword",
			SetEnvVal:        "xyzzy-3",
			FileName:         "../testdata/a.json",
			ExpectedPassword: "xyzzy-3",
		},
		{
			SetEnvName:       "Test2",
			SetEnvVal:        "xyzzy-2",
			FileName:         "../testdata/b.json",
			ExpectedPassword: "xyzzy-2",
		},
	}

	db1 = false // turn on output for debuging in ReadFile
	db2 = false // turn on output for debuging in SetFromEnv

	var home string
	if os.PathSeparator == '/' {
		home = os.Getenv("HOME")
	} else {
		home = "C:\\"
	}

	buf := `{
	"SomePassword": "$ENV$Test2"
}
`
	os.Mkdir(home+"/local", 0755)
	ioutil.WriteFile(home+"/local/b.json", []byte(buf), 0644)

	for ii, test := range tests {
		os.Setenv(test.SetEnvName, test.SetEnvVal)
		ReadFile(test.FileName, &gCfg)
		if gCfg.SomePassword != test.ExpectedPassword {
			t.Errorf("Test %d, expected %s got %s\n", ii, test.ExpectedPassword, gCfg.SomePassword)
		}
		if gCfg.CheckDefault != "dflt-3" {
			t.Errorf("Test %d, expected %s got %s\n", ii, "dflt-3", gCfg.SomePassword)
		}
	}

}
