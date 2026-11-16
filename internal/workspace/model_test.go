package workspace

import (
	"strings"
	"testing"

	"github.com/thakur-vijay/adorbit-go/validation"
)

func TestCreateWorkspaceDTO_Validation(t *testing.T) {
	tests := []struct {
		name  string
		dto   CreateWorkspaceDTO
		valid bool
	}{
		{name: "all valid fields", dto: testWorkspaceDTO(), valid: true},
		{name: "advertiser type", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Type = Advertiser }), valid: true},
		{name: "publisher type", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Type = Publisher }), valid: true},
		{name: "missing name", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Name = "" })},
		{name: "name at maximum length", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Name = strings.Repeat("n", 150) }), valid: true},
		{name: "name exceeds maximum length", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Name = strings.Repeat("n", 151) })},
		{name: "name with leading and trailing whitespace", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Name = " Workspace " }), valid: true},
		{name: "name with Unicode characters", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Name = "日本語" }), valid: true},
		{name: "name with special characters", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Name = "Workspace @ # & Co.!" }), valid: true},
		{name: "unsupported workspace type", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Type = "agency" })},
		{name: "uppercase country code", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Country = "CA" }), valid: true},
		{name: "country with incorrect length", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Country = "USA" })},
		{name: "country with invalid characters passes length-only rule", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Country = "1@" }), valid: true},
		{name: "country with punctuation passes length-only rule", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Country = "I!" }), valid: true},
		{name: "lowercase country passes length-only rule", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Country = "us" }), valid: true},
		{name: "empty required country", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Country = "" })},
		{name: "uppercase currency code", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Currency = "EUR" }), valid: true},
		{name: "currency with incorrect length", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Currency = "EURO" })},
		{name: "currency with invalid characters passes length-only rule", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Currency = "1@!" }), valid: true},
		{name: "unknown currency code passes length-only rule", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Currency = "ZZZ" }), valid: true},
		{name: "lowercase currency passes length-only rule", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Currency = "usd" }), valid: true},
		{name: "empty required currency", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Currency = "" })},
		{name: "timezone exceeds maximum length", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Timezone = strings.Repeat("t", 101) })},
		{name: "timezone at maximum length", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Timezone = strings.Repeat("t", 100) }), valid: true},
		{name: "valid IANA timezone", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Timezone = "Asia/Kolkata" }), valid: true},
		{name: "invalid timezone passes length-only rule", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Timezone = "Invalid/Zone" }), valid: true},
		{name: "empty required timezone", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Timezone = "" })},
		{name: "industry exceeds maximum length", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Industry = stringPointer(strings.Repeat("i", 101)) })},
		{name: "industry at maximum length", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Industry = stringPointer(strings.Repeat("i", 100)) }), valid: true},
		{name: "valid HTTP website", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = stringPointer("http://example.test") }), valid: true},
		{name: "valid HTTPS website", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = stringPointer("https://example.test/path") }), valid: true},
		{name: "FTP website URL", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = stringPointer("ftp://files.example.test/archive.zip") }), valid: true},
		{name: "website URL with query parameters", dto: dtoWith(func(dto *CreateWorkspaceDTO) {
			dto.Website = stringPointer("https://example.test/path?campaign=spring&source=search")
		}), valid: true},
		{name: "website URL with explicit port", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = stringPointer("https://example.test:8443/path") }), valid: true},
		{name: "website without scheme", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = stringPointer("example.test") })},
		{name: "malformed website URL", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = stringPointer("http://[::1") })},
		{name: "description exceeding unrelated field limits", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Description = stringPointer(strings.Repeat("d", 500)) }), valid: true},
		{name: "nil optional pointers", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Industry = nil; dto.Website = nil; dto.Description = nil }), valid: true},
		{name: "empty optional strings", dto: dtoWith(func(dto *CreateWorkspaceDTO) {
			dto.Industry = stringPointer("")
			dto.Website = stringPointer("")
			dto.Description = stringPointer("")
		})},
		{name: "populated optional strings", dto: dtoWith(func(dto *CreateWorkspaceDTO) {
			dto.Industry = stringPointer("Technology")
			dto.Website = stringPointer("https://example.test")
			dto.Description = stringPointer("A workspace")
		}), valid: true},
		{name: "name containing only whitespace passes required tag", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Name = "   " }), valid: true},
		{name: "nil industry pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Industry = nil }), valid: true},
		{name: "empty industry pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Industry = stringPointer("") }), valid: true},
		{name: "whitespace-only industry pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Industry = stringPointer("   ") }), valid: true},
		{name: "nil description pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Description = nil }), valid: true},
		{name: "empty description pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Description = stringPointer("") }), valid: true},
		{name: "whitespace-only description pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Description = stringPointer("   ") }), valid: true},
		{name: "invalid industry pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Industry = stringPointer(strings.Repeat("i", 101)) })},
		{name: "nil website pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = nil }), valid: true},
		{name: "empty website pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = stringPointer("") })},
		{name: "whitespace-only website pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = stringPointer("   ") })},
		{name: "invalid website pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = stringPointer("example.test") })},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validation.Validator(test.dto)
			if test.valid && err != nil {
				t.Errorf("Validator() error = %v, want nil", err)
			}
			if !test.valid && err == nil {
				t.Error("Validator() error = nil, want validation error")
			}
		})
	}
}

func dtoWith(update func(*CreateWorkspaceDTO)) CreateWorkspaceDTO {
	dto := testWorkspaceDTO()
	update(&dto)
	return dto
}
