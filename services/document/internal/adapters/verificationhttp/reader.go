package verificationhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"idenova/document/internal/domain"
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

func (r *Reader) RequirePending(ctx context.Context, verificationID string) error {
	endpoint := r.baseURL + "/verifications/" + url.PathEscape(verificationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("verification lookup: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return domain.ErrVerificationNotFound
	case http.StatusOK:
		var verification struct {
			Status string `json:"status"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&verification); err != nil {
			return fmt.Errorf("decode verification: %w", err)
		}

		if verification.Status != "PENDING" {
			return domain.ErrVerificationNotPending
		}
		return nil
	default:
		return fmt.Errorf("verification lookup returned HTTP %d", resp.StatusCode)
	}
}
