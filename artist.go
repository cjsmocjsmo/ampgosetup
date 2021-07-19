///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////
// LICENSE: GNU General Public License, version 2 (GPLv2)
// Copyright 2016, Charlie J. Smotherman
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License v2
// as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

package ampgosetup

import (
	"os"
	"fmt"
	"log"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "strconv"
)



// //GDistArtist2 exported
func GDistArtist2() (DArtAll []map[string]string) {
	filter := bson.D{}
	opts := options.Distinct().SetMaxTime(2 * time.Second)
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	collection := client.Database("maindb").Collection("maindb")
	DA1, err2 := collection.Distinct(context.TODO(), "artist", filter, opts)
	CheckError(err2, "MongoDB distinct album has failed")
	var DArtist []string
	for _, DA := range DA1 {
		zoo := fmt.Sprintf("%s", DA)
		DArtist = append(DArtist, zoo)
	}


	// fmt.Printf("\n\n\n THIS IS DARTIST %s \n\n\n", DArtist)
	for _, art := range DArtist {
		fmt.Printf("\n\n\n THIS IS ART %s \n\n\n", art)
		filter := bson.M{"artist": art}
		// opts := options.Distinct().SetMaxTime(2 * time.Second)
		client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
		defer Close(client, ctx, cancel)
		CheckError(err, "MongoDB connection has failed")
		collection := client.Database("maindb").Collection("maindb")
		var DArtA map[string]string = make(map[string]string)
		err = collection.FindOne(context.Background(), filter).Decode(&DArtA)
		if err != nil { log.Fatal(err) }
		// fmt.Println("\n\n\n This is DArtA")
		// fmt.Println(DArtA)
		DArtAll = append(DArtAll, DArtA)
	




// 		MAINc := sesC.DB("maindb").C("maindb")
// 		b1 := bson.M{"artist": art}
// 		var DArtA map[string]string = make(map[string]string)
// 		MAINc.Find(b1).One(&DArtA)
// 		DArtAll = append(DArtAll, DArtA)
	}
	return
}

// //GArtInfo2 exported
func GArtInfo2(Dart map[string]string) (string, string) {
	filter := bson.M{"album": Dart["album"]}
	// opts := options.Distinct().SetMaxTime(2 * time.Second)
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	collection := client.Database("maindb").Collection("maindb")
	var ArtInfo2 map[string]string = make(map[string]string)
	err = collection.FindOne(context.Background(), filter).Decode(&ArtInfo2)
	if err != nil { log.Fatal(err) }
	return ArtInfo2["artist"], ArtInfo2["artistID"]


// 	sesC := DBcon()
// 	defer sesC.Close()
// 	MAINc := sesC.DB("maindb").C("maindb")
// 	b1 := bson.M{"artist": Dart["artist"]}
// 	MAINc.Find(b1).One(&ArtInfo2)
	// return
}

// //Ap2 exported
type Ap2 struct {
	Albumz []string
}

// //ArtPipeline exported
func ArtPipeline(dart map[string]string) (AP2 []Ap2) {
	logtxtfile := os.Getenv("AMPGO_LOG_PATH")
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(logtxtfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"artist", dart["artist"]}}}},
		// {"$group": bson.M{"_id": "album", "albumz": bson.M{"$addToSet": "$album"}}},
		{{"$group", bson.M{"_id": "album", "albumz": bson.M{"$addToSet": "$album"}}}},
		{{"$project", bson.D{{"albumz", 1}}}},

	}
	// filter := bson.M{"album": Dart["album"]}
	// opts := options.Distinct().SetMaxTime(2 * time.Second)
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	defer Close(client, ctx, cancel)
	CheckError(err, "MongoDB connection has failed")
	db := client.Database("maindb").Collection("maindb")
	// var AP2 []Ap2
	cur, err := db.Aggregate(context.TODO(), pipeline)
	cur.Decode(&AP2)

	log.Println(AP2)
	log.Printf("%T This is AP2 type", AP2)
	for _, ag := range AP2 {
		fmt.Printf("%v this is ag from AP2", ag)
		fmt.Printf("%T this is ag type", ag)
		log.Printf("%T this is ag type", ag)
		log.Printf("%v this is ag from AP2", ag)
	}
