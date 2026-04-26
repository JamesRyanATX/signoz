package authtypes

import (
	"encoding/json"

	"github.com/SigNoz/signoz/pkg/errors"
	"github.com/SigNoz/signoz/pkg/valuer"
)

// PostableLDAPSession is the request body for LDAP login. The identifier may
// be either an email address or a plain username (for servers that search by
// sAMAccountName or uid).
type PostableLDAPSession struct {
	Identifier string      `json:"identifier"`
	Password   string      `json:"password"`
	OrgID      valuer.UUID `json:"orgId"`
}

func (typ *PostableLDAPSession) UnmarshalJSON(data []byte) error {
	type Alias PostableLDAPSession
	var temp Alias

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.Identifier == "" {
		return errors.New(errors.TypeInvalidInput, errors.CodeInvalidInput, "identifier is required")
	}

	if temp.Password == "" {
		return errors.New(errors.TypeInvalidInput, errors.CodeInvalidInput, "password is required")
	}

	if temp.OrgID.IsZero() {
		return errors.New(errors.TypeInvalidInput, errors.CodeInvalidInput, "orgID is required")
	}

	*typ = PostableLDAPSession(temp)
	return nil
}
