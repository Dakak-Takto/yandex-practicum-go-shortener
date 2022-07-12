package handlers

type (
	userURL_DTO struct {
		Original string `json:"original_url"`
		Short    string `json:"short_url"`
	}

	makeShortRequestDTO struct {
		URL string `json:"url"`
	}

	resultDTO struct {
		Result string `json:"result"`
	}

	batchItemRequestURLs struct {
		CorellationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	batchItemResponse struct {
		CorellationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
)
