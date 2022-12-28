package metrics

import (
	"context"
	"fmt"
	prometheusapi "github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"os"
	"time"
)

var prometheusAddress = os.Getenv("PROMETHEUS_ADDRESS")

func newPrometheusClient(serverAddress string) (prometheus.API, error) {
	client, err := prometheusapi.NewClient(prometheusapi.Config{
		Address: serverAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating metrics client: %v", err)
	}
	return prometheus.NewAPI(client), nil
}

func GetAppTraffic(appGroupName, appName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		sum(
			rate(istio_request_bytes_sum{app_group="`+appGroupName+`", app="`+appName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])
			+
			rate(istio_response_bytes_sum{app_group="`+appGroupName+`", app="`+appName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])
		) by (source_app, destination_app)
		or 
		sum(
			rate(istio_tcp_sent_bytes_total{app_group="`+appGroupName+`", app="`+appName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`]) 
			+ 
			rate(istio_tcp_received_bytes_total{app_group="`+appGroupName+`", app="`+appName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])
		) by (source_app, destination_app)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAllAppTraffic(appGroupName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		sum(
			rate(istio_request_bytes_sum{reporter="source", app_group="`+appGroupName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])
			+
			rate(istio_response_bytes_sum{reporter="source", app_group="`+appGroupName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])
		) by (source_app, destination_app)
		or 
		sum(
			rate(istio_tcp_sent_bytes_total{reporter="source", app_group="`+appGroupName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`]) 
			+ 
			rate(istio_tcp_received_bytes_total{reporter="source", app_group="`+appGroupName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])
		) by (source_app, destination_app)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAppCPU(appGroupName, appName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	containerName := appGroupName + "-" + appName
	result, warnings, err := prometheusClient.Query(ctx, `
		avg by(container) (rate(container_cpu_usage_seconds_total{container="`+containerName+`"}[`+rangeWidth+`])) * 1000
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAllAppCPU(appGroupName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		avg by(container) (rate(container_cpu_usage_seconds_total{container=~"`+appGroupName+`-.*"}[`+rangeWidth+`])) * 1000
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAppMemory(appGroupName, appName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	containerName := appGroupName + "-" + appName
	result, warnings, err := prometheusClient.Query(ctx, `
		avg by(container) (avg_over_time(container_memory_working_set_bytes{container="`+containerName+`"}[`+rangeWidth+`]) / (1024 * 1024))
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAllAppMemory(appGroupName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		avg by(container) (avg_over_time(container_memory_working_set_bytes{container=~"`+appGroupName+`-.*"}[`+rangeWidth+`]) / (1024 * 1024))
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetNodeLatencies(nodeName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		(rate(node_latency_sum{origin_node="`+nodeName+`"}[`+rangeWidth+`]) / rate(node_latency_count{origin_node="`+nodeName+`"}[`+rangeWidth+`])) * 1000
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAllNodeLatencies(rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		(rate(node_latency_sum[`+rangeWidth+`]) / rate(node_latency_count[`+rangeWidth+`])) * 1000
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}
