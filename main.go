package main

import (
	"fmt"
	"sort"
	"strings"
)

// Глобальные переменные
var (
	GameWorld *World
)

// Игровой мир---------------------------------------------------------------
type World struct {
	Rooms           Rooms
	Player          *Player
	InteractionItem InteractionItem
}

func (W *World) Processing(Cmd string, arg ...string) string {
	return W.Player.PerformAction(Cmd, arg)
}

func NewWorld(R Rooms, P *Player, InteractionItem InteractionItem) *World {
	return &World{
		Rooms:           R,
		Player:          P,
		InteractionItem: InteractionItem,
	}
}

// Описание-------------------------------------------------------------------
type Description map[string]string

func (D Description) Add(NameMove, Description string) {
	D[NameMove] = Description
}

func (D Description) GetDescriptionString(Cmd string) string {
	return D[Cmd]
}

func NewDescription() *Description {
	d := make(Description)
	return &d
}

// Задачи---------------------------------------------------------------------
type Task map[string]bool

func (T Task) Add(Name string, Complete bool) {
	T[Name] = Complete
}

func (T Task) GetTaskString(Cmd string, Backpack bool) string {
	if Cmd != "осмотреться" {
		return ""
	}
	for task, complete := range T {
		if Backpack {
			complete = !complete
		}
		if complete {
			return task
		}
	}
	return ""
}

func NewTusk() *Task {
	t := make(Task)
	return &t
}

// Комната--------------------------------------------------------------------
type Room struct {
	Name        string
	RoomItems   *RoomItems
	RoomAllowed []*Room
	Description *Description
	Task        *Task
	CloseDoor   bool
}

func (R Room) MoveAllowedString() string {
	str := ". можно пройти - "
	for _, Room := range R.RoomAllowed {
		str += Room.Name + ", "
	}
	return RemoveLastChar(str)
}

func (R *Room) GetDescription(Cmd string, Backpack bool) string {

	return R.Description.GetDescriptionString(Cmd) + R.RoomItems.GetRoomItemsString(Cmd) + R.Task.GetTaskString(Cmd, Backpack) + R.MoveAllowedString()
}

func (R *Room) SetRoomAllowed(Room ...*Room) {
	R.RoomAllowed = Room
}

func NewRoom(Name string, Items *RoomItems, Task *Task, Desc *Description, CloseDoor bool) *Room {
	return &Room{
		Name:        Name,
		RoomItems:   Items,
		Task:        Task,
		Description: Desc,
		CloseDoor:   CloseDoor,
	}
}

// Комнаты--------------------------------------------------------------------
type Rooms []*Room

func (R Rooms) GetDefaultRoom() (*Room, error) {
	for _, el := range R {
		if el.Name == "кухня" {
			return el, nil
		}
	}
	return &Room{}, fmt.Errorf("Комната по умолчанию не найдена")
}

func (R Rooms) GetRoom(Name string) *Room {
	for _, el := range R {
		if el.Name == Name {
			return el
		}
	}
	return &Room{}
}

func (R *Room) OpenDoor() {
	for _, room := range R.RoomAllowed {
		room.CloseDoor = false // Теперь изменяется оригинальный объект
	}
}

// Предметы комнаты-----------------------------------------------------------
type RoomItems map[string][]Item

func (R RoomItems) Add(Name string, Itm ...string) {
	Items := []Item{}
	for _, el := range Itm {
		Items = append(Items, Item(el))
	}
	R[Name] = Items
}

