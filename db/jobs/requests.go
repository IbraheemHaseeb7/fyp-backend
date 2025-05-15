package jobs

import (
	"fmt"
	"time"

	"github.com/IbraheemHaseeb7/types"
	"gorm.io/gorm"
)

/*
	Job: CheckRequestDeadlines
	Description: This job checks the database for requests that are pending check every minute.
		If a request deadline is passed, it changes the status to "expired".
	Usage: This job is used to ensure that requests that are pending check are not left in the database indefinitely.
		It is important to keep the database clean and to ensure that requests that are no longer valid are removed.
	Note: This job should be run in a separate goroutine to avoid blocking the main thread.
		It is important to ensure that the database connection is properly closed after the job is done.
*/
func CheckRequestDeadlines(db *gorm.DB) {
	for {
		time.Sleep(1 * time.Minute)

		var requests []types.Request
		now := time.Now().UTC()

		if err := db.Where("(status = ? OR status = ?) AND until < ?", "searching", "proposal", now).Find(&requests).Error; err != nil {
			fmt.Println("Error fetching requests:", err)
			continue
		}

		for _, request := range requests {
			request.Status = "expired"
			if err := db.Save(&request).Error; err != nil {
				fmt.Println("Error updating request status:", err)
			}
		}
	}
}
