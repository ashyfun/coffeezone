package coffeezone

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type LocationModel struct {
	Address   string  `json:"address"`
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type CafeModel struct {
	Code     string        `json:"code"`
	Title    string        `json:"title"`
	Location LocationModel `json:"location,omitempty"`
}

func CafesHandler(c *gin.Context) {
	var cafes []*CafeModel
	QueryExec(func(r pgx.Rows, err error) {
		defer r.Close()
		if err != nil {
			log.Printf("Failed to get cafes: %v", err)
		}

		for r.Next() {
			v, err := r.Values()
			if err != nil {
				log.Printf("Failed to get cafe values: %v", err)
				return
			}

			cafes = append(
				cafes,
				&CafeModel{
					v[0].(string),
					v[1].(string),
					LocationModel{v[2].(string), v[3].(float64), v[4].(float64)},
				},
			)
		}
	}, `
	select code, title, address_name, longitude, latitude from cz_cafes
	left join cz_cafe_locations
	on location_id = id
	order by code asc
	`)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    cafes,
	})
}
