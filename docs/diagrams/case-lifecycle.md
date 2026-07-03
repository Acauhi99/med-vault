# Case Lifecycle

## Status Transitions

```mermaid
stateDiagram-v2
    [*] --> Open: Patient creates case

    Open --> Assigned: Admin assigns doctor
    Assigned --> Diagnosed: Doctor writes diagnosis
    Diagnosed --> Closed: Admin closes case

    Open --> Open: Patient adds symptoms
    Open --> Open: Patient uploads images

    Assigned --> Assigned: Doctor reviews images
    Assigned --> Assigned: Doctor reviews symptoms

    Closed --> [*]
```

## State Transition Rules

| From | To | Trigger | Actor | Preconditions |
|------|----|---------|-------|---------------|
| — | Open | CreateCase | Patient | Authenticated, valid tenant |
| Open | Assigned | AssignDoctor | Admin | Case has symptoms; doctor belongs to same tenant |
| Assigned | Diagnosed | WriteDiagnosis | Doctor | Doctor is assigned to case |
| Diagnosed | Closed | CloseCase | Admin | Diagnosis exists |
| Open | Open | AddSymptom | Patient | Case belongs to patient, status = Open |
| Open | Open | UploadImage | Patient | Case belongs to patient, status = Open |

## Invalid Transitions (Blocked)

```mermaid
flowchart LR
    A[Open] -->|WriteDiagnosis| B[❌ BLOCKED]
    C[Assigned] -->|AssignDoctor| D[❌ BLOCKED]
    E[Diagnosed] -->|AddSymptom| F[❌ BLOCKED]
    G[Closed] -->|Any mutation| H[❌ BLOCKED]

    style B fill:#ff6666
    style D fill:#ff6666
    style F fill:#ff6666
    style H fill:#ff6666
```

## Full Case Flow (All Actors)

```mermaid
sequenceDiagram
    actor Patient
    actor Admin
    actor Doctor
    participant Backend
    participant DB

    rect rgb(240, 248, 255)
        Note over Patient, DB: Phase 1 — Case Creation (Patient)
        Patient->>Backend: POST /api/v1/cases {symptoms}
        Backend->>DB: INSERT case (status = Open)
        Backend->>DB: INSERT symptoms
        Backend->>DB: INSERT audit_log (CaseCreated)
        Backend-->>Patient: 201 Created
    end

    rect rgb(255, 248, 240)
        Note over Patient, DB: Phase 2 — Image Upload (Patient)
        Patient->>Backend: POST /api/v1/cases/{id}/images/upload-url
        Backend-->>Patient: Pre-signed URL
        Patient->>Backend: PUT image to S3 via pre-signed URL
        Patient->>Backend: POST /api/v1/cases/{id}/images {metadata}
        Backend->>DB: INSERT image
        Backend->>DB: INSERT audit_log (ImageUploaded)
        Backend-->>Patient: 201 Created
    end

    rect rgb(240, 255, 240)
        Note over Admin, DB: Phase 3 — Doctor Assignment (Admin)
        Admin->>Backend: POST /api/v1/cases/{id}/assign {doctor_id}
        Backend->>DB: UPDATE case SET status = Assigned
        Backend->>DB: INSERT audit_log (DoctorAssigned)
        Backend-->>Admin: 200 OK
    end

    rect rgb(255, 240, 255)
        Note over Doctor, DB: Phase 4 — Diagnosis (Doctor)
        Doctor->>Backend: GET /api/v1/cases/{id}
        Backend-->>Doctor: Case + symptoms + images
        Doctor->>Backend: POST /api/v1/cases/{id}/diagnosis {notes}
        Backend->>DB: UPDATE case SET status = Diagnosed
        Backend->>DB: INSERT diagnosis
        Backend->>DB: INSERT audit_log (DiagnosisWritten)
        Backend-->>Doctor: 200 OK
    end

    rect rgb(255, 255, 240)
        Note over Admin, DB: Phase 5 — Case Closure (Admin)
        Admin->>Backend: POST /api/v1/cases/{id}/close
        Backend->>DB: UPDATE case SET status = Closed
        Backend->>DB: INSERT audit_log (CaseClosed)
        Backend-->>Admin: 200 OK
    end

    rect rgb(240, 240, 255)
        Note over Patient, DB: Phase 6 — View Result (Patient)
        Patient->>Backend: GET /api/v1/cases/{id}
        Backend-->>Patient: Case + symptoms + diagnosis
    end
```
