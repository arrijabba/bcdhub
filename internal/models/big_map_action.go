package models

import (
	"time"

	"github.com/tidwall/gjson"
)

// BigMapAction -
type BigMapAction struct {
	ID             string    `json:"-"`
	Action         string    `json:"action"`
	SourcePtr      *int64    `json:"source_ptr,omitempty"`
	DestinationPtr *int64    `json:"destination_ptr,omitempty"`
	OperationID    string    `json:"operation_id"`
	Level          int64     `json:"level"`
	Address        string    `json:"address"`
	Network        string    `json:"network"`
	IndexedTime    int64     `json:"indexed_time"`
	Timestamp      time.Time `json:"timestamp"`
}

// GetID -
func (b *BigMapAction) GetID() string {
	return b.ID
}

// GetIndex -
func (b *BigMapAction) GetIndex() string {
	return "bigmapaction"
}

// ParseElasticJSON -
func (b *BigMapAction) ParseElasticJSON(hit gjson.Result) {
	b.ID = hit.Get("_id").String()
	b.Action = hit.Get("_source.action").String()

	if hit.Get("_source.source_ptr").Exists() {
		srcPtr := hit.Get("_source.source_ptr").Int()
		b.SourcePtr = &srcPtr
	}

	if hit.Get("_source.destination_ptr").Exists() {
		dstPtr := hit.Get("_source.destination_ptr").Int()
		b.DestinationPtr = &dstPtr
	}

	b.OperationID = hit.Get("_source.operation_id").String()
	b.Level = hit.Get("_source.level").Int()
	b.Address = hit.Get("_source.address").String()
	b.Network = hit.Get("_source.network").String()
	b.IndexedTime = hit.Get("_source.indexed_time").Int()
	b.Timestamp = hit.Get("_source.timestamp").Time().UTC()
}