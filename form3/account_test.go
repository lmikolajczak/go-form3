//go:build integration
// +build integration

package form3_test

import (
	"fmt"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/lmikolajczak/go-form3/form3"
	"testing"
)

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

	attributes.AccountClassification = form3.String("Personal")
	attributes.AccountMatchingOptOut = form3.Bool(false)
	attributes.AccountNumber = "123654"
	attributes.AlternativeNames = []string{"Alternative Names"}
	attributes.BankID = "ABNA"
	attributes.BankIDCode = "ABNANL"
	attributes.BaseCurrency = "EUR"
	attributes.Bic = "ABNANL2A"
	attributes.Iban = "NL91ABNA0417164300"
	attributes.JointAccount = form3.Bool(false)
	attributes.SecondaryIdentification = "Secondary Identification"
	attributes.Status = form3.String("pending")
	attributes.Switched = form3.Bool(false)

	return attributes
}

func accountAttributesInvalid(t *testing.T) *form3.AccountAttributes {
	t.Helper()
	attributes := accountAttributesRequired(t)
	// Invalid account number
	attributes.AccountNumber = "%$#@!123654"

	return attributes
}

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
			wantAttrs:  accountAttributesOptional(t),
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
			organisationID := uuid.NewString()
			account, err := f3.CreateAccount(organisationID, tc.attributes)
			if err != nil {
				testErrorMessage(t, err, tc.wantErr)
			} else {
				testUUID(t, account.OrganisationID, organisationID)
				testAccountAttrs(t, account.Attributes, tc.wantAttrs)
			}
		})
	}
}

func TestClient_DeleteAccount(t *testing.T) {
	f3, teardown := form3.TestClient(t)
	defer teardown()

	account, _ := f3.CreateAccount(uuid.NewString(), accountAttributesRequired(t))

	testcases := []struct {
		name           string
		accountID      string
		accountVersion int64
		wantErr        error
	}{
		{
			name:           "account does not exist",
			accountID:      uuid.NewString(),
			accountVersion: 0,
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
			if err != nil {
				testErrorMessage(t, err, tc.wantErr)
			}
		})
	}
}

func TestClient_FetchAccount(t *testing.T) {
	f3, teardown := form3.TestClient(t)
	defer teardown()

	acc, _ := f3.CreateAccount(uuid.NewString(), accountAttributesRequired(t))
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
			accountID: acc.ID,
			wantAttrs: &form3.AccountAttributes{
				Country: form3.String("NL"),
				Name:    []string{"L. Mikolajczak"},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			account, err := f3.FetchAccount(tc.accountID)
			if err != nil {
				testErrorMessage(t, err, tc.wantErr)
			} else {
				testUUID(t, account.ID, tc.accountID)
				testAccountAttrs(t, account.Attributes, tc.wantAttrs)
			}
		})
	}
}

func testUUID(t *testing.T, uuid string, want string) {
	t.Helper()
	if got := uuid; got != want {
		t.Errorf("uuid = %s; want: %s", got, want)
	}
}

func testAccountAttrs(t *testing.T, attrs *form3.AccountAttributes, want *form3.AccountAttributes) {
	t.Helper()
	if diff := deep.Equal(attrs, want); diff != nil {
		t.Error(diff)
	}
}

func testErrorMessage(t *testing.T, err error, want error) {
	t.Helper()
	if err != nil && want == nil {
		t.Errorf("error message: %s; want: nil", err.Error())
	}
	if err == nil && want != nil {
		t.Errorf("error message: nil; want: %s", want.Error())
	}
	if err != nil && want != nil {
		if got := err.Error(); got != want.Error() {
			t.Errorf("error message: %s; want: %s", got, want)
		}
	}
}
