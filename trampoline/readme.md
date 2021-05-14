# Trampoline
Паттерн Trampoline лёгким движением руки превращает рекурсию в цикл.

# Зачем?
Шаблон «Trampoline» допускает использовать рекурсию без исчерпания стековой памяти.
К примеру при расчёте факториала числа 5000 python создаст исключение:
```
RecursionError: maximum recursion depth exceeded in comparison
```
"Trampoline" позволит обойти эту проблему.