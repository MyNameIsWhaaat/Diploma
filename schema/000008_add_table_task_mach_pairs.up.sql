CREATE TABLE task_match_pairs (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    left_text TEXT NOT NULL,
    right_text TEXT NOT NULL,
    match_key TEXT NOT NULL
);

INSERT INTO tasks (level_id, type, question, xp_reward)
VALUES (16, 'match', 'Сопоставь элементы алгоритма с их описанием', 25);

-- допустим, задача получила id = 9
INSERT INTO task_match_pairs (task_id, left_text, right_text, match_key) VALUES
(10, 'Начало - Конец', 'Обозначение границ алгоритма', '1'),
(10, 'Пошаговость', 'Каждое действие записано последовательно', '2'),
(10, 'Ветвление', 'Условный выбор одного из путей', '3'),
(10, 'Цикл', 'Повторяющееся действие', '4');
