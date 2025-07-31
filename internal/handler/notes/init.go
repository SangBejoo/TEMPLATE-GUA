package notes

import (
	"context"

	"github.com/SangBejoo/Template/internal/entity"
)

type NotesHandler interface {
	CreateNote(ctx context.Context, note *entity.Note) error
	GetNotes(ctx context.Context) ([]*entity.Note, error)
	UpdateNote(ctx context.Context, note *entity.Note) error
	DeleteNote(ctx context.Context, id int) error
}
