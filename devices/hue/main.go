package main

/*
Groups of Lights:

* Bedroom
* Livingroom
* Kitchen
* SophiaRoom
* JenniferRoom

Scenes:
We pre-program scenes and we reference them when we modify the group.

Season   : TimeOfDay   : Group     = Scene
"WINTER" : "BREAKFAST" : "Kitchen" = "MorningWinter"
"SPRING" : "BREAKFAST" : "Kitchen" = "MorningSpring"
"WINTER" : "LUNCH"     : "Kitchen" = "NoonWinter"
"WINTER" : "DINNER"    : "Kitchen" = "SunSetWinter"
"WINTER" : "EVENING"   : "Kitchen" = "EveningWinter"
"WINTER" : "EVENING"   : "Bedroom" = "LateEveningWinter"

All these configurations are stored in REDIS, the key being
the TimeOfDay. When receiving the value we need to get the
matching Season/Group to find out the scene to apply.

We can also scope multiple fields into the key:
     KEY                     |      VALUE
Kitchen:Breakfast:Winter     |  MorningWinter
Kitchen:Breakfast:Summer     |  MorningSummer
Bedroom:Evening:Winter       |  LateEveningWinter

JSON
// KEY = TimeOfDay-BREAKFAST
{
    "scenes" : [
        { "season" : "Winter", "group" : "Kitchen", "scene" : "MorningWinter" },
        { "season" : "Spring", "group" : "Kitchen", "scene" : "MorningSpring" },
    ]
}

Then we need to create the HUE light configurations for the above scenes:
- MorningWinter
- MorningSpring
- NoonWinter
- SunSetWinter
- EveningWinter
- LateEveningWinter
- MorningSummer
- NoonSummer
-

Just need to find a way to create and easily modify these configurations.

*/

func main() {

}
