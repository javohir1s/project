package api

import (
	"github.com/gin-gonic/gin"

	"exam/api/handler"
	"exam/config"
	"exam/storage"

	_ "exam/api/docs"

	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func SetUpApi(r *gin.Engine, cfg *config.Config, strg storage.StorageI) {

	handler := handler.NewHandler(cfg, strg)

	r.POST("/make_pay", handler.MakePay)
	// picking_list ...
	r.POST("/picking_list", handler.CreatePickingList)
	r.GET("/picking_list/:id", handler.GetByIDPickingList)
	r.GET("/picking_list", handler.GetListPickingList)
	r.PUT("/picking_list/:id", handler.UpdatePickingList)
	r.DELETE("/picking_list/:id", handler.DeletePickingList)

	// sale_product ...
	r.POST("/saleproduct", handler.CreateSaleProduct)
	r.GET("/saleproduct/:id", handler.GetByIDSaleProduct)
	r.GET("/saleproduct", handler.GetListSaleProduct)
	r.PUT("/saleproduct/:id", handler.UpdateSaleProduct)
	r.DELETE("/saleproduct/:id", handler.DeleteSaleProduct)

	// sale ...
	r.POST("/sale", handler.CreateSale)
	r.GET("/sale/:id", handler.GetByIDSale)
	r.GET("/sale", handler.GetListSale)
	r.PUT("/sale/:id", handler.UpdateSale)
	r.DELETE("/sale/:id", handler.DeleteSale)

	// product ...
	r.POST("/product", handler.CreateProduct)
	r.GET("/product/:id", handler.GetByIDProduct)
	r.GET("/product", handler.GetListProduct)
	r.PUT("/product/:id", handler.UpdateProduct)
	r.DELETE("/product/:id", handler.DeleteProduct)

	// remainder ...
	r.POST("/remainder", handler.CreateRemainder)
	r.GET("/remainder/:id", handler.GetByIDRemainder)
	r.GET("/remainder", handler.GetListRemainder)
	r.PUT("/remainder/:id", handler.UpdateRemainder)
	r.DELETE("/remainder/:id", handler.DeleteRemainder)

	// client ...
	r.POST("/client", handler.CreateClient)
	r.GET("/client/:id", handler.GetByIDClient)
	r.GET("/client", handler.GetListClient)
	r.PUT("/client/:id", handler.UpdateClient)
	r.DELETE("/client/:id", handler.DeleteClient)

	// branch ...
	r.POST("/branch", handler.CreateBranch)
	r.GET("/branch/:id", handler.GetByIDBranch)
	r.GET("/branch", handler.GetListBranch)
	r.PUT("/branch/:id", handler.UpdateBranch)
	r.DELETE("/branch/:id", handler.DeleteBranch)

	// coming
	r.POST("/coming", handler.CreateComing)
	r.GET("/coming/:id", handler.GetByIDComing)
	r.GET("/coming", handler.GetListComing)
	r.PUT("/coming/:id", handler.UpdateComing)
	r.DELETE("/coming/:id", handler.DeleteComing)

	r.GET("/report", handler.OverallReport)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func customCORSMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE, HEAD")
		c.Header("Access-Control-Allow-Headers", "Password, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
