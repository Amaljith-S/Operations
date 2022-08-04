package main

import (
	deploymentusage "costmap/deployment"
	kindusage "costmap/kind"
	namespaceusage "costmap/namespaces"
	nodeUsage "costmap/nodeusage"
	"fmt"
	"github.com/rivo/tview"
	"strconv"
	"time"
)

func main() {
	app := tview.NewApplication()
	form := tview.NewForm().
		AddDropDown("Start Search from The Year \n", []string{"2017", "2018", "2019", "2020", "2021", "2022"}, 0, nil).
		AddDropDown("Search from The Month", []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}, 0, nil).
		AddDropDown("Search from The Day of the Month", []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31"}, 0, nil).
		AddDropDown("Search From The Add Hour", []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24"}, 0, nil).
		AddDropDown("Search Till The Year", []string{"2017", "2018", "2019", "2020", "2021", "2022"}, 0, nil).
		AddDropDown(" Search Till The Month", []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}, 0, nil).
		AddDropDown("Search Till The Day of the Month", []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31"}, 0, nil).
		AddDropDown("Search Till The Add Hour", []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24"}, 0, nil).
		AddDropDown("Search By Category \n", []string{"Namespace", "Kind", "Deployment","Node"}, 0, nil).
		AddButton("Search Now", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Enter Date to Cost Search").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).SetFocus(form).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
	_, a := form.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
	_, b := form.GetFormItem(1).(*tview.DropDown).GetCurrentOption()
	var k, l string
	switch {
	case b == "January":
		k = "01"
	case b == "February":
		k = "02"
	case b == "March":
		k = "03"
	case b == "April":
		k = "04"
	case b == "May":
		k = "05"
	case b == "June":
		k = "06"
	case b == "July":
		k = "07"
	case b == "August":
		k = "08"
	case b == "September":
		k = "09"
	case b == "October":
		k = "10"
	case b == "November":
		k = "11"
	case b == "December":
		k = "12"
	}

	_, c := form.GetFormItem(2).(*tview.DropDown).GetCurrentOption()

	_, d := form.GetFormItem(3).(*tview.DropDown).GetCurrentOption()

	_, e := form.GetFormItem(4).(*tview.DropDown).GetCurrentOption()

	_, f := form.GetFormItem(5).(*tview.DropDown).GetCurrentOption()

	switch {
	case f == "January":
		l = "01"
	case f == "February":
		l = "02"
	case f == "March":
		l = "03"
	case f == "April":
		l = "04"
	case f == "May":
		l = "05"
	case f == "June":
		l = "06"
	case f == "July":
		l = "07"
	case f == "August":
		l = "08"
	case b == "September":
		l = "09"
	case f == "October":
		l = "10"
	case f == "November":
		l = "11"
	case f == "December":
		l = "12"
	}

	_, g := form.GetFormItem(6).(*tview.DropDown).GetCurrentOption()
	_, h := form.GetFormItem(7).(*tview.DropDown).GetCurrentOption()
	_, i := form.GetFormItem(8).(*tview.DropDown).GetCurrentOption()

	timeFormatString := "-"
	timeCorrecting := ":00:00"

	fmt.Println(a+timeFormatString+k+timeFormatString+c+"T"+d+timeCorrecting, "to", e+timeFormatString+l+timeFormatString+g+"T"+h+timeCorrecting)
	fmt.Println("Calculating the cost of the Type :", i)

	fromDate := a + timeFormatString + k + timeFormatString + c + "T" + d + timeCorrecting
	toDate := e + timeFormatString + l + timeFormatString + g + "T" + h + timeCorrecting

	aa, _ := strconv.Atoi(a)
	kk, _ := strconv.Atoi(k)
	cc, _ := strconv.Atoi(c)
	dd, _ := strconv.Atoi(d)

	ee, _ := strconv.Atoi(e)
	ll, _ := strconv.Atoi(l)
	gg, _ := strconv.Atoi(g)
	hh, _ := strconv.Atoi(h)

	firstDate := time.Date(aa, time.Month(kk), cc, dd, 0, 0, 0, time.UTC)
	secondDate := time.Date(ee, time.Month(ll), gg, hh, 0, 0, 0, time.UTC)
	difference := secondDate.Sub(firstDate)
	timeinHours := difference.Hours()


	switch {
	case i == "Namespace":
		namespaceusage.NamespaceCostFinder(fromDate, toDate, timeinHours)
	case i == "Kind":
		kindusage.KindUsageFinder(fromDate, toDate, timeinHours)
	case i == "Deployment":
		deploymentusage.DeploymentUsageFinder(fromDate, toDate, timeinHours)
	case i== "Node":
		nodeUsage.NodeCostFinder(fromDate, toDate, timeinHours)
		
	}

}
