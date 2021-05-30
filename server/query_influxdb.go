package server

import (
	"context"
	"errors"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/query"
)

type Querier struct {
	c influxdb.Client
	q api.QueryAPI
}

func NewQuerier(serverURL, authToken, org string) *Querier {
	client := influxdb.NewClient(serverURL, authToken)
	if client == nil {
		return nil
	}
	queryAPI := client.QueryAPI(org)
	if queryAPI == nil {
		return nil
	}

	return &Querier{
		c: client,
		q: queryAPI,
	}
}

func (q *Querier) Close() error {
	if q.c != nil {
		q.c.Close()
	}

	return nil
}

func (q *Querier) ExecuteQuery(ctx context.Context, queryToRun string) ([]*query.FluxRecord, error) {
	if q.q == nil {
		return nil, errors.New("query api not available")
	}
	result, err := q.q.Query(ctx, queryToRun)
	if err != nil {
		return nil, err
	}

	var records []*query.FluxRecord
	for result.Next() {
		records = append(records, result.Record())
	}
	if result.Err() != nil {
		return records, result.Err()
	}
	return records, nil
}
