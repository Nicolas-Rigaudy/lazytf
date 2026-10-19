package ui

import (
	"testing"

	"github.com/Nicolas-Rigaudy/lazytf/internal/terraform"
)

func rc(address, module string, actions ...string) terraform.ResourceChange {
	return terraform.ResourceChange{Address: address, ModuleAddress: module, Change: terraform.Change{Actions: actions}}
}

func TestBuildModuleGroups(t *testing.T) {
	plan := &terraform.PlanSummary{ResourceChanges: []terraform.ResourceChange{
		rc("aws_s3_bucket.logs", "", "create"),
		rc(`aws_s3_bucket.k["a.b"]`, "", "update"),
		rc("data.aws_iam_policy.p", "", "read"),
		rc("module.vpc.aws_subnet.a", "module.vpc", "delete", "create"),
		rc("module.vpc.aws_subnet.b", "module.vpc", "create", "delete"),
		rc("module.vpc.module.nat.aws_eip.e", "module.vpc.module.nat", "delete"),
		rc("module.vpc.data.aws_az.all", "module.vpc", "no-op"),
		rc("aws_iam_role.r", "", "no-op"),
		rc("aws_iam_role.empty", ""),
	}}

	type counts struct {
		name                 string
		resources            int
		add, change, destroy int
	}
	want := []counts{
		{name: "(root)", resources: 3, add: 1, change: 1},
		{name: "module.vpc", resources: 2, add: 2, destroy: 2},
		{name: "module.vpc.module.nat", resources: 1, destroy: 1},
	}

	// Run several times: the order used to come from map iteration.
	for run := 0; run < 20; run++ {
		got := buildModuleGroups(plan)
		if len(got) != len(want) {
			t.Fatalf("got %d groups, want %d: %+v", len(got), len(want), got)
		}
		for i, w := range want {
			g := got[i]
			gc := counts{name: g.Name, resources: len(g.Resources), add: g.AddCount, change: g.ChangeCount, destroy: g.DeleteCount}
			if gc != w {
				t.Errorf("run %d group %d = %+v, want %+v", run, i, gc, w)
			}
		}
	}
}

func TestPlanSummaryTitleMatchesTerraformTotals(t *testing.T) {
	// From `terraform plan` on a create_before_destroy replace:
	// "Plan: 1 to add, 0 to change, 1 to destroy."
	plan := &terraform.PlanSummary{ResourceChanges: []terraform.ResourceChange{
		rc("terraform_data.cbd", "", "create", "delete"),
		rc(`terraform_data.k["a.b"]`, "", "no-op"),
	}}
	got := NewPlanSummaryViewModel(plan).Title()
	want := "Plan Summary - 1 to add, 0 to change, 1 to destroy"
	if got != want {
		t.Errorf("Title() = %q, want %q", got, want)
	}
}

func TestActionPrefix(t *testing.T) {
	tests := []struct {
		actions []string
		want    string
	}{
		{[]string{"create"}, "+"},
		{[]string{"delete"}, "-"},
		{[]string{"update"}, "~"},
		{[]string{"read"}, "<="},
		{[]string{"delete", "create"}, "-/+"},
		{[]string{"create", "delete"}, "+/-"},
		{nil, "~"},
	}
	for _, tt := range tests {
		if got, _ := actionPrefix(tt.actions); got != tt.want {
			t.Errorf("actionPrefix(%v) = %q, want %q", tt.actions, got, tt.want)
		}
	}
}

func TestParsePlanJSONReadsModuleAddress(t *testing.T) {
	data := []byte(`{"resource_changes":[
		{"address":"terraform_data.cbd","change":{"actions":["create","delete"]}},
		{"address":"module.vpc.terraform_data.m","module_address":"module.vpc","change":{"actions":["create"]}}
	]}`)
	plan, err := terraform.ParsePlanJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if got := moduleName(plan.ResourceChanges[0]); got != "(root)" {
		t.Errorf("root resource module = %q, want (root)", got)
	}
	if got := moduleName(plan.ResourceChanges[1]); got != "module.vpc" {
		t.Errorf("module resource module = %q, want module.vpc", got)
	}
}
