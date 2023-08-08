package main

import (
	"log"

	"github.com/ashyfun/coffeezone"
	"github.com/jackc/pgx/v5"
)

func CafesHandler(c *Context) {
	var (
		cafes []*coffeezone.CafeModel
		_err  error = nil
	)
	coffeezone.QueryExec(func(r pgx.Rows, err error) {
		defer r.Close()
		if err != nil {
			log.Printf("Failed to get cafes: %v", err)
			_err = err
			return
		}

		for r.Next() {
			cafeModel := &coffeezone.CafeModel{}
			_err = r.Scan(
				&cafeModel.Code,
				&cafeModel.Title,
				&cafeModel.Location.Address,
				&cafeModel.Location.Longitude,
				&cafeModel.Location.Latitude,
			)
			if _err != nil {
				log.Printf("Failed to get cafe values: %v", err)
				cafes = nil
				return
			}

			cafes = append(cafes, cafeModel)
		}
	}, `
	select code, title, address_name, longitude, latitude from cz_cafes
	left join cz_cafe_locations
	on location_id = id
	order by code asc
	`)

	c.Response(cafes, _err)
}