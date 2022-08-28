//go:build integration
// +build integration

package form3_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/lmikolajczak/go-form3/form3"
	"testing"
)

func TestClient_CreateAccount(t *testing.T) {
	f3, teardown := form3.TestClient(t)
	defer teardown()

	testcases := []struct {
		name       string
		attributes *form3.AccountAttributes
		wantAttrs  *form3.AccountAttributes
		wantErr    error
	}{
		{
			name:       "required attributes only",
			attributes: accountAttributesRequired(t),
			wantAttrs: &form3.AccountAttributes{
				Country: form3.String("NL"),
				Name:    []string{"L. Mikolajczak"},
			},
		},
		{
			name:       "with optional attributes",
			attributes: accountAttributesOptional(t),
			wantAttrs: &form3.AccountAttributes{
				Country: form3.String("NL"),
				Name:    []string{"L. Mikolajczak"},
			},
		},
		{
			name: "without attributes",
			wantErr: &form3.F3Error{
				StatusCode:   400,
				ErrorCode:    0,
				ErrorMessage: "validation failure list:\nvalidation failure list:\nattributes in body is required",
			},
		},
		{
			name:       "parse validation errors",
			attributes: accountAttributesInvalid(t),
			wantErr: &form3.F3Error{
				StatusCode:   400,
				ErrorCode:    0,
				ErrorMessage: "validation failure list:\nvalidation failure list:\nvalidation failure list:\naccount_number in body should match '^[A-Z0-9]{0,64}$'",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			account, err := f3.CreateAccount(uuid.NewString(), tc.attributes)

			testErrorMessage(t, err, tc.wantErr)
			if account != nil {
				testAccountAttrs(t, account.Attributes, tc.wantAttrs)
			}
		})
	}
}

func TestClient_DeleteAccount(t *testing.T) {
	f3, teardown := form3.TestClient(t)
	defer teardown()

	organisationId := uuid.NewString()
	account, _ := f3.CreateAccount(organisationId, accountAttributesRequired(t))

	testcases := []struct {
		name           string
		accountID      string
		accountVersion int64
		organisationID string
		wantErr        error
	}{
		{
			name:           "account does not exist",
			accountID:      uuid.NewString(),
			accountVersion: 0,
			organisationID: organisationId,
			wantErr: &form3.F3Error{
				StatusCode:   404,
				ErrorCode:    0,
				ErrorMessage: "",
			},
		},
		{
			name:           "invalid version for existing account",
			accountID:      account.ID,
			accountVersion: 123,
			organisationID: organisationId,
			wantErr: &form3.F3Error{
				StatusCode:   409,
				ErrorCode:    0,
				ErrorMessage: "invalid version",
			},
		},
		{
			name:           "success",
			accountID:      account.ID,
			accountVersion: 0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := f3.DeleteAccount(tc.accountID, tc.accountVersion)

			testErrorMessage(t, err, tc.wantErr)
		})
	}
}

func TestClient_FetchAccount(t *testing.T) {
	f3, teardown := form3.TestClient(t)
	defer teardown()

	account, _ := f3.CreateAccount(uuid.NewString(), accountAttributesRequired(t))
	nonExistingAccountId := uuid.NewString()

	testcases := []struct {
		name      string
		accountID string
		wantAttrs *form3.AccountAttributes
		wantErr   error
	}{
		{
			name:      "account does not exist",
			accountID: nonExistingAccountId,
			wantErr: &form3.F3Error{
				StatusCode:   404,
				ErrorCode:    0,
				ErrorMessage: fmt.Sprintf("record %s does not exist", nonExistingAccountId),
			},
		},
		{
			name:      "success",
			accountID: account.ID,
			wantAttrs: &form3.AccountAttributes{
				Country: form3.String("NL"),
				Name:    []string{"L. Mikolajczak"},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			fetched, err := f3.FetchAccount(tc.accountID)

			testErrorMessage(t, err, tc.wantErr)
			if err == nil {
				testAccountId(t, fetched, tc.accountID)
				testAccountAttrs(t, fetched.Attributes, tc.wantAttrs)
			}
		})
	}
}

func accountAttributesRequired(t *testing.T) *form3.AccountAttributes {
	t.Helper()
	return &form3.AccountAttributes{
		Country: form3.String("NL"),
		Name:    []string{"L. Mikolajczak"},
	}
}

func accountAttributesOptional(t *testing.T) *form3.AccountAttributes {
	t.Helper()
	attributes := accountAttributesRequired(t)

	attributes.AccountNumber = "123654"
	attributes.BaseCurrency = "EUR"
	attributes.Bic = "INGBNL2A"

	return attributes
}

func accountAttributesInvalid(t *testing.T) *form3.AccountAttributes {
	t.Helper()
	attributes := accountAttributesRequired(t)
	// Invalid account number
	attributes.AccountNumber = "%$#@!123654"

	return attributes
}
