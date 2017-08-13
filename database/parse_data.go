package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sync"
	"archive/zip"
	"path/filepath"
	"path"
	"io"

	"github.com/ngalayko/highloadcup/helper"
	"github.com/ngalayko/highloadcup/schema"
)

var (
	fileNameRegex = regexp.MustCompile(`(?P<entity>\w+)_\d+.json`)
)

// ParseData parses data from files
func (db *DB) ParseData(dataPath string) error {
	log.Printf("start loading data from %s", dataPath)

	if err := unzip(dataPath, path.Dir(dataPath)); err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	files, _ := ioutil.ReadDir(path.Dir(dataPath))
	for _, f := range files {

		if f.Name() == path.Base(dataPath) {
			continue
		}

		filename := f.Name()

		wg.Add(1)
		go func() {
			defer wg.Done()

			log.Printf("loading %s\n", filename)
			if err := db.parseFile(fmt.Sprintf("%s/%s", path.Dir(dataPath), filename)); err != nil {
				log.Panic("error when parsing %s: %s", filename, err)
			}

			log.Printf("%s loaded\n", filename)
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
	locations, err := db.LoadAllLocations()
	if err != nil {
		return err
	}

	visits, err := db.LoadAllVisits()
	if err != nil {
		return err
	}

	users, err := db.LoadAllUsers()
	if err != nil {
		return err
	}

	userToVisitsMap := map[uint32][]uint32{}
	locationToMarksMap := map[uint32][]uint8{}
	for _, visit := range visits {
		locationToMarksMap[visit.LocationID] = append(locationToMarksMap[visit.LocationID], visit.Mark)
		userToVisitsMap[visit.UserID] = append(userToVisitsMap[visit.UserID], visit.ID)
	}

	for _, location := range locations {
		location.AvgMark = helper.Avg(locationToMarksMap[location.ID]...)

		if err := db.CreateOrUpdate(location); err != nil {
			return err
		}
	}

	for _, user := range users {
		user.VisitIDs = userToVisitsMap[user.ID]

		if err := db.CreateOrUpdate(user); err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) parseFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	fileEntity, err := parseFileName(file.Name())
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	switch fileEntity {
	case schema.EntityUsers:
		users := &schema.Users{}
		if err := json.Unmarshal(data, users); err != nil {
			return err
		}

		if err := db.CreateUsers(users); err != nil {
			return err
		}

	case schema.EntityVisits:
		visits := &schema.Visits{}
		if err := json.Unmarshal(data, visits); err != nil {
			return err
		}

		if err := db.CreateVisits(visits); err != nil {
			return err
		}

	case schema.EntityLocations:
		locations := &schema.Locations{}
		if err := json.Unmarshal(data, locations); err != nil {
			return err
		}

		if err := db.CreateLocations(locations); err != nil {
			return err
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

func unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}
