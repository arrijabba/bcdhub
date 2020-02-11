package main

import (
	"github.com/aopoltorzhicky/bcdhub/internal/contractparser/consts"
	"github.com/aopoltorzhicky/bcdhub/internal/contractparser/storage"
	"github.com/aopoltorzhicky/bcdhub/internal/elastic"
	"github.com/aopoltorzhicky/bcdhub/internal/noderpc"
	"github.com/tidwall/gjson"
)

func getRichStorage(es *elastic.Elastic, rpc *noderpc.NodeRPC, op gjson.Result, level int64, protocol, operationID string) (storage.RichStorage, error) {
	kind := op.Get("kind").String()

	switch protocol {
	case consts.HashBabylon, consts.HashCarthage:
		parser := storage.NewBabylon(es, rpc)
		switch kind {
		case consts.Transaction:
			return parser.ParseTransaction(op, protocol, level, operationID)
		case consts.Origination:
			return parser.ParseOrigination(op, protocol, level, operationID)
		}
	default:
		parser := storage.NewAlpha(es)
		switch kind {
		case consts.Transaction:
			return parser.ParseTransaction(op, protocol, level, operationID)
		case consts.Origination:
			return parser.ParseOrigination(op, protocol, level, operationID)
		}
	}
	return storage.RichStorage{Empty: true}, nil
}