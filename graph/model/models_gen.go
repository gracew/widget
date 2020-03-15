// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type APIDefinition struct {
	Name       string                 `json:"name"`
	Fields     []*FieldDefinition     `json:"fields"`
	Operations []*OperationDefinition `json:"operations"`
}

type Auth struct {
	ID                 string             `json:"id"`
	APIID              string             `json:"apiID"`
	AuthenticationType AuthenticationType `json:"authenticationType"`
	ReadPolicy         *AuthPolicy        `json:"readPolicy"`
	WritePolicy        *AuthPolicy        `json:"writePolicy"`
}

type AuthAPIInput struct {
	APIID              string             `json:"apiID"`
	AuthenticationType AuthenticationType `json:"authenticationType"`
	ReadPolicy         *AuthPolicyInput   `json:"readPolicy"`
	WritePolicy        *AuthPolicyInput   `json:"writePolicy"`
}

type AuthPolicy struct {
	Type            AuthPolicyType `json:"type"`
	UserAttribute   *string        `json:"userAttribute"`
	ObjectAttribute *string        `json:"objectAttribute"`
}

type AuthPolicyInput struct {
	Type            AuthPolicyType `json:"type"`
	UserAttribute   *string        `json:"userAttribute"`
	ObjectAttribute *string        `json:"objectAttribute"`
}

type Constraint struct {
	MinInt    *int     `json:"minInt"`
	MaxInt    *int     `json:"maxInt"`
	MinFloat  *float64 `json:"minFloat"`
	MaxFloat  *float64 `json:"maxFloat"`
	Regex     *string  `json:"regex"`
	MinLength *int     `json:"minLength"`
	MaxLength *int     `json:"maxLength"`
}

type DefineAPIInput struct {
	RawDefinition string `json:"rawDefinition"`
}

type Deploy struct {
	ID    string      `json:"id"`
	APIID string      `json:"apiID"`
	Env   Environment `json:"env"`
}

type DeployAPIInput struct {
	APIID string      `json:"apiID"`
	Env   Environment `json:"env"`
}

type FieldDefinition struct {
	Name        string      `json:"name"`
	Type        Type        `json:"type"`
	CustomType  *string     `json:"customType"`
	Optional    *bool       `json:"optional"`
	List        *bool       `json:"list"`
	Constraints *Constraint `json:"constraints"`
}

type OperationDefinition struct {
	Type   OperationType     `json:"type"`
	Sort   []*SortDefinition `json:"sort"`
	Filter []string          `json:"filter"`
}

type SaveCustomLogicInput struct {
	APIID         string        `json:"apiID"`
	OperationType OperationType `json:"operationType"`
	Language      Language      `json:"language"`
	Before        *string       `json:"before"`
	After         *string       `json:"after"`
}

type SortDefinition struct {
	Field string    `json:"field"`
	Order SortOrder `json:"order"`
}

type TestToken struct {
	Label string `json:"label"`
	Token string `json:"token"`
}

type TestTokenInput struct {
	Label string `json:"label"`
	Token string `json:"token"`
}

type TestTokenResponse struct {
	TestTokens []*TestToken `json:"testTokens"`
}

type UpdateAPIInput struct {
	APIID         string `json:"apiID"`
	RawDefinition string `json:"rawDefinition"`
}

type AuthPolicyType string

const (
	AuthPolicyTypeCreatedBy      AuthPolicyType = "CREATED_BY"
	AuthPolicyTypeAttributeMatch AuthPolicyType = "ATTRIBUTE_MATCH"
	AuthPolicyTypeCustom         AuthPolicyType = "CUSTOM"
)

var AllAuthPolicyType = []AuthPolicyType{
	AuthPolicyTypeCreatedBy,
	AuthPolicyTypeAttributeMatch,
	AuthPolicyTypeCustom,
}

