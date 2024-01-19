package server_test

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/imirjar/metrx/internal/app/server"
// 	"github.com/stretchr/testify/assert"
// )

// func TestStatusHandler(t *testing.T) {
// 	type want struct {
// 		code        int
// 		response    string
// 		contentType string
// 	}
// 	tests := []struct {
// 		name string
// 		want want
// 	}{
// 		{
// 			name: "positive test #1",
// 			want: want{
// 				code:        200,
// 				response:    `{"status":"ok"}`,
// 				contentType: "application/json",
// 			},
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			server := server.NewServer()
// 			req := httptest.NewRequest(http.MethodPost, "/update/gauge/qwert/12", nil)
// 			// создаём новый Recorder
// 			resp := httptest.NewRecorder()
// 			server.UpdateGauge(resp, req)

// 			res := resp.Result()
// 			// проверяем код ответа
// 			assert.Equal(t, test.want.code, res.StatusCode)
// 			// получаем и проверяем тело запроса
// 			defer res.Body.Close()
// 			// resBody, err := io.ReadAll(res.Body)

// 			// require.NoError(t, err)
// 			// assert.JSONEq(t, test.want.response, string(resBody))
// 			// assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
// 		})
// 	}
// }
