package middleware

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/fatih/color"
)

func Logging(next http.Handler) http.Handler {
	fn := func (w http.ResponseWriter, r *http.Request) {
		warna := []*color.Color{
			color.New(color.FgGreen),
			color.New(color.FgMagenta),
			color.New(color.FgCyan),
			color.New(color.FgYellow),
		}

		if (r.URL.Path != "/favicon.ico" && r.URL.Path != "/.well-known/appspecific/com.chrome.devtools.json") {
			waktu := time.Now().Format("15:04:05")
			random := rand.New(rand.NewSource(time.Now().UnixNano()))
			panjangSlice := len(warna)
			indeksAcak := random.Intn(panjangSlice)
			warnaTerpilih := warna[indeksAcak]
			warnaTerpilih.Printf("[%s] -> IP: %s | PATH: %s | METOD: %s\n", waktu, r.RemoteAddr, r.URL.Path, r.Method)
			next.ServeHTTP(w, r) 
		}
	}
	return http.HandlerFunc(fn)
}

func HandleLog(message string) {
	color.Blue(message)
}