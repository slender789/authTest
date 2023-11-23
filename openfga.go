package auth_library

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	fgaSdk "github.com/openfga/go-sdk"
)

type PermissionApi interface {
	Check(user, relation, object string) error
}

type OpenFgaPermissionApi struct {
	client *fgaSdk.APIClient
}

func NewOpenFgaPermissionApi(config *fgaSdk.Configuration) *OpenFgaPermissionApi {
	return &OpenFgaPermissionApi{fgaSdk.NewAPIClient(config)}
}

func (OFPA OpenFgaPermissionApi) Check(user, relation, object string) error {
	body := fgaSdk.CheckRequest{
		TupleKey: fgaSdk.TupleKey{
			User:     &user,
			Relation: &relation,
			Object:   &object,
		},
	}
	_, response, err := OFPA.client.OpenFgaApi.Check(context.Background()).Body(body).Execute()
	if err != nil {
		return err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var authorizationResponse CheckResponse
	err = json.Unmarshal(responseBody, &authorizationResponse)
	if err != nil {
		return err
	}

	if authorizationResponse.Allowed {
		return nil
	}
	return errors.New("access denied")
}

func (OFPA OpenFgaPermissionApi) DeleteTuple(user, relation, object string) error {
	body := fgaSdk.WriteRequest{
		Deletes: &fgaSdk.TupleKeys{
			TupleKeys: []fgaSdk.TupleKey{
				{
					User:     &user,
					Relation: &relation,
					Object:   &object,
				},
			},
		},
	}

	_, response, err := OFPA.client.OpenFgaApi.WriteExecute(fgaSdk.ApiWriteRequest.Body(fgaSdk.ApiWriteRequest{}, body))
	if err != nil {
		return err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var deleteionResponse TupleMutationResponse
	err = json.Unmarshal(responseBody, &deleteionResponse)
	if err != nil {
		return err
	}

	if deleteionResponse.Code == "" && deleteionResponse.Message == "" {
		return nil
	}
	return errors.New("access denied")
}

func (OFPA OpenFgaPermissionApi) AddTuple(user, relation, object string) error {
	body := fgaSdk.WriteRequest{
		Writes: &fgaSdk.TupleKeys{
			TupleKeys: []fgaSdk.TupleKey{
				{
					User:     &user,
					Relation: &relation,
					Object:   &object,
				},
			},
		},
	}

	_, response, err := OFPA.client.OpenFgaApi.WriteExecute(fgaSdk.ApiWriteRequest.Body(fgaSdk.ApiWriteRequest{}, body))
	if err != nil {
		return err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var additionResponse TupleMutationResponse
	err = json.Unmarshal(responseBody, &additionResponse)
	if err != nil {
		return err
	}

	if additionResponse.Code == "" && additionResponse.Message == "" {
		return nil
	}
	return errors.New("access denied")
}
