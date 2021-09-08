package handlers

import (
	"authserver/common"
	"authserver/data"
	"authserver/models"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type PutClientRolesBody struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (h CoreHandlers) PutClientRoles(req *http.Request, params httprouter.Params, _ *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	var rolesBody []PutClientRolesBody

	//parse the id
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		log.Println(common.ChainError("error parsing id", err))
		return common.NewBadRequestResponse("client id is in an invalid format")
	}

	//parse the body
	err = parseJSONBody(req.Body, &rolesBody)
	if err != nil {
		log.Println(common.ChainError("error parsing PutClientRoles request body", err))
		return common.NewBadRequestResponse("invalid json body")
	}

	//create the user-role models
	userRoles := make([]*models.UserRole, len(rolesBody))
	for index, role := range rolesBody {
		userRoles[index] = models.CreateUserRole(role.Username, id, role.Role)
	}

	//update the roles
	cerr := h.Controllers.UpdateUserRolesForClient(CRUD, id, userRoles)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}
