package awpost052

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

// Connection details
var (
	CHostname = ""
	CPort     = 2345
	CUsername = ""
	CPassword = ""
	CDatabase = ""
)

// Userdata is for holding full user data
// Userdata table + Username
type MSDSCourse struct {
	CID     string `json:"courseID`
	CNAME   string `json:"course_name"`
	CPREREQ string `json:"prerequisite"`
}

func CopenConnection() (*sql.DB, error) {
	// connection string
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		CHostname, CPort, CUsername, CPassword, CDatabase)

	// open database
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// The function returns the CID of a course
// -1 if the course does not exist
func Cexists(coursename string) int {
	coursename = strings.ToLower(coursename)

	db, err := CopenConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	CID := 0
	statement := fmt.Sprintf(`SELECT "cid" FROM "msdscoursecatalog" where cname = '%s'`, coursename)
	rows, err := db.Query(statement)

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return -1
		}
		CID = id
	}
	defer rows.Close()
	return CID
}

// AddCourse adds a new Course to the database
// Returns new Course Name
// -1 if there was an error
func AddCourse(d MSDSCourse) int {
	d.CNAME = strings.ToLower(d.CNAME)

	db, err := CopenConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()
	/*
		courID := Cexists(d.CNAME)
		if courID != -1 {
			fmt.Println("Course Name already exists:", d.CNAME)
			return -1
		}

			insertStatement := `insert into "msdscoursecatalog" ("cname") values ($1)`
			_, err = db.Exec(insertStatement, d.CNAME)
			if err != nil {
				fmt.Println(err)
				return -1
			}

			courID = Cexists(d.CNAME)
			if courID == -1 {
				return courID
			}
	*/
	insertStatement := `insert into "msdscoursecatalog" ("cid", "cname", "cprereq")
	values ($1, $2, $3)`
	_, err = db.Exec(insertStatement, d.CID, d.CNAME, d.CPREREQ)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}
	courID := 0
	return courID
}

// DeleteCourse deletes an existing user
func DeleteCourse(id1 string) error {
	db, err := CopenConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Does the ID exist?
	statement := fmt.Sprintf(`SELECT "cid" FROM "msdscoursecatalog" where cid = %d`, id1)
	rows, err := db.Query(statement)

	var coursename string
	for rows.Next() {
		err = rows.Scan(&coursename)
		if err != nil {
			return err
		}
	}
	defer rows.Close()

	if Cexists(coursename) != id {
		return fmt.Errorf("Course with ID %d does not exist", id)
	}

	// Delete from msdscoursecatalog
	deleteStatement := `delete from "msdscoursecatalog" where cid=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}

	// Delete from Users
	/*deleteStatement = `delete from "users" where id=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}
	*/
	return nil
}

// ListCourses lists all Courses in the database
func ListCourses() ([]MSDSCourse, error) {
	Data := []MSDSCourse{}
	fmt.Println(Data)
	db, err := CopenConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT "cid","cname","cprereq"
		FROM "msdscoursecatalog"`)
	if err != nil {
		return Data, err
	}

	for rows.Next() {
		var cID string
		var cNAME string
		var cPREREQ string
		err = rows.Scan(&cID, &cNAME, &cPREREQ)
		temp1 := MSDSCourse{CID: cID, CNAME: cNAME, CPREREQ: cPREREQ}
		Data = append(Data, temp1)
		if err != nil {
			return Data, err
		}
	}
	defer rows.Close()
	return Data, nil
}

// UpdateCourse is for updating an existing user
func UpdateCourse(d MSDSCourse) error {
	db, err := CopenConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	userID := Cexists(d.CNAME)
	if userID == -1 {
		return errors.New("User does not exist")
	}
	d.CID = string(userID)
	updateStatement := `update "msdscoursecatalog" set "cid"=$1, "cname"=$2, "cprereq"=$3 where "cid"=$4`
	_, err = db.Exec(updateStatement, d.CID, d.CNAME, d.CPREREQ, d.CID)
	if err != nil {
		return err
	}

	return nil
}
