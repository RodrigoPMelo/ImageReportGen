package domain

type Orientation string

const (
	OrientationLandscape Orientation = "Paisagem"
	OrientationPortrait  Orientation = "Retrato"
)

type ImageMetadata struct {
	Path        string
	Orientation Orientation
}

type ReportRequest struct {
	TemplatePath string
	ImagePaths   []string
	OutputPath   string
}

type ReportResult struct {
	OutputPath     string
	TotalImages    int
	LandscapeCount int
	PortraitCount  int
}

type ProcessedInput struct {
	ImagePaths   []string
	IgnoredPaths []string
}
