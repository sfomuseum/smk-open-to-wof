package main

import (
	"flag"
	"log"
	"os"
	"io"
	"encoding/json"
	"encoding/csv"
	"strconv"
	"net/http"
	"bytes"
)

type SPRResults struct {
	Places []*SPRResult `json:"places"`
}

type SPRResult struct {
	Id string `json:"wof:id"`
	ParentId string `json:"wof:parent_id"`
	Name string `json:"wof:name"`
	Placetype string `json:"wof:placetype"`
}

type PointInPolygonRequest struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`	
}

type CollectionItem struct {
	Id int
	Title string
	AccessionNumber string
	Latitude float64
	Longitude float64
	ImageURL string
}

func FetchList(uri string) (io.ReadCloser, error) {

	fh, err := os.Open(uri)

	if err != nil {
		return nil, err
	}
	
	return fh, nil
}

func GetList(uri string) ([]*CollectionItem, error) {

	fh, err := FetchList(uri)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	type response struct {
		Results [][]interface{} `json:"res"`
	}

	var rsp *response

	dec := json.NewDecoder(fh)
	err = dec.Decode(&rsp)

	if err != nil {
		return nil, err
	}

	items := make([]*CollectionItem, len(rsp.Results))
	
	for idx, r := range rsp.Results {

		i, err := GetItem(r)

		if err != nil {
			return nil, err
		}

		items[idx] = i
	}

	return items, nil
}

func GetItem(i []interface{}) (*CollectionItem, error) {

	id := int(i[0].(float64))

	title := i[1].(string)
	accno := i[2].(string)

	lat, err := strconv.ParseFloat(i[3].(string), 10)

	if err != nil {
		return nil, err
	}
	
	lon, err := strconv.ParseFloat(i[4].(string), 10)

	if err != nil {
		return nil, err
	}

	image_url := i[5].(string)
	
	item := &CollectionItem{
		Id: id,
		Title: title,
		AccessionNumber: accno,
		Latitude: lat,
		Longitude: lon,
		ImageURL: image_url,
	}

	return item, nil
}

func GetLocation(client *http.Client, item *CollectionItem) (*SPRResults, error) {

	req := PointInPolygonRequest{
		Latitude: item.Latitude,
		Longitude: item.Longitude,
	}

	body, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	br := bytes.NewReader(body)
	
	http_req, err := http.NewRequest("POST", "http://localhost:8080/api/point-in-polygon", br)

	if err != nil {
		return nil, err
	}

	rsp, err := client.Do(http_req)

	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	var spr *SPRResults

	dec := json.NewDecoder(rsp.Body)
	err = dec.Decode(&spr)

	if err != nil {
		return nil, err
	}

	return spr, nil
}

func main() {

	src := flag.String("source", "", "")
	flag.Parse()

	writers := []io.Writer{
		os.Stdout,
	}

	wr := io.MultiWriter(writers...)

	csv_writer := csv.NewWriter(wr)
	
	items, err := GetList(*src)

	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	
	for idx, i := range items {
		
		spr, err := GetLocation(client, i)

		if err != nil {
			log.Fatal(err)
		}

		if idx == 0 {

			out := []string{
				"smk_id",
				// "smk_title",
				"smk_latitude",
				"smk_longitude",
				"wof_id",
				"wof_parentid",
				"wof_name",
				"wof_placetype",
			}

			csv_writer.Write(out)
		}

		for _, pl := range spr.Places {

			out := []string{
				strconv.Itoa(i.Id),
				// i.Title,
				strconv.FormatFloat(i.Latitude, 'f', -1, 64),
				strconv.FormatFloat(i.Longitude, 'f', -1, 64),				
				pl.Id,
				pl.ParentId,
				pl.Name,
				pl.Placetype,
			}

			csv_writer.Write(out)
		}
	}

	csv_writer.Flush()
}
	
