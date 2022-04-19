package ampgosetup

import (
	// "context"
	// // "crypto/rand"
	// // "encoding/hex"
	// // "encoding/json"
	// "fmt"
	// // "github.com/bogem/id3v2"
	// // "github.com/disintegration/imaging"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "io/ioutil"
	// "log"
	// "os"
	// "path/filepath"
	"strconv"
	// "strings"
	// "time"
)

func InsArtPipeline(AV1 ArtVieW2) {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	CheckError(err, "InsArtPipeline: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "artistview", "artistview", &AV1)
	CheckError(err2, "InsArtPipeline: artistview insertion has failed")
}

func GDistArtist2() (dArtAll []map[string]string) {
	dArtist := AmpgoDistinct("maindb", "maindb", "artist")
	for _, art := range dArtist {
		dArt := AmpgoFindOne("maindb", "maindb", "artist", art)
		dArtAll = append(dArtAll, dArt)
	}
	return dArtAll
}

func ArtPipline(artmap map[string]string, page int, idx int) (MyArView ArtVieW2) {
	dirtyalblist := AmpgoFind("maindb", "maindb", "artistID", artmap["artistID"]) //[]map[string]string
	results2 := get_albums_for_artist(dirtyalblist)
	albc := len(results2)
	albcount := strconv.Itoa(albc)
	MyArView.Artist = artmap["artist"]
	MyArView.ArtistID = artmap["artistID"]
	MyArView.Albums = results2
	MyArView.AlbCount = albcount
	MyArView.Page = strconv.Itoa(page)
	MyArView.Index = strconv.Itoa(idx)
	return
}
