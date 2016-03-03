package database

import "log"

// Setup ...sets up the database
func (d *Database) Setup() {
	log.Println("Race database setup starting...")
	d.Conn.CreateTable(&Deployments{})
	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Deployments{})

	d.Conn.CreateTable(&Resources{})
	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Resources{})

	d.Conn.CreateTable(&Tasks{})
	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Tasks{})

	d.Conn.CreateTable(&TaskResults{})
	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&TaskResults{})

	d.Conn.CreateTable(&VerificationResults{})
	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&VerificationResults{})

	d.Conn.CreateTable(&Verifications{})
	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Verifications{})

	d.Conn.CreateTable(&Workers{})
	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Workers{})

	d.Conn.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&Deployments{}, &Resources{}, &Tasks{}, &TaskResults{}, &VerificationResults{}, &Verifications{}, &Workers{})
	log.Println("Complete!")
}
