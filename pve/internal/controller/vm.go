package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pve/internal/logic"
)

func CreateVM(c *gin.Context) {

	fmt.Println("create vm")
	logic.CreateVM()

}
