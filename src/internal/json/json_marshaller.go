package json_order

import(
	"myapp/src/internal/order_model"
	"encoding/json"
)

func Marshall_order(order *model.Order) ([]byte , error){
	jsonData , err := json.Marshal(*order)
	if err!= nil {
		return nil , err
	}
	return jsonData , nil
}

func Unmarshal_order(json_order []byte) (*model.Order , error){
	var order model.Order
	err:= json.Unmarshal(json_order,&order)
	if err!= nil {
		return  nil, err
	}
	return &order , nil
}