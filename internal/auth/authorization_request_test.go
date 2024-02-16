package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthorizationRequestHandler_checkResponseType(t *testing.T) {
	handler := authorizationRequestHandler{}
	for _, tc := range []struct {
		input    string
		expected func(error) bool
	}{
		{"foobar", func(err error) bool { return err == nil }},
		{"code", func(err error) bool { return err != nil }},
	} {
		if err := handler.checkResponseType(tc.input); tc.expected(err) {
			t.Errorf("test case failed: %v", err)
		}
	}
}

func TestAuthorizationRequestHandler_checkMethod(t *testing.T) {
	handler := authorizationRequestHandler{}
	for _, tc := range []struct {
		subjectTmpl string
		input       *http.Request
		expected    func(error) bool
	}{
		{"expected GET request: %v", httptest.NewRequest("GET", "/foo", nil), func(err error) bool { return err != nil }},
		{"expected POST request: %v", httptest.NewRequest("POST", "/foo", nil), func(err error) bool { return err != nil }},
		{"expected error instead of PUT request: %v", httptest.NewRequest("PUT", "/foo", nil), func(err error) bool { return err == nil }},
	} {
		if err := handler.checkMethod(tc.input); tc.expected(err) {
			t.Errorf(tc.subjectTmpl, err)
		}
	}
}

func TestAuthorizationRequestHandler_checkCodeChallengeMethod(t *testing.T) {
	handler := authorizationRequestHandler{}
	for _, tc := range []struct {
		subjectTmpl string
		input       *http.Request
		expected    func(error) bool
	}{
		{"expected error: %v", httptest.NewRequest("GET", "/foo?code_challenge_method=abc", nil), func(err error) bool { return err == nil }},
		{"expected error: %v", httptest.NewRequest("GET", "/foo", nil), func(err error) bool { return err == nil }},
		{"expected no error: %v", httptest.NewRequest("GET", "/foo?code_challenge_method=S256", nil), func(err error) bool { return err != nil }},
	} {
		if err := handler.checkCodeChallengeMethod(tc.input); tc.expected(err) {
			t.Errorf(tc.subjectTmpl, err)
		}
	}
}
