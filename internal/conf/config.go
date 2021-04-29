package config

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/buger/jsonparser"
)

const ()

// Config struct holds all of the runtime configuration
type Config struct {
	BaseConfig BaseConfig
	// TargetURLs []string
}

// BaseConfig The struct for storing K8s Secret
type BaseConfig struct {
	PN_GLOBAL_ROUTER                string
	PN_GLOBAL_PORTAL                string
	PN_GLOBAL_JWT_PASSPHRASE        string
	MY_POD_NAMESPACE                string
	DB_PATH                         string
	MARVIN_OPERATOR_IMAGE_TAG       string
	CONNECTOR_ADDRESS               string
	NEW_PUBLISHED_POLLING_INTERVAL  string
	CHECK_UPGRADE_STATUS_ENABLED    bool
	UPGRADE_STATUS_POLLING_INTERVAL string
	JOB_TIMEOUT                     string
}

// LoadConfig load configuration from mounted K8s Secret
func LoadConfig(configPath string) ([]byte, error) {
	config, err := ioutil.ReadFile(configPath)
	return config, err
}

// InitConfig Set the BaseConfig with loaded contents
func (baseCfg *BaseConfig) InitConfig(configPath string) {
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	baseCfg.PN_GLOBAL_ROUTER, err = jsonparser.GetString(config, "PN_GLOBAL_ROUTER")
	if err != nil {
		log.Fatalf("PN_GLOBAL_ROUTER %v:", err)
	}

	baseCfg.PN_GLOBAL_PORTAL, err = jsonparser.GetString(config, "PN_GLOBAL_PORTAL")
	if err != nil {
		log.Fatalf("PN_GLOBAL_PORTAL %v:", err)
	}

	baseCfg.PN_GLOBAL_JWT_PASSPHRASE, err = jsonparser.GetString(config, "PN_GLOBAL_JWT_PASSPHRASE")
	if err != nil {
		log.Fatalf("PN_GLOBAL_JWT_PASSPHRASE %v:", err)
	}

	baseCfg.MY_POD_NAMESPACE = os.Getenv("MY_POD_NAMESPACE")

	baseCfg.DB_PATH = os.Getenv("DB_PATH")
	if baseCfg.DB_PATH == "" {
		baseCfg.DB_PATH = "./marvin-connector.db"
	}

	baseCfg.CONNECTOR_ADDRESS = os.Getenv("CONNECTOR_ADDRESS")

	baseCfg.NEW_PUBLISHED_POLLING_INTERVAL = os.Getenv("NEW_PUBLISHED_POLLING_INTERVAL")
	if baseCfg.NEW_PUBLISHED_POLLING_INTERVAL == "" {
		baseCfg.NEW_PUBLISHED_POLLING_INTERVAL = "2h"
	}

	b, err := strconv.ParseBool(os.Getenv("CHECK_UPGRADE_STATUS_ENABLED"))
	if err != nil {
		log.Printf("[Warning] parse bool CHECK_UPGRADE_STATUS_ENABLED failed. Not a boolean")
		baseCfg.CHECK_UPGRADE_STATUS_ENABLED = false
	} else {
		baseCfg.CHECK_UPGRADE_STATUS_ENABLED = b
	}

	baseCfg.UPGRADE_STATUS_POLLING_INTERVAL = os.Getenv("UPGRADE_STATUS_POLLING_INTERVAL")
	if baseCfg.UPGRADE_STATUS_POLLING_INTERVAL == "" {
		baseCfg.UPGRADE_STATUS_POLLING_INTERVAL = "1m"
	}

	baseCfg.JOB_TIMEOUT = os.Getenv("JOB_TIMEOUT")
	if baseCfg.JOB_TIMEOUT == "" {
		baseCfg.JOB_TIMEOUT = "2m"
	}

	log.Printf("SECRET.PN_GLOBAL_ROUTER %v", baseCfg.PN_GLOBAL_ROUTER)
	log.Printf("SECRET.PN_GLOBAL_PORTAL %v", baseCfg.PN_GLOBAL_PORTAL)
	log.Printf("ENV.MY_POD_NAMESPACE %v", baseCfg.MY_POD_NAMESPACE)
	log.Printf("ENV.DB_PATH %v", baseCfg.DB_PATH)
	log.Printf("ENV.CONNECTOR_ADDRESS %v", baseCfg.CONNECTOR_ADDRESS)
	log.Printf("ENV.NEW_PUBLISHED_POLLING_INTERVAL %v", baseCfg.NEW_PUBLISHED_POLLING_INTERVAL)
	log.Printf("ENV.UPGRADE_STATUS_POLLING_INTERVAL %v", baseCfg.UPGRADE_STATUS_POLLING_INTERVAL)
	log.Printf("ENV.JOB_TIMEOUT %v", baseCfg.JOB_TIMEOUT)

}

// Init populates the Config struct based on environment runtime configuration
func Init(baseCfg BaseConfig) Config {
	// urls := []string{}

	appConfig := Config{
		baseCfg,
		// urls,
	}

	return appConfig
}
