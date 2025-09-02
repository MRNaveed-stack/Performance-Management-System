package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Skill struct {
	SkillID   int    `json:"skill_id"`
	SkillName string `json:"skill_name"`
}

type EmployeeSkill struct {
	EmployeeID int `json:"employee_id"`
	SkillID    int `json:"skill_id"`
}

// CreateSkill inserts a new skill into the skills table
func CreateSkill(db *pgxpool.Pool, skillName string) (int, error) {
	var id int
	err := db.QueryRow(context.Background(),
		`INSERT INTO skills (skill_name) VALUES ($1) RETURNING skill_id`,
		skillName).Scan(&id)
	return id, err
}

// AddEmployeeSkill assigns an existing skill to an employee
func AddEmployeeSkill(db *pgxpool.Pool, empID, skillID int) error {
	_, err := db.Exec(context.Background(),
		`INSERT INTO employee_skills (employee_id, skill_id) VALUES ($1, $2)`,
		empID, skillID)
	return err
}

// GetSkillsByEmployee fetches all skill records assigned to a given employee
func GetSkillsByEmployee(db *pgxpool.Pool, empID int) ([]Skill, error) {
	query := `
		SELECT s.skill_id, s.skill_name
		FROM employee_skills es
		JOIN skills s ON es.skill_id = s.skill_id
		WHERE es.employee_id = $1
	`

	rows, err := db.Query(context.Background(), query, empID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []Skill
	for rows.Next() {
		var s Skill
		if err := rows.Scan(&s.SkillID, &s.SkillName); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}
	return skills, nil
}

// DeleteEmployeeSkill removes a skill assignment for an employee
func DeleteEmployeeSkill(db *pgxpool.Pool, empID, skillID int) error {
	_, err := db.Exec(context.Background(),
		`DELETE FROM employee_skills WHERE employee_id = $1 AND skill_id = $2`,
		empID, skillID)
	return err
}
