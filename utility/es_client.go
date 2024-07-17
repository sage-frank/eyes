package utility

import (
	"fmt"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
	"os"
)

func EsInit() (esClient *elasticsearch.Client, err error) {
	esHost := viper.GetString("es.host")
	esPort := viper.GetString("es.port")
	esUser := viper.GetString("es.user")
	esPassWord := viper.GetString("es.password")
	addr := fmt.Sprintf("%s:%s", esHost, esPort)

	cfg := elasticsearch.Config{
		Addresses: []string{addr},
		Username:  esUser,
		Password:  esPassWord,
		Logger: &elastictransport.TextLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: false,
		},
	}

	esClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch.NewClient: %w", err)
	}

	_, err = esClient.Info()
	if err != nil {
		return nil, fmt.Errorf("esClient.Info: %w", err)
	}
	return esClient, nil
}
