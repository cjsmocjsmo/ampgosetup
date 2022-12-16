package ampgosetup

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	// "github.com/bogem/id3v2"
	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/bson"

	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	// "log"
	"os"
	"path/filepath"
	"strconv"
	// "strings"
	// "time"
)

// Tagmap exported
type Tagmap struct {
	Dirpath     string `bson:"dirpath"`
	Filename    string `bson:"filename"`
	Extension   string `bson:"extension"`
	FileID      string `bson:"fileID"`
	Filesize    string `bson:"filesize"`
	Artist      string `bson:"artist"`
	ArtistID    string `bson:"artistID"`
	Album       string `bson:"album"`
	AlbumID     string `bson:"albumID"`
	Title       string `bson:"title"`
	Genre       string `bson:"genre"`
	TitlePage   string `bson:"titlepage"`
	PicID       string `bson:"picID"`
	PicDB       string `bson:"picDB"`
	PicPath     string `bson:"picPath"`
	PicHttpAddr string `bson:"picHttpAddr"`
	Idx         string `bson:"idx"`
	HttpAddr    string `bson:"httpaddr"`
	Duration    string `bson:"duration"`

	ArtStart string `bson:"artstart"`
	AlbStart string `bson:"albstart"`
	TitStart string `bson:"titstart"`
	Howl     string `bson:"howl"`
}

type ArtVieW2 struct {
	Artist   string              `bson:"artist"`
	ArtistID string              `bson:"artistID"`
	Albums   []map[string]string `bson:"albums"`
	AlbCount string              `bson:"albcount"`
	Page     string              `bson:"page"`
	Index    string              `bson:"idx"`
}

type AlbVieW2 struct {
	Artist      string              `bson:"artist"`
	ArtistID    string              `bson:"artistID"`
	Album       string              `bson:"album"`
	AlbumID     string              `bson:"albumID"`
	Songs       []map[string]string `bson:"songs"`
	AlbumPage   string              `bson:"albumpage"`
	NumSongs    string              `bson:"numsongs"`
	PicPath     string              `bson:"picPath"`
	Idx         string              `bson:"idx"`
	PicHttpAddr string              `bson:"picHttpAddr"`
}

type Imageinfomap struct {
	Dirpath       string `bson:"dirpath"`
	Filename      string `bson:"filename"`
	Imagesize     string `bson:"imagesize"`
	ImageHttpAddr string `bson:"imageHttpAddr"`
	Index         string `bson:"index"`
	IType         string `bson:"itype"`
	Page          string `bson:"page"`
}

type Fjpg struct {
	exists bool
	path   string
}

type randDb struct {
	PlayListName  string              `bson:"playlistname"`
	PlayListID    string              `bson:"playlistID"`
	PlayListCount string              `bson:"playlistcount"`
	Playlist      []map[string]string `bson:"playlist"`
}

//////////////////////////////////////////////////////////////////////////

func getFileInfo(apath string) (filename string, size string) {
	ltn, err := os.Open(apath)
	CheckError(err, "getFileInfo: file open has fucked up")
	defer ltn.Close()
	ltnInfo, _ := ltn.Stat()
	filename = ltnInfo.Name()
	size = strconv.FormatInt(ltnInfo.Size(), 10)
	return
}

func UUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[8] = 0x80
	uuid[4] = 0x40
	boo := hex.EncodeToString(uuid)
	return boo, nil
}

func resizeImage(infile string, outfile string) string {
	pic, err := imaging.Open(infile)
	if err != nil {
		return os.Getenv("AMPGO_NO_ART_PIC_PATH")
	}
	sjImage := imaging.Resize(pic, 200, 0, imaging.Lanczos)
	err = imaging.Save(sjImage, outfile)
	CheckError(err, "resizeImage: image save has fucked up")
	return outfile
}

