package handlers

import (
	"hatnote-historical/db"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	db  *db.DB
	log *log.Logger
}

func UnixToStringTimestamp(unix string) string {
	i, err := strconv.ParseInt(unix, 10, 64)
	if err != nil {
		panic(err)
	}
	t := time.Unix(i, 0)
	return t.String()[:19]
}

func (h *Handlers) NetChangePerPeriod(c *gin.Context) {
	langCode := c.Query("langCode")
	if langCode == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No langCode url param",
		})
		return
	}
	from := c.Query("from")
	to := c.Query("to")

	if to == "" || from == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Both 'to' and 'from' url params are required",
		})
	}

	sqlStatement := `
  SELECT byte_change FROM edits WHERE lang_code = $1 AND modified_at > $2 AND modified_at < $3`
	rows, err := h.db.Query(sqlStatement,
		langCode,
		UnixToStringTimestamp(from),
		UnixToStringTimestamp(to))

	defer func() {
		if err = rows.Close(); err != nil {
			panic(err)
		}
	}()

	var changes []int
	for rows.Next() {
		byteChange := new(int)
		if err = rows.
			Scan(&byteChange); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Error scanning db row",
			})
			return
		}
		changes = append(changes, *byteChange)
	}
	if err = rows.Err(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error scanning db row",
		})
		return
	}
	if changes == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "No data available for this time range",
		})
		return
	}

	netChange := 0
	for _, change := range changes {
		netChange += change
	}

	c.JSON(http.StatusOK, gin.H{
		"netChange": netChange,
	})
}

func New(db *db.DB, log *log.Logger) *Handlers {
	return &Handlers{
		db:  db,
		log: log,
	}
}
