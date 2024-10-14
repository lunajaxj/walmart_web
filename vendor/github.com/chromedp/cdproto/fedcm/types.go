package fedcm

// Code generated by cdproto-gen. DO NOT EDIT.

import (
	"fmt"

	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

// LoginState whether this is a sign-up or sign-in action for this account,
// i.e. whether this account has ever been used to sign in to this RP before.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/FedCm#type-LoginState
type LoginState string

// String returns the LoginState as string value.
func (t LoginState) String() string {
	return string(t)
}

// LoginState values.
const (
	LoginStateSignIn LoginState = "SignIn"
	LoginStateSignUp LoginState = "SignUp"
)

// MarshalEasyJSON satisfies easyjson.Marshaler.
func (t LoginState) MarshalEasyJSON(out *jwriter.Writer) {
	out.String(string(t))
}

// MarshalJSON satisfies json.Marshaler.
func (t LoginState) MarshalJSON() ([]byte, error) {
	return easyjson.Marshal(t)
}

// UnmarshalEasyJSON satisfies easyjson.Unmarshaler.
func (t *LoginState) UnmarshalEasyJSON(in *jlexer.Lexer) {
	v := in.String()
	switch LoginState(v) {
	case LoginStateSignIn:
		*t = LoginStateSignIn
	case LoginStateSignUp:
		*t = LoginStateSignUp

	default:
		in.AddError(fmt.Errorf("unknown LoginState value: %v", v))
	}
}

// UnmarshalJSON satisfies json.Unmarshaler.
func (t *LoginState) UnmarshalJSON(buf []byte) error {
	return easyjson.Unmarshal(buf, t)
}

// DialogType the types of FedCM dialogs.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/FedCm#type-DialogType
type DialogType string

// String returns the DialogType as string value.
func (t DialogType) String() string {
	return string(t)
}

// DialogType values.
const (
	DialogTypeAccountChooser  DialogType = "AccountChooser"
	DialogTypeAutoReauthn     DialogType = "AutoReauthn"
	DialogTypeConfirmIdpLogin DialogType = "ConfirmIdpLogin"
	DialogTypeError           DialogType = "Error"
)

// MarshalEasyJSON satisfies easyjson.Marshaler.
func (t DialogType) MarshalEasyJSON(out *jwriter.Writer) {
	out.String(string(t))
}

// MarshalJSON satisfies json.Marshaler.
func (t DialogType) MarshalJSON() ([]byte, error) {
	return easyjson.Marshal(t)
}

// UnmarshalEasyJSON satisfies easyjson.Unmarshaler.
func (t *DialogType) UnmarshalEasyJSON(in *jlexer.Lexer) {
	v := in.String()
	switch DialogType(v) {
	case DialogTypeAccountChooser:
		*t = DialogTypeAccountChooser
	case DialogTypeAutoReauthn:
		*t = DialogTypeAutoReauthn
	case DialogTypeConfirmIdpLogin:
		*t = DialogTypeConfirmIdpLogin
	case DialogTypeError:
		*t = DialogTypeError

	default:
		in.AddError(fmt.Errorf("unknown DialogType value: %v", v))
	}
}

// UnmarshalJSON satisfies json.Unmarshaler.
func (t *DialogType) UnmarshalJSON(buf []byte) error {
	return easyjson.Unmarshal(buf, t)
}

// DialogButton the buttons on the FedCM dialog.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/FedCm#type-DialogButton
type DialogButton string

// String returns the DialogButton as string value.
func (t DialogButton) String() string {
	return string(t)
}

// DialogButton values.
const (
	DialogButtonConfirmIdpLoginContinue DialogButton = "ConfirmIdpLoginContinue"
	DialogButtonErrorGotIt              DialogButton = "ErrorGotIt"
	DialogButtonErrorMoreDetails        DialogButton = "ErrorMoreDetails"
)

// MarshalEasyJSON satisfies easyjson.Marshaler.
func (t DialogButton) MarshalEasyJSON(out *jwriter.Writer) {
	out.String(string(t))
}

// MarshalJSON satisfies json.Marshaler.
func (t DialogButton) MarshalJSON() ([]byte, error) {
	return easyjson.Marshal(t)
}

// UnmarshalEasyJSON satisfies easyjson.Unmarshaler.
func (t *DialogButton) UnmarshalEasyJSON(in *jlexer.Lexer) {
	v := in.String()
	switch DialogButton(v) {
	case DialogButtonConfirmIdpLoginContinue:
		*t = DialogButtonConfirmIdpLoginContinue
	case DialogButtonErrorGotIt:
		*t = DialogButtonErrorGotIt
	case DialogButtonErrorMoreDetails:
		*t = DialogButtonErrorMoreDetails

	default:
		in.AddError(fmt.Errorf("unknown DialogButton value: %v", v))
	}
}

// UnmarshalJSON satisfies json.Unmarshaler.
func (t *DialogButton) UnmarshalJSON(buf []byte) error {
	return easyjson.Unmarshal(buf, t)
}

// AccountURLType the URLs that each account has.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/FedCm#type-AccountUrlType
type AccountURLType string

// String returns the AccountURLType as string value.
func (t AccountURLType) String() string {
	return string(t)
}

// AccountURLType values.
const (
	AccountURLTypeTermsOfService AccountURLType = "TermsOfService"
	AccountURLTypePrivacyPolicy  AccountURLType = "PrivacyPolicy"
)

// MarshalEasyJSON satisfies easyjson.Marshaler.
func (t AccountURLType) MarshalEasyJSON(out *jwriter.Writer) {
	out.String(string(t))
}

// MarshalJSON satisfies json.Marshaler.
func (t AccountURLType) MarshalJSON() ([]byte, error) {
	return easyjson.Marshal(t)
}

// UnmarshalEasyJSON satisfies easyjson.Unmarshaler.
func (t *AccountURLType) UnmarshalEasyJSON(in *jlexer.Lexer) {
	v := in.String()
	switch AccountURLType(v) {
	case AccountURLTypeTermsOfService:
		*t = AccountURLTypeTermsOfService
	case AccountURLTypePrivacyPolicy:
		*t = AccountURLTypePrivacyPolicy

	default:
		in.AddError(fmt.Errorf("unknown AccountURLType value: %v", v))
	}
}

// UnmarshalJSON satisfies json.Unmarshaler.
func (t *AccountURLType) UnmarshalJSON(buf []byte) error {
	return easyjson.Unmarshal(buf, t)
}

// Account corresponds to IdentityRequestAccount.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/FedCm#type-Account
type Account struct {
	AccountID         string     `json:"accountId"`
	Email             string     `json:"email"`
	Name              string     `json:"name"`
	GivenName         string     `json:"givenName"`
	PictureURL        string     `json:"pictureUrl"`
	IdpConfigURL      string     `json:"idpConfigUrl"`
	IdpLoginURL       string     `json:"idpLoginUrl"`
	LoginState        LoginState `json:"loginState"`
	TermsOfServiceURL string     `json:"termsOfServiceUrl,omitempty"` // These two are only set if the loginState is signUp
	PrivacyPolicyURL  string     `json:"privacyPolicyUrl,omitempty"`
}
