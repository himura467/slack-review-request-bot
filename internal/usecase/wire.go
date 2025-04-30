package usecase

import "github.com/google/wire"

var Set = wire.NewSet(
	NewSlackUsecase,
	wire.Bind(new(SlackUsecase), new(*SlackUsecaseImpl)),
)