func (R RoomItems) GetRoomItemsString(Cmd string) string {
	if Cmd != "осмотреться" {
		return ""
	}

	keys := make([]string, 0)

	for key := range R {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	str := ""

	for _, key := range keys {
		str += "на " + key + "е: "
		for _, item := range R[key] {
			str += string(item) + ", "
		}
	}

	if str == "" {
		return "пустая комната"
	}
	return RemoveLastChar(str)

}

func (R RoomItems) HasItem(NameItem string) bool {
	for _, Items := range R {
		for _, Item := range Items {
			if string(Item) == NameItem {
				return true
			}
		}
	}

	return false
}

func (R RoomItems) DeleteItem(Name Item) {
	for key, Items := range R {
		for i, Item := range Items {
			if Item == Name {
				R[key] = append(Items[:i], Items[i+1:]...)
				break
			}
		}
	}
	for key, Items := range R {
		if len(Items) == 0 {
			delete(R, key)
		}
	}
}

func NewRoomItems() *RoomItems {
	r := make(RoomItems)
	return &r

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

func (P *Player) IsChangeRoom(NameRoom string) bool {
	for _, Room := range P.CurrentRoom.RoomAllowed {
		if Room.Name == NameRoom {
			return true
		}
	}
	return false
}

func (P *Player) PerformAction(Cmd string, Arg []string) string {
	switch Cmd {
	case "осмотреться":
		return P.LookAround(Cmd, P.Backpack)
	case "идти":
		return P.Go(Cmd, Arg[0])
	case "надеть":
		return P.PutOn(Cmd, Arg[0])
	case "взять":
		return P.Take(Cmd, Arg[0])
	case "применить":
		return P.Apply(Cmd, Arg)
	default:
		return "неизвестная команда"
	}
}

func (P *Player) Go(Cmd, Arg string) string {
	if !P.IsChangeRoom(Arg) {
		return "нет пути в " + Arg
	}
	Room := GameWorld.Rooms.GetRoom(Arg)
	if Room.CloseDoor {
		return "дверь закрыта"

	}

	P.CurrentRoom = GameWorld.Rooms.GetRoom(Arg)

	return P.CurrentRoom.GetDescription(Cmd, P.Backpack)
}

func (P *Player) Take(Cmd, Arg string) string {
	if !P.CurrentRoom.RoomItems.HasItem(Arg) {
		return "нет такого"
	}
	if !P.Backpack {
		return "некуда класть"
	}

	P.CurrentRoom.RoomItems.DeleteItem(Item(Arg))
	P.Inventory.Add(Item(Arg))
	return "предмет добавлен в инвентарь: " + Arg
}

func (P *Player) PutOn(Cmd, Arg string) string {
	if !P.CurrentRoom.RoomItems.HasItem(Arg) {
		return "нет такого"
	}

	P.Backpack = true
	P.CurrentRoom.RoomItems.DeleteItem(Item(Arg))

	return "вы надели: " + Arg
}

func (P *Player) Apply(Cmd string, Arg []string) string {
	if !P.Inventory.HasItem(Arg[0]) {
		return "нет предмета в инвентаре - " + Arg[0]
	}
	if !GameWorld.InteractionItem.IsApplicationAllowed(Arg) {
		return "не к чему применить"
	}
	P.CurrentRoom.OpenDoor()

	return "дверь открыта"
}

func (P *Player) LookAround(Cmd string, Bacpack bool) string {
	return P.CurrentRoom.GetDescription(Cmd, Bacpack)
}

func NewPlayer(Room *Room) *Player {
	return &Player{
		Name:        "Viktor",
		Backpack:    false,
		Inventory:   make(map[string]Item),
		CurrentRoom: Room,
	}
}

// Инвентарь игрока------------------------------------------------------------
type Inventory map[string]Item

func (I Inventory) Add(Name Item) {
	I[string(Name)] = Name
}

func (I Inventory) HasItem(NameItem string) bool {
	_, err := I[NameItem]

	return err
}

func (I Inventory) DeleteItem(Name string) {
	delete(I, Name)
}

func (I Inventory) IsEmpty() bool {
	if len(I) == 0 {
		return true
	}
	return false
}

func NewInventory() Inventory {
	return make(map[string]Item)
}

// Взаимодействие предметов в мире--------------------------------------------------
type InteractionItem map[string]string

func NewInetactionItem() *InteractionItem {
	I := make(InteractionItem)
	return &I
}

func (I InteractionItem) Add(Item1, Item2 string) {
	I[Item1] = Item2
}

func (I InteractionItem) IsApplicationAllowed(Arg []string) bool {
	if I[Arg[0]] == Arg[1] {
		return true
	}
	return false
}

// Для себя проверка (Знаю, что можно через _test, так удобнее корректировать)
type GameCaseTest [][]string

func NewGameCaseTest() GameCaseTest {
	return GameCaseTest{
		{"идти коридор", "идти комната", "осмотреться", "надеть рюкзак", "взять ключи", "взять конспекты", "идти коридор", "применить ключи дверь", "идти улица"},
	}
}

func initRooms() Rooms {

	var Rooms []*Room // Используем срез указателей на Room

	KitchenDescriptions := NewDescription()
	KitchenDescriptions.Add("осмотреться", "ты находишься на кухне, ")
	KitchenDescriptions.Add("идти", "кухня, ничего интересного")
	KitchenTask := NewTusk()
	KitchenTask.Add(", надо собрать рюкзак и идти в универ", true)
	KitchenTask.Add(", надо идти в универ", false)
	KitchenItems := NewRoomItems()
	KitchenItems.Add("стол", "чай")
	KitchenRoom := NewRoom("кухня", KitchenItems, KitchenTask, KitchenDescriptions, false)

	CorridorDescriptions := NewDescription()
	CorridorDescriptions.Add("осмотреться", "ничего интересного")
	CorridorDescriptions.Add("идти", "ничего интересного")
	CorridorTask := NewTusk()
	CorridorItems := NewRoomItems()
	CorridorRoom := NewRoom("коридор", CorridorItems, CorridorTask, CorridorDescriptions, false)

	MyRoomDescriptions := NewDescription()
	MyRoomDescriptions.Add("идти", "ты в своей комнате")
	MyRoomTask := NewTusk()
	MyRoomItems := NewRoomItems()
	MyRoomItems.Add("стул", "рюкзак")
	MyRoomItems.Add("стол", "ключи", "конспекты")
	MyRoom := NewRoom("комната", MyRoomItems, MyRoomTask, MyRoomDescriptions, false)

	StreetDescriptions := NewDescription()
	StreetDescriptions.Add("идти", "на улице весна")
	StreetTask := NewTusk()
	StreetItems := NewRoomItems()
	Street := NewRoom("улица", StreetItems, StreetTask, StreetDescriptions, true)

	HomeDescriptions := NewDescription()
	HomeDescriptions.Add("идти", "на улице весна")
	HomeTask := NewTusk()
	HomeItems := NewRoomItems()
	Home := NewRoom("домой", HomeItems, HomeTask, HomeDescriptions, false)

	KitchenRoom.SetRoomAllowed(CorridorRoom)
	CorridorRoom.SetRoomAllowed(KitchenRoom, MyRoom, Street)
	MyRoom.SetRoomAllowed(CorridorRoom)
	Street.SetRoomAllowed(Home)

	Rooms = append(Rooms, KitchenRoom, CorridorRoom, MyRoom, Street, Home)

	return Rooms // Теперь возвращаем срез указателей на Room
}

// Инициализация необходимых объектов-------------------------------------------
func initGame() {
	Rooms := initRooms()
	DefaultRoom, err := Rooms.GetDefaultRoom()
	if err != nil {
		panic("Не определена комната по умолчанию")
	}
	Player := NewPlayer(DefaultRoom)

	InteractionItemWorld := NewInetactionItem()
	InteractionItemWorld.Add("ключи", "дверь")
	GameWorld = NewWorld(Rooms, Player, *InteractionItemWorld)

}

// ОбработкаДействия-------------------------------------------------------------
func handleCommand(Cmd string) string {
	str := strings.Split(Cmd, " ")
	return GameWorld.Processing(str[0], str[1:]...)
}

func RemoveLastChar(s string) string {

	runes := []rune(s)

	if len(runes) == 0 {
		return s
	}

	runes = runes[:len(runes)-2]

	return string(runes)
}

func main() {
	CasesTest := NewGameCaseTest()
	initGame()

	for _, Case := range CasesTest {
		for _, el := range Case {
			fmt.Println(handleCommand(el))
		}
	}
}

/* Ссылку или копию должны возвращать конструкторы экземпляров
строкутур/мап/массивов/слайсов/примитивныхт типов?

В Функции main в каком случае следует называть переменные с большой буквы, а в каком
с маленькой? Cases/Case/el и тд?


		{"осмотреться",
"завтракать",
"идти комната",
"идти коридор",
"применить ключи дверь",
"идти комната", "осмотреться",
"взять ключи",
"надеть рюкзак",
"осмотреться",
"взять ключи",
"взять телефон",
"взять ключи",
"осмотреться",
"взять конспекты",
"осмотреться",
"идти коридор",
"идти кухня",
"осмотреться",
"идти коридор",
"идти улица",
"применить ключи дверь",
"применить телефон шкаф",
"применить ключи шкаф",
"идти улица"},

	}


осмотреться", "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"}, // действие осмотреться
		"идти коридор","идти комната","осмотреться","надеть рюкзак","взять ключи","взять конспекты","идти коридор","применить ключи дверь","идти улица"
*/
