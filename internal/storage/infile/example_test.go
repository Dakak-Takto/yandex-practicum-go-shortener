package infile

import (
	"fmt"
	"log"
	"os"
)

func Example() {
	const (
		storageFilename string = "filename.urls"
		originalURL     string = "http://original.url/string/foo=bar"
		shortKey        string = "shortKeyString"
		userID          string = "barnie"
	)

	defer os.Remove(storageFilename)

	// init new storage in file
	storage, err := New(storageFilename)
	if err != nil {
		log.Fatal(err)
	}

	//save url record
	if err := storage.Save(shortKey, originalURL, userID); err != nil {
		log.Fatal(err)
	}

	// get URL by short key
	url, err := storage.GetByShort(shortKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(url.Original)

	// get byt original URL
	url, err = storage.GetByOriginal(originalURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(url.Original)

	//get user urls by userID
	urls, err := storage.SelectByUID(userID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(urls)

	//delete record
	storage.Delete(userID, shortKey)
}
