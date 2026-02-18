// package main

// import (
// 	"errors"
// 	"fmt"
// 	"os"
// )

// // label - уникальное наименование
// type label string

// // command - команда, которую можно выполнять в игре
// type command label

// // список доступных команд
// var (
// 	eat  = command("eat")
// 	take = command("take")
// 	talk = command("talk to")
// )

// // thing - объект, который существует в игре
// type thing struct {
// 	name    label
// 	actions map[command]string
// }

// // supports() возвращает true, если объект
// // поддерживает команду action
// func (t thing) supports(action command) bool {
// 	_, ok := t.actions[action]
// 	return ok
// }

// // String() возвращает описание объекта
// func (t thing) String() string {
// 	return string(t.name)
// }

// // полный список объектов в игре
// var (
// 	apple = thing{"apple", map[command]string{
// 		eat:  "mmm, delicious!",
// 		take: "you have an apple now",
// 	}}
// 	bob = thing{"bob", map[command]string{
// 		talk: "Bob says hello",
// 	}}
// 	coin = thing{"coin", map[command]string{
// 		take: "you have a coin now",
// 	}}
// 	mirror = thing{"mirror", map[command]string{
// 		take: "you have a mirror now",
// 		talk: "mirror does not answer",
// 	}}
// 	mushroom = thing{"mushroom", map[command]string{
// 		eat:  "tastes funny",
// 		take: "you have a mushroom now",
// 	}}
// )

// // step описывает шаг игры: сочетание команды и объекта
// type step struct {
// 	cmd command
// 	obj thing
// }

// // isValid() возвращает true, если объект
// // совместим с командой
// func (s step) isValid() bool {
// 	return s.obj.supports(s.cmd)
// }

// // String() возвращает описание шага
// func (s step) String() string {
// 	return fmt.Sprintf("%s %s", s.cmd, s.obj)
// }

// // начало решения

// // invalidStepError - ошибка, которая возникает,
// // когда команда шага не совместима с объектом
// type invalidStepError struct {
// 	st step
// }

// func (err invalidStepError) Error() string {
// 	return fmt.Sprintf("things like '%s' are impossible", err.st)
// }

// // notEnoughObjectsError - ошибка, которая возникает,
// // когда в игре закончились объекты определенного типа
// type notEnoughObjectsError struct {
// 	obj thing
// }

// func (err notEnoughObjectsError) Error() string {
// 	return fmt.Sprintf("be careful with scarce %ss", err.obj)
// }

// // commandLimitExceededError - ошибка, которая возникает,
// // когда игрок превысил лимит на выполнение команды
// type commandLimitExceededError struct {
// 	cmd command
// }

// func (err commandLimitExceededError) Error() string {
// 	switch err.cmd {
// 	case "eat": return "eat less"
// 	case "talk": return "talk to less"
// 	}
// 	return ""
// }

// // objectLimitExceededError - ошибка, которая возникает,
// // когда игрок превысил лимит на количество объектов
// // определенного типа в инвентаре
// type objectLimitExceededError struct {
// 	obj thing
// }

// func (err objectLimitExceededError) Error() string {
//     return fmt.Sprintf("you already have a %s", err.obj)
// }

// // gameOverError - ошибка, которая произошла в игре
// type gameOverError struct {
// 	// количество шагов, успешно выполненных
// 	// до того, как произошла ошибка
// 	nSteps int
// 	err error
// }

// func (err gameOverError) Error() string {
// 	return err.err.Error()
// }

// func (err gameOverError) Unwrap() error {
// 	return err.err
// }

// // player - игрок
// type player struct {
// 	// количество съеденного
// 	nEaten int
// 	// количество диалогов
// 	nDialogs int
// 	// инвентарь
// 	inventory []thing
// }

// // has() возвращает true, если у игрока
// // в инвентаре есть предмет obj
// func (p *player) has(obj thing) bool {
// 	for _, got := range p.inventory {
// 		if got.name == obj.name {
// 			return true
// 		}
// 	}
// 	return false
// }

// // do() выполняет команду cmd над объектом obj
// // от имени игрока
// func (p *player) do(cmd command, obj thing) error {
// 	// действуем в соответствии с командой
// 	switch cmd {
// 	case eat:
// 		if p.nEaten > 1 {
// 			return commandLimitExceededError{cmd}
// 		}
// 		p.nEaten++
// 	case take:
// 		if p.has(obj) {
// 			return objectLimitExceededError{obj}
// 		}
// 		p.inventory = append(p.inventory, obj)
// 	case talk:
// 		if p.nDialogs > 0 {
// 			return commandLimitExceededError{cmd}
// 		}
// 		p.nDialogs++
// 	}
// 	return nil
// }

// // newPlayer создает нового игрока
// func newPlayer() *player {
// 	return &player{inventory: []thing{}}
// }

// // game описывает игру
// type game struct {
// 	// игрок
// 	player *player
// 	// объекты игрового мира
// 	things map[label]int
// 	// количество успешно выполненных шагов
// 	nSteps int
// }

// // has() проверяет, остались ли в игровом мире указанные предметы
// func (g *game) has(obj thing) bool {
// 	count := g.things[obj.name]
// 	return count > 0
// }

// // execute() выполняет шаг step
// func (g *game) execute(st step) error {
// 	// проверяем совместимость команды и объекта
// 	if !st.isValid() {
// 		return gameOverError{g.nSteps, invalidStepError{st}}
// 	}

// 	// когда игрок берет или съедает предмет,
// 	// тот пропадает из игрового мира
// 	if st.cmd == take || st.cmd == eat {
// 		if !g.has(st.obj) {
// 			return gameOverError{g.nSteps, notEnoughObjectsError{st.obj}}
// 		}
// 		g.things[st.obj.name]--
// 	}

// 	// выполняем команду от имени игрока
// 	if err := g.player.do(st.cmd, st.obj); err != nil {
// 		return gameOverError{g.nSteps, err}
// 	}

// 	g.nSteps++
// 	return nil
// }

// // newGame() создает новую игру
// func newGame() *game {
// 	p := newPlayer()
// 	things := map[label]int{
// 		apple.name:    2,
// 		coin.name:     3,
// 		mirror.name:   1,
// 		mushroom.name: 1,
// 	}
// 	return &game{p, things, 0}
// }

// // giveAdvice() возвращает совет, который
// // поможет игроку избежать ошибки err в будущем
// func giveAdvice(err error) string {
// 	var advice string

// 	var gameErr gameOverError
// 	if errors.As(err, &gameErr) {
// 		err = gameErr.err
// 	}

// 	switch e := err.(type) {
// 	case invalidStepError:
// 		advice = fmt.Sprintf("things like '%s' are impossible", e.st)
// 	case notEnoughObjectsError:
//         advice = fmt.Sprintf("be careful with scarce %ss", e.obj)
// 	case commandLimitExceededError:
//         switch e.cmd {
//         case eat:
//             advice = "eat less"
//         case talk:
//             advice = "talk to less"
//         default:
//             advice = "command limit exceeded"
//         }
// 	case objectLimitExceededError:
//         advice = fmt.Sprintf("don't be greedy, 1 %s is enough", e.obj)
//     default:
//         advice = ""
// 	}

// 	return advice
// }

// // конец решения

// func main() {
// 	m := ZipMap([]string{"one"}, []int{11, 22, 33})
// 	gm := newGame()
// 	steps := []step{
// 		{eat, apple},
// 		{talk, bob},
// 		{take, coin},
// 		{eat, mushroom},
// 	}

// 	for _, st := range steps {
// 		if err := tryStep(gm, st); err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
// 	}
// 	fmt.Println(m)
// }

// // tryStep() выполняет шаг игры и печатает результат
// func tryStep(gm *game, st step) error {
// 	fmt.Printf("trying to %s %s... ", st.cmd, st.obj.name)
// 	if err := gm.execute(st); err != nil {
// 		fmt.Println("FAIL")
// 		return err
// 	}
// 	fmt.Println("OK")
// 	return nil
// }

// // начало решения

// // Produce возвращает срез из n значений val.
// func Produce[T any] (val T, n int) []T {
//     vals := make([]T, n)
//     for i := range n {
//         vals[i] = val
//     }
//     return vals
// }

// // конец решения
// // начало решения

// // ZipMap возвращает карту, где ключи - элементы из keys, а значения - из vals.
// func ZipMap[K comparable, V any](keys []K, vals []V) map[K]V {
//     minLen := min(len(keys), len(vals))
// 	mergedMap := make(map[K]V, minLen)

// 	for i := 0; i < minLen; i++ {
// 		mergedMap[keys[i]] = vals[i]
// 	}

// 	return mergedMap
// }

// // конец решения

// // начало решения

// // Avg - накопительное среднее значение.
// type Avg[T int|float64] struct {
// 	cnt int
// 	sum T
// }

// // Add пересчитывает среднее значение с учетом val.
// func (a *Avg[T]) Add(val T) *Avg[T] {
//     a.cnt++
// 	a.sum += val
// 	return a
// }

// // Val возвращает текущее среднее значение.
// func (a *Avg[T]) Val() T {
// 	if a.cnt > 0 {
//     	return a.sum/T(a.cnt)
// 	}
// 	return T(0)
// }

// // конец решения

// // начало решения

// // Map - карта "ключ-значение".
// type Map[K comparable, V any] map[K]V

// // Set устанавливает значение для ключа.
// func (m *Map[K, V]) Set(key K, val V) {
// 		(*m)[key] = val
// }

// // Get возвращает значение по ключу.
// func (m Map[K, V]) Get(key K) V {
//     if val, ok := m[key]; ok {
// 		return val
// 	}
// 	var zero V
//     return zero
// }

// // Keys возвращает срез ключей карты.
// // Порядок ключей неважен, и не обязан совпадать
// // с порядком значений из метода Values.
// func (m Map[K, V]) Keys() []K {
//     slice := []K{}
// 	for k := range m {
// 		slice = append(slice, k)
// 	}
// 	return slice
// }

// // Values возвращает срез значений карты.
// // Порядок значений неважен, и не обязан совпадать
// // с порядком ключей из метода Keys.
// func (m Map[K, V]) Values() []V {
//     slice := []V{}
// 	for _, v := range m {
// 		slice = append(slice, v)
// 	}
// 	return slice
// }

// // конец решения

// package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"testing"
// 	"time"
// )

// // WeatherService предсказывает погоду.
// type WeatherService struct{}

// // Forecast сообщает ожидаемую дневную температуру на завтра.
// func (ws *WeatherService) Forecast() int {
// 	rand.Seed(time.Now().Unix())
// 	value := rand.Intn(31)
// 	sign := rand.Intn(2)
// 	if sign == 1 {
// 		value = -value
// 	}
// 	return value
// }

// // Weather выдает текстовый прогноз погоды.
// type Weather struct {
// 	service WeatherSERVICE
// }

// // Forecast сообщает текстовый прогноз погоды на завтра.
// func (w Weather) Forecast() string {
// 	deg := w.service.Forecast()
// 	switch {
// 	case deg < 10:
// 		return "холодно"
// 	case deg >= 10 && deg < 15:
// 		return "прохладно"
// 	case deg >= 15 && deg < 20:
// 		return "идеально"
// 	case deg >= 20:
// 		return "жарко"
// 	}
// 	return "инопланетно"
// }

// type testCase struct {
// 	deg  int
// 	want string
// }

// var tests []testCase = []testCase{
// 	{-10, "холодно"},
// 	{0, "холодно"},
// 	{5, "холодно"},
// 	{10, "прохладно"},
// 	{15, "идеально"},
// 	{20, "жарко"},
// }

// type WeatherSERVICE interface {
// 	Forecast() int
// }

// type MockService struct {
// 	deg int
// }

// func (m MockService) Forecast() int {
// 	return m.deg
// }

// func TestForecast(t *testing.T) {
// 	service := &MockService{}
// 	weather := Weather{service}
// 	for _, test := range tests {
// 		service.deg = test.deg
// 		name := fmt.Sprintf("%v", test.deg)
// 		t.Run(name, func(t *testing.T) {
// 			got := weather.Forecast()
// 			if got != test.want {
// 				t.Errorf("%s: got %s, want %s", name, got, test.want)
// 			}
// 		})
// 	}
//
//

// package main

// // не удаляйте импорты, они используются при проверке
// import (
//     "fmt"
//     "math/rand"
//     "os"
//     "testing"
// )

// // IntSet реализует множество целых чисел
// // (элементы множества уникальны).
// type IntSet struct {
//     elems map[int]any
// }

// // MakeIntSet создает пустое множество.
// func MakeIntSet() IntSet {
// 	elems := map[int]any{}
//     return IntSet{elems}
// }

// // Contains проверяет, содержится ли элемент в множестве.
// func (s IntSet) Contains(elem int) bool {
//     if _, ok := s.elems[elem]; ok {
// 		return true
// 	}
// 	return false
// }

// // Add добавляет элемент в множество.
// // Возвращает true, если элемент добавлен,
// // иначе false (если элемент уже содержится в множестве).
// func (s *IntSet) Add(elem int) bool {
// 	if s.Contains(elem) {
// 		return false
// 	}
// 	s.elems[elem] = struct{}{}
//     return true
// }

// package main

// // не удаляйте импорты, они используются при проверке
// import (
//     "fmt"
//     "math/rand"
//     "os"
//     "strings"
//     "testing"
// )

// // Words работает со словами в строке.
// type Words struct {
//     str string
//     words map[string]int
// }

// // MakeWords создает новый экземпляр Words.
// func MakeWords(s string) Words {
// 	size := len(s)
// 	if len(s) > 10000 {
// 		size = 10000
// 	}
// 	newMap := make(map[string]int, size)
// 	words := strings.Fields(s)

// 	for idx, word := range words {
// 		if _, exists := newMap[word]; !exists {
//             newMap[word] = idx
//         }
// 	}
// 	return Words{s, newMap}
// }

// // Index возвращает индекс первого вхождения слова в строке,
// // или -1, если слово не найдено.
// func (w Words) Index(word string) int {
//     if idx, ok := w.words[word]; ok {
// 		return idx
// 	}
//     return -1
// }
// type counter map[string]int
// func countDigitsInWords(phrase string) counter {
//     words := strings.Fields(phrase)
//     counted := make(chan int)
//     var stats counter

//     go func() {
//         for _, word := range words {
//             counted <- countDigits(word)
// 		}
//         // Пройдите по словам,
//         // посчитайте количество цифр в каждом,
//         // и запишите его в канал counted
//     	}()
//         for _, word := range words {
// 			value := <- counted
//             stats[word] = value
//         }
//     // Считайте значения из канала counted
//     // и заполните stats.

