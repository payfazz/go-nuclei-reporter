package model

import (
	"context"
	"errors"
	"github.com/fadhilthomas/go-nuclei-reporter/config"
	"github.com/jomei/notionapi"
	"github.com/rs/zerolog/log"
)

func OpenNotionDB() (client *notionapi.Client) {
	notionToken := config.GetStr(config.NOTION_TOKEN)
	client = notionapi.NewClient(notionapi.Token(notionToken))
	return client
}

func QueryNotionVulnerabilityName(client *notionapi.Client, vulnerabilityName string) (output []notionapi.Page, err error) {
	databaseId := config.GetStr(config.NOTION_DATABASE)
	databaseQueryRequest := &notionapi.DatabaseQueryRequest{
		PropertyFilter: &notionapi.PropertyFilter{
			Property: "Name",
			Text: &notionapi.TextFilterCondition{
				Equals: vulnerabilityName,
			},
		},
	}
	res, err := client.Database.Query(context.Background(), notionapi.DatabaseID(databaseId), databaseQueryRequest)
	if err != nil {
		log.Error().Str("file", "notion").Msg(err.Error())
		return nil, errors.New(err.Error())
	}
	return res.Results, nil
}

func QueryNotionVulnerabilityStatus(client *notionapi.Client, vulnerabilityStatus string) (output []notionapi.Page, err error) {
	databaseId := config.GetStr(config.NOTION_DATABASE)
	databaseQueryRequest := &notionapi.DatabaseQueryRequest{
		PropertyFilter: &notionapi.PropertyFilter{
			Property: "Status",
			Select: &notionapi.SelectFilterCondition{
				Equals: vulnerabilityStatus,
			},
		},
	}
	res, err := client.Database.Query(context.Background(), notionapi.DatabaseID(databaseId), databaseQueryRequest)
	if err != nil {
		log.Error().Str("file", "notion").Msg(err.Error())
		return nil, errors.New(err.Error())
	}
	return res.Results, nil
}

func InsertNotionVulnerability(client *notionapi.Client, vulnerability Output) (output *notionapi.Page, err error) {
	databaseId := config.GetStr(config.NOTION_DATABASE)

	var multiSelect []notionapi.Option
	for _, tag := range vulnerability.Info.Tags {
		selectOption := notionapi.Option{
			Name: tag,
		}
		multiSelect = append(multiSelect, selectOption)
	}

	pageInsertQuery := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(databaseId),
		},
		Properties: notionapi.Properties{
			"Name": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Text: notionapi.Text{
							Content: vulnerability.TemplateID,
						},
					},
				},
			},
			"Severity": notionapi.SelectProperty{
				Select: notionapi.Option{
					Name: vulnerability.Info.Severity,
				},
			},
			"Host": notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Text: notionapi.Text{
							Content: vulnerability.Host,
						},
					},
				},
			},
			"Endpoint": notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Text: notionapi.Text{
							Content: vulnerability.Matched,
						},
					},
				},
			},
			"Status": notionapi.SelectProperty{
				Select: notionapi.Option{
					Name: "open",
				},
			},
			"Tags": notionapi.MultiSelectProperty{
				ID:          "",
				Type:        "",
				MultiSelect: multiSelect,
			},
			"CVSS Score": notionapi.NumberProperty{
				Number: vulnerability.Info.Classification.CvssScore,
			},
		},
	}

	res, err := client.Page.Create(context.Background(), pageInsertQuery)
	if err != nil {
		log.Error().Str("file", "notion").Msg(err.Error())
		return nil, errors.New(err.Error())
	}
	return res, nil
}

func UpdateNotionVulnerabilityStatus(client *notionapi.Client, pageId string, status string) (output *notionapi.Page, err error) {
	pageUpdateQuery := &notionapi.PageUpdateRequest{
		Properties: notionapi.Properties{
			"Status": notionapi.SelectProperty{
				Select: notionapi.Option{
					Name: status,
				},
			},
		},
	}

	res, err := client.Page.Update(context.Background(), notionapi.PageID(pageId), pageUpdateQuery)
	if err != nil {
		log.Error().Str("file", "notion").Msg(err.Error())
		return nil, errors.New(err.Error())
	}
	return res, nil
}
