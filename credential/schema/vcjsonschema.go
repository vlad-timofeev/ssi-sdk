package schema

import (
	"context"
	"fmt"

	"github.com/TBD54566975/ssi-sdk/credential"
	"github.com/TBD54566975/ssi-sdk/schema"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
)

// ValidateCredentialAgainstSchema validates a credential against a schema, returning an error if it is not valid
// The schema is retrieved from the given VCJSONSchemaAccess using the credential's credential schema ID
func ValidateCredentialAgainstSchema(access VCJSONSchemaAccess, cred credential.VerifiableCredential) error {
	jsonSchema, _, err := GetCredentialSchemaFromCredential(access, cred)
	if err != nil {
		return errors.Wrap(err, "getting schema from credential")
	}
	if err = IsCredentialValidForJSONSchema(cred, jsonSchema); err != nil {
		return errors.Wrap(err, "credential not valid for schema")
	}
	return nil
}

// IsCredentialValidForJSONSchema validates a credential against a schema, returning an error if it is not valid
func IsCredentialValidForJSONSchema(cred credential.VerifiableCredential, s JSONSchema) error {
	if !IsSupportedVCJSONSchemaType(cred.CredentialSchema.Type) {
		return fmt.Errorf("credential schema type<%s> is not supported", cred.CredentialSchema.Type)
	}
	if !IsSupportedJSONSchemaVersion(s.Schema()) {
		return fmt.Errorf("schema version<%s> is not supported", s.Schema())
	}
	schemaBytes, err := json.Marshal(s)
	if err != nil {
		return errors.Wrap(err, "marshalling schema")
	}
	credBytes, err := json.Marshal(cred)
	if err != nil {
		return errors.Wrap(err, "marshalling credential")
	}
	if err = schema.IsValidAgainstJSONSchema(string(credBytes), string(schemaBytes)); err != nil {
		return errors.Wrap(err, "credential not valid for schema")
	}
	return nil
}

// GetCredentialSchemaFromCredential returns the credential schema and type for a given credential given
// a credential schema access, which is used to retrieve the schema
func GetCredentialSchemaFromCredential(access VCJSONSchemaAccess, cred credential.VerifiableCredential) (JSONSchema, VCJSONSchemaType, error) {
	if cred.CredentialSchema == nil {
		return nil, "", errors.New("credential does not contain a credential schema")
	}

	t := cred.CredentialSchema.Type
	if !IsSupportedVCJSONSchemaType(t) {
		return nil, "", fmt.Errorf("credential schema type<%s> is not supported", t)
	}

	jsonSchema, err := access.GetVCJSONSchema(context.Background(), VCJSONSchemaType(t), cred.CredentialSchema.ID)
	if err != nil {
		return nil, "", errors.Wrap(err, "getting schema")
	}
	return jsonSchema, VCJSONSchemaType(t), nil
}
