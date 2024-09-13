package common

import (
	"context"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/asdine/storm/v3"
)

type DBService interface {
	GetDB(options DBOptions) storm.Node
}

var defaultDB *storm.DB
var defaultStormNode storm.Node

func SetDefaultStormNode(n storm.Node) {
	defaultStormNode = n
}

type DefaultDBService struct {
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (d *DefaultDBService) GetDB(options DBOptions) storm.Node {
	if options.DB != nil {
		return NewStormDBNodeWrapper(options.DB)
	}
	if defaultStormNode != nil {
		return NewStormDBNodeWrapper(defaultStormNode)
	}

	dbFile := path.Join("./", "kubeone.db")

	exists, _ := PathExists(dbFile)
	if !exists {
		endpoint := os.Getenv("OSS_ENDPOINT")
		endpoint = strings.Replace(endpoint, "\\", "", -1)
		ossclient, err := InitOssClient(endpoint)
		if err != nil {
			panic(err)
		}
		data, err := ossclient.getObject("oss://xx/kubeone.db")
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(dbFile, data, 0644)
		if err != nil {
			panic(err)
		}
	}
	if defaultDB == nil {
		db, err := storm.Open(dbFile)
		if err != nil {
			panic(err)
		}
		defaultDB = db
	}
	return NewStormDBNodeWrapper(defaultDB)
}

type DBOptions struct {
	DB storm.Node
}

type StormDBNodeWrapper struct {
	storm.Node
}

func NewStormDBNodeWrapper(node storm.Node) *StormDBNodeWrapper {
	return &StormDBNodeWrapper{Node: node}
}

func (s *StormDBNodeWrapper) Save(data interface{}) error {
	err := s.Node.Save(data)
	if err != nil {
		log.Printf("s.Node.Save [%v] failed.err:%v ", data, err)
	}

	if err == nil {
		//defaultDB.Close()
		//defaultDB = nil
		// wait for db file ready
		time.Sleep(3 * time.Second)
		dbFile := path.Join("./", "kubeone.db")
		endpoint := os.Getenv("OSS_ENDPOINT")
		endpoint = strings.Replace(endpoint, "\\", "", -1)
		ossclient, err := InitOssClient(endpoint)
		if err != nil {
			log.Println("InitOssClient  ", err)
			return err
		}

		dbFileData, err := os.ReadFile(dbFile)
		if err != nil {
			log.Println("ReadFile  ", err)
			return err
		}
		err = ossclient.PutObject("oss://xx/kubeone.db", dbFileData, 3, "", context.TODO())
		if err != nil {
			log.Println("PutObject  ", err)
		}
	}

	return err
}
