// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package project

import (
	"context"
	"fmt"

	"github.com/azure/azure-dev/cli/azd/pkg/async"
	"github.com/azure/azure-dev/cli/azd/pkg/tools"
)

type ServiceLanguageKind string

const (
	ServiceLanguageDotNet     ServiceLanguageKind = "dotnet"
	ServiceLanguageCsharp     ServiceLanguageKind = "csharp"
	ServiceLanguageFsharp     ServiceLanguageKind = "fsharp"
	ServiceLanguageJavaScript ServiceLanguageKind = "js"
	ServiceLanguageTypeScript ServiceLanguageKind = "ts"
	ServiceLanguagePython     ServiceLanguageKind = "python"
	ServiceLanguageJava       ServiceLanguageKind = "java"
	ServiceLanguageDocker     ServiceLanguageKind = "docker"
)

func parseServiceLanguage(kind ServiceLanguageKind) (ServiceLanguageKind, error) {
	// aliases
	if string(kind) == "py" {
		return ServiceLanguagePython, nil
	}

	switch kind {
	case ServiceLanguageDotNet,
		ServiceLanguageCsharp,
		ServiceLanguageFsharp,
		ServiceLanguageJavaScript,
		ServiceLanguageTypeScript,
		ServiceLanguagePython,
		ServiceLanguageJava:
		// Excluding ServiceLanguageDocker since it is implicitly derived currently, and not an actual language
		return kind, nil
	}

	return ServiceLanguageKind(""), fmt.Errorf("unsupported language '%s'", kind)
}

// FrameworkService is an abstraction for a programming language or framework
// that describe the required tools as well as implementations for
// restore and build commands
type FrameworkService interface {
	// Gets a list of the required external tools for the framework service
	RequiredExternalTools(ctx context.Context) []tools.ExternalTool

	// Initializes the framework service for the specified service configuration
	// This is useful if the framework needs to subscribe to any service events
	Initialize(ctx context.Context, serviceConfig *ServiceConfig) error

	// Restores dependencies for the framework service
	Restore(
		ctx context.Context,
		serviceConfig *ServiceConfig,
	) *async.TaskWithProgress[*ServiceRestoreResult, ServiceProgress]

	// Builds the source for the framework service
	Build(
		ctx context.Context,
		serviceConfig *ServiceConfig,
		restoreOutput *ServiceRestoreResult,
	) *async.TaskWithProgress[*ServiceBuildResult, ServiceProgress]

	// Packages the source suitable for publishing
	// This may optionally perform a rebuild internally depending on the language/framework requirements
	Package(
		ctx context.Context,
		serviceConfig *ServiceConfig,
		buildOutput *ServiceBuildResult,
	) *async.TaskWithProgress[*ServicePackageResult, ServiceProgress]
}

// CompositeFrameworkService is a framework service that requires a nested
// framework service. An example would be a Docker project that uses another
// framework services such as NPM or Python as a dependency. This supports
// local inner-loop as well as release restore & package support.
type CompositeFrameworkService interface {
	FrameworkService
	SetSource(inner FrameworkService)
}
