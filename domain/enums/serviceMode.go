package enums

type ServiceMode string

const (
	ServiceModeDevelopment ServiceMode = "development"
	ServiceModeStage       ServiceMode = "stage"
	ServiceModeProduction  ServiceMode = "production"
)

func (r ServiceMode) String() string {
	return string(r)
}

func (r ServiceMode) IsValid() bool {
	switch r {
	case ServiceModeDevelopment,
		ServiceModeStage,
		ServiceModeProduction:
		return true
	default:
		return false
	}
}

func (r ServiceMode) List() []string {
	return []string{
		ServiceModeDevelopment.String(),
		ServiceModeStage.String(),
		ServiceModeProduction.String(),
	}
}
