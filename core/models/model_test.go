package models_test

import (
	"encoding/json"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/core/models"
)

func TestViewJSONMarshaling(t *testing.T) {
	query := "Log {address topics data transactionHash blockNumber}"

	sdl := "type FilteredAndDecodedLogs @materialized(if: false) {hash: String block: String address: String signature: String }"

	lensA := models.Lens{
		Label: "filter_usdt",
		Path:  "assets/lens_filter_usdt.wasm",
		Arguments: map[string]any{
			"src":   "address",
			"value": "0xdac17f958d2ee523a2206206994597c13d831ec7",
		},
	}

	lensB := models.Lens{
		Label: "decode_inputs",
		Path:  "assets/lens_decode_inputs.wasm",
		Arguments: map[string]any{
			"abi": `{"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}]}`,
		},
	}

	transform := models.Transform{
		Lenses: []models.Lens{lensA, lensB},
	}

	revisionA := models.Revision{
		Version:   0,
		Timestamp: "",
		Diff:      "{\"name\":\"example\",\"query\":null,\"sdl\":null,\"transform\":{\"lenses\":[]},\"metadata\":{\"_v\":0,\"_t\":0,\"revisions\":[],\"createdAt\":\"1749746283\",\"updatedAt\":\"1749746283\"}}",
	}

	metadata := models.Metadata{
		Version:   2,
		Total:     2,
		Revisions: []models.Revision{revisionA},
		CreatedAt: "1749746283",
		UpdatedAt: "1749748820",
	}

	view := models.View{
		Name:      "example",
		Query:     &query,
		Sdl:       &sdl,
		Transform: transform,
		Metadata:  metadata,
	}

	jsonBytes, err := json.Marshal(view)
	if err != nil {
		t.Fatalf("Failed to marshal view: %v", err)
	}

	expected := `{"name":"example","query":"Log {address topics data transactionHash blockNumber}","sdl":"type FilteredAndDecodedLogs @materialized(if: false) {hash: String block: String address: String signature: String }","transform":{"lenses":[{"label":"filter_usdt","path":"assets/lens_filter_usdt.wasm","arguments":{"src":"address","value":"0xdac17f958d2ee523a2206206994597c13d831ec7"}},{"label":"decode_inputs","path":"assets/lens_decode_inputs.wasm","arguments":{"abi":"{\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}]}"}}]},"metadata":{"_v":2,"_t":2,"revisions":[{"version":0,"timestamp":"","diff":"{\"name\":\"example\",\"query\":null,\"sdl\":null,\"transform\":{\"lenses\":[]},\"metadata\":{\"_v\":0,\"_t\":0,\"revisions\":[],\"createdAt\":\"1749746283\",\"updatedAt\":\"1749746283\"}}"}],"createdAt":"1749746283","updatedAt":"1749748820"}}`

	if string(jsonBytes) != expected {
		t.Errorf("Unexpected Json Mismatch \n Expected:\n %s \n Got: \n %s", expected, string(jsonBytes))
	}
}

func TestViewJSONUnmarshaling(t *testing.T) {
	jsonString := `{
		"name": "example",
		"query": "Log {address topics data transactionHash blockNumber}",
		"sdl": "type FilteredAndDecodedLogs @materialized(if: false) {hash: String block: String address: String signature: String }",
		"transform": {
			"lenses": [
			{
				"label": "filter_usdt",
				"path": "assets/lens_filter_usdt.wasm",
				"arguments": {
				"src": "address",
				"value": "0xdac17f958d2ee523a2206206994597c13d831ec7"
				}
			},
			{
				"label": "decode_inputs",
				"path": "assets/lens_decode_inputs.wasm",
				"arguments": {
				"abi": "{\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}]}"
				}
			}
			]
		},
		"metadata": {
			"_v": 2,
			"_t": 2,
			"revisions": [
			{
				"diff": "{\"name\":\"example\",\"query\":null,\"sdl\":null,\"transform\":{\"lenses\":[]},\"metadata\":{\"_v\":0,\"_t\":0,\"revisions\":[],\"createdAt\":\"1749746283\",\"updatedAt\":\"1749746283\"}}"
			}
			],
			"createdAt": "1749746283",
			"updatedAt": "1749748820"
		}
	}
	`

	var view models.View
	if err := json.Unmarshal([]byte(jsonString), &view); err != nil {
		t.Fatalf("Failed to Unmarshal JSON: %v", err)
	}

	if view.Name != "example" {
		t.Errorf("Name Mismatch Error")
	}
}

func TestViewAllowsNullQueryAndSdlUnmarshalling(t *testing.T) {
	jsonString := `{
		"name": "the-good-view",
		"query": null,
		"sdl": null,
		"transform": {
			"lenses": []
		},
		"metadata": {
			"_v": 0,
			"_t": 0,
			"revisions": [],
			"createdAt": "1749746283",
			"updatedAt": "1749746283"
		}
	}
	`
	var view models.View

	if err := json.Unmarshal([]byte(jsonString), &view); err != nil {
		t.Fatalf("Unexpected Error Unmarshalling JSON: %v", err)
	}

	if view.Query != nil {
		t.Errorf("Failed to nullify query")
	}

	if view.Sdl != nil {
		t.Errorf("Failed to nullify sdl")
	}
}

func TestViewAllowsNullQueryAndSdlMarshalling(t *testing.T) {
	view := models.View{
		Name:  "the-good-view",
		Query: nil,
		Sdl:   nil,
		Transform: models.Transform{
			Lenses: []models.Lens{},
		},
		Metadata: models.Metadata{
			Version:   0,
			Total:     0,
			Revisions: []models.Revision{},
			CreatedAt: "1749746283",
			UpdatedAt: "1749746283",
		},
	}

	jsonBytes, err := json.Marshal(view)
	if err != nil {
		t.Fatalf("Error Marshalling JSON: %v", err)
	}

	expected := `{"name":"the-good-view","query":null,"sdl":null,"transform":{"lenses":[]},"metadata":{"_v":0,"_t":0,"revisions":[],"createdAt":"1749746283","updatedAt":"1749746283"}}`

	if string(jsonBytes) != expected {
		t.Errorf("Unexpected MisMatch \nExpected: \n%s \nGot: \n%s \n", expected, string(jsonBytes))
	}
}
