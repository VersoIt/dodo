package common

// Money - денежный тип для всей системы.
// Используем float64 для соответствия формулам в ТЗ (деление, проценты).
type Money float64

// Метод для округления до 2 знаков (для красоты в JSON/БД)
func (m Money) Round() float64 {
	return float64(int(m*100)) / 100
}
