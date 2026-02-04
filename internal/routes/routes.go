package routes

import (
	"financas/internal/controllers"
	"net/http"
)

// secureHandler adiciona cabeçalhos de segurança básicos
func secureHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' https://cdn.jsdelivr.net; "+
				"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; "+
				"font-src 'self' https://fonts.gstatic.com; "+
				"img-src 'self' data:; "+
				"connect-src 'self';",
		)

		h.ServeHTTP(w, r)
	}
}

func RegisterRoutes(controller *controllers.ExpenseController) {
	http.HandleFunc("/", secureHandler(controller.Index))
	http.HandleFunc("/create", secureHandler(controller.Create))
	http.HandleFunc("/edit", secureHandler(controller.Edit))
	http.HandleFunc("/update", secureHandler(controller.Update))
	http.HandleFunc("/delete", secureHandler(controller.Delete))
	http.HandleFunc("/insights", secureHandler(controller.Insights))
	http.HandleFunc("/rateio", secureHandler(controller.Rateio))
}
