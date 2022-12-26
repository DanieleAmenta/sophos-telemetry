package main

import (
	"github.com/amarchese96/sophos-telemetry/metrics"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAvgAppTraffic(c *gin.Context) {
	appGroupName := c.Query("app-group")
	appName := c.Query("app")
	rangeWidth := c.Query("range-width")

	if rangeWidth == "" {
		rangeWidth = "5m"
	}

	if appName != "" {
		trafficValues := map[string]float64{}
		results, _, err := metrics.GetAvgAppTraffic(appGroupName, appName, rangeWidth)

		//fmt.Println(warnings)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
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
		trafficValues := map[string]map[string]float64{}
		results, _, err := metrics.GetAllAvgAppTraffic(appGroupName, rangeWidth)

		//fmt.Println(warnings)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
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

func getAvgNodeLatencies(c *gin.Context) {
	nodeName := c.Query("node")

	rangeWidth := c.Query("range-width")

	if rangeWidth == "" {
		rangeWidth = "5m"
	}

	if nodeName != "" {
		latencyValues := map[string]float64{}
		results, _, err := metrics.GetAvgNodeLatencies(nodeName, rangeWidth)

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
		results, _, err := metrics.GetAllAvgNodeLatencies(rangeWidth)

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
	router.GET("/metrics/app/avg-traffic", getAvgAppTraffic)
	router.GET("/metrics/node/avg-latencies", getAvgNodeLatencies)

	router.Run("0.0.0.0:8080")
}
