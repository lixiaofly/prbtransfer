package prbtransfer

import (
	//"prbtransfer/prbtransfer"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

type TblStaProbInfos struct {
	//ID                    int64 `orm:"pk;auto"`
	//ApMac                 string
	//NowTime               time.Time `orm:"type(datetime)"`
	//WlanMode              string
	WlanMac               string
	StaBrand              string
	StaCacheSsid          string
	CaptureTime           int64
	WlanRssi              int32
	Identification        int32
	CertificateCode       string
	WlanEssid             string
	AccessAPMac           string
	WlanChannel           int32
	WlanEncrypion         string
	XCoordinate           int64
	YCoordinate           int64
	NetbarWacode          string
	CollectionEquipmentId string
	Longitude             float64
	Latitude              float64
}
type STblStaProbInfos []TblStaProbInfos

func TestProbe(t *testing.T) {

	pi := TblStaProbInfos{
		WlanMac:               "11-11-11-11-11-11",
		StaBrand:              "huawei",
		StaCacheSsid:          "hao123",
		CaptureTime:           1458021145,
		WlanRssi:              -22,
		Identification:        0,
		CertificateCode:       "",
		WlanEssid:             "hao123",
		AccessAPMac:           "22-22-22-22-22-22",
		WlanChannel:           6,
		WlanEncrypion:         "",
		XCoordinate:           5,
		YCoordinate:           1,
		NetbarWacode:          "22331234567890",
		CollectionEquipmentId: "333333333333101909571",
		Longitude:             123.23000,
		Latitude:              133.000000,
	}
	piSlice := []TblStaProbInfos{pi}
	spiSlice := STblStaProbInfos(piSlice)
	fmt.Println("original content=", spiSlice)
	content,_ := ProbDataDeal(&spiSlice)
	//tcpAddr, err := net.ResolveTCPAddr("tcp4", "61.141.136.6:18040")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "121.43.231.237:7777")
	if err != nil {
		fmt.Println("ResolveTCPAddr failed!")
		return
	}
	fmt.Println("Send:", content,"len=", len(content))
	//for {
		TcpUpload(content, tcpAddr)
		time.Sleep(1)
	//}
	return
}

func (this *STblStaProbInfos) ProbFormatComb() []byte {
	jaSlice := make([]string, 0)
	for _, v := range *this {
		a := []string{v.WlanMac, v.StaBrand, v.StaCacheSsid, toString(v.CaptureTime), toString(v.WlanRssi),
			toString(v.Identification), v.CertificateCode, v.WlanEssid, v.AccessAPMac, toString(v.WlanChannel),
			v.WlanEncrypion, toString(v.XCoordinate), toString(v.YCoordinate), v.NetbarWacode, v.CollectionEquipmentId,
			toString(v.Longitude), toString(v.Latitude)}
		ja := strings.Join(a, "\t")
		jaSlice = append(jaSlice, ja)
		//fmt.Println("byte1:", jaSlice)
	}
	jasj := strings.Join(jaSlice, "\r\n")
	//fmt.Println("byte2:[", jasj, "]")
	//fmt.Println("len(jasj)=",len(jasj))
	//fmt.Println("cap(jasj)=",cap(jasj))
	return []byte(jasj)
}

func (this *STblStaProbInfos) ProbDataType() uint32 {
	//return []byte{'1', '0', '0', '1'}
	return 1001
}
func (this *STblStaProbInfos) DesKey() []byte {
	return []byte("pk$@gtjt")
}
func (this *STblStaProbInfos) DesIv() []byte {
	return []byte("thvn#&@@")
}

func toString(a interface{}) string {
	switch a.(type) {
	case int32:
		return fmt.Sprintf("%d", a.(int32))
	case int64:
		return fmt.Sprintf("%d", a.(int64))
	case float64:
		return fmt.Sprintf("%f", a.(float64))
	case string:
		return a.(string)
	default:
		return ""
	}
}
