package brpc

import (
	"context"
	"net/http"
	"testing"
)

func TestBrpcServiceClient_Do(t *testing.T) {
	ctx := context.Background()
	url := "http://127.0.0.1:8081/recall/get"
	http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
}
