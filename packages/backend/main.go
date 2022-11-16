package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/boltdb/bolt"
	"github.com/botvid/webapp/abbreviation"
	"github.com/botvid/webapp/backup"
	"github.com/botvid/webapp/database/boltdb"
	"github.com/botvid/webapp/document"
	"github.com/botvid/webapp/repository"
	"github.com/botvid/webapp/user"
	"github.com/botvid/webapp/ws"
	"github.com/gorilla/mux"
)

var cfg *envVariables
var repo repository.Repository

func init() {
	cfg = getEnvVariables(".env")
}

func main() {
	if cfg.mode == "desktop" {
		if cfg.version != "0.0.0" {
			log.Println("Build version", cfg.version)
			if cfg.updateURL != "" {
				doUpdate(cfg.updateURL)
			}
		} else {
			log.Println("Test build 0.0.0")
		}
	}
	repo = boltdb.NewBoltRepository(openBoltDB(cfg.dbFileName))
	err := repo.Init()
	if err != nil {
		log.Fatalf("Couldn't init Bolt DB: %s\n", err)
	}
	r := mux.NewRouter().StrictSlash(true)
	if cfg.mode == "desktop" {
		r.HandleFunc("/ip", getLocalIP).Methods("GET")
		r.PathPrefix("/assets").Handler(
			http.StripPrefix("/assets",
				http.FileServer(
					http.Dir("public/assets"),
				),
			),
		)
		a := App{pools: make(map[string]*ws.Pool)}
		a.Initialize(r)
		log.Println("Run local interpretation")
		go http.ListenAndServe(":8080", nil)
	}

	r.HandleFunc("/", apiHelper)
	http.Handle("/", accessControl(r))
	log.Println("Running as: ", cfg.mode)
	abbService := abbreviation.NewAbbService(repo, cfg.globalStandard)
	abbHandler := abbreviation.NewAbbHandler(abbService)
	abbreviation.Endpoints(r, abbHandler)

	docService := document.NewDocService(repo)
	docHandler := document.NewDocHandler(docService)
	document.Endpoints(r, docHandler)

	userService := user.NewUserService(repo, cfg.mode, abbService, docService)
	userHandler := user.NewUserHandler(userService, cfg.mode)
	user.Endpoints(r, userHandler)
	users, err := user.GetUsers(userService)

	abbService.InitCache(users)

	backupService := backup.NewBackupService(repo)
	backupHandler := backup.NewBackupHandler(backupService)
	backup.Endpoints(r, backupHandler)

	printEndpoints(r)
	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port:", cfg.port)
		errs <- http.ListenAndServe(cfg.port, nil)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Printf("terminated: %s", <-errs)
}

func openBoltDB(fileName string) *bolt.DB {
	log.Println("Opening BoltDB", fileName)
	bolt, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		log.Fatalf("Failed opening BoltDB %s, err: %s\n", fileName, err)
	}
	return bolt

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		///origin := req.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Id-Token, Cache-Control")
		//w.Header().Set("Access-Control-Allow-Origin", "*")

		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, req)
	})
}

func staticFiles(w http.ResponseWriter, r *http.Request) {
	log.Println(r)

}

func apiHelper(w http.ResponseWriter, r *http.Request) {
	if cfg.mode == "desktop" {
		log.Println("Serve consumer")
		log.Println(r)
		http.ServeFile(w, r, "public/")
		return
	}
	w.Write([]byte("This is the API helper speaking, available endpoints/methods are:"))
}

func getLocalIP(w http.ResponseWriter, r *http.Request) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if isPrivateIP(ip) {
				resp := fmt.Sprintf(`{ "ip": "%s" }`, ip)
				w.Write([]byte(resp))
				return
			}

		}
	}
	w.WriteHeader(http.StatusInternalServerError)
	return
}

func isPrivateIP(ip net.IP) bool {
	var privateIPBlocks []*net.IPNet
	for _, cidr := range []string{
		// don't check loopback ips
		//"127.0.0.0/8",    // IPv4 loopback
		//"::1/128",        // IPv6 loopback
		//"fe80::/10",      // IPv6 link-local
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This function is not implemented"))
}
