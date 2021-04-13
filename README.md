# smk-open-to-wof

The National Gallery of Denmark (SMK) wrote a blog post about geolocating works in their collection:

* [We’re all over the map — how we geolocated SMK artworks with the kind help of humans and machines](https://medium.com/smk-open/were-all-over-the-map-how-we-geolocated-smk-artworks-with-the-kind-help-of-humans-and-machines-91f652295563)

And then they plotted all those locations on a map:

* https://danmarkskort.open.smk.dk/

And I thought "Wouldn't it be nice if all those locations (coordinates) had corresponding Who's On First (WOF) IDs because [sometime place is not always spatial](https://whosonfirst.org/what)". Conveniently SFO Museum has recently published some tooling to do just that:

* [Reverse-Geocoding in Time at SFO Museum](https://millsfield.sfomuseum.org/blog/2021/03/26/spatial/)

This repository contains a mapping of SMK object IDs and their latitude and longitude coordinates to corresponding WOF ID, name and placetype properties.

* [smk-open-wof.csv](smk-open-wof.csv)

And now it's possible to query for objects in the SMK collection using a WOF ID. For example, here are works in the neighbourhood of [Niels Juels Gade](https://spelunker.whosonfirst.org/id/85906061) (85906061):

```
$> grep 85906061 ./smk-open-wof.csv 
70,55.67634666739197,12.587949512316898,85906061,101749159,Niels Juels Gade,neighbourhood
441,55.675375135166256,12.587992427661137,85906061,101749159,Niels Juels Gade,neighbourhood
1421,55.67605987856784,12.586535647521977,85906061,101749159,Niels Juels Gade,neighbourhood
2778,55.67678132295374,12.589279887988285,85906061,101749159,Niels Juels Gade,neighbourhood
3193,55.67739609142182,12.587052924731443,85906061,101749159,Niels Juels Gade,neighbourhood
```

Here are all the airlines from Denmark (85633121) that hold hands with the SFO Museum collection:

* https://millsfield.sfomuseum.org/countries/85633121/airlines/

## Do it yourself

These are the steps for building a mapping of SMK object IDs and their latitude and longitude coordinates to corresponding WOF ID, name and placetype properties from scratch.

### Step 1

Download [Who's On First data](https://github.com/whosonfirst-data/whosonfirst-data-admin-dk) for Denmark.

```
$> git clone git@github.com:whosonfirst-data/whosonfirst-data-admin-dk.git /path/to/whosonfirst-data-admin-dk
```

### Step 2

Create a spatially-enabled SQLite database for the WOF data for Denmark using the tools in the [go-whosonfirst-sqlite-features-index](https://github.com/whosonfirst/go-whosonfirst-sqlite-features-index) package.

```
$> cd go-whosonfirst-sqlite-features-index

$> make cli
go build -mod vendor -o bin/wof-sqlite-index-features cmd/wof-sqlite-index-features/main.go

$> ./bin/wof-sqlite-index-features -all -dsn /path/to/whosonfirst-dk.sqlite /path/to/whosonfirst-data-admin-dk
```

### Step 3

Expose the SQLite database with WOF data for Denmark as a simple HTTP service using the tools in the [go-whosonfirst-spatial-www-sqlite](https://github.com/whosonfirst/go-whosonfirst-spatial-www-sqlite) package.

```
$> cd go-whosonfirst-sqlite-spatial-www-sqlite
go build -mod vendor -o bin/server cmd/server/main.go

$> make cli

$> ./bin/server -spatial-database-uri 'sqlite://?dsn=/path/to/whosonfirst-dk.sqlite'
```

### Step 4

Get the SMK Open map data.

```
$> curl -o smk-open.json https://danmarkskort.open.smk.dk/get_list
```

### Step 5

Run the `to-wof` tool included in this repository. This will read the SMK Open map data and look up the WOF ID(s) that contain the latitude and longitude coordinate for each work. The merged results will be emitted as CSV data to STDOUT.

```
$> cd to-wof
$> go run main.go -source ../smk-open.json > ../smk-open-wof.csv
```

## See also

* https://spelunker.whosonfirst.org/id/85633121 (Denmark)
* https://github.com/whosonfirst-data/whosonfirst-data-admin-dk
* https://github.com/whosonfirst/go-whosonfirst-spatial-www-sqlite
* https://github.com/whosonfirst/go-whosonfirst-sqlite-features-index