func TaGmap(apath string, apage int, idx int) (TaGmaP Tagmap) {
	artist, album, title, genre, picpath := DumpArtToFile(apath)
	if artist != "None" && album != "None" && title != "None" {
		fmt.Println(apath)
		page := strconv.Itoa(apage)
		index := strconv.Itoa(idx)
		uuid, _ := UUID()
		pichttpaddr := os.Getenv("AMPGO_SERVER_ADDRESS") + ":" + os.Getenv("AMPGO_SERVER_PORT") + picpath[5:]
		fname, size := getFileInfo(apath)
		httpaddr := os.Getenv("AMPGO_SERVER_ADDRESS") + ":" + os.Getenv("AMPGO_SERVER_PORT") + apath[5:]
		TaGmaP.Dirpath = filepath.Dir(apath)
		TaGmaP.Filename = fname
		TaGmaP.Extension = filepath.Ext(apath)
		TaGmaP.FileID = uuid
		TaGmaP.Filesize = size
		TaGmaP.Artist = artist
		TaGmaP.ArtistID = "None"
		TaGmaP.Album = album
		TaGmaP.AlbumID = "None"
		TaGmaP.Title = title
		TaGmaP.Genre = genre
		TaGmaP.TitlePage = page
		TaGmaP.PicID = uuid
		TaGmaP.PicDB = "None"
		TaGmaP.PicPath = picpath
		TaGmaP.PicHttpAddr = pichttpaddr
		TaGmaP.Idx = index
		TaGmaP.HttpAddr = httpaddr
		TaGmaP.Duration = "None"
		TaGmaP.ArtStart = "None"
		TaGmaP.AlbStart = "None"
		TaGmaP.TitStart = "None"
		TaGmaP.Howl = ""
		client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
		CheckError(err, "TaGmap: Connections has failed")
		defer Close(client, ctx, cancel)
		_, err2 := InsertOne(client, ctx, "tempdb1", "meta1", &TaGmaP)
		CheckError(err2, "TaGmap: Tempdb1 insertion has failed")
		return
	} else {
		os.Remove(apath)
	}
	return
}

/////////////////////////////////////////////////////////////////////////////////////////////

// func GetAllObjects() (Main2SL []map[string]string) {

func GetAllObjects() (Main2SL []JsonMP3) {
	filter := bson.D{}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgodb")
	defer Close(client, ctx, cancel)
	CheckError(err, "GetAllObjects: MongoDB connection has failed")
	collection := client.Database("tempdb1").Collection("meta1")
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		fmt.Println(err)
	}
	if err = cur.All(context.Background(), &Main2SL); err != nil {
		fmt.Println("GetAllObjects: cur.All has failed")
		fmt.Println(err)
	}
	return
}

func Unique(arr []string) []string {
	occured := map[string]bool{}
	result := []string{}
	for e := range arr {
		if !occured[arr[e]] {
			occured[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result
}

func create_just_albumID_list(alist []map[string]string) (just_albumID_list []string) {
	for _, albID := range alist {
		just_albumID_list = append(just_albumID_list, albID["albumID"])
	}
	return
}

func get_albums_for_artist(fullalblist []map[string]string) (final_alblist []map[string]string) {
	just_albumID_list := create_just_albumID_list(fullalblist)
	//remove double albumid entries
	unique_items := Unique(just_albumID_list)
	for _, uitem := range unique_items {
		albINFO := AmpgoFindOne("maindb", "maindb", "albumID", uitem)
		final_alblist = append(final_alblist, albINFO)
	}
	return
}

func get_songs_for_album(fullsonglist []map[string]string) (final_songlist []map[string]string) {
	//a list of just albumid's
	var just_songID_list []string
	for _, song := range fullsonglist {
		just_songID_list = append(just_songID_list, song["fileID"])
	}
	//remove double songID entries
	unique_items := Unique(just_songID_list)
	for _, uitem := range unique_items {
		songINFO := AmpgoFindOne("maindb", "maindb", "fileID", uitem)
		final_songlist = append(final_songlist, songINFO)
	}
	return final_songlist
}

func CreateRandomPlaylistDB() string {
	var ranDBInfo randDb
	var emptylist []map[string]string
	var emptyitem map[string]string = map[string]string{"None": "No Songs Found"}
	emptylist = append(emptylist, emptyitem)
	uuid, _ := UUID()
	ranDBInfo.PlayListName = "EmptyRandomPlaylist"
	ranDBInfo.PlayListID = uuid
	ranDBInfo.PlayListCount = "0"
	ranDBInfo.Playlist = emptylist
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "CreateRandomPlaylistDB: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "randplaylists", "randplaylists", ranDBInfo)
	CheckError(err2, "CreateRandomPlaylistDB: randplaylists insertion has failed")
	return "Created"
}

func ReadDurationFile(apath string) map[string]string {
	data, err := ioutil.ReadFile(apath)
	CheckError(err, "ReadDurationFile: mp3info read has failed")
	var mp3info map[string]string
	err2 := json.Unmarshal(data, &mp3info)
	CheckError(err2, "ReadDurationFile: json unmarshal has failed")
	return mp3info
}

func InsertDurationInfo(apath string) string {
	mp3 := ReadDurationFile(apath)
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "InsertDurationInfo: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "durdb", "durdb", mp3)
	CheckError(err2, "InsertDurationInfo: durdb insertion has failed")
	return "durdb Created"
}

func CreateCurrentPlayListNameDB() string {
	var curPlayListName map[string]string = map[string]string{"record": "1", "curplaylistname": "None", "curplaylistID": "None"}
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "InsertDurationInfo: Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "curplaylistname", "curplaylistname", &curPlayListName)
	CheckError(err2, "InsertDurationInfo: curplaylistname insertion has failed")
	return "curplaylistname Created"
}
