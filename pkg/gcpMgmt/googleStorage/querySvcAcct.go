package googleStorage

import (
	"context"
	"fmt"
	"io"

	iam "google.golang.org/api/iam/v1"
)

// listServiceAccounts lists a project's service accounts.
func listServiceAccounts(w io.Writer, projectID string) ([]*iam.ServiceAccount, error) {
	ctx := context.Background()
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %v", err)
	}

	response, err := service.Projects.ServiceAccounts.List("projects/" + projectID).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.List: %v", err)
	}
	for _, account := range response.Accounts {
		fmt.Fprintf(w, "Listing service account: %v\n", account.Name)
	}
	return response.Accounts, nil
}
