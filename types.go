package main

import (
	"github.com/buzzer13/healthpose/checks/dns"
)

type CheckCassandra struct {
	Hosts    []string `koanf:"hosts"`
	Keyspace string   `koanf:"keyspace"`
}

type CheckDNS struct {
	Address        string         `koanf:"address"`
	Server         string         `koanf:"server"`
	Type           dns.RecordType `koanf:"type"`
	RequestTimeout float64        `koanf:"request_timeout"`
	FallbackDelay  float64        `koanf:"fallback_delay"`
}

//type CheckGRPC struct {
//	Target       string        `koanf:"target"`
//	Service      string        `koanf:"service"`
//	DialOptions  map[string]interface{}
//	CheckTimeout float64 `koanf:"check_timeout"`
//}

type CheckHTTP struct {
	URL            string  `koanf:"url"`
	RequestTimeout float64 `koanf:"request_timeout"`
}

type CheckICMP struct {
	Address        string  `koanf:"address"`
	Count          int     `koanf:"count"`
	Interval       float64 `koanf:"interval"`
	RequestTimeout float64 `koanf:"request_timeout"`
}

type CheckInfluxDB struct {
	URL string `koanf:"url"`
}

type CheckMemcached struct {
	DSN string `koanf:"dsn"`
}

type CheckMongo struct {
	DSN               string  `koanf:"dsn"`
	TimeoutConnect    float64 `koanf:"timeout_connect"`
	TimeoutDisconnect float64 `koanf:"timeout_disconnect"`
	TimeoutPing       float64 `koanf:"timeout_ping"`
}

type CheckMySQL struct {
	DSN string `koanf:"dsn"`
}

type CheckNATS struct {
	DSN string `koanf:"dsn"`
}

type CheckPostgres struct {
	DSN string `koanf:"dsn"`
}

type CheckRabbitMQ struct {
	DSN            string  `koanf:"dsn"`
	Exchange       string  `koanf:"exchange"`
	RoutingKey     string  `koanf:"routing_key"`
	Queue          string  `koanf:"queue"`
	ConsumeTimeout float64 `koanf:"consume_timeout"`
	DialTimeout    float64 `koanf:"dial_timeout"`
}

type CheckRedis struct {
	DSN string `koanf:"dsn"`
}

type ServiceCheck struct {
	Name      string          `koanf:"name"`
	Timeout   float64         `koanf:"timeout"`
	SkipOnErr bool            `koanf:"optional"`
	Cassandra *CheckCassandra `koanf:"cassandra,omitempty"`
	DNS       *CheckDNS       `koanf:"dns,omitempty"`
	//GRPC      *CheckGRPC      `koanf:"grpc,omitempty"`
	HTTP      *CheckHTTP      `koanf:"http,omitempty"`
	ICMP      *CheckICMP      `koanf:"icmp,omitempty"`
	InfluxDB  *CheckInfluxDB  `koanf:"influxdb,omitempty"`
	Memcached *CheckMemcached `koanf:"memcached,omitempty"`
	Mongo     *CheckMongo     `koanf:"mongo,omitempty"`
	MySQL     *CheckMySQL     `koanf:"mysql,omitempty"`
	NATS      *CheckNATS      `koanf:"nats,omitempty"`
	Postgres  *CheckPostgres  `koanf:"postgres,omitempty"`
	RabbitMQ  *CheckRabbitMQ  `koanf:"rabbitmq,omitempty"`
	Redis     *CheckRedis     `koanf:"redis,omitempty"`
}

type ConfigService struct {
	Name    string         `koanf:"name"`
	Version string         `koanf:"version"`
	Checks  []ServiceCheck `koanf:"checks"`
}

type ConfigHTTP struct {
	Listen string `koanf:"listen"`
}

type Config struct {
	HTTP     ConfigHTTP               `koanf:"http"`
	Services map[string]ConfigService `koanf:"services,omitempty"`
}
