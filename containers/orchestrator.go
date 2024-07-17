package containers

import (
	"context"
	"log"

	cn "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func Tick(client *client.Client, ctx context.Context) error {
	requests, err := GetCache()
	if err != nil {
		return err
	}

	for _, request := range requests {
		f := filters.NewArgs(filters.KeyValuePair{Key: "name", Value: request.Name})

		c, err := client.ContainerList(ctx, cn.ListOptions{
			Filters: f,
		})
		if err != nil {
			log.Printf("Err finding container: %v\n", err)
			continue
		}

		if len(c) == 0 {
			err := Create(request, client, ctx)
			if err != nil {
				log.Printf("Err creating new container during tick: %v\n", err)
			}

			continue
		}

		container, err := client.ContainerInspect(ctx, c[0].ID)
		if err != nil {
			log.Println(err)
			continue
		}

		// TODO: make the restartcount configurable by the user
		if container.State.Health.Status == "Unhealthy" && container.RestartCount < 10 {
			err = client.ContainerRestart(ctx, container.ID, cn.StopOptions{})
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}

	return nil
}
