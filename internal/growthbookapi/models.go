package growthbookapi

import (
	"encoding/json"
)

// Project represents a GrowthBook project object.
type Project struct {
	ID          string          `json:"id,omitempty"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Settings    ProjectSettings `json:"settings,omitempty"`
	DateCreated string          `json:"dateCreated,omitempty"`
	DateUpdated string          `json:"dateUpdated,omitempty"`
}

// ProjectSettings holds the settings for a GrowthBook project.
type ProjectSettings struct {
	StatsEngine string `json:"statsEngine,omitempty"`
}

// Feature represents a GrowthBook feature object.
type Feature struct {
	ID            string                              `json:"id,omitempty"`
	Archived      bool                                `json:"archived"`
	Description   string                              `json:"description,omitempty"`
	Owner         string                              `json:"owner,omitempty"`
	Project       string                              `json:"project,omitempty"`
	ValueType     string                              `json:"valueType,omitempty"`
	DefaultValue  string                              `json:"defaultValue,omitempty"`
	Tags          []string                            `json:"tags"`
	Environments  map[string]FeatureEnvironmentConfig `json:"environments,omitempty"`
	Prerequisites []string                            `json:"prerequisites"`
}

// FeatureEnvironmentConfig holds the configuration for a GrowthBook environment.
type FeatureEnvironmentConfig struct {
	Enabled      bool          `json:"enabled"`
	DefaultValue string        `json:"defaultValue,omitempty"`
	Definition   string        `json:"definition,omitempty"`
	Rules        []FeatureRule `json:"rules"`
}

// FeatureRule represents a targeting rule for a feature in GrowthBook.
type FeatureRule struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description,omitempty"`
	Condition   string `json:"condition,omitempty"`
	// force / rollout
	Value string `json:"value,omitempty"`
	// rollout only
	Coverage      *float64 `json:"coverage,omitempty"`
	HashAttribute string   `json:"hashAttribute,omitempty"`
	// experiment-ref only
	ExperimentID        string                       `json:"experimentId,omitempty"`
	Variations          []FeatureVariation           `json:"variations,omitempty"`
	SavedGroupTargeting []FeatureSavedGroupTargeting `json:"savedGroupTargeting,omitempty"`
	Prerequisites       []FeaturePrerequisite        `json:"prerequisites,omitempty"`
}

// FeatureVariation represents a single variation in an experiment-ref rule.
type FeatureVariation struct {
	Value       string `json:"value"`
	VariationID string `json:"variationId"`
}

// FeatureSavedGroupTargeting represents targeting configuration based on saved groups.
type FeatureSavedGroupTargeting struct {
	MatchType   string   `json:"matchType,omitempty"`
	SavedGroups []string `json:"savedGroups"`
}

// FeaturePrerequisite represents a prerequisite for a rule or variation.
type FeaturePrerequisite struct {
	ID        string `json:"id"`
	Condition string `json:"condition"`
}

// FeatureDraft represents a draft configuration for a feature.
type FeatureDraft struct {
	Enabled    bool          `json:"enabled"`
	Rules      []FeatureRule `json:"rules"`
	Definition string        `json:"definition,omitempty"`
}

// Environment represents a GrowthBook environment object.
type Environment struct {
	ID           string   `json:"id,omitempty"`
	Description  string   `json:"description,omitempty"`
	ToggleOnList bool     `json:"toggleOnList"`
	DefaultState bool     `json:"defaultState"`
	Projects     []string `json:"projects,omitempty"`
}

type Attribute struct {
	Property    string   `json:"property"`
	DataType    string   `json:"datatype"`
	Format      string   `json:"format"`
	EnumValues  string   `json:"enum"`
	Projects    []string `json:"projects"`
	Archived    bool     `json:"archived"`
	Description string   `json:"description"`
}

