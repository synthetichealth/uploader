package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/intervention-engine/fhir/models"
)

func main() {
	start := time.Now()

	fhir := flag.String("fhir", "", "Endpoint for the FHIR server")
	path := flag.String("path", "", "Path to the folder containing records to upload")
	flag.Parse()

	if *fhir == "" || *path == "" {
		fmt.Println("Must provide parameter values for fhir and path")
		os.Exit(1)
	}

	endpoint := *fhir
	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}

	files, err := ioutil.ReadDir(*path)
	if err != nil {
		fmt.Printf("Error reading directory %s: %s\n", *path, err.Error())
		os.Exit(1)
	}

	fmt.Printf("Uploading bundles from %s to %s ", *path, endpoint)
	var total, success int
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		total++

		fpath := filepath.Join(*path, file.Name())

		fmt.Print(".")
		bundle, err := loadBundle(fpath)
		if err != nil {
			fmt.Printf("Could not load bundle from %s: %s\n", fpath, err.Error())
			continue
		}

		updateBundleForTx(bundle)

		if err := postBundle(bundle, endpoint); err != nil {
			fmt.Printf("Error posting %s: %s\n", fpath, err.Error())
			continue
		}
		success++
	}
	fmt.Println()

	duration := time.Now().Sub(start)
	fmt.Printf("Uploaded %d of %d records in %s\n", success, total, duration.String())
	if success > 0 {
		average := duration / time.Duration(int64(success))
		fmt.Printf("Average upload time: %s\n", average.String())
	}
}

func loadBundle(fpath string) (*models.Bundle, error) {
	// Open the file
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Deserialize it to a bundle
	bundle := new(models.Bundle)
	if err := json.NewDecoder(f).Decode(bundle); err != nil {
		return nil, err
	}

	return bundle, nil
}

func updateBundleForTx(bundle *models.Bundle) {
	bundle.Type = "transaction"
	for i := range bundle.Entry {
		bundle.Entry[i].Request = &models.BundleEntryRequestComponent{
			Method: "POST",
			Url:    reflect.TypeOf(bundle.Entry[i].Resource).Elem().Name(),
		}
	}
}

func postBundle(bundle *models.Bundle, endpoint string) error {
	// Encode the bundle and post it
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(bundle)
	res, err := http.Post(endpoint, "application/json", b)
	if err != nil {
		return err
	}

	// Drain the response and check the status code
	defer res.Body.Close()
	io.Copy(ioutil.Discard, res.Body)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected response code 200 but got %d", res.StatusCode)
	}

	return nil
}
