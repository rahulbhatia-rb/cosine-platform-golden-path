package gate

import (
	"fmt"
	"strings"
)

type SLO struct {
	AvailabilityPercent float64 `json:"availability_percent"`
	P95LatencyMS        int     `json:"p95_latency_ms"`
}

type Observability struct {
	Metrics bool `json:"metrics"`
	Logs    bool `json:"logs"`
	Traces  bool `json:"traces"`
}

type Rollout struct {
	Strategy        string `json:"strategy"`
	RollbackEnabled bool   `json:"rollback_enabled"`
}

type Contract struct {
	Service             string        `json:"service"`
	Environment         string        `json:"environment"`
	Tenant              string        `json:"tenant"`
	Replicas            int           `json:"replicas"`
	CPURequest          string        `json:"cpu_request"`
	CPULimit            string        `json:"cpu_limit"`
	MemoryRequest       string        `json:"memory_request"`
	MemoryLimit         string        `json:"memory_limit"`
	MinReplicas         int           `json:"min_replicas"`
	MaxReplicas         int           `json:"max_replicas"`
	ReadinessProbe      bool          `json:"readiness_probe"`
	LivenessProbe       bool          `json:"liveness_probe"`
	PDB                  bool          `json:"pod_disruption_budget"`
	SecretSource        string        `json:"secret_source"`
	EgressPolicy        string        `json:"egress_policy"`
	PublicAdminIngress  bool          `json:"public_admin_ingress"`
	BackupRestore       bool          `json:"backup_restore"`
	Stateful            bool          `json:"stateful"`
	IsolatedEnvironment bool          `json:"isolated_environment"`
	SLO                 SLO           `json:"slo"`
	Observability       Observability `json:"observability"`
	Rollout             Rollout       `json:"rollout"`
}

type Finding struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type Result struct {
	Allowed  bool      `json:"allowed"`
	Findings []Finding `json:"findings"`
}

func Evaluate(c Contract) Result {
	var f []Finding
	prod := c.Environment == "production"

	if prod {
		if c.Replicas < 2 && !c.Stateful {
			f = append(f, Finding{"high", "production stateless workloads require at least 2 replicas"})
		}
		if c.CPURequest == "" || c.CPULimit == "" || c.MemoryRequest == "" || c.MemoryLimit == "" {
			f = append(f, Finding{"high", "production workload requires CPU and memory requests/limits"})
		}
		if !c.ReadinessProbe || !c.LivenessProbe {
			f = append(f, Finding{"high", "production workload requires readiness and liveness probes"})
		}
		if !c.PDB {
			f = append(f, Finding{"medium", "PodDisruptionBudget required for production"})
		}
		if c.MinReplicas < 1 || c.MaxReplicas < c.MinReplicas {
			f = append(f, Finding{"high", "autoscaling bounds are invalid"})
		}
		if c.SLO.AvailabilityPercent < 99.0 || c.SLO.P95LatencyMS <= 0 {
			f = append(f, Finding{"high", "production workload requires explicit availability and latency SLOs"})
		}
		if !c.Observability.Metrics || !c.Observability.Logs || !c.Observability.Traces {
			f = append(f, Finding{"high", "metrics, logs, and traces are required"})
		}
		if strings.Contains(strings.ToLower(c.SecretSource), "plain") || c.SecretSource == "" {
			f = append(f, Finding{"critical", "secrets must come from a managed secret provider"})
		}
		if c.EgressPolicy == "" || c.EgressPolicy == "*" {
			f = append(f, Finding{"high", "explicit bounded egress policy required"})
		}
		if c.PublicAdminIngress {
			f = append(f, Finding{"critical", "administrative ingress cannot be public"})
		}
		if !c.IsolatedEnvironment {
			f = append(f, Finding{"high", "production environment must be isolated"})
		}
		if c.Stateful && !c.BackupRestore {
			f = append(f, Finding{"critical", "stateful production workloads require backup/restore policy"})
		}
	}

	switch c.Rollout.Strategy {
	case "canary", "blue-green", "staged":
	default:
		f = append(f, Finding{"high", fmt.Sprintf("unsupported rollout strategy %q", c.Rollout.Strategy)})
	}
	if !c.Rollout.RollbackEnabled {
		f = append(f, Finding{"high", "rollback must be enabled"})
	}

	allowed := true
	for _, x := range f {
		if x.Severity == "high" || x.Severity == "critical" {
			allowed = false
		}
	}
	return Result{Allowed: allowed, Findings: f}
}
