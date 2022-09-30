package main

import (
	"errors"
	"log"
	"time"
)

func dbGetUser(id int64) (*User, error) {
	for _, u := range users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func dbNewManifest(manifest *Manifest) (*SemVer, error) {
	m := *manifest
	semver, err := StringToSemVer(m.Version)
	if err != nil {
		return semver, errors.New("creating manifest failed, faulty semver:" + err.Error())
	}
	m.Major = semver.Major
	m.Minor = semver.Minor
	m.Patch = semver.Patch
	m.Date = time.Now()
	hash, err := checkFileIntegrity(semver)
	initChangeLog(&m)
	log.Println("post:", manifest.Hash)
	log.Println("api:", hash)
	if err != nil {
		return nil, errors.New("creating manifest failed, checksum failed:" + err.Error())

	}
	if manifest.Hash != hash {
		return nil, errors.New("creating manifest failed, checksum not matching")
	}

	m.Hash = hash
	manifests[m.Version] = &m
	return semver, err

}

func dbGetManifest(s *SemVer) (*Manifest, error) {
	log.Printf("dbGetManifests: %+v\n", s)
	var manifest *Manifest
	for _, m := range manifests {
		log.Println(m)
		if m.Major == s.Major {
			if m.Minor == s.Minor {
				if m.Patch == s.Patch {
					manifest = m
				}
			}
		}
	}

	return manifest, nil
}

func dbgetManifests(major, minor, patch int) ([]*Manifest, error) {
	return nil, nil

}
