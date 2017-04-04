package apis

import "fmt"
import "flag"
import "time"
import "github.com/adbjesus/apis/parser"
import "github.com/adbjesus/apis/database"

var agencyFlag string
var routeFlag string
var whatToGetFlag int

func init() {
	flag.StringVar(&agencyFlag, "agency", "mbta", "Nextbus API: Agency Tag")
	flag.StringVar(&routeFlag, "route", "1", "Nextbus API: Route Tag")
	flag.IntVar(&whatToGetFlag, "obtain", 0, "What to obtain? 0 for vehicle Locations, 1 for Stops, 2 for predictions")
}

type VehicleLocation struct {
	Id          string  `xml:"id,attr"`
	RouteTag    string  `xml:"routeTag,attr"`
	DirTag      string  `xml:"dirTag,attr"`
	Lat         float64 `xml:"lat,attr"`
	Lon         float64 `xml:"lon,attr"`
	SinceReport int     `xml:"secsSinceReport,attr"`
	Predictable bool    `xml:"predictable,attr"`
	Heading     int     `xml:"heading,attr"`
	Speed       int     `xml:"speedKmHr,attr"`
}

type LocationQuery struct {
	VehicleLocation []VehicleLocation `xml:"vehicle"`
}

type Stop struct {
	Tag        string  `xml:"tag,attr"`
	Title      string  `xml:"title,attr"`
	ShortTitle string  `xml:"shortTitle,attr"`
	Lat        float64 `xml:"lat,attr"`
	Lon        float64 `xml:"lon,attr"`
	StopId     string  `xml:"stopId,attr"`
}

type StopsQuery struct {
	Stops []Stop `xml:"route>stop"`
}

type Prediction struct {
	Time        int    `xml:"epochTime,attr"`
	Seconds     int    `xml:"seconds,attr"`
	IsDeparture bool   `xml:"isDeparture,attr"`
	Vehicle     string `xml:"vehicle,attr"`
}

type PredictionDirection struct {
	Direction  string       `xml:"title,attr"`
	Prediction []Prediction `xml:"prediction"`
}

type PredictionStop struct {
	StopTag             string                `xml:"stopTag,attr"`
	PredictionDirection []PredictionDirection `xml:"direction"`
}

type PredictionList struct {
	Predictions []PredictionStop `xml:"predictions"`
}

var db database.DB

func Nextbus() {
	flag.Parse()

	var stops StopsQuery

	err := database.Connect(&db)
	if err != nil {
		fmt.Print(time.Now().Format("2006/01/02 15:04:05 : "))
		fmt.Println(err)
		return
	}

	switch whatToGetFlag {
	case 0:
		getLocation(true)
		break
	case 1:
		getStops(true)
		break
	case 2:
		stops = getStops(false)
		getPredictions(stops, true)
		break
	}

}

func getLocation(saveToDb bool) LocationQuery {
	var query LocationQuery
	var querypage = fmt.Sprintf("http://webservices.nextbus.com/service/publicXMLFeed?command=vehicleLocations&a=%s&r=%s&t=0", agencyFlag, routeFlag)
	var err error

	err = parser.ParseWebXML(querypage, &query)
	if err != nil {
		fmt.Print(time.Now().Format("2006/01/02 15:04:05 : "))
		fmt.Println(err)
		return query
	}
	t := time.Now()
	t = t.UTC()
	if saveToDb {
		for _, v := range query.VehicleLocation {
			var arr []string
			arr = append(arr, t.Format("2006/01/02 15:04:05 UTC"), agencyFlag, routeFlag, v.Id, v.DirTag, fmt.Sprintf("%f", v.Lat), fmt.Sprintf("%f", v.Lon), fmt.Sprintf("%t", v.Predictable))
			database.InsertArray(&db, "vehicle_locations", arr)
		}
	}
	return query
}

func getStops(saveToDb bool) StopsQuery {
	var query StopsQuery
	var querypage = fmt.Sprintf("http://webservices.nextbus.com/service/publicXMLFeed?command=routeConfig&a=%s&r=%s", agencyFlag, routeFlag)
	var err error

	err = parser.ParseWebXML(querypage, &query)
	if err != nil {
		fmt.Print(time.Now().Format("2006/01/02 15:04:05 : "))
		fmt.Println(err)
		return query
	}

	t := time.Now()
	t = t.UTC()

	if saveToDb {
		for _, v := range query.Stops {
			var arr []string
			arr = append(arr, t.Format("2006/01/02 15:04:05 UTC"), agencyFlag, routeFlag, v.Tag, v.Title, v.ShortTitle, fmt.Sprintf("%f", v.Lat), fmt.Sprintf("%f", v.Lon), v.StopId)
			database.InsertArray(&db, "stops", arr)
		}
	}
	return query
}

func getPredictions(stops StopsQuery, saveToDb bool) PredictionList {
	var query PredictionList
	var querypage string
	var err error

	querypage = fmt.Sprintf("http://webservices.nextbus.com/service/publicXMLFeed?command=predictionsForMultiStops&a=%s", agencyFlag)
	for _, s := range stops.Stops {
		querypage += fmt.Sprintf("&stops=%s|%s", routeFlag, s.Tag)
	}

	err = parser.ParseWebXML(querypage, &query)
	if err != nil {
		fmt.Print(time.Now().Format("2006/01/02 15:04:05 : "))
		fmt.Println(err)
		return query
	}

	t := time.Now()
	t = t.UTC()

	if saveToDb {
		var arr []string
		arr = append(arr, t.Format("2006/01/02 15:04:05 UTC"), agencyFlag, routeFlag, "bus_id", "stop", "direction", "seconds")
		for _, i := range query.Predictions {
			arr[4] = i.StopTag
			for _, j := range i.PredictionDirection {
				arr[5] = j.Direction
				for _, k := range j.Prediction {
					arr[3] = k.Vehicle
					arr[6] = fmt.Sprintf("%d", k.Seconds)
					database.InsertArray(&db, "predictions", arr)
				}
			}
		}
	}
	return query
}
