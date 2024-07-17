package es

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"eyes/internal/domain"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const (
	IndexName string = "monitor"
)

type MonitorESDAO interface {
	Count(ctx context.Context, info domain.Monitor) (int64, error)
	Detail(ctx context.Context, info domain.Monitor) (*domain.Monitor, error)
	List(ctx context.Context, page, size int64, info domain.Monitor) ([]*domain.Monitor, int64, error)
	Query(ctx context.Context, page, size int64, args string, info domain.Monitor) ([]*domain.Monitor, int64, error)
}

var _ MonitorESDAO = new(monitorDAO)

func NewMonitorDAO(es *elasticsearch.Client) MonitorESDAO {
	return &monitorDAO{
		es: es,
	}
}

type monitorDAO struct {
	es *elasticsearch.Client
}

func (m *monitorDAO) Count(ctx context.Context, info domain.Monitor) (int64, error) {
	var (
		buf []byte
		ret map[string]any
		cnt int64
	)

	resp, err := m.es.Count(
		func(request *esapi.CountRequest) {
			request.Index = []string{IndexName}
		},
	)
	if err != nil {
		return 0, fmt.Errorf("es.Count, err; %v", err)
	}

	if resp.StatusCode == http.StatusOK {

		buf, err = io.ReadAll(resp.Body)
		if err != nil {
			return 0, fmt.Errorf("io.ReadAll: %v", err)
		}

		err = json.Unmarshal(buf, &ret)
		if err != nil {
			return 0, fmt.Errorf("json.unmarshal: %v", err)
		}

		value, ok := ret["count"].(float64)
		if ok {
			cnt = int64(value)
		} else {
			return 0, fmt.Errorf("int64(value): %v", ok)
		}

		return cnt, nil

	} else {
		return 0, fmt.Errorf("resp: %v", resp)
	}
}

func (m *monitorDAO) Detail(ctx context.Context, info domain.Monitor) (*domain.Monitor, error) {
	ret := new(domain.Monitor)

	resp, err := m.es.GetSource(IndexName, info.ID)
	if err != nil {
		return nil, fmt.Errorf("m.es.GetSource(IndexName, info.ID): %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("resp: %v", resp)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll(resp.Body): %v", err)
	}
	err = json.Unmarshal(buf, ret)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(buf, ret): %v", err)
	}
	ret.ID = info.ID
	return ret, nil
}

func (m *monitorDAO) List(ctx context.Context, page, size int64, info domain.Monitor) ([]*domain.Monitor, int64, error) {
	resp, err := m.es.Search(
		m.es.Search.WithIndex(IndexName),
		m.es.Search.WithPretty(),
		m.es.Search.WithBody(strings.NewReader(
			fmt.Sprintf(`
				{
					"from": %d,
					"size": %d
				}`, (page-1)*size, size)),
		),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("es_bak.Search: %v", err)
	}

	return _parseResponse(resp)
}

func (m *monitorDAO) Query(ctx context.Context, page, size int64, args string, info domain.Monitor) ([]*domain.Monitor, int64, error) {
	resp, err := m.es.Search(
		m.es.Search.WithIndex(IndexName),
		m.es.Search.WithPretty(),
		m.es.Search.WithBody(strings.NewReader(args)),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("m.es.Search(: %v", err)
	}

	return _parseResponse(resp)
}

func _parseResponse(resp *esapi.Response) ([]*domain.Monitor, int64, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("resp: %v", resp)
	}

	var respMap map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&respMap); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response body: %v", err)
	}

	hits, ok := respMap["hits"].(map[string]any)
	if !ok {
		return nil, 0, fmt.Errorf("respMap[\"hits\"].(map[string]any) %+v", respMap)
	}

	totalMap, ok := hits["total"].(map[string]any)
	if !ok {
		return nil, 0, fmt.Errorf("hits[\"total\"].(map[string]any)%+v", respMap)
	}
	fmt.Println("totalMap::", totalMap)

	total := int64(totalMap["value"].(float64))

	hitsArray, ok := hits["hits"].([]any)
	if !ok {
		return nil, 0, fmt.Errorf("invalid 'hits' field type in response: %+v", hitsArray)
	}

	ret := make([]*domain.Monitor, len(hitsArray))
	for i, hit := range hitsArray {

		hitMap, ok := hit.(map[string]any)
		if !ok {
			return nil, 0, fmt.Errorf("invalid hit structure in JSON response: %v", hitMap)
		}

		_ID, ok := hitMap["_id"].(string)
		if !ok {
			return nil, 0, fmt.Errorf("invalid '_id' field type in hit: %v", hitMap["_id"])
		}

		source, ok := hitMap["_source"].(map[string]any)
		if !ok {
			return nil, 0, fmt.Errorf("invalid '_source' field in hit: %v", source)
		}

		b, err := json.Marshal(source)
		if err != nil {
			return nil, 0, fmt.Errorf("json.Marshal(source): %v", err)
		}

		tmp := domain.Monitor{}
		err = json.Unmarshal(b, &tmp)
		if err != nil {
			return nil, 0, fmt.Errorf("json.Unmarshal(b, &tmp): %v", err)
		}

		//var tmp models.Monitor
		//if err := mapstructure.Decode(source, &tmp); err != nil {
		//	return nil, fmt.Errorf("failed to decode '_source' field: %v", err)
		//}

		tmp.ID = _ID
		ret[i] = &tmp
	}
	return ret, total, nil
}
