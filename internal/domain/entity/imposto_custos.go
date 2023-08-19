package entity

type Imposto struct {
	Id    int
	Uf    string
	Icms  float64
	Ipi   float64
	Difal float64
}

type Custo struct {
	Id       int
	Hmaquina float64
	Zincagem  float64
	Retifica float64
	Tempera  float64
	Dobra    float64
}
