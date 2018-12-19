package storage_test

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/oktadeveloper/okta-go-vue-example/pkg/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"os"
	"testing"
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
