CREATE TYPE task_type AS ENUM (
    'input',
    'choice_one',
    'choice_many',
    'code',
    'match'
);

ALTER TABLE tasks
ALTER COLUMN type TYPE task_type
USING type::task_type;

INSERT INTO tasks (level_id, type, question, correct_answer, xp_reward)
VALUES (1, 'input', 'Напиши "hello world"', 'hello world', 5);

INSERT INTO tasks (level_id, type, question, xp_reward)
VALUES (1, 'choice_one', 'Выбери правильный вариант', 10);

INSERT INTO task_variants (task_id, content, is_correct)
VALUES 
(6, 'Неправильно', false),
(6, 'Правильно', true),
(6, 'Тоже мимо', false);