//     // В результате stats должна содержать слова
//     // и количество цифр в каждом.

//     return stats
// }

// Канал с результатами.
// package main

// import (
// 	"fmt"
// 	"strings"
// 	"unicode"
// )

// counter хранит количество цифр в каждом слове.
// Ключ карты - слово, а значение - количество цифр в слове.
// type counter map[string]int

// // начало решения

// // countDigitsInWords считает количество цифр в словах фразы.
// func countDigitsInWords(phrase string) counter {
//     words := strings.Fields(phrase)
//     counted := make(chan int)
//     stats := make(counter)

//     go func() {
//         for _, word := range words {
//             counted <- countDigits(word)
//         }
//         // Пройдите по словам,
//         // посчитайте количество цифр в каждом,
//         // и запишите его в канал counted
//     }()
//         for _, word := range words {
//             value := <- counted
//             stats[word] = value
//         }
//     // Считайте значения из канала counted
//     // и заполните stats.

//     // В результате stats должна содержать слова
//     // и количество цифр в каждом.

//     return stats
// }

// // конец решения
// // countDigits возвращает количество цифр в строке.
// func countDigits(str string) int {
// 	count := 0
// 	for _, char := range str {
// 		if unicode.IsDigit(char) {
// 			count++
// 		}
// 	}
// 	return count
// }

// // printStats печатает количество цифр в словах.
// func printStats(stats counter) {
// 	for word, count := range stats {
// 		fmt.Printf("%s: %d\n", word, count)
// 	}
// }

// func main() {
// 	phrase := "0ne 1wo thr33 4068"
// 	stats := countDigitsInWords(phrase)
// 	printStats(stats)
// }

// Выборка из генератора.
// package main

// import (
//     "fmt"
//     "strings"
//     "unicode"
// )

// // nextFunc возвращает следующее слово из генератора.
// type nextFunc func() string

// // counter хранит количество цифр в каждом слове.
// // Ключ карты - слово, а значение - количество цифр в слове.
// type counter map[string]int

// // pair хранит слово и количество цифр в нем.
// type pair struct {
//     word  string
//     count int
// }

// // начало решения

// // countDigitsInWords считает количество цифр в словах,
// // выбирая очередные слова с помощью next().
// func countDigitsInWords(next nextFunc) counter {
//     counted := make(chan pair)
// 	stats := make(counter)

//     go func() {
// 		for {
// 			word := next()
// 			counted <- pair{word, countDigits(word)}
// 			if word == "" {
// 				break
// 			}
// 		}
//         // Пройдите по словам,
//         // посчитайте количество цифр в каждом,
//         // и запишите его в канал counted
//     }()

// 	for {
// 		// values, ok := <- counted
// 		// if !ok {
// 		// 	break
// 		// }
// 		// stats[values.word] = values.count
// 		p := <-counted
//         if p.word == "" {
//             break
//         }
//         stats[p.word] = p.count
// 	}
//     // Считайте значения из канала counted
//     // и заполните stats.

//     // В результате stats должна содержать слова
//     // и количество цифр в каждом.

//     return stats
// }

// // конец решения

// // countDigits возвращает количество цифр в строке.
// func countDigits(str string) int {
//     count := 0
//     for _, char := range str {
//         if unicode.IsDigit(char) {
//             count++
//         }
//     }
//     return count
// }

// // printStats печатает количество цифр в словах.
// func printStats(stats counter) {
//     for word, count := range stats {
//         fmt.Printf("%s: %d\n", word, count)
//     }
// }

// // wordGenerator возвращает генератор,
// // который выдает слова из фразы.
// func wordGenerator(phrase string) nextFunc {
//     words := strings.Fields(phrase)
//     idx := 0
//     return func() string {
//         if idx == len(words) {
//             return ""
//         }
//         word := words[idx]
//         idx++
//         return word
//     }
// }

// func main() {
//     phrase := "0ne 1wo thr33 4068"
//     next := wordGenerator(phrase)
//     stats := countDigitsInWords(next)
//     printStats(stats)
// }

// Читатель и счетовод.
// package main

// import (
// 	"fmt"
// 	"strings"
// 	"unicode"
// )

// // nextFunc возвращает следующее слово из генератора.
// type nextFunc func() string

// // counter хранит количество цифр в каждом слове.
// // Ключ карты - слово, а значение - количество цифр в слове.
// type counter map[string]int

// // pair хранит слово и количество цифр в нем.
// type pair struct {
// 	word  string
// 	count int
// }

// // начало решения

// // countDigitsInWords считает количество цифр в словах,
// // выбирая очередные слова с помощью next().
// func countDigitsInWords(next nextFunc) counter {
// 	pending := make(chan string)
// 	counted := make(chan pair)
// 	stats := make(counter)

// 	// отправляет слова на подсчет
// 	go func() {
// 		for {
// 			word := next()
// 			pending <- word
// 			if word == "" {
// 				break
// 			}
// 		}
// 		// Пройдите по словам и отправьте их
// 		// в канал pending
// 	}()

// 	// считает цифры в словах
// 	go func() {
// 		for {
// 			word := <- pending
// 			counted <- pair{word, countDigits(word)}
// 			if word == "" {
// 				break
// 			}
// 		}

// 		// Считайте слова из канала pending,
// 		// посчитайте количество цифр в каждом,
// 		// и запишите его в канал counted
// 	}()

// 	for {
// 		p := <- counted
// 		if p.word == "" {
// 			break
// 		}
// 		stats[p.word] = p.count
// 	}
// 	// Считайте значения из канала counted
// 	// и заполните stats.

// 	// В результате stats должна содержать слова
// 	// и количество цифр в каждом.

// 	return stats
// }

// // конец решения

// // countDigits возвращает количество цифр в строке.
// func countDigits(str string) int {
// 	count := 0
// 	for _, char := range str {
// 		if unicode.IsDigit(char) {
// 			count++
// 		}
// 	}
// 	return count
// }

// // printStats печатает количество цифр в словах.
// func printStats(stats counter) {
// 	for word, count := range stats {
// 		fmt.Printf("%s: %d\n", word, count)
// 	}
// }

// // wordGenerator возвращает генератор,
// // который выдает слова из фразы.
// func wordGenerator(phrase string) nextFunc {
// 	words := strings.Fields(phrase)
// 	idx := 0
// 	return func() string {
// 		if idx == len(words) {
// 			return ""
// 		}
// 		word := words[idx]
// 		idx++
// 		return word
// 	}
// }

// func main() {
// 	phrase := "0ne 1wo thr33 4068"
// 	next := wordGenerator(phrase)
// 	stats := countDigitsInWords(next)
// 	printStats(stats)
// }

// package main

// import (
// 	"fmt"
// 	"strings"
// 	"unicode"
// )

// // nextFunc возвращает следующее слово из генератора
// type nextFunc func() string

// // counter хранит количество цифр в каждом слове.
// // ключ карты - слово, а значение - количество цифр в слове.
// type counter map[string]int

// // pair хранит слово и количество цифр в нем
// type pair struct {
// 	word  string
// 	count int
// }

// // countDigitsInWords считает количество цифр в словах,
// // выбирая очередные слова с помощью next()
// func countDigitsInWords(next nextFunc) counter {
// 	pending := make(chan string)
// 	go submitWords(next, pending)

// 	counted := make(chan pair)
// 	go countWords(pending, counted)

// 	return fillStats(counted)
// }

// // начало решения

// // submitWords отправляет слова на подсчет
// func submitWords(next func() string, pending chan string) {
// 	for {
// 		word := next()
// 		pending <- word
// 		if word == "" {
// 			break
// 		}
// 	}

// }
// // countWords считает цифры в словах
// func countWords(pending chan string, counted chan pair) {
// 	for {
// 		word := <- pending
// 		counted <- pair{word, countDigits(word)}
// 		if word == "" {
// 			break
// 		}
// 	}
// }
// // fillStats готовит итоговую статистику
// func fillStats(counted chan pair) counter {
// 	stats := make(counter)
// 	for {
// 		p := <- counted
// 		if p.word == "" {
// 			break
// 		}
// 		stats[p.word] = p.count
// 	}
// 	return stats
// }
// // конец решения

// // countDigits возвращает количество цифр в строке
// func countDigits(str string) int {
// 	count := 0
// 	for _, char := range str {
// 		if unicode.IsDigit(char) {
// 			count++
// 		}
// 	}
// 	return count
// }

// // printStats печатает слова и количество цифр в каждом
// func printStats(stats counter) {
// 	for word, count := range stats {
// 		fmt.Printf("%s: %d\n", word, count)
// 	}
// }

// // wordGenerator возвращает генератор, который выдает слова из фразы
// func wordGenerator(phrase string) nextFunc {
// 	words := strings.Fields(phrase)
// 	idx := 0
// 	return func() string {
// 		if idx == len(words) {
// 			return ""
// 		}
// 		word := words[idx]
// 		idx++
// 		return word
// 	}
// }

// func main() {
// 	phrase := "0ne 1wo thr33 4068"
// 	next := wordGenerator(phrase)
// 	stats := countDigitsInWords(next)
// 	printStats(stats)
// }

// package main

// import (
// 	"fmt"
// 	"strings"
// 	"unicode"
// )

// // nextFunc возвращает следующее слово из генератора
// type nextFunc func() string

// // counter хранит количество цифр в каждом слове.
// // ключ карты - слово, а значение - количество цифр в слове.
// type counter map[string]int

// // pair хранит слово и количество цифр в нем
// type pair struct {
// 	word  string
// 	count int
// }

// // countDigitsInWords считает количество цифр в словах,
// // выбирая очередные слова с помощью next()
// func countDigitsInWords(next nextFunc) counter {
// 	pending := make(chan string)
// 	go submitWords(next, pending)

// 	counted := make(chan pair)
// 	go countWords(pending, counted)

// 	return fillStats(counted)
// }

// // начало решения

// // submitWords отправляет слова на подсчет
// func submitWords(next nextFunc, out chan string) {
// 	for {
// 		defer close(out)
// 		word := next()
//         if word == "" {
//             break
//         }
// 		out <- word

// 	}
// }

// // countWords считает цифры в словах
// func countWords(in chan string, out chan pair) {
// 	for word := range in {
// 		defer close(out)
// 		out <- pair{word, countDigits(word)}
// 	}
// }

// // fillStats готовит итоговую статистику
// func fillStats(in chan pair) counter {
// 	stats := make(counter)
// 	for stat := range in {
// 		stats[stat.word] = stat.count
// 	}
// 	return stats
// }

// // конец решения

// // countDigits возвращает количество цифр в строке
// func countDigits(str string) int {
// 	count := 0
// 	for _, char := range str {
// 		if unicode.IsDigit(char) {
// 			count++
// 		}
// 	}
// 	return count
// }

// // printStats печатает слова и количество цифр в каждом
// func printStats(stats counter) {
// 	for word, count := range stats {
// 		fmt.Printf("%s: %d\n", word, count)
// 	}
// }

// // wordGenerator возвращает генератор, который выдает слова из фразы
// func wordGenerator(phrase string) nextFunc {
// 	words := strings.Fields(phrase)
// 	idx := 0
// 	return func() string {
// 		if idx == len(words) {
// 			return ""
// 		}
// 		word := words[idx]
// 		idx++
// 		return word
// 	}
// }

// func main() {
// 	phrase := "0ne 1wo thr33 4068"
// 	next := wordGenerator(phrase)
// 	stats := countDigitsInWords(next)
// 	printStats(stats)
// }

// package main

// import (
// 	"fmt"
// 	"strings"
// )

// // encode кодирует строку шифром Цезаря
// func encode(str string) string {
// 	// начало решения

// 	submitter := func(str string) <-chan string {
// 		ch := make(chan string)
// 		go func() {
// 			words := strings.Fields(str)
// 			for _, word := range words {
// 				ch <- word
// 			}
// 			close(ch)
// 		}()
// 		return ch
// 	}

// 	encoder := func(ch1 <-chan string) <-chan string {
// 		ch2 := make(chan string)
// 		go func() {
// 			for word := range ch1 {
// 				ch2 <- encodeWord(word)
// 			}
// 			close(ch2)
// 		}()
// 		return ch2
// 	}

// 	receiver := func(ch <-chan string) []string {
// 		words := []string{}
// 		for word := range ch {
// 			words = append(words, word)
// 		}
// 		return words
// 	}

// 	// конец решения

// 	pending := submitter(str)
// 	encoded := encoder(pending)
// 	words := receiver(encoded)
// 	return strings.Join(words, " ")
// }

// // encodeWord кодирует слово шифром Цезаря
// func encodeWord(word string) string {
// 	const shift = 13
// 	const char_a byte = 'a'
// 	encoded := make([]byte, len(word))
// 	for idx, char := range []byte(word) {
// 		delta := (char - char_a + shift) % 26
// 		encoded[idx] = char_a + delta
// 	}
// 	return string(encoded)
// }

// func main() {
// 	src := "go is awesome"
// 	res := encode(src)
// 	fmt.Println(res)
// }

// Четыре счетовода.
// package main

// import (
// 	"fmt"
// 	"strings"
// 	"unicode"
// )

// // nextFunc возвращает следующее слово из генератора
// type nextFunc func() string

// // counter хранит количество цифр в каждом слове.
// // Ключ карты - слово, а значение - количество цифр в слове.
// type counter map[string]int

// // pair хранит слово и количество цифр в нем
// type pair struct {
// 	word  string
// 	count int
// }

// // countDigitsInWords считает количество цифр в словах,
// // выбирая очередные слова с помощью next()
// func countDigitsInWords(next nextFunc) counter {
// 	pending := make(chan string)
// 	go submitWords(next, pending)

// 	done := make(chan struct{})
// 	counted := make(chan pair)

// 	// начало решения
// 	N := 4
// 	for range(N) {
// 		go func() {
// 			countWords(done, pending, counted)
// 		}()
// 	}
// 	go func() {
// 		for range(N) {
// 			<-done
// 		}
// 		close(counted)
// 	}()
// 	// Запустите четыре горутины countWords()
// 	// вместо одной.

// 	// Используйте канал завершения, чтобы дождаться
// 	// окончания обработки и закрыть канал counted.

// 	// конец решения

// 	return fillStats(counted)
// }

// // submitWords отправляет слова на подсчет
// func submitWords(next nextFunc, out chan<- string) {
// 	for {
// 		word := next()
// 		if word == "" {
// 			break
// 		}
// 		out <- word
// 	}
// 	close(out)
// }

// // countWords считает цифры в словах
// func countWords(done chan<- struct{}, in <-chan string, out chan<- pair) {
// 	for word := range in {
// 		out <- pair{word, countDigits(word)}
// 	}
// 	done <- struct{}{}
// }

