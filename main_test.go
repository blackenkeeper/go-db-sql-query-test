package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

func Test_SelectClient_WhenOk(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		t.Errorf("Не удалось подключиться к БД. Причина: %s\n", err)
		return
	}
	defer db.Close()

	clientID := 1

	client, err := selectClient(db, clientID)

	require.NoError(t, err)

	assert.Equal(t, clientID, client.ID)
	assert.NotEmpty(t, client.FIO)
	assert.NotEmpty(t, client.Birthday)
	assert.NotEmpty(t, client.Email)
	assert.NotEmpty(t, client.Login)
}

func Test_SelectClient_WhenNoClient(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		t.Errorf("Не удалось подключиться к БД. Причина: %s\n", err)
	}
	defer db.Close()

	clientID := -1
	client, err := selectClient(db, clientID)

	require.Equal(t, sql.ErrNoRows, err)

	assert.Empty(t, client.ID)
	assert.Empty(t, client.FIO)
	assert.Empty(t, client.Birthday)
	assert.Empty(t, client.Email)
	assert.Empty(t, client.Login)
}

func Test_InsertClient_ThenSelectAndCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		t.Errorf("Не удалось подключиться к БД. Причина: %s\n", err)
	}
	defer db.Close()

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	cl.ID, err = insertClient(db, cl)

	require.NotEmpty(t, cl.ID)
	require.NoError(t, err)

	selectClient, err := selectClient(db, cl.ID)
	require.NoError(t, err)

	assert.Equal(t, cl, selectClient)

}

func Test_InsertClient_DeleteClient_ThenCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		t.Errorf("Не удалось подключиться к БД. Причина: %s\n", err)
	}
	defer db.Close()

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	cl.ID, err = insertClient(db, cl)

	require.NotEmpty(t, cl.ID)
	require.NoError(t, err)

	_, err = selectClient(db, cl.ID)
	require.NoError(t, err)

	err = deleteClient(db, cl.ID)
	require.NoError(t, err)

	_, err = selectClient(db, cl.ID)
	require.Equal(t, sql.ErrNoRows, err)

}
