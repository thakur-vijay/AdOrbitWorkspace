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
		{name: "unsupported workspace type", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Type = "agency" })},
		{name: "country with two characters", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Country = "CA" }), valid: true},
		{name: "country with incorrect length", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Country = "USA" })},
		{name: "lowercase country is accepted by length rule", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Country = "us" }), valid: true},
		{name: "empty required country", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Country = "" })},
		{name: "currency with three characters", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Currency = "EUR" }), valid: true},
		{name: "currency with incorrect length", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Currency = "EURO" })},
		{name: "lowercase currency is accepted by length rule", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Currency = "usd" }), valid: true},
		{name: "empty required currency", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Currency = "" })},
		{name: "timezone exceeds maximum length", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Timezone = strings.Repeat("t", 101) })},
		{name: "timezone at maximum length", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Timezone = strings.Repeat("t", 100) }), valid: true},
		{name: "empty required timezone", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Timezone = "" })},
		{name: "industry exceeds maximum length", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Industry = stringPointer(strings.Repeat("i", 101)) })},
		{name: "industry at maximum length", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Industry = stringPointer(strings.Repeat("i", 100)) }), valid: true},
		{name: "valid HTTP website", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = stringPointer("http://example.test") }), valid: true},
		{name: "valid HTTPS website", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = stringPointer("https://example.test/path") }), valid: true},
		{name: "website without scheme", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = stringPointer("example.test") })},
		{name: "malformed website URL", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = stringPointer("http://[::1") })},
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
		{name: "name with leading and trailing spaces", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Name = " Workspace " }), valid: true},
		{name: "nil industry pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Industry = nil }), valid: true},
		{name: "empty industry pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Industry = stringPointer("") }), valid: true},
		{name: "empty description pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Description = stringPointer("") }), valid: true},
		{name: "invalid industry pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Industry = stringPointer(strings.Repeat("i", 101)) })},
		{name: "nil website pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = nil }), valid: true},
		{name: "empty website pointer", dto: dtoWith(func(dto *CreateWorkspaceDTO) { dto.Website = stringPointer("") })},
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
