package section

import "github.com/AlexBond702/catalog-service/internal/app/util"

type (
	Repository struct {
		Postgres RepositoryPostgres
	}
	RepositoryPostgres struct {
		Address      string        `required:"true"`
		Username     string        `required:"true"`
		Password     string        `required:"true"`
		Name         string        `required:"true"`
		ReadTimeout  util.Duration `split_words:"true" required:"true"`
		WriteTimeout util.Duration `split_words:"true" required:"true"`
	}
)
