package runner

import (
	"authserver/common"
	"authserver/controllers"
	"authserver/data"
)

// Run runs the admin creator. Returns any errors
func Run(sf data.ScopeFactory, c controllers.UserController, username string, password string) error {
	return sf.CreateDataExecutorScope(func(exec data.DataExecutor) error {
		return sf.CreateTransactionScope(exec, func(tx data.Transaction) (bool, error) {
			//create the user
			_, rerr := c.CreateUser(tx, username, password)
			if rerr.Type != common.ErrorTypeNone {
				return false, common.ChainError("error creating user", rerr)
			}

			return true, nil
		})
	})
}
