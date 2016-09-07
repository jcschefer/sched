// Jack Schefer, began 9/3/16
// Purpose: to print out the schedule using the ION API. 
//
package main
//
import (
   "os"
   "fmt"
   "time"
   "strconv"
   "strings"
   "net/http"
   "encoding/json"
   "io/ioutil"
)
//
const (
   CYAN  string = "\033[96m"
   GREEN string = "\033[92m"
   BOLD  string = "\033[1m"
   END   string = "\033[0m"
)
//
func main() {
   if len(os.Args) > 1 && os.Args[1] == "week" {
      // show the order of days in the week
      fmt.Println("\nRemainder of the week's schedule.\n")
      _, this_week := time.Now().ISOWeek()
      var i int = 1
      next_name, next_date := parseJSON(1)
      t, err := time.Parse("2006-01-02", next_date)
      check(err)
      _, week := t.ISOWeek()
      for week == this_week {
         fmt.Printf("%s: %s\n", t.Weekday(), next_name)
         i = i + 1
         next_name, next_date = parseJSON(i)
         t, err = time.Parse("2006-01-02", next_date)
         check(err)
         _, week = t.ISOWeek()
         // so that next_name isn't unused
         foo := strings.Index(next_name, "foo")
         if foo > -1 {
            fmt.Println("name contains foo!")
         }
      }
      //
} else {
      //
      data := map[string]interface{}{}
      //
      // 1. Get and print out the current date and time.
      //
      fmt.Println()
      now := time.Now()
      var tsuffix string
      var zero string = ""
      fmt.Printf("%v, %v %v\n", now.Weekday(), now.Month(), now.Day())
      if now.Hour() / 12 == 0 {
         tsuffix = "am"
      } else {
         tsuffix = "pm"
      }
      if now.Minute() < 10 {
         zero = "0"
      }
      hr := now.Hour() % 12
      if hr == 0 {
         hr = 12
      }
      fmt.Printf("Time: %v:%s%v %s\n", hr, zero, now.Minute(), tsuffix)
      fmt.Println()
      //
      // 2. Get the schedule data from Ion
      res, err := http.Get("https://ion.tjhsst.edu/api/schedule?format=json")
      check(err)
      defer res.Body.Close()
      body, err := ioutil.ReadAll(res.Body)
      check(err)
      json.Unmarshal(body, &data)
      //
      // 3. check if there's school today
      schedule_date := data["results"].([]interface{})[0].(map[string]interface{})["date"].(string)
      today := false
      if now.Format("2006-01-02") == schedule_date {
         today = true
      }
      //
      // 4. if so, print it out
      if today {
         title := data["results"].([]interface{})[0].(map[string]interface{})["day_type"].(map[string]interface{})["name"]
         fmt.Println(title)
         //
         blocks := data["results"].([]interface{})[0].(map[string]interface{})["day_type"].(map[string]interface{})["blocks"].([]interface{})
         for _, b := range blocks {
            name  := b.(map[string]interface{})["name"].(string)
            start := b.(map[string]interface{})["start"].(string)
            end   := b.(map[string]interface{})["end"].(string)
            if strings.Contains(name, "<br>") {
               name = name[0: strings.Index(name, "<br>")]
            }
            //
            shrs,serr := strconv.Atoi(start[0:strings.Index(start, ":")])
            smin,smer := strconv.Atoi(start[strings.Index(start, ":") + 1:])
            ehrs,eerr := strconv.Atoi(end[0:strings.Index(end, ":")])
            emin,emer := strconv.Atoi(end[strings.Index(end, ":") + 1:])
            if serr == nil && shrs > 12 {
               start = strings.Replace(start, strconv.Itoa(shrs), strconv.Itoa(shrs - 12), 1)
            }
            if eerr == nil && ehrs > 12 {
               end = strings.Replace(end, strconv.Itoa(ehrs), strconv.Itoa(ehrs - 12), 1)
            }
            var isCurrentBlock bool
            if serr == nil && eerr == nil && smer == nil && emer == nil {
               isCurrentBlock = withinTimes(now.Hour(), now.Minute(), shrs, smin, ehrs, emin)
            }
            //
            if isCurrentBlock {
               fmt.Printf("%s", CYAN)
            }
            fmt.Printf("%s:  \t%s\t-   %s\n", name, start, end)
            if isCurrentBlock {
               fmt.Printf("%s", END)
            }
         }
      } else {
         fmt.Println("No schedule available for today.")
      }
   }
   fmt.Println()
}
//
///////////////////////////////////////////////////////////
//
//  HELPER METHODS
//
///////////////////////////////////////////////////////////
//
func parseJSON(index int) (string, string) {
   data := map[string]interface{}{}
   res, err := http.Get("https://ion.tjhsst.edu/api/schedule?format=json&page=" + strconv.Itoa(index))
   check(err)
   defer res.Body.Close()
   body, err := ioutil.ReadAll(res.Body)
   check(err)
   json.Unmarshal(body, &data)
   //
   name := data["results"].([]interface{})[0].(map[string]interface{})["day_type"].(map[string]interface{})["name"].(string)
   date := data["results"].([]interface{})[0].(map[string]interface{})["date"].(string)
   return name, date
}
//
///////////////////////////////////////////////////////////
//
func withinTimes(nhrs int, nmin int, shrs int, smin int, ehrs int, emin int) bool {
   return shrs * 60 + smin <= nhrs * 60 + nmin && nhrs * 60 + nmin <= ehrs * 60 + emin
}
//
///////////////////////////////////////////////////////////
//
func check(e error) {
   if e != nil {
      panic(e.Error())
   }
}
//
// End of file.