// SDKConnection represents a GrowthBook SDK Connection object.
type SDKConnection struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Language    string `json:"language"`
	Environment string `json:"environment"`

	// optional
	Languages                   []string `json:"languages,omitempty"`
	SdkVersion                  string   `json:"sdkVersion,omitempty"`
	Projects                    []string `json:"projects,omitempty"`
	EncryptPayload              bool     `json:"encryptPayload"`
	IncludeVisualExperiments    bool     `json:"includeVisualExperiments"`
	IncludeDraftExperiments     bool     `json:"includeDraftExperiments"`
	IncludeExperimentNames      bool     `json:"includeExperimentNames"`
	IncludeRedirectExperiments  bool     `json:"includeRedirectExperiments"`
	IncludeRuleIDs              bool     `json:"includeRuleIds"`
	ProxyEnabled                bool     `json:"proxyEnabled"`
	ProxyHost                   string   `json:"proxyHost,omitempty"`
	HashSecureAttributes        bool     `json:"hashSecureAttributes"`
	RemoteEvalEnabled           bool     `json:"remoteEvalEnabled"`
	SavedGroupReferencesEnabled bool     `json:"savedGroupReferencesEnabled"`

	// computed
	Organization    string `json:"organization,,omitempty"`
	EncryptionKey   string `json:"encryptionKey,omitempty"`
	Key             string `json:"key,omitempty"`
	ProxySigningKey string `json:"proxySigningKey,omitempty"`
	SseEnabled      bool   `json:"sseEnabled,omitempty"`
	DateCreated     string `json:"dateCreated,omitempty"`
	DateUpdated     string `json:"dateUpdated,omitempty"`
}

// SavedGroup represents a GrowthBook saved group.
type SavedGroup struct {
	ID           string   `json:"id,omitempty"`
	Name         string   `json:"name"`
	Type         string   `json:"type,omitempty"`
	Condition    string   `json:"condition,omitempty"`
	AttributeKey string   `json:"attributeKey,omitempty"`
	Values       []string `json:"values,omitempty"`
	Owner        string   `json:"owner,omitempty"`
	Projects     []string `json:"projects,omitempty"`
	Description  string   `json:"description,omitempty"`
	DateCreated  string   `json:"dateCreated,omitempty"`
	DateUpdated  string   `json:"dateUpdated,omitempty"`
}

// Segment represents a GrowthBook segment.
type Segment struct {
	ID             string   `json:"id,omitempty"`
	Name           string   `json:"name"`
	Owner          string   `json:"owner,omitempty"`
	Description    string   `json:"description,omitempty"`
	DatasourceID   string   `json:"datasourceId,omitempty"`
	IdentifierType string   `json:"identifierType,omitempty"`
	Type           string   `json:"type,omitempty"`
	Query          string   `json:"query,omitempty"`
	FactTableID    string   `json:"factTableId,omitempty"`
	Projects       []string `json:"projects,omitempty"`
	ManagedBy      string   `json:"managedBy,omitempty"`
	DateCreated    string   `json:"dateCreated,omitempty"`
	DateUpdated    string   `json:"dateUpdated,omitempty"`
}

// Metric represents a GrowthBook (legacy) metric.
type Metric struct {
	ID           string          `json:"id,omitempty"`
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	Owner        string          `json:"owner,omitempty"`
	DatasourceID string          `json:"datasourceId,omitempty"`
	Type         string          `json:"type,omitempty"`
	Tags         []string        `json:"tags,omitempty"`
	Projects     []string        `json:"projects,omitempty"`
	Archived     bool            `json:"archived"`
	Behavior     json.RawMessage `json:"behavior,omitempty"`
	SQL          json.RawMessage `json:"sql,omitempty"`
	ManagedBy    string          `json:"managedBy,omitempty"`
	DateCreated  string          `json:"dateCreated,omitempty"`
	DateUpdated  string          `json:"dateUpdated,omitempty"`
}

// Dimension represents a GrowthBook dimension.
type Dimension struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Owner          string `json:"owner,omitempty"`
	DatasourceID   string `json:"datasourceId,omitempty"`
	IdentifierType string `json:"identifierType,omitempty"`
	Query          string `json:"query,omitempty"`
	ManagedBy      string `json:"managedBy,omitempty"`
	DateCreated    string `json:"dateCreated,omitempty"`
	DateUpdated    string `json:"dateUpdated,omitempty"`
}

