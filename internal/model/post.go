/*
 * Otto user service
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

import "time"

type Post struct {
	Uuid      string    `json:"uuid"`
	Text      string    `json:"text,omitempty"`
	Draft     bool      `json:"draft"`
	Likes     uint      `json:"likes"`
	Replies   uint      `json:"replies"`
	SelfLiked bool      `json:"selfLiked,omitempty"`
	Share     *Post     `json:"share,omitempty"`
	User      User      `json:"user,omitempty"`
	Images    []Image   `json:"images,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
