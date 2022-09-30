package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type SemVer struct {
	Major int
	Minor int
	Patch int
}

func StringToSemVer(v string) (*SemVer, error) {
	semver := &SemVer{}
	log.Println("string (len:", len(v), ") to semver:", v)
	fields := strings.Split(v, ".")
	if len(fields) != 3 {
		return nil, errors.New("not 3 fields (major/minor/patch) in version")
	}
	major, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, errors.New("major version field faulty")
	}

	minor, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil, errors.New("minor version field faulty")
	}
	patch, err := strconv.Atoi(fields[2])
	if err != nil {
		return nil, errors.New("patch version field faulty")
	}
	log.Println("major:", major, "minor:", minor, "patch:", patch)
	semver.Major = major
	semver.Minor = minor
	semver.Patch = patch
	return semver, nil
}

func SemVerToString(version *SemVer) (string, error) {
	if version.Major < 0 {
		return "", errors.New("Version missing major field")
	}
	if version.Minor < 0 {
		return "", errors.New("Version missing minor field")
	}
	if version.Patch < 0 {
		return "", errors.New("Version missing patch field")
	}
	major := strconv.Itoa(version.Major)
	minor := strconv.Itoa(version.Minor)
	patch := strconv.Itoa(version.Patch)
	return major + "." + minor + "." + patch, nil

}

func checkFileIntegrity(v *SemVer) (string, error) {
	version, err := SemVerToString(v)
	if err != nil {
		return "", errors.New("checkFileIntegrity failed, bad semver:" + err.Error())
	}
	h := sha256.New()
	f, err := os.Open("updates/sttylus_update_" + version + ".zip")
	if err != nil {
		return "", errors.New("checkFileIntegrity failed, couldn't open file:" + err.Error())
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil

}
