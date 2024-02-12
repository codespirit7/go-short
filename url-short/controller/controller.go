package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/codespirit7/url-shortner/model"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URLJson struct {
	URL string
}

const dbName = "shorturl"
const collectionName = "urls"

// create reference to mongoDb collection
var collection *mongo.Collection

// connect with mongoDB
func init() {
	godotenv.Load()
	var token = os.Getenv("MONGO_URL")

	var connectionString = token
	clienOption := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clienOption)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connection succesful")

	collection = client.Database(dbName).Collection(collectionName)
}

// insert a record
func addUrl(url model.URLSchema) *mongo.InsertOneResult {
	data, err := collection.InsertOne(context.Background(), url)

	if err != nil {
		log.Fatal(err)
	}

	return data
}

// find one record
func getOriginalUrl(shortId string) (model.URLSchema, error) {
	filter := bson.M{"shortid": shortId}

	var result model.URLSchema

	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.URLSchema{}, errors.New("URL not found")
		}
		return model.URLSchema{}, err
	}

	return result, nil

}

// generate short key
func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 5

	shortKey := make([]byte, keyLength)

	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}

	return string(shortKey)
}

// generate short url
func HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var urlInfo URLJson

	//Decoding url from the request body and attaching it to variable url
	err := json.NewDecoder(r.Body).Decode(&urlInfo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Generating a unique short id
	shortId := generateShortKey()

	// data contains {shortId, originalURL}
	data := model.URLSchema{}

	data.ShortId = shortId
	data.OriginalURL = urlInfo.URL

	//Add to mongoDB database
	addUrl(data)

	var resp = struct {
		URLShortId string `json:"short-url"`
	}{
		URLShortId: fmt.Sprintf("http://localhost:8080/short/%s", shortId),
	}

	//sending response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)

}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len("/short/"):]
	fmt.Println("key", shortKey)
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	data, err := getOriginalUrl(shortKey)
	if err != nil {
		http.Error(w, "Short url not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, data.OriginalURL, http.StatusMovedPermanently)
}