func (e AuthPolicyType) IsValid() bool {
	switch e {
	case AuthPolicyTypeCreatedBy, AuthPolicyTypeAttributeMatch, AuthPolicyTypeCustom:
		return true
	}
	return false
}

func (e AuthPolicyType) String() string {
	return string(e)
}

func (e *AuthPolicyType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AuthPolicyType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AuthPolicyType", str)
	}
	return nil
}

func (e AuthPolicyType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type AuthenticationType string

const (
	AuthenticationTypeBuiltIn AuthenticationType = "BUILT_IN"
)

var AllAuthenticationType = []AuthenticationType{
	AuthenticationTypeBuiltIn,
}

func (e AuthenticationType) IsValid() bool {
	switch e {
	case AuthenticationTypeBuiltIn:
		return true
	}
	return false
}

func (e AuthenticationType) String() string {
	return string(e)
}

func (e *AuthenticationType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AuthenticationType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AuthenticationType", str)
	}
	return nil
}

func (e AuthenticationType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Environment string

const (
	EnvironmentSandbox    Environment = "SANDBOX"
	EnvironmentStaging    Environment = "STAGING"
	EnvironmentProduction Environment = "PRODUCTION"
)

var AllEnvironment = []Environment{
	EnvironmentSandbox,
	EnvironmentStaging,
	EnvironmentProduction,
}

func (e Environment) IsValid() bool {
	switch e {
	case EnvironmentSandbox, EnvironmentStaging, EnvironmentProduction:
		return true
	}
	return false
}

func (e Environment) String() string {
	return string(e)
}

func (e *Environment) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Environment(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Environment", str)
	}
	return nil
}

func (e Environment) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Language string

const (
	LanguageJavascript Language = "JAVASCRIPT"
	LanguagePython     Language = "PYTHON"
)

var AllLanguage = []Language{
	LanguageJavascript,
	LanguagePython,
}

func (e Language) IsValid() bool {
	switch e {
	case LanguageJavascript, LanguagePython:
		return true
	}
	return false
}

func (e Language) String() string {
	return string(e)
}

func (e *Language) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Language(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Language", str)
	}
	return nil
}

func (e Language) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type OperationType string

const (
	OperationTypeCreate OperationType = "CREATE"
	OperationTypeUpdate OperationType = "UPDATE"
	OperationTypeRead   OperationType = "READ"
	OperationTypeList   OperationType = "LIST"
)

var AllOperationType = []OperationType{
	OperationTypeCreate,
	OperationTypeUpdate,
	OperationTypeRead,
	OperationTypeList,
}

func (e OperationType) IsValid() bool {
	switch e {
	case OperationTypeCreate, OperationTypeUpdate, OperationTypeRead, OperationTypeList:
		return true
	}
	return false
}

func (e OperationType) String() string {
	return string(e)
}

func (e *OperationType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OperationType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OperationType", str)
	}
	return nil
}

func (e OperationType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

var AllSortOrder = []SortOrder{
	SortOrderAsc,
	SortOrderDesc,
}

func (e SortOrder) IsValid() bool {
	switch e {
	case SortOrderAsc, SortOrderDesc:
		return true
	}
	return false
}

func (e SortOrder) String() string {
	return string(e)
}

func (e *SortOrder) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortOrder(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortOrder", str)
	}
	return nil
}

func (e SortOrder) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Type string

const (
	TypeFloat   Type = "FLOAT"
	TypeInt     Type = "INT"
	TypeBoolean Type = "BOOLEAN"
	TypeString  Type = "STRING"
)

var AllType = []Type{
	TypeFloat,
	TypeInt,
	TypeBoolean,
	TypeString,
}

func (e Type) IsValid() bool {
	switch e {
	case TypeFloat, TypeInt, TypeBoolean, TypeString:
		return true
	}
	return false
}

func (e Type) String() string {
	return string(e)
}

func (e *Type) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Type(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Type", str)
	}
	return nil
}

func (e Type) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