// 	sesC := DBcon()
// 	defer sesC.Close()
// 	AMPc := sesC.DB("maindb").C("maindb")
// 	pipeline2 := AMPc.Pipe([]bson.M{
// 		{"$match": bson.M{"artist": dart["artist"]}},
// 		{"$group": bson.M{"_id": "album", "albumz": bson.M{"$addToSet": "$album"}}},
// 		{"$project": bson.M{"albumz": 1}},
// 	}).Iter()
// 	err := pipeline2.All(&AP2)
// 	if err != nil {
// 		fmt.Printf("\n this is Agg artist pipeline2 fucked up %v %v %T", dart, err, dart)
// 	}
	return
}
// 
// //AddAlbumID exported
// func AddAlbumID(PL2 []Ap2) (AAID []map[string]string) {
	
// 	// client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
// 	// defer Close(client, ctx, cancel)
// 	// CheckError(err, "MongoDB connection has failed")
// 	// collection := client.Database("maindb").Collection("maindb")

// 	for _, aaid := range PL2 {
// 		fmt.Printf("\n %v This aaid from PL2 \n\n\n", aaid)
		
	

// 		// cur, err := collection.Find(context.Background(), filter)
// 		// if err != nil { log.Fatal(err) }
// 		// filter := bson.D{}
// 		// var AAID []map[string]string
// 		// if err = cur.All(context.Background(), &AAID); err != nil {
// 		// 	log.Fatal(err)
// 		// }
// 	}
	
// 	return




// // 	sesC := DBcon()
// // 	defer sesC.Close()
// // 	AMP2c := sesC.DB("maindb").C("maindb")
// // 	var AAID []map[string]string
// // 	for _, boo := range PL2 {
// // 		for _, boo2 := range boo.Albumz {
// // 			var AAid map[string]string = make(map[string]string)
// // 			b1 := bson.M{"album": boo2}
// // 			b2 := bson.M{}
// // 			AMP2c.Find(b1).Select(b2).One(&AAid)
// // 			AAID = append(AAID, AAid)
// // 		}
// // 	}
// 	return AAID
// }

// //ArtVIEW exported
type ArtVIEW struct {
	Artist   string              `bson:"artist"`
	ArtistID string              `bson:"artistID"`
	Albums   []Ap2               `bson:"albums"`
	Page     string              `bson:"page"`
	Idx      string              `bson:"idx"`
}

// //InsArtIPipe2 exported
func InsArtIPipe2(AV1 ArtVIEW) {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "Connections has failed")
	defer Close(client, ctx, cancel)
	_, err2 := InsertOne(client, ctx, "artistview", "artistview", &AV1)
	CheckError(err2, "artistview insertion has failed")
// 	sesC := DBcon()
// 	defer sesC.Close()
// 	ARTV3c := sesC.DB("artistview").C("artistviews")
// 	err := ARTV3c.Insert(&AV1)
// 	if err != nil {
// 		fmt.Printf("this is ARTV3c Insert err %v", err)
// 	}
}

// //GAVAll exported
// func GAVAll() (Artview []ArtVIEW) {
// 	sesC := DBcon()
// 	defer sesC.Close()
// 	ARTc := sesC.DB("artistview").C("artistviews")
// 	err := ARTc.Find(nil).All(&Artview)
// 	CheckError(err, "GAVAll: artview fucked up")
// 	return
// }

// // ArtistOffset exported
// func ArtistOffset() {
// 	ArtVieW := GAVAll()
// 	var page int = 1
// 	for i, art := range ArtVieW {
// 		var BOO ArtVIEW
// 		if i < Offset {
// 			BOO.Page = strconv.Itoa(page)
// 		} else if i % Offset == 0 {
// 			page++
// 			BOO.Page = strconv.Itoa(page)
// 		} else {
// 			BOO.Page = strconv.Itoa(page)
// 		}
// 		BOO.Artist = art.Artist
// 		BOO.ArtistID = art.ArtistID
// 		BOO.Albums = art.Albums
// 		BOO.Idx = art.Idx
// 		// sesC := DBcon()
// 		// defer sesC.Close()
// 		// ARTc := sesC.DB("artistview").C("artistviews")
// 		// ARTc.Update(bson.M{"ArtistID": art.ArtistID, "Page": art.Page}, BOO)
// 		ses := DBcon()
// 		defer ses.Close()
// 		tagz := ses.DB("tempdb2").C("artistviewoffset")
// 		tagz.Insert(BOO)

// 	}
// }
