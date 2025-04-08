package seeder

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v4/pgxpool"
)

func SeedTools(db *pgxpool.Pool, total int) {
	ctx := context.Background()

	fmt.Printf("Seeding %d tools...\n", total)

	tx, err := db.Begin(ctx)
	if err != nil {
		log.Fatal("Failed to begin transaction:", err)
	}
	defer tx.Rollback(ctx)

	for i := 0; i < total; i++ {
		name := gofakeit.AppName() + " " + gofakeit.HackerNoun()
		desc := gofakeit.Paragraph(1, 3, 20, " ")

		_, err := tx.Exec(ctx, `
			INSERT INTO tools (name, description, created_at)
			VALUES ($1, $2, $3)
		`, name, desc, time.Now())

		if err != nil {
			log.Fatalf("Insert failed at %d: %v", i, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal("Failed to commit:", err)
	}

	fmt.Println("✅ Tools seeded successfully!")
}
