package room

import "testing"

func TestHTTPMethod_String(t *testing.T) {
	tests := []struct {
		method   HTTPMethod
		expected string
	}{
		{GET, "GET"},
		{POST, "POST"},
		{PUT, "PUT"},
		{PATCH, "PATCH"},
		{DELETE, "DELETE"},
		{HEAD, "HEAD"},
		{"INVALID", "GET"}, // Test case for invalid method, should default to GET
	}

	for _, test := range tests {
		result := test.method.String()
		if result != test.expected {
			t.Errorf("HTTPMethod.String() returned %s, expected %s", result, test.expected)
		}
	}
}

func TestHTTPProtocol_String(t *testing.T) {
	tests := []struct {
		protocol HTTPProtocol
		expected string
	}{
		{Http, "http://"},
		{Https, "https://"},
		{HTTPProtocol(999), "https://"}, // Test case for invalid protocol, should default to HTTPS
	}

	for _, test := range tests {
		result := test.protocol.String()
		if result != test.expected {
			t.Errorf("HTTPProtocol.String() returned %s, expected %s", result, test.expected)
		}
	}
}
