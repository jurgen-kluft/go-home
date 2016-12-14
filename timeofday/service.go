// SHANGHAI : 
// - Lattitude = 31.2222200
// - Longitude = 121.4580600
// http://api.sunrise-sunset.org/json?lat=31.2222200&lng=121.4580600

// Time is in UTC
// RESPONSE = 
// {
// 	"results": {
// 		"sunrise": "10:45:01 PM",
// 		"sunset": "8:53:12 AM",
// 		"solar_noon": "3:49:06 AM",
// 		"day_length": "10:08:11",
// 		"civil_twilight_begin": "10:18:24 PM",
// 		"civil_twilight_end": "9:19:48 AM",
// 		"nautical_twilight_begin": "9:48:10 PM",
// 		"nautical_twilight_end": "9:50:03 AM",
// 		"astronomical_twilight_begin": "9:18:35 PM",
// 		"astronomical_twilight_end": "10:19:38 AM"
// 	},
// 	"status": "OK"
// }

func main() {

	// Open REDIS client
	// Command-Line or Config file should specify the connection details of REDIS

	timeofday_config_key := "Go-Home-TimeOfDay-Config"
	timeofday_state_key := "Go-Home-TimeOfDay-State"

	channel_name := "Go-Home"

	//  - name  : SUN
	//    start : $(SUNRISE)
	//    end   : $(SUNSET)

	// Load config from key 'timeofday_config_key'
	// Config is a YAML configuration like this:
	//
	// --- TimeOfDay Configuration
	//   URL : http://api.sunrise-sunset.org/json?lat=31.2222200&lng=121.4580600
	//   update_every : 300
    //   timeofday:
	//       - name  : BREAKFAST
	//         start :  7:00
	//         end   :  9:30
	//       - name  : MORNING
	//         start :  6:00
	//         end   : 12:00
	//       - name  : NOON
	//         start : 12:00
	//         end   : 13:00
	//       - name  : LUNCH
	//         start : 11:45
	//         end   : 12:45
	//       - name  : AFTERNOON
	//         start : 13:00
	//         end   : 18:00
	//       - name  : DINNER
	//         start : 18:00
	//         end   : 20:00
	//       - name  : NIGHT
	//         start : 20:00
	//         end   :  6:00
	//       - name  : SLEEPING
	//         start : 22:00
	//         end   :  6:00
	//       - name  : EVENING
	//         start : ($(SUN).start - 0:30)
	//         end   : SLEEPING.start

	// 	Every N seconds 
	//     
	//     Determine the TimeOfDay elements we are in
	//     Send as YAML to REDIS channel 'channel_name', like:
	//     
	//          ---- TIMEOFDAY
	//          timeofday : Morning & Breakfast

}