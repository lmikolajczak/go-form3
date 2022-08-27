package form3

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

// Account represents an account in the form3 org section.
type Account struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
}

// String returns a simple string representation of the Account.
func (a Account) String() string {
	return fmt.Sprintf("Account(id=%s, version=%d)", a.ID, *a.Version)
}

// AccountAttributes represents attributes of a single account.
type AccountAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
}

// AccountJSON represents request payload to account resource.
type AccountJSON struct {
	Data Account `json:"data,omitempty"`
}

// FetchAccount returns account with the given identifier.
func (c *Client) FetchAccount(id string) (*Account, error) {
	endpoint := fmt.Sprintf("/v1/organisation/accounts/%s", id)
	headers := map[string]string{"Accept": "application/vnd.api+json"}

	request, err := c.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	accountJSON := new(AccountJSON)
	if err = c.Request(accountJSON, request, headers); err != nil {
		return nil, err
	}
	return &accountJSON.Data, nil
}

// CreateAccount creates account with the given attributes.
func (c *Client) CreateAccount(organisationId string, attributes *AccountAttributes) (*Account, error) {
	endpoint := "/v1/organisation/accounts"
	headers := map[string]string{
		"Accept":       "application/vnd.api+json",
		"Content-Type": "application/vnd.api+json",
	}

	payload := AccountJSON{
		Data: Account{
			Attributes:     attributes,
			ID:             uuid.NewString(),
			OrganisationID: organisationId,
			Type:           "accounts",
		},
	}
	request, err := c.NewRequest(http.MethodPost, endpoint, payload)
	if err != nil {
		return nil, err
	}

	accountJSON := new(AccountJSON)
	if err = c.Request(accountJSON, request, headers); err != nil {
		return nil, err
	}
	return &accountJSON.Data, nil
}

// DeleteAccount deletes the account with the given identifier.
func (c *Client) DeleteAccount(id string, version int64) error {
	endpoint := fmt.Sprintf("/v1/organisation/accounts/%s?version=%d", id, version)
	headers := map[string]string{"Accept": "application/vnd.api+json"}

	request, err := c.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	return c.Request(nil, request, headers)
}
