package users

import (
	"testing"
	"golang.org/x/crypto/bcrypt"

)

func TestMongoStore(t *testing.T) {
	store := NewMongoStore(nil, "test", "test")
	if store == nil {
		t.Fatal("Failed creating User Storage")
	}
	user1 := &NewUser{
		Email        : "user1@user1.com",
		Password     : "password",
		PasswordConf : "password",
		UserName     : "username1",
		FirstName    : "user1",
		LastName     : "one",
	}
	user2 := &NewUser {
		Email : "user2@user2.com",
		Password : "password",
		PasswordConf: "password",
		UserName : "username2",
		FirstName : "user2",
		LastName: "two",
	}

	userone, err := store.Insert(user1)
	if err != nil {
		t.Fatalf("Failed inserting new user: %v\n", err)
	}
	if userone == nil {
		t.Fatalf("user not inserted correctly: %v\n", err)
	}

	usertwo, err := store.Insert(user2)
	if err != nil {
		t.Fatalf("Failed inserting new user: %v\n", err)
	}
	if usertwo == nil {
		t.Fatalf("user not inserted correctly: %v\n", err)
	}
	users, err := store.GetAll()
	if err != nil {
		t.Fatalf("Failed retrieving all Users %v\n", err)
	}
	if users == nil || len(users) != 2{
		t.Fatalf("users not retrieved correctly %v\n", err)
	}

	userthree, err := store.Insert(user2)
	if err == nil {
		t.Fatal("Same Users shouldn't be inserted twice")
	}
	if userthree != nil {
		t.Fatal("Same Users shouldn't be inserted twice")
	}

	iuser, err := store.GetByID(userone.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve user with given ID: %v\n", err)
	}
	if iuser.FirstName != userone.FirstName {
		t.Fatalf("Retrieved user does not match the given ID: %v\n", err)
	}

	euser, err := store.GetByEmail("user1@user1.com")
	if err != nil {
		t.Fatalf("Failed to retrieve user with given email: %v\n", err)
	} 
	if euser.FirstName != user1.FirstName {
		t.Fatalf("Retrieved user does not match the given email")
	}

	update := &UserUpdates{
		FirstName: "UPDATE",
		LastName: "UPDATE",
	}
	if err := store.Update(update, userone); err != nil {
		t.Fatalf("Failed to retrieve user with given ID: %v\n", err)
	}
	user1U, err := store.GetByID(userone.ID)
	if err != nil {
		t.Fatalf("Error retrieiving user with given ID")
	}
	if user1U.FirstName != "UPDATE" {
		t.Fatalf("Incorrect Update! Value Retrieved: %v Should be UPDATE", userone.FirstName)
	}
	if user1U.LastName != "UPDATE" {
		t.Fatalf("Incorrect Update! Value Retrieved: %v Should be UPDATE", userone.LastName)
	}

	if err := store.ResetPassword("user2@user2.com", "NEWPASSWORD"); err != nil {
		t.Fatalf("Error occurred while updating password")
	}
	user2R, err := store.GetByEmail("user2@user2.com")
	if err != nil {
		t.Fatalf("Failed to retrieve user with given email: %v\n", err)
	}
	if err := bcrypt.CompareHashAndPassword(user2R.PassHash, []byte("NEWPASSWORD")); err != nil {
		t.Fatalf("PASSWORD NOT RESET CORRECTLY")
	}
}