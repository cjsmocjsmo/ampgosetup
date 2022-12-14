package ampgosetup

import (
	"context"
	// "crypto/rand"
	// "encoding/hex"
	// "encoding/json"
	"fmt"
	// "github.com/bogem/id3v2"
	// "github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "io/ioutil"
	// "log"
	// "os"
	// "path/filepath"
	// "strconv"
	// "strings"
	"time"
)

func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func Connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

func InsertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {
	collection := client.Database(dataBase).Collection(col)
	result, err := collection.InsertOne(ctx, doc)
	return result, err
}

func Query(client *mongo.Client, ctx context.Context, dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

func AmpgoFindOne(db string, coll string, filtertype string, filterstring string) map[string]string {
	filter := bson.M{filtertype: filterstring}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "AmpgoFindOne: MongoDB connection has failed")
	collection := client.Database(db).Collection(coll)
	var results map[string]string = make(map[string]string)
	err = collection.FindOne(context.Background(), filter).Decode(&results)
	if err != nil {
		fmt.Println("AmpgoFindOne: find one has fucked up")
		fmt.Println(err)
	}
	return results
}

func AmpgoFind(dbb string, collb string, filtertype string, filterstring string) []map[string]string {
	filter := bson.M{}
	if filtertype != "None" && filterstring != "None" {
		filter = bson.M{filtertype: filterstring}
	}
	// filter := bson.D{{filtertype, filterstring}}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "AmpgoFind: MongoDB connection has failed")
	coll := client.Database(dbb).Collection(collb)
	cur, err := coll.Find(context.TODO(), filter)
	CheckError(err, "AmpgoFind: ArtPipeline find has failed")
	var results []map[string]string //all albums for artist to include double entries
	if err = cur.All(context.TODO(), &results); err != nil {
		fmt.Println("AmpgoFind: cur.All has fucked up")
		fmt.Println(err)
	}
	return results
}

func AmpgoInsertOne(db string, coll string, ablob map[string]string) {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	CheckError(err, "AmpgoInsertOne: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, db, coll, ablob)
	CheckError(err2, "AmpgoInsertOne has failed")
}

func InsertMP3Json(db string, coll string, ablob JsonMP3) {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	CheckError(err, "InsertMP3Json: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, db, coll, ablob)
	CheckError(err2, "InsertMP3Json has failed")
}

func InsertJPGJson(db string, coll string, ablob JsonJPG) {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	CheckError(err, "InsertJPGJson: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, db, coll, ablob)
	CheckError(err2, "InsertJPGJson has failed")
}

func InsertPagesJson(db string, coll string, ablob JsonPage) {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	CheckError(err, "InsertPagesJson: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, db, coll, ablob)
	CheckError(err2, "InsertPagesJson has failed")
}

func InsertArtistIDS(db string, coll string, ablob ArtistIDS) {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	CheckError(err, "InsertArtistIDS: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, db, coll, ablob)
	CheckError(err2, "InsertArtistIDS has failed")
}

func InsertAlbumIDS(db string, coll string, ablob AlbumIDS) {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	CheckError(err, "InsertAlbumIDS: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, db, coll, ablob)
	CheckError(err2, "InsertAlbumIDS has failed")
}