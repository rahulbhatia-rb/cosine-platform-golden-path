package gate

import "testing"

func good() Contract {
	return Contract{
		Service: "agent-api", Environment: "production", Tenant: "cosine",
		Replicas: 3, CPURequest: "500m", CPULimit: "2", MemoryRequest: "1Gi", MemoryLimit: "4Gi",
		MinReplicas: 3, MaxReplicas: 20, ReadinessProbe: true, LivenessProbe: true, PDB: true,
		SecretSource: "aws-secrets-manager", EgressPolicy: "allowlisted", PublicAdminIngress: false,
		IsolatedEnvironment: true,
		SLO: SLO{AvailabilityPercent: 99.9, P95LatencyMS: 500},
		Observability: Observability{Metrics: true, Logs: true, Traces: true},
		Rollout: Rollout{Strategy: "canary", RollbackEnabled: true},
	}
}

func TestGoodProductionContract(t *testing.T) {
	r := Evaluate(good())
	if !r.Allowed {
		t.Fatalf("expected allowed: %+v", r.Findings)
	}
}

func TestRejectsUnsafeSecrets(t *testing.T) {
	c := good()
	c.SecretSource = "plaintext-env"
	if Evaluate(c).Allowed {
		t.Fatal("expected secret policy rejection")
	}
}

func TestRejectsNoRollback(t *testing.T) {
	c := good()
	c.Rollout.RollbackEnabled = false
	if Evaluate(c).Allowed {
		t.Fatal("expected rollback rejection")
	}
}
