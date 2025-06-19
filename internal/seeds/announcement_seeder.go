package seeds

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mariopaath23/backend-jte-ticketing/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SeedAnnouncements populates the database with initial announcements.
func SeedAnnouncements(db *mongo.Database) {
	collection := db.Collection("announcements")

	announcements := []models.Announcement{
		{
			Title:            "Perbaikan pada ruangan JTE-2",
			Author:           "Admin JTE",
			DatePublished:    parseDateTime("2025-04-10T14:17:30Z"),
			Content:          "Sehubungan dengan terjadinya kebocoran pada plafon ruangan JTE-2, maka pengajuan peminjaman tidak dapat dilakukan untuk sementara khusus ruangan JTE-2. Terima kasih",
			Tags:             []string{"urgent", "maintenance"},
			AnnouncementType: "public",
		},
		{
			Title:            "Maintenance Sistem",
			Author:           "Tim Developer - Mario",
			DatePublished:    parseDateTime("2025-03-05T14:17:30Z"),
			Content:          "Sistem akan mengalami maintenance pada tanggal 10 Maret 2025 pukul 14:00 - 16:00 WIB. Mohon untuk tidak melakukan peminjaman pada waktu tersebut.",
			Tags:             []string{"critical", "maintenance"},
			AnnouncementType: "public",
		},
		{
			Title:            "Libur Hari Raya Idul Fitri",
			Author:           "Admin JTE",
			DatePublished:    parseDateTime("2025-04-01T08:00:00Z"),
			Content:          "Diberitahukan kepada seluruh civitas akademika, bahwa libur Hari Raya Idul Fitri akan dimulai pada tanggal 8 April hingga 15 April 2025. Aktivitas perkuliahan akan dimulai kembali pada tanggal 16 April 2025.",
			Tags:             []string{"holiday", "info"},
			AnnouncementType: "public",
		},
		{
			Title:            "Rapat Internal Staff",
			Author:           "Kepala Jurusan",
			DatePublished:    parseDateTime("2025-06-10T09:00:00Z"),
			Content:          "Diharapkan kehadiran seluruh staff administrasi untuk rapat internal pada hari Jumat, 13 Juni 2025 pukul 15:00 WITA.",
			Tags:             []string{"internal", "meeting"},
			AnnouncementType: "private",
		},
	}

	for _, announcement := range announcements {
		// Use title as a unique key for seeding to prevent duplicates
		var existing models.Announcement
		err := collection.FindOne(context.TODO(), bson.M{"title": announcement.Title}).Decode(&existing)
		if err == mongo.ErrNoDocuments {
			announcement.ID = primitive.NewObjectID()
			_, insertErr := collection.InsertOne(context.TODO(), announcement)
			if insertErr != nil {
				log.Printf("Failed to seed announcement '%s': %v", announcement.Title, insertErr)
			} else {
				fmt.Printf("Successfully seeded announcement: '%s'\n", announcement.Title)
			}
		}
	}
}

// Helper to parse date-time strings
func parseDateTime(dateStr string) time.Time {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		log.Printf("Could not parse date-time: %s", dateStr)
		return time.Time{}
	}
	return t
}
