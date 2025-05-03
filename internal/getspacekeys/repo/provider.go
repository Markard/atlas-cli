package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/Markard/atlas-cli/internal/getspacekeys/entity"
	"github.com/ctreminiom/go-atlassian/confluence/v2"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/ztrue/tracerr"
	"net/url"
)

type ConfluenceProvider struct {
	url      string
	username string
	token    string
}

func NewConfluenceProvider(url string, username string, token string) *ConfluenceProvider {
	return &ConfluenceProvider{url: url, username: username, token: token}
}

func (cp *ConfluenceProvider) GetSpaces(batchSize int) ([]*entity.Space, error) {
	instance, err := confluence.New(nil, cp.url)
	if err != nil {
		return nil, tracerr.Wrap(fmt.Errorf("check atlassian url\nError: %w\n", err))
	}
	instance.Auth.SetBasicAuth(cp.username, cp.token)

	var (
		options = &models.GetSpacesOptionSchemeV2{
			IDs:               nil,
			Keys:              nil,
			Type:              "",
			Status:            "",
			Labels:            nil,
			Sort:              "",
			DescriptionFormat: "",
			SerializeIDs:      false,
		}
		result []*entity.Space
		cursor string
	)
	for {
		spaces, response, err := instance.Space.Bulk(context.Background(), options, cursor, batchSize)
		if err != nil && response != nil {
			err = fmt.Errorf("response code: %v - check atlassian username and token\nError: %w\n",
				response.Code,
				err)
			return nil, tracerr.Wrap(err)
		}

		if spaces == nil {
			return nil, tracerr.Wrap(errors.New("no spaces returned"))
		}

		for _, space := range spaces.Results {
			result = append(result, entity.NewSpace(space.ID, space.Name, space.Key))
		}

		if spaces.Links.Next == "" {
			break
		}

		values, err := url.ParseQuery(spaces.Links.Next)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}

		_, containsCursor := values["cursor"]
		if containsCursor {
			cursor = values["cursor"][0]
		}
	}

	return result, nil
}
