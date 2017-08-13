package database

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"sync"

	"github.com/ngalayko/highloadcup/schema"
)

var (
	fileNameRegex = regexp.MustCompile(`(?P<entity>\w+)_\d+.json`)
)

// ParseData parses data from files
func (db *DB) ParseData(dataPath string) error {
	log.Printf("start loading data from %s", dataPath)

	reader, err := zip.OpenReader(dataPath)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for _, file := range reader.File {
		if file.Name == path.Base(dataPath) {
			continue
		}

		fileName := file.Name

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			log.Printf("loading %s\n", fileName)
			if err := db.parseFile(fileReader, fileName); err != nil {
				log.Panic("error when parsing %s: %s", fileName, err)
			}

			log.Printf("%s loaded\n", fileName)
		}()
	}

	wg.Wait()

	log.Println("updating generic values")
	if err := db.updateGenericValues(); err != nil {
		return err
	}

	log.Println("data loaded and updated")

	return nil
}

func (db *DB) updateGenericValues() (err error) {
	locationMap := db.mapByEntity(schema.EntityLocations)
	usersMap := db.mapByEntity(schema.EntityUsers)
	visitsMap := db.mapByEntity(schema.EntityVisits)

	userToVisitsMap := map[uint32][]uint32{}
	locationToVisitsMap := map[uint32][]uint32{}

	visitsMap.Range(func(k, v interface{}) bool {
		if visit, ok := v.(*schema.Visit); ok {
			locationToVisitsMap[visit.LocationID] = append(locationToVisitsMap[visit.LocationID], visit.ID)
			userToVisitsMap[visit.UserID] = append(userToVisitsMap[visit.UserID], visit.ID)
		}
		return true
	})

	locationMap.Range(func(k, v interface{}) bool {
		if location, ok := v.(*schema.Location); ok {
			location.VisitIDs = locationToVisitsMap[location.ID]
		}
		return true
	})

	usersMap.Range(func(k, v interface{}) bool {
		if user, ok := v.(*schema.User); ok {
			user.VisitIDs = userToVisitsMap[user.ID]
		}
		return true
	})

	return nil
}

func (db *DB) parseFile(fileReader io.Reader, fileName string) error {
	fileEntity, err := parseFileName(fileName)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(fileReader)
	if err != nil {
		return err
	}

	switch fileEntity {
	case schema.EntityUsers:

		users := &schema.Users{}
		if err := json.Unmarshal(data, users); err != nil {
			return err
		}

		usersMap := db.mapByEntity(schema.EntityUsers)
		for _, user := range users.Users {
			usersMap.Store(user.ID, user)
		}

	case schema.EntityVisits:
		visits := &schema.Visits{}
		if err := json.Unmarshal(data, visits); err != nil {
			return err
		}

		visitsMap := db.mapByEntity(schema.EntityVisits)
		for _, visit := range visits.Visits {
			visitsMap.Store(visit.ID, visit)
		}

	case schema.EntityLocations:
		locations := &schema.Locations{}
		if err := json.Unmarshal(data, locations); err != nil {
			return err
		}

		locationMap := db.mapByEntity(schema.EntityLocations)
		for _, location := range locations.Locations {
			locationMap.Store(location.ID, location)
		}
	}

	return nil
}

func parseFileName(fileName string) (schema.Entity, error) {
	match := fileNameRegex.FindStringSubmatch(fileName)

	var entity schema.Entity
	var found bool

	for i, name := range fileNameRegex.SubexpNames() {
		if i > len(match) {
			continue
		}

		if name != "entity" {
			continue
		}

		if err := entity.UnmarshalText([]byte(match[i])); err != nil {
			return entity, err
		}

		found = true
	}

	if !found {
		return entity, errors.New("file name not valid")
	}

	return entity, nil
}
