package node

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/bf_playground_backend/config"
	"github.com/tsinghua-cel/bf_playground_backend/models/dbmodel"
	"github.com/tsinghua-cel/bf_playground_backend/openapi"
)

type Node struct {
	api *openapi.OpenAPI
}

func NewNode(conf *config.Config) (*Node, error) {
	n := new(Node)
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.Database.User, conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.DBName)
	dbmodel.DbInit(conn)

	api := openapi.NewOpenAPI(conf)
	n.api = api
	return n, nil
}

func (n *Node) Start() error {
	if err := n.api.Run(); err != nil {
		log.WithError(err).Error("start openapi server failed")
		return err
	}

	return nil
}

func (n *Node) Stop() {
}