// // fillStats готовит итоговую статистику
// func fillStats(in <-chan pair) counter {
// 	stats := counter{}
// 	for p := range in {
// 		stats[p.word] = p.count
// 	}
// 	return stats
// }

// // countDigits возвращает количество цифр в строке
// func countDigits(str string) int {
// 	count := 0
// 	for _, char := range str {
// 		if unicode.IsDigit(char) {
// 			count++
// 		}
// 	}
// 	return count
// }

// // printStats печатает количество цифр в словах
// func printStats(stats counter) {
// 	for word, count := range stats {
// 		fmt.Printf("%s: %d\n", word, count)
// 	}
// }

// // wordGenerator возвращает генератор, который выдает слова из фразы
// func wordGenerator(phrase string) nextFunc {
// 	words := strings.Fields(phrase)
// 	idx := 0
// 	return func() string {
// 		if idx == len(words) {
// 			return ""
// 		}
// 		word := words[idx]
// 		idx++
// 		return word
// 	}
// }

// func main() {
// 	phrase := "1 22 333 4444 55555 666666 7777777 88888888"
// 	next := wordGenerator(phrase)
// 	stats := countDigitsInWords(next)
// 	printStats(stats)
// }

// Promise.all()
// package main

// import (
// 	"fmt"
// 	"time"
// )

// // начало решения

// // gather выполняет переданные функции одновременно
// // и возвращает срез с результатами, когда они готовы
// func gather(funcs []func() any) []any {
// 	// Выполните все переданные функции,
// 	// соберите результаты в срез и верните его.
// 	results := make([]any, len(funcs))
// 	done := make(chan struct{
// 		index int
// 		value any
// 	})
// 	for idx, function := range funcs {
// 		go func() {
// 			done <- struct{index int; value any}{idx, function()}
// 		}()
// 	}
// 	for range(len(funcs)) {
// 		result := <-done
// 		results[result.index] = result.value
// 	}
// 	return results
// }

// // конец решения

// // squared возвращает функцию,
// // которая считает квадрат n
// func squared(n int) func() any {
// 	return func() any {
// 		time.Sleep(time.Duration(n) * 100 * time.Millisecond)
// 		return n * n
// 	}
// }

// func main() {
// 	funcs := []func() any{squared(1), squared(2), squared(3), squared(4), squared(5)}

// 	start := time.Now()
// 	nums := gather(funcs)
// 	elapsed := float64(time.Since(start)) / 1_000_000

// 	fmt.Println(nums)
// 	fmt.Printf("Took %.0f ms\n", elapsed)
// }

// package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"strings"
// 	"time"
// )

// // say печатает фразу от имени обработчика
// func say(id int, phrase string) {
// 	for _, word := range strings.Fields(phrase) {
// 		fmt.Printf("Worker #%d says: %s...\n", id, word)
// 		dur := time.Duration(rand.Intn(100)) * time.Millisecond
// 		time.Sleep(dur)
// 	}
// }

// // начало решения

// // makePool создает пул на n обработчиков
// // возвращает функции handle и wait
// func makePool(n int, handler func(int, string)) (func(string), func()) {
// 	pool := make(chan int, n)
// 	done := make(chan struct{})

// 	for i := range n {
// 		pool <- i
// 	}

// 	var handleCount int

// 	handle := func(s string) {
// 		handleCount++
// 		go func() {
// 			idx := <-pool
// 			handler(idx, s)
// 			pool <- idx
// 			done <- struct{}{}
// 		}()
// 	}

// 	wait := func() {
// 		for range handleCount {
// 			<-done
// 		}
// 	}

// 	return handle, wait
// }

// // конец решения

// func main() {
// 	phrases := []string{
// 		"go is awesome",
// 		"cats are cute",
// 		"rain is wet",
// 		"channels are hard",
// 		"floor is lava",
// 	}

// 	handle, wait := makePool(2, say)
// 	for _, phrase := range phrases {
// 		handle(phrase)
// 	}
// 	wait()
// }

// package main

// import (
// 	"fmt"
// )

// // начало решения

// // count отправляет в канал числа от start до бесконечности
// func count(cancel chan struct{}, start int) <-chan int {
// 	out := make(chan int)
// 	go func() {
// 		defer close(out)
// 		for i := start; ; i++ {
// 			select {
// 			case out <- i:
// 			case <-cancel: {
// 				return
// 			}
// 			}
// 		}
// 	}()
// 	return out
// }

// // take выбирает первые n чисел из in и отправляет в выходной канал
// func take(cancel chan struct{}, in <-chan int, n int) <-chan int {
// 	out := make(chan int)
// 	go func() {
// 		defer close(out)
// 		for i := 0; i < n; i++ {
// 			select {
// 			case val, ok := <-in: if !ok {return} else {out <- val}
// 			case <-cancel: {
// 				return
// 			}
// 			}
// 		}
// 	}()
// 	return out
// }

// // конец решения

// func main() {
// 	cancel := make(chan struct{})
// 	defer close(cancel)

// 	stream := take(cancel, count(cancel, 10), 5)
// 	first := <-stream
// 	second := <-stream
// 	third := <-stream

// 	fmt.Println(first, second, third)
// }

// Promise.all()
// package main

// import (
// 	"fmt"
// 	"time"
// )

// // // начало решения
// // type result = struct {
// // 		value any
// // 		idx int
// // 	}
// // // gather выполняет переданные функции одновременно
// // // и возвращает срез с результатами, когда они готовы
// // func gather(funcs []func() any) []any {
// // 	results := make([]any, len(funcs))
// // 	done := make(chan result)

// // 	for idx, function := range funcs {
// // 		go func() {
// // 			done <- result{function, idx}
// // 		}()
// // 	}

// // 	for range(len(funcs)) {
// // 		res := <- done
// // 		results[res.idx] = res.value
// // 	}
// // 	return results
// // 	// Выполните все переданные функции,
// // 	// соберите результаты в срез и верните его.
// // }

// // // конец решения

// // squared возвращает функцию,
// // которая считает квадрат n
// func squared(n int) func() any {
// 	return func() any {
// 		time.Sleep(time.Duration(n) * 100 * time.Millisecond)
// 		return n * n
// 	}
// }

// func main() {
// 	funcs := []func() any{squared(2), squared(3), squared(4)}

// 	start := time.Now()
// 	nums := gather(funcs)
// 	elapsed := float64(time.Since(start)) / 1_000_000

// 	fmt.Println(nums)
// 	fmt.Printf("Took %.0f ms\n", elapsed)
// }

// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// // rangeGen отправляет в канал числа от start до stop-1
// func rangeGen(start, stop int) <-chan int {
// 	out := make(chan int)
// 	go func() {
// 		defer close(out)
// 		for i := start; i < stop; i++ {
// 			time.Sleep(50 * time.Millisecond)
// 			out <- i
// 		}

// 	}()
// 	return out
// }

// // начало решения

// // merge выбирает числа из входных каналов и отправляет в выходной
// func merge(channels ...<-chan int) <-chan int {
// 	var wg sync.WaitGroup
// 	out := make(chan int)
// 	defer close(out)

// 	for channel := range channels {
// 		wg.Go(func() {
// 			for val := range channel {
// 				out <- val
// 			}
// 		})
// 	}
// 	go func() {
// 		wg.Wait()
// 		close(out)
// 	}()

// 	return out
// }

// // конец решения

// func main() {
// 	in1 := rangeGen(11, 15)
// 	in2 := rangeGen(21, 25)
// 	in3 := rangeGen(31, 35)

// 	start := time.Now()
// 	merged := merge(in1, in2, in3)
// 	for val := range merged {
// 		fmt.Print(val, " ")
// 	}
// 	fmt.Println()
// 	fmt.Println("Took", time.Since(start))
// }

// package main

// import (
// 	"fmt"
// 	"math/rand"
// )

// // начало решения

// // генерит случайные слова из 5 букв
// // с помощью randomWord(5)
// func generate(cancel <-chan struct{}) <-chan string{
// 	out := make(chan string)
// 	go func() {
// 		defer close(out)
// 		for {
// 			select {
// 			case out <- randomWord(5):
// 			case <-cancel: return
// 			}
// 		}
// 	}()
// 	return out
// }

// // выбирает слова, в которых не повторяются буквы,
// // abcde - подходит
// // abcda - не подходит
// func takeUnique(cancel <-chan struct{}, in <-chan string) <-chan string {
// 	out := make(chan string)
// 	go func() {
// 		defer close(out)
// 		for {
// 			select {
// 			case word, ok := <- in:
// 				if !ok {
// 					return
// 				}
// 				seen := map[rune]struct{}{}
// 				flag := false
// 				for _, letter := range word {
// 					if _, exists := seen[letter]; exists{
// 						flag = true
// 						break
// 					}
// 					seen[letter] = struct{}{}
// 				}
// 				if !flag {
// 					select {
// 					case out <- word:
// 					case <-cancel:
// 						return
// 					}
// 				}
// 			case <-cancel: return
// 			}
// 		}
// 	}()
// 	return out
// }

// // переворачивает слова
// // abcde -> edcba
// func reverse(cancel <-chan struct{}, in <-chan string) <-chan string{
// 	out := make(chan string)
// 	go func() {
// 		defer close(out)
// 		for {
// 			select {
// 			case word, ok := <-in:
// 				if !ok {
// 					return
// 				}
// 				reverseWord := func(s string) string {
// 					runeStr := []rune(s)
// 					reversed := runeStr
// 					for _, rune := range " -> " {
// 						reversed = append(reversed, rune)
// 					}
// 					for i := len(runeStr)-1; i >= 0; i-- {
// 						reversed = append(reversed, runeStr[i])
// 					}
// 					return string(reversed)
// 				}
// 				select {
// 				case out <- reverseWord(word):
// 				case <-cancel:
// 					return
// 				}
// 			case <-cancel: return
// 			}
// 		}
// 	}()
// 	return out
// }

// func merge(cancel <-chan struct{}, c1, c2 <-chan string) <-chan string {
//     out := make(chan string)
//     go func() {
//         defer close(out)
//         for c1 != nil || c2 != nil {
//             select {
//             case v, ok := <-c1:
//                 if !ok {
//                     c1 = nil
//                     continue
//                 }
//                 select {
//                 case out <- v:
//                 case <-cancel:
//                     return
//                 }
//             case v, ok := <-c2:
//                 if !ok {
//                     c2 = nil
//                     continue
//                 }
//                 select {
//                 case out <- v:
//                 case <-cancel:
//                     return
//                 }
//             case <-cancel:
//                 return
//             }
//         }
//     }()
//     return out
// }
// // печатает первые n результатов
// func print(cancel <-chan struct{}, in <-chan string, n int) {
// 	for range(n) {
// 		select {
// 		case val, ok := <- in:
// 			if !ok {
// 				return
// 			}
// 			fmt.Println(val + " ")
// 		case <- cancel: return
// 		}
// 	}
// }

// // конец решения

// // генерит случайное слово из n букв
// func randomWord(n int) string {
// 	const letters = "aeiourtnsl"
// 	chars := make([]byte, n)
// 	for i := range chars {
// 		chars[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return string(chars)
// }

// func main() {
// 	cancel := make(chan struct{})
// 	defer close(cancel)

// 	c1 := generate(cancel)
// 	c2 := takeUnique(cancel, c1)
// 	c3_1 := reverse(cancel, c2)
// 	c3_2 := reverse(cancel, c2)
// 	c4 := merge(cancel, c3_1, c3_2)
// 	print(cancel, c4, 10)
// }

// package main

// import (
// 	"errors"
// 	"fmt"
// )

// var ErrFull = errors.New("Queue is full")
// var ErrEmpty = errors.New("Queue is empty")

// // начало решения

// // Queue - FIFO-очередь на n элементов
// type Queue struct {
//     ch chan int
// }

// // Get возвращает очередной элемент.
// // Если элементов нет и block = false -
// // возвращает ошибку.
// func (q Queue) Get(block bool) (int, error) {
//     if block {
//         return <-q.ch, nil
//     } else {
//         select {
//         case val := <-q.ch:
//             return val, nil
//         default:
//             return 0, ErrEmpty
//         }
//     }
// }

// // Put помещает элемент в очередь.
// // Если очередь заполнения и block = false -
// // возвращает ошибку.
// func (q Queue) Put(val int, block bool) error {
//     if block {
//         q.ch <- val
//         return nil
//     } else {
//         select {
//         case q.ch <- val:
//             return nil
//         default:
//             return ErrFull
//         }
//     }
// }

// // MakeQueue создает новую очередь
// func MakeQueue(n int) Queue {
//     return Queue{
//         ch: make(chan int, n),
//     }
// }

// // конец решения

// func main() {
// 	q := MakeQueue(2)

// 	err := q.Put(1, false)
// 	fmt.Println("put 1:", err)

// 	err = q.Put(2, false)
// 	fmt.Println("put 2:", err)

// 	err = q.Put(3, false)
// 	fmt.Println("put 3:", err)

// 	res, err := q.Get(false)
// 	fmt.Println("get:", res, err)

// 	res, err = q.Get(false)
// 	fmt.Println("get:", res, err)

// 	res, err = q.Get(false)
// 	fmt.Println("get:", res, err)
// }

// package main

// import (
// 	"errors"
// 	"fmt"
// 	"math/rand"
// 	"time"
// )

// // выполняет какую-то операцию,
// // обычно быстро, но иногда медленно
// func work() int {
// 	if rand.Intn(10) < 8 {
// 		time.Sleep(10 * time.Millisecond)
// 	} else {
// 		time.Sleep(200 * time.Millisecond)
// 	}
// 	return 42
// }

// // выполняет функцию fn() c таймаутом timeout и возвращает результат
// // если в течение timeout функция не вернула ответ - возвращает ошибку
// func withTimeout(fn func() int, timeout time.Duration) (int, error) {
// 	var result int

// 	done := make(chan struct{})
// 	go func() {
// 		result = fn()
// 		close(done)
// 	}()

// 	select {
// 	case <-done:
// 		return result, nil
// 	case <-after(timeout):
// 		return 0, errors.New("timeout")
// 	}
// }

// // начало решения

// // возвращает канал, в котором появится значение
// // через промежуток времени dur
// func after(dur time.Duration) <-chan time.Time {
//     ch := make(chan time.Time, 1) // Буфер на 1 элемент

//     go func() {
//         time.Sleep(dur)
//         ch <- time.Now()
//         close(ch)
//     }()

//     return ch
// }

// // конец решения

