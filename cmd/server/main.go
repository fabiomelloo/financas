package main

import (
	"financas/database"
	"financas/internal/controllers"
	"financas/internal/repositories"
	"financas/internal/routes"
	"financas/internal/services"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// ============================================
	// Inicializar conexão com o banco de dados
	// ============================================
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Falha ao conectar ao banco de dados:", err)
	}
	defer db.Close()

	// ============================================
	// Inicializar Repositories (Acesso a Dados)
	// ============================================
	expenseRepo := repositories.NewExpenseRepository(db)
	userRepo := repositories.NewUserRepository(db)
	purchaseRepo := repositories.NewPurchaseRepository(db)
	achievementRepo := repositories.NewAchievementRepository(db)

	// ============================================
	// Inicializar Services (Regras de Negócio)
	// ============================================
	expenseService := services.NewExpenseService(expenseRepo)
	userService := services.NewUserService(userRepo)
	purchaseService := services.NewPurchaseService(purchaseRepo, userRepo)
	gamificationService := services.NewGamificationService(userRepo, purchaseRepo, achievementRepo)

	// ============================================
	// Inicializar Controllers (HTTP Handlers)
	// ============================================
	expenseController := controllers.NewExpenseController(expenseService)
	userController := controllers.NewUserController(userService)
	purchaseController := controllers.NewPurchaseController(purchaseService, userService, gamificationService)
	gamificationController := controllers.NewGamificationController(gamificationService, purchaseService)

	// ============================================
	// Registrar Rotas
	// ============================================
	allControllers := &routes.Controllers{
		Expense:      expenseController,
		User:         userController,
		Purchase:     purchaseController,
		Gamification: gamificationController,
	}
	routes.RegisterRoutes(allControllers)

	// Servir arquivos estáticos (CSS, JS, imagens)
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// ============================================
	// Iniciar servidor HTTP
	// ============================================
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║  🚀 Servidor Finanças + Rateio rodando!                      ║")
	fmt.Println("║                                                              ║")
	fmt.Println("║  📊 Dashboard:      http://localhost:8080                    ║")
	fmt.Println("║  👥 Membros:        http://localhost:8080/users              ║")
	fmt.Println("║  🥪 Compras:        http://localhost:8080/purchases          ║")
	fmt.Println("║  🏆 Ranking:        http://localhost:8080/ranking            ║")
	fmt.Println("║  🏅 Conquistas:     http://localhost:8080/achievements       ║")
	fmt.Println("║  📈 Relatórios:     http://localhost:8080/insights           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Falha ao iniciar servidor:", err)
	}
}
