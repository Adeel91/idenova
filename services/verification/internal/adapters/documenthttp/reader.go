package documenthttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"idenova/verification/internal/ports"
)

type Reader struct {
	baseURL string
	client  *http.Client
}

func NewReader(baseURL string, client ...*http.Client) *Reader {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	if len(client) > 0 && client[0] != nil {
		httpClient = client[0]
	}

	return &Reader{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  httpClient,
	}
}

func (r *Reader) ListByVerificationID(
	ctx context.Context,
	verificationID string,
) ([]ports.Document, error) {
	endpoint := fmt.Sprintf(
		"%s/verifications/%s/documents",
		r.baseURL,
		url.PathEscape(verificationID),
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("document service request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"document service returned HTTP %d",
			resp.StatusCode,
		)
	}

	var records []struct {
		ID     string               `json:"ID"`
		Status ports.DocumentStatus `json:"Status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return nil, fmt.Errorf("decode documents: %w", err)
	}

	documents := make([]ports.Document, 0, len(records))

	for _, record := range records {
		documents = append(documents, ports.Document{
			ID:     record.ID,
			Status: record.Status,
		})
	}

	return documents, nil
}
