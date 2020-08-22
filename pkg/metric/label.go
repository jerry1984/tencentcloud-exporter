package metric

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/tencentyun/tencentcloud-exporter/pkg/instance"
	"sort"
)

type Labels map[string]string

func (l *Labels) Md5() (string, error) {
	h := md5.New()
	jb, err := json.Marshal(l)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(jb)), nil
}

type TcmLabels struct {
	queryLableNames    []string
	instanceLabelNames []string
	constLabels        Labels
	Names              []string
}

func (l *TcmLabels) GetValues(filters map[string]string, ins instance.TcInstance) (values []string, err error) {
	nameValues := map[string]string{}
	for _, name := range l.queryLableNames {
		v, ok := filters[name]
		if ok {
			nameValues[name] = v
		} else {
			nameValues[name] = ""
		}
	}
	for _, name := range l.instanceLabelNames {
		v, e := ins.GetFieldValueByName(name)
		if e != nil {
			nameValues[name] = ""
		} else {
			nameValues[name] = v
		}
	}
	for name, value := range l.constLabels {
		nameValues[name] = value
	}
	for _, name := range l.Names {
		values = append(values, nameValues[name])
	}
	return
}

func NewTcmLabels(qln []string, iln []string, cl Labels) (*TcmLabels, error) {
	var labelNames []string
	labelNames = append(labelNames, qln...)
	labelNames = append(labelNames, iln...)
	for lname := range cl {
		labelNames = append(labelNames, lname)
	}
	var uniq = map[string]bool{}
	for _, name := range labelNames {
		uniq[name] = true
	}
	var uniqLabelNames []string
	for n := range uniq {
		uniqLabelNames = append(uniqLabelNames, n)
	}
	sort.Strings(uniqLabelNames)

	l := &TcmLabels{
		queryLableNames:    qln,
		instanceLabelNames: iln,
		constLabels:        cl,
		Names:              uniqLabelNames,
	}
	return l, nil
}
