package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	fb "firebase.google.com/go"
)

// https://github.com/GoogleCloudPlatform/golang-samples/blob/HEAD/run/helloworld/main.go
func main() {
	log.Print("starting server...")
	http.HandleFunc("/", handler)
	http.HandleFunc("/firestore", firestoreTest)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)

	// Secret managerのdemo
	sc := os.Getenv("TEST_SECRET")
	if sc != "" {
		fmt.Fprintf(w, "TEST_SECRET_VALUE:  %s!\n", sc)
	}
}

// firestoreのdemo
func firestoreTest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conf := &fb.Config{ProjectID: os.Getenv("PROJECT_ID")}
	app, err := fb.NewApp(ctx, conf)
	if err != nil {
		panic(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	type Hoge struct {
		Val string
	}

	_, err = client.Collection("testCollection").Doc(time.Now().String()).Create(ctx, &Hoge{
		Val: "test_val",
	})
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "ok")
	return
}
