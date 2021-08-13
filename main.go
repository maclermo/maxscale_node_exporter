package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	configUsername string
	configPassword string
	configHost     string
	configPort     int
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type Server struct {
	Data []ServerData `json:"data"`
}

type ServerData struct {
	ID         string           `json:"id"`
	Attributes ServerAttributes `json:"attributes"`
}

type ServerAttributes struct {
	Statistics ServerStatistics `json:"statistics"`
}

type ServerStatistics struct {
	ActiveOperations      int    `json:"active_operations"`
	AdaptiveAvgSelectTime string `json:"adaptive_avg_select_time"`
	ConnectionPoolEmpty   int    `json:"connection_pool_empty"`
	Connections           int    `json:"connections"`
	MaxConnections        int    `json:"max_connections"`
	MaxPoolSize           int    `json:"max_pool_size"`
	PersistentConnections int    `json:"persistent_connections"`
	ReusedConnections     int    `json:"reused_connections"`
	RoutedPackets         int    `json:"routed_packets"`
	TotalConnections      int    `json:"total_connections"`
}

type ServerCollector struct {
	ActiveOperations      *prometheus.Desc
	AdaptiveAvgSelectTime *prometheus.Desc
	ConnectionPoolEmpty   *prometheus.Desc
	Connections           *prometheus.Desc
	MaxConnections        *prometheus.Desc
	MaxPoolSize           *prometheus.Desc
	PersistentConnections *prometheus.Desc
	ReusedConnections     *prometheus.Desc
	RoutedPackets         *prometheus.Desc
	TotalConnections      *prometheus.Desc
}

type Service struct {
	Data []ServiceData `json:"data"`
}

type ServiceData struct {
	Attributes ServiceAttributes `json:"attributes"`
	ID         string            `json:"id"`
}

type ServiceAttributes struct {
	RouterDiagnostics ServiceRouterDiagnostics `json:"router_diagnostics"`
	Statistics        ServiceStatisticts       `json:"statistics"`
}

type ServiceRouterDiagnostics struct {
	Queries              int `json:"queries"`
	ReplayedTransactions int `json:"replayed_transactions"`
	ROTransactions       int `json:"ro_transactions"`
	RouteAll             int `json:"route_all"`
	RouteMaster          int `json:"route_master"`
	RouteSlave           int `json:"route_slave"`
	RWTransactions       int `json:"rw_transactions"`
}

type ServiceStatisticts struct {
	ActiveOperations int `json:"active_operations"`
	Connections      int `json:"connections"`
	MaxConnections   int `json:"max_connections"`
	RoutedPackets    int `json:"routed_packets"`
	TotalConnections int `json:"total_connections"`
}

type ServiceCollector struct {
	Queries              *prometheus.Desc
	ReplayedTransactions *prometheus.Desc
	ROTransactions       *prometheus.Desc
	RouteAll             *prometheus.Desc
	RouteMaster          *prometheus.Desc
	RouteSlave           *prometheus.Desc
	RWTransactions       *prometheus.Desc
	ActiveOperations     *prometheus.Desc
	Connections          *prometheus.Desc
	MaxConnections       *prometheus.Desc
	RoutedPackets        *prometheus.Desc
	TotalConnections     *prometheus.Desc
}

func newServerCollector() *ServerCollector {
	return &ServerCollector{
		ActiveOperations:      prometheus.NewDesc("s_active_operations", "", []string{"server"}, nil),
		AdaptiveAvgSelectTime: prometheus.NewDesc("s_adaptive_avg_select_time", "", []string{"server"}, nil),
		ConnectionPoolEmpty:   prometheus.NewDesc("s_connection_pool_empty", "", []string{"server"}, nil),
		Connections:           prometheus.NewDesc("s_connections", "", []string{"server"}, nil),
		MaxConnections:        prometheus.NewDesc("s_max_connections", "", []string{"server"}, nil),
		MaxPoolSize:           prometheus.NewDesc("s_max_pool_size", "", []string{"server"}, nil),
		PersistentConnections: prometheus.NewDesc("s_persistent_connections", "", []string{"server"}, nil),
		ReusedConnections:     prometheus.NewDesc("s_reused_connections", "", []string{"server"}, nil),
		RoutedPackets:         prometheus.NewDesc("s_routed_packets", "", []string{"server"}, nil),
		TotalConnections:      prometheus.NewDesc("s_total_connections", "", []string{"server"}, nil),
	}
}

