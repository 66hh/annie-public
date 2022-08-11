package constant

type ClimateType struct {
	CLIMATE_NONE         uint16
	CLIMATE_SUNNY        uint16
	CLIMATE_CLOUDY       uint16
	CLIMATE_RAIN         uint16
	CLIMATE_THUNDERSTORM uint16
	CLIMATE_SNOW         uint16
	CLIMATE_MIST         uint16
}

func GetClimateTypeConst() (r *ClimateType) {
	r = new(ClimateType)
	r.CLIMATE_NONE = 0
	r.CLIMATE_SUNNY = 1
	r.CLIMATE_CLOUDY = 2
	r.CLIMATE_RAIN = 3
	r.CLIMATE_THUNDERSTORM = 4
	r.CLIMATE_SNOW = 5
	r.CLIMATE_MIST = 6
	return r
}
