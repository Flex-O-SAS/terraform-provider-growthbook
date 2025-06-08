package growthbookapi

type Project struct {
	ID          string          `json:"id,omitempty"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Settings    ProjectSettings `json:"settings,omitempty"`
	DateCreated string          `json:"dateCreated,omitempty"`
	DateUpdated string          `json:"dateUpdated,omitempty"`
}

type ProjectSettings struct {
	StatsEngine string `json:"statsEngine,omitempty"`
}

type Feature struct {
	ID            string                       `json:"id,omitempty"`
	Archived      bool                         `json:"archived"`
	Description   string                       `json:"description,omitempty"`
	Owner         string                       `json:"owner,omitempty"`
	Project       string                       `json:"project,omitempty"`
	ValueType     string                       `json:"valueType,omitempty"`
	DefaultValue  string                       `json:"defaultValue,omitempty"`
	Tags          []string                     `json:"tags,omitempty"`
	Environments  map[string]EnvironmentConfig `json:"environments,omitempty"`
	Prerequisites []string                     `json:"prerequisites,omitempty"`
}

type EnvironmentConfig struct {
	Enabled    bool   `json:"enabled"`
	Rules      []Rule `json:"rules"`
	Definition string `json:"definition,omitempty"`
	Draft      *Draft `json:"draft,omitempty"`
}

type Rule struct {
	Type                   string                `json:"type"`
	Value                  string                `json:"value,omitempty"`
	Coverage               float64               `json:"coverage,omitempty"`
	HashAttribute          string                `json:"hashAttribute,omitempty"`
	Condition              string                `json:"condition,omitempty"`
	Description            string                `json:"description,omitempty"`
	SavedGroupTargeting    []SavedGroupTargeting `json:"savedGroupTargeting,omitempty"`
	ID                     string                `json:"id,omitempty"`
	Enabled                bool                  `json:"enabled"`
	Prerequisites          []Prerequisite        `json:"prerequisites,omitempty"`
	Variations             []Variation           `json:"variations,omitempty"`
	ExperimentID           string                `json:"experimentId,omitempty"`
	TrackingKey            string                `json:"trackingKey,omitempty"`
	FallbackAttribute      string                `json:"fallbackAttribute,omitempty"`
	DisableStickyBucketing bool                  `json:"disableStickyBucketing"`
	BucketVersion          float64               `json:"bucketVersion,omitempty"`
	MinBucketVersion       float64               `json:"minBucketVersion,omitempty"`
	Namespace              *Namespace            `json:"namespace,omitempty"`
	Values                 []Value               `json:"values,omitempty"`
}

type SavedGroupTargeting struct {
	MatchType   string   `json:"matchType"`
	SavedGroups []string `json:"savedGroups"`
}

type Prerequisite struct {
	ID        string `json:"id"`
	Condition string `json:"condition"`
}

type Variation struct {
	Value       string `json:"value"`
	VariationID string `json:"variationId"`
}

type Namespace struct {
	Enabled bool      `json:"enabled"`
	Name    string    `json:"name"`
	Range   []float64 `json:"range"`
}

type Value struct {
	Value  string  `json:"value"`
	Weight float64 `json:"weight"`
	Name   string  `json:"name,omitempty"`
}

type Draft struct {
	Enabled    bool   `json:"enabled"`
	Rules      []Rule `json:"rules"`
	Definition string `json:"definition,omitempty"`
}

type Environment struct {
	ID           string   `json:"id,omitempty"`
	Description  string   `json:"description,omitempty"`
	ToggleOnList bool     `json:"toggleOnList"`
	DefaultState bool     `json:"defaultState"`
	Projects     []string `json:"projects,omitempty"`
}

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
	IncludeRuleIds              bool     `json:"includeRuleIds"`
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
