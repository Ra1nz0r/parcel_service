package services

import (
	"database/sql"
	"errors"
	"log"

	"github.com/ra1nz0r/parcel_service/internal/coretypes"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p coretypes.Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	lastAddID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastAddID), nil
}

func (s ParcelStore) Get(number int) (coretypes.Parcel, error) {
	resQuery := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))

	p := coretypes.Parcel{}
	err := resQuery.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			log.Println("Record does not exist in the database.")
			return p, err
		}
		log.Fatal(err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]coretypes.Parcel, error) {
	resQuery, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer resQuery.Close()

	var res []coretypes.Parcel

	for resQuery.Next() {
		p := coretypes.Parcel{}
		errQuery := resQuery.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if errQuery != nil {
			return nil, err
		}
		res = append(res, p)
	}

	if err = resQuery.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", coretypes.ParcelStatusRegistered))

	return err
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", coretypes.ParcelStatusRegistered))

	return err
}