// func main() {
// 	for i := 0; i < 10; i++ {
// 		start := time.Now()
// 		timeout := 50 * time.Millisecond
// 		if answer, err := withTimeout(work, timeout); err != nil {
// 			fmt.Printf("Took %v. Error: %v\n", time.Since(start), err)
// 		} else {
// 			fmt.Printf("Took %v. Result: %v\n", time.Since(start), answer)
// 		}
// 	}
// }

// package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"time"
// )

// // начало решения

// func delay(dur time.Duration, fn func()) func() {
// 	done := make(chan struct{})
// 	cancel_ch := make(chan struct{})
// 	var cancel_cnt int

// 	timer := time.NewTimer(dur)

// 	cancel := func() {
// 		if timer.Stop() && cancel_cnt == 0 {
// 			cancel_cnt++
// 			cancel_ch <- struct{}{}
// 		}
// 	}

// 	go func() {
// 		defer close(done)
// 		select {
// 		case <-timer.C:
// 			fn()
// 		case <- cancel_ch:
// 		{
// 			close(cancel_ch)
// 			return
// 		}

// 		}
// 	}()

// 	return cancel
// }

// // конец решения

// func main() {
// 	work := func() {
// 		fmt.Println("work done")
// 	}

// 	cancel := delay(100*time.Millisecond, work)

// 	time.Sleep(10 * time.Millisecond)
// 	if rand.Float32() < 0.5 {
// 		cancel()
// 		fmt.Println("delayed function canceled")
// 	}
// 	time.Sleep(100 * time.Millisecond)
// }

// package main

// import (
// 	"fmt"
// 	"time"
// )

// // начало решения

// func schedule(dur time.Duration, fn func()) func() {
// 	ticker := time.NewTicker(dur)

// 	cancel_ch := make(chan struct{})
// 	var cancel_cnt int
// 	cancel := func() {
// 		if cancel_cnt == 0 {
// 			cancel_cnt++
// 			close(cancel_ch)
// 			ticker.Stop()
// 		}
// 	}
// 	go func() {
// 		defer ticker.Stop()
// 		for {
// 			select {
// 			case <-ticker.C:
// 				fn()
// 			case <-cancel_ch:
// 				return
// 			}
// 		}
// 	}()

// 	return cancel
// }

// // конец решения

// func main() {
// 	work := func() {
// 		at := time.Now()
// 		fmt.Printf("%s: work done\n", at.Format("15:04:05.000"))
// 	}

// 	cancel := schedule(50*time.Millisecond, work)
// 	defer cancel()

// 	// хватит на 5 тиков
// 	time.Sleep(260 * time.Millisecond)
// }

// Ограничитель скорости
// package main

// import (
// 	"errors"
// 	"fmt"
// 	"time"
// )

// var ErrCanceled error = errors.New("canceled")

// // начало решения

// // throttle следит, чтобы функция fn выполнялась не более limit раз в секунду.
// // Возвращает функции handle (выполняет fn с учетом лимита) и cancel (останавливает ограничитель).
// func throttle(limit int, fn func()) (handle func() error, cancel func()) {
// 	ticker := time.NewTicker(time.Second/time.Duration(limit))

// 	doFunc := make(chan struct{}, 15)

// 	errCh := make(chan error)
// 	defer close(errCh)

// 	cancel_ch := make(chan struct{})
// 	cancel_cnt := 0
// 	cancel = func() {
// 		select {
// 		case <-cancel_ch:
// 			return
// 		default:
// 			ticker.Stop()
// 			cancel_cnt++
// 			close(cancel_ch)
// 			close(doFunc)
// 		}

// 	}
// 	handle = func() error {
// 		select {
// 		case <-cancel_ch:
// 			return ErrCanceled
// 		default:
// 			doFunc <- struct{}{}
// 			return <-errCh
// 		}

// 	}

// 	go func() {
// 		select {
// 		case <- doFunc:
// 			select {
// 			case <-ticker.C:
// 				fn()
// 			case <-cancel_ch:
// 				if cancel_cnt == 0 {
// 					errCh <- nil
// 				} else {
// 					errCh <- ErrCanceled
// 				}
// 			}

// 		}
// 	}()

// 	return handle, cancel
// }

// // конец решения

// func main() {
// 	work := func() {
// 		fmt.Print(".")
// 	}

// 	handle, cancel := throttle(5, work)
// 	defer cancel()

// 	start := time.Now()
// 	const n = 10
// 	for i := 0; i < n; i++ {
// 		handle()
// 	}
// 	fmt.Println()
// 	fmt.Printf("%d queries took %v\n", n, time.Since(start))
// }

// package main

// import (
// 	"context"
// 	"fmt"
// )

// // начало решения

// // генерит целые числа от start и до бесконечности
// func generate(ctx context.Context, start int) <-chan int {
// 	out := make(chan int)
// 	go func() {
// 		defer close(out)
// 		for i := start; ; i++ {
// 			select {
// 			case out <- i:
// 			case <-ctx.Done():
// 				return
// 			}
// 		}
// 	}()
// 	return out
// }

// // конец решения

// func main() {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	generated := generate(ctx, 11)
// 	for num := range generated {
// 		fmt.Print(num, " ")
// 		if num > 14 {
// 			break
// 		}
// 	}
// 	fmt.Println()
// }

// package main

// import (
// 	"context"
// 	"fmt"
// 	"strings"
// 	"unicode"
// )

// // информация о количестве цифр в каждом слове
// type counter map[string]int

// // слово и количество цифр в нем
// type pair struct {
// 	word  string
// 	count int
// }

// // начало решения

// // считает количество цифр в словах
// func countDigitsInWords(ctx context.Context, words []string) counter {
// 	childCtx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	select {
// 	case <-ctx.Done():
// 		cancel()
// 		return counter{}
// 	default:
// 		pending := submitWords(childCtx, words)
// 		counted := countWords(childCtx, pending)
// 		return fillStats(childCtx, counted)
// 	}

// }

// // отправляет слова на подсчет
// func submitWords(ctx context.Context, words []string) <-chan string {
// 	out := make(chan string)
// 	go func() {
// 		for _, word := range words {
// 			select {
// 			case <- ctx.Done():
// 				close(out)
// 				return
// 			case out <- word:
// 			}
// 		}
// 		close(out)
// 	}()
// 	return out
// }

// // считает цифры в словах
// func countWords(ctx context.Context, in <-chan string) <-chan pair {
// 	out := make(chan pair)
// 	go func() {
// 		for word := range in {
// 			count := countDigits(word)
// 			select {
// 			case <- ctx.Done():
// 				close(out)
// 				return
// 			case out <- pair{word, count}:
// 			}
// 		}
// 		close(out)
// 	}()
// 	return out
// }

// // готовит итоговую статистику
// func fillStats(ctx context.Context, in <-chan pair) counter {
// 	stats := counter{}
// 	for p := range in {
// 		select {
// 		case <-ctx.Done():
// 			return counter{}
// 		default:
// 			stats[p.word] = p.count
// 		}
// 	}
// 	return stats
// }

// // конец решения

// // считает количество цифр в слове
// func countDigits(str string) int {
// 	count := 0
// 	for _, char := range str {
// 		if unicode.IsDigit(char) {
// 			count++
// 		}
// 	}
// 	return count
// }

// func main() {
// 	phrase := "0ne 1wo thr33 4068"
// 	words := strings.Fields(phrase)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	stats := countDigitsInWords(ctx, words)
// 	fmt.Println(stats)
// }

// package main

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"time"
// )

// // начало решения

// // ErrFailed и ErrManual - причины остановки цикла.
// var ErrFailed = errors.New("failed")
// var ErrManual = errors.New("manual")

// // Worker выполняет заданную функцию в цикле, пока не будет остановлен.
// // Гарантируется, что Worker используется только в одной горутине.
// type Worker struct {
// 	fn      func() error
// 	ctx     context.Context
// 	cancel  context.CancelFunc
// 	stopCnt int
// 	started bool
// 	err error
// 	// TODO: добавить поля
// }

// // NewWorker создает новый экземпляр Worker с заданной функцией.
// // Но пока не запускает цикл с функцией.
// func NewWorker(fn func() error) *Worker {
// 	cont, cancelfunc := context.WithCancel(context.Background())
// 	return &Worker{fn: fn, ctx: cont, cancel: cancelfunc, stopCnt: 0, started: false, err: nil}
// }

// // Start запускает отдельную горутину, в которой циклически
// // выполняет заданную функцию, пока не будет вызван метод Stop,
// // либо пока функция не вернет ошибку.
// // Повторные вызовы Start игнорируются.
// func (w *Worker) Start() {
// 	// TODO: реализовать требования
// 	if !w.started {
// 		w.started = true
// 		go func() {
// 			for {
// 				select {
// 				case <-w.ctx.Done():
// 					return
// 				default:
// 					err := w.fn()
// 					if err != nil {
// 						w.err = ErrFailed
// 						return
// 					}
// 				}
// 			}
// 		}()
// 	}
// }

// // Stop останавливает выполнение цикла.
// // Вызов Stop до Start игнорируется.
// // Повторные вызовы Stop игнорируются.
// func (w *Worker) Stop() {
// 	if w.stopCnt == 0 && w.started {
// 		w.err = ErrManual
// 		w.stopCnt++
// 		w.cancel()
// 	}
// }

// // AfterStop регистрирует функцию, которая
// // будет вызвана после остановки цикла.
// // Можно зарегистрировать несколько функций.
// // Вызовы AfterStop после Start игнорируются.
// func (w *Worker) AfterStop(fn func()) {
// 	// TODO: реализовать требования
// 	if !w.started {
// 		context.AfterFunc(w.ctx, fn)
// 	}
// }

// // Err возвращает причину остановки цикла:
// // - ErrManual - вручную через метод Stop;
// // - ErrFailed - из-за ошибки, которую вернула функция.
// func (w *Worker) Err() error {
// 	// TODO: реализовать требования
// 	return w.err
// }

// // конец решения

// func main() {
// 	{
// 		// Start-Stop
// 		count := 9
// 		fn := func() error {
// 			fmt.Print(count, " ")
// 			count--
// 			time.Sleep(105 * time.Millisecond)
// 			return nil
// 		}

// 		worker := NewWorker(fn)
// 		worker.Start()
// 		time.Sleep(10 * time.Millisecond)
// 		worker.Stop()

// 		fmt.Println()
// 		// 9 8 7 6 5 4 3 2 1 0
// 	}
// 	{
// 		// ErrFailed
// 		count := 3
// 		fn := func() error {
// 			fmt.Print(count, " ")
// 			count--
// 			if count == 0 {
// 				return errors.New("count is zero")
// 			}
// 			time.Sleep(10 * time.Millisecond)
// 			return nil
// 		}

// 		worker := NewWorker(fn)
// 		worker.Start()
// 		time.Sleep(35 * time.Millisecond)
// 		worker.Stop()

// 		fmt.Println(worker.Err())
// 		// 3 2 1 failed
// 	}
// 	{
// 		// AfterStop
// 		fn := func() error { return nil }

// 		worker := NewWorker(fn)
// 		worker.AfterStop(func() {
// 			fmt.Println("called after stop")
// 		})

// 		worker.Start()
// 		worker.Stop()

//			time.Sleep(10 * time.Millisecond)
//			// called after stop
//		}
//	}
// package main

// import (
// 	"fmt"
// 	"time"
// )

// func delay(duration time.Duration, fn func()) func() {
//     canceled := false          // (1)

//     go func() {
//         time.Sleep(duration)
//         if !canceled {         // (2)
//             fn()
//         }
//     }()

//     cancel := func() {
//         canceled = true        // (3)
//     }
//     return cancel              // (4)
// }

// func main() {
// 	work := func() {
// 		fmt.Println("work done")
// 	}

// 	cancel := delay(50*time.Millisecond, work)
// 	time.Sleep(50 * time.Millisecond)
// 	go cancel()
// }

// package main

// import (
// 	"fmt"
// 	"sync"
// )

// // начало решения

// type Counter struct {
// 	freqMap map[string]int
// 	mtx sync.Mutex
// }

// func (c *Counter) Increment(str string) {
// 	c.mtx.Lock()
// 	c.freqMap[str]++
// 	c.mtx.Unlock()
// }

// func (c *Counter) Value(str string) int {
// 	c.mtx.Lock()
// 	defer c.mtx.Unlock()
// 	if val, ok := c.freqMap[str]; ok {
// 		return val
// 	}
// 	return 0
// }

// func (c *Counter) Range(fn func(key string, val int)) {
// 	c.mtx.Lock()
// 	for key, value := range c.freqMap {
// 		fn(key, value)
// 	}
// 	c.mtx.Unlock()
// }

// func NewCounter() *Counter {
// 	return &Counter{map[string]int{}, sync.Mutex{}}
// }

// // конец решения

// func main() {
// 	counter := NewCounter()

// 	var wg sync.WaitGroup
// 	wg.Add(3)

// 	increment := func(key string, val int) {
// 		defer wg.Done()
// 		for ; val > 0; val-- {
// 			counter.Increment(key)
// 		}
// 	}

// 	go increment("one", 100)
// 	go increment("two", 200)
// 	go increment("three", 300)

// 	wg.Wait()

// 	fmt.Println("two:", counter.Value("two"))

// 	fmt.Print("{ ")
// 	counter.Range(func(key string, val int) {
// 		fmt.Printf("%s:%d ", key, val)
// 	})
// 	fmt.Println("}")
// }

// package main

// import (
// 	"fmt"
// 	"sync"
// )

// // начало решения

// type Counter struct {
// 	lock sync.RWMutex
// 	freqMap map[string]int
// }

// func (c *Counter) Increment(str string) {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()
// 	c.freqMap[str]++
// }

// func (c *Counter) Value(str string) int {
// 	c.lock.RLock()
// 	defer c.lock.RUnlock()
// 	return c.freqMap[str]
// }

// func (c *Counter) Range(fn func(key string, val int)) {
// 	c.lock.RLock()
// 	defer c.lock.RUnlock()
// 	for key, value := range c.freqMap {
// 		fn(key, value)
// 	}
// }

// func NewCounter() *Counter {
// 	return &Counter{
// 		lock: sync.RWMutex{},
// 		freqMap: map[string]int{},
// 	}
// }

// // конец решения

// func main() {
// 	counter := NewCounter()

// 	var wg sync.WaitGroup
// 	wg.Add(3)

// 	increment := func(key string, val int) {
// 		defer wg.Done()
// 		for ; val > 0; val-- {
// 			counter.Increment(key)
// 		}
// 	}

// 	go increment("one", 100)
// 	go increment("two", 200)
// 	go increment("three", 300)

// 	wg.Wait()

// 	fmt.Println("two:", counter.Value("two"))

// 	fmt.Print("{ ")
// 	counter.Range(func(key string, val int) {
// 		fmt.Printf("%s:%d ", key, val)
// 	})
// 	fmt.Println("}")
// }

