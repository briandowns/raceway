package database

import "log"

// Setup ...sets up the database
func (d *Database) Setup() {
	log.Println("Race database setup starting...")
	d.Conn.CreateTable(&Deployment{})
	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Deployment{})

	d.Conn.CreateTable(&Resource{})
	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Resource{})

	d.Conn.CreateTable(&Task{})
	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Task{})

	d.Conn.CreateTable(&TaskResult{})
	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&TaskResult{})

	d.Conn.CreateTable(&VerificationResult{})
	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&VerificationResult{})

	d.Conn.CreateTable(&Verification{})
	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Verification{})

	d.Conn.CreateTable(&Worker{})
	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Worker{})

	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&Deployment{}, &Resource{}, &Task{}, &TaskResult{}, &VerificationResult{}, &Verification{}, &Worker{})
	log.Println("Complete!")
}
