package storage_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/klebervirgilio/vue-crud-app-with-golang/pkg/core"
	. "github.com/klebervirgilio/vue-crud-app-with-golang/pkg/storage"
)

var _ = Describe("Mongo", func() {
	var repo core.Repository
	os.Setenv("MONGO_URL", "mongodb://mongo_user:mongo_secret@0.0.0.0:27017/kudos_test")
	repo = NewMongoRepository()

	Describe(".Create", func() {
		It("inserts a kudo into the database", func() {
			// Arrange
			var err error
			err = repo.Create(&core.Kudo{RepoID: "some-id"})
			if err != nil {
				Fail(err.Error())
			}

			// Act
			kudo, err := repo.Find("some-id")
			if err != nil {
				Fail(err.Error())
			}

			// Assert
			Expect(kudo.RepoID).To(Equal("some-id"))
		})

		Context("when document already exists", func() {
			var kudo *core.Kudo
			var err error

			BeforeEach(func() {
				err = repo.Create(&core.Kudo{RepoID: "some-id"})
				err = repo.Create(&core.Kudo{RepoID: "some-id", Language: "golang"})
				kudo, err = repo.Find("some-id")
				if err != nil {
					Fail(err.Error())
				}
			})

			It("updates kudo", func() {
				Expect(kudo.Language).To(Equal("golang"))
			})

			It("does not create duplicated kudos", func() {
				Expect(repo.Count()).To(Equal(1))
			})
		})
	})

	Describe(".Count", func() {
		It("counts kudos", func() {
			// Arrange
			var err error
			err = repo.Create([]*core.Kudo{
				{RepoID: "some-id"},
				{RepoID: "some-other-id"},
				{RepoID: "another-id"},
			}...)

			// Act
			count, err := repo.Count()
			if err != nil {
				Fail(err.Error())
			}

			// Assert
			Expect(count).To(Equal(3))
		})
	})

	Describe(".Find", func() {
		It("finds a kudo", func() {
			// Arrange
			var err error
			err = repo.Create(&core.Kudo{RepoID: "some-id"})
			if err != nil {
				Fail(err.Error())
			}

			// Act
			_, err = repo.Find("some-id")

			// Assert
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe(".FindAll", func() {
		It("finds all kudos which match with the selector", func() {
			// Arrange
			var err error
			err = repo.Create([]*core.Kudo{
				{RepoID: "some-id", Language: "golang"},
				{RepoID: "some-other-id", Language: "golang"},
				{RepoID: "another-id", Language: "erlang"},
			}...)
			if err != nil {
				Fail(err.Error())
			}

			// Act
			kudos, err := repo.FindAll(map[string]interface{}{"language": "golang"})

			// Assert
			Expect(kudos).To(Equal([]*core.Kudo{
				{RepoID: "some-id", Language: "golang"},
				{RepoID: "some-other-id", Language: "golang"},
			}))
		})
	})
	Describe(".Delete", func() {
		It("deletes a kudo", func() {
			// Arrange
			var err error
			err = repo.Create([]*core.Kudo{
				{RepoID: "some-id", Language: "golang"},
				{RepoID: "some-other-id", Language: "golang"},
				{RepoID: "another-id", Language: "erlang"},
			}...)
			if err != nil {
				Fail(err.Error())
			}

			// Act
			repo.Delete(&core.Kudo{RepoID: "some-id"})

			// Assert
			kudos, err := repo.FindAll(map[string]interface{}{})
			Expect(kudos).To(Equal([]*core.Kudo{
				{RepoID: "some-other-id", Language: "golang"},
				{RepoID: "another-id", Language: "erlang"},
			}))
		})
	})

	Describe(".Update", func() {
		It("updates a kudo", func() {
			// Arrange
			var err error
			err = repo.Create(&core.Kudo{RepoID: "some-id"})
			if err != nil {
				Fail(err.Error())
			}

			// Act
			err = repo.Update(&core.Kudo{RepoID: "some-id", Language: "golang"})

			// Assert
			kudo, err := repo.Find("some-id")
			if err != nil {
				Fail(err.Error())
			}
			Expect(kudo.Language).To(Equal("golang"))
		})
	})
})
