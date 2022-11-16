package gqlMgmt

//go:generate go run github.com/Khan/genqlient

import (
	"context"
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"net/http"
	"time"
)

type authedTransport struct {
	key     string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Auth-Token", t.key)
	return t.wrapped.RoundTrip(req)
}

func IndexToGraph(objUrl string, bucket string, cameraName string) {
	GraphqlEndpoint := "https://blue-surf-640094.us-east-1.aws.cloud.dgraph.io/graphql"

	key := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjJVQjY1NzZXLVhmM0VvWjJNb2I1aCJ9.eyJodHRwczovL2RncmFwaC5pby9qd3QvY2xhaW1zIjp7IlVTRVIiOiJuQGdyZWVuY3ViZXMuaW8ifSwiZ2l2ZW5fbmFtZSI6Ik5hdGhhbiIsImZhbWlseV9uYW1lIjoiQXJnZXRzaW5nZXIiLCJuaWNrbmFtZSI6Im4iLCJuYW1lIjoiTmF0aGFuIEFyZ2V0c2luZ2VyIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hL0FMbTV3dTNzYnBCRm5iS1Z2V29xQkZYQzlBeFRRVThNZjRpSUdwbnhSTzczPXM5Ni1jIiwibG9jYWxlIjoiZW4iLCJ1cGRhdGVkX2F0IjoiMjAyMi0xMC0xOFQwMDoxMToyMS45MzhaIiwiZW1haWwiOiJuQGdyZWVuY3ViZXMuaW8iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiaXNzIjoiaHR0cHM6Ly9kZXYtdWk1aDZ3bXAudXMuYXV0aDAuY29tLyIsInN1YiI6Imdvb2dsZS1vYXV0aDJ8MTAwNDMwMTUzNTQ4MjIyOTU1MzI5IiwiYXVkIjoiSk1hTzRXNlNVVHlTaGxla2c4YlFCMzV5a1N0dHFiaXUiLCJpYXQiOjE2NjYwNTE4ODIsImV4cCI6MTY2NjA4Nzg4Miwic2lkIjoiTDhoMlV6Nk9UVWctV2ROWEdmUVN3Q0FfdGJ4MFRLUVUiLCJub25jZSI6Ill6ZCtZWFIwTFc1YVpFZHRWRkkzV2tWWk1rZHZWbmx2TFVsRVpXdEhaa2RHVFdVNFgySnVNRU5HVGc9PSJ9.YP1gzvwDSWm1_BNaN_LKWn4miaFE8hcA7y87fdSqAasyKONckZ1rKSM8g0BsL_c7N5biSof4d9meA8RJUUY7xKTluHPMqN2971ufEOBsx-jeSymgcli9VyVBtidNiL7_btjZcjlTCRPawbHs1GaYaQ9KyJs1JJZTi3CBHH4xeIU5uEV3W0jCCLMNXTx5N5a3r-Tr_m13DpoH0MznE_AFsEFTZFVPAhsu4ciRUQe39zjCscK6Y3xXOxJ6rd3ZP2ac-Ov9W1U6-chrgrfpHrKRwTuYqdeASEh4qN-uVmor91lHmVrk59v8B1E3yEjwJILQ9vmC4KdQX-GEZyRFV8rmfQ"

	httpClient := http.Client{
		Transport: &authedTransport{
			key:     key,
			wrapped: http.DefaultTransport,
		},
	}

	ctx := context.Background()
	//client := graphql.NewClient(GraphqlEndpoint, http.DefaultClient)
	client := graphql.NewClient(GraphqlEndpoint, &httpClient)
	resp, err := MyQuery(ctx, client)
	fmt.Println(resp.QueryUser, err)

	//mutresp, err := backend_graphql_client.AddPosts(ctx, client

	iso_8601_fmt := fmt.Sprint(time.Now().Format(time.RFC3339))
	println(iso_8601_fmt)

	println("now using simply... time.Now()")

	stillResp, err := IndexImage(ctx, client, objUrl, bucket, cameraName, time.Now())
	fmt.Println(stillResp)
}
