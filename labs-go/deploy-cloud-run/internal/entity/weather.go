package entity

type ViaCepOutput struct {
	Localidade string `json:"localidade"`
	Erro       string `json:"erro"`
}

func (v *ViaCepOutput) IsErro() bool {
	return v.Erro == "true"
}

type CurrentWeather struct {
	TempC float64 `json:"temp_c"`
}

type WeatherApiOutput struct {
	Current CurrentWeather `json:"current"`
}

type WeatherOutput struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}
