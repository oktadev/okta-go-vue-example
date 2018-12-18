package storage_test

import (
	"log"
	"os"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/klebervirgilio/vue-crud-app-with-golang/pkg/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storage Suite")
}

var _ = BeforeEach(func() {
	mongoURL := os.Getenv("MONGO_URL")
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	session.DB("").C(storage.GetCollectionName()).RemoveAll(bson.M{})
})
