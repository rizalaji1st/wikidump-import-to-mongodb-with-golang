package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Redirect struct {
	Namespace int64  `json:"namespace"`
	Title     string `json:"title"`
}

type Coordinates struct {
	Coord   Coord  `json:"coord"`
	Dim     int64  `json:"dim"`
	Globe   string `json:"globe"`
	Name    string `json:"name"`
	Primary bool   `json:"primary"`
	Region  string `json:"region"`
	Type    string `json:"type"`
}

type Coord struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type WikiIndex struct {
	Index Index `json:"index"`
}
type Index struct {
	ID   string `json:"_id,omitempty"`
	Type string `json:"_type,omitempty"`
}
type WikiContent struct {
	ID              string        `json:"id"`
	AuxiliaryText   []string      `json:"auxiliary_text"`
	Category        []string      `json:"category"`
	Coordinates     []Coordinates `json:"coordinates"`
	DefaultSort     interface{}   `json:"defaultsort,omitempty"`
	ExternalLink    []string      `json:"external_link"`
	Heading         []string      `json:"heading"`
	IncomingLinks   int64         `json:"incoming_links"`
	Language        string        `json:"language"`
	Namespace       int64         `json:"namespace"`
	NamespaceText   string        `json:"namespace_text"`
	OpeningText     string        `json:"opening_text"`
	OutgoingLink    []string      `json:"outgoing_link"`
	PopularityScore float64       `json:"popularity_score"`
	Redirect        []Redirect    `json:"redirect"`
	Score           float64       `json:"score"`
	SourceText      string        `json:"source_text"`
	Template        []string      `json:"template"`
	Text            string        `json:"text"`
	TextBytes       float64       `json:"text_bytes"`
	Timestamp       time.Time     `json:"timestamp"`
	Title           string        `json:"title"`
	Version         int64         `json:"version"`
	VersionType     string        `json:"version_type"`
	Wiki            string        `json:"wiki"`
	WikibaseItem    string        `json:"wikibase_item"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()

	//connect db
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://username:password@localhost:27017"))
	if err != nil {
		panic(err)
	}
	//disconnect
	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	//collections
	collections := client.Database("idwiki").Collection("content")

	f, err := os.Open("dump.json")
	if err != nil {
		fmt.Print("Error open file")
	}
	d := json.NewDecoder(f)
	for {
		var v WikiIndex
		if err := d.Decode(&v); err == io.EOF {
			break // done decoding file
		} else if err != nil {
			fmt.Println("Error decode index")
		}
		// out, _ := json.Marshal(v)
		// fmt.Println(string(out))
		var w WikiContent
		if err := d.Decode(&w); err == io.EOF {
			break // done decoding file
		} else if err != nil {
			fmt.Println("error in " + v.Index.ID)
			fmt.Println(err)
		}
		w.ID = v.Index.ID
		if w.PopularityScore == 0 {
			w.PopularityScore = math.Pow(10, -20)
		}
		// oute, _ := json.Marshal(w)
		// fmt.Println(string(oute))
		data := bson.D{
			// {"id", w.ID},
			{Key: "id", Value: w.ID},
			{Key: "auxiliary_text", Value: w.AuxiliaryText},
			{Key: "category", Value: w.Category},
			{Key: "coordinates", Value: w.Coordinates},
			{Key: "default_sort", Value: w.DefaultSort},
			{Key: "external_link", Value: w.ExternalLink},
			{Key: "heading", Value: w.Heading},
			{Key: "incoming_links", Value: w.IncomingLinks},
			{Key: "language", Value: w.Language},
			{Key: "namespace", Value: w.Namespace},
			{Key: "namespace_text", Value: w.NamespaceText},
			{Key: "opening_text", Value: w.OpeningText},
			{Key: "outgoing_link", Value: w.OutgoingLink},
			{Key: "popularity_score", Value: w.PopularityScore},
			{Key: "redirect", Value: w.Redirect},
			{Key: "score", Value: w.Score},
			{Key: "source_text", Value: w.SourceText},
			{Key: "template", Value: w.Template},
			{Key: "text", Value: w.Text},
			{Key: "text_bytes", Value: w.TextBytes},
			{Key: "timestamp", Value: w.Timestamp},
			{Key: "title", Value: w.Title},
			{Key: "version", Value: w.Version},
			{Key: "version_type", Value: w.VersionType},
			{Key: "wiki", Value: w.Wiki},
			{Key: "wikibase_item", Value: w.WikibaseItem},
		}
		resultSet, err := collections.InsertOne(ctx, data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resultSet.InsertedID)
	}
}
