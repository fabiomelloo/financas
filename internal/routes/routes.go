package routes

import (
	"financas/internal/controllers"
	"net/http"
)

func RegisterRoutes(controller *controllers.ExpenseController) {
	http.HandleFunc("/", controller.Index)
	http.HandleFunc("/create", controller.Create)
	http.HandleFunc("/edit", controller.Edit)
	http.HandleFunc("/update", controller.Update)
	http.HandleFunc("/delete", controller.Delete)
}
