package clients

import (
	pClient "github.com/cantylv/authorization-service/client"
	aClient "github.com/cantylv/authorization-service/microservices/archive_manager/client"
)

type Cluster struct {
	ArchiveClient   *aClient.Client
	PrivelegeClient *pClient.Client
}

func InitCluster() *Cluster {
	privelegeClient := pClient.NewClient(&pClient.ClientOpts{
		Host:   "localhost",
		Port:   8010,
		UseSsl: false,
	})
	privelegeClient.CheckConnection()

	archiveClient := aClient.NewClient(&aClient.ClientOpts{
		Host:   "localhost",
		Port:   8011,
		UseSsl: false,
	})
	archiveClient.CheckConnection()
	return &Cluster{
		ArchiveClient:   archiveClient,
		PrivelegeClient: privelegeClient,
	}
}
