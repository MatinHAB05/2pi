package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MatinHAB05/2pi/internal/domain/exception"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/internal/infrastructure/database"
	dsl "github.com/elastic/go-elasticsearch/v9/typedapi/esdsl"
)

type articleESTypedRepository struct {
	client    database.Elastic
	aliasName string
}

func NewArticleESTypedRepository(client database.Elastic) repository_contract.ArticleESRepository {
	return &articleESTypedRepository{
		client:    client,
		aliasName: "articles",
	}
}

func (r *articleESTypedRepository) Index(ctx context.Context, article *repository_contract.ArticleDocument) (*string, error) {
	res, err := r.client.GetClient().Index(r.aliasName).
		Id(article.ID).
		Document(article).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w [id=%s]: %v", exception.ErrIndexFailed, article.ID, err)
	}

	res_name := res.Result.Name
	return &res_name, nil
}

func (r *articleESTypedRepository) Delete(ctx context.Context, id string) (*string, error) {
	res, err := r.client.GetClient().Delete(r.aliasName, id).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("%w [id=%s]: %v", exception.ErrDeleteFailed, id, err)
	}

	res_name := res.Result.Name
	return &res_name, nil
}

// serach after : is better way to do pagiantion(speccilaly for more than 10,000 articles) but the telggram-client can't pass search_after parameter...
func (r *articleESTypedRepository) Search(ctx context.Context, keyword string, page, size int) (*repository_contract.ArticleSearchResult, error) {
	if page < 1 {
		page = 1
	}
	from := (page - 1) * size

	boost3 := float32(3.0)
	boost1_5 := float32(1.5)
	boost1 := float32(1.0)
	boost0_5 := float32(0.5)
	minShouldMatch := 1
	fuzziness := "AUTO"

	res, err := r.client.GetClient().Search().
		Index(r.aliasName).
		From(from).
		Size(size).
		Query(dsl.NewQuery().Bool(dsl.NewBoolQuery().MinimumShouldMatch(dsl.NewMinimumShouldMatch().Int(minShouldMatch)).
			Should(
				dsl.NewMatchQuery("title", keyword).Boost(boost3),
				dsl.NewMatchQuery("description", keyword).Boost(boost1_5),
				dsl.NewMatchQuery("title.fuzzy", keyword).Fuzziness(dsl.NewFuzziness().String(fuzziness)).Boost(boost1),
				dsl.NewMatchQuery("description.fuzzy", keyword).Fuzziness(dsl.NewFuzziness().String(fuzziness)).Boost(boost0_5),
			),
		)).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrSearchFailed, err)
	}

	articles := make([]*repository_contract.ArticleDocument, 0, min(len(res.Hits.Hits), size))
	for i := range res.Hits.Hits {
		var doc repository_contract.ArticleDocument
		if err := json.Unmarshal(res.Hits.Hits[i].Source_, &doc); err != nil {
			return nil, fmt.Errorf("%w [index=%d]: %v", exception.ErrUnmarshalFailed, i, err)
		}
		articles = append(articles, &doc)
	}

	return &repository_contract.ArticleSearchResult{
		Articles: articles,
		Total:    len(res.Hits.Hits),
	}, nil
}
