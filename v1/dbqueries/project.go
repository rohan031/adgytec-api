package dbqueries

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

// create project
// const CreateProject = `INSERT INTO project (project_name)
//
//	VALUES (@projectName)`
const CreateProject = `WITH inserted_row AS (
    INSERT INTO project (project_name)
    VALUES (@projectName)
    RETURNING project_id
)
INSERT INTO client_token (token, project_id)
SELECT @clientToken, project_id
FROM inserted_row;`

func CreateProjectArgs(projectName, clientToken string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"projectName": projectName,
		"clientToken": clientToken,
	}
}

// add services to project
func AddServicesToProject(projectId string, services []string) string {
	query := "INSERT INTO project_to_service (project_id, service_id) VALUES "
	var values []string

	for _, id := range services {
		value := fmt.Sprintf("('%v', '%v')", projectId, id)
		values = append(values, value)
	}

	query += strings.Join(values, ", ")

	return query
}

// add a user to project
const AddUserToProject = `INSERT INTO user_to_project (user_id, project_id)
	VALUES (@userId, @projectId)
`

func AddUserToProjectArgs(userId, projectId string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"userId":    userId,
		"projectId": projectId,
	}
}

// get project id  by token
const GetProjectIdByClientToken = `select project_id from client_token where token=@clientToken`

func GetProjectIdByClientTokenArgs(clientToken string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"clientToken": clientToken,
	}
}

const GetProjectIdByUserId = `select project_id from user_to_project where user_id=@userId`

func GetProjectIdByUserIdArgs(userId string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"userId": userId,
	}
}

const GetProjectById = `select project_name from project where project_id=@projectId`

func GetProjectByIdArgs(projectId string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"projectId": projectId,
	}
}