// package main

// import (
// 	"fmt"
// 	"sync"
// 	"sync/atomic"
// )

// // начало решения

// type Total struct {
// 	atom atomic.Int32
// }

// func (t *Total) Increment() {
// 	t.atom.Add(1)
// }

// func (t *Total) Value() int {
// 	return int(t.atom.Load())
// }

// // конец решения

// func main() {
// 	var wg sync.WaitGroup

// 	var total Total

// 	for i := 0; i < 5; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			for i := 0; i < 10000; i++ {
// 				total.Increment()
// 			}
// 		}()
// 	}

// 	wg.Wait()
// 	fmt.Println("total", total.Value())
// }

// package main

// import (
// 	"fmt"
// 	"strings"
// )

// // начало решения

// // slugify возвращает "безопасный" вариант заголовока:
// // только латиница, цифры и дефис
// func slugify(src string) string {
// 	safeFunc := func(r rune) bool {
// 		// a-z и A-Z, цифры 0-9 и дефис -
// 		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' {
// 			return false
// 		}
// 		return true
// 	}
// 	//srcRune := []rune(src)
// 	res := strings.ToLower(strings.Join(strings.FieldsFunc(src, safeFunc), "-"))
// 	return res
// }

// // конец решения

// func main() {
//     phrase := "Go Is Awesome!"
//     fmt.Println(slugify(phrase))
//     // go-is-awesome

//     phrase = "Tabs are all we've got"
//     fmt.Println(slugify(phrase))
//     // tabs-are-all-we-ve-got
// }

// package main

// import (
// 	"fmt"
// 	"strconv"
// 	"strings"
// )

// // начало решения

// // calcDistance возвращает общую длину маршрута в метрах
// func calcDistance(directions []string) int {
// 	distance := 0
// 	for _, s := range directions {
// 		if strings.ContainsAny(s, "0123456789") {
// 			for word := range strings.SplitSeq(s, " ") {
// 				if strings.ContainsAny(word, "0123456789") {
// 					var letterPos int
// 					var flag bool
// 					if strings.Contains(word, "k") {
// 							letterPos = strings.Index(word, "k")
// 							flag = true
// 						} else {
// 							letterPos = strings.Index(word, "m")

// 					}
// 					if strings.Contains(word, ".") {
// 						d, _ := strconv.ParseFloat(word[:letterPos], 32)
// 						distance += int(d*1000)
// 					} else {
// 						d, _ := strconv.Atoi(word[:letterPos])
// 						if flag {
// 							distance += d*1000
// 						} else {
// 							distance += d
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return distance
// }

// // конец решения

// func main() {
// 	directions := []string{
// 		"100m to intersection",
// 		"turn right",
// 		"straight 300m",
// 		"enter motorway",
// 		"straight 5km",
// 		"exit motorway",
// 		"500m straight",
// 		"turn sharp left",
// 		"continue 100m to destination",
// 	}
// 	const want = 6000
// 	got := calcDistance(directions)
// 	fmt.Print(got)
// }

// package main

// import (
// 	"fmt"
// 	"slices"
// 	"strings"
// 	"strconv"
// )

// // начало решения

// // prettify возвращает отформатированное
// // строковое представление карты
// func prettify(m map[string]int) string {
// 	if len(m) == 0 {
// 		return "{}"
// 	}
// 	if len(m) == 1 {
// 		for key, val := range m {
// 			return fmt.Sprintf("{ %s: %d }", key, val)
// 		}
// 	}
// 	var b strings.Builder
// 	b.WriteRune('{')
// 	b.WriteRune('\n')
// 	sortedKeys := []string{}
// 	for key := range m {
// 		sortedKeys = append(sortedKeys, key)
// 	}
// 	slices.Sort(sortedKeys)
// 	sortedLen := len(sortedKeys)
// 	for idx, key := range sortedKeys {
// 		b.WriteString("    ")
// 		b.WriteString(key)
// 		b.WriteString(": ")
// 		b.WriteString(strconv.Itoa(m[key]))
// 		b.WriteRune(',')
// 		if !(idx == sortedLen - 1) {
// 			b.WriteRune('\n')
// 		}
// 	}
// 	b.WriteRune('\n')
// 	b.WriteRune('}')
// 	return b.String()
// }

// // конец решения

// func main() {
// 	m := map[string]int{"one": 1, "two": 2, "three": 3}
// 	const want = "{\n    one: 1,\n    three: 3,\n    two: 2,\n}"
// 	got := prettify(m)
// 	print(got)
// }

// package main

// import (
// 	"regexp"
// 	"strings"
// )

// // начало решения

// // slugify возвращает "безопасный" вариант заголовока:
// // только латиница, цифры и дефис
// func slugify(src string) string {
// //	re := regexp.MustCompile(`[abcdefghijklmnopqrstuvwxyz0123456789\- ]*\s`)
// 	re := regexp.MustCompile(`[^a-z0-9]+`)
// 	src = strings.ToLower(src)
// //	return strings.Join(re.Split(src, -1), "-")
// 	res := re.ReplaceAllString(src, "-")
// 	return res[:len(res)-1]
// }

// // конец решения

// func main() {
// 	const phrase = "Go Is Awesome!"
// 	const want = "go-is-awesome"
// 	got := slugify(phrase)
// 	print(got)
// }

// package main

// import (
// 	"bytes"
// 	"text/template"
// )

// // начало решения
// // Алиса, добрый день! Ваш баланс - 1000₽. Все в порядке.

// // const txt = `Сейчас {{.Time}}, {{.Day}}.
// // {{if .Sunny -}} Солнечно! {{- else -}} Пасмурно :-/ {{- end}}
// // `

// var templateText = `{{.Name}}, добрый день! Ваш баланс - {{.Balance}} {{if ge .Balance 100 -}}₽. Все в порядке. {{- end}}{{if eq .Balance 0 -}} Доступ заблокирован. {{- else -}}{{if le .Balance 100 -}} Пора пополнить. {{- end}}{{- end}}`

// // конец решения

// type User struct {
// 	Name    string
// 	Balance int
// }

// // renderToString рендерит данные по шаблону в строку
// func renderToString(tpl *template.Template, data any) string {
// 	var buf bytes.Buffer
// 	tpl.Execute(&buf, data)
// 	return buf.String()
// }

// func main() {
// 	tpl := template.New("message")
// 	tpl = template.Must(tpl.Parse(templateText))

// 	user := User{"Алиса", 500}
// 	got := renderToString(tpl, user)

// 	const want = "Алиса, добрый день! Ваш баланс - 500₽. Все в порядке."
// 	print(got)
// }
// package main

// import (
// 	"fmt"
// 	"slices"
// 	"strings"
// 	"strconv"
// )
// // начало решения

// // slugify возвращает "безопасный" вариант заголовока:
// // только латиница, цифры и дефис
// func slugify(src string) string {
//     //var b strings.Builder
// 	safeFunc := func(r rune) bool {
// 		// a-z и A-Z, цифры 0-9 и дефис -
// 		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' {
// 			return false
// 		}
// 		return true
// 	}

// 	//srcRune := []rune(src)
// 	res := strings.ToLower(strings.Join(strings.FieldsFunc(src, safeFunc), "-"))
// 	return res
// }

// // конец решения

// package main

// import (
// 	"errors"
// 	"fmt"

// 	"time"
// )

// // начало решения

// // TimeOfDay описывает время в пределах одного дня
// type TimeOfDay struct {
// 	hour, min, sec int
// 	loc *time.Location
// }

// // Hour возвращает часы в пределах дня
// func (t TimeOfDay) Hour() int {
// 	return t.hour
// }

// // Minute возвращает минуты в пределах часа
// func (t TimeOfDay) Minute() int {
// 	return t.min
// }

// // Second возвращает секунды в пределах минуты
// func (t TimeOfDay) Second() int {
// 	return t.sec
// }

// // String возвращает строковое представление времени
// // в формате чч:мм:сс TZ (например, 12:34:56 UTC)
// func (t TimeOfDay) String() string {
// 	var h, m, s string
// 	if t.hour < 10 {
// 		h = fmt.Sprintf("0%d", t.hour)
// 	} else {
// 		h = fmt.Sprintf("%d", t.hour)
// 	}

// 	if t.min < 10 {
// 		m = fmt.Sprintf("0%d", t.min)
// 	} else {
// 		m = fmt.Sprintf("%d", t.min)
// 	}

// 	if t.sec < 10 {
// 		s = fmt.Sprintf("0%d", t.sec)
// 	} else {
// 		s = fmt.Sprintf("%d", t.sec)
// 	}

// 	return fmt.Sprintf("%s:%s:%s %v", h, m, s, t.loc)
// }

// // Equal сравнивает одно время с другим.
// // Если у t и other разные локации - возвращает false.
// func (t TimeOfDay) Equal(other TimeOfDay) bool {
// 	if t.loc.String() != other.loc.String() {
// 		return false
// 	}
// 	if t.hour == other.hour && t.min == other.min && t.sec == other.sec {
// 		return true
// 	}
// 	return false
// }

// // Before возвращает true, если время t предшествует other.
// // Если у t и other разные локации - возвращает ошибку.
// func (t TimeOfDay) Before(other TimeOfDay) (bool, error) {
// 	if t.loc.String() != other.loc.String() {
// 		return false, errors.New("dif locs")
// 	}
// 	if t.Equal(other) {
// 		return false, nil
// 	}
// 	if t.hour < other.hour {
// 		return true, nil
// 	} else {
// 		if t.hour > other.hour {
// 			return false, nil
// 		}
// 		if t.min < other.min {
// 			return true, nil
// 		} else {
// 			if t.sec < other.sec {
// 				return true, nil
// 			}
// 			return false, nil
// 		}
// 	}
// }

// // After возвращает true, если время t идет после other.
// // Если у t и other разные локации - возвращает ошибку.
// func (t TimeOfDay) After(other TimeOfDay) (bool, error) {
// 	if t.loc.String() != other.loc.String() {
// 		return false, errors.New("dif locs")
// 	}
// 	if t.Equal(other) {
// 		return false, nil
// 	}
// 	if t.hour > other.hour {
// 		return true, nil
// 	} else {
// 		if t.hour > other.hour {
// 			return false, nil
// 		}
// 		if t.min > other.min {
// 			return true, nil
// 		} else {
// 			if t.sec > other.sec {
// 				return true, nil
// 			}
// 			return false, nil
// 		}
// 	}
// }

// // MakeTimeOfDay создает время в пределах дня
// func MakeTimeOfDay(hour, min, sec int, loc *time.Location) TimeOfDay {
// 	if hour > 60 || min > 60 || sec > 60 || hour < 0 || sec < 0 || min < 0 {
// 		return TimeOfDay{}
// 	}
// 	return TimeOfDay{hour, min, sec, loc}
// }

// // конец решения

// func main() {
// 	t1 := MakeTimeOfDay(17, 45, 22, time.UTC)
// 	t2 := MakeTimeOfDay(20, 3, 4, time.UTC)

// 	if t1.Equal(t2) {
// 		fmt.Printf("%v should not be equal to %v", t1, t2)
// 	}

// 	before, _ := t1.Before(t2)
// 	if !before {
// 		fmt.Printf("%v should be before %v", t1, t2)
// 	}

// 	after, _ := t1.After(t2)
// 	if after {
// 		fmt.Printf("%v should NOT be after %v", t1, t2)
// 	}
// }

// package main

// import (
// 	"errors"
// 	"fmt"
// 	"regexp"
// 	"slices"

// 	//	"slices"
// 	"strings"
// 	"time"
// )

// // начало решения

// // Task описывает задачу, выполненную в определенный день
// type Task struct {
// 	Date  time.Time
// 	Dur   time.Duration
// 	Title string
// }

// // ParsePage разбирает страницу журнала
// // и возвращает задачи, выполненные за день
// func ParsePage(src string) ([]Task, error) {
// 	lines := strings.Split(src, "\n")
// 	date, err := parseDate(lines[0])
// 	tasks, err := parseTasks(date, lines[1:])
// 	sortTasks(tasks)
// 	return tasks, err
// }

// // parseDate разбирает дату в формате дд.мм.гггг
// func parseDate(src string) (time.Time, error) {
// 	return time.Parse("02.01.2006", src)
// }

// var re = regexp.MustCompile(`(\d+:\d+) - (\d+:\d+) (.+)`)
// // parseTasks разбирает задачи из записей журнала
// func parseTasks(date time.Time, lines []string) ([]Task, error) {
// 	taskIdx := make(map[string]int, len(lines))
// 	var tasks []Task

// 	for _, task := range lines {
// 		groups := re.FindStringSubmatch(task)
// 		if len(groups) < 4 {
// 			return []Task{}, errors.New("error while parsing regexp")
// 		}
// 		start, err := time.Parse("15:04", groups[1])
// 		if err != nil {
// 			return []Task{}, errors.New("error while parsing")
// 		}
// 		end, err := time.Parse("15:04", groups[2])
// 		if err != nil {
// 			return []Task{}, errors.New("error while parsing")
// 		}
// 		title := groups[3]
// 		dur := end.Sub(start)
// 		if idx, exists := taskIdx[title]; !exists {
// 			taskIdx[title] = len(tasks)
// 			tasks = append(tasks, Task{
//                 Date:  date,
//                 Dur:   dur,
//                 Title: title,
//             })
// 		} else {
// 			tasks[idx].Dur += dur
// 		}
// 	}
// 	return tasks, nil
// }

// // sortTasks упорядочивает задачи по убыванию длительности
// func sortTasks(tasks []Task) {
// 	sortFunc := func(a, b Task) int {
// 		if a.Dur > b.Dur {
// 			return -1
// 		}
// 		if a.Dur < b.Dur {
// 			return 1
// 		}
// 		return 0
// 	}
// 	slices.SortFunc(tasks, sortFunc)
// }

// // конец решения
// // ::footer

// func main() {
// 	page := `15.04.2022
// 8:00 - 8:30 Завтрак
// 8:30 - 9:30 Оглаживание кота
// 9:30 - 10:00 Интернеты
// 10:00 - 14:00 Напряженная работа
// 14:00 - 14:45 Обед
// 14:45 - 15:00 Оглаживание кота
// 15:00 - 19:00 Напряженная работа
// 19:00 - 19:30 Интернеты
// 19:30 - 22:30 Безудержное веселье
// 22:30 - 23:00 Оглаживание кота`

// 	entries, err := ParsePage(page)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Мои достижения за", entries[0].Date.Format("2006-01-02"))
// 	for _, entry := range entries {
// 		fmt.Printf("- %v: %v\n", entry.Title, entry.Dur)
// 	}

