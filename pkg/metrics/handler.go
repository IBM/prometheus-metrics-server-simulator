package metrics

import (
	"encoding/json"
	"net/http"
)

type label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type valueBody struct {
	Name   string  `json:"name"`
	Labels []label `json:"labels,omitempty"`
	Value  float64 `json:"value"`
}

func SetValue(w http.ResponseWriter, req *http.Request) {
	if generator == nil {
		http.Error(w, "Server is not initialized yet", http.StatusTooEarly)
		return
	}
	m := req.Method
	if m != http.MethodPut {
		http.Error(w, "Put method is supported only", http.StatusMethodNotAllowed)
		return
	}
	var v valueBody
	err := json.NewDecoder(req.Body).Decode(&v)
	if err != nil {
		http.Error(w, "failed to parse body: "+err.Error(), http.StatusBadRequest)
	}
	generator.SetValue(v.Name, convertLabelsToMap(v.Labels), v.Value)
	w.WriteHeader(http.StatusOK)

}
func convertLabelsToMap(labels []label) map[string]string {
	labelmap := make(map[string]string)
	if labels == nil || len(labels) == 0 {
		return labelmap
	}
	for _, l := range labels {
		labelmap[l.Name] = l.Value
	}
	return labelmap

}
