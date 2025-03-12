package main

import (
	"fmt"
	"sort"
	"strings"
)

// Глобальные переменные

var gameWorld *World

// Игровой мир---------------------------------------------------------------
type World struct {
	Rooms           Rooms
	Player          *Player
	InteractionItem InteractionItem
}

func (w *World) Processing(cmd string, arg ...string) string {
	return w.Player.PerformAction(cmd, arg)
}

func NewWorld(r Rooms, p *Player, interactionItem InteractionItem) *World {
	return &World{
		Rooms:           r,
		Player:          p,
		InteractionItem: interactionItem,
	}
}

// Описание-------------------------------------------------------------------
type Description map[string]string

func (d Description) Add(nameMove, description string) {
	d[nameMove] = description
}

func (d Description) GetDescriptionString(cmd string) string {
	return d[cmd]
}

func NewDescription(desc map[string]string) *Description {
	d := Description(desc)
	return &d
}

// Задачи---------------------------------------------------------------------
type Tasks map[string]bool

func (T Tasks) Add(name string, complete bool) {
	T[name] = complete
}

func (T Tasks) GetTaskString(cmd string, backpack bool) string {
	if cmd != "осмотреться" {
		return ""
	}
	for task, complete := range T {
		if backpack {
			complete = !complete
		}
		if complete {
			return task
		}
	}
	return ""
}

func NewTusks(m map[string]bool) *Tasks {
	t := Tasks(m)
	return &t
}

// Комната--------------------------------------------------------------------
type Room struct {
	Name        string
	RoomItems   *RoomItems
	RoomAllowed []*Room
	Description *Description
	Task        *Tasks
	CloseDoor   bool
}

func (r Room) MoveAllowedString() string {
	str := ". можно пройти - "
	for _, room := range r.RoomAllowed {
		str += room.Name + ", "
	}
	return RemoveLastChar(str)
}

func (r *Room) GetDescription(cmd string, backpack bool) string {
	return r.Description.GetDescriptionString(cmd) + r.RoomItems.GetRoomItemsString(cmd) + r.Task.GetTaskString(cmd, backpack) + r.MoveAllowedString()
}

func (r *Room) SetRoomAllowed(room ...*Room) {
	r.RoomAllowed = room
}

func NewRoom(name string, items *RoomItems, task *Tasks, desc *Description, closeDoor bool) *Room {
	return &Room{
		Name:        name,
		RoomItems:   items,
		Task:        task,
		Description: desc,
		CloseDoor:   closeDoor,
	}
}

// Комнаты--------------------------------------------------------------------
type Rooms []*Room

func (r Rooms) GetDefaultRoom() (*Room, error) {
	for _, el := range r {
		if el.Name == "кухня" {
			return el, nil
		}
	}
	return &Room{}, fmt.Errorf("Комната по умолчанию не найдена")
}

func (r Rooms) GetRoom(name string) *Room {
	for _, el := range r {
		if el.Name == name {
			return el
		}
	}
	return &Room{}
}

func (r *Room) OpenDoor() {
	for _, room := range r.RoomAllowed {
		room.CloseDoor = false // Теперь изменяется оригинальный объект
	}
}

// Предметы комнаты-----------------------------------------------------------
type RoomItems map[string][]Item

func (r RoomItems) Add(name string, itm ...string) {
	Items := make([]Item, 0, 3)
	for _, el := range itm {
		Items = append(Items, Item(el))
	}
	r[name] = Items
}