// 	// ожидаемый результат
// 	/*
// 		Мои достижения за 2022-04-15
// 		- Напряженная работа: 8h0m0s
// 		- Безудержное веселье: 3h0m0s
// 		- Оглаживание кота: 1h45m0s
// 		- Интернеты: 1h0m0s
// 		- Обед: 45m0s
// 		- Завтрак: 30m0s
// 	*/
// }

// package main

// import (
// 	"errors"
// 	"fmt"
// 	"strconv"
// 	"strings"
// //	"testing"
// 	"time"
// )

// // начало решения

// // asLegacyDate преобразует время в легаси-дату
// func asLegacyDate(t time.Time) string {
// 	sec := t.Unix()
// 	nsec := t.UnixNano()
// 	return fmt.Sprintf("%d.%09d", sec, nsec-sec)
// }

// // parseLegacyDate преобразует легаси-дату во время.
// // Возвращает ошибку, если легаси-дата некорректная.
// func parseLegacyDate(d string) (time.Time, error) {
// 	if len(d) < 3 || strings.Index(d, ",") == len(d)-1 || !strings.Contains(d, ".") || d[0] == '.' || strings.Contains(d, ",") {
// 		return time.Time{}, errors.New("invalid legacy date")
// 	}
// 	times := strings.Split(d, ".")
// 	sec, err := strconv.Atoi(times[0])
// 	if err != nil {
// 		return time.Time{}, errors.New("error while atoi sec")
// 	}
// 	nsec, err := strconv.Atoi(times[1])
// 	if err != nil {
// 		return time.Time{}, errors.New("error while atoi nsec")
// 	}
// 	unixTime := time.Unix(int64(sec), int64(nsec))
// 	for i := len(times[1]); i < 9; i++ {
// 		nsec *= 10
// 	}
// 	return unixTime, nil
// }

// // конец решения

// func main() {
// 	samples := map[time.Time]string{
// 		time.Date(1970, 1, 1, 1, 0, 0, 123456789, time.UTC): "3600.123456789",
// 		time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC):         "3600.0",
// 		time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC):         "0.0",
// 	}
// 	for src, want := range samples {
// 		got := asLegacyDate(src)
// 		if got != want {
// 			fmt.Printf("%v: got %v, want %v", src, got, want)
// 		}
// 	}

// 	samples1 := map[string]time.Time{
// 		"3600.123456789": time.Date(1970, 1, 1, 1, 0, 0, 123456789, time.UTC),
// 		"3600.0":         time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC),
// 		"0.0":            time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
// 		"1.123456789":    time.Date(1970, 1, 1, 0, 0, 1, 123456789, time.UTC),
// 	}
// 	for src, want := range samples1 {
// 		got, err := parseLegacyDate(src)
// 		if err != nil {
// 			fmt.Printf("%v: unexpected error", src)
// 		}
// 		if !got.Equal(want) {
// 			fmt.Printf("%v: got %v, want %v\n", src, got, want)
// 		}

// }
// }

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// //	"strings"
// )

// // начало решения

// // readLines возвращает все строки из указанного файла
// func readLines(name string) ([]string, error) {
// 	file, err := os.Open(name)
// 	if err != nil {
// 		panic(err)
// 	}
// 	var res []string
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		res = append(res, line)
// 	}
// 	// text := string(fileText)
// 	// res := strings.Split(text, "\n")
// 	// if res[len(res)-1] == "" {
// 	// 	return res[:len(res)-1], nil
// 	// }
// 	return res, nil
// }

// // конец решения

// func main() {
// 	lines, err := readLines("task.data")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Print(lines)
// 	// for idx, line := range lines {
// 	// 	fmt.Printf("%d: %s\n", idx+1, line)
// 	// }
// }

// package main

// import (
// 	"bufio"
// 	"os"
// 	"strings"
// 	"fmt"
// )

// func main() {
// 	var str []string
// 	scanner := bufio.NewScanner(os.Stdin)
// 	scanner.Split(bufio.ScanWords)
// 	for scanner.Scan() {
// 		word := []rune(scanner.Text())
// 		titleword := strings.ToUpper(string(word[0])) + strings.ToLower(string(word[1:]))
// 		//titleword := strings.ToUpper(string(rune(word[0]))) + word[1:]
// 		str = append(str, titleword)
// 	}
// 	if err := scanner.Err(); err != nil {
// 		panic(err)
// 	}
// 	for i, letter := range str {
// 		if i == (len(str)-1) {
// 			fmt.Print(letter)
// 			continue
// 		}
// 		fmt.Print(letter)
// 		fmt.Print(" ")
// 	}
// }

// package main

// import (
// 	"bufio"
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"crypto/rand"
// )

// // начало решения

// // RandomReader создает читателя, который возвращает случайные байты,
// // но не более max штук
// func RandomReader(max int) io.Reader {
// 	b := make([]byte, max)
// 	_, _ = rand.Read(b)
// 	reader := bytes.NewReader(b)
// 	rd := bufio.NewReader(reader)
// 	return rd
// }

// // конец решения

// func main() {
// 	rnd := RandomReader(5)
// 	rd := bufio.NewReader(rnd)
// 	for {
// 		b, err := rd.ReadByte()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Printf("%d ", b)
// 	}
// 	fmt.Println()
// 	// 1 148 253 194 250
// 	// (значения могут отличаться)
// }

// package main

// import (
// 	"fmt"
// 	"io"
// 	"strings"
// )

// // начало решения

// // AbyssWriter пишет данные в никуда,
// // но при этом считает количество записанных байт
// type AbyssWriter struct{
// 	wrtnBytes int
// }

// func (w *AbyssWriter) Write(p []byte) (n int, err error) {
// 	w.wrtnBytes += len(p)
// 	return len(p), nil
// }
// // Total возвращает общее количество записанных байт
// func (w *AbyssWriter) Total() int {
// 	return w.wrtnBytes
// }

// // NewAbyssWriter создает новый AbyssWriter
// func NewAbyssWriter() *AbyssWriter {
// 	return &AbyssWriter{0}
// }

// // конец решения

// func main() {
// 	r := strings.NewReader("go is awesome")
// 	w := NewAbyssWriter()
// 	written, err := io.Copy(w, r)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("written %d bytes\n", written)
// 	fmt.Println(written == int64(w.Total()))
// }

// package main

// import (
// 	"fmt"
// 	"io"
// )

// // TokenReader начитывает токены из источника
// type TokenReader interface {
// 	// ReadToken считывает очередной токен
// 	// Если токенов больше нет, возвращает ошибку io.EOF
// 	ReadToken() (string, error)
// }

// // TokenWriter записывает токены в приемник
// type TokenWriter interface {
// 	// WriteToken записывает очередной токен
// 	WriteToken(s string) error
// }

// // начало решения

// // FilterTokens читает все токены из src и записывает в dst тех,
// // кто проходит проверку predicate
// func FilterTokens(dst TokenWriter, src TokenReader, predicate func(s string) bool) (int, error) {
// 	total := 0
// 	for {
// 		token, err := src.ReadToken()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return total, err
// 		}
// 		if predicate(token) {
// 			err := dst.WriteToken(token)
// 			if err != nil {
// 				return total, err
// 			}
// 			total++
// 		}
// 	}
// 	return total, nil
// }

// // конец решения
// type WordReader struct{}
// type WordWriter struct{}
// func main() {
// 	// Для проверки придется создать конкретные типы,
// 	// которые реализуют интерфейсы TokenReader и TokenWriter.

// 	// Ниже для примера используются NewWordReader и NewWordWriter,
// 	// но вы можете сделать любые на свое усмотрение.

// 	r := NewWordReader("go is awesome")
// 	w := NewWordWriter()
// 	predicate := func(s string) bool {
// 		return s != "is"
// 	}
// 	n, err := FilterTokens(w, r, predicate)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("%d tokens: %v\n", n, w.Words())
// 	// 2 tokens: [go awesome]
// }

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	mathrand "math/rand"
// 	"os"
// 	"path/filepath"
// 	"strings"
// )

// // алфавит планеты Нибиру
// const alphabet = "aeiourtnsl"

// // Census реализует перепись населения.
// // Записи о рептилоидах хранятся в каталоге census, в отдельных файлах,
// // по одному файлу на каждую букву алфавита.
// // В каждом файле перечислены рептилоиды, чьи имена начинаются
// // на соответствующую букву, по одному рептилоиду на строку.
// type Census struct{
// 	openedFiles []*os.File
// 	Total int
// }

// // Count возвращает общее количество переписанных рептилоидов.
// func (c *Census) Count() int {
// 	return c.Total
// }

// // Add записывает сведения о рептилоиде.
// func (c *Census) Add(name string) {
// 	var file *os.File
// 	var err error
// 	firstLetter := string(name[0])
// 	fileName := firstLetter + ".txt"
// 	index := strings.Index(alphabet, firstLetter)
// 	if !(c.openedFiles[index] != nil) {
// 		file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0777)
// 		if err != nil {
// 		panic(err)
// 		}
// 		c.openedFiles[strings.Index(alphabet, firstLetter)] = file
// 	} else {
// 		file = c.openedFiles[index]
// 	}
// 	w := bufio.NewWriter(file)
// 	_, err = w.WriteString(name)
// 	if err != nil {
// 		panic(err)
// 	}
// 	w.WriteByte('\n')
// 	w.Flush()
// 	c.Total++
// }

// // Close закрывает файлы, использованные переписью.
// func (c *Census) Close() {
// 	for _, file := range c.openedFiles {
// 		err := file.Close()
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

