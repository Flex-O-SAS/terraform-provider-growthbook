package growthbookapi

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
