package main

import (
	"log"

	"github.com/ashyfun/coffeezone"
	"github.com/jackc/pgx/v5"
)

type Pagination struct {
	Limit int `json:"limit" form:"limit,omitempty" binding:"numeric"`
	Page  int `json:"page" form:"page" binding:"numeric"`
	Total int `json:"total"`
}

func CafesHandler(c *Context) {
	pag := &Pagination{
		Limit: 30,
		Page:  1,
	}
	if err := c.ShouldBindQuery(pag); err != nil {
		log.Printf("Failed to bind query: %v", err)
		c.BadRequest(err)
		return
	}

	var cafes []*coffeezone.CafeModel
	coffeezone.QueryExec(func(r pgx.Rows, err error) {
		defer r.Close()
		if err != nil {
			log.Printf("Failed to get cafes: %v", err)
			c.BadRequest(err)
			return
		}

		for r.Next() {
			cafeModel := &coffeezone.CafeModel{}
			err := r.Scan(
				&cafeModel.Code,
				&cafeModel.Title,
				&cafeModel.Location.Address,
				&cafeModel.Location.Longitude,
				&cafeModel.Location.Latitude,
				&cafeModel.Topics,
				&pag.Total,
			)
			if err != nil {
				log.Printf("Failed to get cafe values: %v", err)
				c.BadRequest(err)
				return
			}

			cafes = append(cafes, cafeModel)
		}
	}, `
	select code,
		title,
		address_name,
		longitude,
		latitude,
		string_agg(tpcs.feature, ', ' order by tpcs.id) as topics,
		count(*) over() as total_rows
	from cz_cafes
	left join cz_cafe_locations as cafe_lcts on location_id = cafe_lcts.id
	left join cz_cafes_topics as cafe_tpcs on code = cafe_tpcs.cafe_code
	left join cz_topics as tpcs on cafe_tpcs.topic_id = tpcs.id
	group by code, address_name, longitude, latitude
	order by code asc
	limit $1
	offset $2
	`, pag.Limit, pag.Page*pag.Limit)

	c.Response(cafes, pag)
}
