package dependency

import (
	"github.com/SangBejoo/Template/init/infra"
	"google.golang.org/grpc"

	base "github.com/SangBejoo/Template/gen/proto"
	baseGrpcServer "github.com/SangBejoo/Template/internal/handler/base"
	notesGrpcHandler "github.com/SangBejoo/Template/internal/handler/notes"
	NotesRepository "github.com/SangBejoo/Template/internal/repository/notes"
	NotesUseCase "github.com/SangBejoo/Template/internal/usecase/notes"
)

func InitGrpcDependency(server *grpc.Server, repo infra.Repository) {
	baseServer := baseGrpcServer.NewBaseHandler()
	base.RegisterBaseServer(server, baseServer)
	notesRepository := NotesRepository.NewNotesRepository(repo.DB)
	notesUseCase := NotesUseCase.NewNotesUseCase(notesRepository)
	notesServer := notesGrpcHandler.NewNotesHandler(notesUseCase)
	base.RegisterNotesServiceServer(server, notesServer)
}
