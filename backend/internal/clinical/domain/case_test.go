package domain

import "testing"

func TestValidateTransition(t *testing.T) {
	tests := []struct {
		name    string
		from    CaseStatus
		to      CaseStatus
		wantErr bool
	}{
		{"open→assigned", CaseStatusOpen, CaseStatusAssigned, false},
		{"open→diagnosed", CaseStatusOpen, CaseStatusDiagnosed, true},
		{"open→closed", CaseStatusOpen, CaseStatusClosed, true},
		{"assigned→diagnosed", CaseStatusAssigned, CaseStatusDiagnosed, false},
		{"assigned→open", CaseStatusAssigned, CaseStatusOpen, true},
		{"assigned→closed", CaseStatusAssigned, CaseStatusClosed, true},
		{"diagnosed→closed", CaseStatusDiagnosed, CaseStatusClosed, false},
		{"diagnosed→open", CaseStatusDiagnosed, CaseStatusOpen, true},
		{"closed→open", CaseStatusClosed, CaseStatusOpen, true},
		{"closed→assigned", CaseStatusClosed, CaseStatusAssigned, true},
		{"closed→diagnosed", CaseStatusClosed, CaseStatusDiagnosed, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTransition(tt.from, tt.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTransition(%s, %s) error = %v, wantErr %v", tt.from, tt.to, err, tt.wantErr)
			}
		})
	}
}
