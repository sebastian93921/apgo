package system

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Config represents the global configuration
type Settings struct {
	CACertificate []byte `json:"ca_cert"`
	CAPrivateKey  []byte `json:"ca_private"`
}

type GlobalSettings struct {
}

var SettingsFileName = "settings.json"

func LoadGlobalSettings(path string) *Settings {

	// if path doesn't exists, create it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}

	var cfg *Settings
	jsonSettingsPath := filepath.Join(path, SettingsFileName)
	jsonSettingsFile, err := os.Open(jsonSettingsPath)
	defer jsonSettingsFile.Close()
	if err != nil {
		cfg = initGlobalSettings()
		// Generate a settings file if needed...
		// saveGlobalSettings(cfg, jsonSettingsPath)
		return cfg
	}

	byteValue, err := io.ReadAll(jsonSettingsFile)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		cfg = initGlobalSettings()
		// saveGlobalSettings(cfg, jsonSettingsPath)
	}

	return cfg

}

func initGlobalSettings() *Settings {
	// generate a new CA
	rawPvt, rawCA, _ := CreateCA()
	pemPvt := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: rawPvt})
	pemCA := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rawCA})
	cfg := &Settings{
		CACertificate: pemCA,
		CAPrivateKey:  pemPvt,
	}
	return cfg
}

func saveGlobalSettings(cfg *Settings, path string) error {
	jsonSettings, _ := json.MarshalIndent(cfg, "", " ")
	return os.WriteFile(path, jsonSettings, 0700)
}
