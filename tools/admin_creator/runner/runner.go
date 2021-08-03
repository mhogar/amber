package runner

import (
	"authserver/common"
	"authserver/controllers"
	"authserver/data"
	"fmt"
)

// Run runs the admin creator. Returns any errors
func Run(sf data.ScopeFactory, c controllers.UserController, username string, password string) error {
	return sf.CreateDataExecutorScope(func(exec data.DataExecutor) error {
		return sf.CreateTransactionScope(exec, func(tx data.Transaction) (bool, error) {
			//save the user
			user, rerr := c.CreateUser(tx, username, password)
			if rerr.Type != common.ErrorTypeNone {
				return false, common.ChainError("error creating user", rerr)
			}

			fmt.Println("Created user:", user.ID.String())
			return true, nil
		})
	})
}
