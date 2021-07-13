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
	// "os"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"strconv"
)



//GDistArtist2 exported
func GDistArtist2() (DArtAll []map[string]string) {
	sesC := DBcon()
	defer sesC.Close()
	MAINc := sesC.DB("maindb").C("maindb")
	var DArtist []string
	MAINc.Find(nil).Distinct("artist", &DArtist)
	// fmt.Printf("\n\n\n THIS IS DARTIST %s \n\n\n", DArtist)
	for _, art := range DArtist {
		fmt.Printf("\n\n\n THIS IS ART %s \n\n\n", art)
		MAINc := sesC.DB("maindb").C("maindb")
		b1 := bson.M{"artist": art}
		var DArtA map[string]string = make(map[string]string)
		MAINc.Find(b1).One(&DArtA)
		DArtAll = append(DArtAll, DArtA)
	}
	return
}

//GArtInfo2 exported
func GArtInfo2(Dart map[string]string) (ArtInfo2 map[string]string) {
	sesC := DBcon()
	defer sesC.Close()
	MAINc := sesC.DB("maindb").C("maindb")
	b1 := bson.M{"artist": Dart["artist"]}
	MAINc.Find(b1).One(&ArtInfo2)
	return
}

//Ap2 exported
type Ap2 struct {
	Albumz []string
}

//ArtPipeline exported
func ArtPipeline(dart map[string]string) (AP2 []Ap2) {
	sesC := DBcon()
	defer sesC.Close()
	AMPc := sesC.DB("maindb").C("maindb")
	pipeline2 := AMPc.Pipe([]bson.M{
		{"$match": bson.M{"artist": dart["artist"]}},
		{"$group": bson.M{"_id": "album", "albumz": bson.M{"$addToSet": "$album"}}},
		{"$project": bson.M{"albumz": 1}},
	}).Iter()
	err := pipeline2.All(&AP2)
	if err != nil {
		fmt.Printf("\n this is Agg artist pipeline2 fucked up %v %v %T", dart, err, dart)
	}
	return
}

//AddAlbumID exported
func AddAlbumID(PL2 []Ap2) []map[string]string {
	sesC := DBcon()
	defer sesC.Close()
	AMP2c := sesC.DB("maindb").C("maindb")
	var AAID []map[string]string
	for _, boo := range PL2 {
		for _, boo2 := range boo.Albumz {
			var AAid map[string]string = make(map[string]string)
			b1 := bson.M{"album": boo2}
			b2 := bson.M{}
			AMP2c.Find(b1).Select(b2).One(&AAid)
			AAID = append(AAID, AAid)
		}
	}
	return AAID
}

//ArtVIEW exported
type ArtVIEW struct {
	Artist   string              `bson:"artist"`
	ArtistID string              `bson:"artistID"`
	Albums   []map[string]string `bson:"albums"`
	Page     string              `bson:"page"`
	Idx      string              `bson:"idx"`
}

//InsArtIPipe2 exported
func InsArtIPipe2(AI2 map[string]string, aAID []map[string]string, idxx int) {
	// page, _ := strconv.Atoi(AI2["page"])
	var AV1 ArtVIEW
	AV1.Artist = AI2["artist"]
	AV1.ArtistID = AI2["artistID"]
	AV1.Albums = aAID
	// AV1.Page = strconv.Itoa(page)
	AV1.Page = AI2["page"]
	AV1.Idx = strconv.Itoa(idxx)
	sesC := DBcon()
	defer sesC.Close()
	ARTV3c := sesC.DB("artistview").C("artistviews")
	err := ARTV3c.Insert(&AV1)
	if err != nil {
		fmt.Printf("this is ARTV3c Insert err %v", err)
	}
}

//GAVAll exported
func GAVAll() (Artview []ArtVIEW) {
	sesC := DBcon()
	defer sesC.Close()
	ARTc := sesC.DB("artistview").C("artistviews")
	err := ARTc.Find(nil).All(&Artview)
	CheckError(err, "GAVAll: artview fucked up")
	return
}

// ArtistOffset exported
func ArtistOffset() {
	ArtVieW := GAVAll()
	var page int = 1
	for i, art := range ArtVieW {
		var BOO ArtVIEW
		if i < Offset {
			BOO.Page = strconv.Itoa(page)
		} else if i % Offset == 0 {
			page++
			BOO.Page = strconv.Itoa(page)
		} else {
			BOO.Page = strconv.Itoa(page)
		}
		BOO.Artist = art.Artist
		BOO.ArtistID = art.ArtistID
		BOO.Albums = art.Albums
		BOO.Idx = art.Idx
		// sesC := DBcon()
		// defer sesC.Close()
		// ARTc := sesC.DB("artistview").C("artistviews")
		// ARTc.Update(bson.M{"ArtistID": art.ArtistID, "Page": art.Page}, BOO)
		ses := DBcon()
		defer ses.Close()
		tagz := ses.DB("tempdb2").C("artistviewoffset")
		tagz.Insert(TAGmap)

	}
}
