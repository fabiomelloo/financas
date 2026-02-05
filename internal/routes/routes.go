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

// Controllers contém todos os controllers da aplicação
type Controllers struct {
	Expense      *controllers.ExpenseController
	User         *controllers.UserController
	Purchase     *controllers.PurchaseController
	Gamification *controllers.GamificationController
}

func RegisterRoutes(c *Controllers) {
	// ============================================
	// Rotas de Despesas/Receitas (Finanças Pessoais)
	// ============================================
	http.HandleFunc("/", secureHandler(c.Expense.Index))
	http.HandleFunc("/create", secureHandler(c.Expense.Create))
	http.HandleFunc("/edit", secureHandler(c.Expense.Edit))
	http.HandleFunc("/update", secureHandler(c.Expense.Update))
	http.HandleFunc("/delete", secureHandler(c.Expense.Delete))
	http.HandleFunc("/insights", secureHandler(c.Expense.Insights))

	// ============================================
	// Rotas de Membros/Usuários (Equipe do Rateio)
	// ============================================
	http.HandleFunc("/users", secureHandler(c.User.Index))
	http.HandleFunc("/users/create", secureHandler(c.User.Create))
	http.HandleFunc("/users/delete", secureHandler(c.User.Delete))

	// ============================================
	// Rotas de Compras de Lanche (Rateio)
	// ============================================
	http.HandleFunc("/purchases", secureHandler(c.Purchase.Index))
	http.HandleFunc("/purchases/create", secureHandler(c.Purchase.Create))
	http.HandleFunc("/purchases/delete", secureHandler(c.Purchase.Delete))
	http.HandleFunc("/purchases/process", secureHandler(c.Purchase.ProcessMonth))

	// ============================================
	// Rotas de Gamificação (Ranking e Conquistas)
	// ============================================
	http.HandleFunc("/ranking", secureHandler(c.Gamification.Ranking))
	http.HandleFunc("/achievements", secureHandler(c.Gamification.Achievements))
}
