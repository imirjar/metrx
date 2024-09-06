package truster

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func Truster(cidr string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if cidr != "" {
				ip := r.Header.Get("X-Real-IP")
				log.Println("cidr--->", cidr)

				// Преобразуем строку с IP-адресом в тип net.IP
				parsedIP := net.ParseIP(ip)
				if parsedIP == nil {
					fmt.Println("Неверный IP-адрес")
					return
				}

				// Преобразуем строку с CIDR в тип net.IPNet и базовый IP-адрес
				_, ipNet, err := net.ParseCIDR(cidr)
				if err != nil {
					fmt.Println("Неверный CIDR")
					return
				}

				// Проверяем, входит ли IP в диапазон CIDR
				if ipNet.Contains(parsedIP) {
					// fmt.Println("IP-адрес входит в диапазон CIDR")
					next.ServeHTTP(w, r)
				} else {
					fmt.Println("IP-адрес не входит в диапазон CIDR")
					w.WriteHeader(http.StatusForbidden)
					return
				}
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