func (r RoomItems) GetRoomItemsString(cmd string) string {
	if cmd != "осмотреться" {
		return ""
	}

	keys := make([]string, 0, len(r))

	for key := range r {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	str := ""

	for _, key := range keys {
		str += "на " + key + "е: "
		for _, item := range r[key] {
			str += string(item) + ", "
		}
	}

	if str == "" {
		return "пустая комната"
	}
	return RemoveLastChar(str)

}

func (r RoomItems) HasItem(nameItem string) bool {
	for _, items := range r {
		for _, item := range items {
			if string(item) == nameItem {
				return true
			}
		}
	}
	return false
}

func (r RoomItems) DeleteItem(name Item) {
	for key, items := range r {
		for i, item := range items {
			if item == name {
				r[key] = append(items[:i], items[i+1:]...)
				break
			}
		}
	}
	for key, items := range r {
		if len(items) == 0 {
			delete(r, key)
		}
	}
}

func NewRoomItems(i map[string][]Item) *RoomItems {
	r := RoomItems(i)
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

func (p *Player) IsChangeRoom(nameRoom string) bool {
	for _, room := range p.CurrentRoom.RoomAllowed {
		if room.Name == nameRoom {
			return true
		}
	}
	return false
}

func (p *Player) PerformAction(cmd string, arg []string) string {
	switch cmd {
	case "осмотреться":
		return p.LookAround(cmd, p.Backpack)
	case "идти":
		return p.Go(cmd, arg[0])
	case "надеть":
		return p.PutOn(arg[0])
	case "взять":
		return p.Take(arg[0])
	case "применить":
		return p.Apply(arg)
	default:
		return "неизвестная команда"
	}
}

func (p *Player) Go(cmd, arg string) string {
	if !p.IsChangeRoom(arg) {
		return "нет пути в " + arg
	}
	Room := gameWorld.Rooms.GetRoom(arg)
	if Room.CloseDoor {
		return "дверь закрыта"

	}

	p.CurrentRoom = gameWorld.Rooms.GetRoom(arg)

	return p.CurrentRoom.GetDescription(cmd, p.Backpack)
}

func (p *Player) Take(arg string) string {
	if !p.CurrentRoom.RoomItems.HasItem(arg) {
		return "нет такого"
	}
	if !p.Backpack {
		return "некуда класть"
	}

	p.CurrentRoom.RoomItems.DeleteItem(Item(arg))
	p.Inventory.Add(Item(arg))
	return "предмет добавлен в инвентарь: " + arg
}

func (p *Player) PutOn(arg string) string {
	if !p.CurrentRoom.RoomItems.HasItem(arg) {
		return "нет такого"
	}

	p.Backpack = true
	p.CurrentRoom.RoomItems.DeleteItem(Item(arg))

	return "вы надели: " + arg
}

func (p *Player) Apply(arg []string) string {
	if !p.Inventory.HasItem(arg[0]) {
		return "нет предмета в инвентаре - " + arg[0]
	}
	if !gameWorld.InteractionItem.IsApplicationAllowed(arg) {
		return "не к чему применить"
	}
	p.CurrentRoom.OpenDoor()

	return "дверь открыта"
}

func (p *Player) LookAround(cmd string, bacpack bool) string {
	return p.CurrentRoom.GetDescription(cmd, bacpack)
}

func NewPlayer(room *Room) *Player {
	return &Player{
		Name:        "Viktor",
		Backpack:    false,
		Inventory:   make(Inventory),
		CurrentRoom: room,
	}
}

// Инвентарь игрока------------------------------------------------------------
type Inventory map[string]Item

func (i Inventory) Add(name Item) {
	i[string(name)] = name
}

func (i Inventory) HasItem(nameItem string) bool {
	_, err := i[nameItem]
	return err
}

func (i Inventory) DeleteItem(name string) {
	delete(i, name)
}

func (i Inventory) IsEmpty() bool {
	if len(i) == 0 {
		return true
	}
	return false
}

// Взаимодействие предметов в мире--------------------------------------------------
type InteractionItem map[string]string

func NewInetactionItem() *InteractionItem {
	i := make(InteractionItem)
	return &i
}

func (i InteractionItem) Add(item1, item2 string) {
	i[item1] = item2
}

func (i InteractionItem) IsApplicationAllowed(arg []string) bool {
	if i[arg[0]] == arg[1] {
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

	kitchenRoom := NewRoom("кухня",
		NewRoomItems(map[string][]Item{
			"стол": {"чай"},
		}),
		NewTusks(map[string]bool{
			", надо собрать рюкзак и идти в универ": true,
			", надо идти в универ":                  false,
		}),
		NewDescription(map[string]string{
			"осмотреться": "ты находишься на кухне, ",
			"идти":        "кухня, ничего интересного",
		}),
		false)

	corridorRoom := NewRoom("коридор",
		NewRoomItems(make(map[string][]Item)),
		NewTusks(make(map[string]bool)),
		NewDescription(map[string]string{
			"осмотреться": "ничего интересного",
			"идти":        "ничего интересного",
		}),
		false)

	myRoom := NewRoom("комната",
		NewRoomItems(map[string][]Item{
			"стул": {"рюкзак"},
			"стол": {"ключи", "конспекты"},
		}),
		NewTusks(make(map[string]bool)),
		NewDescription(map[string]string{
			"идти": "ты в своей комнате",
		}),
		false)

	street := NewRoom("улица",
		NewRoomItems(make(map[string][]Item)),
		NewTusks(make(map[string]bool)),
		NewDescription(map[string]string{
			"идти": "на улице весна"}),
		true)

	home := NewRoom("домой",
		NewRoomItems(make(map[string][]Item)),
		NewTusks(make(map[string]bool)),
		NewDescription(map[string]string{
			"идти": "на улице весна"}),
		false)

	kitchenRoom.SetRoomAllowed(corridorRoom)
	corridorRoom.SetRoomAllowed(kitchenRoom, myRoom, street)
	myRoom.SetRoomAllowed(corridorRoom)
	street.SetRoomAllowed(home)

	Rooms := make([]*Room, 0, 6)
	Rooms = append(Rooms, kitchenRoom, corridorRoom, myRoom, street, home)

	return Rooms
}

// Инициализация необходимых объектов-------------------------------------------
func initGame() {
	rooms := initRooms()
	defaultRoom, err := rooms.GetDefaultRoom()
	if err != nil {
		panic("Не определена комната по умолчанию")
	}
	player := NewPlayer(defaultRoom)

	interactionItemWorld := NewInetactionItem()
	interactionItemWorld.Add("ключи", "дверь")
	gameWorld = NewWorld(rooms, player, *interactionItemWorld)

}

// ОбработкаДействия-------------------------------------------------------------
func handleCommand(cmd string) string {
	str := strings.Split(cmd, " ")
	return gameWorld.Processing(str[0], str[1:]...)
}

func RemoveLastChar(s string) string {
	runes := []rune(s)

	if len(runes) < 2 {
		return s
	}

	runes = runes[:len(runes)-2]

	return string(runes)
}

func main() {
	casesTest := NewGameCaseTest()
	initGame()

	for _, c := range casesTest {
		for _, el := range c {
			fmt.Println(handleCommand(el))
		}
	}
}
