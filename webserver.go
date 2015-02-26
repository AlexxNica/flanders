package flanders

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/weave-lab/flanders/db"
	"github.com/weave-lab/flanders/log"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebServer(ip string, port int) {

	goji.Use(CORS)

	goji.Get("/search", func(c web.C, w http.ResponseWriter, r *http.Request) {
		filter := db.NewFilter()
		options := &db.Options{}

		r.ParseForm()
		startDate := r.Form.Get("startDate")
		endDate := r.Form.Get("endDate")
		limit := r.Form.Get("limit")
		touser := r.Form.Get("touser")
		todomain := r.Form.Get("todomain")
		fromuser := r.Form.Get("fromuser")
		fromdomain := r.Form.Get("fromdomain")

		if startDate != "" {
			filter.StartDate = startDate
		}

		if endDate != "" {
			filter.EndDate = endDate
		}

		if touser != "" {
			filter.Equals["touser"] = touser
		}

		if todomain != "" {
			filter.Equals["todomain"] = fromuser
		}

		if fromuser != "" {
			filter.Equals["fromuser"] = touser
		}

		if fromdomain != "" {
			filter.Equals["fromdomain"] = fromuser
		}

		if limit == "" {
			options.Limit = 50
		} else {
			limitUint64, err := strconv.Atoi(limit)
			if err != nil {
				options.Limit = 50
			} else {
				options.Limit = limitUint64
			}
		}

		order := r.Form["orderby"]
		if len(order) == 0 {
			options.Sort = append(options.Sort, "-datetime")
		} else {
			options.Sort = order
		}

		var results []db.DbObject

		db.Db.Find(&filter, options, &results)
		jsonResults, err := json.Marshal(results)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		fmt.Fprintf(w, "%s", string(jsonResults))
	})

	goji.Get("/call/:id", func(c web.C, w http.ResponseWriter, r *http.Request) {
		callId := c.URLParams["id"]
		log.Debug("Call ID: " + callId)
		filter := db.NewFilter()
		options := &db.Options{}

		filter.Equals["callid"] = callId

		options.Sort = append(options.Sort, "datetime")

		var results []db.DbObject
		db.Db.Find(&filter, options, &results)

		jsonResults, err := json.Marshal(results)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		fmt.Fprintf(w, "%s", string(jsonResults))

	})

	goji.Get("/ws", func(c web.C, w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		filter := r.Form.Get("filter")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Err(err.Error())
			return
		}
		go func() {
			for {
				if _, _, err := conn.NextReader(); err != nil {
					conn.Close()
					break
				}
			}
		}()
		newfilter := RegisterListener(filter)
	forloop:
		for {
			select {
			case dbObject := <-newfilter.Broadcast:
				//jsonMessage, err := json.Marshal(dbObject)
				conn.WriteJSON(dbObject)
			case <-newfilter.quit:
				conn.Close()
				break forloop
			}
		}
	})

	goji.Get("/*", http.FileServer(http.Dir("www")))
	flag.Set("bind", ip+":"+strconv.Itoa(port))
	goji.Serve()
}

func CORS(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST,PUT,GET,HEAD,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Max-Age", "1728000")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
