package app

import dg "github.com/jst-Frenzy/iot/backend/internal/usecase/data_generator"

type Usecases struct {
	DataGenerator *dg.DataGenerator
}

func NewUsecase() *Usecases {
	dataGenerator := dg.NewDataGenerator(&dg.Deps{})

	return &Usecases{
		DataGenerator: dataGenerator,
	}
}
