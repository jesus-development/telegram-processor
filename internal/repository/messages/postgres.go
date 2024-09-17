package messages

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log/slog"
	"strconv"
	"telegram-processor/internal/db"
	"telegram-processor/internal/models"
	"telegram-processor/pkg/models/json"
	"time"
	"unicode/utf8"
)

const (
	MESSAGES_TABLENAME   = "messages"
	EMBEDDINGS_TABLENAME = "embeddings_3large"

	SQL_GET_MESSAGES_WITHOUT_VECTORS = `select m.id, m.text from messages m 
    									left join embeddings_3large emb on m.id = emb.message_id 
										where emb.embedding is null;`
	SQL_GET_CLOSEST = `with emb as (select message_id, embedding <=> $1 as similarity 
										from embeddings_3large order by similarity limit $2
										) 
										select emb.message_id, m.text, emb.similarity 
										from emb left join messages m on emb.message_id = m.id;`
)

type PGMessagesRepository struct {
	db *sql.DB
}

func NewPGMessagesRepository(db *sql.DB) *PGMessagesRepository {
	return &PGMessagesRepository{db: db}
}

func (r *PGMessagesRepository) GetMessagesWithoutVectors(ctx context.Context) ([]*models.Message, error) {
	var messages []*models.Message
	rows, err := r.db.QueryContext(ctx, SQL_GET_MESSAGES_WITHOUT_VECTORS)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return messages, nil
		}
		return nil, fmt.Errorf("r.db.QueryContext -> %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var msg = &models.Message{}
		if err = rows.Scan(&msg.ID, &msg.Text); err != nil {
			return nil, fmt.Errorf("rows.Scan -> %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// Receives json model for better performance
func (r *PGMessagesRepository) ImportChatFromJson(ctx context.Context, chat *json.ChatJson) error {
	if len(chat.Messages) == 0 {
		slog.Info("Empty chat")
		return nil
	}

	txn, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("r.db.Begin -> %w", err)
	}

	stmt, err := txn.PrepareContext(ctx, pq.CopyIn(MESSAGES_TABLENAME, "tg_message_id", "text", "length", "date", "chat_id", "sender_id", "sender_username"))
	if err != nil {
		return fmt.Errorf("txn.PrepareContext -> %w", err)
	}

	for i, msg := range chat.Messages {
		unixdate, err := strconv.ParseInt(msg.DateUnix, 10, 64)
		if err != nil {
			rollback(txn)
			return fmt.Errorf("strconv.ParseInt [i:%d, string:%s] -> %w", i, msg.DateUnix, err)
		}

		_, err = stmt.ExecContext(ctx, msg.TelegramID, msg.String(), utf8.RuneCountInString(msg.String()), time.Unix(unixdate, 0), chat.ID, msg.UserID, msg.Username)
		if err != nil {
			rollback(txn)
			return fmt.Errorf("stmt.ExecContext [i:%d, TelegramID:%d] -> %w", i, msg.TelegramID, err)
		}
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		rollback(txn)
		return fmt.Errorf("stmt.ExecContext [empty] -> %w", err)
	}

	err = stmt.Close()
	if err != nil {
		rollback(txn)
		return fmt.Errorf("stmt.Close -> %w", err)
	}

	err = txn.Commit()
	if err != nil {
		rollback(txn)
		return fmt.Errorf("txn.Commit -> %w", err)
	}

	return nil
}

func rollback(txn *sql.Tx) {
	if err := txn.Rollback(); err != nil {
		slog.Error("Failed to rollback transaction", "error", err)
	}
}

func (r *PGMessagesRepository) InsertEmbeddings(ctx context.Context, embeddings []*models.MessageEmbedding) error {
	if len(embeddings) == 0 {
		return nil
	}

	txn, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("r.db.Begin -> %w", err)
	}

	stmt, err := txn.PrepareContext(ctx, pq.CopyIn(EMBEDDINGS_TABLENAME, "message_id", "embedding"))
	if err != nil {
		return fmt.Errorf("txn.PrepareContext -> %w", err)
	}

	for _, emb := range embeddings {
		_, err = stmt.ExecContext(ctx, emb.MessageID, db.Vector32Float(emb.Embedding))
		if err != nil {
			rollback(txn)
			return fmt.Errorf("stmt.ExecContext [message_id:%d] -> %w", emb.MessageID, err)
		}
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		rollback(txn)
		return fmt.Errorf("stmt.ExecContext [empty] -> %w", err)
	}

	err = stmt.Close()
	if err != nil {
		rollback(txn)
		return fmt.Errorf("stmt.Close -> %w", err)
	}

	err = txn.Commit()
	if err != nil {
		rollback(txn)
		return fmt.Errorf("txn.Commit -> %w", err)
	}

	return nil
}

func (r *PGMessagesRepository) GetClosest(ctx context.Context, search []float32, limit int64) ([]*models.Message, error) {
	var messages []*models.Message

	rows, err := r.db.QueryContext(ctx, SQL_GET_CLOSEST, db.Vector32Float(search), limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return messages, nil
		}
		return nil, fmt.Errorf("r.db.QueryContext -> %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var msg = &models.Message{}
		if err = rows.Scan(&msg.ID, &msg.Text, &msg.Similarity); err != nil {
			return nil, fmt.Errorf("rows.Scan -> %w", err)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
