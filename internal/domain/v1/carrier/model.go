package carrier

type Carrier struct {
	Id                     string `json:"id"`
	Key                    string `json:"key"`
	Name                   string `json:"name"`
	AbroudInfo             string `json:"abroudInfo"`
	NanoSim                string `json:"nanoSim"`
	NanoSimDescription     string `json:"nanoSimDescription"`
	MnpSim                 string `json:"mnpSim"`
	MnpSimDescription      string `json:"mnpSimDescription"`
	MnpNanoSim             string `json:"mnpNanoSim"`
	MnpNanoSimDescription  string `json:"mnpNanoSimDescription"`
	StandartSim            string `json:"standardSim"`
	StandartSimDescription string `json:"standardSimDescription"`
	SpKey                  string `json:"spKey"`
	LogoUrl                string `json:"logoUrl"`
}
