//go:build integration
// +build integration

package main_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const url string = "http://localhost:8080"

func TestIntegration(t *testing.T) {
	c := &http.Client{Timeout: 30 * time.Second}

	testGetMatches(t, c)
	testGetPartnerById(t, c)
}

func testGetMatches(t *testing.T, c *http.Client) {
	getMatchBody := `
	{
		"materials": [1, 2],
		"address": {
			"lat": 1.1,
			"long": 1.1
		},
		"square_meters": 5,
		"phone_number": "+351912345678"
	}
	`
	req, _ := http.NewRequest(http.MethodPost, url+"/partners/match", strings.NewReader(getMatchBody))
	resp, _ := c.Do(req)
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)

	expectedPartners := `
	[
		{
			"id": 2,
			"categories": [
				{
					"id": 1,
					"description": "Flooring materials"
				}
			],
			"materials": [
				{
					"id": 1,
					"description": "Wood"
				},
				{
					"id": 2,
					"description": "Carpet"
				},
				{
					"id": 3,
					"description": "Tile"
				}
			],
			"address": {
				"lat": 1.2,
				"long": 1.2
			},
			"radius": 200,
			"rating": 3
		},
		{
			"id": 4,
			"categories": [
				{
					"id": 1,
					"description": "Flooring materials"
				}
			],
			"materials": [
				{
					"id": 1,
					"description": "Wood"
				},
				{
					"id": 2,
					"description": "Carpet"
				}
			],
			"address": {
				"lat": 1.4,
				"long": 1.4
			},
			"radius": 200,
			"rating": 2
		},
		{
			"id": 3,
			"categories": [
				{
					"id": 1,
					"description": "Flooring materials"
				}
			],
			"materials": [
				{
					"id": 1,
					"description": "Wood"
				},
				{
					"id": 2,
					"description": "Carpet"
				},
				{
					"id": 3,
					"description": "Tile"
				}
			],
			"address": {
				"lat": 1.1,
				"long": 1.1
			},
			"radius": 200,
			"rating": 1
		},
		{
			"id": 1,
			"categories": [
				{
					"id": 1,
					"description": "Flooring materials"
				}
			],
			"materials": [
				{
					"id": 1,
					"description": "Wood"
				},
				{
					"id": 2,
					"description": "Carpet"
				},
				{
					"id": 3,
					"description": "Tile"
				}
			],
			"address": {
				"lat": 1.3,
				"long": 1.3
			},
			"radius": 200,
			"rating": 1
		}
	]
	`
	buffer := new(bytes.Buffer)
	_ = json.Compact(buffer, []byte(expectedPartners))

	if diff := cmp.Diff(buffer.String(), string(b)); diff != "" {
		t.Errorf("guest list mismatch (-want +got):\n%s", diff)
	}
}

func testGetPartnerById(t *testing.T, c *http.Client) {
	req, _ := http.NewRequest(http.MethodGet, url+"/partners/1", nil)
	resp, _ := c.Do(req)
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)

	expectedPartner := `
	{
		"id": 1,
		"categories": [
			{
				"id": 1,
				"description": "Flooring materials"
			}
		],
		"materials": [
			{
				"id": 1,
				"description": "Wood"
			},
			{
				"id": 2,
				"description": "Carpet"
			},
			{
				"id": 3,
				"description": "Tile"
			}
		],
		"address": {
			"lat": 1.3,
			"long": 1.3
		},
		"radius": 200,
		"rating": 1
	}
	`

	buffer := new(bytes.Buffer)
	_ = json.Compact(buffer, []byte(expectedPartner))

	if diff := cmp.Diff(buffer.String(), string(b)); diff != "" {
		t.Errorf("guest list mismatch (-want +got):\n%s", diff)
	}
}
