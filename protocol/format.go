//don't get it complication, but I can't do that
//btw, add error handle

package protocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func toInt(v any) int { // I have no fucking idea
	switch v := v.(type) {
	case int:
		return int(v)
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	default:
		return 0
	}
}

func Encode(a ...any) (string, error) {
	var result string = "" //result string
	var lena int = len(a)  //len of a slice

	for i, v := range a {
		switch v := v.(type) {
		case string:
			vs := strings.ReplaceAll(v, "|", "\\|")
			vs = strings.ReplaceAll(vs, "\"", "\\\"")

			if i != lena-1 {
				result += "\"" + vs + "\"" + "|"
				continue
			}

			result += "\"" + vs + "\""

		case int, int8, int16, int32, uint, uint8, uint16, uint32:
			vint := strconv.Itoa(toInt(v))

			if i != lena-1 {
				result += vint + "|"
				continue
			}
			result += vint

		case float32, float64:
			switch v := v.(type) {
			case float32:
				vfloat := strconv.FormatFloat(float64(v), 'f', -1, 64)

				if i != lena-1 {
					result += vfloat + "|"
					continue
				}
				result += vfloat

			case float64:
				vfloat := strconv.FormatFloat(v, 'f', -1, 64)

				if i != lena-1 {
					result += vfloat + "|"
					continue
				}
				result += vfloat

			}

		case bool:
			vbool := strconv.FormatBool(v)

			if i != lena-1 {
				result += vbool + "|"
				continue
			}
			result += vbool

		default:
			return "", errors.New("unsupported type")
		}
	}

	return result, nil
}

func Decode(s string) (a []any, e error) {
	//No nil string
	if len(s) <= 0 {
		return []any{}, errors.New("nil string")
	}

	runes := []rune(s) //golang default use utf-8 in []string
	lens := len(runes) //runes's len
	i, head := 0, 0    //iter and head
	strMode := false   //to know this " is first or last

	//spilt
	for i < lens {
		switch {
		case i == lens-1 && runes[i] != '"': //the last thing
			if (i-head <= 0 && len(runes) < 1) || runes[i] == '|' { //avoid nil things
				return []any{}, errors.New("nil type")
			}

			part := runes[head:]
			a = append(a, part)

		case runes[i] == '|':
			if i-1 < 0 || i+1 >= lens {
				i++
				continue //although it's by manual
			}

			if runes[i-1] == '"' {
				head = i + 1
				i++
				continue //don't append string twice
			}

			if (runes[i-1] == '"' && runes[i+1] == '"') || runes[i-1] == '\\' {
				i++
				continue //igonre it
			}

			//normal
			if i-head <= 0 { //avoid nil things
				return []any{}, errors.New("nil type")
			}

			part := runes[head:i]
			head = i + 1
			a = append(a, part)

		case runes[i] == '"': //string
			if !(i-1 < 0) { //although it's by manual
				if runes[i-1] == '\\' {
					i++
					continue //igonre it, too
				}
			}

			if strMode { //end of string
				if i-head <= 0 { //avoid nil things
					return []any{}, errors.New("nil type")
				}

				part := runes[head : i+1]
				a = append(a, part)
				strMode = false
			} else { //start of string
				head = i
				strMode = true
			}
		}

		i++
	}

	// it's an error
	if strMode {
		return []any{}, errors.New("unclosed string")
	}

	//resolve
	for i, v := range a {
		vs := string(v.([]rune))

		if vi, err := strconv.Atoi(vs); err == nil { //int
			a[i] = vi
		} else if vs[0] == '"' && vs[len(vs)-1] == '"' { //string
			vs = vs[1 : len(vs)-1]
			vs = strings.ReplaceAll(vs, "\\|", "|")
			a[i] = strings.ReplaceAll(vs, "\\\"", "\"")
		} else if vs == "true" { //bool
			a[i] = true
		} else if vs == "false" { //bool
			a[i] = false
		} else if vf, err := strconv.ParseFloat(vs, 64); err == nil { //float
			a[i] = vf
		} else {
			return []any{}, errors.New("unexpect type")
		}
	}

	return a, nil
}

func Run() {
	// testencode := Encode("\\\"asdasa\"", 1, "test", "asdasda", 1145, "你好|", true, 3.1415926, "testwow\\|")
	// s := ""
	// for range 10000 {
	// 	s += "\\\"asdasa\""
	// }

	// start := time.Now()

	testencode := "100|1|\"2b\"|\"\""
	// testencode := "1"
	// testencode, err := Encode(s)
	// testencode, err := Encode("\\\"asdasa\"", 1, "test", "asdasda", 1145, "你好|", true, 3.1415926, "testwow\\|")
	// testencode, err := Encode(300)
	// if err != nil {
	// 	fmt.Println("Error", err)
	// 	return
	// }
	fmt.Println(testencode)
	decode, err := Decode(testencode)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	fmt.Println(decode...)

	// end := time.Since(start)

	for _, v := range decode {
		fmt.Printf("%T ", v)
	}
	fmt.Print("\n")
	fmt.Println(len(decode))

	// fmt.Print("encode and decode time:", end)

	// there are float string conv examples
	// fmt.Println(strconv.ParseFloat("3.1415926", 64))
	// fmt.Println(strconv.FormatFloat(3.1415926, 'f', -1, 64))
}
