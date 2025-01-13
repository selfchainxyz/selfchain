# Identity Module Test Implementation Plan

## Overview

This document outlines the comprehensive test plan for the Self Chain Identity Module. The plan covers unit tests, integration tests, end-to-end tests, performance tests, and security validations.

## Table of Contents
- [1. Unit Tests](#1-unit-tests)
- [2. Integration Tests](#2-integration-tests)
- [3. End-to-End Tests](#3-end-to-end-tests)
- [4. Performance Tests](#4-performance-tests)
- [5. Security Tests](#5-security-tests)
- [6. CLI Tests](#6-cli-tests)

## 1. Unit Tests

### 1.1 Types Package (`/types`)

#### Basic Types Validation

##### DID Document Validation
```go
func TestDIDDocumentValidation(t *testing.T) {
    tests := []struct {
        name    string
        doc     types.DIDDocument
        wantErr bool
    }{
        {
            name: "valid DID document",
            doc: types.DIDDocument{
                Id:         "did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS",
                Controller: []string{"did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS"},
                VerificationMethod: []types.VerificationMethod{
                    {
                        Id:              "key-1",
                        Type:            "Ed25519VerificationKey2020",
                        Controller:      "did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS",
                        PublicKeyBase58: "H3C2AVvLMv6gmMNam3uVAjZpfkcJCwDwnZn6z3wXmqPV",
                    },
                },
            },
            wantErr: false,
        },
        {
            name: "invalid DID document - missing ID",
            doc: types.DIDDocument{
                Controller: []string{"did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS"},
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.doc.ValidateBasic()
            if (err != nil) != tt.wantErr {
                t.Errorf("DIDDocument.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

##### Credential Validation
```go
func TestCredentialValidation(t *testing.T) {
    tests := []struct {
        name    string
        cred    types.Credential
        wantErr bool
    }{
        {
            name: "valid credential",
            cred: types.Credential{
                Id:     "cred-1",
                Type:   []string{"VerifiableCredential", "KYCCredential"},
                Issuer: "did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS",
                Subject: map[string]string{
                    "id":   "did:selfchain:1234567890",
                    "name": "John Doe",
                    "age":  "25",
                },
                Status: types.Status_STATUS_ACTIVE,
            },
            wantErr: false,
        },
    }
    // Test implementation
}
```

### 1.2 Keeper Package (`/keeper`)

#### DID Operations
```go
func TestDIDOperations(t *testing.T) {
    keeper, ctx := setupKeeper(t)
    
    // Test StoreDID
    did := types.DIDDocument{
        Id: "did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS",
        // ... other fields
    }
    
    err := keeper.StoreDID(ctx, did)
    require.NoError(t, err)
    
    // Test GetDID
    storedDID, found := keeper.GetDID(ctx, did.Id)
    require.True(t, found)
    require.Equal(t, did, storedDID)
    
    // Test UpdateDID
    did.Controller = append(did.Controller, "did:selfchain:newcontroller")
    err = keeper.UpdateDID(ctx, did)
    require.NoError(t, err)
    
    // Test DeactivateDID
    err = keeper.DeactivateDID(ctx, did.Id)
    require.NoError(t, err)
    
    deactivatedDID, found := keeper.GetDID(ctx, did.Id)
    require.True(t, found)
    require.Equal(t, types.Status_STATUS_DISABLED, deactivatedDID.Status)
}
```

#### OAuth & Social Identity
```go
func TestOAuthVerification(t *testing.T) {
    keeper, ctx := setupKeeper(t)
    
    // Test OAuth token verification
    token := "valid-oauth-token"
    provider := "google"
    
    socialID, err := keeper.VerifyOAuthToken(ctx, provider, token)
    require.NoError(t, err)
    require.NotEmpty(t, socialID)
    
    // Test social identity linking
    did := "did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS"
    err = keeper.LinkSocialIdentity(ctx, did, socialID, provider)
    require.NoError(t, err)
    
    // Test getting linked DID
    linkedDID, found := keeper.GetLinkedDID(ctx, provider, socialID)
    require.True(t, found)
    require.Equal(t, did, linkedDID)
}
```

### 1.3 ZKP Operations Testing

#### Test ZKP Generation and Verification
```go
func TestZKPOperations(t *testing.T) {
    // Test data
    claims := map[string]string{
        "name":     "John Doe",
        "age":      "25",
        "country":  "US",
        "verified": "true",
    }
    
    disclosedClaims := []string{"name", "verified"}
    verificationKey := "test-verification-key"
    
    // Generate proof
    proof, err := types.GenerateZKProof(claims, disclosedClaims, verificationKey)
    require.NoError(t, err)
    require.NotNil(t, proof)
    
    // Verify the proof structure
    err = proof.ValidateBasic()
    require.NoError(t, err)
    
    // Verify with disclosed claims
    disclosed := map[string]string{
        "name":     "John Doe",
        "verified": "true",
    }
    
    valid, err := types.VerifyZKProof(proof, disclosed)
    require.NoError(t, err)
    require.True(t, valid)
    
    // Test invalid proof
    invalidDisclosed := map[string]string{
        "name":     "Jane Doe", // Wrong value
        "verified": "true",
    }
    
    valid, err = types.VerifyZKProof(proof, invalidDisclosed)
    require.NoError(t, err)
    require.False(t, valid)
}
```

#### Test Selective Disclosure
```go
func TestSelectiveDisclosure(t *testing.T) {
    credential := types.Credential{
        Id:     "cred-1",
        Type:   []string{"VerifiableCredential", "KYCCredential"},
        Issuer: "did:selfchain:issuer",
        Subject: map[string]string{
            "id":       "did:selfchain:subject",
            "name":     "John Doe",
            "age":      "25",
            "country":  "US",
            "verified": "true",
        },
    }
    
    // Test different disclosure patterns
    testCases := []struct {
        name             string
        disclosedClaims  []string
        expectedSuccess  bool
        verificationFunc func(proof *types.ZKProof) bool
    }{
        {
            name:            "disclose basic info",
            disclosedClaims: []string{"name", "verified"},
            expectedSuccess: true,
            verificationFunc: func(proof *types.ZKProof) bool {
                return len(proof.DisclosedIndices) == 2
            },
        },
        {
            name:            "disclose sensitive info",
            disclosedClaims: []string{"age", "country"},
            expectedSuccess: true,
            verificationFunc: func(proof *types.ZKProof) bool {
                return len(proof.DisclosedIndices) == 2
            },
        },
        {
            name:            "disclose all info",
            disclosedClaims: []string{"name", "age", "country", "verified"},
            expectedSuccess: true,
            verificationFunc: func(proof *types.ZKProof) bool {
                return len(proof.DisclosedIndices) == 4
            },
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            proof, err := types.GenerateZKProof(
                credential.Subject,
                tc.disclosedClaims,
                "test-key",
            )
            
            if tc.expectedSuccess {
                require.NoError(t, err)
                require.True(t, tc.verificationFunc(proof))
            } else {
                require.Error(t, err)
            }
        })
    }
}
```

### 1.4 MFA Implementation Tests

#### Test TOTP Generation and Validation
```go
func TestTOTPOperations(t *testing.T) {
    keeper, ctx := setupKeeper(t)
    
    // Test TOTP setup
    did := "did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS"
    
    // Add MFA method
    method := types.MFAMethod{
        Type:   "totp",
        Status: types.Status_STATUS_ACTIVE,
    }
    
    err := keeper.AddMFAMethod(ctx, did, method)
    require.NoError(t, err)
    
    // Generate TOTP challenge
    challenge, err := keeper.CreateMFAChallenge(ctx, did, "totp")
    require.NoError(t, err)
    require.NotNil(t, challenge)
    
    // Verify TOTP code
    validCode := "123456" // Mock valid TOTP code
    err = keeper.VerifyMFAChallenge(ctx, did, "totp", validCode)
    require.NoError(t, err)
    
    // Test invalid code
    invalidCode := "999999"
    err = keeper.VerifyMFAChallenge(ctx, did, "totp", invalidCode)
    require.Error(t, err)
}
```

#### Test MFA Recovery Flow
```go
func TestMFARecovery(t *testing.T) {
    keeper, ctx := setupKeeper(t)
    
    did := "did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS"
    
    // Setup MFA with backup codes
    method := types.MFAMethod{
        Type:   "backup_codes",
        Status: types.Status_STATUS_ACTIVE,
        BackupCodes: []string{
            "code1", "code2", "code3",
        },
    }
    
    err := keeper.AddMFAMethod(ctx, did, method)
    require.NoError(t, err)
    
    // Test recovery with backup code
    err = keeper.VerifyMFAMethod(ctx, did, "backup_codes", "code1")
    require.NoError(t, err)
    
    // Verify backup code is consumed
    config, found := keeper.GetMFAConfig(ctx, did)
    require.True(t, found)
    
    method, err = keeper.GetMFAMethod(ctx, did, "backup_codes")
    require.NoError(t, err)
    require.NotContains(t, method.BackupCodes, "code1")
    
    // Test invalid backup code
    err = keeper.VerifyMFAMethod(ctx, did, "backup_codes", "invalid")
    require.Error(t, err)
}
```

### 1.5 Rate Limiting Tests

#### Test Rate Limiting Implementation
```go
func TestRateLimiting(t *testing.T) {
    keeper, ctx := setupKeeper(t)
    
    did := "did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS"
    action := "mfa_verify"
    
    // Test within limits
    for i := 0; i < 5; i++ {
        allowed, err := keeper.CheckRateLimit(ctx, did, action)
        require.NoError(t, err)
        require.True(t, allowed)
        
        err = keeper.UpdateRateLimit(ctx, did, action)
        require.NoError(t, err)
    }
    
    // Test exceeding limits
    allowed, err := keeper.CheckRateLimit(ctx, did, action)
    require.NoError(t, err)
    require.False(t, allowed)
    
    // Test reset after window
    ctx = ctx.WithBlockTime(ctx.BlockTime().Add(1 * time.Hour))
    allowed, err = keeper.CheckRateLimit(ctx, did, action)
    require.NoError(t, err)
    require.True(t, allowed)
}
```

### 1.6 Audit Logging Tests

#### Test Audit Trail
```go
func TestAuditLogging(t *testing.T) {
    keeper, ctx := setupKeeper(t)
    
    did := "did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS"
    
    // Log various events
    events := []struct {
        action string
        data   map[string]string
    }{
        {
            action: "did_create",
            data: map[string]string{
                "did": did,
            },
        },
        {
            action: "mfa_enable",
            data: map[string]string{
                "did":        did,
                "mfa_type":   "totp",
                "mfa_status": "active",
            },
        },
        {
            action: "credential_issue",
            data: map[string]string{
                "did":            did,
                "credential_id":  "cred-1",
                "credential_type": "KYCCredential",
            },
        },
    }
    
    for _, event := range events {
        err := keeper.LogAuditEvent(ctx, did, event.action, event.data)
        require.NoError(t, err)
    }
    
    // Test audit log retrieval
    logs, err := keeper.GetAuditLogs(ctx, did)
    require.NoError(t, err)
    require.Len(t, logs, len(events))
    
    // Test audit log filtering
    filtered, err := keeper.FilterAuditLogs(ctx, did, "mfa_enable")
    require.NoError(t, err)
    require.Len(t, filtered, 1)
    require.Equal(t, "mfa_enable", filtered[0].Action)
}
```

## 2. Integration Tests

### 2.1 DID Lifecycle
```go
func TestCompleteDIDLifecycle(t *testing.T) {
    app := setupTestApp(t)
    ctx := app.BaseApp.NewContext(false, tmproto.Header{})
    
    // Create DID
    msg := types.MsgCreateDID{
        Creator: testAccounts[0].Address.String(),
        VerificationMethod: []types.VerificationMethod{
            {
                Id:              "key-1",
                Type:            "Ed25519VerificationKey2020",
                PublicKeyBase58: "H3C2AVvLMv6gmMNam3uVAjZpfkcJCwDwnZn6z3wXmqPV",
            },
        },
    }
    
    resp, err := app.IdentityKeeper.CreateDID(sdk.WrapSDKContext(ctx), &msg)
    require.NoError(t, err)
    require.NotEmpty(t, resp.Did)
    
    // Update DID
    updateMsg := types.MsgUpdateDID{
        Creator: testAccounts[0].Address.String(),
        Did:     resp.Did,
        Service: []types.Service{
            {
                Id:              "service-1",
                Type:            "LinkedDomains",
                ServiceEndpoint: "https://example.com",
            },
        },
    }
    
    _, err = app.IdentityKeeper.UpdateDID(sdk.WrapSDKContext(ctx), &updateMsg)
    require.NoError(t, err)
    
    // Verify DID state
    did, found := app.IdentityKeeper.GetDID(ctx, resp.Did)
    require.True(t, found)
    require.Equal(t, 1, len(did.Service))
}
```

### 2.2 OAuth & Identity Flow
```go
func TestOAuthIdentityFlow(t *testing.T) {
    app := setupTestApp(t)
    ctx := app.BaseApp.NewContext(false, tmproto.Header{})
    
    // Setup mock OAuth provider
    mockProvider := setupMockOAuthProvider(t)
    
    // Verify OAuth token
    msg := types.MsgVerifyOAuthToken{
        Creator:  testAccounts[0].Address.String(),
        Provider: "google",
        Token:    "mock-token",
    }
    
    resp, err := app.IdentityKeeper.VerifyOAuthToken(sdk.WrapSDKContext(ctx), &msg)
    require.NoError(t, err)
    require.NotEmpty(t, resp.Id)
    
    // Link social identity to DID
    linkMsg := types.MsgLinkSocialIdentity{
        Creator:   testAccounts[0].Address.String(),
        Did:       "did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS",
        Provider:  "google",
        SocialId:  resp.Id,
    }
    
    _, err = app.IdentityKeeper.LinkSocialIdentity(sdk.WrapSDKContext(ctx), &linkMsg)
    require.NoError(t, err)
}
```

## 3. End-to-End Tests

### 3.1 Complete Identity Workflow
```go
func TestCompleteIdentityWorkflow(t *testing.T) {
    app := setupTestApp(t)
    ctx := app.BaseApp.NewContext(false, tmproto.Header{})
    
    // 1. Create DID with OAuth
    didResp := createDIDWithOAuth(t, app, ctx)
    
    // 2. Setup MFA
    setupMFA(t, app, ctx, didResp.Did)
    
    // 3. Issue credential
    cred := issueCredential(t, app, ctx, didResp.Did)
    
    // 4. Create ZKP presentation
    zkp := createZKPPresentation(t, app, ctx, cred)
    
    // 5. Verify presentation
    verifyZKPPresentation(t, app, ctx, zkp)
}

func createDIDWithOAuth(t *testing.T, app *app.App, ctx sdk.Context) *types.MsgCreateDIDResponse {
    // Implementation
}

func setupMFA(t *testing.T, app *app.App, ctx sdk.Context, did string) {
    // Implementation
}

func issueCredential(t *testing.T, app *app.App, ctx sdk.Context, did string) *types.Credential {
    // Implementation
}
```

## 4. Performance Tests

### 4.1 Load Testing
```go
func TestConcurrentOperations(t *testing.T) {
    app := setupTestApp(t)
    ctx := app.BaseApp.NewContext(false, tmproto.Header{})
    
    const numOperations = 100
    var wg sync.WaitGroup
    
    // Test concurrent DID creations
    for i := 0; i < numOperations; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            msg := types.MsgCreateDID{
                Creator: testAccounts[i%len(testAccounts)].Address.String(),
            }
            _, err := app.IdentityKeeper.CreateDID(sdk.WrapSDKContext(ctx), &msg)
            require.NoError(t, err)
        }()
    }
    
    wg.Wait()
}
```

## 5. Security Tests

### 5.1 Access Control Tests
```go
func TestAccessControl(t *testing.T) {
    app := setupTestApp(t)
    ctx := app.BaseApp.NewContext(false, tmproto.Header{})
    
    // Test unauthorized access
    msg := types.MsgUpdateDID{
        Creator: testAccounts[1].Address.String(), // Non-owner
        Did:     "did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS",
    }
    
    _, err := app.IdentityKeeper.UpdateDID(sdk.WrapSDKContext(ctx), &msg)
    require.Error(t, err)
    require.Contains(t, err.Error(), "unauthorized")
}
```

## 6. CLI Tests

### 6.1 Query Commands
```go
func TestQueryCommands(t *testing.T) {
    net := network.New(t)
    
    // Test DID query
    val := net.Validators[0]
    didQuery := fmt.Sprintf(`%s query identity did %s`, val.ClientCtx.BinaryName, testDID)
    
    out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, didQuery)
    require.NoError(t, err)
    
    var resp types.QueryDIDResponse
    require.NoError(t, val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp))
    require.NotNil(t, resp.DidDocument)
}
```

## Implementation Priority

1. Core Types Validation Tests
   - Basic validation for all types
   - Message validation
   - ZKP operations

2. Keeper Basic Operations Tests
   - DID operations
   - OAuth operations
   - MFA operations
   - Credential operations

3. Integration Tests
   - DID lifecycle
   - OAuth flow
   - Credential issuance and verification

4. Security & Performance Tests
   - Access control
   - Rate limiting
   - Concurrent operations

5. CLI & Query Tests
   - Command validation
   - Query responses
   - Error handling

## Test Coverage Goals

- Types Package: 90%
- Keeper Package: 85%
- CLI Package: 80%
- Overall Module: 85%

## Running Tests

```bash
# Run all tests
go test ./x/identity/... -v

# Run specific test
go test ./x/identity/keeper -run TestDIDOperations -v

# Run with coverage
go test ./x/identity/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Contributing

When adding new tests:
1. Follow the existing test structure
2. Include both positive and negative test cases
3. Add proper documentation
4. Ensure proper cleanup in teardown
5. Maintain idempotency

#### Edge Cases for ZKP Operations
```go
func TestZKPEdgeCases(t *testing.T) {
    testCases := []struct {
        name          string
        claims        map[string]string
        disclosed     []string
        expectErr     bool
        errorContains string
    }{
        {
            name: "empty claims",
            claims: map[string]string{},
            disclosed: []string{"name"},
            expectErr: true,
            errorContains: "empty claims",
        },
        {
            name: "nil claims",
            claims: nil,
            disclosed: []string{"name"},
            expectErr: true,
            errorContains: "nil claims",
        },
        {
            name: "disclosed claim not in claims",
            claims: map[string]string{
                "name": "John",
            },
            disclosed: []string{"age"},
            expectErr: true,
            errorContains: "claim not found",
        },
        {
            name: "extremely large claim value",
            claims: map[string]string{
                "name": strings.Repeat("a", 1000000),
            },
            disclosed: []string{"name"},
            expectErr: true,
            errorContains: "claim value too large",
        },
        {
            name: "special characters in claims",
            claims: map[string]string{
                "name": "John@#$%^&*()",
                "age":  "25",
            },
            disclosed: []string{"name"},
            expectErr: false,
        },
        {
            name: "unicode characters in claims",
            claims: map[string]string{
                "name": "जॉन डो",
                "age":  "२५",
            },
            disclosed: []string{"name", "age"},
            expectErr: false,
        },
        {
            name: "empty disclosed claims",
            claims: map[string]string{
                "name": "John",
            },
            disclosed: []string{},
            expectErr: true,
            errorContains: "no claims to disclose",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            proof, err := types.GenerateZKProof(tc.claims, tc.disclosed, "test-key")
            
            if tc.expectErr {
                require.Error(t, err)
                require.Contains(t, err.Error(), tc.errorContains)
            } else {
                require.NoError(t, err)
                require.NotNil(t, proof)
            }
        })
    }
}
```

#### Edge Cases for MFA Operations
```go
func TestMFAEdgeCases(t *testing.T) {
    keeper, ctx := setupKeeper(t)
    did := "did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS"

    testCases := []struct {
        name          string
        setup         func() error
        test         func() error
        expectErr     bool
        errorContains string
    }{
        {
            name: "multiple MFA methods exceed limit",
            setup: func() error {
                // Try to add more than allowed MFA methods
                for i := 0; i < 6; i++ { // Max is 5
                    method := types.MFAMethod{
                        Type:   fmt.Sprintf("method_%d", i),
                        Status: types.Status_STATUS_ACTIVE,
                    }
                    if err := keeper.AddMFAMethod(ctx, did, method); err != nil {
                        return err
                    }
                }
                return nil
            },
            expectErr: true,
            errorContains: "maximum MFA methods exceeded",
        },
        {
            name: "concurrent MFA challenges",
            setup: func() error {
                method := types.MFAMethod{
                    Type:   "totp",
                    Status: types.Status_STATUS_ACTIVE,
                }
                return keeper.AddMFAMethod(ctx, did, method)
            },
            test: func() error {
                // Create multiple challenges concurrently
                var wg sync.WaitGroup
                errs := make(chan error, 5)
                
                for i := 0; i < 5; i++ {
                    wg.Add(1)
                    go func() {
                        defer wg.Done()
                        _, err := keeper.CreateMFAChallenge(ctx, did, "totp")
                        if err != nil {
                            errs <- err
                        }
                    }()
                }
                
                wg.Wait()
                close(errs)
                
                if len(errs) > 0 {
                    return <-errs
                }
                return nil
            },
            expectErr: true,
            errorContains: "concurrent challenge creation not allowed",
        },
        {
            name: "expired TOTP challenge",
            setup: func() error {
                method := types.MFAMethod{
                    Type:   "totp",
                    Status: types.Status_STATUS_ACTIVE,
                }
                if err := keeper.AddMFAMethod(ctx, did, method); err != nil {
                    return err
                }
                
                _, err := keeper.CreateMFAChallenge(ctx, did, "totp")
                return err
            },
            test: func() error {
                // Move time forward beyond expiry
                ctx = ctx.WithBlockTime(ctx.BlockTime().Add(11 * time.Minute))
                return keeper.VerifyMFAChallenge(ctx, did, "totp", "123456")
            },
            expectErr: true,
            errorContains: "challenge expired",
        },
        {
            name: "reuse backup code",
            setup: func() error {
                method := types.MFAMethod{
                    Type:   "backup_codes",
                    Status: types.Status_STATUS_ACTIVE,
                    BackupCodes: []string{"code1"},
                }
                if err := keeper.AddMFAMethod(ctx, did, method); err != nil {
                    return err
                }
                
                // Use the backup code once
                return keeper.VerifyMFAMethod(ctx, did, "backup_codes", "code1")
            },
            test: func() error {
                // Try to use the same backup code again
                return keeper.VerifyMFAMethod(ctx, did, "backup_codes", "code1")
            },
            expectErr: true,
            errorContains: "backup code already used",
        },
        {
            name: "deactivated MFA method",
            setup: func() error {
                method := types.MFAMethod{
                    Type:   "totp",
                    Status: types.Status_STATUS_DISABLED,
                }
                return keeper.AddMFAMethod(ctx, did, method)
            },
            test: func() error {
                return keeper.VerifyMFAMethod(ctx, did, "totp", "123456")
            },
            expectErr: true,
            errorContains: "MFA method is disabled",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            if tc.setup != nil {
                err := tc.setup()
                if tc.expectErr && err != nil {
                    require.Contains(t, err.Error(), tc.errorContains)
                    return
                }
                require.NoError(t, err)
            }

            if tc.test != nil {
                err := tc.test()
                if tc.expectErr {
                    require.Error(t, err)
                    require.Contains(t, err.Error(), tc.errorContains)
                } else {
                    require.NoError(t, err)
                }
            }
        })
    }
}
```

#### Edge Cases for Rate Limiting
```go
func TestRateLimitingEdgeCases(t *testing.T) {
    keeper, ctx := setupKeeper(t)
    did := "did:selfchain:14zCHv2Pq6aYJwDJTE8MYKQyJ3LhNEtKaS"

    testCases := []struct {
        name          string
        setup         func() error
        test         func() (bool, error)
        expectAllowed bool
        expectErr     bool
        errorContains string
    }{
        {
            name: "burst requests",
            test: func() (bool, error) {
                // Simulate burst of requests
                results := make([]bool, 100)
                var wg sync.WaitGroup
                var mu sync.Mutex
                
                for i := 0; i < 100; i++ {
                    wg.Add(1)
                    go func(idx int) {
                        defer wg.Done()
                        allowed, _ := keeper.CheckRateLimit(ctx, did, "burst_action")
                        mu.Lock()
                        results[idx] = allowed
                        mu.Unlock()
                    }(i)
                }
                
                wg.Wait()
                
                // Count allowed requests
                allowed := 0
                for _, r := range results {
                    if r {
                        allowed++
                    }
                }
                
                // Should not exceed rate limit
                return allowed <= 5, nil
            },
            expectAllowed: true,
        },
        {
            name: "different actions same window",
            setup: func() error {
                // Max out one action
                for i := 0; i < 5; i++ {
                    if err := keeper.UpdateRateLimit(ctx, did, "action1"); err != nil {
                        return err
                    }
                }
                return nil
            },
            test: func() (bool, error) {
                // Try different action
                return keeper.CheckRateLimit(ctx, did, "action2")
            },
            expectAllowed: true,
        },
        {
            name: "window reset edge",
            setup: func() error {
                // Max out the limit
                for i := 0; i < 5; i++ {
                    if err := keeper.UpdateRateLimit(ctx, did, "edge_action"); err != nil {
                        return err
                    }
                }
                
                // Move time to just before window reset
                ctx = ctx.WithBlockTime(ctx.BlockTime().Add(59 * time.Minute + 59 * time.Second))
                return nil
            },
            test: func() (bool, error) {
                return keeper.CheckRateLimit(ctx, did, "edge_action")
            },
            expectAllowed: false,
        },
        {
            name: "malformed action name",
            test: func() (bool, error) {
                return keeper.CheckRateLimit(ctx, did, strings.Repeat("a", 1000))
            },
            expectErr: true,
            errorContains: "action name too long",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            if tc.setup != nil {
                err := tc.setup()
                require.NoError(t, err)
            }

            allowed, err := tc.test()
            
            if tc.expectErr {
                require.Error(t, err)
                require.Contains(t, err.Error(), tc.errorContains)
            } else {
                require.NoError(t, err)
                require.Equal(t, tc.expectAllowed, allowed)
            }
        })
    }
}
```

#### Edge Cases for DID Operations
```go
func TestDIDOperationsEdgeCases(t *testing.T) {
    keeper, ctx := setupKeeper(t)

    testCases := []struct {
        name          string
        setup         func() (string, error)
        test         func(string) error
        expectErr     bool
        errorContains string
    }{
        {
            name: "circular controller reference",
            setup: func() (string, error) {
                // Create first DID
                did1 := types.DIDDocument{
                    Id:         "did:selfchain:1",
                    Controller: []string{"did:selfchain:2"}, // Points to second DID
                }
                if err := keeper.StoreDID(ctx, did1); err != nil {
                    return "", err
                }

                // Create second DID pointing back to first
                did2 := types.DIDDocument{
                    Id:         "did:selfchain:2",
                    Controller: []string{"did:selfchain:1"}, // Points back to first DID
                }
                if err := keeper.StoreDID(ctx, did2); err != nil {
                    return "", err
                }

                return did1.Id, nil
            },
            test: func(did string) error {
                // Try to verify controller chain
                return keeper.VerifyControllerChain(ctx, did)
            },
            expectErr: true,
            errorContains: "circular controller reference detected",
        },
        {
            name: "deep controller chain",
            setup: func() (string, error) {
                // Create a chain of 11 DIDs (exceeding max depth of 10)
                var lastDID string
                for i := 1; i <= 11; i++ {
                    did := types.DIDDocument{
                        Id: fmt.Sprintf("did:selfchain:%d", i),
                    }
                    if i < 11 {
                        did.Controller = []string{fmt.Sprintf("did:selfchain:%d", i+1)}
                    }
                    if err := keeper.StoreDID(ctx, did); err != nil {
                        return "", err
                    }
                    lastDID = did.Id
                }
                return "did:selfchain:1", nil
            },
            test: func(did string) error {
                return keeper.VerifyControllerChain(ctx, did)
            },
            expectErr: true,
            errorContains: "controller chain too deep",
        },
        {
            name: "concurrent DID updates",
            setup: func() (string, error) {
                did := types.DIDDocument{
                    Id: "did:selfchain:concurrent",
                    VerificationMethod: []types.VerificationMethod{
                        {
                            Id:              "key-1",
                            Type:            "Ed25519VerificationKey2020",
                            Controller:      "did:selfchain:concurrent",
                            PublicKeyBase58: "initial",
                        },
                    },
                }
                return did.Id, keeper.StoreDID(ctx, did)
            },
            test: func(did string) error {
                var wg sync.WaitGroup
                errs := make(chan error, 5)

                // Try to update the same DID concurrently
                for i := 0; i < 5; i++ {
                    wg.Add(1)
                    go func(idx int) {
                        defer wg.Done()
                        doc, found := keeper.GetDID(ctx, did)
                        if !found {
                            errs <- fmt.Errorf("DID not found")
                            return
                        }

                        doc.VerificationMethod[0].PublicKeyBase58 = fmt.Sprintf("key-%d", idx)
                        if err := keeper.UpdateDID(ctx, doc); err != nil {
                            errs <- err
                        }
                    }(i)
                }

                wg.Wait()
                close(errs)

                if len(errs) > 0 {
                    return <-errs
                }
                return nil
            },
            expectErr: true,
            errorContains: "concurrent modification detected",
        },
        {
            name: "malformed DID document",
            test: func(string) error {
                did := types.DIDDocument{
                    Id: "did:selfchain:malformed",
                    VerificationMethod: []types.VerificationMethod{
                        {
                            Id:              strings.Repeat("a", 1000), // Too long ID
                            Type:            "InvalidType",
                            Controller:      "did:selfchain:malformed",
                            PublicKeyBase58: "invalid base58",
                        },
                    },
                }
                return keeper.StoreDID(ctx, did)
            },
            expectErr: true,
            errorContains: "invalid verification method",
        },
        {
            name: "deactivate non-existent DID",
            test: func(string) error {
                return keeper.DeactivateDID(ctx, "did:selfchain:nonexistent")
            },
            expectErr: true,
            errorContains: "DID not found",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            var did string
            var err error

            if tc.setup != nil {
                did, err = tc.setup()
                require.NoError(t, err)
            }

            err = tc.test(did)
            if tc.expectErr {
                require.Error(t, err)
                require.Contains(t, err.Error(), tc.errorContains)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

#### Edge Cases for OAuth Flow
```go
func TestOAuthFlowEdgeCases(t *testing.T) {
    keeper, ctx := setupKeeper(t)

    testCases := []struct {
        name          string
        setup         func() error
        test         func() error
        expectErr     bool
        errorContains string
    }{
        {
            name: "expired oauth token",
            test: func() error {
                // Create an expired token
                token := &types.OAuthToken{
                    Token:     "expired_token",
                    Provider:  "google",
                    IssuedAt:  time.Now().Add(-2 * time.Hour).Unix(),
                    ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
                }
                return keeper.VerifyOAuthToken(ctx, token)
            },
            expectErr: true,
            errorContains: "token expired",
        },
        {
            name: "multiple social identities same provider",
            setup: func() error {
                // Link first social identity
                err := keeper.LinkSocialIdentity(ctx, "did:selfchain:1", types.SocialIdentity{
                    Provider:   "google",
                    ProviderId: "user1@gmail.com",
                })
                if err != nil {
                    return err
                }

                // Try to link another social identity from same provider
                return keeper.LinkSocialIdentity(ctx, "did:selfchain:1", types.SocialIdentity{
                    Provider:   "google",
                    ProviderId: "user2@gmail.com",
                })
            },
            expectErr: true,
            errorContains: "social identity for provider already exists",
        },
        {
            name: "link to multiple DIDs",
            setup: func() error {
                // Link social identity to first DID
                err := keeper.LinkSocialIdentity(ctx, "did:selfchain:1", types.SocialIdentity{
                    Provider:   "google",
                    ProviderId: "user@gmail.com",
                })
                if err != nil {
                    return err
                }

                // Try to link same social identity to another DID
                return keeper.LinkSocialIdentity(ctx, "did:selfchain:2", types.SocialIdentity{
                    Provider:   "google",
                    ProviderId: "user@gmail.com",
                })
            },
            expectErr: true,
            errorContains: "social identity already linked to another DID",
        },
        {
            name: "unsupported oauth provider",
            test: func() error {
                return keeper.VerifyOAuthToken(ctx, &types.OAuthToken{
                    Token:    "valid_token",
                    Provider: "unsupported_provider",
                })
            },
            expectErr: true,
            errorContains: "unsupported OAuth provider",
        },
        {
            name: "malformed oauth token",
            test: func() error {
                return keeper.VerifyOAuthToken(ctx, &types.OAuthToken{
                    Token:    strings.Repeat("a", 10000), // Too long token
                    Provider: "google",
                })
            },
            expectErr: true,
            errorContains: "invalid token format",
        },
        {
            name: "concurrent social identity linking",
            setup: func() error {
                var wg sync.WaitGroup
                errs := make(chan error, 5)

                // Try to link social identities concurrently
                for i := 0; i < 5; i++ {
                    wg.Add(1)
                    go func(idx int) {
                        defer wg.Done()
                        err := keeper.LinkSocialIdentity(ctx, fmt.Sprintf("did:selfchain:%d", idx),
                            types.SocialIdentity{
                                Provider:   "google",
                                ProviderId: "user@gmail.com",
                            })
                        if err != nil {
                            errs <- err
                        }
                    }(i)
                }

                wg.Wait()
                close(errs)

                if len(errs) > 0 {
                    return <-errs
                }
                return nil
            },
            expectErr: true,
            errorContains: "concurrent social identity linking",
        },
        {
            name: "rate limited oauth verification",
            setup: func() error {
                // Max out rate limit
                for i := 0; i < 10; i++ {
                    _ = keeper.VerifyOAuthToken(ctx, &types.OAuthToken{
                        Token:    fmt.Sprintf("token_%d", i),
                        Provider: "google",
                    })
                }
                return nil
            },
            test: func() error {
                return keeper.VerifyOAuthToken(ctx, &types.OAuthToken{
                    Token:    "another_token",
                    Provider: "google",
                })
            },
            expectErr: true,
            errorContains: "rate limit exceeded",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            if tc.setup != nil {
                err := tc.setup()
                if tc.expectErr && err != nil {
                    require.Contains(t, err.Error(), tc.errorContains)
                    return
                }
                require.NoError(t, err)
            }

            if tc.test != nil {
                err := tc.test()
                if tc.expectErr {
                    require.Error(t, err)
                    require.Contains(t, err.Error(), tc.errorContains)
                } else {
                    require.NoError(t, err)
                }
            }
        })
    }
}
```

#### Edge Cases for Specific OAuth Providers
```go
func TestSpecificOAuthProvidersEdgeCases(t *testing.T) {
    keeper, ctx := setupKeeper(t)

    testCases := []struct {
        name          string
        setup         func() error
        test         func() error
        expectErr     bool
        errorContains string
    }{
        // Google OAuth Edge Cases
        {
            name: "google invalid audience",
            test: func() error {
                token := &types.OAuthToken{
                    Token:    "valid_token",
                    Provider: "google",
                    Claims: map[string]interface{}{
                        "aud": "wrong_client_id",
                    },
                }
                return keeper.VerifyGoogleToken(ctx, token)
            },
            expectErr: true,
            errorContains: "invalid audience",
        },
        {
            name: "google email not verified",
            test: func() error {
                token := &types.OAuthToken{
                    Token:    "valid_token",
                    Provider: "google",
                    Claims: map[string]interface{}{
                        "aud": "correct_client_id",
                        "email_verified": false,
                    },
                }
                return keeper.VerifyGoogleToken(ctx, token)
            },
            expectErr: true,
            errorContains: "email not verified",
        },
        {
            name: "google hosted domain mismatch",
            test: func() error {
                token := &types.OAuthToken{
                    Token:    "valid_token",
                    Provider: "google",
                    Claims: map[string]interface{}{
                        "aud": "correct_client_id",
                        "email_verified": true,
                        "hd": "wrongdomain.com",
                    },
                }
                return keeper.VerifyGoogleToken(ctx, token)
            },
            expectErr: true,
            errorContains: "hosted domain mismatch",
        },

        // Facebook OAuth Edge Cases
        {
            name: "facebook invalid app secret proof",
            test: func() error {
                token := &types.OAuthToken{
                    Token:    "valid_token",
                    Provider: "facebook",
                    Claims: map[string]interface{}{
                        "appsecret_proof": "invalid_proof",
                    },
                }
                return keeper.VerifyFacebookToken(ctx, token)
            },
            expectErr: true,
            errorContains: "invalid app secret proof",
        },
        {
            name: "facebook required scope missing",
            test: func() error {
                token := &types.OAuthToken{
                    Token:    "valid_token",
                    Provider: "facebook",
                    Claims: map[string]interface{}{
                        "scopes": []string{"public_profile"},  // email scope missing
                    },
                }
                return keeper.VerifyFacebookToken(ctx, token)
            },
            expectErr: true,
            errorContains: "required scope missing",
        },
        {
            name: "facebook token debug data mismatch",
            test: func() error {
                token := &types.OAuthToken{
                    Token:    "valid_token",
                    Provider: "facebook",
                    Claims: map[string]interface{}{
                        "app_id": "wrong_app_id",
                    },
                }
                return keeper.VerifyFacebookToken(ctx, token)
            },
            expectErr: true,
            errorContains: "token debug data mismatch",
        },

        // Twitter OAuth Edge Cases
        {
            name: "twitter invalid signature",
            test: func() error {
                token := &types.OAuthToken{
                    Token:    "valid_token",
                    Provider: "twitter",
                    Claims: map[string]interface{}{
                        "oauth_signature": "invalid_signature",
                    },
                }
                return keeper.VerifyTwitterToken(ctx, token)
            },
            expectErr: true,
            errorContains: "invalid signature",
        },
        {
            name: "twitter expired nonce",
            test: func() error {
                token := &types.OAuthToken{
                    Token:    "valid_token",
                    Provider: "twitter",
                    Claims: map[string]interface{}{
                        "oauth_timestamp": fmt.Sprintf("%d", time.Now().Add(-2*time.Hour).Unix()),
                    },
                }
                return keeper.VerifyTwitterToken(ctx, token)
            },
            expectErr: true,
            errorContains: "expired nonce",
        },
    }

    // Test execution code...
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            if tc.setup != nil {
                err := tc.setup()
                require.NoError(t, err)
            }

            err := tc.test()
            if tc.expectErr {
                require.Error(t, err)
                require.Contains(t, err.Error(), tc.errorContains)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

#### Edge Cases for Credential Management
```go
func TestCredentialManagementEdgeCases(t *testing.T) {
    keeper, ctx := setupKeeper(t)
    issuerDID := "did:selfchain:issuer"
    holderDID := "did:selfchain:holder"

    testCases := []struct {
        name          string
        setup         func() error
        test         func() error
        expectErr     bool
        errorContains string
    }{
        {
            name: "credential with future issuance date",
            test: func() error {
                credential := types.VerifiableCredential{
                    Id:           "cred1",
                    Issuer:       issuerDID,
                    IssuanceDate: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
                }
                return keeper.IssueCredential(ctx, credential)
            },
            expectErr: true,
            errorContains: "future issuance date",
        },
        {
            name: "credential with expired validity period",
            test: func() error {
                credential := types.VerifiableCredential{
                    Id:           "cred2",
                    Issuer:       issuerDID,
                    IssuanceDate: time.Now().Format(time.RFC3339),
                    ExpirationDate: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
                }
                return keeper.IssueCredential(ctx, credential)
            },
            expectErr: true,
            errorContains: "credential expired",
        },
        {
            name: "revoked credential verification",
            setup: func() error {
                // Issue and then revoke a credential
                credential := types.VerifiableCredential{
                    Id:           "cred3",
                    Issuer:       issuerDID,
                    IssuanceDate: time.Now().Format(time.RFC3339),
                }
                if err := keeper.IssueCredential(ctx, credential); err != nil {
                    return err
                }
                return keeper.RevokeCredential(ctx, credential.Id, issuerDID)
            },
            test: func() error {
                return keeper.VerifyCredential(ctx, "cred3")
            },
            expectErr: true,
            errorContains: "credential revoked",
        },
        {
            name: "credential with invalid proof",
            test: func() error {
                credential := types.VerifiableCredential{
                    Id:           "cred4",
                    Issuer:       issuerDID,
                    IssuanceDate: time.Now().Format(time.RFC3339),
                    Proof: &types.Proof{
                        Type:               "Ed25519Signature2020",
                        VerificationMethod: "invalid_key",
                        SignatureValue:     "invalid_signature",
                    },
                }
                return keeper.VerifyCredential(ctx, credential.Id)
            },
            expectErr: true,
            errorContains: "invalid proof",
        },
        {
            name: "credential schema validation failure",
            test: func() error {
                credential := types.VerifiableCredential{
                    Id:           "cred5",
                    Issuer:       issuerDID,
                    IssuanceDate: time.Now().Format(time.RFC3339),
                    CredentialSchema: &types.CredentialSchema{
                        Id:   "schema1",
                        Type: "JsonSchemaValidator2018",
                    },
                    CredentialSubject: map[string]interface{}{
                        "age": "invalid_number", // should be integer
                    },
                }
                return keeper.IssueCredential(ctx, credential)
            },
            expectErr: true,
            errorContains: "schema validation failed",
        },
        {
            name: "credential status verification failure",
            setup: func() error {
                // Create a credential with a status that points to a non-existent registry
                credential := types.VerifiableCredential{
                    Id:           "cred6",
                    Issuer:       issuerDID,
                    IssuanceDate: time.Now().Format(time.RFC3339),
                    CredentialStatus: &types.CredentialStatus{
                        Id:   "status1",
                        Type: "RevocationList2020Status",
                    },
                }
                return keeper.IssueCredential(ctx, credential)
            },
            test: func() error {
                return keeper.VerifyCredentialStatus(ctx, "cred6")
            },
            expectErr: true,
            errorContains: "status verification failed",
        },
        {
            name: "concurrent credential updates",
            setup: func() error {
                // Issue initial credential
                credential := types.VerifiableCredential{
                    Id:           "cred7",
                    Issuer:       issuerDID,
                    IssuanceDate: time.Now().Format(time.RFC3339),
                }
                return keeper.IssueCredential(ctx, credential)
            },
            test: func() error {
                var wg sync.WaitGroup
                errs := make(chan error, 5)

                // Try to update the same credential concurrently
                for i := 0; i < 5; i++ {
                    wg.Add(1)
                    go func(idx int) {
                        defer wg.Done()
                        credential := types.VerifiableCredential{
                            Id:           "cred7",
                            Issuer:       issuerDID,
                            IssuanceDate: time.Now().Format(time.RFC3339),
                            CredentialSubject: map[string]interface{}{
                                "update": idx,
                            },
                        }
                        if err := keeper.UpdateCredential(ctx, credential); err != nil {
                            errs <- err
                        }
                    }(i)
                }

                wg.Wait()
                close(errs)

                if len(errs) > 0 {
                    return <-errs
                }
                return nil
            },
            expectErr: true,
            errorContains: "concurrent modification",
        },
        {
            name: "credential chain validation",
            setup: func() error {
                // Create a chain of credentials where one in the middle is invalid
                cred1 := types.VerifiableCredential{
                    Id:           "chain1",
                    Issuer:       issuerDID,
                    IssuanceDate: time.Now().Format(time.RFC3339),
                }
                if err := keeper.IssueCredential(ctx, cred1); err != nil {
                    return err
                }

                cred2 := types.VerifiableCredential{
                    Id:           "chain2",
                    Issuer:       holderDID,
                    IssuanceDate: time.Now().Format(time.RFC3339),
                    Evidence: []types.Evidence{
                        {
                            Id:   "chain1",
                            Type: "VerifiableCredential",
                        },
                    },
                }
                return keeper.IssueCredential(ctx, cred2)
            },
            test: func() error {
                return keeper.VerifyCredentialChain(ctx, "chain2")
            },
            expectErr: true,
            errorContains: "invalid credential chain",
        },
    }

    // Test execution code...
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            if tc.setup != nil {
                err := tc.setup()
                require.NoError(t, err)
            }

            err := tc.test()
            if tc.expectErr {
                require.Error(t, err)
                require.Contains(t, err.Error(), tc.errorContains)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

{{ ... }}
