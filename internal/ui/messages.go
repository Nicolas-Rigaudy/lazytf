package ui

import (
	"github.com/Nicolas-Rigaudy/lazytf/internal/aws"
	"github.com/Nicolas-Rigaudy/lazytf/internal/terraform"
)

type ProjectSelectedMsg struct {
	ProjectName string
	Index       int
}

type VarFileSelectedMsg struct {
	VarFileName string
	Index       int
}

type InitEnvironmentSelectedMsg struct {
	EnvName  string
	Backends []terraform.BackendVarFile
}

type InitBackendSelectedMsg struct {
	EnvName string
	Backend terraform.BackendVarFile
}

type RunInitMsg struct {
	ProjectPath string
	Options     terraform.InitOptions
}

type RunAWSSSOLoginMsg struct {
	Session *aws.SSOSession
}
