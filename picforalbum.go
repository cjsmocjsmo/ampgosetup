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
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "io/ioutil"
	"log"
	// "os"
	// "path/filepath"
	// "strconv"
	// "strings"
	// "time"
)

func GetPicForAlbum(alb string) map[string]string {
	// startLibLogging()
	log.Printf("GetPicForAlbum: %s this is alb", alb)
	filter := bson.M{"album": alb}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "GetPicForAlbum: MongoDB connection has failed")
	collection := client.Database("maindb").Collection("maindb")
	var albuminfo Tagmap
	err = collection.FindOne(context.Background(), filter).Decode(&albuminfo)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("GetPicForAlbum: %s this is album", alb)
	log.Printf("GetPicForAlbum: %s this is AlbumID", albuminfo.AlbumID)
	log.Printf("GetPicForAlbum: %s this is PicHttpAddr", albuminfo.PicHttpAddr)

	var albinfo map[string]string = make(map[string]string)
	albinfo["Album"] = alb
	albinfo["AlbumID"] = albuminfo.AlbumID
	albinfo["PicPath"] = albuminfo.PicHttpAddr
	AmpgoInsertOne("tempdb2", "artidpic", albinfo)
	fmt.Println(albinfo)
	fmt.Println(albinfo)
	return albinfo
}
