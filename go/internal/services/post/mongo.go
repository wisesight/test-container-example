package post

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const COLLECTION = "posts"

type MongoRepository struct {
	coll *mongo.Collection
}

type RawMongoPost struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Title    string             `bson:"title,omitempty"`
	Body     string             `bson:"body,omitempty"`
	DateTime time.Time          `bson:"date_time,omitempty"`
}

func (rp *RawMongoPost) ToPost() *Post {
	return &Post{
		ID:       rp.ID.Hex(),
		Title:    rp.Title,
		Body:     rp.Body,
		DateTime: rp.DateTime,
	}
}

func (rp *RawMongoPost) FromPost(post *Post) error {
	objID, err := primitive.ObjectIDFromHex(post.ID)
	if err != nil {
		return err
	}
	rp.ID = objID
	rp.Title = post.Title
	rp.Body = post.Body
	rp.DateTime = post.DateTime
	return nil
}

func (r *MongoRepository) New(ctx context.Context, post *Post) error {
	rawMongoPost := &RawMongoPost{}
	err := rawMongoPost.FromPost(post)
	if err != nil {
		return err
	}
	_, err = r.coll.InsertOne(ctx, rawMongoPost)
	return err
}

func (r *MongoRepository) FindByID(ctx context.Context, id string) (*Post, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	rawMongoPost := &RawMongoPost{}
	result := r.coll.FindOne(ctx, bson.M{"_id": objectID})
	err = result.Decode(rawMongoPost)
	return rawMongoPost.ToPost(), err
}

func (r *MongoRepository) FindByDateTime(ctx context.Context, dateTime time.Time) ([]*Post, error) {
	result, err := r.coll.Find(ctx, bson.M{"date_time": dateTime})
	if err != nil {
		return nil, err
	}
	var rawMongoPosts []*RawMongoPost
	err = result.All(ctx, &rawMongoPosts)
	if err != nil {
		return nil, err
	}
	var posts []*Post
	for _, rawMongoPost := range rawMongoPosts {
		posts = append(posts, rawMongoPost.ToPost())
	}
	return posts, nil
}

func NewMongoRepository(client *mongo.Client, databaseName string) *MongoRepository {
	return &MongoRepository{
		coll: client.Database(databaseName).Collection(COLLECTION),
	}
}
