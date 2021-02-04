package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func addFile(tw *tar.Writer, path, arpath string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if arpath == "" {
		arpath = path
	}
	if stat, err := file.Stat(); err == nil {
		// now lets create the header as needed for this file within the tarball
		header := new(tar.Header)
		header.Name = arpath
		header.Size = stat.Size()
		header.Mode = int64(stat.Mode())
		header.ModTime = stat.ModTime()
		// write the header to the tarball archive
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// copy the file data to the tarball
		if _, err := io.Copy(tw, file); err != nil {
			return err
		}
	}
	return nil
}

//CompressFile compress files to .tar.gz
func CompressFile(output string, paths, arpaths []string) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()
	// set up the gzip writer
	gw := gzip.NewWriter(file)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	// grab the paths that need to be added in
	// add each file as needed into the current tar archive
	for i := range paths {
		if err := addFile(tw, paths[i], arpaths[i]); err != nil {
			return err
		}
	}
	return nil
}

func hash(obj interface{}) [32]byte {
	val := fmt.Sprintf("%v", obj)
	return sha256.Sum256([]byte(val))
}

func hash64(obj interface{}) string {
	b := hash(obj)
	return base64.StdEncoding.EncodeToString(b[:])
}

//HashFile calculate file hash
func HashFile(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write(b)
	md := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(md), nil
}

//ObjectToFile write object to file
func ObjectToFile(filepath string, obj interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, b, 0644)
}

//FileExists check if file exist
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
