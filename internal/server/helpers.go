package server

import "net/http"

func getFormValues(r *http.Request, key ...string) (map[string]string, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	values := make(map[string]string)
	for _, k := range key {
		values[k] = r.PostForm.Get(k)
	}

	return values, err
}
