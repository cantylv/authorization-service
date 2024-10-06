package clients

import (
	pClient "github.com/cantylv/authorization-service/client"
	aClient "github.com/cantylv/authorization-service/microservices/archive_manager/client"
	"github.com/spf13/viper"
)

type Cluster struct {
	ArchiveClient   *aClient.Client
	PrivelegeClient *pClient.Client
}

func InitCluster() *Cluster {
	privelegeClient := pClient.NewClient(&pClient.ClientOpts{
		Host:   viper.GetString("microservice_privelege.host"),
		Port:   viper.GetInt("microservice_privelege.port"),
		UseSsl: false,
	})
	privelegeClient.CheckConnection()

	archiveClient := aClient.NewClient(&aClient.ClientOpts{
		Host:   viper.GetString("microservice_archive.host"),
		Port:   viper.GetInt("microservice_archive.port"),
		UseSsl: false,
	})
	archiveClient.CheckConnection()
	return &Cluster{
		ArchiveClient:   archiveClient,
		PrivelegeClient: privelegeClient,
	}
}
