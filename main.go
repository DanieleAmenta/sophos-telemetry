package main

import (
	"fmt"
	"github.com/amarchese96/sophos-telemetry/metrics"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAppTraffic(c *gin.Context) {
	appGroupName := c.Query("app-group")
	appName := c.Query("app")
	rangeWidth := c.Query("range-width")

	if rangeWidth == "" {
		rangeWidth = "5m"
	}

	if appName != "" {
		results, _, err := metrics.GetAppTraffic(appGroupName, appName, rangeWidth)

		//fmt.Println(warnings)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			trafficValues := map[string]float64{}
			for _, result := range results {
				if string(result.Metric["source_app"]) == appName {
					trafficValues[string(result.Metric["destination_app"])] = float64(result.Value)
				} else if string(result.Metric["destination_app"]) == appName {
					trafficValues[string(result.Metric["source_app"])] = float64(result.Value)
				}
			}
			c.IndentedJSON(http.StatusOK, trafficValues)
		}
	} else {
		results, _, err := metrics.GetAllAppTraffic(appGroupName, rangeWidth)

		//fmt.Println(warnings)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			trafficValues := map[string]map[string]float64{}
			for _, result := range results {
				_, ok := trafficValues[string(result.Metric["source_app"])]
				if !ok {
					trafficValues[string(result.Metric["source_app"])] = map[string]float64{}
				}
				trafficValues[string(result.Metric["source_app"])][string(result.Metric["destination_app"])] = float64(result.Value)

				_, ok = trafficValues[string(result.Metric["destination_app"])]
				if !ok {
					trafficValues[string(result.Metric["destination_app"])] = map[string]float64{}
				}
				trafficValues[string(result.Metric["destination_app"])][string(result.Metric["source_app"])] = float64(result.Value)
			}
			c.IndentedJSON(http.StatusOK, trafficValues)
		}
	}
}

func getAppCPU(c *gin.Context) {
	appGroupName := c.Query("app-group")
	appName := c.Query("app")
	rangeWidth := c.Query("range-width")

	if rangeWidth == "" {
		rangeWidth = "5m"
	}

	if appName != "" {
		results, _, err := metrics.GetAppCPU(appGroupName, appName, rangeWidth)

		//fmt.Println(warnings)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			if len(results) < 1 {
				c.IndentedJSON(http.StatusNotFound, fmt.Errorf("cpu usage metrics for app %s not found", appName))
			} else {
				c.IndentedJSON(http.StatusOK, float64(results[0].Value))
			}
		}
	} else {
		results, _, err := metrics.GetAllAppCPU(appGroupName, rangeWidth)

		//fmt.Println(warnings)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			cpuValues := map[string]float64{}
			for _, result := range results {
				cpuValues[string(result.Metric["container"])] = float64(result.Value)
			}
			c.IndentedJSON(http.StatusOK, cpuValues)
		}
	}
}

func getAppMemory(c *gin.Context) {
	appGroupName := c.Query("app-group")
	appName := c.Query("app")
	rangeWidth := c.Query("range-width")

	if rangeWidth == "" {
		rangeWidth = "5m"
	}

	if appName != "" {
		results, _, err := metrics.GetAppMemory(appGroupName, appName, rangeWidth)

		//fmt.Println(warnings)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			if len(results) < 1 {
				c.IndentedJSON(http.StatusNotFound, fmt.Errorf("memory usage metrics for app %s not found", appName))
			} else {
				c.IndentedJSON(http.StatusOK, float64(results[0].Value))
			}
		}
	} else {
		results, _, err := metrics.GetAllAppMemory(appGroupName, rangeWidth)

		//fmt.Println(warnings)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			memoryValues := map[string]float64{}
			for _, result := range results {
				memoryValues[string(result.Metric["container"])] = float64(result.Value)
			}
			c.IndentedJSON(http.StatusOK, memoryValues)
		}
	}
}

func getNodeLatencies(c *gin.Context) {
	nodeName := c.Query("node")

	rangeWidth := c.Query("range-width")

	if rangeWidth == "" {
		rangeWidth = "5m"
	}

	if nodeName != "" {
		latencyValues := map[string]float64{}
		results, _, err := metrics.GetNodeLatencies(nodeName, rangeWidth)

		//fmt.Println(warnings)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			for _, result := range results {
				latencyValues[string(result.Metric["destination_node"])] = float64(result.Value)
			}
			c.IndentedJSON(http.StatusOK, latencyValues)
		}
	} else {
		latencyValues := map[string]map[string]float64{}
		results, _, err := metrics.GetAllNodeLatencies(rangeWidth)

		//fmt.Println(warnings)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			for _, result := range results {
				_, ok := latencyValues[string(result.Metric["origin_node"])]
				if !ok {
					latencyValues[string(result.Metric["origin_node"])] = map[string]float64{}
				}
				latencyValues[string(result.Metric["origin_node"])][string(result.Metric["destination_node"])] = float64(result.Value)
			}
			c.IndentedJSON(http.StatusOK, latencyValues)
		}
	}
}

func main() {
	router := gin.Default()
	router.GET("/metrics/app/traffic", getAppTraffic)
	router.GET("/metrics/app/cpu", getAppCPU)
	router.GET("/metrics/app/memory", getAppMemory)
	router.GET("/metrics/node/latencies", getNodeLatencies)

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Printf("Exiting because of error: %s", err.Error())
		return
	}
}
