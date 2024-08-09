package handler

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"rest.gtld.test/realTimeApp/app/model"
	"rest.gtld.test/realTimeApp/app/usecases"
)

type weatherHandler struct {
	WeatherUsecaseImp *usecases.WeatherUsecaseImp
}

func NewWeatherHandler(weatherUsecase *usecases.WeatherUsecaseImp) *weatherHandler {
	return &weatherHandler{
		WeatherUsecaseImp: weatherUsecase,
	}
}

func (w *weatherHandler) HandleWebSocketConnection(conn *websocket.Conn) {
	for {
		var weather model.Weather
		_, message, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			break
		}
		json.Unmarshal(message, &weather)
		w.WeatherUsecaseImp.WeatherDataProcessing(&weather)
	}
}