// FactTable represents a GrowthBook fact table.
type FactTable struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Owner       string   `json:"owner,omitempty"`
	Projects    []string `json:"projects,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Datasource  string   `json:"datasource,omitempty"`
	UserIDTypes []string `json:"userIdTypes,omitempty"`
	SQL         string   `json:"sql,omitempty"`
	EventName   string   `json:"eventName,omitempty"`
	ManagedBy   string   `json:"managedBy,omitempty"`
	Archived    bool     `json:"archived"`
	DateCreated string   `json:"dateCreated,omitempty"`
	DateUpdated string   `json:"dateUpdated,omitempty"`
}

// FactMetric represents a GrowthBook fact metric.
type FactMetric struct {
	ID                           string          `json:"id,omitempty"`
	Name                         string          `json:"name"`
	Description                  string          `json:"description,omitempty"`
	Owner                        string          `json:"owner,omitempty"`
	Projects                     []string        `json:"projects,omitempty"`
	Tags                         []string        `json:"tags,omitempty"`
	Datasource                   string          `json:"datasource,omitempty"`
	MetricType                   string          `json:"metricType,omitempty"`
	Numerator                    json.RawMessage `json:"numerator,omitempty"`
	Denominator                  json.RawMessage `json:"denominator,omitempty"`
	Inverse                      bool            `json:"inverse"`
	CappingSettings              json.RawMessage `json:"cappingSettings,omitempty"`
	WindowSettings               json.RawMessage `json:"windowSettings,omitempty"`
	PriorSettings                json.RawMessage `json:"priorSettings,omitempty"`
	RegressionAdjustmentSettings json.RawMessage `json:"regressionAdjustmentSettings,omitempty"`
	RiskThresholdSuccess         *float64        `json:"riskThresholdSuccess,omitempty"`
	RiskThresholdDanger          *float64        `json:"riskThresholdDanger,omitempty"`
	MinPercentChange             *float64        `json:"minPercentChange,omitempty"`
	MaxPercentChange             *float64        `json:"maxPercentChange,omitempty"`
	MinSampleSize                *float64        `json:"minSampleSize,omitempty"`
	TargetMDE                    *float64        `json:"targetMDE,omitempty"`
	ManagedBy                    string          `json:"managedBy,omitempty"`
	Archived                     bool            `json:"archived"`
	DateCreated                  string          `json:"dateCreated,omitempty"`
	DateUpdated                  string          `json:"dateUpdated,omitempty"`
}

// Namespace represents a GrowthBook namespace.
type Namespace struct {
	ID            string `json:"id,omitempty"`
	DisplayName   string `json:"displayName,omitempty"`
	Description   string `json:"description,omitempty"`
	Status        string `json:"status,omitempty"`
	Format        string `json:"format,omitempty"`
	HashAttribute string `json:"hashAttribute,omitempty"`
	Seed          string `json:"seed,omitempty"`
}

// Experiment represents a GrowthBook experiment.
type Experiment struct {
	ID                string          `json:"id,omitempty"`
	Name              string          `json:"name"`
	TrackingKey       string          `json:"trackingKey,omitempty"`
	Type              string          `json:"type,omitempty"`
	Project           string          `json:"project,omitempty"`
	Hypothesis        string          `json:"hypothesis,omitempty"`
	Description       string          `json:"description,omitempty"`
	Tags              []string        `json:"tags,omitempty"`
	Owner             string          `json:"owner,omitempty"`
	Archived          bool            `json:"archived"`
	Status            string          `json:"status,omitempty"`
	AutoRefresh       bool            `json:"autoRefresh"`
	HashAttribute     string          `json:"hashAttribute,omitempty"`
	FallbackAttribute string          `json:"fallbackAttribute,omitempty"`
	DatasourceID      string          `json:"datasourceId,omitempty"`
	AssignmentQueryID string          `json:"assignmentQueryId,omitempty"`
	SegmentID         string          `json:"segmentId,omitempty"`
	Metrics           []string        `json:"metrics,omitempty"`
	SecondaryMetrics  []string        `json:"secondaryMetrics,omitempty"`
	GuardrailMetrics  []string        `json:"guardrailMetrics,omitempty"`
	StatsEngine       string          `json:"statsEngine,omitempty"`
	Variations        json.RawMessage `json:"variations,omitempty"`
	Phases            json.RawMessage `json:"phases,omitempty"`
	DateCreated       string          `json:"dateCreated,omitempty"`
	DateUpdated       string          `json:"dateUpdated,omitempty"`
}

// DataSource represents a GrowthBook data source (read-only).
type DataSource struct {
	ID           string   `json:"id,omitempty"`
	Name         string   `json:"name,omitempty"`
	Type         string   `json:"type,omitempty"`
	Description  string   `json:"description,omitempty"`
	ProjectIDs   []string `json:"projectIds,omitempty"`
	EventTracker string   `json:"eventTracker,omitempty"`
	DateCreated  string   `json:"dateCreated,omitempty"`
	DateUpdated  string   `json:"dateUpdated,omitempty"`
}
