package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ValueRequest struct {
	AccessKey int64 `json:"access_key"`
	Value     int   `json:"value"`
}

type Request struct {
	DepreciationId int64 `json:"depreciation_id"`
	Mileage        int   `json:"mileage"`
	InitialPrice   int   `json:"initial_price"`
}

type DepreciationData struct {
	Mileage      int `json:"mileage"`
	InitialPrice int `json:"initial_price"`
}

func (h *Handler) issueValue(c *gin.Context) {
	var input Request
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Printf("handler.issueValue: DepreciationId=%d, Mileage=%d, InitialPrice=%d\n",
		input.DepreciationId, input.Mileage, input.InitialPrice)

	c.Status(http.StatusOK)

	go func() {
		time.Sleep(3 * time.Second)
		sendValueRequest(input)
	}()
}

func sendValueRequest(request Request) {
	var value int

	// Вычисляем значение по формуле (mileage * initial_price) / 500000
	if request.Mileage > 0 && request.InitialPrice > 0 {
		// Умножаем перед делением для сохранения точности
		value = (request.Mileage * request.InitialPrice) / 500000
		fmt.Printf("Расчет по формуле: (%d * %d) / 500000 = %d\n",
			request.Mileage, request.InitialPrice, value)
	}

	answer := ValueRequest{
		AccessKey: 123,
		Value:     value,
	}

	client := &http.Client{}

	jsonAnswer, _ := json.Marshal(answer)
	bodyReader := bytes.NewReader(jsonAnswer)

	requestURL := fmt.Sprintf("http://django:8000/api/depreciations/%d/update_summ/", request.DepreciationId)

	req, _ := http.NewRequest(http.MethodPut, requestURL, bodyReader)
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending PUT request:", err)
		return
	}

	defer response.Body.Close()

	fmt.Println("Результат вычислений:", value)
	fmt.Println("PUT Request Status:", response.Status)
}
