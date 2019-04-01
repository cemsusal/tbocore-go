package core

// Import net/http and elastic
import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"golang.org/x/sync/errgroup"
	"io"
	"reflect"
	"strings"
)

// ElasticSearch helps you to handle elasticsearch requests
type ElasticSearch struct {
	Client       *elastic.Client
	Config       *Config
	Aggregations []Aggregation
}

type Aggregation struct {
	Agg        string
	SubAgg     string
	Name       string
	SubAggName string
}

// InitClient initiates an instance of ElasticSearch type
func (e *ElasticSearch) InitClient(uri []string) {
	opts := e.clientOptions()
	client, err := elastic.NewClient(opts...)

	if err != nil {
		// Handle error
		panic(err)
	}
	e.Client = client
}

func (e *ElasticSearch) clientOptions() []elastic.ClientOptionFunc {
	confFunctions := []elastic.ClientOptionFunc{}
	elasticConf := e.Config.Elastic

	if elasticConf.Uris != nil {
		confFunctions = append(confFunctions, elastic.SetURL(strings.Join(elasticConf.Uris, ", ")))
	}

	if len(elasticConf.Username) > 0 && len(elasticConf.Password) > 0 {
		confFunctions = append(confFunctions, elastic.SetBasicAuth(elasticConf.Username, elasticConf.Password))
	}

	return confFunctions
}

// DeleteIndex deletes an existing index by its name
func (e *ElasticSearch) DeleteIndex(index string) (bool, error) {
	// Delete an index
	ctx := context.Background()
	deleteIndex, err := e.Client.DeleteIndex(index).Do(ctx)
	if err != nil {
		return false, err
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
	}
	return true, err
}

// IndexExists checks if index exists or not
func (e *ElasticSearch) IndexExists(index string) {
	exists, err := e.Client.IndexExists(index).Do(context.Background())
	if err != nil {
		// Handle error
	}
	if !exists {
		// Index does not exist yet.
	}
}

// Ping returns health information for an elastic node
func (e *ElasticSearch) Ping(uri string) (*elastic.PingResult, int, error) {
	ctx := context.Background()
	info, code, err := e.Client.Ping(uri).Do(ctx)
	if err != nil {
		panic(err)
	}
	return info, code, err
}

// GetVersion retursn the elasticsearch version for a node
func (e *ElasticSearch) GetVersion(uri string) (string, error) {
	version, err := e.Client.ElasticsearchVersion(uri)
	if err != nil {
		panic(err)
	}
	return version, err
}

// Indices returns the index names for a client
func (e *ElasticSearch) Indices(uri string) ([]string, error) {
	names, err := e.Client.IndexNames()
	return names, err
}

// Insert creates a document
func (e *ElasticSearch) Insert(doc interface{}) (*elastic.IndexResponse, error) {
	// Index a tweet (using JSON serialization)
	ctx := context.Background()
	t := reflect.TypeOf(doc).String()
	document, err := e.Client.Index().
		Index(t).
		Type(t).
		Id("1").
		BodyJson(t).
		Do(ctx)
	if err != nil {
		return document, err
	}
	return document, err
}

// GetByID returns a document for an ID
func (e *ElasticSearch) GetByID(id string, index string) (interface{}, error) {
	ctx := context.Background()
	doc, err := e.Client.Get().
		Index(index).
		Type(index).
		Id(id).
		Do(ctx)
	return doc, err
}

// Delete removes a document by ID from an index
func (e *ElasticSearch) Delete(id string, index string) (*elastic.DeleteResponse, error) {
	ctx := context.Background()
	res, err := e.Client.Delete().
		Index(index).
		Type(index).
		Id(id).
		Do(ctx)
	return res, err
}

// CreateIndex creates an index (DB) on elastic client
func (e *ElasticSearch) CreateIndex(mapping string, indexName string) (bool, error) {
	ctx := context.Background()
	exists, err := e.Client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return false, err
	}
	if exists {
		_, err = e.Client.DeleteIndex(indexName).Do(ctx)
		if err != nil {
			return false, err
		}
	}
	// Create index with mapping
	_, err = e.Client.CreateIndex(indexName).Body(mapping).Do(ctx)

	createIndex, err := e.Client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
	return createIndex.Acknowledged, err
}

