package imaging

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	clinicaldomain "github.com/Acauhi99/med-vault/internal/clinical/domain"
	"github.com/Acauhi99/med-vault/internal/generated"
	imagingdomain "github.com/Acauhi99/med-vault/internal/imaging/domain"
	sharedauth "github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/google/uuid"
)

type imagingAPIResponse struct {
	Data struct {
		UploadURL string `json:"upload_url"`
		S3Key     string `json:"s3_key"`
		ExpiresIn int    `json:"expires_in"`
	} `json:"data"`
	Error *struct {
		Code string `json:"code"`
	} `json:"error"`
}

type mockImageRepo struct {
	images map[uuid.UUID]*imagingdomain.Image
}

func newMockImageRepo() *mockImageRepo {
	return &mockImageRepo{images: make(map[uuid.UUID]*imagingdomain.Image)}
}

func (m *mockImageRepo) Create(_ context.Context, img *imagingdomain.Image) error {
	m.images[img.ID] = img
	return nil
}

func (m *mockImageRepo) ListByCase(_ context.Context, tenantID, caseID uuid.UUID) ([]imagingdomain.Image, error) {
	var out []imagingdomain.Image
	for _, img := range m.images {
		if img.TenantID == tenantID && img.CaseID == caseID {
			out = append(out, *img)
		}
	}
	return out, nil
}

func (m *mockImageRepo) GetByID(_ context.Context, tenantID, imageID uuid.UUID) (*imagingdomain.Image, error) {
	img, ok := m.images[imageID]
	if !ok || img.TenantID != tenantID {
		return nil, errors.New("image not found")
	}
	return img, nil
}

func (m *mockImageRepo) FindByS3Key(_ context.Context, tenantID, caseID uuid.UUID, s3Key string) (*imagingdomain.Image, error) {
	for _, img := range m.images {
		if img.TenantID == tenantID && img.CaseID == caseID && img.S3Key == s3Key {
			return img, nil
		}
	}
	return nil, errors.New("image not found")
}

func (m *mockImageRepo) Delete(_ context.Context, tenantID, imageID uuid.UUID) error {
	img, ok := m.images[imageID]
	if !ok || img.TenantID != tenantID {
		return errors.New("image not found")
	}
	delete(m.images, imageID)
	return nil
}

type mockCaseRepo struct {
	cases map[uuid.UUID]*clinicaldomain.Case
}

func newMockCaseRepo() *mockCaseRepo {
	return &mockCaseRepo{cases: make(map[uuid.UUID]*clinicaldomain.Case)}
}

func (m *mockCaseRepo) Create(context.Context, *clinicaldomain.Case) error { return nil }

func (m *mockCaseRepo) GetByID(_ context.Context, tenantID, caseID uuid.UUID) (*clinicaldomain.Case, error) {
	c, ok := m.cases[caseID]
	if !ok || c.TenantID != tenantID {
		return nil, errors.New("case not found")
	}
	return c, nil
}

func (m *mockCaseRepo) ListByPatient(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ string, _, _ int) ([]clinicaldomain.Case, int, error) {
	return nil, 0, nil
}

func (m *mockCaseRepo) ListByDoctor(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ string, _, _ int) ([]clinicaldomain.Case, int, error) {
	return nil, 0, nil
}

func (m *mockCaseRepo) ListByTenant(_ context.Context, _ uuid.UUID, _ string, _, _ int) ([]clinicaldomain.Case, int, error) {
	return nil, 0, nil
}

func (m *mockCaseRepo) Update(context.Context, *clinicaldomain.Case) error { return nil }
func (m *mockCaseRepo) AddSymptom(context.Context, uuid.UUID, uuid.UUID, *clinicaldomain.Symptom) error {
	return nil
}
func (m *mockCaseRepo) WriteDiagnosis(context.Context, uuid.UUID, uuid.UUID, *clinicaldomain.Diagnosis) error {
	return nil
}

type mockStorage struct{}

func (m *mockStorage) GenerateUploadURL(_ context.Context, key, _ string, _ time.Duration) (string, error) {
	return "https://s3.example.com/upload/" + key, nil
}

func (m *mockStorage) GenerateDownloadURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return "https://s3.example.com/download/" + key, nil
}

func TestAPI_RequestUploadURL(t *testing.T) {
	imgRepo := newMockImageRepo()
	caseRepo := newMockCaseRepo()
	caseID := uuid.New()
	tenantID := uuid.New()
	patientID := uuid.New()
	caseRepo.cases[caseID] = &clinicaldomain.Case{ID: caseID, TenantID: tenantID, PatientID: patientID, Status: clinicaldomain.CaseStatusOpen}
	api := NewAPI(imgRepo, caseRepo, &mockStorage{})

	body, err := json.Marshal(generated.UploadURLRequest{FileName: "xray.png", ContentType: generated.UploadURLRequestContentType("image/png"), FileSize: 1024})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	patient := sharedauth.Principal{UserID: patientID, TenantID: tenantID, Role: sharedauth.RolePatient}
	doctor := sharedauth.Principal{UserID: uuid.New(), TenantID: tenantID, Role: sharedauth.RoleDoctor}

	tests := []struct {
		name       string
		principal  sharedauth.Principal
		wantCode   string
		wantStatus int
	}{
		{name: "doctor forbidden", principal: doctor, wantStatus: http.StatusForbidden, wantCode: "FORBIDDEN"},
		{name: "patient gets url", principal: patient, wantStatus: http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/cases/"+caseID.String()+"/images/upload-url", bytes.NewReader(body))
			req = req.WithContext(sharedauth.ContextWithPrincipal(req.Context(), tc.principal))

			api.RequestUploadURL(rec, req, caseID)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
			}

			var resp imagingAPIResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			if tc.wantCode != "" {
				if resp.Error == nil || resp.Error.Code != tc.wantCode {
					t.Fatalf("error code = %#v, want %q", resp.Error, tc.wantCode)
				}
				return
			}

			if resp.Error != nil {
				t.Fatalf("unexpected error: %#v", resp.Error)
			}
			if resp.Data.UploadURL == "" || resp.Data.S3Key == "" {
				t.Fatalf("missing upload payload: %#v", resp.Data)
			}
			if resp.Data.ExpiresIn != 900 {
				t.Fatalf("expires_in = %d, want 900", resp.Data.ExpiresIn)
			}
		})
	}
}

var _ imagingdomain.Repository = (*mockImageRepo)(nil)
var _ clinicaldomain.Repository = (*mockCaseRepo)(nil)
