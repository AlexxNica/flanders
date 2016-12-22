package api

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/goji/param"
	"github.com/gorilla/websocket"
	"github.com/weave-lab/flanders/capture"
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

func StartWebServer(address string, assetfolder string) error {

	goji.Use(CORS)

	goji.Get("/search", func(c web.C, w http.ResponseWriter, r *http.Request) {
		filter := db.NewFilter()
		options := &db.Options{}

		r.ParseForm()
		startDate := r.Form.Get("startDate")
		endDate := r.Form.Get("endDate")
		limit := r.Form.Get("limit")
		touser := r.Form.Get("touser")
		fromuser := r.Form.Get("fromuser")
		sourceip := r.Form.Get("sourceip")
		destip := r.Form.Get("destip")

		if startDate != "" {
			filter.StartDate = startDate
		}

		if endDate != "" {
			filter.EndDate = endDate
		}

		if touser != "" {
			filter.Equals["touser"] = touser
		}

		if fromuser != "" {
			filter.Equals["fromuser"] = fromuser
		}

		if sourceip != "" {
			filter.Equals["sourceip"] = sourceip
		}

		if destip != "" {
			filter.Equals["destinationip"] = destip
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

/*		order := r.Form["orderby"]
		if len(order) == 0 {
			options.Sort = append(options.Sort, "-datetime")
		} else {
			options.Sort = order
		}*/

		results, err := db.Db.Find(&filter, options)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//fmt.Print(results)
		jsonResults, err := json.Marshal(results)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		fmt.Fprintf(w, "%s", string(jsonResults))
	})

	goji.Get("/call/:id", func(c web.C, w http.ResponseWriter, r *http.Request) {
		callID := c.URLParams["id"]

		dbResults, err := packetsByCallID(callID, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sort.Sort(dbResults)

		var results db.DbResult

		for key, val := range dbResults {
			if key == 0 || val != dbResults[key-1] {
				results = append(results, val)
			}
		}

		enc := json.NewEncoder(w)
		err = enc.Encode(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	goji.Get("/call/:id/dump", func(c web.C, w http.ResponseWriter, r *http.Request) {
		callId := c.URLParams["id"]
		r.ParseForm()

		ip := r.Form.Get("ip")

		finalResults, err := packetsByCallID(callId, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sort.Sort(finalResults)

		dedupResults := make(db.DbResult, 0, len(finalResults))

		for key, val := range finalResults {
			if key == 0 || val != finalResults[key-1] {
				dedupResults = append(dedupResults, val)
			}
		}

		var dump string

		for _, packet := range dedupResults {
			if packet.SourceIp == ip || packet.DestinationIp == ip {
				dump += "U " + packet.SourceIp + ":" + strconv.Itoa(int(packet.SourcePort)) + " -> " + packet.DestinationIp + ":" + strconv.Itoa(int(packet.DestinationPort)) + "\n"
				dump += packet.Msg
				dump += "\n\n"
			}
		}

		w.Header().Set("Content-Disposition", "attachment; filename=dump.txt")
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

		fmt.Fprintf(w, "%s", dump)

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
		newfilter := capture.RegisterListener(filter)
	forloop:
		for {
			select {
			case dbObject := <-newfilter.Broadcast:
				//jsonMessage, err := json.Marshal(dbObject)
				conn.WriteJSON(dbObject)
			case <-newfilter.Quit:
				conn.Close()
				break forloop
			}
		}
	})

	goji.Get("/settings/:group", func(c web.C, w http.ResponseWriter, r *http.Request) {
		group := c.URLParams["group"]

		results, err := db.Db.GetSettings(group)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResults, err := json.Marshal(results)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		fmt.Fprintf(w, "%s", string(jsonResults))

	})

	goji.Post("/settings/:group", func(c web.C, w http.ResponseWriter, r *http.Request) {
		group := c.URLParams["group"]
		r.ParseForm()
		fmt.Println(r.PostForm)
		setting := db.SettingObject{}
		err := param.Parse(r.Form, &setting)
		fmt.Printf("%#v", setting)
		if err != nil {
			log.Err(err.Error())
			return
		}
		dberr := db.Db.SetSetting(group, setting)
		if dberr != nil {
			log.Err(dberr.Error())
			return
		}

		jsonResults, err := json.Marshal(setting)
		fmt.Fprintf(w, "%s", string(jsonResults))
	})

	goji.Delete("/settings/:group/:key", func(c web.C, w http.ResponseWriter, r *http.Request) {
		group := c.URLParams["group"]
		key := c.URLParams["key"]

		dberr := db.Db.DeleteSetting(group, key)
		if dberr != nil {
			log.Err(dberr.Error())
			fmt.Fprintf(w, "%s", "{ \"result\": false, \"error\": \""+dberr.Error()+"\" }")
			return
		}
		fmt.Fprintf(w, "%s", "{ \"result\": true }")
	})

	// goji.Post("/settings/alias", func(c web.C, w http.ResponseWriter, r *http.Request) {

	// 	//r.ParseForm()
	// 	name := r.FormValue("name")
	// 	ip := r.FormValue("ip")

	// 	var results []db.DbObject
	// 	db.Db.Find(&filter, options, &results)

	// 	jsonResults, err := json.Marshal(results)
	// 	if err != nil {
	// 		fmt.Fprint(w, err)
	// 		return
	// 	}

	// 	fmt.Fprintf(w, "%s", string(jsonResults))
	// })

	// goji.Put("/settings/alias/:id", func(c web.C, w http.ResponseWriter, r *http.Request) {
	// 	//r.ParseForm()
	// 	aliusId := c.URLParams["id"]
	// 	name := r.FormValue("name")
	// 	ip := r.FormValue("ip")
	// })

	// goji.Delete("/settings/alias/:id", func(c web.C, w http.ResponseWriter, r *http.Request) {
	// 	aliusId := c.URLParams["id"]
	// })

	goji.Get("/*", http.FileServer(http.Dir(assetfolder)))
	flag.Set("bind", address)
	goji.Serve()

	return nil
}

func packetsByCallID(callID string, excludeCallID string) (db.DbResult, error) {
	var results db.DbResult
	filter := db.NewFilter()
	options := &db.Options{}
	callIdMap := make(map[string]interface{})
	callIdALegMap := make(map[string]interface{})

	callIdMap["callid"] = callID
	callIdALegMap["callidaleg"] = callID

	filter.Equals["$or"] = []interface{}{
		callIdMap,
		callIdALegMap,
	}

	options.Sort = append(options.Sort, "datetime")
	options.Sort = append(options.Sort, "microseconds")
	results, err := db.Db.Find(&filter, options)
	if err != nil {
		return nil, err
	}

	altCallIds := make(map[string]bool)
	for _, msg := range results {
		if excludeCallID != "" && msg.CallIdAleg == excludeCallID {
			continue
		}
		if msg.CallId != callID {
			_, ok := altCallIds[msg.CallId]
			if !ok {
				altCallIds[msg.CallId] = true
			}
		}
	}
	for newCallID, _ := range altCallIds {
		p, err := packetsByCallID(newCallID, callID)
		if err != nil {
			return nil, err
		}

		results = append(results, p...)
	}

	return results, nil
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
