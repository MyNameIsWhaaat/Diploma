CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,             -- Отображаемое имя
    username VARCHAR(255) NOT NULL UNIQUE,  -- Уникальный логин
    password_hash VARCHAR(255) NOT NULL,    -- Хэш пароля
    email VARCHAR(255) UNIQUE,              -- Почта (для восстановления)
    avatar_url TEXT,                        -- Ссылка на аватар
    role VARCHAR(50) DEFAULT 'student',     -- Роль пользователя (student, admin)
    is_active BOOLEAN DEFAULT TRUE,         -- Аккаунт активен или нет
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,                   -- Название курса
    short_description TEXT NOT NULL,               -- Краткое описание
    full_description TEXT,                         -- Подробное описание
    image_url TEXT,                                -- Картинка (обложка)
    difficulty VARCHAR(50),                        -- Сложность (beginner/intermediate/advanced)
    xp_reward INTEGER DEFAULT 0,                   -- Сколько XP за курс
    is_active BOOLEAN DEFAULT TRUE,                -- Статус курса
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE levels (
    id SERIAL PRIMARY KEY,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE, -- К какому курсу относится
    title VARCHAR(255) NOT NULL,                 -- Название уровня
    description TEXT,                            -- Описание (по желанию)
    order_index INTEGER DEFAULT 0,               -- Порядок отображения
    xp_reward INTEGER DEFAULT 0,                 -- Сколько XP за прохождение уровня
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    level_id INTEGER NOT NULL REFERENCES levels(id) ON DELETE CASCADE, -- К какому уровню принадлежит
    type VARCHAR(50) NOT NULL,              -- Тип задания (input, choice_one, choice_many, code, match и т.д.)
    question TEXT NOT NULL,                 -- Текст задания или JSON-описание
    correct_answer TEXT,                    -- Правильный ответ (если применимо)
    xp_reward INTEGER DEFAULT 0,            -- Очки за выполнение задания
    order_index INTEGER DEFAULT 0,          -- Порядок в уровне
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE task_variants (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE, -- К какому заданию принадлежит
    content TEXT NOT NULL,                    -- Текст или содержимое варианта
    is_correct BOOLEAN DEFAULT FALSE,         -- Правильный ли это ответ
    order_index INTEGER DEFAULT 0             -- Для сортировки/перетаскивания
);

CREATE TABLE user_courses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    xp_earned INTEGER DEFAULT 0,                  -- Сколько XP пользователь заработал в этом курсе
    current_level_id INTEGER REFERENCES levels(id), -- Текущий уровень, на котором остановился
    completed BOOLEAN DEFAULT FALSE,              -- Завершён ли курс
    started_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE (user_id, course_id)                   -- Один пользователь — один курс
);

CREATE TABLE user_levels (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    level_id INTEGER NOT NULL REFERENCES levels(id) ON DELETE CASCADE,
    completed BOOLEAN DEFAULT FALSE,            -- Завершён ли уровень
    is_current BOOLEAN DEFAULT FALSE,           -- Является ли текущим активным уровнем
    xp_earned INTEGER DEFAULT 0,                -- Сколько XP получил пользователь на этом уровне
    started_at TIMESTAMP DEFAULT now(),
    completed_at TIMESTAMP,                     -- Время завершения, если завершён
    UNIQUE (user_id, level_id)
);

CREATE TABLE progress (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    is_completed BOOLEAN DEFAULT FALSE,            -- Выполнено ли задание
    last_answer TEXT,                              -- Последний ответ пользователя (если применимо)
    attempts INTEGER DEFAULT 0,                    -- Количество попыток
    xp_earned INTEGER DEFAULT 0,                   -- Получено XP
    time_spent INTEGER,                            -- Время выполнения (в секундах, опционально)
    completed_at TIMESTAMP,                        -- Дата выполнения (если завершено)
    needs_review BOOLEAN DEFAULT FALSE,            -- Флаг: нужно повторить (ошибка была)
    UNIQUE (user_id, task_id)
);

CREATE TABLE profile_levels (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	min_xp INT NOT NULL,
	max_xp INT NOT NULL
);

CREATE TABLE user_profile_levels (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	profile_level_id INT NOT NULL REFERENCES profile_levels(id),
	total_xp INT NOT NULL DEFAULT 0,
	last_level_up TIMESTAMP,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);



INSERT INTO profile_levels (title, min_xp, max_xp) VALUES
('Новичок', 0, 99),
('Ученик', 100, 299),
('Профессионал', 300, 599),
('Мастер', 600, 999),
('Гуру', 1000, 1999),
('Легенда', 2000, 999999);

INSERT INTO courses (title, short_description, full_description, image_url, difficulty, xp_reward, is_active)
VALUES
('Алгоритмы для начинающих', 'Базовые алгоритмы', 'Описание полного курса по базовым алгоритмам', 'https://example.com/img1.png', 'beginner', 100, true),
('Структуры данных', 'Массивы, списки и деревья', 'Углублённое изучение структур данных', 'https://example.com/img2.png', 'intermediate', 150, true),
('Графы и динамика', 'Алгоритмы на графах и динамическое программирование', 'Курс по продвинутым алгоритмам', 'https://example.com/img3.png', 'advanced', 200, true);

INSERT INTO levels (course_id, title, description, order_index, xp_reward)
VALUES
(1, 'Введение в алгоритмы', 'Базовые понятия и зачем всё это нужно', 1, 10),
(1, 'Сложность алгоритмов', 'Что такое Big O и как её считать', 2, 15),
(1, 'Поиск и сортировка', 'Простейшие алгоритмы поиска и сортировки', 3, 20);

INSERT INTO levels (course_id, title, description, order_index, xp_reward)
VALUES
(2, 'Массивы и списки', 'Изучаем динамические и статические структуры', 1, 15),
(2, 'Стеки и очереди', 'ФИФО и ЛИФО в жизни и коде', 2, 20),
(2, 'Деревья и графы', 'Первые шаги в сложные структуры', 3, 25);