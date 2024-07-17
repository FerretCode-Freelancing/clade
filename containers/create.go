package containers

import (
	"context"
	"encoding/base64"
	"io"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/ferretcode-freelancing/clade/registry"
)

func Create(request Request, client *client.Client, ctx context.Context) error {
	reader, err := client.ImagePull(ctx, request.ImageURL, image.PullOptions{
		PrivilegeFunc: func(ctx context.Context) (string, error) {
			username := registry.RegistryStore.Username
			secret := registry.RegistryStore.Secret

			auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + secret))

			return auth, nil
		},
	})
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	log.Println(string(bytes))

	_, err = client.ContainerCreate(
		ctx, &container.Config{
			Image: request.Image,
		}, nil, nil, nil, request.Name)
	if err != nil {
		return err
	}

	err = WriteCache([]Request{
		request,
	})
	if err != nil {
		return err
	}

	return nil
}
