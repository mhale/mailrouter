package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"mime"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mhale/smtpd"
	"github.com/streadway/simpleuuid"
)

var (
	config Config  // Filters & routes
	stats  Stats   // Statistics of sent and dropped mail
	logs   LogList // Recent mail log for Dashboard
)

var httpAddr *string = flag.String("http", ":8080", "Address & port for HTTP server")
var smtpAddr *string = flag.String("smtp", ":2525", "Address & port for SMTP server")
var confFile *string = flag.String("conf", "/etc/mailrouter.conf", "Full path to configuration file")

// Handler for handling incoming mail messages.
func mailHandler(origin net.Addr, from string, to []string, data []byte) {
	originIPStr, _, _ := net.SplitHostPort(origin.String())
	originIP := net.ParseIP(originIPStr)

	// Parse the message to get the Subject header.
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		log.Printf("Failed to parse message: %s\n", err)
		log.Printf("Aborting processing of message.")
		return
	}
	subject := msg.Header.Get("Subject")

	// Check each filter in order.
	var filterName string
	var routeId string
	for _, filter := range SortedFilters() {
		if filter.Match(from, to, subject, originIP) {
			filterName = filter.Name
			routeId = filter.RouteId
			break
		}
	}

	// Use the default route if no filters were matched.
	if routeId == "" {
		for _, route := range config.Routes {
			if route.IsDefault {
				routeId = route.Id
			}
		}
	}

	// If the message is to be dropped, record the drop and return.
	if routeId == "DROP" {
		stats.Dropped(len(data))
		logs.Add(originIP, from, to, subject, filterName, "Drop")
		return
	}

	// Otherwise, deliver the mail to the selected route and record the delivery.
	route := config.Routes[routeId]
	addr := route.Hostname + ":" + strconv.Itoa(route.Port)

	// Override the recipient if To field is set.
	if route.To != "" {
		to = []string{route.To}
	}

	// Deliver the mail.
	var auth smtp.Auth
	auth = nil
	if route.AuthType == "plain" {
		auth = smtp.PlainAuth("", route.Username, route.Password, route.Hostname)
	} else if route.AuthType == "crammd5" {
		auth = smtp.CRAMMD5Auth(route.Username, route.Password)
	}

	err = smtp.SendMail(addr, auth, from, to, data)
	if err != nil {
		log.Printf("Failed to deliver mail to route %s (%s): %s", route.Name, addr, err)
		stats.Failed(len(data))
		logs.Add(originIP, from, to, subject, filterName, "Failed")
		return
	}

	// Record the successful delivery.
	stats.Sent(len(data))
	logs.Add(originIP, from, to, subject, filterName, route.Name)
}