// // NewCensus создает новую перепись и пустые файлы
// // для будущих записей о населении.
// func NewCensus() *Census {
// 	os.RemoveAll("census")
// 	err := os.Mkdir("census", 0777)
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = os.Chdir(filepath.FromSlash("census"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	touch := func(path string){
// 		p := filepath.FromSlash(path)
// 		data := []byte{}
// 		err := os.WriteFile(p, data, 0777)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	//aeiourtnsl
// 	for _, letter := range alphabet {
// 		path := string(letter) + ".txt"
// 		touch(path)
// 	}
// 	return &Census{
// 		openedFiles: make([]*os.File, len(alphabet)),
// 		Total: 0,
// 	}
// }

// // ┌─────────────────────────────────┐
// // │ не меняйте код ниже этой строки │
// // └─────────────────────────────────┘

// var rand = mathrand.New(mathrand.NewSource(0))

// // randomName возвращает имя очередного рептилоида.
// func randomName(n int) string {
// 	chars := make([]byte, n)
// 	for i := range chars {
// 		chars[i] = alphabet[rand.Intn(len(alphabet))]
// 	}
// 	return string(chars)
// }

// func main() {
// 	census := NewCensus()
// 	defer census.Close()
// 	for i := 0; i < 1024; i++ {
// 		reptoid := randomName(5)
// 		census.Add(reptoid)
// 	}
// 	fmt.Println(census.Count())
// }

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strings"
// 	"time"
// )

// // начало решения

// // Duration описывает продолжительность фильма
// type Duration time.Duration

// func (d Duration) MarshalJSON() ([]byte, error) {
// 	var b strings.Builder
// 	b.WriteByte('"')
// 	hours := int(time.Duration(d).Hours())
// 	mins := int(time.Duration(d).Minutes()) - 60*hours
// 	if hours == 0 && mins == 0 {
// 		return nil, nil
// 	} else if hours == 0 {
// 		t := fmt.Sprintf("%dm", mins)
// 		b.WriteString(t)
// 	} else if mins == 0 {
// 		t := fmt.Sprintf("%dh", hours)
// 		b.WriteString(t)
// 	} else {
// 		t := fmt.Sprintf("%dh%dm", hours, mins)
// 		b.WriteString(t)
// 	}
// 	b.WriteByte('"')
// 	return []byte(b.String()), nil
// }

// // Rating описывает рейтинг фильма
// type Rating int

// func (r Rating) MarshalJSON() ([]byte, error) {
// 	var b strings.Builder
// 	b.WriteByte('"')
// 	switch r {
// 		// ☆★
// 	case 0:
// 		b.WriteString("☆☆☆☆☆")
// 	case 1:
// 		b.WriteString("★☆☆☆☆")
// 	case 2:
// 		b.WriteString("★★☆☆☆")
// 	case 3:
// 		b.WriteString("★★★☆☆")
// 	case 4:
// 		b.WriteString("★★★★☆")
// 	case 5:
// 		b.WriteString("★★★★★")
// 	}
// 	b.WriteByte('"')
// 	return []byte(b.String()), nil
// }

// // Movie описывает фильм
// type Movie struct {
// 	Title string
// 	Year int
// 	Director string
// 	Genres []string
// 	Duration Duration
// 	Rating Rating
// }

// // MarshalMovies кодирует фильмы в JSON.
// // - если indent = 0 - использует json.Marshal
// // - если indent > 0 - использует json.MarshalIndent
// //   с отступом в указанное количество пробелов.
// func MarshalMovies(indent int, movies ...Movie) (string, error) {
// 	if indent == 0 {
// 		b, err := json.Marshal(movies)
// 		if err != nil {
// 			return "", err
// 		}
// 		return string(b), nil
// 	} else {
// 		b, err := json.MarshalIndent(movies, "", strings.Repeat(" ", indent))
// 		if err != nil {
// 			return "", err
// 		}
// 		return string(b), nil
// 	}
// }

// // конец решения

// func main() {
// 	m1 := Movie{
// 		Title:    "Interstellar",
// 		Year:     2014,
// 		Director: "Christopher Nolan",
// 		Genres:   []string{"Adventure", "Drama", "Science Fiction"},
// 		Duration: Duration(2*time.Hour + 49*time.Minute),
// 		Rating:   5,
// 	}
// 	m2 := Movie{
// 		Title:    "Sully",
// 		Year:     2016,
// 		Director: "Clint Eastwood",
// 		Genres:   []string{"Drama", "History"},
// 		Duration: Duration(time.Hour + 36*time.Minute),
// 		Rating:   4,
// 	}

// 	s, err := MarshalMovies(4, m1, m2)
// 	fmt.Println(err)
// 	// nil
// 	fmt.Println(s)
// 	/*
// 		[
// 		    {
// 		        "Title": "Interstellar",
// 		        "Year": 2014,
// 		        "Director": "Christopher Nolan",
// 		        "Genres": [
// 		            "Adventure",
// 		            "Drama",
// 		            "Science Fiction"
// 		        ],
// 		        "Duration": "2h49m",
// 		        "Rating": "★★★★★"
// 		    },
// 		    {
// 		        "Title": "Sully",
// 		        "Year": 2016,
// 		        "Director": "Clint Eastwood",
// 		        "Genres": [
// 		            "Drama",
// 		            "History"
// 		        ],
// 		        "Duration": "1h36m",
// 		        "Rating": "★★★★☆"
// 		    }
// 		]
// 	*/
// }

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// )

// // начало решения

// // Genre описывает жанр фильма
// type Genre string

// func (an *Genre) UnmarshalJSON(data []byte) error {
//     // Go рекомендует игнорировать значения null
//     if string(data) == "null" {
//         return nil
//     }
//     // декодируем исходное число
// 	var tempGenre struct{
// 		Name string `json:"name"`
// 	}

//     err := json.Unmarshal(data, &tempGenre)
//     if err != nil {
// 		fmt.Println(err)
//         return err
//     }
// 	*an = Genre(tempGenre.Name)
//     // преобразуем в значение типа AncientNumber
//     return nil
// }

// // Movie описывает фильм
// type Movie struct {
// 	Title  string `json:"name"`
// 	Year   int `json:"released_at"`
// 	Genres []Genre `json:"tags"`
// }

// // конец решения

// func main() {
// 	const src = `{
// 		"name": "Interstellar",
// 		"released_at": 2014,
// 		"director": "Christopher Nolan",
// 		"tags": [
// 			{ "name": "Adventure" },
// 			{ "name": "Drama" },
// 			{ "name": "Science Fiction" }
// 		],
// 		"duration": "2h49m",
// 		"rating": "★★★★★"
// 	}`

// 	var m Movie
// 	err := json.Unmarshal([]byte(src), &m)
// 	fmt.Println(err)
// 	// nil
// 	fmt.Println(m)
// 	// {Interstellar 2014 [Adventure Drama Science Fiction]}
// }

// package main

// import (
// 	"fmt"
// 	"io"
// 	"os"
// 	"strings"
// 	"encoding/json"
// )

// // начало решения

// // Email описывает письмо
// type Email struct {
// 	From string `json:"from"`
// 	To string `json:"to"`
// 	Subject string `json:"subject"`
// }

// // FilterEmails читает все письма из src и записывает в dst тех,
// // кто проходит проверку predicate
// func FilterEmails(dst io.Writer, src io.Reader, predicate func(e Email) bool) (int, error) {
// 	emailCount := 0
// 	en := json.NewEncoder(dst)
// 	dec := json.NewDecoder(src)
// 	for {
// 		var email Email
// 		err := dec.Decode(&email)
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return emailCount, err
// 		}
// 		if predicate(email) {
// 			if err := en.Encode(email); err != nil {
// 				return emailCount, err
// 			}
// 			emailCount++
// 		}
// 	}
// 	return emailCount, nil
// }

// // конец решения

// func main() {
// 	src := strings.NewReader(`
// 		{ "from": "alice@go.dev",      "to": "zet@php.net",              "subject": "How are you?" }
// 		{ "from": "bob@temp-mail.org", "to": "yolanda@java.com",         "subject": "Re: Indonesia" }
// 		{ "from": "cindy@go.dev",      "to": "xavier@rust-lang.org",     "subject": "Go vs Rust" }
// 		{ "from": "diane@dart.dev",    "to": "wanda@typescriptlang.org", "subject": "Our crypto startup" }
// 	`)
// 	dst := os.Stdout

// 	predicate := func(email Email) bool {
// 		return !strings.Contains(email.Subject, "crypto")
// 	}

// 	n, err := FilterEmails(dst, src, predicate)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(n, "good emails")

// 	// {"from":"alice@go.dev","to":"zet@php.net","subject":"How are you?"}
// 	// {"from":"bob@temp-mail.org","to":"yolanda@java.com","subject":"Re: Indonesia"}
// 	// {"from":"cindy@go.dev","to":"xavier@rust-lang.org","subject":"Go vs Rust"}
// 	// 3 good emails
// }

// package main

// import (
// 	"io"
// 	"os"
// 	"strings"
// 	"encoding/xml"
// 	"encoding/csv"
// 	"strconv"
// )

// // начало решения
// type Employee struct {
// 	Id int `xml:"id,attr"`
// 	Name string `xml:"name"`
// 	City string `xml:"city"`
// 	Salary int `xml:"salary"`

// }
// type Department struct {
// 	Code string `xml:"code"`
// 	Employees []Employee `xml:"employees>employee"`
// }

// type Organization struct {
// 	Departments []Department `xml:"department"`
// }

// func (e *Employee) Slice(depCode string) []string {
// 	return []string{strconv.Itoa(e.Id), e.Name, e.City, depCode, strconv.Itoa(e.Salary)}
// }
// // ConvertEmployees преобразует XML-документ с информацией об организации
// // в плоский CSV-документ с информацией о сотрудниках
// func ConvertEmployees(outCSV io.Writer, inXML io.Reader) error {
// 	var org Organization
// 	dec := xml.NewDecoder(inXML)
// 	if err := dec.Decode(&org); err != nil {
// 		return err
// 	}
// 	w := csv.NewWriter(outCSV)
// 	w.Write([]string{"id", "name" ,"city", "department" ,"salary"})
// 	for _, d := range org.Departments {
// 		for _, e := range d.Employees {
// 			err := w.Write(e.Slice(d.Code))
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	w.Flush()
// 	if w.Error() != nil {
// 		return w.Error()
// 	}
// 	return nil
// }

// // конец решения

// func main() {
// 	src := `<organization>
//     <department>
//         <code>hr</code>
//         <employees>
//             <employee id="11">
//                 <name>Дарья</name>
//                 <city>Самара</city>
//                 <salary>70</salary>
//             </employee>
//             <employee id="12">
//                 <name>Борис</name>
//                 <city>Самара</city>
//                 <salary>78</salary>
//             </employee>
//         </employees>
//     </department>
//     <department>
//         <code>it</code>
//         <employees>
//             <employee id="21">
//                 <name>Елена</name>
//                 <city>Самара</city>
//                 <salary>84</salary>
//             </employee>
//         </employees>
//     </department>
// </organization>`

// 	in := strings.NewReader(src)
// 	out := os.Stdout
// 	ConvertEmployees(out, in)
// 	/*
// 		id,name,city,department,salary
// 		11,Дарья,Самара,hr,70
// 		12,Борис,Самара,hr,78
// 		21,Елена,Самара,it,84
// 	*/
// }

// package main

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/url"
// 	"time"
// 	"encoding/json"
// )

// // StatusErr описывает ситуацию, когда на запрос
// // пришел ответ с HTTP-статусом, отличным от 2xx.
// type StatusErr struct {
// 	Code   int
// 	Status string
// }

// func (e StatusErr) Error() string {
// 	return "invalid response status: " + e.Status
// }

// // начало решения

// // httpGet выполняет GET-запрос с указанными заголовками и параметрами,
// // парсит ответ как JSON и возвращает получившуюся карту.
// //
// // Считает ошибкой любые ответы с HTTP-статусом, отличным от 2xx.
// func httpGet(uri string, headers map[string]string, params map[string]string, timeout int) (map[string]any, error) {
// 	client := http.Client{Timeout: time.Duration(timeout)*time.Millisecond}
// 	req, err := http.NewRequest(http.MethodGet, uri, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(headers) > 0 {
// 		for key, val := range headers {
// 			req.Header.Add(key, val)
// 		}
// 	}
// 	if len(params) > 0 {
// 		reqParams := url.Values{}
// 		for key, val := range params {
// 			reqParams.Add(key, val)
// 		}
// 		req.URL.RawQuery = reqParams.Encode()
// 	}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	code := resp.StatusCode
// 	if code >= 200 && code <= 299 {
// 		return nil, StatusErr{Code: resp.StatusCode, Status: resp.Status[:4]}
// 	}
// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var data map[string]any
// 	if err := json.Unmarshal(body, &data); err != nil {
// 		return nil, err
// 	}

// 	return data, nil
// }

// // конец решения

// func main() {
// 	{
// 		// GET-запрос
// 		const uri = "https://httpbingo.org/json"
// 		data, err := httpGet(uri, nil, nil, 3000)
// 		fmt.Printf("GET %v\n", uri)
// 		fmt.Println(data, err)
// 		fmt.Println()
// 		// GET https://httpbingo.org/json
// 		// map[slideshow:map[author:Yours Truly date:date of publication slides:[map[title:Wake up to WonderWidgets! type:all] map[items:[Why <em>WonderWidgets</em> are great Who <em>buys</em> WonderWidgets] title:Overview type:all]] title:Sample Slide Show]] <nil>
// 	}

// 	{
// 		// 404 Not Found
// 		const uri = "https://httpbingo.org/whatever"
// 		data, err := httpGet(uri, nil, nil, 3000)
// 		fmt.Printf("GET %v\n", uri)
// 		fmt.Println(data, err)
// 		fmt.Println()
// 		// GET https://httpbingo.org/whatever
// 		// map[] invalid response status: 404 Not Found
// 	}

// 	{
// 		// С заголовками
// 		const uri = "https://httpbingo.org/headers"
// 		headers := map[string]string{
// 			"accept": "application/xml",
// 			"host":   "httpbingo.org",
// 		}
// 		data, err := httpGet(uri, headers, nil, 3000)
// 		fmt.Printf("GET %v\n", uri)
// 		respHeaders := data["headers"].(map[string]any)
// 		fmt.Println(respHeaders["Accept"], respHeaders["Host"], err)
// 		fmt.Println()
// 		// GET https://httpbingo.org/headers
// 		// [application/xml] [httpbingo.org] <nil>
// 	}

// 	{
// 		// С URL-параметрами
// 		const uri = "https://httpbingo.org/get"
// 		params := map[string]string{"id": "42"}
// 		data, err := httpGet(uri, nil, params, 3000)
// 		fmt.Printf("GET %v\n", uri)
// 		fmt.Println(data["args"], err)
// 		fmt.Println()
// 		// GET https://httpbingo.org/get
// 		// map[id:[42]] <nil>
// 	}
// }

// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/url"
// )

// // начало решения

// // Handy предоставляет удобный интерфейс
// // для выполнения HTTP-запросов
// type Handy struct{
// 	req *http.Request
// 	client *http.Client
// 	url string
// 	err error
// }

// // NewHandy создает новый экземпляр Handy
// func NewHandy() *Handy {
// 	request, err := http.NewRequest("", "", nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return &Handy{request, http.DefaultClient, "", nil}
// }

// // URL устанавливает URL, на который пойдет запрос
// func (h *Handy) URL(uri string) *Handy {
// 	parsedURL, err := url.Parse(uri)
// 	if err != nil {
// 		panic(err)
// 	}
// 	h.req.URL = parsedURL
// 	h.url = uri
// 	return h
// }

// // Client устанавливает HTTP-клиента
// // вместо умолчательного http.DefaultClient
// func (h *Handy) Client(client *http.Client) *Handy {
// 	h.client = client
// 	return h
// }

// // Header устанавливает значение заголовка
// func (h *Handy) Header(key, value string) *Handy {
// 	h.req.Header.Add(key, value)
// 	return h
// }

// // Param устанавливает значение URL-параметра
// func (h *Handy) Param(key, value string) *Handy {
// 	q := h.req.URL.Query()
//     q.Add(key, value)
//     h.req.URL.RawQuery = q.Encode()
// 	return h
// }

// // Form устанавливает данные, которые будут закодированы
// // как application/x-www-form-urlencoded и отправлены в теле запроса
// // с соответствующим content-type
// func (h *Handy) Form(form map[string]string) *Handy {
// 	h.req.Header = make(http.Header)
// 	data := url.Values{}
// 	for key, val := range form {
// 		data.Add(key, val)
// 	}
// 	bodyStr := data.Encode()
// 	h.req.Body = io.NopCloser(bytes.NewReader([]byte(bodyStr)))
// 	h.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	return h
// }

// // JSON устанавливает данные, которые будут закодированы
// // как application/json и отправлены в теле запроса
// // с соответствующим content-type
// func (h *Handy) JSON(v any) *Handy {
// 	h.req.Header = make(http.Header)
// 	for key := range h.req.Header {
// 		h.req.Header.Del(key)
// 	}
// 	bodyBytes, err := json.Marshal(v)
// 	if err != nil {
// 		h.err = err
// 	}
// 	h.req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
// 	h.req.Header.Set("Content-Type", "application/json") 
// 	return h
// }

// // Get выполняет GET-запрос с настроенными ранее параметрами
// func (h *Handy) Get() *HandyResponse {
// 	if h.err != nil {
// 		return &HandyResponse{0, nil, h.err, nil}
// 	}
// 	h.req.Method = http.MethodGet
// 	resp, err := h.client.Do(h.req)
// 	if resp == nil {
// 		return &HandyResponse{0, nil, err, nil}
// 	}
// 	return &HandyResponse{resp.StatusCode, resp, err, nil}
// }

// // Post выполняет POST-запрос с настроенными ранее параметрами
// func (h *Handy) Post() *HandyResponse {
// 	if h.err != nil {
// 		return &HandyResponse{0, nil, h.err, nil}
// 	}
// 	h.req.Method = http.MethodPost
// 	resp, err := h.client.Do(h.req)
// 	if resp == nil {
// 		return &HandyResponse{0, nil, err, nil}
// 	}
// 	return &HandyResponse{resp.StatusCode, resp, err, nil}
// }

// // HandyResponse представляет ответ на HTTP-запрос
// type HandyResponse struct {
// 	StatusCode int
// 	Resp *http.Response
// 	Error error
// 	bodyBytes  []byte
// }

// // OK возвращает true, если во время выполнения запроса
// // не произошло ошибок, а код HTTP-статуса ответа равен 200
// func (r *HandyResponse) OK() bool {
// 	if (r.Error != nil || r.Resp == nil){
// 		return false
// 	}
// 	return r.StatusCode == http.StatusOK
// }

// // Bytes возвращает тело ответа как срез байт
// func (r *HandyResponse) Bytes() []byte {
// 	if r.bodyBytes != nil {
//         return r.bodyBytes
//     }
// 	if r.Resp == nil || r.Resp.Body == nil {
//         return nil
//     }
// 	defer r.Resp.Body.Close()
// 	bytesResp, err := io.ReadAll(r.Resp.Body)
// 	r.Error = err
// 	r.bodyBytes = bytesResp
// 	return bytesResp
// }

// // String возвращает тело ответа как строку
// func (r *HandyResponse) String() string {
// 	return string(r.Bytes())
// }

// // JSON декодирует тело ответа из JSON и сохраняет
// // результат по адресу, на который указывает v
// func (r *HandyResponse) JSON(v any) {
// 	// работает аналогично json.Unmarshal()
// 	// если при декодировании произошла ошибка,
// 	// она должна быть доступна через r.Err()
// 	defer r.Resp.Body.Close()
// 	bytesResp := r.Bytes()
// 	if len(bytesResp) == 0 {
// 		return
// 	}
// 	err := json.Unmarshal(bytesResp, &v)
// 	r.Error = err
// }

// // Err возвращает ошибку, которая возникла при выполнении запроса
// // или обработке ответа
// func (r *HandyResponse) Err() error {
// 	return r.Error
// }

// // конец решения

// func main() {
// 	{
// 		// примеры запросов

// 		// GET-запрос с параметрами
// 		NewHandy().URL("https://httpbingo.org/get").Param("id", "42").Get()

// 		// HTTP-заголовки
// 		NewHandy().
// 			URL("https://httpbingo.org/get").
// 			Header("Accept", "text/html").
// 			Header("Authorization", "Bearer 1234567890").
// 			Get()

// 		// POST формы
// 		params := map[string]string{
// 			"brand":    "lg",
// 			"category": "tv",
// 		}
// 		NewHandy().URL("https://httpbingo.org/post").Form(params).Post()

// 		// POST JSON-документа
// 		NewHandy().URL("https://httpbingo.org/post").JSON(params).Post()
// 	}

// 	{
// 		// пример обработки ответа

// 		// отправляем GET-запрос с параметрами
// 		resp := NewHandy().URL("https://httpbingo.org/get").Param("id", "42").Get()
// 		if !resp.OK() {
// 			panic(resp.String())
// 		}

// 		// декодируем ответ в JSON
// 		var data map[string]any
// 		resp.JSON(&data)

// 		fmt.Println(data["url"])
// 		// "https://httpbingo.org/get"
// 		fmt.Println(data["args"])
// 		// map[id:[42]]
// 	}
// }

// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"strconv"
// )

// // начало решения

// // statusHandler возвращает ответ с кодом, который передан
// // в заголовке X-Status. Например:
// //
// //	X-Status = 200 -> вернет ответ с кодом 200
// //	X-Status = 404 -> вернет ответ с кодом 404
// //	X-Status = 503 -> вернет ответ с кодом 503
// //
// // Если заголовок отстутствует, возвращает ответ с кодом 200.
// // Тело ответа пустое.
// func statusHandler(w http.ResponseWriter, r *http.Request) {
// 	header := r.Header.Get("X-Status")
// 	if header == "" {
// 		w.WriteHeader(http.StatusOK)
// 	}
// 	code, err := strconv.Atoi(header)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	w.WriteHeader(code)
// }

// // echoHandler возвращает ответ с тем же телом
// // и заголовком Content-Type, которые пришли в запросе
// func echoHandler(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	reqBody, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	reqHeaderCT := r.Header.Get("Content-Type")  
// 	w.Header().Set("Content-Type", reqHeaderCT)
// 	w.Write(reqBody)
// }

// // jsonHandler проверяет, что Content-Type = application/json,
// // а в теле запроса пришел валидный JSON,
// // после чего возвращает ответ с кодом 200.
// // Если какая-то проверка не прошла — возвращает ответ с кодом 400.
// // Тело ответа пустое.
// func jsonHandler(w http.ResponseWriter, r *http.Request) {
// 	reqHeaderCT := r.Header.Get("Content-Type")
// 	if reqHeaderCT != "application/json" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close()
// 	reqJSON, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	if !json.Valid(reqJSON){
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// }

// // конец решения

// func startServer() *httptest.Server {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/status", statusHandler)
// 	mux.HandleFunc("/echo", echoHandler)
// 	mux.HandleFunc("/json", jsonHandler)
// 	return httptest.NewServer(mux)
// }

// func main() {
// 	server := startServer()
// 	defer server.Close()
// 	client := server.Client()

// 	{
// 		uri := server.URL + "/status"
// 		req, _ := http.NewRequest(http.MethodGet, uri, nil)
// 		req.Header.Add("X-Status", "201")
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(resp.Status)
// 		// 201 Created
// 	}

// 	{
// 		uri := server.URL + "/echo"
// 		reqBody := []byte("hello world")
// 		resp, err := client.Post(uri, "text/plain", bytes.NewReader(reqBody))
// 		if err != nil {
// 			panic(err)
// 		}

// 		defer resp.Body.Close()
// 		respBody, _ := io.ReadAll(resp.Body)
// 		fmt.Println(resp.Status)
// 		fmt.Println(string(respBody))
// 		// 200 OK
// 		// hello world
// 	}

// 	{
// 		uri := server.URL + "/json"
// 		reqBody, _ := json.Marshal(map[string]bool{"ok": true})
// 		resp, err := client.Post(uri, "application/json", bytes.NewReader(reqBody))
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(resp.Status)
// 		// 200 OK
// 	}
// }

// package main

// import (
// 	"database/sql"
// 	"fmt"

// 	_ "github.com/mattn/go-sqlite3"
// )

// // начало решения

// // SQLMap представляет карту, которая хранится в SQL-базе данных
// type SQLMap struct{
// 	Map map[string]any
// 	db *sql.DB
// }

// // NewSQLMap создает новую SQL-карту в указанной базе
// func NewSQLMap(db *sql.DB) (*SQLMap, error) {
// 	query := `create table if not exists map(key text primary key, val blob)`
// 	_, err := db.Exec(query)
// 	return &SQLMap{db: db}, err
// }

// // Get возвращает значение для указанного ключа.
// // Если такого ключа нет - возвращает ошибку sql.ErrNoRows.
// func (m *SQLMap) Get(key string) (any, error) {
// 	query := `select val from map where key = ?`
// 	row := m.db.QueryRow(query, key)
// 	var val any
// 	err := row.Scan(&val)
// 	if err == sql.ErrNoRows {
// 		return nil, err
// 	} else if err != nil {
// 		return nil, err
// 	}
// 	return val, nil
// }

// // Set устанавливает значение для указанного ключа.
// // Если такой ключ уже есть - затирает старое значение (это не считается ошибкой).
// func (m *SQLMap) Set(key string, val any) error {
// 	query := `insert into map(key, val) values (?, ?)
// on conflict (key) do update set val = excluded.val`
// 	_, err := m.db.Exec(query, key, val)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Delete удаляет запись карты с указанным ключом.
// // Если такого ключа нет - ничего не делает (это не считается ошибкой).
// func (m *SQLMap) Delete(key string) error {
// 	query := `delete from map where key = ?`
// 	res, err := m.db.Exec(query, key)
// 	if err != nil {
// 		return err
// 	}
// 	if count, _ := res.RowsAffected(); count == 0 {
// 		return nil
// 	}
// 	return nil
// }

// // конец решения

// func main() {
// 	db, err := sql.Open("sqlite3", ":memory:")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	m, err := NewSQLMap(db)
// 	if err != nil {
// 		panic(err)
// 	}

// 	m.Set("name", "Alice")
// 	m.Set("age", 42)

// 	name, err := m.Get("name")
// 	fmt.Printf("name = %v, err = %v\n", name, err)
// 	// name = Alice, err = <nil>

// 	age, err := m.Get("age")
// 	fmt.Printf("age = %v, err = %v\n", age, err)
// 	// age = 42, err = <nil>

// 	m.Set("name", "Bob")
// 	name, err = m.Get("name")
// 	fmt.Printf("name = %v, err = %v\n", name, err)
// 	// name = Bob, err = <nil>

// 	m.Delete("name")
// 	name, err = m.Get("name")
// 	fmt.Printf("name = %v, err = %v\n", name, err)
// 	// name = <nil>, err = sql: no rows in result set
// }

// package main

// import (
// 	"database/sql"
// 	"fmt"

// 	_ "github.com/mattn/go-sqlite3"
// )

// //начало решения

// //SQLMap представляет карту, которая хранится в SQL-базе данных
// // SQLMap представляет карту, которая хранится в SQL-базе данных
// type SQLMap struct{
// 	Map map[string]any
// 	db *sql.DB
// 	GetStmt *sql.Stmt
// 	SetStmt *sql.Stmt
// 	DelStmt *sql.Stmt
// }

// // NewSQLMap создает новую SQL-карту в указанной базе
// func NewSQLMap(db *sql.DB) (*SQLMap, error) {
// 	query := `create table if not exists map(key text primary key, val blob)`
// 	_, err := db.Exec(query)
// 	get, err := db.Prepare(`select val from map where key = ?`) 
// 	if err != nil {
// 		return nil, err
// 	}
// 	set, err := db.Prepare(`insert into map(key, val) values (?, ?)
// on conflict (key) do update set val = excluded.val`) 
// 	if err != nil {
// 		return nil, err
// 	}
// 	del, err := db.Prepare(`delete from map where key = ?`) 
// 	if err != nil {
// 		return nil, err
// 	}
	
// 	return &SQLMap{
// 		Map: map[string]any{},
// 		db: db,
// 		GetStmt: get,
// 		SetStmt: set,
// 		DelStmt: del,
// 	}, err
// }

// // Get возвращает значение для указанного ключа.
// // Если такого ключа нет - возвращает ошибку sql.ErrNoRows.
// func (m *SQLMap) Get(key string) (any, error) {
// 	row := m.GetStmt.QueryRow(key)
// 	var val any
// 	err := row.Scan(&val)
// 	if err == sql.ErrNoRows {
// 		return nil, err
// 	} else if err != nil {
// 		return nil, err
// 	}
// 	return val, nil
// }

// // Set устанавливает значение для указанного ключа.
// // Если такой ключ уже есть - затирает старое значение (это не считается ошибкой).
// func (m *SQLMap) Set(key string, val any) error {
// 	_, err := m.SetStmt.Exec(key, val)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Delete удаляет запись карты с указанным ключом.
// // Если такого ключа нет - ничего не делает (это не считается ошибкой).
// func (m *SQLMap) Delete(key string) error {
// 	res, err := m.DelStmt.Exec(key)
// 	if err != nil {
// 		return err
// 	}
// 	if count, _ := res.RowsAffected(); count == 0 {
// 		return nil
// 	}
// 	return nil
// }

// // SetItems устанавливает значения указанных ключей.
// func (m *SQLMap) SetItems(items map[string]any) error {
// 	tx, err := m.db.Begin()
// 	if err != nil {
//         return err
//     }
// 	defer tx.Rollback()

// 	txStmt := tx.Stmt(m.SetStmt)
// 	for key, val := range items {
// 		_, err := txStmt.Exec(key, val)
// 	if err != nil {
// 		return err
// 	}
// 	}
// 	return tx.Commit()
// }

// // Close освобождает ресурсы, занятые картой в базе.
// func (m *SQLMap) Close() error {
// 	err := m.GetStmt.Close()
// 	if err != nil {
// 		return err
// 	}
// 	err = m.SetStmt.Close()
// 	if err != nil {
// 		return err
// 	}
// 	err = m.DelStmt.Close()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // конец решения

// func main() {
// 	db, err := sql.Open("sqlite3", ":memory:")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	m, err := NewSQLMap(db)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer m.Close()

// 	m.Set("name", "Alice")

// 	items := map[string]any{
// 		"name": "Bob",
// 		"age":  42,
// 	}
// 	m.SetItems(items)

// 	name, err := m.Get("name")
// 	fmt.Printf("name = %v, err = %v\n", name, err)
// 	// name = Bob, err = <nil>

// 	age, err := m.Get("age")
// 	fmt.Printf("age = %v, err = %v\n", age, err)
// 	// age = 42, err = <nil>
// }

package main

import (
	"database/sql"
	"time"
	"context"

	_ "github.com/mattn/go-sqlite3"
)

//начало решения

//SQLMap представляет карту, которая хранится в SQL-базе данных
// SQLMap представляет карту, которая хранится в SQL-базе данных
type SQLMap struct{
	Map map[string]any
	db *sql.DB
	GetStmt *sql.Stmt
	SetStmt *sql.Stmt
	DelStmt *sql.Stmt
	Timeout time.Duration
}

// NewSQLMap создает новую SQL-карту в указанной базе
func NewSQLMap(db *sql.DB) (*SQLMap, error) {
	query := `create table if not exists map(key text primary key, val blob)`
	_, err := db.Exec(query)
	if err != nil {
		return nil, err
	}
	get, err := db.Prepare(`select val from map where key = ?`) 
	if err != nil {
		return nil, err
	}
	set, err := db.Prepare(`insert into map(key, val) values (?, ?)
on conflict (key) do update set val = excluded.val`) 
	if err != nil {
		return nil, err
	}
	del, err := db.Prepare(`delete from map where key = ?`) 
	if err != nil {
		return nil, err
	}
	return &SQLMap{
		Map: map[string]any{},
		db: db,
		GetStmt: get,
		SetStmt: set,
		DelStmt: del,
		Timeout: 60 * time.Second,
	}, nil
}

// SetTimeout устанавливает максимальное время выполнения
// отдельного метода карты.
func (m *SQLMap) SetTimeout(d time.Duration) {
	m.Timeout = d
}

// Get возвращает значение для указанного ключа.
// Если такого ключа нет - возвращает ошибку sql.ErrNoRows.
func (m *SQLMap) Get(key string) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()

	row := m.GetStmt.QueryRowContext(ctx, key)
	var val any
	err := row.Scan(&val)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

// Set устанавливает значение для указанного ключа.
// Если такой ключ уже есть - затирает старое значение (это не считается ошибкой).
func (m *SQLMap) Set(key string, val any) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()
	_, err := m.SetStmt.ExecContext(ctx, key, val)
	if err != nil {
		return err
	}
	return nil
}

// Delete удаляет запись карты с указанным ключом.
// Если такого ключа нет - ничего не делает (это не считается ошибкой).
func (m *SQLMap) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()
	res, err := m.DelStmt.ExecContext(ctx, key)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return nil
	}
	return nil
}

// SetItems устанавливает значения указанных ключей.
func (m *SQLMap) SetItems(items map[string]any) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
        return err
    }
	defer tx.Rollback()

	txStmt := tx.StmtContext(ctx, m.SetStmt)
	for key, val := range items {
		_, err := txStmt.ExecContext(ctx, key, val)
	if err != nil {
		return err
	}
	}
	return tx.Commit()
}

// Close освобождает ресурсы, занятые картой в базе.
func (m *SQLMap) Close() error {
	err := m.GetStmt.Close()
	if err != nil {
		return err
	}
	err = m.SetStmt.Close()
	if err != nil {
		return err
	}
	err = m.DelStmt.Close()
	if err != nil {
		return err
	}
	return nil
}

// конец решения

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	m, err := NewSQLMap(db)
	if err != nil {
		panic(err)
	}
	defer m.Close()

	m.SetTimeout(10 * time.Millisecond)

	m.Set("name", "Alice")
	m.Get("name")
}
