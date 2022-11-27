/*
 * Otto user service
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

import (
	"encoding/json"
	"net/http"
)

type NewShare struct {
	Text string `json:"text,omitempty"`

	User User `json:"user"`

	Post Post `json:"post"`
}

func DecodeRequestToNewShare(r *http.Request) *NewShare {
	decoder := json.NewDecoder(r.Body)
	var data *NewShare
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}
	return data
}
