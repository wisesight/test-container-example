package post_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	mongodbTestContainer "github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/wisesight/test-container-example/internal/services/post"
	"github.com/wisesight/test-container-example/internal/testutils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MONGO_VERSION = "6"
	DB            = "testdb"
	COLLECTION    = "posts"
)

var _ = Describe("Mongo", Ordered, func() {
	var mongoDBContainerURI string
	var mongodbContainer *mongodbTestContainer.MongoDBContainer
	var mongoClient *mongo.Client
	BeforeAll(func() {
		var err error
		mongodbContainer, mongoDBContainerURI, err = testutils.SpawnMongoDBContainer(context.Background(), MONGO_VERSION)
		if err != nil {
			panic(err)
		}
		mongoClient, err = testutils.NewMongoDBClient(mongoDBContainerURI)
		if err != nil {
			panic(err)
		}
	})

	BeforeEach(func() {
		if err := testutils.EmptyCollection(context.Background(), mongoClient.Database(DB).Collection(COLLECTION)); err != nil {
			panic(err)
		}
	})

	Describe("New", func() {
		It("should add a new post to the database", func() {
			repo := post.NewMongoRepository(mongoClient, DB)
			now := time.Now().UTC().Truncate(time.Second)
			newPost := &post.Post{
				ID:       primitive.NewObjectID().Hex(),
				Title:    "Test Post",
				Body:     "This is a test post",
				DateTime: now,
			}

			Expect(repo.New(context.Background(), newPost)).To(Succeed())

			resultFromDB, err := testutils.FindOneByID(context.Background(), mongoClient.Database(DB).Collection(COLLECTION), newPost.ID)
			Expect(err).ToNot(HaveOccurred())
			resultFromDB["_id"] = resultFromDB["_id"].(primitive.ObjectID).Hex()
			resultFromDB["date_time"] = resultFromDB["date_time"].(primitive.DateTime).Time().UTC().Truncate(time.Second)
			Expect(resultFromDB["title"]).To(Equal(newPost.Title))
			Expect(resultFromDB["body"]).To(Equal(newPost.Body))
			Expect(resultFromDB["date_time"]).To(Equal(newPost.DateTime))
			Expect(resultFromDB["_id"]).To(Equal(newPost.ID))
		})

		It("should got error when adding a new post with invalid ID", func() {
			repo := post.NewMongoRepository(mongoClient, DB)
			now := time.Now().UTC().Truncate(time.Second)
			newPost := &post.Post{
				ID:       "invalid id",
				Title:    "Test Post",
				Body:     "This is a test post",
				DateTime: now,
			}

			Expect(repo.New(context.Background(), newPost)).To(HaveOccurred())
		})
	})

	Describe("FindByID", func() {
		var repo *post.MongoRepository
		var newPost map[string]interface{}

		BeforeEach(func() {
			repo = post.NewMongoRepository(mongoClient, DB)
			now := time.Now().UTC().Truncate(time.Second)
			// Using insertOne to insert a post to the database
			newPost = map[string]interface{}{
				"_id":       primitive.NewObjectID(),
				"title":     "Test Post",
				"body":      "This is a test post",
				"date_time": now,
			}
			Expect(testutils.InsertOne(context.Background(), mongoClient.Database(DB).Collection(COLLECTION), newPost)).To(Succeed())
		})
		It("should find a post by ID", func() {
			resultFromDB, err := repo.FindByID(context.Background(), newPost["_id"].(primitive.ObjectID).Hex())
			Expect(err).ToNot(HaveOccurred())
			Expect(resultFromDB.ID).To(Equal(newPost["_id"].(primitive.ObjectID).Hex()))
			Expect(resultFromDB.Title).To(Equal(newPost["title"].(string)))
			Expect(resultFromDB.Body).To(Equal(newPost["body"].(string)))
			Expect(resultFromDB.DateTime).To(Equal(newPost["date_time"].(time.Time).UTC().Truncate(time.Second)))
		})

		It("should got error when finding a post with invalid ID", func() {
			resultFromDB, err := repo.FindByID(context.Background(), "invalid id")
			Expect(err).To(HaveOccurred())
			Expect(resultFromDB).To(BeNil())
		})
	})

	Describe("FindByDateTime", func() {
		var repo *post.MongoRepository
		var newPostOne map[string]interface{}
		var newPostTwo map[string]interface{}
		var now time.Time

		BeforeEach(func() {
			repo = post.NewMongoRepository(mongoClient, DB)
			now = time.Now().UTC().Truncate(time.Second)
			// Using insertOne to insert a post to the database
			newPostOne = map[string]interface{}{
				"_id":       primitive.NewObjectID(),
				"title":     "Test Post",
				"body":      "This is a test post",
				"date_time": now,
			}
			newPostTwo = map[string]interface{}{
				"_id":       primitive.NewObjectID(),
				"title":     "Test Post",
				"body":      "This is a test post",
				"date_time": now,
			}
			Expect(testutils.InsertOne(context.Background(), mongoClient.Database(DB).Collection(COLLECTION), newPostOne)).To(Succeed())
			Expect(testutils.InsertOne(context.Background(), mongoClient.Database(DB).Collection(COLLECTION), newPostTwo)).To(Succeed())
		})
		It("should find posts by date time", func() {
			resultFromDB, err := repo.FindByDateTime(context.Background(), now)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(resultFromDB)).To(Equal(2))

			// Sort the result from DB by ID
			if resultFromDB[0].ID > resultFromDB[1].ID {
				resultFromDB[0], resultFromDB[1] = resultFromDB[1], resultFromDB[0]
			}

			Expect(resultFromDB[0].ID).To(Equal(newPostOne["_id"].(primitive.ObjectID).Hex()))
			Expect(resultFromDB[0].Title).To(Equal(newPostOne["title"].(string)))
			Expect(resultFromDB[0].Body).To(Equal(newPostOne["body"].(string)))
			Expect(resultFromDB[0].DateTime).To(Equal(newPostOne["date_time"].(time.Time).UTC().Truncate(time.Second)))

			Expect(resultFromDB[1].ID).To(Equal(newPostTwo["_id"].(primitive.ObjectID).Hex()))
			Expect(resultFromDB[1].Title).To(Equal(newPostTwo["title"].(string)))
			Expect(resultFromDB[1].Body).To(Equal(newPostTwo["body"].(string)))
			Expect(resultFromDB[1].DateTime).To(Equal(newPostTwo["date_time"].(time.Time).UTC().Truncate(time.Second)))
		})

		It("should find no posts by datetime is invalid", func() {
			resultFromDB, err := repo.FindByDateTime(context.Background(), time.Time{})
			Expect(err).ToNot(HaveOccurred())
			Expect(len(resultFromDB)).To(Equal(0))
		})
	})

	AfterAll(func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
		if err := testutils.TerminateMongoDBContainer(context.Background(), mongodbContainer); err != nil {
			panic(err)
		}
	})
})
