package application

import (
	"context"
	"testing"
	"time"

	clinicaldomain "github.com/Acauhi99/med-vault/internal/clinical/domain"
	imagingdomain "github.com/Acauhi99/med-vault/internal/imaging/domain"
	sharedauth "github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/google/uuid"
)

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
		if img.CaseID == caseID {
			out = append(out, *img)
		}
	}
	return out, nil
}

func (m *mockImageRepo) GetByID(_ context.Context, tenantID, imageID uuid.UUID) (*imagingdomain.Image, error) {
	img, ok := m.images[imageID]
	if !ok || img.TenantID != tenantID {
		return nil, ErrImageNotFound
	}
	return img, nil
}

func (m *mockImageRepo) FindByS3Key(_ context.Context, tenantID, caseID uuid.UUID, s3Key string) (*imagingdomain.Image, error) {
	for _, img := range m.images {
		if img.TenantID == tenantID && img.CaseID == caseID && img.S3Key == s3Key {
			return img, nil
		}
	}
	return nil, ErrImageNotFound
}

func (m *mockImageRepo) Delete(_ context.Context, tenantID, imageID uuid.UUID) error {
	img, ok := m.images[imageID]
	if !ok || img.TenantID != tenantID {
		return ErrImageNotFound
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

func (m *mockCaseRepo) Create(_ context.Context, _ *clinicaldomain.Case) error {
	return nil
}

func (m *mockCaseRepo) GetByID(_ context.Context, tenantID, caseID uuid.UUID) (*clinicaldomain.Case, error) {
	c, ok := m.cases[caseID]
	if !ok || c.TenantID != tenantID {
		return nil, ErrCaseNotFound
	}
	return c, nil
}

func (m *mockCaseRepo) ListByPatient(_ context.Context, _, _ uuid.UUID, _ string, _, _ int) ([]clinicaldomain.Case, int, error) {
	return nil, 0, nil
}

func (m *mockCaseRepo) ListByDoctor(_ context.Context, _, _ uuid.UUID, _ string, _, _ int) ([]clinicaldomain.Case, int, error) {
	return nil, 0, nil
}

func (m *mockCaseRepo) ListByTenant(_ context.Context, _ uuid.UUID, _ string, _, _ int) ([]clinicaldomain.Case, int, error) {
	return nil, 0, nil
}

func (m *mockCaseRepo) Update(_ context.Context, _ *clinicaldomain.Case) error {
	return nil
}

func (m *mockCaseRepo) AddSymptom(_ context.Context, _, _ uuid.UUID, _ *clinicaldomain.Symptom) error {
	return nil
}

func (m *mockCaseRepo) WriteDiagnosis(_ context.Context, _, _ uuid.UUID, _ *clinicaldomain.Diagnosis) error {
	return nil
}

type mockStorage struct {
	uploadURLs   map[string]string
	downloadURLs map[string]string
}

func newMockStorage() *mockStorage {
	return &mockStorage{
		uploadURLs:   make(map[string]string),
		downloadURLs: make(map[string]string),
	}
}

func (m *mockStorage) GenerateUploadURL(_ context.Context, key, contentType string, _ time.Duration) (string, error) {
	url := "https://s3.example.com/upload/" + key
	m.uploadURLs[key] = url
	return url, nil
}

func (m *mockStorage) GenerateDownloadURL(_ context.Context, key string, _ time.Duration) (string, error) {
	url := "https://s3.example.com/download/" + key
	m.downloadURLs[key] = url
	return url, nil
}

func newPrincipal(userID, tenantID uuid.UUID, role sharedauth.Role) sharedauth.Principal {
	return sharedauth.Principal{UserID: userID, TenantID: tenantID, Role: role}
}

func setupTestRepos() (*mockImageRepo, *mockCaseRepo, *mockStorage) {
	imgRepo := newMockImageRepo()
	caseRepo := newMockCaseRepo()
	storage := newMockStorage()
	return imgRepo, caseRepo, storage
}

func seedOpenCase(caseRepo *mockCaseRepo, id, tenantID, patientID uuid.UUID) {
	caseRepo.cases[id] = &clinicaldomain.Case{
		ID:        id,
		TenantID:  tenantID,
		PatientID: patientID,
		Status:    clinicaldomain.CaseStatusOpen,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func seedClosedCase(caseRepo *mockCaseRepo, id, tenantID, patientID uuid.UUID) {
	caseRepo.cases[id] = &clinicaldomain.Case{
		ID:        id,
		TenantID:  tenantID,
		PatientID: patientID,
		Status:    clinicaldomain.CaseStatusClosed,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func seedImage(imgRepo *mockImageRepo, id, tenantID, patientID, caseID uuid.UUID) {
	imgRepo.images[id] = &imagingdomain.Image{
		ID:          id,
		TenantID:    tenantID,
		PatientID:   patientID,
		CaseID:      caseID,
		FileName:    "xray.png",
		ContentType: "image/png",
		S3Key:       tenantID.String() + "/" + caseID.String() + "/" + id.String() + "/xray.png",
		UploadedAt:  time.Now().UTC(),
	}
}

func TestRequestUploadURLSuccess(t *testing.T) {
	imgRepo, caseRepo, storage := setupTestRepos()
	patientID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	seedOpenCase(caseRepo, caseID, tenantID, patientID)

	cmd := NewRequestUploadURLCommand(imgRepo, caseRepo, storage)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	url, s3Key, expiresIn, err := cmd.Execute(context.Background(), p, caseID, "xray.png", "image/png", 1024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url == "" {
		t.Error("expected non-empty upload URL")
	}
	if s3Key == "" {
		t.Error("expected non-empty S3 key")
	}
	if expiresIn != 900 {
		t.Errorf("expiresIn = %d, want 900", expiresIn)
	}
}

func TestRequestUploadURLDoctorRejected(t *testing.T) {
	imgRepo, caseRepo, storage := setupTestRepos()
	doctorID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	seedOpenCase(caseRepo, caseID, tenantID, uuid.New())

	cmd := NewRequestUploadURLCommand(imgRepo, caseRepo, storage)
	p := newPrincipal(doctorID, tenantID, sharedauth.RoleDoctor)

	_, _, _, err := cmd.Execute(context.Background(), p, caseID, "xray.png", "image/png", 1024)
	if err != ErrInvalidRole {
		t.Errorf("expected ErrInvalidRole, got %v", err)
	}
}

func TestRequestUploadURLCaseNotFound(t *testing.T) {
	imgRepo, caseRepo, storage := setupTestRepos()
	patientID := uuid.New()
	tenantID := uuid.New()

	cmd := NewRequestUploadURLCommand(imgRepo, caseRepo, storage)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	_, _, _, err := cmd.Execute(context.Background(), p, uuid.New(), "xray.png", "image/png", 1024)
	if err != ErrCaseNotFound {
		t.Errorf("expected ErrCaseNotFound, got %v", err)
	}
}

func TestRequestUploadURLAccessDenied(t *testing.T) {
	imgRepo, caseRepo, storage := setupTestRepos()
	otherPatient := uuid.New()
	patientID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	seedOpenCase(caseRepo, caseID, tenantID, otherPatient)

	cmd := NewRequestUploadURLCommand(imgRepo, caseRepo, storage)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	_, _, _, err := cmd.Execute(context.Background(), p, caseID, "xray.png", "image/png", 1024)
	if err != ErrAccessDenied {
		t.Errorf("expected ErrAccessDenied, got %v", err)
	}
}

func TestRequestUploadURLCaseNotOpen(t *testing.T) {
	imgRepo, caseRepo, storage := setupTestRepos()
	patientID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	seedClosedCase(caseRepo, caseID, tenantID, patientID)

	cmd := NewRequestUploadURLCommand(imgRepo, caseRepo, storage)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	_, _, _, err := cmd.Execute(context.Background(), p, caseID, "xray.png", "image/png", 1024)
	if err != ErrCaseNotOpen {
		t.Errorf("expected ErrCaseNotOpen, got %v", err)
	}
}

func TestRequestUploadURLTooLarge(t *testing.T) {
	imgRepo, caseRepo, storage := setupTestRepos()
	patientID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	seedOpenCase(caseRepo, caseID, tenantID, patientID)

	cmd := NewRequestUploadURLCommand(imgRepo, caseRepo, storage)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	_, _, _, err := cmd.Execute(context.Background(), p, caseID, "xray.png", "image/png", maxUploadSizeBytes+1)
	if err != ErrFileTooLarge {
		t.Errorf("expected ErrFileTooLarge, got %v", err)
	}
}

func TestConfirmUploadSuccess(t *testing.T) {
	imgRepo, caseRepo, _ := setupTestRepos()
	patientID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	seedOpenCase(caseRepo, caseID, tenantID, patientID)
	imgID := uuid.New()
	seedImage(imgRepo, imgID, tenantID, patientID, caseID)

	cmd := NewConfirmUploadCommand(imgRepo, caseRepo)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	img, err := cmd.Execute(context.Background(), p, caseID, tenantID.String()+"/"+caseID.String()+"/"+imgID.String()+"/xray.png", "xray.png", "image/png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if img.CaseID != caseID {
		t.Errorf("CaseID = %v, want %v", img.CaseID, caseID)
	}
	if img.FileName != "xray.png" {
		t.Errorf("FileName = %v, want xray.png", img.FileName)
	}
}

func TestConfirmUploadDoctorRejected(t *testing.T) {
	imgRepo, caseRepo, _ := setupTestRepos()
	doctorID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	seedOpenCase(caseRepo, caseID, tenantID, uuid.New())

	cmd := NewConfirmUploadCommand(imgRepo, caseRepo)
	p := newPrincipal(doctorID, tenantID, sharedauth.RoleDoctor)

	_, err := cmd.Execute(context.Background(), p, caseID, "key", "xray.png", "image/png")
	if err != ErrInvalidRole {
		t.Errorf("expected ErrInvalidRole, got %v", err)
	}
}

func TestListImagesSuccess(t *testing.T) {
	imgRepo, caseRepo, _ := setupTestRepos()
	patientID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	seedOpenCase(caseRepo, caseID, tenantID, patientID)

	imgID1 := uuid.New()
	imgID2 := uuid.New()
	seedImage(imgRepo, imgID1, tenantID, patientID, caseID)
	seedImage(imgRepo, imgID2, tenantID, patientID, caseID)

	cmd := NewListImagesQuery(imgRepo, caseRepo)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	images, err := cmd.Execute(context.Background(), p, caseID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(images) != 2 {
		t.Errorf("len = %d, want 2", len(images))
	}
}

func TestListImagesAccessDenied(t *testing.T) {
	imgRepo, caseRepo, _ := setupTestRepos()
	otherPatient := uuid.New()
	patientID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	seedOpenCase(caseRepo, caseID, tenantID, otherPatient)

	cmd := NewListImagesQuery(imgRepo, caseRepo)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	_, err := cmd.Execute(context.Background(), p, caseID)
	if err != ErrAccessDenied {
		t.Errorf("expected ErrAccessDenied, got %v", err)
	}
}

func TestGetDownloadURLSuccess(t *testing.T) {
	imgRepo, caseRepo, storage := setupTestRepos()
	patientID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	seedOpenCase(caseRepo, caseID, tenantID, patientID)

	imgID := uuid.New()
	seedImage(imgRepo, imgID, tenantID, patientID, caseID)

	cmd := NewGetDownloadURLCommand(imgRepo, caseRepo, storage)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	url, expiresIn, err := cmd.Execute(context.Background(), p, imgID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url == "" {
		t.Error("expected non-empty download URL")
	}
	if expiresIn != 900 {
		t.Errorf("expiresIn = %d, want 900", expiresIn)
	}
}

func TestGetDownloadURLImageNotFound(t *testing.T) {
	imgRepo, caseRepo, storage := setupTestRepos()
	patientID := uuid.New()
	tenantID := uuid.New()

	cmd := NewGetDownloadURLCommand(imgRepo, caseRepo, storage)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	_, _, err := cmd.Execute(context.Background(), p, uuid.New())
	if err != ErrImageNotFound {
		t.Errorf("expected ErrImageNotFound, got %v", err)
	}
}

func TestGetDownloadURLAccessDenied(t *testing.T) {
	imgRepo, caseRepo, storage := setupTestRepos()
	otherPatient := uuid.New()
	patientID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	seedOpenCase(caseRepo, caseID, tenantID, otherPatient)

	imgID := uuid.New()
	seedImage(imgRepo, imgID, tenantID, otherPatient, caseID)

	cmd := NewGetDownloadURLCommand(imgRepo, caseRepo, storage)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	_, _, err := cmd.Execute(context.Background(), p, imgID)
	if err != ErrAccessDenied {
		t.Errorf("expected ErrAccessDenied, got %v", err)
	}
}
