package mysql

import (
	"database/sql"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"strconv"
)

type ArticleModel struct {
	DB *sql.DB
}

func (m *ArticleModel) SelectByCategory(categoryID string) ([]models.Article, error) {
	id, _ := strconv.Atoi(categoryID)

	articles := []models.Article{
		models.Article{
			1, 1, id, 1, "Maru wins yet another GSL", "John Doe", "2019-04-01", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur posuere elit et ligula cursus, sed fermentum erat faucibus. Morbi sollicitudin et quam sagittis mattis. Donec pellentesque, lorem id congue interdum, lorem tellus gravida diam, in mattis lorem quam ac erat. Sed enim est, condimentum eu pretium eget, pretium eu massa. Cras vestibulum porttitor tellus, vel ullamcorper dui scelerisque ultrices. Aliquam ultrices justo ligula, quis elementum dui interdum non. Nullam sed viverra neque. In hac habitasse platea dictumst. Vestibulum vel justo non erat condimentum ornare ut ac ante. Donec condimentum lobortis risus in facilisis. Aliquam imperdiet tellus ut lectus varius, non tincidunt ipsum vehicula. Vestibulum nec arcu non leo dapibus semper id vel risus.\n\nPhasellus elementum aliquet nisi, eu fermentum diam commodo sit amet. Etiam eu tristique risus. Suspendisse condimentum finibus congue. Etiam porttitor lorem id massa auctor aliquet. Pellentesque congue porta purus at rhoncus. Suspendisse nec ipsum et sem dignissim finibus ac at est. Donec sed ex odio. Duis quis interdum metus. Morbi ac ultricies neque. Mauris aliquam velit nec ligula mollis, et condimentum est dapibus. Mauris tincidunt odio at malesuada mattis. Praesent porta iaculis lectus, non ultrices risus vehicula sed. Donec quis dolor felis. Ut lacus tellus, suscipit nec mauris vel, laoreet iaculis libero. Pellentesque a nunc quis ligula tristique dapibus. Nunc velit diam, congue vitae feugiat vitae, congue ac elit.", "https://deepmind.com/blog/deepmind-and-blizzard-open-starcraft-ii-ai-research-environment/",
		},
		models.Article{
			2, 1, id, 1, "Koreans continue to dominate", "A sad foreigner", "2019-03-25", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur posuere elit et ligula cursus, sed fermentum erat faucibus. Morbi sollicitudin et quam sagittis mattis. Donec pellentesque, lorem id congue interdum, lorem tellus gravida diam, in mattis lorem quam ac erat. Sed enim est, condimentum eu pretium eget, pretium eu massa. Cras vestibulum porttitor tellus, vel ullamcorper dui scelerisque ultrices. Aliquam ultrices justo ligula, quis elementum dui interdum non. Nullam sed viverra neque. In hac habitasse platea dictumst. Vestibulum vel justo non erat condimentum ornare ut ac ante. Donec condimentum lobortis risus in facilisis. Aliquam imperdiet tellus ut lectus varius, non tincidunt ipsum vehicula. Vestibulum nec arcu non leo dapibus semper id vel risus.\n\nPhasellus elementum aliquet nisi, eu fermentum diam commodo sit amet. Etiam eu tristique risus. Suspendisse condimentum finibus congue. Etiam porttitor lorem id massa auctor aliquet. Pellentesque congue porta purus at rhoncus. Suspendisse nec ipsum et sem dignissim finibus ac at est. Donec sed ex odio. Duis quis interdum metus. Morbi ac ultricies neque. Mauris aliquam velit nec ligula mollis, et condimentum est dapibus. Mauris tincidunt odio at malesuada mattis. Praesent porta iaculis lectus, non ultrices risus vehicula sed. Donec quis dolor felis. Ut lacus tellus, suscipit nec mauris vel, laoreet iaculis libero. Pellentesque a nunc quis ligula tristique dapibus. Nunc velit diam, congue vitae feugiat vitae, congue ac elit.", "https://deepmind.com/blog/deepmind-and-blizzard-open-starcraft-ii-ai-research-environment/",
		},
		models.Article{
			3, 1, id, 1, "What does it take to get to silver league in Starcraft 2", "A platinum veteran", "2019-04-01", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur posuere elit et ligula cursus, sed fermentum erat faucibus. Morbi sollicitudin et quam sagittis mattis. Donec pellentesque, lorem id congue interdum, lorem tellus gravida diam, in mattis lorem quam ac erat. Sed enim est, condimentum eu pretium eget, pretium eu massa. Cras vestibulum porttitor tellus, vel ullamcorper dui scelerisque ultrices. Aliquam ultrices justo ligula, quis elementum dui interdum non. Nullam sed viverra neque. In hac habitasse platea dictumst. Vestibulum vel justo non erat condimentum ornare ut ac ante. Donec condimentum lobortis risus in facilisis. Aliquam imperdiet tellus ut lectus varius, non tincidunt ipsum vehicula. Vestibulum nec arcu non leo dapibus semper id vel risus.\n\nPhasellus elementum aliquet nisi, eu fermentum diam commodo sit amet. Etiam eu tristique risus. Suspendisse condimentum finibus congue. Etiam porttitor lorem id massa auctor aliquet. Pellentesque congue porta purus at rhoncus. Suspendisse nec ipsum et sem dignissim finibus ac at est. Donec sed ex odio. Duis quis interdum metus. Morbi ac ultricies neque. Mauris aliquam velit nec ligula mollis, et condimentum est dapibus. Mauris tincidunt odio at malesuada mattis. Praesent porta iaculis lectus, non ultrices risus vehicula sed. Donec quis dolor felis. Ut lacus tellus, suscipit nec mauris vel, laoreet iaculis libero. Pellentesque a nunc quis ligula tristique dapibus. Nunc velit diam, congue vitae feugiat vitae, congue ac elit.", "https://deepmind.com/blog/deepmind-and-blizzard-open-starcraft-ii-ai-research-environment/",
		},
	}

	return articles, nil
}
