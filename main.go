package main

import (
	"context"
	"flag"
	"fmt"
	"hash/crc32"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gobwas/glob"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v2"
)

var (
	credentialsFile = flag.String("credentials-file", "", "The Google service account credentials file")
	bucketName      = flag.String("bucket-name", "", "The Google Cloud Storage bucket name")

	configFile = flag.String("config-file", "", "The configuration file")
)

func main() {
	flag.Parse()

	if *bucketName == "" {
		log.Fatalf("Bucket name is not configured")
	}

	now := time.Now()

	config, err := newConfig()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}

	matcher, err := newMatcher(config)
	if err != nil {
		log.Fatalf("Error creating matcher: %v", err)
	}

	client, err := newClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	var paths []string

	err = filepath.WalkDir(".", func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() {
			return nil
		}

		if !matcher.match(path) {
			return nil
		}

		paths = append(paths, path)

		return nil
	})
	if err != nil {
		log.Fatalf("Error walking files: %v", err)
		return
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Error reading file: %v", err)
			continue
		}

		err = client.updateObject(path, data, now)
		if err != nil {
			log.Printf("Error updating object: %v", err)
			continue
		}
	}
}

type Config struct {
	Files []string `json:"files"`
}

func newConfig() (*Config, error) {
	config := Config{}

	if *configFile == "" {
		return &config, nil
	}

	data, err := os.ReadFile(*configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading configuration file '%v': %v", *configFile, err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing configuration file '%v': %v", *configFile, err)
	}

	return &config, nil
}

type Matcher struct {
	globs []glob.Glob
}

func newMatcher(config *Config) (*Matcher, error) {
	matcher := Matcher{}

	for _, pattern := range config.Files {
		glob, err := glob.Compile(pattern, '/')
		if err != nil {
			return nil, fmt.Errorf("error compiling pattern '%v': %v", pattern, err)
		}

		matcher.globs = append(matcher.globs, glob)
	}

	return &matcher, nil
}

func (m *Matcher) match(value string) bool {
	for _, glob := range m.globs {
		if glob.Match(value) {
			return true
		}
	}

	return false
}

type Client struct {
	client *storage.Client
}

func newClient() (*Client, error) {
	ctx := context.Background()

	var options []option.ClientOption
	if *credentialsFile != "" {
		options = append(options, option.WithCredentialsFile(*credentialsFile))
	}

	client, err := storage.NewClient(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) updateObject(path string, data []byte, now time.Time) error {
	ctx := context.Background()

	object := c.client.Bucket(*bucketName).Object(path)

	attrs, err := object.Attrs(ctx)
	if err != nil && err != storage.ErrObjectNotExist {
		return fmt.Errorf("error reading attributes for object '%v': %v", path, err)
	}

	CRC32C := crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))

	if attrs == nil || attrs.CRC32C != CRC32C {
		writer := object.NewWriter(ctx)
		writer.CustomTime = now
		writer.CRC32C = CRC32C
		writer.SendCRC32C = true

		_, err = writer.Write(data)
		if err != nil {
			return fmt.Errorf("error writing data for object '%v': %v", path, err)
		}

		err = writer.Close()
		if err != nil {
			return fmt.Errorf("error writing object '%v': %v", path, err)
		}

		log.Printf("Wrote object '%v'", path)
	} else {
		_, err = object.Update(ctx, storage.ObjectAttrsToUpdate{
			CustomTime: now,
		})
		if err != nil {
			return fmt.Errorf("error updating object attributes for '%v': %v", path, err)
		}

		log.Printf("Updated attributes for object '%v'", path)
	}

	return nil
}
