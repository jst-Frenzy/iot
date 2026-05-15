package data_generator

type Deps struct {
}

type DataGenerator struct {
}

func NewDataGenerator(d *Deps) *DataGenerator { //nolint: revive
	return &DataGenerator{}
}

func (g *DataGenerator) Generate() {}
