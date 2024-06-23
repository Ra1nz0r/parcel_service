package services

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/ra1nz0r/parcel_service/internal/coretypes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func connDB(t *testing.T) (*sql.DB, ParcelStore, coretypes.Parcel) {
	db, err := sql.Open("sqlite", "../database/tracker.db")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	parcel.Number, err = store.Add(parcel)

	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	return db, store, parcel
}

// getTestParcel возвращает тестовую посылку
func getTestParcel() coretypes.Parcel {
	return coretypes.Parcel{
		Client:    1040,
		Status:    coretypes.ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	db, store, parcel := connDB(t)
	defer db.Close()

	stored, err := store.Get(parcel.Number)
	require.NoError(t, err)
	assert.Equal(t, parcel, stored)

	err = store.Delete(parcel.Number)
	require.NoError(t, err)

	_, err = store.Get(parcel.Number)
	require.ErrorIs(t, sql.ErrNoRows, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, store, parcel := connDB(t)
	defer db.Close()

	newAddress := "new test address"
	err := store.SetAddress(parcel.Number, newAddress)

	require.NoError(t, err)

	resAddress, err := store.Get(parcel.Number)
	require.NoError(t, err)
	require.Equal(t, newAddress, resAddress.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, store, parcel := connDB(t)
	defer db.Close()

	err := store.SetStatus(parcel.Number, coretypes.ParcelStatusDelivered)
	require.NoError(t, err)

	resStatus, err := store.Get(parcel.Number)
	require.NoError(t, err)
	require.Equal(t, coretypes.ParcelStatusDelivered, resStatus.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "../database/tracker.db")
	require.NoError(t, err)
	defer db.Close()

	parcels := []coretypes.Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]coretypes.Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	store := NewParcelStore(db)
	for i := 0; i < len(parcels); i++ {

		id, errID := store.Add(parcels[i])
		require.NoError(t, errID)
		require.NotEmpty(t, id)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)

	require.Equal(t, len(parcelMap), len(storedParcels))

	for _, parcel := range storedParcels {
		assert.Contains(t, parcelMap, parcel.Number)
		assert.Equal(t, parcelMap[parcel.Number], parcel)
	}
}
