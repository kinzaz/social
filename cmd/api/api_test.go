package main

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	ratelimiter "github.com/kinzaz/social/internal/rateLimiter"
// )

// func TestRatelimiterMiddleware(t *testing.T) {
// 	cfg := config{
// 		rateLimiter: ratelimiter.Config{
// 			RequestsPerTimeFrame: 20,
// 			TimeFrame:            time.Second * 5,
// 			Enabled:              true,
// 		},
// 		addr: ":8080",
// 	}

// 	app := newTestApplication(t, cfg)
// 	ts := httptest.NewServer(app.mount())
// 	defer ts.Close()

// 	client := &http.Client{}
// 	mockIP := "192.168.1.1"
// 	marginOfError := 2

// 	for i := 0; i < cfg.rateLimiter.ReadHeaderTimeout+marginOfError; i++ {
// 		req, err := http.NewRequest("GET", ts.URL+"/v1/health", nil)
// 		if err != nil {
// 			t.Fatalf("could not create request: %v", err)
// 		}

// 		req.Header.Set("X-Forwarded-For", mockIP)
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			t.Fatalf("could not send request: %v", err)
// 		}
// 		defer resp.Body.Close()

// 		if i < cfg.rateLimiter.RequestsPerTimeFrame {
// 			if resp.StatusCode != http.StatusOK {
// 				t.Errorf("expected status OK, but got %v", resp.Status)
// 			}
// 		} else {
// 			if resp.StatusCode != http.StatusTooManyRequests {
// 				t.Errorf("expected status too many request, but got %v", resp.Status)
// 			}
// 		}

// 	}

// }
