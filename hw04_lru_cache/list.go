package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value any       // Значение
	Next  *ListItem // Следующий элемент
	Prev  *ListItem // Предыдущий элемент
}

type linkedList struct {
	len   int       // Количество элементов
	front *ListItem // Начало
	back  *ListItem // Конец
}

// Создать новый список.
func NewList() List {
	return new(linkedList)
}

// Получить длину списка.
func (l *linkedList) Len() int {
	return l.len
}

// Возвращает первый элемент списка linkedList.
func (l *linkedList) Front() *ListItem {
	return l.front
}

// Возвращает последний элемент списка linkedList.
func (l *linkedList) Back() *ListItem {
	return l.back
}

// Поместить значение в начало списка.
func (l *linkedList) PushFront(v any) *ListItem {
	newNode := &ListItem{
		Value: v,
		Prev:  nil,
		Next:  l.front,
	}

	if l.front == nil {
		l.back = newNode
	} else {
		l.front.Prev = newNode
	}
	l.front = newNode
	l.len++
	return newNode
}

// Поместить значение в конец списка.
func (l *linkedList) PushBack(v any) *ListItem {
	newNode := &ListItem{
		Value: v,
		Prev:  l.back,
		Next:  nil,
	}

	if l.back == nil {
		l.front = newNode
	} else {
		l.back.Next = newNode
	}
	l.back = newNode
	l.len++
	return newNode
}

// Удалить значение из списка.
func (l *linkedList) Remove(e *ListItem) {
	if e == nil {
		return
	}

	if e.Next != nil {
		e.Next.Prev = e.Prev
	} else {
		l.back = e.Prev
	}

	if e.Prev != nil {
		e.Prev.Next = e.Next
	} else {
		l.front = e.Next
	}

	l.len--
}

// Переместить значение в начало списка.
func (l *linkedList) MoveToFront(e *ListItem) {
	l.Remove(e)
	l.PushFront(e.Value)
}
