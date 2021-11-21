package mysql

import (
	"reflect"
	"testing"
	"time"

	"kerseeeHuang.com/snippetbox/pkg/models"
)

func TestUserModelGet(t *testing.T) {
	// Skip the integration test if the -test.short flag is set.
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	// Create test cases.
	tests := []struct{
		name		string
		userID		int
		wantUser	*models.User
		wantError	error
	}{
		{
			name:		"Valid ID",
			userID:		1,
			wantUser:	&models.User {
				ID:			1,
				Name: 		"Alice Jones",
				Email:		"alice@example.com",
				Created:	time.Date(2021, 11, 21, 17, 8, 0, 0, time.UTC),
				Active: 	true,
			},
			wantError: nil,
		},
		{
			name:		"Zero ID",
			userID:		0,
			wantUser: 	nil,
			wantError: 	models.ErrNoRecord,
		},
		{
			name:		"Non-existent ID",
			userID:		2,
			wantUser: 	nil,
			wantError: 	models.ErrNoRecord,
		},
	}

	// Run test cases.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){
			// Initialize the connection poll to the test DB.
			db, teardown := newTestDB(t)
			defer teardown()

			// Initialize a userModel.
			m := UserModel{db}

			// Test the test case.
			user, err := m.Get(test.userID)
			if err != test.wantError {
				t.Errorf("want %v; got %v", test.wantError, err)
			}
			if !reflect.DeepEqual(user, test.wantUser) {
				t.Errorf("want %v; got %v", test.wantUser, user)
			}
		})
	}
}