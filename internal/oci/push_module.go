package oci

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/credentials"
	"oras.land/oras-go/v2/registry/remote/retry"
)

const (
	ARTIFACT_MANIFEST_MT = "application/vnd.opentofu.module.v1+json"
)

func PushModule(args []string) int {
	var ref, filepath string
	ctx := context.Background()

	if len(args) == 2 {
		ref = args[0]
		filepath = args[1]
	} else {
		fmt.Println("Invalid arguments")
		return 1
	}

	fmt.Println("Pushing module ", ref, " from ", filepath)

	// Prepare layer
	desc, data, err := prepareLayer(filepath, "application/vnd.opentofu.module.v1+zip")
	if err != nil {
		fmt.Println("Failed to prepare layer: ", err)
		return 1
	}

	// Push layer
	store := memory.New()
	if err := store.Push(ctx, desc, data); err != nil {
		fmt.Println("Failed to push layer: ", err)
		return 1
	}

	repo, err := remote.NewRepository(ref)
	if err != nil {
		fmt.Println("Failed to create repository: ", err)
		return 1
	}

	storeOpts := credentials.StoreOptions{}
	credStore, err := credentials.NewStoreFromDocker(storeOpts)
	if err != nil {
		fmt.Println("Failed to create credential store: ", err)
		return 1
	}

	repo.Client = &auth.Client{
		Client:     retry.DefaultClient,
		Cache:      auth.NewCache(),
		Credential: credentials.Credential(credStore),
	}

	tag := "latest"
	manifestDesc, err := oras.Copy(ctx, repo, tag, store, tag, oras.CopyOptions{})
	if err != nil {
		fmt.Println("Failed to copy: ", err)
		return 1
	}

	fmt.Println("Pushed module ", ref, " with digest ", manifestDesc.Digest)
	return 0
}

func prepareLayer(path string, mediaType string) (ocispec.Descriptor, io.Reader, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ocispec.Descriptor{}, nil, err
	}

	desc := ocispec.Descriptor{
		MediaType: mediaType,
		Size:      int64(len(data)),
		Digest:    content.NewDescriptorFromBytes(mediaType, data).Digest,
	}

	return desc, bytes.NewReader(data), nil
}