func (ServerCollector *ServerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- ServerCollector.ActiveOperations
	ch <- ServerCollector.AdaptiveAvgSelectTime
	ch <- ServerCollector.ConnectionPoolEmpty
	ch <- ServerCollector.Connections
	ch <- ServerCollector.MaxConnections
	ch <- ServerCollector.MaxPoolSize
	ch <- ServerCollector.PersistentConnections
	ch <- ServerCollector.ReusedConnections
	ch <- ServerCollector.RoutedPackets
	ch <- ServerCollector.TotalConnections
}

func (ServerCollector *ServerCollector) Collect(ch chan<- prometheus.Metric) {
	data := getServer()

	for i := 0; i < len(data.Data); i++ {
		ch <- prometheus.MustNewConstMetric(ServerCollector.ActiveOperations, prometheus.GaugeValue, float64(data.Data[i].Attributes.Statistics.ActiveOperations), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServerCollector.ConnectionPoolEmpty, prometheus.GaugeValue, float64(data.Data[i].Attributes.Statistics.ConnectionPoolEmpty), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServerCollector.Connections, prometheus.GaugeValue, float64(data.Data[i].Attributes.Statistics.Connections), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServerCollector.MaxConnections, prometheus.GaugeValue, float64(data.Data[i].Attributes.Statistics.MaxConnections), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServerCollector.MaxPoolSize, prometheus.GaugeValue, float64(data.Data[i].Attributes.Statistics.MaxPoolSize), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServerCollector.PersistentConnections, prometheus.GaugeValue, float64(data.Data[i].Attributes.Statistics.PersistentConnections), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServerCollector.ReusedConnections, prometheus.GaugeValue, float64(data.Data[i].Attributes.Statistics.ReusedConnections), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServerCollector.RoutedPackets, prometheus.GaugeValue, float64(data.Data[i].Attributes.Statistics.RoutedPackets), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServerCollector.TotalConnections, prometheus.GaugeValue, float64(data.Data[i].Attributes.Statistics.TotalConnections), data.Data[i].ID)
	}
}

func newServiceCollector() *ServiceCollector {
	return &ServiceCollector{
		Queries:              prometheus.NewDesc("r_queries", "", []string{"service"}, nil),
		ReplayedTransactions: prometheus.NewDesc("r_replayed_transactions", "", []string{"service"}, nil),
		ROTransactions:       prometheus.NewDesc("r_ro_transactions", "", []string{"service"}, nil),
		RouteAll:             prometheus.NewDesc("r_route_all", "", []string{"service"}, nil),
		RouteMaster:          prometheus.NewDesc("r_route_master", "", []string{"service"}, nil),
		RouteSlave:           prometheus.NewDesc("r_route_slave", "", []string{"service"}, nil),
		RWTransactions:       prometheus.NewDesc("r_rw_transactions", "", []string{"service"}, nil),
		ActiveOperations:     prometheus.NewDesc("r_active_operations", "", []string{"service"}, nil),
		Connections:          prometheus.NewDesc("r_connections", "", []string{"service"}, nil),
		MaxConnections:       prometheus.NewDesc("r_max_connections", "", []string{"service"}, nil),
		RoutedPackets:        prometheus.NewDesc("r_routed_packets", "", []string{"service"}, nil),
		TotalConnections:     prometheus.NewDesc("r_total_connections", "", []string{"service"}, nil),
	}
}

