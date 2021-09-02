package integration_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserRoleCRUDTestSuite struct {
	CRUDTestSuite
}

func TestUserRoleCRUDTestSuite(t *testing.T) {
	suite.Run(t, &UserRoleCRUDTestSuite{})
}
