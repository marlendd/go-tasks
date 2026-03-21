package main

import (
	"encoding/csv"
	"encoding/xml"
	"io"
	"os"
	"strconv"
	"strings"
)

// начало решения
type Employee struct {
	Id int `xml:"id,attr"`
	Name string `xml:"name"`
	City string `xml:"city"`
	Salary int `xml:"salary"`
	
}
type Department struct {
	Code string `xml:"code"`
	Employees []Employee `xml:"employees>employee"`
}

type Organization struct {
	Departments []Department `xml:"department"`
}

func (e *Employee) Slice(depCode string) []string {
	return []string{strconv.Itoa(e.Id), e.Name, e.City, depCode, strconv.Itoa(e.Salary)}
}
// ConvertEmployees преобразует XML-документ с информацией об организации
// в плоский CSV-документ с информацией о сотрудниках
func ConvertEmployees(outCSV io.Writer, inXML io.Reader) error {
	var org Organization
	dec := xml.NewDecoder(inXML)
	if err := dec.Decode(&org); err != nil {
		return err
	}
	w := csv.NewWriter(outCSV)
	w.Write([]string{"id", "name" ,"city", "department" ,"salary"})
	for _, d := range org.Departments {
		for _, e := range d.Employees {
			err := w.Write(e.Slice(d.Code)) 
			if err != nil {
				return err
			}
		}
	}
	w.Flush()
	if w.Error() != nil {
		return w.Error()
	}
	return nil
}

// конец решения


func main() {
	src := `<organization>
    <department>
        <code>hr</code>
        <employees>
            <employee id="11">
                <name>Дарья</name>
                <city>Самара</city>
                <salary>70</salary>
            </employee>
            <employee id="12">
                <name>Борис</name>
                <city>Самара</city>
                <salary>78</salary>
            </employee>
        </employees>
    </department>
    <department>
        <code>it</code>
        <employees>
            <employee id="21">
                <name>Елена</name>
                <city>Самара</city>
                <salary>84</salary>
            </employee>
        </employees>
    </department>
</organization>`

	in := strings.NewReader(src)
	out := os.Stdout
	ConvertEmployees(out, in)
	/*
		id,name,city,department,salary
		11,Дарья,Самара,hr,70
		12,Борис,Самара,hr,78
		21,Елена,Самара,it,84
	*/
}