func (ServiceCollector *ServiceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- ServiceCollector.Queries
	ch <- ServiceCollector.ReplayedTransactions
	ch <- ServiceCollector.ROTransactions
	ch <- ServiceCollector.RouteAll
	ch <- ServiceCollector.RouteMaster
	ch <- ServiceCollector.RouteSlave
	ch <- ServiceCollector.RWTransactions
	ch <- ServiceCollector.ActiveOperations
	ch <- ServiceCollector.Connections
	ch <- ServiceCollector.MaxConnections
	ch <- ServiceCollector.RoutedPackets
	ch <- ServiceCollector.TotalConnections
}

func (ServiceCollector *ServiceCollector) Collect(ch chan<- prometheus.Metric) {
	data := getService()

	for i := 0; i < len(data.Data); i++ {
		ch <- prometheus.MustNewConstMetric(ServiceCollector.Queries, prometheus.GaugeValue, float64(data.Data[i].Attributes.RouterDiagnostics.Queries), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServiceCollector.ReplayedTransactions, prometheus.GaugeValue, float64(data.Data[i].Attributes.RouterDiagnostics.ReplayedTransactions), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServiceCollector.ROTransactions, prometheus.GaugeValue, float64(data.Data[i].Attributes.RouterDiagnostics.ROTransactions), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServiceCollector.RouteAll, prometheus.GaugeValue, float64(data.Data[i].Attributes.RouterDiagnostics.RouteAll), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServiceCollector.RouteMaster, prometheus.GaugeValue, float64(data.Data[i].Attributes.RouterDiagnostics.RouteMaster), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServiceCollector.RouteSlave, prometheus.GaugeValue, float64(data.Data[i].Attributes.RouterDiagnostics.RouteSlave), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServiceCollector.RWTransactions, prometheus.GaugeValue, float64(data.Data[i].Attributes.RouterDiagnostics.RWTransactions), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServiceCollector.ActiveOperations, prometheus.GaugeValue, float64(data.Data[i].Attributes.Statistics.ActiveOperations), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServiceCollector.Connections, prometheus.GaugeValue, float64(data.Data[i].Attributes.Statistics.Connections), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServiceCollector.MaxConnections, prometheus.GaugeValue, float64(data.Data[i].Attributes.Statistics.MaxConnections), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServiceCollector.RoutedPackets, prometheus.GaugeValue, float64(data.Data[i].Attributes.Statistics.RoutedPackets), data.Data[i].ID)
		ch <- prometheus.MustNewConstMetric(ServiceCollector.TotalConnections, prometheus.GaugeValue, float64(data.Data[i].Attributes.Statistics.TotalConnections), data.Data[i].ID)
	}
}

func getServer() Server {
	data := getHttp("servers")
	var server Server
	json.Unmarshal(data, &server)
	return server
}

func getService() Service {
	data := getHttp("services")
	var service Service
	json.Unmarshal(data, &service)
	return service
}

func getHttp(path string) []byte {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	url := configHost + ":" + strconv.Itoa(configPort) + "/v1/" + path

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(configUsername, configPassword)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("ERROR: Cannot open path ", url)
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	log.Println("[LOG]   Served", url)

	return body
}

func main() {
	var pathArg string

	flag.StringVar(&pathArg, "path", "", "Path to json configuration file")
	flag.Parse()

	jsonConfig, err := os.Open(pathArg)
	if err != nil {
		log.Fatal("[ERROR] Cannot open file", pathArg)
		return
	}

	defer jsonConfig.Close()

	configBytes, _ := ioutil.ReadAll(jsonConfig)

	var config Config

	errJson := json.Unmarshal(configBytes, &config)
	if errJson != nil {
		log.Fatal("[ERROR] Cannot parse json configuration file")
	}

	configUsername = config.Username
	configPassword = config.Password
	configHost = config.Host
	configPort = config.Port

	ServerCollector := newServerCollector()
	ServiceCollector := newServiceCollector()
	prometheus.MustRegister(ServerCollector)
	prometheus.MustRegister(ServiceCollector)

	http.Handle("/metrics", promhttp.Handler())
	log.Println("[LOG]   Serving on port 9104 endpoint /metrics")
	log.Fatal(http.ListenAndServe(":9104", nil))
}
