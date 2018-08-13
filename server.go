package main

import (
	"flag"
	"sample-room-golang/routers"
	// "sampleroomgolangnew/plugins"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"crypto/tls"
	"fmt"
	"os"
)

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	return ":" + port
}

// func HystrixHandler(command string) gin.HandlerFunc {
//   return func(c *gin.Context) {
//     hystrix.Do(command, func() error {
//       c.Next()
//       return nil
//     }, func(err error) error {
//       c.String(http.StatusInternalServerError, "500 Internal Server Error")
//       return err
//     })
//   }
// }

func RequestTracker(counter *prometheus.CounterVec) gin.HandlerFunc {
	return func(c *gin.Context) {
		labels := map[string]string{"Route": c.Request.URL.Path, "Method": c.Request.Method}
		counter.With(labels).Inc()
		c.Next()
	}
}

// func OpenTracing() gin.HandlerFunc {
//   return func(c *gin.Context) {
//     wireCtx, _ := opentracing.GlobalTracer().Extract(
//       opentracing.HTTPHeaders,
//       opentracing.HTTPHeadersCarrier(c.Request.Header))

//     serverSpan := opentracing.StartSpan(c.Request.URL.Path,
//       ext.RPCServerOption(wireCtx))
//     defer serverSpan.Finish()
//     c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), serverSpan))
//     c.Next()
//   }
// }

// type LogrusAdapter struct{}

// func (l LogrusAdapter) Error(msg string) {
//   log.Errorf(msg)
// }

// func (l LogrusAdapter) Infof(msg string, args ...interface{}) {
//   log.Infof(msg, args)
// }

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// Adding Route Counter via Prometheus Metrics
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "counters",
		Subsystem: "page_requests",
		Name:      "request_count",
		Help:      "Number of requests received",
	}, []string{"Route", "Method"})
	prometheus.MustRegister(counter)

	// Hystrix configuration
	// hystrix.ConfigureCommand("timeout", hystrix.CommandConfig{
	//   Timeout: 1000,
	//   MaxConcurrentRequests: 100,
	//   ErrorPercentThreshold: 25,
	// })
	//Add Hystrix to prometheus metrics
	// collector := plugins.InitializePrometheusCollector(plugins.PrometheusCollectorConfig{
	//   Namespace: "sampleroomgolangnew",
	// })
	// metricCollector.Registry.Register(collector.NewPrometheusCollector)

	//And jaeger metrics and reporting to prometheus route
	// logAdapt := LogrusAdapter{}
	// jaeger.NewLoggingReporter(logAdapt)
	// factory := jaegerprom.New()
	// metrics := jaeger.NewMetrics(factory, map[string]string{"lib": "jaeger"})

	// transport, err := jaeger.NewUDPTransport("localhost:5775", 0)
	// if err != nil {
	//   log.Errorln(err.Error())
	// }

	// reporter := jaeger.NewCompositeReporter(
	//   jaeger.NewLoggingReporter(logAdapt),
	//   jaeger.NewRemoteReporter(transport,
	//     jaeger.ReporterOptions.Metrics(metrics),
	//     jaeger.ReporterOptions.Logger(logAdapt),
	//   ),
	// )
	// defer reporter.Close()

	// sampler := jaeger.NewConstSampler(true)
	// tracer, closer := jaeger.NewTracer("sampleroomgolangnew",
	//   sampler,
	//   reporter,
	//   jaeger.TracerOptions.Metrics(metrics),
	// )
	// defer closer.Close()

	// opentracing.SetGlobalTracer(tracer)

	router := gin.Default()

	router.Use(RequestTracker(counter))
	// router.Use(OpenTracing())
	// router.Use(HystrixHandler("timeout"))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Use(static.Serve("/", static.LocalFile("./public", false)))
	router.GET("/health", routers.HealthGET)

	locus := "MAIN"
	checkpoint(locus, "processCommandLine")
	err := processCommandline()
	if err != nil {
		log.Errorln(err.Error())
		flag.Usage()
		return
	}
	printConfig(&config)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	if len(config.roomToDelete) > 0 {
		checkpoint(locus, fmt.Sprintf("deleteWithRetries %s", config.roomToDelete))
		err = deleteWithRetries(client, config.roomToDelete)
		if err != nil {
			checkpoint(locus, fmt.Sprintf("DELETE.FAILED err=%s", err.Error()))
		}
		return
	}

	checkpoint(locus, "registerWithRetries")
	err = registerWithRetries(client)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	checkpoint(locus, "startServer")

	locus = "WS.SERVER"
	go TrackPlayers()
	go InjectConversations()
	checkpoint(locus, fmt.Sprintf("Listening to port %d", config.listeningPort))
	router.GET("/ws", func(c *gin.Context) {
		log.Println("Got something...")
		roomHandler(c.Writer, c.Request)
	})
	router.Run(port())
}

// Prints a simple checkpoint message.
func checkpoint(locus, s string) {
	log.Printf("CHECKPOINT: %s.%s\n", locus, s)
}
