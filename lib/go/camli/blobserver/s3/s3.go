/*
Copyright 2011 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package s3

import (
	"fmt"
	"http"
	"os"

	"camli/blobserver"
	"camli/misc/amazon/s3"
)

type s3Storage struct {
	*blobserver.SimpleBlobHubPartitionMap
	s3Client *s3.Client
	bucket   string
}

func newFromConfig(config blobserver.JSONConfig) (storage blobserver.Storage, err os.Error) {
	client := &s3.Client{
		Auth: &s3.Auth{
			AccessKey:       config.RequiredString("aws_access_key"),
			SecretAccessKey: config.RequiredString("aws_secret_access_key"),
		},
		HttpClient: http.DefaultClient,
	}
	sto := &s3Storage{
		SimpleBlobHubPartitionMap: &blobserver.SimpleBlobHubPartitionMap{},
		s3Client:                  client,
		bucket:                    config.RequiredString("bucket"),
	}
	skipStartupCheck := config.OptionalBool("skipStartupCheck", false)
	if err := config.Validate(); err != nil {
		return nil, err
	}
	if !skipStartupCheck {
		// TODO: skip this check if a file
		// ~/.camli/.configcheck/sha1-("IS GOOD: s3: sha1(access key +
		// secret key)") exists and has recent time?
		if _, err := client.Buckets(); err != nil {
			return nil, fmt.Errorf("Failed to get bucket list from S3: %v", err)
		}
	}
	return sto, nil
}

func init() {
	blobserver.RegisterStorageConstructor("s3", blobserver.StorageConstructor(newFromConfig))
}

