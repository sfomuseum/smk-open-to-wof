# smk-open-to-wof

The really short version is that the National Gallery of Denmark (SMK) wrote a blog post about geolocating works in their collection:

* [We’re all over the map — how we geolocated SMK artworks with the kind help of humans and machines](https://medium.com/smk-open/were-all-over-the-map-how-we-geolocated-smk-artworks-with-the-kind-help-of-humans-and-machines-91f652295563)

And then they plotted all those locations on a map:

* https://danmarkskort.open.smk.dk/

And I thought "Wouldn't it be nice if all those locations (coordinates) had corresponding Who's On First (WOF) IDs because [sometime place is not always spatial](https://whosonfirst.org/what)". Conveniently SFO Museum has recently published some tooling to do just that:

* [Reverse-Geocoding in Time at SFO Museum](https://millsfield.sfomuseum.org/blog/2021/03/26/spatial/)

This repository contains a mapping of SMK object IDs and their latitude and longitude coordinates to corresponding WOF ID, name and placetype properties.

* [smk-open-wof.csv](smk-open-wof.csv)

It also contains the raw SMK data source used to build this mapping:

* [smk-open.json](smk-open.json)

I was going to include the raw WOF (SQLite) data as well but it's 385MB. Instructions for building it from scratch are included below.

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

Run the `to-wof` tool included in this repository. This will read the SMK Open map data

```
$> cd to-wof
$> go run main.go -source ../smk-open.json > ../smk-open-wof.csv
```

## See also

* https://github.com/whosonfirst-data/whosonfirst-data-admin-dk
* https://github.com/whosonfirst/go-whosonfirst-spatial-www-sqlite
* https://github.com/whosonfirst/go-whosonfirst-sqlite-features-index