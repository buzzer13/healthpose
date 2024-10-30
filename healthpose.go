package main

import (
	"fmt"
	"github.com/hellofresh/health-go/v5"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"log"
	"net/http"
	"os"
	"time"

	checkCassandra "github.com/hellofresh/health-go/v5/checks/cassandra"
	//checkGRPC "github.com/hellofresh/health-go/v5/checks/grpc"
	checkDNS "github.com/buzzer13/healthpose/checks/dns"
	checkICMP "github.com/buzzer13/healthpose/checks/icmp"
	checkHTTP "github.com/hellofresh/health-go/v5/checks/http"
	checkInfluxDB "github.com/hellofresh/health-go/v5/checks/influxdb"
	checkMemcached "github.com/hellofresh/health-go/v5/checks/memcached"
	checkMongo "github.com/hellofresh/health-go/v5/checks/mongo"
	checkMySQL "github.com/hellofresh/health-go/v5/checks/mysql"
	checkNATS "github.com/hellofresh/health-go/v5/checks/nats"
	checkPostgres "github.com/hellofresh/health-go/v5/checks/postgres"
	checkRabbitMQ "github.com/hellofresh/health-go/v5/checks/rabbitmq"
	checkRedis "github.com/hellofresh/health-go/v5/checks/redis"
)

var defaultConfig = Config{
	HTTP: ConfigHTTP{
		Listen: ":8080",
	},
}

func main() {
	k := koanf.NewWithConf(koanf.Conf{Delim: "."})

	_ = k.Load(structs.Provider(defaultConfig, "koanf"), nil)
	_ = k.Load(file.Provider("/etc/healthpose/healthpose.yaml"), yaml.Parser())
	_ = k.Load(file.Provider("/config/healthpose.yaml"), yaml.Parser())
	_ = k.Load(file.Provider("healthpose.yaml"), yaml.Parser())
	_ = k.Load(file.Provider(os.Getenv("CONFIG_FILE")), yaml.Parser())
	_ = k.Load(env.Provider("", "__", nil), nil)

	cfg := Config{}
	err := k.Unmarshal("", &cfg)

	if err != nil {
		log.Fatalln("failed to parse config:", err)
	}

	m := http.NewServeMux()

	m.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "OK")

		if err != nil {
			log.Println("failed to write healthcheck response")
		}
	})

	if cfg.Services == nil {
		log.Fatalln("no health checks configured")
	}

	for svcKey, svc := range cfg.Services {
		h, err := health.New(
			health.WithComponent(health.Component{
				Name:    svc.Name,
				Version: svc.Version,
			}),
		)

		if err != nil {
			log.Fatalf("failed to set up service %s: %s\n", svc.Name, err)
		}

		log.Printf("setting up %d health checks(s) for service \"%s\"\n", len(svc.Checks), svc.Name)

		for _, chk := range svc.Checks {
			var timeout = secondsToDuration(chk.Timeout)
			var checker health.CheckFunc

			switch {
			case chk.Cassandra != nil:
				checker = checkCassandra.New(checkCassandra.Config{
					Hosts:    chk.Cassandra.Hosts,
					Keyspace: chk.Cassandra.Keyspace,
				})
			case chk.DNS != nil:
				checker = checkDNS.New(checkDNS.Config{
					Address:        chk.DNS.Address,
					Server:         chk.DNS.Server,
					Type:           chk.DNS.Type,
					RequestTimeout: secondsToDuration(chk.DNS.RequestTimeout),
					FallbackDelay:  secondsToDuration(chk.DNS.FallbackDelay),
				})
			//case chk.GRPC != nil:
			//	checker = checkGRPC.New(checkGRPC.Config{})
			case chk.HTTP != nil:
				checker = checkHTTP.New(checkHTTP.Config{
					URL:            chk.HTTP.URL,
					RequestTimeout: secondsToDuration(chk.HTTP.RequestTimeout),
				})
			case chk.ICMP != nil:
				checker = checkICMP.New(checkICMP.Config{
					Address:        chk.ICMP.Address,
					Count:          chk.ICMP.Count,
					Interval:       secondsToDuration(chk.ICMP.Interval),
					RequestTimeout: secondsToDuration(chk.ICMP.RequestTimeout),
				})
			case chk.InfluxDB != nil:
				checker = checkInfluxDB.New(checkInfluxDB.Config(*chk.InfluxDB))
			case chk.Memcached != nil:
				checker = checkMemcached.New(checkMemcached.Config(*chk.Memcached))
			case chk.Mongo != nil:
				checker = checkMongo.New(checkMongo.Config{
					DSN:               chk.Mongo.DSN,
					TimeoutConnect:    secondsToDuration(chk.Mongo.TimeoutConnect),
					TimeoutDisconnect: secondsToDuration(chk.Mongo.TimeoutDisconnect),
					TimeoutPing:       secondsToDuration(chk.Mongo.TimeoutPing),
				})
			case chk.MySQL != nil:
				checker = checkMySQL.New(checkMySQL.Config(*chk.MySQL))
			case chk.NATS != nil:
				checker = checkNATS.New(checkNATS.Config(*chk.NATS))
			case chk.Postgres != nil:
				checker = checkPostgres.New(checkPostgres.Config(*chk.Postgres))
			case chk.RabbitMQ != nil:
				checker = checkRabbitMQ.New(checkRabbitMQ.Config{
					DSN:            chk.RabbitMQ.DSN,
					Exchange:       chk.RabbitMQ.Exchange,
					RoutingKey:     chk.RabbitMQ.RoutingKey,
					Queue:          chk.RabbitMQ.Queue,
					ConsumeTimeout: secondsToDuration(chk.RabbitMQ.ConsumeTimeout),
					DialTimeout:    secondsToDuration(chk.RabbitMQ.DialTimeout),
				})
			case chk.Redis != nil:
				checker = checkRedis.New(checkRedis.Config(*chk.Redis))
			default:
				log.Fatalf("invalid health check \"%s\" config for service %s\n", chk.Name, svc.Name)
			}

			if timeout <= 0 {
				timeout = 60 * time.Second
			}

			err = h.Register(health.Config{
				Name:      chk.Name,
				Timeout:   timeout,
				SkipOnErr: chk.SkipOnErr,
				Check:     checker,
			})

			if err != nil {
				log.Fatalf("failed to set up health check \"%s\" for service \"%s\": %s\n", chk.Name, svc.Name, err)
			}
		}

		m.Handle(fmt.Sprintf("/v1/health/%s", svcKey), h.Handler())
	}

	log.Println("listening on", cfg.HTTP.Listen)

	err = http.ListenAndServe(cfg.HTTP.Listen, m)

	if err != nil {
		log.Fatalln("failed to start server:", err)
	}
}

func secondsToDuration(seconds float64) time.Duration {
	return time.Duration(seconds * float64(time.Second))
}
