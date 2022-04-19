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
	// "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "io/ioutil"
	// "log"
	// "os"
	// "path/filepath"
	// "strconv"
	// "strings"
	"time"
)

func AmpgoDistinct(db string, coll string, fieldd string) []string {
	filter := bson.D{}
	opts := options.Distinct().SetMaxTime(2 * time.Second)
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "AmpgoDistinct: MongoDB connection has failed")
	collection := client.Database(db).Collection(coll)
	DD1, err2 := collection.Distinct(context.TODO(), fieldd, filter, opts)
	CheckError(err2, "AmpgoDistinct: MongoDB distinct album has failed")
	var DAlbum1 []string
	for _, DD := range DD1 {
		zoo := fmt.Sprintf("%s", DD)
		DAlbum1 = append(DAlbum1, zoo)
	}
	return DAlbum1
}

func InsAlbumID(alb string) {
	uuid, _ := UUID()
	Albid := map[string]string{"album": alb, "albumID": uuid}
	AmpgoInsertOne("tempdb2", "albumid", Albid)
}

func InsArtistID(art string) {
	uuid, _ := UUID()
	Artid := map[string]string{"artist": art, "artistID": uuid}
	AmpgoInsertOne("tempdb2", "artistid", Artid)
}