// TermQuery performs a term query on elastic client
func (e *ElasticSearch) TermQuery(fieldName string, query string, indexName string, sortBy string, from int, size int) (*elastic.SearchResult, error) {
	ctx := context.Background()
	termQuery := elastic.NewTermQuery(fieldName, query)
	searchResult, err := e.Client.Search().
		Index(indexName).
		Query(termQuery).
		Sort(sortBy, true).
		From(from).Size(size).
		Pretty(true).
		Do(ctx)
	return searchResult, err
}

func (e *ElasticSearch) Search(indexName string, theType string, query elastic.Query, aggregation elastic.Aggregation, sort []string, from int, size int) {
	// Create service and use query, aggregations, sort, filter, pagination funcs
	search := e.Client.Search().Index(indexName).Type(theType).Pretty(true)
	search = search.Query(query)
	search = e.aggregate(search)
	search = e.sort(search, sort)
	search = e.paginate(search, from, size)

}

// aggs sets up the aggregations in the service.
func (e *ElasticSearch) aggregate(service *elastic.SearchService) *elastic.SearchService {
	for _, agg := range e.Aggregations {
		aggregation := elastic.NewTermsAggregation().Field(agg.Agg)
		if len(agg.SubAgg) > 0 && len(agg.SubAggName) > 0 {
			subAggregation := elastic.NewTermsAggregation().Field(agg.SubAgg)
			aggregation.SubAggregation(agg.SubAggName, subAggregation)
		}
		service = service.Aggregation(agg.Name, aggregation)
	}
	return service
}

// paginate sets up pagination in the service.
func (e *ElasticSearch) paginate(service *elastic.SearchService, from int, size int) *elastic.SearchService {
	if from > 0 {
		service = service.From(from)
	}
	if size > 0 {
		service = service.Size(size)
	}
	return service
}

// sorting applies sort to the service.
func (e *ElasticSearch) sort(service *elastic.SearchService, sort []string) *elastic.SearchService {
	if len(sort) == 0 {
		// Sort by score by default
		service = service.Sort("_score", false)
		return service
	}

	// Sort by fields; prefix of "-" means: descending sort order.
	for _, s := range sort {
		s = strings.TrimSpace(s)

		var field string
		var asc bool

		if strings.HasPrefix(s, "-") {
			field = s[1:]
			asc = false
		} else {
			field = s
			asc = true
		}

		// Maybe check for permitted fields to sort

		service = service.Sort(field, asc)
	}
	return service
}

func (e *ElasticSearch) Scroll(indexName string, theType interface{}, size int) {
	// This example illustrates how to use goroutines to iterate
	// through a result set via ScrollService.
	//
	// It uses the excellent golang.org/x/sync/errgroup package to do so.
	//
	// The first goroutine will Scroll through the result set and send
	// individual documents to a channel.
	//
	// The second cluster of goroutines will receive documents from the channel and
	// deserialize them.
	//
	// Feel free to add a third goroutine to do something with the
	// deserialized results.
	//
	// Let's go.

	// 1st goroutine sends individual hits to channel.
	t := reflect.TypeOf(theType).String()
	hits := make(chan json.RawMessage)
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		defer close(hits)
		// Initialize scroller. Just don't call Do yet.
		scroll := e.Client.Scroll(indexName).Type(t).Size(size)
		for {
			results, err := scroll.Do(ctx)
			if err == io.EOF {
				return nil // all results retrieved
			}
			if err != nil {
				return err // something went wrong
			}

			// Send the hits to the hits channel
			for _, hit := range results.Hits.Hits {
				select {
				case hits <- *hit.Source:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	})

	// 2nd goroutine receives hits and deserializes them.
	//
	// If you want, setup a number of goroutines handling deserialization in parallel.
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			for hit := range hits {
				// Deserialize
				var arbitrary_json map[string]theType
				json.Unmarshal([]byte(hit), &arbitrary_json)
				if err != nil {
					return err
				}

				// Do something with the product here, e.g. send it to another channel
				// for further processing.
				_ = tp

				// Terminate early?
				select {
				default:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})
	}

	// Check whether any goroutines failed.
	if err := g.Wait(); err != nil {
		panic(err)
	}
}