func routeHandler(w http.ResponseWriter, req *http.Request) {
	var msg string
	_, id, action := ParsePath(req.URL.Path)
	method := req.FormValue("_method")

	if req.Method == "GET" {
		data := make(map[string]interface{})
		data["list"] = SortedRoutes()

		// Populate the form if requested.
		if id != "" && action == "edit" {
			data["id"] = id
			data["edit"] = config.Routes[id]
		}

		// Check for info and error messages passed via cookies. Clear any that are displayed.
		msg = GetCookie(w, req, "info")
		if msg != "" {
			data["info"] = msg
		}
		msg = GetCookie(w, req, "error")
		if msg != "" {
			data["error"] = msg
		}

		// Render the page. Reparsing the template every time eases development at the expense of performance.
		html, _ := Asset("views/routes.html")
		tmpl, err := template.New("routes").Parse(string(html))
		if err != nil {
			log.Println(err)
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}

	if req.Method == "POST" {
		config.Lock()
		defer config.Unlock()

		if method == "delete" {
			msg = fmt.Sprintf("Deleted route %s.", config.Routes[id].Name)
			// If a default route was deleted, make Drop the default.
			if config.Routes[id].IsDefault == true {
				route := config.Routes["DROP"]
				route.IsDefault = true
				config.Routes["DROP"] = route
				msg = fmt.Sprintf("%s The Drop route is now the default route.", msg)
			}
			delete(config.Routes, id)
		}

		if method == "default" {
			msg = fmt.Sprintf("The %s route is now the default route.", config.Routes[id].Name)
			for id, route := range config.Routes {
				route.IsDefault = false
				config.Routes[id] = route
			}
			route := config.Routes[id]
			route.IsDefault = true
			config.Routes[id] = route
		}

		if method == "save" {
			// Unset id means a new route is being added
			if id == "" {
				msg = fmt.Sprintf("Added route %s to host %s.", req.FormValue("routename"), req.FormValue("hostname"))
				uuid, _ := simpleuuid.NewTime(time.Now())
				id = uuid.String()
			} else {
				msg = fmt.Sprintf("Updated route %s.", req.FormValue("routename"))
			}

			// Create a new Route from the form submission.
			port, _ := strconv.Atoi(req.FormValue("port"))
			isDefault, _ := strconv.ParseBool(req.FormValue("isdefault"))
			route := Route{
				Id:        id,
				Name:      req.FormValue("routename"),
				To:        req.FormValue("to"),
				Hostname:  req.FormValue("hostname"),
				Port:      port,
				AuthType:  req.FormValue("authentication"),
				Username:  req.FormValue("username"),
				Password:  req.FormValue("password"),
				IsDefault: isDefault,
			}
			config.Routes[id] = route
		}

		if msg != "" {
			log.Printf(msg)
			SetCookie(w, "info", msg)
			err := SaveConfig()
			if err != nil {
				msg = fmt.Sprintf("Failed to save configuration to file: %v", err)
				log.Printf(msg)
				SetCookie(w, "error", msg)
			}
		}

		http.Redirect(w, req, "/routes/", http.StatusFound)
	}
}

func filterHandler(w http.ResponseWriter, req *http.Request) {
	var msg string
	_, id, action := ParsePath(req.URL.Path)
	method := req.FormValue("_method")

	if req.Method == "GET" {
		data := make(map[string]interface{})
		data["list"] = SortedFilters()
		data["routes"] = SortedRoutes()

		if len(config.Routes) == 1 {
			data["info"] = "No routes are defined. It is recommended to define routes before filters to populate the route drop-down menu below."
		}

		// Populate the form if requested.
		if id != "" && action == "edit" {
			data["id"] = id
			data["edit"] = config.Filters[id]
		}

		// Check for info and error messages passed via cookies. Clear any that are displayed.
		msg = GetCookie(w, req, "info")
		if msg != "" {
			data["info"] = msg
		}
		msg = GetCookie(w, req, "error")
		if msg != "" {
			data["error"] = msg
		}

		// Render the page. Reparsing the template every time eases development at the expense of performance.
		html, _ := Asset("views/filters.html")
		tmpl, err := template.New("filters").Parse(string(html))
		if err != nil {
			log.Println(err)
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}

	if req.Method == "POST" {
		config.Lock()
		defer config.Unlock()

		if method == "delete" {
			msg = fmt.Sprintf("Deleted filter %s.", config.Filters[id].Name)
			delete(config.Filters, id)
		}

		if method == "save" {
			// Unset id means a new filter is being added
			if id == "" {
				msg = fmt.Sprintf("Added filter %s.", req.FormValue("filtername"))
				uuid, _ := simpleuuid.NewTime(time.Now())
				id = uuid.String()
			} else {
				msg = fmt.Sprintf("Updated filter %s.", req.FormValue("filtername"))
			}

			// Create a new Filter from the form submission.
			order, _ := strconv.Atoi(req.FormValue("order"))
			filter := Filter{
				Id:      id,
				Order:   order,
				Name:    req.FormValue("filtername"),
				To:      req.FormValue("to"),
				From:    req.FormValue("from"),
				Origin:  req.FormValue("origin"),
				Subject: req.FormValue("subject"),
				RouteId: req.FormValue("route-id"),
			}
			filter.Summary = filter.Summarise()
			filter.RouteName = config.Routes[filter.RouteId].Name
			config.Filters[id] = filter
		}

		if msg != "" {
			log.Printf(msg)
			SetCookie(w, "info", msg)
			err := SaveConfig()
			if err != nil {
				msg = fmt.Sprintf("Failed to save configuration to file: %v", err)
				log.Printf(msg)
				SetCookie(w, "error", msg)
			}
		}

		http.Redirect(w, req, "/filters/", http.StatusFound)
	}
}

// Handler for serving the Dashboard.
func indexHandler(w http.ResponseWriter, req *http.Request) {
	// Catch bad URLs.
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	data := make(map[string]interface{})
	data["stats"] = stats
	data["logs"] = logs.Logs
	data["maxLogs"] = MaxLogs

	html, _ := Asset("views/index.html")
	tmpl, err := template.New("index").Parse(string(html))
	if err != nil {
		log.Println(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

// Handler for serving static assets (CSS/JS).
func assetHandler(w http.ResponseWriter, req *http.Request) {
	path := string(req.URL.Path[1:]) // Strip leading slash
	data, err := Asset(path)
	if err != nil {
		log.Printf("Asset not found: %s", path)
		http.NotFound(w, req)
		return
	}

	ext := filepath.Ext(path)
	contentType := mime.TypeByExtension(ext)
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

func main() {
	flag.Parse()

	// Load filters & routes from configuration file.
	err := LoadConfig()
	if err != nil {
		log.Printf("Could not load configuration file: %s", err)
		log.Printf("No routes or filters are defined. All incoming mail will be dropped.")
	} else {
		// Subtract 1 from config.Routes to account for Drop route.
		log.Printf("Loaded %d routes and %d filters.", len(config.Routes)-1, len(config.Filters))
	}

	// Create a PID file.
	if config.Options["PIDFile"] != "" {
		err := CreatePIDFile()
		if err != nil {
			log.Printf("Could not create PID file: %s", err)
		} else {
			defer RemovePIDFile()
		}
	}

	// Run HTTP server in the background.
	log.Printf("Mailrouter serving HTTP on %s", *httpAddr)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/assets/", assetHandler)
	http.HandleFunc("/routes/", routeHandler)
	http.HandleFunc("/filters/", filterHandler)
	go http.ListenAndServe(*httpAddr, nil)

	// Run SMTP server in the foreground to force an exit if it fails.
	log.Printf("Mailrouter serving SMTP on %s", *smtpAddr)
	err = smtpd.ListenAndServe(*smtpAddr, mailHandler, "Mailrouter", "")
	if err != nil {
		log.Printf("smtpd.ListenAndServe error: %v", err)
	}

	log.Println("Exiting.")
}
