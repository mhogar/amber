package runner

import (
	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/controllers"
	"github.com/mhogar/amber/data"
)

// Run runs the admin creator and returns any errors.
func Run(sf data.ScopeFactory, c controllers.UserController, username string, password string, rank int) error {
	return sf.CreateDataExecutorScope(func(exec data.DataExecutor) error {
		return sf.CreateTransactionScope(exec, func(tx data.Transaction) (bool, error) {
			//create the user
			_, cerr := c.CreateUser(tx, username, password, rank)
			if cerr.Type != common.ErrorTypeNone {
				return false, common.ChainError("error creating user", cerr)
			}

			return true, nil
		})
	})
}
