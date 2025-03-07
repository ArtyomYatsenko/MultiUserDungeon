package main

import (
	"fmt"
)

// Глобальные переменные
var (
	GameWorld *World
)

// Игровой мир---------------------------------------------------------------
type World struct {
	Rooms  []Room
	Player Player
}

func (W *World) Handling(Cmd string) {

}

func NewWorld(R *[]Room, P *Player) *World {
	return &World{
		Rooms:  *R,
		Player: *P,
	}
}

// Комната--------------------------------------------------------------------
type Room struct {
	Name      string
	RoomItems RoomItems
}

func (R *Room) Processing(Cmd string) {

}

func NewRoom(Name string) (*Room, error) {

	switch Name {
	case "Кухня":
		return &Room{
			Name:      "Кухня",
			RoomItems: NewRoomItems(Name),
		}, nil
	case "Коридор":
		return &Room{
			Name:      "Коридор",
			RoomItems: NewRoomItems(Name),
		}, nil
	case "Комната":
		return &Room{
			Name:      "Комната",
			RoomItems: NewRoomItems(Name),
		}, nil
	case "Улица":
		return &Room{
			Name:      "Улица",
			RoomItems: NewRoomItems(Name),
		}, nil
	}

	return &Room{}, fmt.Errorf("Error create room")
}

func NewRooms(Names ...string) []Room {
	Rooms := []Room{}
	for _, Name := range Names {
		Room, err := NewRoom(Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		Rooms = append(Rooms, *Room)
	}
	return Rooms
}

// Предметы комнаты-----------------------------------------------------------
type RoomItems map[string][]Item

func NewRoomItems(Name string) RoomItems {

	switch Name {
	case "Кухня":
		return RoomItems{"Стол": {"Чай"}}
	case "Коридор":
		return RoomItems{"Стол": {"Чай"}}
	case "Комната":
		return RoomItems{"Стол": {"Чай"}}
	case "Улица":
		return RoomItems{"Стол": {"Чай"}}
	default:
		return make(map[string][]Item)
	}

}

// Предметы-------------------------------------------------------------------
type Item string

// Игрок----------------------------------------------------------------------
type Player struct {
	Name        string
	Backpack    bool
	Inventory   Inventory
	CurrentRoom *Room
}

func NewPlayer() *Player {
	return &Player{
		Name:        "Viktor",
		Backpack:    false,
		Inventory:   make(map[string]Item),
		CurrentRoom: &Room{},
	}
}

// Инвентарь игрока------------------------------------------------------------
type Inventory map[string]Item

// На развитие...
func (I Inventory) IsEmpty() bool {
	if len(I) == 0 {
		return true
	}
	return false
}

func NewInventory() Inventory {
	return make(map[string]Item)
}

// Для себя проверка (Знаю, что можно через _test, просто практикуюсь в написании кода)
type GameCaseTest [][]string

func NewGameCaseTest() GameCaseTest {
	return GameCaseTest{
		{"осмотреться",
			"идти коридор",
			"идти комната",
			"осмотреться",
			"надеть рюкзак",
			"взять ключи",
			"взять конспекты",
			"идти коридор",
			"применить ключи дверь",
			"идти улица"},
		{"осмотреться",
			"идти коридор",
			"идти комната",
			"осмотреться",
			"надеть рюкзак",
			"взять ключи",
			"взять конспекты",
			"идти коридор",
			"применить ключи дверь",
			"идти улица"},
	}
}

// Инициализация необходимых объектов-------------------------------------------
func initGame() {
	Player := NewPlayer()
	Rooms := NewRooms()
	GameWorld = NewWorld(&Rooms, Player)
}

// ОбработкаДействия-------------------------------------------------------------
func handleCommand(Cmd string) string {

	return ""
}

func main() {
	CasesTest := NewGameCaseTest()
	initGame()
	for _, Case := range CasesTest {
		for _, el := range Case {
			handleCommand(el)
		}
	}
}

/* Ссылку или копию должны возвращать конструкторы экземпляров
строкутур/мап/массивов/слайсов/примитивныхт типов?

В Функции main в каком случае следует называть переменные с большой буквы, а в каком
с маленькой? Cases/Case/el и тд?


*/
