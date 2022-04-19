package ampgosetup

func AlbumFirst(astring string) string {
	client, ctx, cancel, err := Connect("mongodb://db:27017/ampgo")
	CheckError(err, "AlbumFirst:  Connections has failed")
	defer Close(client, ctx, cancel)

	char := StartsWith(astring)

	switch {
	case char == "A":
		var item map[string]string = map[string]string{"album": astring}
		_, erra := InsertOne(client, ctx, "albumalpha", "A", item)
		CheckError(erra, "AlbumFirst: A insertion has failed")
		return "A Created"

	case char == "B":
		var item map[string]string = map[string]string{"album": astring}
		_, errb := InsertOne(client, ctx, "albumalpha", "B", item)
		CheckError(errb, "AlbumFirst: B insertion has failed")
		return "B Created"

	case char == "C":
		var item map[string]string = map[string]string{"album": astring}
		_, errc := InsertOne(client, ctx, "albumalpha", "C", item)
		CheckError(errc, "AlbumFirst: C insertion has failed")
		return "C Created"

	case char == "D":
		var item map[string]string = map[string]string{"album": astring}
		_, errd := InsertOne(client, ctx, "albumalpha", "D", item)
		CheckError(errd, "AlbumFirst: D insertion has failed")
		return "D Created"

	case char == "E":
		var item map[string]string = map[string]string{"album": astring}
		_, erre := InsertOne(client, ctx, "albumalpha", "E", item)
		CheckError(erre, "AlbumFirst: E insertion has failed")
		return "E Created"

	case char == "F":
		var item map[string]string = map[string]string{"album": astring}
		_, errf := InsertOne(client, ctx, "albumalpha", "F", item)
		CheckError(errf, "AlbumFirst: F insertion has failed")
		return "F Created"

	case char == "G":
		var item map[string]string = map[string]string{"album": astring}
		_, errg := InsertOne(client, ctx, "albumalpha", "G", item)
		CheckError(errg, "AlbumFirst: G insertion has failed")
		return "G Created"

	case char == "H":
		var item map[string]string = map[string]string{"album": astring}
		_, errh := InsertOne(client, ctx, "albumalpha", "H", item)
		CheckError(errh, "AlbumFirst: H insertion has failed")
		return "H Created"

	case char == "I":
		var item map[string]string = map[string]string{"album": astring}
		_, erri := InsertOne(client, ctx, "albumalpha", "I", item)
		CheckError(erri, "AlbumFirst: I insertion has failed")
		return "I Created"

	case char == "J":
		var item map[string]string = map[string]string{"album": astring}
		_, errj := InsertOne(client, ctx, "albumalpha", "J", item)
		CheckError(errj, "AlbumFirst: J insertion has failed")
		return "J Created"

	case char == "K":
		var item map[string]string = map[string]string{"album": astring}
		_, errk := InsertOne(client, ctx, "albumalpha", "K", item)
		CheckError(errk, "AlbumFirst: K insertion has failed")
		return "K Created"

	case char == "L":
		var item map[string]string = map[string]string{"album": astring}
		_, errl := InsertOne(client, ctx, "albumalpha", "L", item)
		CheckError(errl, "AlbumFirst: L insertion has failed")
		return "L Created"

	case char == "M":
		var item map[string]string = map[string]string{"album": astring}
		_, errm := InsertOne(client, ctx, "albumalpha", "M", item)
		CheckError(errm, "AlbumFirst: M insertion has failed")
		return "M Created"

	case char == "N":
		var item map[string]string = map[string]string{"album": astring}
		_, errn := InsertOne(client, ctx, "albumalpha", "N", item)
		CheckError(errn, "AlbumFirst: N insertion has failed")
		return "N Created"

	case char == "O":
		var item map[string]string = map[string]string{"album": astring}
		_, erro := InsertOne(client, ctx, "albumalpha", "O", item)
		CheckError(erro, "AlbumFirst: O insertion has failed")
		return "O Created"

	case char == "P":
		var item map[string]string = map[string]string{"album": astring}
		_, errp := InsertOne(client, ctx, "albumalpha", "P", item)
		CheckError(errp, "AlbumFirst: P insertion has failed")
		return "P Created"

	case char == "Q":
		var item map[string]string = map[string]string{"album": astring}
		_, errq := InsertOne(client, ctx, "albumalpha", "Q", item)
		CheckError(errq, "AlbumFirst: Q insertion has failed")
		return "Q Created"

	case char == "R":
		var item map[string]string = map[string]string{"album": astring}
		_, errr := InsertOne(client, ctx, "albumalpha", "R", item)
		CheckError(errr, "AlbumFirst: R insertion has failed")
		return "R Created"

	case char == "S":
		var item map[string]string = map[string]string{"album": astring}
		_, errs := InsertOne(client, ctx, "albumalpha", "S", item)
		CheckError(errs, "AlbumFirst: S insertion has failed")
		return "S Created"

	case char == "T":
		var item map[string]string = map[string]string{"album": astring}
		_, errt := InsertOne(client, ctx, "albumalpha", "T", item)
		CheckError(errt, "AlbumFirst: T insertion has failed")
		return "T Created"

	case char == "U":
		var item map[string]string = map[string]string{"album": astring}
		_, erru := InsertOne(client, ctx, "albumalpha", "U", item)
		CheckError(erru, "AlbumFirst: U insertion has failed")
		return "U Created"

	case char == "V":
		var item map[string]string = map[string]string{"album": astring}
		_, errv := InsertOne(client, ctx, "albumalpha", "V", item)
		CheckError(errv, "AlbumFirst: V insertion has failed")
		return "V Created"

	case char == "W":
		var item map[string]string = map[string]string{"album": astring}
		_, errw := InsertOne(client, ctx, "albumalpha", "W", item)
		CheckError(errw, "AlbumFirst: W insertion has failed")
		return "W Created"

	case char == "X":
		var item map[string]string = map[string]string{"album": astring}
		_, errx := InsertOne(client, ctx, "albumalpha", "X", item)
		CheckError(errx, "AlbumFirst: X insertion has failed")
		return "X Created"

	case char == "Z":
		var item map[string]string = map[string]string{"album": astring}
		_, errz := InsertOne(client, ctx, "albumalpha", "Z", item)
		CheckError(errz, "AlbumFirst: Z insertion has failed")
		return "Z Created"
	}
	return "None"
}