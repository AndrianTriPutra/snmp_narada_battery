package domain

type Battery struct {
	Pack_ID uint8
	Data    DataBattery
}

type DataBattery struct {
	Current float32
	Voltage float32
	Soc     float32
	Cells   Cell
}

type Cell struct {
	Cell1  float32
	Cell2  float32
	Cell3  float32
	Cell4  float32
	Cell5  float32
	Cell6  float32
	Cell7  float32
	Cell8  float32
	Cell9  float32
	Cell10 float32
	Cell11 float32
	Cell12 float32
	Cell13 float32
	Cell14 float32
	Cell15 float32
}
